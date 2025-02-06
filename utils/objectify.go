package utils

import "fmt"

func Objectify(blob any) (map[string]any, error) {
	data, ok := blob.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("blob is not an object")
	}

	return data, nil
}
