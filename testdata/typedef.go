package test

type MyInt int
type MyInt2 [4]int

func main() {
	type t1 int
	var x1 t1 = 1
	assert(1, x1, "type t1 int;var x t1=1; x")
	type t2 struct{ a int }
	var x2 t2
	x2.a = 1
	assert(1, x2.a, "type t2 struct {a int;};var x2 t2; x2.a=1; x2.a")
	type t3 int
	var t3 t3 = 1
	assert(1, t3, "type t3 int;var t3 t3=1; t3")
	type t4 struct{ a int }
	{
		type t4 int
	}
	var x4 t4
	x4.a = 2
	assert(2, x4.a, "type t4 struct { a int;}; { type t4 int; };var x4 t4; x4.a=2; x4.a")
	var x5 MyInt = 3
	assert(3, x5, "var x5 MyInt=3; x5")
	var x6 MyInt2
	assert(16, Sizeof(x6), "var x6 MyInt2; Sizeof(x6)")

	println("OK")
}
