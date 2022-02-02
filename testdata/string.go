package test

func assert(want int, act int, code string)
func println(format string)

func main() {
	assert(0, ""[0], "\"\"[0]")
	assert(1, Sizeof(""), "Sizeof(\"\")")

	assert(97, "abc"[0], "\"abc\"[0]")
	assert(98, "abc"[1], "\"abc\"[1]")
	assert(99, "abc"[2], "\"abc\"[2]")
	assert(0, "abc"[3], "\"abc\"[3]")
	assert(4, Sizeof("abc"), "Sizeof(\"abc\")")

	assert(7, "\a"[0], "\"\\a\"[0]")
	assert(8, "\b"[0], "\"\\b\"[0]")
	assert(9, "\t"[0], "\"\\t\"[0]")
	assert(10, "\n"[0], "\"\\n\"[0]")
	assert(11, "\v"[0], "\"\\v\"[0]")
	assert(12, "\f"[0], "\"\\f\"[0]")
	assert(13, "\r"[0], "\"\\r\"[0]")
	assert(27, "\e"[0], "\"\\e\"[0]");

	assert(106, "\j"[0], "\"\\j\"[0]");
	assert(107, "\k"[0], "\"\\k\"[0]");
	assert(108, "\l"[0], "\"\\l\"[0]");

	assert(7, "\ax\ny"[0], "\"\\ax\\ny\"[0]");
	assert(120, "\ax\ny"[1], "\"\\ax\\ny\"[1]");
	assert(10, "\ax\ny"[2], "\"\\ax\\ny\"[2]");
	assert(121, "\ax\ny"[3], "\"\\ax\\ny\"[3]");

	assert(0, "\0"[0], "\"\\0\"[0]");
	assert(16, "\20"[0], "\"\\20\"[0]");
	assert(65, "\101"[0], "\"\\101\"[0]");
	assert(104, "\1500"[0], "\"\\1500\"[0]");
	assert(0, "\x00"[0], "\"\\x00\"[0]");
	assert(119, "\x77"[0], "\"\\x77\"[0]");

	println("OK")
}
