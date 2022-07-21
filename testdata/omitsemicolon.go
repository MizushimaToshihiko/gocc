package test_omitsemicolon

func assert(want int, act int, code string)
func println(format ...string)

func foo()            { return }
func bar() (int, int) { return 1, 2 }

func main() {
	foo()
	a01, b01 := bar()
	assert(1, a01, "a01")
	assert(2, b01, "b01")
	println("OK")
}
