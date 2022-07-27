package test_decl

func assert(want int, act int, code string)
func println(format string)

func strcmp(s1 string, s2 string) int

#include "test.h"

var g1, g2, g3 bool
var g4, g5, g6 = 2.0, 8, "foo"

var (
	g7          int
	g8, g9, g10 = 2.0, 3.0, "bar"
)

func main() {
	var x1 byte
	ASSERT(1, Sizeof(x1))
	var x2 int16
	ASSERT(2, Sizeof(x2))
	var x3 int
	ASSERT(4, Sizeof(x3))
	var x4 int64
	ASSERT(8, Sizeof(x4))

	var x5 bool = 0
	ASSERT(0, x5)
	var x6 bool = 1
	ASSERT(1, x6)
	var x7 bool = 2
	ASSERT(1, x7)
	ASSERT(1, bool(1))
	ASSERT(1, bool(2))
	ASSERT(0, bool(byte(256)))

	var x8, x9 int
	ASSERT(0, x8)
	ASSERT(4, Sizeof(x8))
	ASSERT(0, x9)

	var x10, x11, x12, x13 int = 1, 2, 3, 4
	ASSERT(1, x10)
	ASSERT(4, Sizeof(x10))
	ASSERT(2, x11)
	ASSERT(3, x12)
	ASSERT(4, x13)

	var x14, x15, x16, x17 string = "1", "2", "3", "4"
	ASSERT(8, Sizeof(x14))
	ASSERT(0, strcmp(x14, "1"))
	ASSERT(0, strcmp(x15, "2"))
	ASSERT(0, strcmp(x16, "3"))
	ASSERT(0, strcmp(x17, "4"))

	x18, x19, x20 := 1, 2, 3
	ASSERT(1, x18)
	ASSERT(2, x19)
	ASSERT(3, x20)

	x21, x22, x23, x24 := "1", "2", "3", "4"
	ASSERT(8, Sizeof(x21))
	ASSERT(0, strcmp(x21, "1"))
	ASSERT(0, strcmp(x22, "2"))
	ASSERT(0, strcmp(x23, "3"))
	ASSERT(0, strcmp(x24, "4"))

	var (
		i25           int
		u25, v25, s25 = 2.0, 3.0, "bar"
	)
	ASSERT(0, i25)
	ASSERT(2.0, u25)
	ASSERT(3.0, v25)
	ASSERT(0, strcmp(s25, "bar"))

	ASSERT(0, g1)
	ASSERT(0, g2)
	ASSERT(0, g3)
	ASSERT(2.0, g4)
	ASSERT(8, g5)
	ASSERT(0, strcmp(g6, "foo"))

	ASSERT(0, g7)
	ASSERT(2.0, g8)
	ASSERT(3.0, g9)
	ASSERT(0, strcmp(g10, "bar"))

	println("OK")
}
