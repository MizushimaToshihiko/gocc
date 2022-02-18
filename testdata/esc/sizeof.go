package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	assert(1, Sizeof(byte), "Sizeof(byte)")
	assert(2, Sizeof(int16), "Sizeof(int16)")
	assert(4, Sizeof(int), "Sizeof(int)")
	assert(8, Sizeof(int64), "Sizeof(int64)")
	type T1 struct {
		a int
		b int
	}
	assert(8, Sizeof(T1), "Sizeof(type T1 struct {a int; b int;};)")

	var x int = 0
	assert(4, Sizeof(x+1), "Sizeof(x+1)")
	assert(8, Sizeof(-10+int64(5)), "Sizeof(-10+int64(5))")
	assert(8, Sizeof(-10 - int64(5)), "Sizeof(-10 - int64(5))")
	assert(8, Sizeof(-10 * int64(5)), "Sizeof(-10 * int64(5)")
	assert(8, Sizeof(-10 / int64(5)), "Sizeof(-10 / int64(5)")
	assert(8, Sizeof(int64(-10) + 5), "Sizeof(int64(-10) + 5)")
	assert(8, Sizeof(int64(-10) - 5), "Sizeof(int64(-10) - 5)")
	assert(8, Sizeof(int64(-10) * 5), "Sizeof(int64(-10) * 5)")
	assert(8, Sizeof(int64(-10) / 5), "Sizeof(int64(-10) / 5)")
	var i byte
	assert(1, Sizeof(i++), "var i byte; Sizeof(i++)")

	assert(1, Sizeof(int8)<<31>>31, "Sizeof(int8)<<31>>31")
	assert(1, Sizeof(int8)<<63>>63, "Sizeof(int8)<<63>>63")

  assert(8, Sizeof(1.0+2), "Sizeof(1.0+2)");
  assert(8, Sizeof(1.0-2), "Sizeof(1.0-2)");
  assert(8, Sizeof(1.0*2), "Sizeof(1.0*2)");
  assert(8, Sizeof(1.0/2), "Sizeof(1.0/2)");

	println("OK")
}
