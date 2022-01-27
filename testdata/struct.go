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
	assert(16, Sizeof(x7), "type T7 struct {a int;}; var x [4]T7; Sizeof(x7)")
	// assert(24, ({ struct {int a[3];} x[2]; sizeof(x); }));
	// assert(2, ({ struct {char a; char b;} x; sizeof(x); }));
	// assert(0, ({ struct {} x; sizeof(x); }));
	// assert(8, ({ struct {char a; int b;} x; sizeof(x); }));
	// assert(8, ({ struct {int a; char b;} x; sizeof(x); }));

	// assert(8, ({ struct t {int a; int b;} x; struct t y; sizeof(y); }));
	// assert(8, ({ struct t {int a; int b;}; struct t y; sizeof(y); }));
	// assert(2, ({ struct t {char a[2];}; { struct t {char a[4];}; } struct t y; sizeof(y); }));
	// assert(3, ({ struct t {int x;}; int t=1; struct t y; y.x=2; t+y.x; }));

	// assert(3, ({ struct t {char a;} x; struct t *y = &x; x.a=3; y->a; }));
	// assert(3, ({ struct t {char a;} x; struct t *y = &x; y->a=3; x.a; }));

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
