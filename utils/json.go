package utils

import "encoding/json"

func Unmarshal[T any](raw *json.RawMessage) (out T, err error) {
	if raw == nil {
		return out, nil
	}

	if err := json.Unmarshal(*raw, &out); err != nil {
		return out, err
	}

	return out, nil
}
