package mathutils

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
        return FloorDivInt(-a, b) + 1
    }

    return a / b
}

func RemEuclidInt(a int, b int) int {
    return ((a % b) + b) % b;
}
