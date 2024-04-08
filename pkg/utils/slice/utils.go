package slice

import "math"

func Filter[T any](ss []T, test func(T) bool) []T {
	res := make([]T, 0, len(ss))

	for _, s := range ss {
		if test(s) {
			res = append(res, s)
		}
	}

	return res
}

func Map[T any, V any](s []T, f func(T) V) []V {
	res := make([]V, len(s))

	for i, e := range s {
		res[i] = f(e)
	}

	return res
}

func Unique[T comparable](s []T) []T {
	visited := make(map[T]bool)
	res := make([]T, 0, len(s))

	for _, elem := range s {
		_, ok := visited[elem]
		if !ok {
			visited[elem] = true

			res = append(res, elem)
		}
	}

	return res
}

func Slice[T any](s []T, offset int, limit int) []T {
	if limit == 0 {
		return nil
	}

	if offset > len(s) {
		return nil
	}

	if limit == math.MaxInt64 {
		limit = len(s)
	}

	leftBound := offset
	rightBound := leftBound + limit

	if rightBound >= len(s) {
		rightBound = len(s)
	}

	return s[leftBound:rightBound]
}

func ContainsValue[T comparable](ss []*T, v T) bool {
	for _, s := range ss {
		if *s == v {
			return true
		}
	}

	return false
}

func PointerAndValueSlicesEquals[T comparable](p []*T, v []T) bool {
	if len(p) != len(v) {
		return false
	}

	for i := range p {
		if *p[i] != v[i] {
			return false
		}
	}

	return true
}
