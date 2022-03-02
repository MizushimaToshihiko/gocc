package test_mytest

func assert(want int, act int, code string)
func println(format string)

func main() {
	var a = [2]int{1, 2}
	var s = a[0:2]
	assert(1, s[0], "s[0]")
	assert(2, s[1], "s[1]")
	println("OK")
}
