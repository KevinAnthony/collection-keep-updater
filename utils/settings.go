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

func GetPtr[T any](data map[string]interface{}, key string) (out *T) {
	value, found := data[key]
	if !found {
		return out
	}

	cast, ok := value.(T)
	if !ok {
		return out
	}

	return &cast
}

func GetArray[T any](data map[string]interface{}, key string) (out []T) {
	value, found := data[key]
	if !found {
		return out
	}

	castArray, ok := value.([]interface{})
	if !ok {
		return out
	}

	tArray := make([]T, 0, len(castArray))
	for _, v := range castArray {
		cast, ok := v.(T)
		if !ok {
			return out
		}

		tArray = append(tArray, cast)
	}

	return tArray
}
