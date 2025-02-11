package proxy

import (
	"net/http"
	"time"

	"github.com/v-e-r-n/bagman"
)

type BodyMod func([]byte) ([]byte, error)
type HeaderMod func(http.Header) (http.Header, error)

func OriginalBody(original []byte) ([]byte, error) {
	return original, nil
}

func OriginalHeaders(original http.Header) (http.Header, error) {
	return original, nil
}

type proxyOptions struct {
	// The path prefix to strip from the request path
	pathPrefix string

	// The base URL to send the request to
	baseURL string

	// The upstream endpoint to send the request to
	upstreamEndpoint string

	// The function to modify the request body
	requestBodyMod BodyMod

	// The function to modify the request headers
	requestHeaderMod HeaderMod

	// The function to modify the response headers
	responseHeaderMod HeaderMod

	// The http client to use
	driver *http.Client

	// The logger to use
	logger bagman.Logger
}

func Proxy() *proxyOptions {
	return &proxyOptions{
		requestBodyMod:    OriginalBody,
		requestHeaderMod:  OriginalHeaders,
		responseHeaderMod: OriginalHeaders,
		driver:            &http.Client{Timeout: 5 * time.Minute},
		logger:            getDefaultLogger(),
	}
}

func (p *proxyOptions) WithPathPrefix(pathPrefix string) *proxyOptions {
	p.pathPrefix = pathPrefix
	return p
}

func (p *proxyOptions) WithBaseURL(baseURL string) *proxyOptions {
	p.baseURL = baseURL
	return p
}

func (p *proxyOptions) WithUpstreamEndpoint(upstreamEndpoint string) *proxyOptions {
	p.upstreamEndpoint = upstreamEndpoint
	return p
}

func (p *proxyOptions) WithRequestBodyMod(requestBodyMod BodyMod) *proxyOptions {
	p.requestBodyMod = requestBodyMod
	return p
}

func (p *proxyOptions) WithRequestHeaderMod(requestHeaderMod HeaderMod) *proxyOptions {
	p.requestHeaderMod = requestHeaderMod
	return p
}

func (p *proxyOptions) WithResponseHeaderMod(responseHeaderMod HeaderMod) *proxyOptions {
	p.responseHeaderMod = responseHeaderMod
	return p
}

func (p *proxyOptions) WithDriver(driver *http.Client) *proxyOptions {
	p.driver = driver
	return p
}

func (p *proxyOptions) WithLogger(logger bagman.Logger) *proxyOptions {
	p.logger = logger
	return p
}
