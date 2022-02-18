package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	type T1 struct {
		a int
		b int
	}
	var x1 T1
	x1.a = 1
	x1.b = 2
	assert(1, x1.a, "type T1 struct {a int; b int;}; var x1 T1; x1.a=1; x1.b=2; x1.a")
	assert(2, x1.b, "type T1 struct {a int; b int;}; var x1 T1; x1.a=1; x1.b=2; x1.b")
	type T2 struct {
		a byte
		b byte
		c byte
	}
	var x2 T2
	x2.a = 1
	x2.b = 2
	x2.c = 3
	assert(1, x2.a, "type T2 struct {a byte; b byte; c byte;}; var x2 T2; x2.a=1; x2.b=2; x2.c=3; x2.a")
	assert(2, x2.b, "type T2 struct {a byte; b byte; c byte;}; var x2 T2; x2.a=1; x2.b=2; x2.c=3; x2.b")
	assert(3, x2.c, "type T2 struct {a byte; b byte; c byte;}; var x2 T2; x2.a=1; x2.b=2; x2.c=3; x2.c")
	type T3 struct {
		a struct {
			b byte
		}
	}
	var x3 T3
	x3.a.b = 6
	assert(6, x3.a.b, "type T3 struct { a struct { b byte; }; }; var x T3; x.a.b=6; x.a.b")

	type T4 struct {
		a int
	}
	var x4 T4
	assert(4, Sizeof(x4), "type T4 struct {a int;}; var x4 T4; Sizeof(x4)")
	type T5 struct {
		a int
		b int
	}
	var x5 T5
	assert(8, Sizeof(x5), "type T5 struct {a int; b int;};var x5 T5; Sizeof(x5)")
	type T6 struct{ a [3]int }
	var x6 T6
	assert(12, Sizeof(x6), "type T6 struct {int a[3];}; var x6 T6; Sizeof(x6)")
	type T7 struct{ a int }
	var x7 [4]T7
	assert(16, Sizeof(x7), "type T7 struct {a int}; var x [4]T7; Sizeof(x7)")
	type T8 struct{ a [3]int }
	var x8 [2]T8
	assert(24, Sizeof(x8), "type T8 struct {a [3]int};var x8 [2]T8; Sizeof(x8)")
	type T9 struct {
		a byte
		b byte
	}
	var x9 T9
	assert(2, Sizeof(x9), "type T9 struct {a byte; b byte;}; var x9 T9; Sizeof(x9)")
	type T10 struct{}
	var x10 T10
	assert(0, Sizeof(x10), "type T10 struct {};var x10 T10; Sizeof(x10)")
	type T11 struct {
		a byte
		b byte
	}
	var x11 T11
	assert(2, Sizeof(x11), "type T11 struct {a byte;b byte;};var x11 T11; Sizeof(x11)")
	type T12 struct {
		a int
		b byte
	}
	var x12 T12
	assert(8, Sizeof(x12), "type T12 struct {a int; b byte;};var x12 T12; Sizeof(x12)")
	type T13 struct {
		a int
		b int
	}
	var y13 T13
	assert(8, Sizeof(y13), "type T13 struct {a int;b int;};var y13 T13; Sizeof(y13)")

	type T14 struct{ a byte }
	var x14 T14
	var y14 *T14 = &x14
	x14.a = 3
	assert(3, y14.a, "type T14 struct {a byte;};var x14 T14; var y14 *T14= &x14; x14.a=3; y14.a")
	type T15 struct{ a byte }
	var x15 T15
	var y15 *T15 = &x15
	y15.a = 3
	assert(3, x15.a, "type T15 struct {a byte;};var x15 T15; var y15 *T15=&x; y15.a=3; x15.a")

	type T16 struct {
		a int
		b int
	}
	var x16 T16
	var y16 T16
	x16.a = 3
	y16 = x16
	assert(3, y16.a, "type T16 struct {a int;b int;};var x16 T16;var y16 T16; x16.a=3; y16=x16; y16.a")
	type T17 struct {
		a int
		b int
	}
	var x17 T17
	x17.a = 7
	var y17 T17
	var z17 *T17 = &y17
	*z17 = x17
	assert(7, y17.a, "type T17 struct {a int;b int;};var x17 T17; x17.a=7;var y17 T17;var z17 *T17=&y17; *z17=x17; y17.a")
	type T18 struct {
		a int
		b int
	}
	var x18 T18
	x18.a = 7
	var y18 T18
	var p18 *T18 = &x18
	var q18 *T18 = &y18
	*q18 = *p18
	assert(7, y18.a, "type T18 struct {a int;b int;};var x18 T18; x18.a=7;var y18 T18;var p18 *T18=&x18;var q18 *T18=&y18; *q18=*p18; y18.a")
	type T19 struct {
		a byte
		b byte
	}
	var x19 T19
	var y19 T19
	x19.a = 5
	y19 = x19
	assert(5, y19.a, "type T19 struct {a byte; b byte;};var x19 T19;var y19 T19; x19.a=5; y19=x19; y19.a")

	// the belows(comma) is not supported yet.
	// type T20 struct{ a, b int }
	// var x20, y20 T20
	// x20.a = 3
	// y20 = x20
	// assert(3, y20.a, "type T20 struct {a,b int;};var x20,y20 T20; x20.a=3; y20=x20; y20.a")
	// x20.a = 7
	// var z20 *T20 = &y
	// *z20 = x20
	// assert(7, y20.a, "type T20 struct {a,b int;};var x20 T20;var y20 T20;x20.a=7;var z20 *T20=&y; *z20=x20; y20.a")
	// var p20, q20 *T20 = &x20, &y20
	// *q20 = *p20
	// assert(7, y20.a, "type T20 struct {a,b int;};var x20 T20;var y20 T20; x20.a=7;var p20,q20 *T20=&x20,&y20; *q20=*p20; y20.a")
	// type T21 struct{ a, b byte }
	// var x21, y21 T21
	// x21.a = 5
	// y21 = x21
	// assert(5, y21.a, "type T21 struct {a,b byte;};var x21,y21 T21; x21.a=5; y21=x21; y21.a")

	type T22 struct {
		a int
		b int
	}
	var y22 T22
	assert(8, Sizeof(y22), "type T22 struct {a int;b int;};var y22 T22; Sizeof(y22)")
	type T23 struct {
		a byte
		b int64
	}
	assert(16, Sizeof(T23), "type T23 struct {a byte;b int64;}; Sizeof(T23)")
	type T24 struct {
		a byte
		b int16
	}
	assert(4, Sizeof(T24), "type T24 struct {a byte;b int16;}; Sizeof(T24)")
	var foo *struct{ x int }
	assert(8, Sizeof(foo), "var foo *struct {x int;}; Sizeof(foo)")
	type T25 struct {
		next *T25
		x    int
	}
	var a25 T25
	var b25 T25
	b25.x = 1
	a25.next = &b25
	assert(1, a25.next.x, "type T25 struct { next *T25;x int; };var a25 T25;var b25 T25; b25.x=1; a25.next=&b25; a25.next.x")
	type T26 struct{ x int }
	assert(4, Sizeof(T26), "type T26 struct{ x int; }; Sizeof(T26)")

	println("OK")
}
