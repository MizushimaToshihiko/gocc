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

func main() {
	Hanoi(4, "A", "B", "C")
}
