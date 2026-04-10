package utils

// RemoveElement removes the first occurrence of a specified element from a slice
func RemoveElement[T comparable](slice []T, element T) []T {
	for i, v := range slice {
		if v == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
