package test_typedef

func assert(want int, act int, code string)
func println(format string)

#include "test.h"

type MyInt int
type MyInt2 [4]int

func main() {
	type t1 int
	var x1 t1 = 1
	ASSERT(1, x1)
	type t2 struct{ a int }
	var x2 t2
	x2.a = 1
	ASSERT(1, x2.a)
	type t3 int
	var t3 t3 = 1
	ASSERT(1, t3)
	type t4 struct{ a int }
	{
		type t4 int
	}
	var x4 t4
	x4.a = 2
	ASSERT(2, x4.a)
	var x5 MyInt = 3
	ASSERT(3, x5)
	var x6 MyInt2
	ASSERT(16, Sizeof(x6))

	println("OK")
}
