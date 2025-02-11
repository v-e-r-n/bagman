package proxy

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/v-e-r-n/bagman"
)

type HeaderModifier func(http.Header)
type RequestModifier func(bagman.Request)

type Handler interface {
	WithRequestModifier(RequestModifier) Handler
	WithBaseURL(string) Handler
	WithUpstreamEndpoint(string) Handler
	WithAuthifier(Authifier) Handler
	WithPathPrefix(string) Handler
	Finalize() Handler
	Handle(http.ResponseWriter, *http.Request, ...bagman.Logger) error
}

func getDefaultLogger() bagman.Logger {
	loggerOnce.Do(func() {
		logger = &defaultLogger{}
	})

	return logger
}

type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, args ...any) {}
func (l *defaultLogger) Error(msg string, args ...any) {}
func (l *defaultLogger) Info(msg string, args ...any)  {}
func (l *defaultLogger) Warn(msg string, args ...any)  {}

var loggerOnce = &sync.Once{}
var logger bagman.Logger = &defaultLogger{}

func UnfinalizedHandlerError(kind string) error {
	return &unfinalizedHandlerError{kind}
}

type unfinalizedHandlerError struct {
	kind string
}

func (e *unfinalizedHandlerError) Error() string {
	return fmt.Sprintf("%s handler is not finalized", e.kind)
}
