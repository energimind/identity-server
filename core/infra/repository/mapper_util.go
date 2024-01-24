package repository

// mapSlice maps a slice of type T to a slice of type M using the provided
// mapper function.
func mapSlice[T, M any](s []T, mapper func(T) M) []M {
	m := make([]M, len(s))

	for i, v := range s {
		m[i] = mapper(v)
	}

	return m
}
