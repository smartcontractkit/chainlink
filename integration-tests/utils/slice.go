package utils

func DivideSlice[T any](slice []T, parts int) [][]T {
	var divided [][]T
	if parts == 1 {
		return [][]T{slice}
	}

	sliceLength := len(slice)
	baseSize := sliceLength / parts
	remainder := sliceLength % parts

	start := 0
	for i := 0; i < parts; i++ {
		end := start + baseSize
		if i < remainder { // Distribute the remainder among the first slices
			end++
		}

		divided = append(divided, slice[start:end])
		start = end
	}

	return divided
}
