package test_slice

func assert(want int, act int, code string)
func println(format ...string)

var g01 = 1
var g02 = 3

func ret3() int {
	return 3
}

func main() {
	// var a01 = [2]int{1, 2}
	// var s01 = a01[0:2]
	// assert(1, s01[0], "s01[0]")
	// assert(2, s01[1], "s01[1]")
	// assert(8, Sizeof(s01), "Sizeof(x01)")
	// assert(2, len(s01), "len(s01)")
	// assert(2, cap(s01), "cap(s01)")

	var a02 = [6]int{1, 2, 3, 4, 5, 6}
	var s02 = a02[2:5]
	// assert(3, s02[0], "s02[0]")
	// assert(4, s02[1], "s02[1]")
	// assert(5, s02[2], "s02[2]")
	// assert(8, Sizeof(s02), "Sizeof(x02)")
	// assert(3, len(s02), "len(s02)")
	// assert(4, cap(s02), "cap(s02)")
	s02[0] = 100
	// assert(100, a02[2], "a02[2]")

	// var x021, x0211 int
	// x021, x0211 = 2.0, 5.0
	// s021 := a02[x021:x0211]
	// assert(100, s021[0], "s021[0]")
	// assert(4, s021[1], "s021[1]")
	// assert(5, s021[2], "s021[2]")
	// assert(8, Sizeof(s021), "Sizeof(x021)")
	// assert(3, len(s021), "len(s021)")
	// assert(4, cap(s021), "cap(s021)")

	// var x022, x0221 = 2.0, 5.0
	// s022 := a02[x022:x0221]
	// assert(100, s022[0], "s022[0]")
	// assert(4, s022[1], "s022[1]")
	// assert(5, s022[2], "s022[2]")
	// assert(8, Sizeof(s022), "Sizeof(x022)")
	// assert(3, len(s022), "len(s022)")
	// assert(4, cap(s022), "cap(s022)")

	s023 := a02[g01:ret3()]
	assert(2, s023[0], "s023[0]")
	assert(100, s023[1], "s023[1]")
	assert(2, len(s023), "len(s023)")
	assert(5, cap(s023), "cap(s023)")

	println("OK")
}
