package pkg

func Map[T1 any, T2 any](arr *[]T1, f func(*T1) T2) *[]T2 {
	var result []T2
	for _, v := range *arr {
		result = append(result, f(&v))
	}
	return &result
}

func MapA[T1 any, T2 any](arr *[]T1, f func(*T1) T2) *[]T2 {
	var result []T2
	for _, v := range *arr {
		result = append(result, f(&v))
	}
	return &result
}

func ConvertT[T1 any, T2 any](arr *T1, f func(*T1) *T2) *T2 {
	return f(arr)
}

func Reverse[T any](arr []T) []T {
	n := len(arr)
	result := make([]T, n)
	for i, v := range arr {
		result[n-i-1] = v
	}
	return result
}

func Reduce[T any, U any](arr *[]T, f func(U, T) U, init U) U {
	result := init
	for _, v := range *arr {
		result = f(result, v)
	}
	return result
}

func IsInclude[T any](arr *[]T, f func(t1 *T) bool) bool {
	for _, v := range *arr {
		if f(&v) {
			return true
		}
	}
	return false
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

func FirstIndexOf[T1 any](in *[]T1, f func(*T1) bool) int {
	for i, v := range *in {
		if f(&v) {
			return i
		}
	}
	return -1
}

func LastIndexOf[T1 any](in *[]T1, f func(*T1) bool) int {
	length := len(*in)
	for i := length - 1; i >= 0; i-- {
		if f(&(*in)[i]) {
			return i
		}
	}

	return -1
}

func RemoveElementAt[T1 any](in []T1, index int) []T1 {
	s := in
	i := index
	copy(s[i:], s[i+1:])
	s = s[:len(s)-1]
	return s
}

func RemoveElementsIn[T1 any](source []T1, toBeRemove []T1, contains func(e1 *T1, e2 *T1) bool) []T1 {
	var result []T1 = make([]T1, 0)
	for _, v := range source {
		if !ContainsFuncV2(toBeRemove, &v, contains) {
			result = append(result, v)
		}
	}
	return result
}

func RemoveElements[T1 any](source []T1, contains func(e1 *T1) bool) []T1 {
	var result []T1 = make([]T1, 0)
	for _, v := range source {
		if !contains(&v) {
			result = append(result, v)
		}
	}
	return result
}

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsFunc[T any](slice []T, f func(e *T) bool) bool {
	for _, v := range slice {
		if f(&v) {
			return true
		}
	}
	return false
}

func ContainsFuncV2[T any](slice []T, e2 *T, f func(e1, e2 *T) bool) bool {
	for _, v := range slice {
		if f(&v, e2) {
			return true
		}
	}
	return false
}

func HasRelationshipFunc[T1 any, T2 any](slice []T1, e2 *T2, relationship func(e1 *T1, e2 *T2) bool) bool {
	for _, v := range slice {
		if relationship(&v, e2) {
			return true
		}
	}
	return false
}
