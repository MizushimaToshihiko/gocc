package test_alignof

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

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
	ASSERT(1, Alignof(g3))
	ASSERT(2, Alignof(g4))
	ASSERT(4, Alignof(g5))
	ASSERT(8, Alignof(g6))
	ASSERT(8, Alignof(g7))
	ASSERT(1, Alignof(g8))
	ASSERT(8, Alignof(g9))

	var x int8
	ASSERT(1, Alignof(x)<<31>>31)
	ASSERT(1, Alignof(x)<<63>>63)

	println("OK")
}
