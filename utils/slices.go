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
