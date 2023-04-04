package maps

func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) M1 {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
