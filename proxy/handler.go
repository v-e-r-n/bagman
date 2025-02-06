package proxy

import (
	"fmt"
	"net/http"

	"github.com/v-e-r-n/baagman"
)

type Authifier func(string) string
type HeaderModifier func(http.Header)
type RequestModifier func(bagman.Request)

type Handler interface {
	WithRequestModifier(RequestModifier) Handler
	WithBaseURL(string) Handler
	WithAuthifier(Authifier) Handler
	Finalize() Handler
	Handle(http.ResponseWriter, *http.Request, ...bagman.Logger) error
}

type defaultLogger struct{}

func (l *defaultLogger) Debug(msg string, args ...any) {}
func (l *defaultLogger) Error(msg string, args ...any) {}
func (l *defaultLogger) Info(msg string, args ...any)  {}
func (l *defaultLogger) Warn(msg string, args ...any)  {}

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
