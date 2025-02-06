package utils

import "fmt"

func Find[T any](data map[string]any, key string) (T, bool, error) {
	var zero T

	maybe, found := data[key]
	if !found {
		var zero T
		return zero, false, nil
	}

	value, ok := maybe.(T)
	if !ok {
		return zero, true, fmt.Errorf("key %s is not of type %T", key, value)
	}

	return value, true, nil
}
