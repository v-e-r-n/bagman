package utils

import "fmt"

func Get[T any](data map[string]any, key string) (T, error) {
	maybe, found := data[key]
	if !found {
		var zero T
		return zero, fmt.Errorf("key %s not found", key)
	}

	value, ok := maybe.(T)
	if !ok {
		return value, fmt.Errorf("key %s is not of type %T", key, value)
	}

	return value, nil
}

func GetMap(data map[string]any, key string) (map[string]any, error) {
	maybe, found := data[key]
	if !found {
		return nil, fmt.Errorf("key %s not found", key)
	}

	value, ok := maybe.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("key %s is not an object", key)
	}

	return value, nil
}
