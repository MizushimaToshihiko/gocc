package test_nameresolve

#include "test.h"

func main() {
	ASSERT(10, x)
	ASSERT(3, add2(1, 2))

	println("OK")
}

var x int = 10

func add2(x int, y int) int {
	return x + y
}


func assert(want int, act int, code string)
func println(format string)
func strcmp(s1, s2 string)
