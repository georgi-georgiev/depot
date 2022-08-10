package main

func AbsInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func absDiffUint(x, y uint) uint {
	if x < y {
		return y - x
	}
	return x - y
}
