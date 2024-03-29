package test_macro

func assert(want int, act int, code string)
func println(format ...string)
func strcmp(s1 string, s2 string) int

#include "test.h"
#include "include1.h"

var main_filename1 = __FILE__
var main_line1 = __LINE__
#define LINE() __LINE__
var main_line2 = LINE()

#

/* */ #

func ret3() int { return 3 }
func dbl(x int) int { return x*x }
func add2(x, y int) { return x+y }
func add6(a,b,c,d,e,f int) { return a+b+c+d+e+f }

func main() {
	ASSERT(5, include1)
	ASSERT(7, include2)

#if 0
#include "/no/such/file"
	ASSERT(0, 1)
#if nested
#endif
#endif

	var m int = 0

#if 1
	m = 5
#endif
	ASSERT(5, m)

#if 1
# if 0
#  if 1
	foo bar
#  endif
# endif
	m = 3
#endif
	ASSERT(3, m)

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
	ASSERT(3, m)

#if 1
	m = 2
#else
	m = 3
#endif
	ASSERT(2, m)

#if 1
	m = 2
#else
	m = 3
#endif
	ASSERT(2, m)

#if 0
	m = 1
#elif 0
	m = 2
#elif 3+5
	m = 3
#elif 1*5
	m = 4
#endif
	ASSERT(3, m)

#if 1+5
	m = 1
#elif 1
	m = 2
#elif 3
	m = 2
#endif
	ASSERT(1, m)

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
	ASSERT(2, m)

	var M1 int = 5

#define M1 3
	ASSERT(3, M1)
#define M1 4
	ASSERT(4, M1)

#define M1 3+4+
	ASSERT(12, M1 5)

#define M1 3+4
	ASSERT(23, M1*5)

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
	ASSERT(5, m)

#define M 5
#if M-5
	m = 6
#elif M
	m = 5
#endif
	ASSERT(5, m)

	var M2 int = 6
#define M2 M2 + 3
	ASSERT(9, M2)

#define M3 M2 + 3
	ASSERT(12, M3)

	var M4 int = 3
#define M4 M5 * 5
#define M5 M4 + 2
	ASSERT(13, M4)

#ifdef M6
	m = 5
#else
	m = 3
#endif
	ASSERT(3, m)

#define M6
#ifdef M6
	m = 5
#else
	m = 3
#endif
	ASSERT(5, m)

#ifndef M7
	m = 3
#else
	m = 5
#endif
	ASSERT(3, m)

#define M7
#ifndef M7
	m = 3
#else
	m = 5
#endif
	ASSERT(5, m)

#if 0
#ifdef NO_SUCH_MACRO
#endif
#ifndef NO_SUCH_MACRO
#endif
#else
#endif

#define M7() 1
	var M7 int = 5
	ASSERT(1, M7())
	ASSERT(5, M7)

#define M7 ()
	ASSERT(3, ret3 M7)

#define M8(x,y) x+y
	ASSERT(7, M8(3,4))

#define M8(x,y) x*y
	ASSERT(24, M8(3+4,4+5))

#define M8(x,y) (x)*(y)
	ASSERT(63, M8(3+4,4+5))

#define M8(x,y) x y
	ASSERT(9, M8(,4+5))

#define M8(x,y) x*y
	ASSERT(20, M8((2+3),4))

#define M8(x,y) x*y
	ASSERT(12, M8((2,3),4))

#define dbl(x) M10(x) * x
#define M10(x) dbl(x) + 3
	ASSERT(10, dbl(2))

#define M11(x) #x
	ASSERT('a', M11( a!b `""c)[0])
	ASSERT('!', M11( a!b `""c)[1])
	ASSERT('b', M11( a!b `""c)[2])
	ASSERT(' ', M11( a!b `""c)[3])
	ASSERT('`', M11( a!b `""c)[4])
	ASSERT('"', M11( a!b `""c)[5])
	ASSERT('"', M11( a!b `""c)[6])
	ASSERT('c', M11( a!b `""c)[7])
	ASSERT(0, M11( a!b `""c)[8])

#define paste(x,y) x##y
	ASSERT(15, paste(1,5))
	ASSERT(255, paste(0,xff))
	foobar := 3 
	ASSERT(3, paste(foo,bar))
	ASSERT(5, paste(5,))
	ASSERT(5, paste(,5))

#define i 5
	i3 := 100
	ASSERT(101, paste(1+i,3))
#undef i

#define paste2(x) x##5
	ASSERT(26, paste2(1+2))

#define paste3(x) 2##x
	ASSERT(23, paste3(1+2))

#define paste4(x,y,z) x##y##z
	ASSERT(123, paste4(1,2,3))

#define M12
#if defined(M12)
	m = 3
#else
	m = 4
#endif
	ASSERT(3, m)

#define M12
#if defined M12
	m = 3
#else
	m = 4
#endif
	ASSERT(3, m)

#if defined(M12) -1
	m = 3
#else
	m = 4
#endif
	ASSERT(4, m)

#if defined(NO_SUSH_MACRO)
	m = 3
#else
	m = 4
#endif
	ASSERT(4, m)

#if no_such_symbol == 0
	m = 5
#else
	m = 5
#endif
	ASSERT(5, m)

#define STR(x) #x
#define M12(x) STR(x)
#define M13(x) M12(foo.x)
	ASSERT(0, strcmp(M13(bar), "foo.bar"))

#define M12 foo
#define M13(x) STR(x)
#define M14(x) M13(x.M12)
	ASSERT(0, strcmp(M14(bar), "bar.foo"))

#define M14(x) M13(x. M12)
	ASSERT(0, strcmp(M14(bar), "bar. foo"))

#include "include3.h"
	ASSERT(3, foo)

#include "include4.h"
	ASSERT(4, foo)

#define M13 "include3.h"
#include M13
	ASSERT(3, foo)

#define M13 < include4.h
#include M13 >
	ASSERT(4, foo)

#undef foo

// #error error directive test

	ASSERT(1, __STDC__)

	ASSERT(0, strcmp(main_filename1, "testdata/macro.go"))
	ASSERT(11, main_line1)
	ASSERT(13, main_line2)
	ASSERT(0, strcmp(include1_filename, "testdata/include1.h"))
	ASSERT(4, include1_line)

#define M14(...) 3
	ASSERT(3, M14())

#define M14(...) __VA_ARGS__
	ASSERT(2, M14() 2)
	ASSERT(5, M14(5))

#define M14(...) add2(__VA_ARGS__)
	ASSERT(8, M14(2, 6))

#define M14(...) add6(1,2,__VA_ARGS__,6)
	ASSERT(21, M14(3,4,5))

#define M14(x, ...) add6(1,2,x,__VA_ARGS__,6)
	ASSERT(21, M14(3,4,5))

#define M14(x, ...) x
	ASSERT(5, M14(5))

#define CONCAT(x,y) x##y
	var f0zz int = 5
	ASSERT(5, CONCAT(f,0zz))
	ASSERT(5, CONCAT(4,.57)+0.5)

	println("OK\n")
}