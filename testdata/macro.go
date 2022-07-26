package test_macro

func assert(want int, act int, code string)
func println(format ...string)

#include "include1.h"

#

/* */ #

func ret3() int { return 3 }
func dbl(x int) int { return x*x }

func main() {
	assert(5, include1, "include1")
	assert(7, include2, "include2")

#if 0
#include "/no/such/file"
	assert(0, 1, "1")
#if nested
#endif
#endif

	var m int = 0

#if 1
	m = 5
#endif
	assert(5, m, "m")

#if 1
# if 0
#  if 1
	foo bar
#  endif
# endif
	m = 3
#endif
	assert(3, m, "m")

#if 1-1
# if 1
# endif
# if 1
# else
# endif
# if 0
# else
# endif
	m = 2
#else
# if 1
	m = 3
# endif
#endif
	assert(3, m, "m")

#if 1
	m = 2
#else
	m = 3
#endif
	assert(2, m, "m")

#if 1
	m = 2
#else
	m = 3
#endif
	assert(2, m, "m")

#if 0
	m = 1
#elif 0
	m = 2
#elif 3+5
	m = 3
#elif 1*5
	m = 4
#endif
	assert(3, m, "m")

#if 1+5
	m = 1
#elif 1
	m = 2
#elif 3
	m = 2
#endif
	assert(1, m, "m")

#if 0
	m = 1
#elif 1
# if 1
	m = 2
# else
	m = 3
# endif
#else
	m = 5
#endif
	assert(2, m, "m")

	var M1 int = 5

#define M1 3
	assert(3, M1, "M1")
#define M1 4
	assert(4, M1, "M1")

#define M1 3+4+
	assert(12, M1 5, "M1 5")

#define M1 3+4
	assert(23, M1*5, "M1*5")

#define ASSERT_ assert(
#define if 5
#define five "5"
#define END )
	ASSERT_ 5, if, five END

#undef ASSERT_
#undef if
#undef five
#undef END

	if 0 {}

#define M 5
#if M
	m = 5
#else
	m = 6
#endif
	assert(5, m, "m")

#define M 5
#if M-5
	m = 6
#elif M
	m = 5
#endif
	assert(5, m, "m")

	var M2 int = 6
#define M2 M2 + 3
	assert(9, M2, "M2")

#define M3 M2 + 3
	assert(12, M3, "M3")

	var M4 int = 3
#define M4 M5 * 5
#define M5 M4 + 2
	assert(13, M4, "M4")

#ifdef M6
	m = 5
#else
	m = 3
#endif
	assert(3, m, "m")

#define M6
#ifdef M6
	m = 5
#else
	m = 3
#endif
	assert(5, m, "m")

#ifndef M7
	m = 3
#else
	m = 5
#endif
	assert(3, m, "m")

#define M7
#ifndef M7
	m = 3
#else
	m = 5
#endif
	assert(5, m, "m")

#if 0
#ifdef NO_SUCH_MACRO
#endif
#ifndef NO_SUCH_MACRO
#endif
#else
#endif

#define M7() 1
	var M7 int = 5
	assert(1, M7(), "M7()")
	assert(5, M7, "M7")

#define M7 ()
	assert(3, ret3 M7, "ret3 M7")

#define M8(x,y) x+y
	assert(7, M8(3,4), "M8(3,4)")

#define M8(x,y) x*y
	assert(24, M8(3+4,4+5), "M8(3+4,4+5)")

#define M8(x,y) (x)*(y)
	assert(63, M8(3+4,4+5), "M8(3+4,4+5)")

#define M8(x,y) x y
	assert(9, M8(,4+5), "M8(,4+5)")

#define M8(x,y) x*y
	assert(20, M8((2+3),4), "M8((2+3),4)")

#define M8(x,y) x*y
	assert(12, M8((2,3),4), "M8((2,3),4)")

#define dbl(x) M10(x) * x
#define M10(x) dbl(x) + 3
	assert(10, dbl(2), "dbl(2)")

	println("OK\n")
}