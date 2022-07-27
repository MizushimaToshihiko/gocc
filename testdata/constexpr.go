package test_constexpr

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

var g40 float32 = 1.5
var g41 float64 = 1 + 1*5.0/2*float64(2)*int(2.0)

func main() {
	var i int = 0
	switch 3 {
	case 5 - 2 + 0*3:
		i++
	}
	ASSERT(1, i)
	var x1 [1 + 1]int
	ASSERT(8, Sizeof(x1))
	var x2 [8 - 2]byte
	ASSERT(6, Sizeof(x2))
	var x3 [2 * 3]byte
	ASSERT(6, Sizeof(x3))
	var x4 [12 / 4]byte
	ASSERT(3, Sizeof(x4))
	var x5 [12 % 10]byte
	ASSERT(2, Sizeof(x5))
	var x6 [0b110 & 0b101]byte
	ASSERT(0b100, Sizeof(x6))
	var x7 [0b110 | 0b101]byte
	ASSERT(0b111, Sizeof(x7))
	var x8 [0b111 ^ 0b001]byte
	ASSERT(0b110, Sizeof(x8))

	var x9 [1 << 2]byte
	ASSERT(4, Sizeof(x9))
	var x10 [4 >> 1]byte
	ASSERT(2, Sizeof(x10))
	var x11 [(1 == 1) + 1]byte
	ASSERT(2, Sizeof(x11))
	var x12 [(1 != 1) + 1]byte
	ASSERT(1, Sizeof(x12))
	var x13 [(1 < 1) + 1]byte
	ASSERT(1, Sizeof(x13))
	var x14 [(1 <= 1) + 1]byte
	ASSERT(2, Sizeof(x11))
	var x15 [!0 + 1]byte
	ASSERT(2, Sizeof(x15))
	var x16 [!1 + 1]byte
	ASSERT(1, Sizeof(x16))
	var x17 [^-3]byte
	ASSERT(2, Sizeof(x17))
	var x18 [(5 || 6) + 1]byte
	ASSERT(2, Sizeof(x18))
	var x19 [(0 || 0) + 1]byte
	ASSERT(1, Sizeof(x19))
	var x20 [(1 && 1) + 1]byte
	ASSERT(2, Sizeof(x20))
	var x21 [(1 && 0) + 1]byte
	ASSERT(1, Sizeof(x21))
	var x22 [int(3)]byte
	ASSERT(3, Sizeof(x22))

	var x23 [(1,3)]byte
	ASSERT(3, Sizeof(x23))
	var x24 [byte(0xffffff0f)]byte
	ASSERT(15, Sizeof(x24))
	var x25 [int16(0xffff010f)]byte
	ASSERT(0x10f, Sizeof(x25))

	// error occures
	// var x26 [int(0xfffffffffff)+5]byte
	// ASSERT(4, Sizeof(x26));

	// Below is not supported in Go.
	// var x26 [(*int)(0) + 2]byte
	// ASSERT(8, Sizeof(x26))
	// assert(12, ({ char x[(int*)16-1]; Sizeof(x); }));
	// assert(3, ({ char x[(int*)16-(int*)4]; Sizeof(x); }));

	var x26 [(-1>>31)+5]int8
	ASSERT(4, Sizeof(x26));
	var x27 [uint8(0xffffffff)]int8
	ASSERT(255, Sizeof(x27));
	var x28 [uint16(0xffff800f)]int8
	ASSERT(0x800f, Sizeof(x28));
	var x29 [uint(0xfffffffffff)>>31]int8
	ASSERT(1, Sizeof(x29));
	var x30 [int64(-1)/(int64(1)<<62)+1]int8
	ASSERT(1, Sizeof(x30));
	// var x31 [uint64(-1)/(int64(1)<<62)+1]int8
	// ASSERT(4, Sizeof(x31));
	var x32 [uint(1)<-1]int8
	ASSERT(1, Sizeof(x32));
	var x33 [uint(1)<=-1]int8
	ASSERT(1, Sizeof(x33));

	ASSERT(1, g40 == 1.5)
	ASSERT(1, g41 == 11)

	println("OK")
}
