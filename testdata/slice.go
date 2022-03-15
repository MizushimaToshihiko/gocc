package test_slice

func assert(want int, act int, code string)
func println(format string)

func main() {
	var a01 = [2]int{1, 2}
	var s01 = a01[0:2]
	assert(1, s01[0], "s01[0]")
	assert(2, s01[1], "s01[1]")
	assert(8, Sizeof(s01), "Sizeof(x01)")

	var a02 = [6]int{1, 2, 3, 4, 5, 6}
	var s02 = a02[2:4]
	assert(3, s02[0], "s02[0]")
	assert(4, s02[1], "s02[1]")
	assert(5, s02[2], "s02[2]") // out of range
	assert(8, Sizeof(s02), "Sizeof(x02)")

	println("OK")
}
