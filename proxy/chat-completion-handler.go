package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/v-e-r-n/bagman"
	"github.com/v-e-r-n/bagman/chat"
)

const chatCompletionBaseURL = "https://api.openai.com/v1"
const chatCompletionPath = "chat/completions"

type chatCompletionHandler struct {
	authifier       Authifier
	requestModifier RequestModifier
	baseURL         string
	finalized       bool
	driver          *http.Client
}

func ChatCompletionHandler() Handler {
	return &chatCompletionHandler{
		// By default, don't modify the authentication
		authifier: func(key string) string {
			return key
		},
		// By default, don't modify the request
		requestModifier: func(request bagman.Request) {},

		baseURL: chatCompletionBaseURL,

		driver: &http.Client{Timeout: 5 * time.Minute},
	}
}

func (h *chatCompletionHandler) WithAuthifier(authifier Authifier) Handler {
	if h.finalized {
		return h
	}

	h.authifier = authifier

	return h
}

func (h *chatCompletionHandler) WithRequestModifier(modifier RequestModifier) Handler {
	if h.finalized {
		return h
	}

	h.requestModifier = modifier

	return h
}

func (h *chatCompletionHandler) WithBaseURL(baseURL string) Handler {
	if h.finalized {
		return h
	}

	h.baseURL = baseURL

	return h
}

func (h *chatCompletionHandler) Finalize() Handler {
	h.finalized = true

	return h
}

func (h *chatCompletionHandler) Handle(w http.ResponseWriter, r *http.Request, loggers ...bagman.Logger) error {
	hlog := logger
	if len(loggers) > 0 {
		hlog = loggers[0]
	}

	if !h.finalized {
		hlog.Error("fatal exception", "error", "handler not finalized")
		return UnfinalizedHandlerError("chat completion")
	}

	upstream, err := url.Parse(h.baseURL)
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

	upstreamRequest := &chat.CompletionRequest{}
	err = json.Unmarshal(data, upstreamRequest)
	if err != nil {
		hlog.Error("failed to unmarshal request", "error", err)
		return err
	}

	h.requestModifier(upstreamRequest)

	proxyBody, err := json.Marshal(upstreamRequest)
	if err != nil {
		hlog.Error("failed to marshal upstream request", "error", err)
		return err
	}

	bodyReader := bytes.NewReader(proxyBody)
	proxyRequest, err := http.NewRequest(
		r.Method,
		upstream.JoinPath(chatCompletionPath).String(),
		bodyReader,
	)
	if err != nil {
		hlog.Error("failed to create proxy request", "error", err)
		return err
	}

	// copy the headers from the original request
	for key, values := range r.Header {
		for _, value := range values {

			proxyRequest.Header.Add(key, value)

		}
	}

	// Set the authorization header
	originalKey := strings.TrimPrefix("Bearer ", r.Header.Get("Authorization"))
	proxyRequest.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", h.authifier(originalKey)),
	)

	response, err := h.driver.Do(proxyRequest)
	if err != nil {
		hlog.Error("failed to proxy request", "error", err)
		return err
	}

	// Copy the response headers to our response
	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if upstreamRequest.IsStreaming() {
		return h.streamResponse(w, response)
	}

	hlog.Info("proxying to upstream LLM")
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		hlog.Error("failed to read response body", "error", err)
		return err
	}

	w.WriteHeader(response.StatusCode)
	_, err = w.Write(responseBody)

	return err
}

func (h *chatCompletionHandler) streamResponse(w http.ResponseWriter, response *http.Response, loggers ...bagman.Logger) error {
	hlog := logger
	if len(loggers) > 0 {
		hlog = loggers[0]
	}

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
