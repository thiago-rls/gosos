package utils

// RemoveElement returns a new slice with the first occurrence of element
// removed. The input slice is not modified.
func RemoveElement[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			result := make([]T, 0, len(slice)-1)
			result = append(result, slice[:i]...)
			result = append(result, slice[i+1:]...)
			return result
		}
	}
	return slice
}
