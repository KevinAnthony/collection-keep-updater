package utils

func Contains[T comparable](list []T, find T, cmp func(T, T) bool) bool {
	for _, l := range list {
		if cmp(l, find) {
			return true
		}
	}

	return false
}
