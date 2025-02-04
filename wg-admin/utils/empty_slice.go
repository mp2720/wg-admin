package utils

// Required for REST API with JSON, where the empty array is preferred over null.
func NilToEmptySlice[T any](slice []T) []T {
	if slice == nil {
		return make([]T, 0)
	} else {
		return slice
	}
}
