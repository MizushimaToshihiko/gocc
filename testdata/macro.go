package test_macro

func assert(want int, act int, code string)
func println(format ...string)

#include "include1.h"

#

/* */ #

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

	println("OK\n")
}