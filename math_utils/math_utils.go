package mathutils

import (
	"cmp"
	"fmt"
)

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func FloorDivInt(a int, b int) int {
	if b <= 0 {
		panic(nil)
	}

	if a < 0 {
		if -a < 0 {
			panic(fmt.Sprintf("number too big: %d", -a))
		}
		return -FloorDivInt(-a, b) + 1
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

