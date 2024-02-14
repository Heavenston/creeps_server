package mathutils

import (
	"fmt"
	"math"
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

func MinInt(a ...int) int {
	min := math.MinInt
	for _, x := range a {
		if x < min {
			min = x
		}
	}
	return min
}

func MaxInt(a ...int) int {
	max := math.MaxInt
	for _, x := range a {
		if x > max {
			max = x
		}
	}
	return max
}

func MinFloat64(a ...float64) float64 {
	min := math.Inf(1)
	for _, x := range a {
		if x < min {
			min = x
		}
	}
	return min
}

func MaxFloat64(a ...float64) float64 {
	max := math.Inf(-1)
	for _, x := range a {
		if x > max {
			max = x
		}
	}
	return max
}
