package utils

func Unique[T comparable](ss []T) []T {
	seen := make(map[T]struct{}, len(ss))
	result := make([]T, 0, max(1, len(ss)))
	for _, v := range ss {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}

func MergeSliceUnique[T comparable](ss ...[]T) []T {
	size := 0
	for _, s := range ss {
		size += len(s)
	}

	seen := make(map[T]struct{}, size)
	result := make([]T, 0, size)
	for _, s := range ss {
		for _, v := range s {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
