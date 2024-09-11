package pkg

func Map[T1 any, T2 any](arr *[]T1, f func(*T1) T2) *[]T2 {
	var result []T2
	for _, v := range *arr {
		result = append(result, f(&v))
	}
	return &result
}

func Filter[T any](arr []T, f func(T) bool) []T {
	var result []T
	for _, v := range arr {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
