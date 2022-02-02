package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	assert(int64(-5), -10+int64(5), "-10+int64(5)")
	assert(int64(-15), -10-int64(5), "-10-int64(5)")
	assert(int64(-50), -10*int64(5), "-10*int64(5)")
	assert(int64(-2), -10/int64(5), "-10/int64(5)")

	assert(1, -2 < int64(-1), "-2 < int64(-1)")
	assert(1, -2 <= int64(-1), "-2 <= int64(-1)")
	assert(0, -2 > int64(-1), "-2 > int64(-1)")
	assert(0, -2 >= int64(-1), "-2 >= int64(-1)")

	assert(1, int64(-2 < -1), "int64(-2 < -1)")
	assert(1, int64(-2 <= -1), "int64(-2 <= -1)")
	assert(0, int64(-2 > -1), "int64(-2 > -1)")
	assert(0, int64(-2 >= -1), "int64(-2 >= -1)")

	assert(0, 2147483647+2147483647+2, "2147483647+2147483647+2")
	var x1 int64
	x1 = -1
	assert(int64(-1), x1, "var x int64; x=-1; x")

	var x2 [3]byte
	x2[0] = 0
	x2[1] = 1
	x2[2] = 2
	var y2 *byte = x2 + 1
	assert(1, y2[0], "var x2 [3]byte; x2[0]=0; x2[1]=1; x2[2]=2;var y2 *byte=x2+1; y2[0]")
	assert(0, y2[-1], "var x2 [3]byte; x2[0]=0; x2[1]=1; x2[2]=2;var y2 *byte=x2+1; y2[-1]")
	type t3 struct{ a byte }
	var x3 t3
	var y3 t3
	x3.a = 5
	y3 = x3
	assert(5, y3.a, "type t3 struct{a byte;};var x t3,var y3 t3; x3.a=5; y3=x3; y3.a")

	println("OK")
}
