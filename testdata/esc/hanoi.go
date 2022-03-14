package test_hanoi

func println(format ...string)

func Hanoi(n int, from string, work string, dest string) {
	if n >= 2 {
		Hanoi(n-1, from, dest, work)
	}
	println("%d: from %s to %s", n, from, dest)
	if n >= 2 {
		Hanoi(n-1, work, from, dest)
	}
}

func Hanoi2(n int, from string, work string, dest string) {
	if n >= 1 {
		Hanoi2(n-1, from, dest, work)
		println("%d: from %s to %s", n, from, dest)
		Hanoi2(n-1, work, from, dest)
	}
}

func main() {
	println("Hanoi1:")
	Hanoi(4, "A", "B", "C")
	println("Hanoi2:")
	Hanoi2(4, "A", "B", "C")
}
