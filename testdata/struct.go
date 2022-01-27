package test

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
	type T6 struct {
		a [3]int
	}
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

	// assert(3, ({ struct {int a,b;} x,y; x.a=3; y=x; y.a; }));
	// assert(7, ({ struct t {int a,b;}; struct t x; x.a=7; struct t y; struct t *z=&y; *z=x; y.a; }));
	// assert(7, ({ struct t {int a,b;}; struct t x; x.a=7; struct t y, *p=&x, *q=&y; *q=*p; y.a; }));
	// assert(5, ({ struct t {char a, b;} x, y; x.a=5; y=x; y.a; }));

	// assert(3, ({ struct {int a,b;} x,y; x.a=3; y=x; y.a; }));
	// assert(7, ({ struct t {int a,b;}; struct t x; x.a=7; struct t y; struct t *z=&y; *z=x; y.a; }));
	// assert(7, ({ struct t {int a,b;}; struct t x; x.a=7; struct t y, *p=&x, *q=&y; *q=*p; y.a; }));
	// assert(5, ({ struct t {char a, b;} x, y; x.a=5; y=x; y.a; }));

	// assert(8, ({ struct t {int a; int b;} x; struct t y; sizeof(y); }));
	// assert(8, ({ struct t {int a; int b;}; struct t y; sizeof(y); }));

	// assert(16, ({ struct {char a; long b;} x; sizeof(x); }));
	// assert(4, ({ struct {char a; short b;} x; sizeof(x); }));

	// assert(8, ({ struct foo *bar; sizeof(bar); }));
	// assert(4, ({ struct T *foo; struct T {int x;}; sizeof(struct T); }));
	// assert(1, ({ struct T { struct T *next; int x; } a; struct T b; b.x=1; a.next=&b; a.next->x; }));
	// assert(4, ({ typedef struct T T; struct T { int x; }; sizeof(T); }));

	println("OK")
}
