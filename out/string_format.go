package out

import (
	"fmt"
	"strings"
)

func ValueOrEmpty[T any](t *T) (out T) {
	if t == nil {
		return out
	}

	return *t
}

func IntSliceToStrOrEmpty(i []int) string {
	if len(i) == 0 {
		return ""
	}

	return strings.Trim(strings.Replace(fmt.Sprint(i), " ", ", ", -1), "[]")
}

func Partial(s string, n int) string {
	if len(s) < n {
		return s
	}

	return s[0:n] + "..."
}
