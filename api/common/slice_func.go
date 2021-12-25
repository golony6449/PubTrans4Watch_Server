package common

func insert[T](origin []T, idx int, target T) []T {
	result := origin[:idx]

	for _, elem := range origin[idx:]{
		result = append(result, elem)
	}

	return result
}