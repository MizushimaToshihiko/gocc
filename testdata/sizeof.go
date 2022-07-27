package test_sizeof

func assert(want int, act int, code string)
func println(format string)

#include "test.h"

func main() {
	ASSERT(1, Sizeof(byte))
	ASSERT(2, Sizeof(int16))
	ASSERT(4, Sizeof(int))
	ASSERT(8, Sizeof(int64))
	type T1 struct {
		a int
		b int
	}
	ASSERT(8, Sizeof(T1))

	var x int = 0
	ASSERT(4, Sizeof(x+1))
	ASSERT(8, Sizeof(-10+int64(5)))
	ASSERT(8, Sizeof(-10 - int64(5)))
	ASSERT(8, Sizeof(-10 * int64(5)))
	ASSERT(8, Sizeof(-10 / int64(5)))
	ASSERT(8, Sizeof(int64(-10) + 5))
	ASSERT(8, Sizeof(int64(-10) - 5))
	ASSERT(8, Sizeof(int64(-10) * 5))
	ASSERT(8, Sizeof(int64(-10) / 5))
	var i byte
	ASSERT(1, Sizeof(i++))

	ASSERT(1, Sizeof(int8)<<31>>31)
	ASSERT(1, Sizeof(int8)<<63>>63)

  ASSERT(8, Sizeof(1.0+2));
  ASSERT(8, Sizeof(1.0-2));
  ASSERT(8, Sizeof(1.0*2));
  ASSERT(8, Sizeof(1.0/2));

	println("OK")
}
