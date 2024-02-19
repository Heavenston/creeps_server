package mathutils

import (
	"cmp"
)

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func FloorDivInt(a int, b int) int {
	if b <= 0 {
		panic("division by non-positive number")
	}

	if a < 0 {
		return -((-a) / b) - 1
	}

	return a / b
}

func RemEuclidInt(a int, b int) int {
	return ((a % b) + b) % b
}

func Min[T cmp.Ordered](a T, rest ...T) T {
	min := a
	for _, x := range rest {
		if x < min {
			min = x
		}
	}
	return min
}

func Max[T cmp.Ordered](a T, rest ...T) T {
	max := a
	for _, x := range rest {
		if x > max {
			max = x
		}
	}
	return max
}

