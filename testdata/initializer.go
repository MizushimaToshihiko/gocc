package test_initializer

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

var g3 byte = 3
var g4 int16 = 4
var g5 int = 5
var g6 int64 = 6
var g9 [3]int = [3]int{0, 1, 2}

type gT11 struct {
	a byte
	b int
}

var g11 [2]gT11 = [2]gT11{gT11{1, 2}, gT11{3, 4}}

type gT12 struct{ a [2]int }

var g12 [2]gT12 = [2]gT12{{{1, 2}}}
var g17 string = "foobar"
var g17_2 [7]byte = "foobar"
var g18 string

var g24 int = 3
var g25 *int = &g24
var g26 [3]int = [3]int{1, 2, 3}

var g27 *int = g26 + 1
var g28 *int = &g11[1].a
var g30 = struct{ a struct{ a [3]int } }{{[3]int{1, 2, 3}}}

var g31 *int = g30.a.a
var g031 = g26
var g032 = g30.a.a

var g40 [2]struct{ a int } = [2]struct{ a int }{{1}, {3}}
var g41 [3]struct {
	a int
	b int
} = [3]struct {
	a int
	b int
}{{1, 2}, {3, 4}, {5, 6}}

var g01 = 3
var g02 = 'a'
var g03 = [3]int{1, 2, 3}

var g04 = [2][3]int{{1, 2, 3}, {4, 5, 6}}
var g04_2 = g04[1]

var g05 = "abc"
var g06 = [2]string{"abc", "def"}

var g07 = [2]struct {
	a int
	b int
	c int
}{
	{1, 2, 3},
	{4, 5, 6},
}

var g08 = [2]struct {
	a int
	b int
	c int
}{
	{a: 1, b: 2, c: 3},
	{a: 4, b: 5, c: 6},
}

var g42 int
var g43 [4]int
var g44 [4]string
var g45 bool
var g46 float64
var g47 gT11

func main() {
	// comparing strings => unimplement yet.
	// ASSERT(1, x034[0].b == "abc")

	var x01 = 3 + 1*2
	ASSERT(5, x01)
	var x02 = [3]int{1, 2, 3}
	ASSERT(1, x02[0])
	ASSERT(2, x02[1])
	ASSERT(3, x02[2])
	var x03 = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	ASSERT(1, x03[0][0])
	ASSERT(2, x03[0][1])
	ASSERT(3, x03[0][2])
	ASSERT(4, x03[1][0])
	ASSERT(5, x03[1][1])
	ASSERT(6, x03[1][2])

	var x04 = "abc"
	ASSERT('a', x04[0])

	var x05 = [2]string{"abc", "def"}
	ASSERT('a', x05[0][0])
	ASSERT('b', x05[0][1])
	ASSERT('c', x05[0][2])
	ASSERT('d', x05[1][0])
	ASSERT('e', x05[1][1])
	ASSERT('f', x05[1][2])

	var x06 = struct {
		a int
		b int
		c int
	}{1, 2, 3}
	ASSERT(1, x06.a)
	ASSERT(2, x06.b)
	ASSERT(3, x06.c)
	var x07 = [2]struct {
		a int
		b int
		c int
	}{
		{1, 2, 3},
		{4, 5, 6},
	}
	ASSERT(1, x07[0].a)
	ASSERT(2, x07[0].b)
	ASSERT(3, x07[0].c)
	ASSERT(4, x07[1].a)
	ASSERT(5, x07[1].b)
	ASSERT(6, x07[1].c)

	var x08 [3]int = [3]int{1, 2, 3}
	var x08_2 = x08
	ASSERT(1, x08_2[0])
	ASSERT(2, x08_2[1])
	ASSERT(3, x08_2[2])

	x09 := 3
	ASSERT(3, x09)

	x010 := [3]int{1, 2, 3}
	ASSERT(1, x010[0])
	ASSERT(2, x010[1])
	ASSERT(3, x010[2])

	x011 := [2]struct {
		a int
		b string
	}{{1, "abc"}, {2, "def"}}
	ASSERT(1, x011[0].a)
	ASSERT('a', x011[0].b[0])
	ASSERT('b', x011[0].b[1])
	ASSERT('c', x011[0].b[2])
	ASSERT(2, x011[1].a)
	ASSERT('d', x011[1].b[0])
	ASSERT('e', x011[1].b[1])
	ASSERT('f', x011[1].b[2])

	ASSERT(3, g01)
	ASSERT('a', g02)

	ASSERT(1, g03[0])
	ASSERT(2, g03[1])
	ASSERT(3, g03[2])

	ASSERT(1, g04[0][0])
	ASSERT(2, g04[0][1])
	ASSERT(3, g04[0][2])
	ASSERT(4, g04[1][0])
	ASSERT(5, g04[1][1])
	ASSERT(6, g04[1][2])

	ASSERT(4, g04_2[0])
	ASSERT(5, g04_2[1])
	ASSERT(6, g04_2[2])

	g04_2[2]=100
	ASSERT(100, g04_2[2])
	ASSERT(6, g04[1][2])

	ASSERT('a', g05[0])
	ASSERT('b', g05[1])
	ASSERT('c', g05[2])

	ASSERT('a', g06[0][0])
	ASSERT('b', g06[0][1])
	ASSERT('c', g06[0][2])
	ASSERT('d', g06[1][0])
	ASSERT('e', g06[1][1])
	ASSERT('f', g06[1][2])

	ASSERT(1, g07[0].a)
	ASSERT(2, g07[0].b)
	ASSERT(3, g07[0].c)
	ASSERT(4, g07[1].a)
	ASSERT(5, g07[1].b)
	ASSERT(6, g07[1].c)

	ASSERT(1, g08[0].a)
	ASSERT(2, g08[0].b)
	ASSERT(3, g08[0].c)
	ASSERT(4, g08[1].a)
	ASSERT(5, g08[1].b)
	ASSERT(6, g08[1].c)

	var x1 [3]int = [3]int{1, 2, 3}
	ASSERT(1, x1[0])
	ASSERT(2, x1[1])
	ASSERT(3, x1[2])

	var x2 [2][3]int = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	ASSERT(2, x2[0][1])
	ASSERT(4, x2[1][0])
	ASSERT(6, x2[1][2])

	var x3 [3]int = [3]int{}
	ASSERT(0, x3[0])
	ASSERT(0, x3[1])
	ASSERT(0, x3[2])

	var x4 [2][3]int = [2][3]int{{1, 2}}
	ASSERT(2, x4[0][1])
	ASSERT(0, x4[1][0])
	ASSERT(0, x4[1][2])
	var x5 [4]byte = "abc"
	ASSERT('a', x5[0])
	ASSERT('c', x5[2])
	ASSERT(0, x5[3])

	var x6 string = "def"
	ASSERT('d', x6[0])
	ASSERT('f', x6[2])
	ASSERT(0, x6[3])

	var x7 [2][4]byte = [2][4]byte{"abc", "def"}
	ASSERT('a', x7[0][0])
	ASSERT(0, x7[0][3])
	ASSERT('d', x7[1][0])
	ASSERT('f', x7[1][2])

	var x8 [2]string = [2]string{"abc", "def"}
	ASSERT('a', x8[0][0])
	ASSERT(0, x8[0][3])
	ASSERT('d', x8[1][0])
	ASSERT('f', x8[1][2])

	// assert(4, ({ int x[]={1,2,3,4}; x[3]; }));
	// assert(16, ({ int x[]={1,2,3,4}; sizeof(x); }));
	// assert(4, ({ char x[]="foo"; sizeof(x); }));

	type T9 string
	var x9 T9 = "foo"
	var y9 T9 = "x"
	ASSERT(8, Sizeof(x9))
	ASSERT(8, Sizeof(y9))
	var x10 T9 = "foo"
	var y10 T9 = "x"
	ASSERT(8, Sizeof(x10))
	ASSERT(8, Sizeof(y10))

	// assert(4, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(x); }));
	// assert(2, ({ typedef char T[]; T x="foo"; T y="x"; sizeof(y); }));
	// assert(2, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(x); }));
	// assert(4, ({ typedef char T[]; T x="x"; T y="foo"; sizeof(y); }));

	type T11 struct {
		a int
		b int
		c int
	}
	var x11 T11 = T11{1, 2, 3}
	ASSERT(1, x11.a)
	ASSERT(2, x11.b)
	ASSERT(3, x11.c)
	var x12 T11 = T11{1}
	ASSERT(1, x12.a)
	ASSERT(0, x12.b)
	ASSERT(0, x12.c)
	type T13 struct {
		a int
		b int
	}
	var x13 [2]T13 = [2]T13{T13{1, 2}, T13{3, 4}}
	ASSERT(1, x13[0].a)
	ASSERT(2, x13[0].b)
	ASSERT(3, x13[1].a)
	ASSERT(4, x13[1].b)
	type T14 struct {
		a int
		b int
	}
	var x14 [2]T14 = [2]T14{T14{1, 2}}
	ASSERT(0, x14[1].b)
	type T15 struct {
		a int
		b int
	}
	var x15 T15 = T15{}
	ASSERT(0, x15.a)
	ASSERT(0, x15.b)
	type T16 struct {
		a int
		b int
		c int
		d int
		e int
		f int
	}
	var x16 T16 = T16{1, 2, 3, 4, 5, 6}
	var y16 T16
	y16 = x16
	ASSERT(5, y16.e)
	type T17 struct {a int;b int;}; var x17 T17=T17{1,2};var y17 T17;var z17 T17; z17=y17=x17;
	ASSERT(2, z17.b)
	type T18 struct {a int;b int;}; var x18 T18=T18{1,2};var y18 T18=x18;
	ASSERT(1, y18.a);

	ASSERT(3, g3)
	ASSERT(4, g4)
	ASSERT(5, g5)
	ASSERT(6, g6)

	ASSERT(0, g9[0])
	ASSERT(1, g9[1])
	ASSERT(2, g9[2])

	ASSERT(1, g11[0].a)
	ASSERT(2, g11[0].b)
	ASSERT(3, g11[1].a)
	ASSERT(4, g11[1].b)

	ASSERT(1, g12[0].a[0])
	ASSERT(2, g12[0].a[1])
	ASSERT(0, g12[1].a[0])
	ASSERT(0, g12[1].a[1])

	ASSERT(8, Sizeof(g17))

	ASSERT('f', g17[0])
	ASSERT('o', g17[1])
	ASSERT('o', g17[2])
	ASSERT('b', g17[3])
	ASSERT('a', g17[4])
	ASSERT('r', g17[5])

	ASSERT('f', g17_2[0])
	ASSERT('o', g17_2[1])
	ASSERT('o', g17_2[2])
	ASSERT('b', g17_2[3])
	ASSERT('a', g17_2[4])
	ASSERT('r', g17_2[5])

	g18 = "foo"
	ASSERT('f', g18[0])
	ASSERT('o', g18[1])
	ASSERT('o', g18[2])

	ASSERT(3, g24)
	ASSERT(3, *g25)
	ASSERT(2, *g27)
	ASSERT(3, *g28)

	ASSERT(1, g31[0])
	ASSERT(2, g31[1]);
	ASSERT(3, g31[2]);

	ASSERT(1, g031[0])
	ASSERT(2, g031[1])
	ASSERT(3, g031[2])

	ASSERT(1, g032[0])
	ASSERT(2, g032[1])
	ASSERT(3, g032[2])

	ASSERT(1, g40[0].a)
	ASSERT(3, g40[1].a)

	ASSERT(1, g41[0].a)
	ASSERT(2, g41[0].b)
	ASSERT(3, g41[1].a)
	ASSERT(4, g41[1].b)
	ASSERT(5, g41[2].a)
	ASSERT(6, g41[2].b)
	var a [3]int=[3]int{1,2,3,};
	ASSERT(3, a[2]);
	var x19 struct {a int;b int;c int;}={1,2,3,};
	ASSERT(1, x19.a);

	type T20 struct { a int; b int; }; var x20 = T20{b:3,a:4};
	ASSERT(4, x20.a);
	ASSERT(3, x20.b);

	var x020 = struct { a int; b int; }{b:3,a:4};
	ASSERT(4, x020.a);
	ASSERT(3, x020.b);

	type T21 struct { c struct{ a int; b int; }; }; var x21 = T21{c: struct {a int;b int;}{a: 1, b: 2}};
	ASSERT(1, x21.c.a);
	ASSERT(2, x21.c.b);

	type T21 struct { c struct{ a int; b int; }; }; var x21 = &T21{c: struct {a int;b int;}{a: 1, b: 2}};
	ASSERT(1, x21.c.a);
	ASSERT(2, x21.c.b);

	var x22 = T21{c: struct {a int;b int;}{b:1},};
	ASSERT(0, x22.c.a);
	ASSERT(1, x22.c.b);

	var x23 = struct{ a [2]int;}{a: [2]int{1, 2}}
	ASSERT(1, x23.a[0]);
	ASSERT(2, x23.a[1]);

	var x24 = struct{ a [2]int }{a: [2]int{1}}
	ASSERT(1, x24.a[0]);
	ASSERT(0, x24.a[1]);

	// Initializing with 0
	var x25 string
	ASSERT(0, x25[0])

	var x26 [4]int
	ASSERT(0, x26[0])
	ASSERT(0, x26[1])
	ASSERT(0, x26[2])
	ASSERT(0, x26[3])

	var x27 [10]int
	ASSERT(0, x27[0])
	ASSERT(0, x27[1])
	ASSERT(0, x27[2])
	ASSERT(0, x27[3])
	ASSERT(0, x27[4])
	ASSERT(0, x27[5])
	ASSERT(0, x27[6])
	ASSERT(0, x27[7])
	ASSERT(0, x27[8])
	ASSERT(0, x27[9])

	var x28 [4]string
	ASSERT(0, x28[0][0])
	ASSERT(0, x28[1][0])
	ASSERT(0, x28[2][0])
	ASSERT(0, x28[3][0])

	var x29 struct {
		a int
		b string
	}
	ASSERT(0, x29.a)
	ASSERT(0, x29.b[0])

	var x30 struct {
		a [4]int
		b [4]string
	}
	ASSERT(0, x30.a[0])
	ASSERT(0, x30.a[1])
	ASSERT(0, x30.a[2])
	ASSERT(0, x30.a[3])
	ASSERT(0, x30.b[0][0])
	ASSERT(0, x30.b[1][0])
	ASSERT(0, x30.b[2][0])
	ASSERT(0, x30.b[3][0])

	var x31 int64
	ASSERT(1, x31 == 0)
	var x32 byte
	ASSERT(1, x32 == 0)
	var x33 bool
	ASSERT(1, x33 == 0)
	var x34 float64
	ASSERT(1, x34 == 0)

	ASSERT(1, g42 == 0)
	ASSERT(0, g43[0])
	ASSERT(0, g43[1])
	ASSERT(0, g43[2])
	ASSERT(0, g43[3])
	ASSERT(0, g44[0][0])
	ASSERT(0, g44[1][0])
	ASSERT(0, g44[2][0])
	ASSERT(0, g44[3][0])
	ASSERT(1, g45 == 0)
	ASSERT(1, g46 == 0)
	ASSERT(1, g47.a==0)
	ASSERT(1, g47.b==0)

	// 以下後回し案件
	// assert(3, ({ struct { int a,b; } x[]={[1].b=1,2,[0]=3,4,}; x[0].a; }));
	// assert(4, ({ struct { int a,b; } x[]={[1].b=1,2,[0]=3,4,}; x[0].b; }));
	// assert(0, ({ struct { int a,b; } x[]={[1].b=1,2,[0]=3,4,}; x[1].a; }));
	// assert(1, ({ struct { int a,b; } x[]={[1].b=1,2,[0]=3,4,}; x[1].b; }));
	// assert(2, ({ struct { int a,b; } x[]={[1].b=1,2,[0]=3,4,}; x[2].a; }));
	// assert(0, ({ struct { int a,b; } x[]={[1].b=1,2,[0]=3,4,}; x[2].b; }));

	// assert(1, ({ typedef struct { int a,b; } T; T x={1,2}; T y[]={x}; y[0].a; }));
	// assert(2, ({ typedef struct { int a,b; } T; T x={1,2}; T y[]={x}; y[0].b; }));
	// assert(0, ({ typedef struct { int a,b; } T; T x={1,2}; T y[]={x, [0].b=3}; y[0].a; }));
	// assert(3, ({ typedef struct { int a,b; } T; T x={1,2}; T y[]={x, [0].b=3}; y[0].b; }));

	// assert(5, ((struct { int a,b,c; }){ .c=5 }).c);
	// assert(0, ((struct { int a,b,c; }){ .c=5 }).a);

	println("OK")
}
