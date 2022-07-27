package test_cast

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

func main() {
	ASSERT(131585, int(8590066177))
	ASSERT(513, int16(8590066177))
	ASSERT(1, int8(8590066177))
	ASSERT(1, int64(1))
	// var x int=512; *(*byte)(&x)=1;
	// ASSERT(513, x);
	// assert(5, ({ int x=5; long y=(long)&x; *(int*)y; }));

	ASSERT(-1, int8(255))
	ASSERT(255, uint8(255))
	ASSERT(-1, int16(65535))
	ASSERT(65535, uint16(65535))
	ASSERT(-1, int(0xffffffff))
	ASSERT(-1, int32(0xffffffff))
	ASSERT(0xffffffff, uint(0xffffffff))
	ASSERT(0xffffffff, uint32(0xffffffff))

	ASSERT(1, -1 < 1)
	ASSERT(0, -1 < uint(1))
	ASSERT(254, int8(127)+int8(127))
	ASSERT(65534, int16(32767)+int16(32767))
	ASSERT(-1, -1>>1)
	ASSERT(-1, uint64(-1))
	ASSERT(2147483647, uint(-1)>>1)
	ASSERT(-50, (-100)/2)
	ASSERT(2147483598, uint(-100)/2)
	// Floating point exception???
	// ASSERT(9223372036854775758, uint64(-100)/2)
	ASSERT(0, int64(-1)/uint(100))
	ASSERT(-2, (-100)%7)
	ASSERT(2, uint(-100)%7)
	// Floating point exception???
	// ASSERT(6, uint64(-100)%9)

	ASSERT(65535, (int(uint16(65535))))
	var x uint16 = 65535
	ASSERT(65535, x)
	var x uint16 = 65535
	ASSERT(65535, int(x))

	type T1 int16
	var x T1 = 65535
	ASSERT(-1, int(x))
	type T2 uint16
	var x T2 = 65535
	ASSERT(65535, int(x))

	ASSERT(0, bool(0.0))
	ASSERT(1, bool(0.1))
	ASSERT(3, int8(3.0))
	ASSERT(1000, int16(1000.3))
	ASSERT(3, int(3.99))
	ASSERT(2000000000000000, int64(2e15))
	ASSERT(3, float32(3.5))
	ASSERT(5, float64(float32(5.5)))
	ASSERT(3, float32(3))
	ASSERT(3, float64(3))

	println("OK")
}
