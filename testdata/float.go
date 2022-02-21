package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	assert(35, float32(int8(35)), "float32(int8(35))")
	assert(35, float32(int16(35)), "float32(int16(35))")
	assert(35, float32(int(35)), "float32(int(35))")
	assert(35, float32(int64(35)), "float32(int64(35))")
	assert(35, float32(uint8(35)), "float32(uint8(35))")
	assert(35, float32(uint16(35)), "float32(uint16(35))")
	assert(35, float32(uint(35)), "float32(uint(35))")
	assert(35, float32(uint64(35)), "float32(uint64(35))")

	assert(35, float64(int8(35)), "float64(int8(35))")
	assert(35, float64(int16(35)), "float64(int16(35))")
	assert(35, float64(int(35)), "float64(int(35))")
	assert(35, float64(int64(35)), "float64(int64(35))")
	assert(35, float64(uint8(35)), "float64(uint8(35))")
	assert(35, float64(uint16(35)), "float64(uint16(35))")
	assert(35, float64(uint(35)), "float64(uint(35))")
	assert(35, float64(uint64(35)), "float64(uint64(35))")

	assert(35, int8(float32(35)), "int8(float32(35))")
	assert(35, int16(float32(35)), "int16(float32(35))")
	assert(35, int(float32(35)), "int(float32(35))")
	assert(35, int64(float32(35)), "int64(float32(35))")
	assert(35, uint8(float32(35)), "uint8(float32(35))")
	assert(35, uint16(float32(35)), "uint16(float32(35))")
	assert(35, uint(float32(35)), "uint(float32(35))")
	assert(35, uint64(float32(35)), "uint64(float32(35))")

	assert(35, int8(float64(35)), "int8(float64(35))")
	assert(35, int16(float64(35)), "int16(float64(35))")
	assert(35, int(float64(35)), "int(float64(35))")
	assert(35, int64(float64(35)), "int64(float64(35))")
	assert(35, uint8(float64(35)), "uint8(float64(35))")
	assert(35, uint16(float64(35)), "uint16(float64(35))")
	assert(35, uint(float64(35)), "uint(float64(35))")
	assert(35, uint64(float64(35)), "uint64(float64(35))")

	assert(-2147483648, float64(uint64(int64(-1))), "float64(uint64(int64(-1))")

	assert(1, 2e3 == 2e3, "2e3==2e3")
	assert(0, 2e3 == 2e5, "2e3==2e5")
	assert(1, 2.0 == 2, "2.0==2")
	assert(0, 5.1 < 5, "5.1<5")
	assert(0, 5.0 < 5, "5.0<5")
	assert(1, 4.9 < 5, "4.9<5")
	assert(0, 5.1 <= 5, "5.1<=5")
	assert(1, 5.0 <= 5, "5.0<=5")
	assert(1, 4.9 <= 5, "4.9<=5")

	assert(6, 2.3+3.8, "2.3+3.8")
	assert(-1, 2.3-3.8, "2.3-3.8")
	assert(-3, -3.8, "-3.8")
	assert(13, 3.3*4, "3.3*4")
	assert(2, 5.0/2, "5.0/2")

	assert(0, 0.0/0.0 == 0.0/0.0, "0.0/0.0==0.0/0.0")
	assert(1, 0.0/0.0 != 0.0/0.0, "0.0/0.0!=0.0/0.0")

	assert(0, 0.0/0.0 < 0, "0.0/0.0<0")
	assert(0, 0.0/0.0 <= 0, "0.0/0.0<=0")
	assert(0, 0.0/0.0 > 0, "0.0/0.0>0")
	assert(0, 0.0/0.0 >= 0, "0.0/0.0>=0")

	assert(0, !3., "!3.")
	assert(1, !0., "!0.")

	println("OK")
}
