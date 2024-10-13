package util

func PointTo[T any](value T) *T {
	return &value
}
