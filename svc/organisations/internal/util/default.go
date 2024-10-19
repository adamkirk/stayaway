package util

func Default[T any](value *T, def *T) *T {
	if value != nil {
		return value
	}

	return def
}
