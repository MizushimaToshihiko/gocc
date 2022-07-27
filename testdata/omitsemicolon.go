package test_omitsemicolon

func assert(want int, act int, code string)
func println(format ...string)

#include "test.h"

func foo()            { return }
func bar() (int, int) { return 1, 2 }

func main() {
	foo()
	a01, b01 := bar()
	ASSERT(1, a01)
	ASSERT(2, b01)
	println("OK")
}
