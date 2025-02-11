package proxy

import (
	"fmt"
	"testing"
)

func TestExtractCanonicalEndpoint(t *testing.T) {
	prefix := "/some/random/prefix"
	canonical := "/v1/engines/davinci/completions"

	t.Run("with a single /v1/", func(t *testing.T) {
		path := prefix + canonical

		actual := extractCanonicalEndpoint(prefix, path)

		t.Run("it returns the expected path", func(t *testing.T) {
			if actual != canonical {
				t.Errorf("expected path %s, got %s", canonical, actual)
			}
		})

	})

	t.Run("with multiple /v1/", func(t *testing.T) {
		canonical := "/v1/chat/completions/v1/engines/davinci/completions"
		path := prefix + canonical

		actual := extractCanonicalEndpoint(prefix, path)

		t.Run("it returns everything after the prefix", func(t *testing.T) {
			if actual != canonical {
				t.Errorf("expected path %s, got %s", canonical, actual)
			}
		})

	})

	t.Run("with nothing after the prefix", func(t *testing.T) {
		canonical := ""
		path := prefix + canonical

		actual := extractCanonicalEndpoint(prefix, path)
		t.Run("it returns the root endpoint", func(t *testing.T) {
			fmt.Println("fuck thaat prefix?")
			if actual != "/" {
				t.Errorf("expected path /, got %s", actual)
			}
		})

	})

	t.Run("with no prefix", func(t *testing.T) {
		prefix := ""
		path := prefix + canonical

		actual := extractCanonicalEndpoint(prefix, path)

		t.Run("it returns the expected path", func(t *testing.T) {
			if actual != canonical {
				t.Errorf("expected path %s, got %s", canonical, actual)
			}
		})

	})

	t.Run("with a root prefix", func(t *testing.T) {
		prefix := "/"
		path := prefix + canonical

		actual := extractCanonicalEndpoint(prefix, path)

		t.Run("it returns the expected path", func(t *testing.T) {
			if actual != canonical {
				t.Errorf("expected path %s, got %s", canonical, actual)
			}
		})

	})

	t.Run("with a prefix that has a trailing slash", func(t *testing.T) {
		prefix := "/some/random/prefix/"
		path := prefix + canonical

		actual := extractCanonicalEndpoint(prefix, path)

		t.Run("it returns the expected path", func(t *testing.T) {
			if actual != canonical {
				t.Errorf("expected path %s, got %s", canonical, actual)
			}
		})

	})
}
