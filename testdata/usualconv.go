package test_usualconv

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

func ret10() int {
	return 10
}

func main() {
	ASSERT(int64(-5), -10+int64(5))
	ASSERT(int64(-15), -10-int64(5))
	ASSERT(int64(-50), -10*int64(5))
	ASSERT(int64(-2), -10/int64(5))

	ASSERT(1, -2 < int64(-1))
	ASSERT(1, -2 <= int64(-1))
	ASSERT(0, -2 > int64(-1))
	ASSERT(0, -2 >= int64(-1))

	ASSERT(1, int64(-2) < -1)
	ASSERT(1, int64(-2) <= -1)
	ASSERT(0, int64(-2) > -1)
	ASSERT(0, int64(-2) >= -1)

	ASSERT(0, 2147483647+2147483647+2)
	var x1 int64
	x1 = -1
	ASSERT(int64(-1), x1)

	var x2 [3]byte
	x2[0] = 0
	x2[1] = 1
	x2[2] = 2
	var y2 *byte = x2 + 1
	ASSERT(1, y2[0])
	// ASSERT(0, y2[-1], "var x2 [3]byte; x2[0]=0; x2[1]=1; x2[2]=2;var y2 *byte=x2+1; y2[-1]")
	type t3 struct{ a byte }
	var x3 t3
	var y3 t3
	x3.a = 5
	y3 = x3
	ASSERT(5, y3.a)

	var fn = ret10
	ASSERT(10, fn())

	println("OK")
}
