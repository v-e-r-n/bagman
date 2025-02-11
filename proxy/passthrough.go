package proxy

import (
	"net/http"

	"github.com/v-e-r-n/bagman"
)

const passthroughPath = "v1/chat/completions"

type passthroughHandler struct {
	finalized bool
	opts      *proxyOptions
}

func PassthroughHandler() Handler {
	return &passthroughHandler{
		opts: Proxy().WithBaseURL(
			defaultUpstreamURL,
		),
	}
}

func (h *passthroughHandler) WithUpstreamEndpoint(endpoint string) Handler {
	if h.finalized {
		return h
	}

	h.opts.WithUpstreamEndpoint(endpoint)

	return h
}

func (h *passthroughHandler) WithPathPrefix(prefix string) Handler {
	if h.finalized {
		return h
	}

	h.opts.WithPathPrefix(prefix)

	return h
}

func (h *passthroughHandler) WithAuthifier(authifier Authifier) Handler {
	if h.finalized {
		return h
	}

	h.opts.WithRequestHeaderMod(authify(authifier))

	return h
}

func (h *passthroughHandler) WithRequestModifier(modifier RequestModifier) Handler {
	if h.finalized {
		return h
	}

	hlog := h.opts.logger

	h.opts.WithRequestBodyMod(requestify(NewGenericRequest(), modifier, hlog))

	// h.opts.WithRequestBodyMod(func(data []byte) ([]byte, error) {
	// 	upstreamRequest := &chat.CompletionRequest{}
	// 	err := json.Unmarshal(data, upstreamRequest)
	// 	if err != nil {
	// 		hlog.Error("failed to unmarshal request", "error", err)
	// 		return nil, err
	// 	}
	//
	// 	modifier(upstreamRequest)
	//
	// 	proxyBody, err := json.Marshal(upstreamRequest)
	// 	if err != nil {
	// 		hlog.Error("failed to marshal upstream request", "error", err)
	// 		return nil, err
	// 	}
	//
	// 	return proxyBody, nil
	// })

	return h
}

func (h *passthroughHandler) WithBaseURL(baseURL string) Handler {
	if h.finalized {
		return h
	}

	h.opts.WithBaseURL(baseURL)

	return h
}

func (h *passthroughHandler) Finalize() Handler {
	h.finalized = true

	return h
}

func (h *passthroughHandler) Handle(w http.ResponseWriter, r *http.Request, loggers ...bagman.Logger) error {
	hlog := logger
	if len(loggers) > 0 {
		hlog = loggers[0]
	}

	if !h.finalized {
		hlog.Error("fatal exception", "error", "handler not finalized")
		return UnfinalizedHandlerError("chat completion")
	}

	h.opts.logger = hlog

	return Handle(w, r, h.opts)
}
