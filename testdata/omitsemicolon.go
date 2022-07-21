package test_omitsemicolon

func assert(want int, act int, code string)
func println(format ...string)

func foo() { return }
func bar() { return 1 }

func main() {
	foo()
	assert(1, bar(), "bar()")
	println("OK")
}
