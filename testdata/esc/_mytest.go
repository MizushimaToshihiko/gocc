package test_mytest

func assert(want int, act int, code string)
func println(format ...string)

#include "../test.h"

func main() {

#define M1(x,y) x##y
	// ASSERT(0, M1("""", "\\\\"))
	println("\\\\")

	println("OK")
}
