package proxy

import (
	"fmt"
	"net/http"
)

type Authifier func(string) string

func authify(authifier Authifier) HeaderMod {
	return func(headers http.Header) (http.Header, error) {
		originalKey := headers.Get("Authorization")
		headers.Set(
			"Authorization",
			fmt.Sprintf("Bearer %s", authifier(originalKey)),
		)
		return headers, nil
	}
}
