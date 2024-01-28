package utils

func MergePointer[T comparable](a *T, b *T) *T {
	if a == nil && b == nil {
		return nil
	} else if a == nil && b != nil {
		return b
	} else if a != nil && b == nil {
		return a
	} else if *a == *b {
		return a
	}
	c := MergeValue(*a, *b)
	return &c
}

func MergeValue[T comparable](a, b T) T {
	var zero T
	if a == zero && b == zero {
		return zero
	} else if a == zero && b != zero {
		return b
	} else if a != zero && b == zero {
		return a
	}
	return b
}
