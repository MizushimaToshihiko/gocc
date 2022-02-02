package test

func assert(want int, act int, code string)
func println(format string)

var g3 byte
var g4 int16
var g5 int
var g6 int64
var g7 = "abcdef"
var g8 = [2]struct {
	a byte
	b byte
}{{1, 2}}
var g9 = [2]struct {
	a byte
	b int64
}{{1, 2}}

func main() {
	assert(1, Alignof(g3), "Alignof(g3)")
	assert(2, Alignof(g4), "Alignof(g4)")
	assert(4, Alignof(g5), "Alignof(g5)")
	assert(8, Alignof(g6), "Alignof(g6)")
	assert(8, Alignof(g7), "Alignof(g7)")
	assert(1, Alignof(g8), "Alignof(g8)")
	assert(8, Alignof(g9), "Alignof(g9)")

	println("OK")
}
