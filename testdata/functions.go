package test_function

func assert(want int, act int, code string)
func println(frmt ...string)

#include "test.h"

func ret3() int {
	return 3
	return 5
}

func add2(x int, y int) int {
	return x + y
}

func sub2(x int, y int) int {
	return x - y
}

func add6(a int, b int, c int, d int, e int, f int) int {
	return a + b + c + d + e + f
}

func addx(x *int, y int) int {
	return *x + y
}

func subChar(a byte, b byte, c byte) int {
	return a - b - c
}

func fib(x int) int {
	if x <= 1 {
		return 1
	}
	return fib(x-1) + fib(x-2)
}

func subLong(a int64, b int64, c int64) int {
	return a - b - c
}

func subShort(a int16, b int16, c int16) int {
	return a - b - c
}

var g1 int

func g1Ptr() *int {
	return &g1
}

func intToChar(x int) byte {
	return x
}

func divLong(a int64, b int64) int {
	return a / b
}

func boolFnAdd(x bool) bool {
	return x + 1
}

func boolFnSub(x bool) bool {
	return x - 1
}

func paramDecay(x []int) int {
	return x[0]
}

func retNone() {
	return
}

func falseFn() bool

func trueFn() bool

func charFn() byte

func shortFn() int16

// sliceの追加後
// func addAll(n ...int) int
// func printAll(s ...string) {
// 	println(s)
// }

func add_double(x float64, y float64) float64
func add_float(x float32, y float32) float32

func add_float3(x float32, y float32, z float32) float32 {
	return x + y + z
}

func add_double3(x float64, y float64, z float64) float64 {
	return x + y + z
}

func sprintf(buf string, format ...string) string
func strcmp(s1 string, s2 string) int

func fnptr(fn func(int, int) int, a int, b int) int {
	return fn(a, b)
}

func add10_int(x1 int, x2 int, x3 int, x4 int, x5 int, x6 int, x7 int, x8 int, x9 int, x10 int) int
func add10_float(x1 float32, x2 float32, x3 float32, x4 float32, x5 float32, x6 float32, x7 float32, x8 float32, x9 float32, x10 float32) float32

func add10_double(x1 float64, x2 float64, x3 float64, x4 float64, x5 float64, x6 float64, x7 float64, x8 float64, x9 float64, x10 float64) float64

func many_args1(a int, b int, c int, d int, e int, f int, g int, h int) int {
	return g / h
}

func many_args2(a float64, b float64, c float64, d float64, e float64,
	f float64, g float64, h float64, i float64, j float64) float64 {
	return i / j
}

func many_args3(a int, b float64, c int, d int, e float64, f int,
	g float64, h int, i float64, j float64, k float64,
	l float64, m float64, n int, o int, p float64) int {
	return o / p
}

type Ty4 struct {
	a int
	b int
	c int16
	d int8
}

type Ty5 struct {
	a int
	b float32
	c float64
}

type Ty6 struct {
	a [3]uint8
}

type Ty7 struct {
	a int64
	b int64
	c int64
}

func struct_test5(x Ty5, n int) int
func struct_test4(x Ty4, n int) int
func struct_test6(x Ty6, n int) int
func struct_test7(x Ty7, n int) int

func structTest14(x Ty4, n int) int {
	switch n {
	case 0:
		return x.a
	case 1:
		return x.b
	case 2:
		return x.c
	default:
		return x.d
	}
}

func structTest15(x Ty5, n int) int {
	switch n {
	case 0:
		return x.a
	case 1:
		return x.b
	default:
		return x.c
	}
}

type Ty20 struct {
	a [10]int8
}

type Ty21 struct {
	a [20]int8
}

func struct_test24() Ty4

func struct_test25() Ty5
func struct_test26() Ty6
func struct_test27() Ty20
func struct_test28() Ty21

func struct_test34() Ty4 {
	return Ty4{10, 20, 30, 40}
}

func struct_test35() Ty5 {
	return Ty5{10, 20, 30}
}

func struct_test36() Ty6 {
	return Ty6{10, 20, 30}
}

func struct_test37() Ty20 {
	return Ty20{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
}

func struct_test38() Ty21 {
	return Ty21{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
}

type Ty021 struct {
	a string
	b string
}

func struct_test39() Ty021 {
	return Ty021{"aaa", "bbb"}
}

func add10_identList_int(x1, x2, x3, x4, x5, x6, x7, x8, x9, x10 int) int
func add10_identList_float(x1, x2, x3, x4, x5, x6, x7, x8, x9, x10 float32) float32
func add10_identList_double(x1, x2, x3, x4, x5, x6, x7, x8, x9, x10 float64) float64

func many_args_list1(a, b, c, d, e, f, g, h int) int {
	return g / h
}

func many_args_list2(a, b, c, d, e, f, g, h, i, j float64) float64 {
	return i / j
}

func many_args_list3(a int, b float64, c, d int, e float64, f int,
	g float64, h int, i, j, k float64,
	l, m float64, n, o int, p float64) int {
	return o / p
}

func multi_return() (int, int, int) {
	return 3, 5, 6
}

func funcFn() { return __func__ }

func main() {
	ASSERT(3, ret3())
	ASSERT(8, add2(3, 5))
	ASSERT(2, sub2(5, 3))
	ASSERT(21, add6(1, 2, 3, 4, 5, 6))
	ASSERT(66, add6(1, 2, add6(3, 4, 5, 6, 7, 8), 9, 10, 11))
	ASSERT(136, add6(1, 2, add6(3, add6(4, 5, 6, 7, 8, 9), 10, 11, 12, 13), 14, 15, 16))

	ASSERT(7, add2(3, 4))
	ASSERT(1, sub2(4, 3))
	ASSERT(55, fib(9))

	ASSERT(1, subChar(7, 3, 3))

	ASSERT(1, subLong(7, 3, 3))
	ASSERT(1, subShort(7, 3, 3))

	g1 = 3

	ASSERT(3, *g1Ptr())
	ASSERT(5, intToChar(261))
	ASSERT(5, intToChar(261))
	ASSERT(-5, divLong(-10, 2))

	ASSERT(1, boolFnAdd(3))
	ASSERT(0, boolFnSub(3))
	ASSERT(1, boolFnAdd(-3))
	ASSERT(0, boolFnSub(-3))
	ASSERT(1, boolFnAdd(0))
	ASSERT(1, boolFnSub(0))
	var x [2]int
	x[0] = 3
	ASSERT(3, paramDecay(x))

	retNone()

	ASSERT(1, trueFn())
	ASSERT(0, falseFn())
	ASSERT(3, charFn())
	ASSERT(5, shortFn())

	// sliceの追加後
	// ASSERT(6, addAll(3, 1, 2, 3))
	// ASSERT(5, addAll(4, 1, 2, 3, -1))
	// printAll("1", "2", "3", "4")
	// printAll("1", "2", "3", "4", "5", "6")

	ASSERT(6, int(add_float(2.3, 3.8)))
	ASSERT(6, int(add_double(2.3, 3.8)))

	ASSERT(7, int(add_float3(2.5, 2.5, 2.5)))
	ASSERT(7, int(add_double3(2.5, 2.5, 2.5)))

	var buf string
	sprintf(buf, "%.1f", float32(3.5))
	ASSERT(0, strcmp(buf, "3.5"))

	ASSERT(&ret3, ret3)
	var fn func() int = ret3
	ASSERT(3, fn())
	var fn01 = ret3
	ASSERT(3, fn01())
	fn02 := ret3
	ASSERT(3, fn02())
	fn03 := add2
	ASSERT(3, fn03(1, 2))
	ASSERT(3, fnptr(add2, 1, 2))

	ASSERT(55, add10_int(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	ASSERT(55, add10_float(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	ASSERT(55, add10_double(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))

	var buf2 string
	sprintf(buf2, "%d %.1f %.1f %.1f %d %d %.1f %d %d %d %d %.1f %d %d %.1f %.1f %.1f %.1f %d", 1, 1.0, 1.0, 1.0, 1, 1, 1.0, 1, 1, 1, 1, 1.0, 1, 1, 1.0, 1.0, 1.0, 1.0, 1)
	ASSERT(0, strcmp(buf2, "1 1.0 1.0 1.0 1 1 1.0 1 1 1 1 1.0 1 1 1.0 1.0 1.0 1.0 1"))

	ASSERT(4, many_args1(1, 2, 3, 4, 5, 6, 40, 10))
	ASSERT(4, many_args2(1, 2, 3, 4, 5, 6, 7, 8, 40, 10))
	ASSERT(8, many_args3(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 80, 10))

	x4 := Ty4{10, 20, 30, 40}
	ASSERT(10, x4.a)
	ASSERT(10, struct_test4(x4, 0))
	ASSERT(20, struct_test4(x4, 1))
	ASSERT(30, struct_test4(x4, 2))
	ASSERT(40, struct_test4(x4, 3))

	x5 := Ty5{10, 20, 30}
	ASSERT(10, struct_test5(x5, 0))
	ASSERT(20, struct_test5(x5, 1))
	ASSERT(30, struct_test5(x5, 2))

	x6 := Ty6{10, 20, 30}
	ASSERT(10, struct_test6(x6, 0))
	ASSERT(20, struct_test6(x6, 1))
	ASSERT(30, struct_test6(x6, 2))

	x7 := Ty7{10, 20, 30}
	ASSERT(10, struct_test7(x7, 0))
	ASSERT(20, struct_test7(x7, 1))
	ASSERT(30, struct_test7(x7, 2))

	x8 := Ty4{10, 20, 30, 40}
	ASSERT(10, structTest14(x8, 0))
	ASSERT(20, structTest14(x8, 1))
	ASSERT(30, structTest14(x8, 2))
	ASSERT(40, structTest14(x8, 3))

	x9 := Ty5{10, 20, 30}
	ASSERT(10, structTest15(x9, 0))
	ASSERT(20, structTest15(x9, 1))
	ASSERT(30, structTest15(x9, 2))

	ASSERT(10, struct_test24().a)
	ASSERT(20, struct_test24().b)
	ASSERT(30, struct_test24().c)
	ASSERT(40, struct_test24().d)

	ASSERT(10, struct_test25().a)
	ASSERT(20, struct_test25().b)
	ASSERT(30, struct_test25().c)

	ASSERT(10, struct_test26().a[0])
	ASSERT(20, struct_test26().a[1])
	ASSERT(30, struct_test26().a[2])

	ASSERT(10, struct_test27().a[0])
	ASSERT(60, struct_test27().a[5])
	ASSERT(100, struct_test27().a[9])

	ASSERT(1, struct_test28().a[0])
	ASSERT(5, struct_test28().a[4])
	ASSERT(10, struct_test28().a[9])
	ASSERT(15, struct_test28().a[14])
	ASSERT(20, struct_test28().a[19])

	ASSERT(10, struct_test34().a)
	ASSERT(20, struct_test34().b)
	ASSERT(30, struct_test34().c)
	ASSERT(40, struct_test34().d)

	ASSERT(10, struct_test35().a)
	ASSERT(20, struct_test35().b)
	ASSERT(30, struct_test35().c)

	ASSERT(10, struct_test36().a[0])
	ASSERT(20, struct_test36().a[1])
	ASSERT(30, struct_test36().a[2])

	ASSERT(10, struct_test37().a[0])
	ASSERT(60, struct_test37().a[5])
	ASSERT(100, struct_test37().a[9])

	ASSERT(1, struct_test38().a[0])
	ASSERT(5, struct_test38().a[4])
	ASSERT(10, struct_test38().a[9])
	ASSERT(15, struct_test38().a[14])
	ASSERT(20, struct_test38().a[19])

	ASSERT(55, add10_identList_int(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	ASSERT(55, add10_identList_float(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	ASSERT(55, add10_identList_double(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))

	ASSERT(4, many_args_list1(1, 2, 3, 4, 5, 6, 40, 10))
	ASSERT(4, many_args_list2(1, 2, 3, 4, 5, 6, 7, 8, 40, 10))
	ASSERT(8, many_args_list3(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 80, 10))

	var x10, x11, x12 = multi_return()
	ASSERT(3, x10)
	ASSERT(5, x11)
	ASSERT(6, x12)

	ASSERT(0, strcmp(struct_test39().a, "aaa"))
	ASSERT(0, strcmp(struct_test39().b, "bbb"))

	ASSERT(5, Sizeof(__func__))
	ASSERT(0, strcmp("main", __func__))
	ASSERT(0, strcmp("funcFn", funcFn()))

	println("OK")
}
