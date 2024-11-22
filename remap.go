package csv2structs

// remap any type of slice to another type of slice via the provided function
func remap[T any, K any](objs []T, fn func(T) K) []K {
	var result []K
	for _, obj := range objs {
		result = append(result, fn(obj))
	}
	return result
}
