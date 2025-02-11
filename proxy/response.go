package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/v-e-r-n/bagman"
)

func Handle(w http.ResponseWriter, r *http.Request, opts *proxyOptions) error {
	hlog := opts.logger
	baseURL := opts.baseURL
	pathPrefix := opts.pathPrefix
	responseHeaderMod := opts.responseHeaderMod

	var endpoint string

	if opts.upstreamEndpoint != "" {
		endpoint = canonicalize(opts.upstreamEndpoint)
	} else {

		// Extract the endpoint from the request path
		endpoint = extractCanonicalEndpoint(pathPrefix, r.URL.Path)
	}

	upstream, err := url.Parse(baseURL)
	if err != nil {
		hlog.Error("failed to parse upstream url", "error", err)
		return err
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		hlog.Error("failed to read request body", "error", err)
		return err
	}

	defer r.Body.Close()

	data, err = opts.requestBodyMod(data)
	if err != nil {
		hlog.Error("failed to modify request body", "error", err)
		return err
	}

	upstreamHeaders, err := opts.requestHeaderMod(r.Header)
	if err != nil {
		hlog.Error("failed to modify request headers", "error", err)
		return err
	}

	bodyReader := bytes.NewReader(data)

	proxyRequest, err := http.NewRequest(
		r.Method,
		upstream.JoinPath(endpoint).String(),
		bodyReader,
	)

	if err != nil {
		hlog.Error("failed to create proxy request", "error", err)
		return err
	}

	// copy the upstream headers to the proxy request
	for key, values := range upstreamHeaders {
		for _, value := range values {
			proxyRequest.Header.Add(key, value)
		}
	}

	response, err := opts.driver.Do(proxyRequest)
	if err != nil {
		hlog.Error("failed to send proxy request", "error", err)
		return err
	}

	responseHeaders, err := responseHeaderMod(response.Header)
	if err != nil {
		hlog.Error("failed to modify response headers", "error", err)
		return err
	}

	// copy the response headers to the response writer
	for key, values := range responseHeaders {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if detectStreamingResponse(response) {
		return doStreamResponse(w, response, hlog)
	}

	return doResponse(w, response, hlog)
}

func detectStreamingResponse(response *http.Response) bool {
	// Check if the content-type is text/event-stream
	contentType := response.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		return true
	}

	// Check if the content-type is application/stream+json
	if strings.Contains(contentType, "application/stream+json") {
		return true
	}

	// Check if the content-type is application/x-ndjson
	if strings.Contains(contentType, "application/x-ndjson") {
		return true
	}

	return false
}

var multislashRegex = regexp.MustCompile(`/+`)

func canonicalize(path string) string {
	return multislashRegex.ReplaceAllString("/"+path, "/")
}

func extractCanonicalEndpoint(prefix string, path string) string {
	// No prefix is basically a no-op
	if prefix == "" || prefix == "/" {
		fmt.Println("no prefix found")
		return canonicalize(path)
	}

	if strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	return canonicalize(strings.TrimPrefix(path, prefix))
}

func doResponse(w http.ResponseWriter, response *http.Response, hlog bagman.Logger) error {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		hlog.Error("failed to read response body", "error", err)
		return err
	}

	w.WriteHeader(response.StatusCode)
	_, err = w.Write(responseBody)

	return err
}

func doStreamResponse(w http.ResponseWriter, response *http.Response, hlog bagman.Logger) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		hlog.Error("response writer does not support flushing")
		return fmt.Errorf("response writer does not support flushing")
	}

	body := response.Body
	defer body.Close()

	w.WriteHeader(response.StatusCode)

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		output := scanner.Text() + "\n"

		_, err := w.Write([]byte(output))
		if err != nil {
			hlog.Error("failed to write response", "error", err)
			return err
		}
		flusher.Flush()
	}

	return nil

}
