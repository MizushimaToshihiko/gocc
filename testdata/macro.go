package test_macro

func assert(want int, act int, code string)
func println(format ...string)

#include "include1.h"

#

/* */ #

func main() {
	assert(5, include1, "include1")
	assert(7, include2, "include2")

	println("OK\n")
}