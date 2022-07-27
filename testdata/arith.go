package test_arith

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

func main() {
	ASSERT(0, 0)
	ASSERT(42, 42)
	ASSERT(5, 5)
	ASSERT(41, 12+34-5)
	ASSERT(15, 5*(9-6))
	ASSERT(4, (3+5)/2)
	ASSERT(10, -10+20)
	ASSERT(10, - -10)
	ASSERT(10, - -+10)

	ASSERT(0, 0 == 1)
	ASSERT(1, 42 == 42)
	ASSERT(1, 0 != 1)
	ASSERT(0, 42 != 42)

	ASSERT(1, 0 < 1)
	ASSERT(0, 1 < 1)
	ASSERT(0, 2 < 1)
	ASSERT(1, 0 <= 1)
	ASSERT(1, 1 <= 1)
	ASSERT(0, 2 <= 1)

	ASSERT(1, 1 > 0)
	ASSERT(0, 1 > 1)
	ASSERT(0, 1 > 2)
	ASSERT(1, 1 >= 0)
	ASSERT(1, 1 >= 1)
	ASSERT(0, 1 >= 2)

	ASSERT(4294967297, 4294967297)
	ASSERT(0, 1073741824*100/100)

	var i int
	i = 2
	i += 5
	ASSERT(7, i)
	i = 5
	i -= 2
	ASSERT(3, i)
	i = 3
	i *= 2
	ASSERT(6, i)
	i = 6
	i /= 2
	ASSERT(3, i)
	i = 2
	i++
	ASSERT(3, i)
	i = 2
	i--
	ASSERT(1, i)

	ASSERT(0, !1)
	ASSERT(0, !2)
	ASSERT(1, !0)
	ASSERT(1, !byte(0))
	ASSERT(3, int64(3))
	ASSERT(4, Sizeof(!byte(0)))
	ASSERT(4, Sizeof(!int64(0)))

	ASSERT(-1, ^0)
	ASSERT(0, ^-1)

	ASSERT(5, 17%6)
	ASSERT(5, (int64(17))%6)
	i = 10
	i %= 4
	ASSERT(2, i)
	var i int64
	i = 10
	i %= 4
	ASSERT(2, i)

	ASSERT(0, 0&1)
	ASSERT(1, 3&1)
	ASSERT(3, 7&3)
	ASSERT(10, -1&10)

	ASSERT(1, 0|1)
	ASSERT(0b10011, 0b10000|0b00011)

	ASSERT(0, 0^0)
	ASSERT(0, 0b1111^0b1111)
	ASSERT(0b110100, 0b111000^0b001100)

	var x int
	var p *int = &x
	println("var x int;var p *int=&x")
	ASSERT(20, p+20-p)
	ASSERT(1, p+20-p > 0)
	ASSERT(-20, p-20-p)
	ASSERT(1, p-20-p < 0)

	var x01, x02 int
	x01, x02 = 1, 2
	ASSERT(1, x01)
	ASSERT(2, x02)

	println("OK")
}
