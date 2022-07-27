package test_float

func assert(want int, act int, code string)
func println(format string)

#include "test.h"

func main() {
	ASSERT(35, float32(int8(35)))
	ASSERT(35, float32(int16(35)))
	ASSERT(35, float32(int(35)))
	ASSERT(35, float32(int64(35)))
	ASSERT(35, float32(uint8(35)))
	ASSERT(35, float32(uint16(35)))
	ASSERT(35, float32(uint(35)))
	ASSERT(35, float32(uint64(35)))

	ASSERT(35, float64(int8(35)))
	ASSERT(35, float64(int16(35)))
	ASSERT(35, float64(int(35)))
	ASSERT(35, float64(int64(35)))
	ASSERT(35, float64(uint8(35)))
	ASSERT(35, float64(uint16(35)))
	ASSERT(35, float64(uint(35)))
	ASSERT(35, float64(uint64(35)))

	ASSERT(35, int8(float32(35)))
	ASSERT(35, int16(float32(35)))
	ASSERT(35, int(float32(35)))
	ASSERT(35, int64(float32(35)))
	ASSERT(35, uint8(float32(35)))
	ASSERT(35, uint16(float32(35)))
	ASSERT(35, uint(float32(35)))
	ASSERT(35, uint64(float32(35)))

	ASSERT(35, int8(float64(35)))
	ASSERT(35, int16(float64(35)))
	ASSERT(35, int(float64(35)))
	ASSERT(35, int64(float64(35)))
	ASSERT(35, uint8(float64(35)))
	ASSERT(35, uint16(float64(35)))
	ASSERT(35, uint(float64(35)))
	ASSERT(35, uint64(float64(35)))

	ASSERT(-2147483648, float64(uint64(int64(-1))))

	ASSERT(1, 2e3 == 2e3)
	ASSERT(0, 2e3 == 2e5)
	ASSERT(1, 2.0 == 2)
	ASSERT(0, 5.1 < 5)
	ASSERT(0, 5.0 < 5)
	ASSERT(1, 4.9 < 5)
	ASSERT(0, 5.1 <= 5)
	ASSERT(1, 5.0 <= 5)
	ASSERT(1, 4.9 <= 5)

	ASSERT(6, 2.3+3.8)
	ASSERT(-1, 2.3-3.8)
	ASSERT(-3, -3.8)
	ASSERT(13, 3.3*4)
	ASSERT(2, 5.0/2)

	ASSERT(0, 0.0/0.0 == 0.0/0.0)
	ASSERT(1, 0.0/0.0 != 0.0/0.0)

	ASSERT(0, 0.0/0.0 < 0)
	ASSERT(0, 0.0/0.0 <= 0)
	ASSERT(0, 0.0/0.0 > 0)
	ASSERT(0, 0.0/0.0 >= 0)

	ASSERT(0, !3.)
	ASSERT(1, !0.)

	println("OK")
}
