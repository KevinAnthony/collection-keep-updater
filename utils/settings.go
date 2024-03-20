package utils

func Get[T any](data map[string]interface{}, key string) (out T) {
	value, found := data[key]
	if !found {
		return out
	}

	cast, ok := value.(T)
	if !ok {
		return out
	}

	return cast
}

func GetPtr[T any](data map[string]interface{}, key string) *T {
	value, found := data[key]
	if !found {
		return nil
	}

	cast, ok := value.(T)
	if !ok {
		return nil
	}

	return &cast
}

func GetArray[T any](data map[string]interface{}, key string) []T {
	value, found := data[key]
	if !found {
		return nil
	}

	castArray, ok := value.([]interface{})
	if !ok {
		return nil
	}

	tArray := make([]T, 0, len(castArray))
	for _, v := range castArray {
		cast, ok := v.(T)
		if !ok {
			return nil
		}

		tArray = append(tArray, cast)
	}

	return tArray
}
