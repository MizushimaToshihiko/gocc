package test_slice

func assert(want int, act int, code string)
func println(format ...string)

var g01 = 1
var g02 = 3

func strcmp(s1, s2 string) int

func ret3() int {
	return 3
}

func retf3() float64 {
	return 3.5
}

func main() {
	var a01 = [2]int{1, 2}
	var s01 = a01[0:2]
	assert(1, s01[0], "s01[0]")
	assert(2, s01[1], "s01[1]")
	assert(8, Sizeof(s01), "Sizeof(x01)")
	assert(2, len(s01), "len(s01)")
	assert(2, cap(s01), "cap(s01)")

	var a02 = [6]int{1, 2, 3, 4, 5, 6}
	var s02 = a02[2:5]
	assert(3, s02[0], "s02[0]")
	assert(4, s02[1], "s02[1]")
	assert(5, s02[2], "s02[2]")
	assert(8, Sizeof(s02), "Sizeof(x02)")
	assert(3, len(s02), "len(s02)")
	assert(4, cap(s02), "cap(s02)")
	s02[0] = 100
	assert(100, a02[2], "a02[2]")

	var x021, x0211 int
	x021, x0211 = 2.0, 5.0
	s021 := a02[x021:x0211]
	assert(100, s021[0], "s021[0]")
	assert(4, s021[1], "s021[1]")
	assert(5, s021[2], "s021[2]")
	assert(8, Sizeof(s021), "Sizeof(x021)")
	assert(3, len(s021), "len(s021)")
	assert(4, cap(s021), "cap(s021)")

	var x022, x0221 = 2.0, 5.0
	s022 := a02[x022:x0221]
	assert(100, s022[0], "s022[0]")
	assert(4, s022[1], "s022[1]")
	assert(5, s022[2], "s022[2]")
	assert(8, Sizeof(s022), "Sizeof(x022)")
	assert(3, len(s022), "len(s022)")
	assert(4, cap(s022), "cap(s022)")

	s023 := a02[g01 : ret3()+3]
	assert(2, s023[0], "s023[0]")
	assert(100, s023[1], "s023[1]")
	assert(5, len(s023), "len(s023)")
	assert(5, cap(s023), "cap(s023)")

	s024 := a02[g01 : retf3()+3]
	assert(2, s024[0], "s024[0]")
	assert(100, s024[1], "s024[1]")
	assert(5, len(s024), "len(s024)")
	assert(5, cap(s024), "cap(s024)")

	var s025 []int
	assert(0, len(s025), "len(s025)")
	assert(0, cap(s025), "cap(s025)")
	assert(8, Sizeof(s025), "Sizeof(s025)")

	var s026 []int = []int{1, 2, 3, 4, 5, 6}
	assert(0, len(s026), "len(s026)")
	assert(0, cap(s026), "cap(s026)")
	assert(8, Sizeof(s026), "Sizeof(s026)")
	assert(1, s026[0], "s026[0]")
	assert(2, s026[1], "s026[1]")
	assert(3, s026[2], "s026[2]")
	assert(4, s026[3], "s026[3]")
	assert(5, s026[4], "s026[4]")
	assert(6, s026[5], "s026[5]")
	s026[0], s026[1], s026[2], s026[3], s026[4], s026[5] = 100, 101, 102, 103, 104, 105
	assert(100, s026[0], "s026[0]")
	assert(101, s026[1], "s026[1]")
	assert(102, s026[2], "s026[2]")
	assert(103, s026[3], "s026[3]")
	assert(104, s026[4], "s026[4]")
	assert(105, s026[5], "s026[5]")

	var s027 = []int{1, 2, 3, 4, 5, 6}
	assert(0, len(s027), "len(s027)")
	assert(0, cap(s027), "cap(s027)")
	assert(8, Sizeof(s027), "Sizeof(s027)")
	assert(1, s027[0], "s027[0]")
	assert(2, s027[1], "s027[1]")
	assert(3, s027[2], "s027[2]")
	assert(4, s027[3], "s027[3]")
	assert(5, s027[4], "s027[4]")
	assert(6, s027[5], "s027[5]")

	s028 := []string{"abc", "def", "ghi"}
	assert(8, Sizeof(s028), "Sizeof(s028)")
	assert(0, strcmp(s028[0], "abc"), "strcmp(s028[0], \"abc\")")
	assert(0, strcmp(s028[1], "def"), "strcmp(s028[1], \"def\")")
	assert(0, strcmp(s028[2], "ghi"), "strcmp(s028[2], \"ghi\")")

	println("OK")
}
