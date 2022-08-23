package main

// Encode a given character in UTF-8.
// Don't use it now
func encodeUft8(buf *[]byte, c int, idx int) int {
	if c <= 0x7F { // 127
		(*buf)[idx+0] = byte(c)
		return 1
	}

	if c <= 0x7FF { // 2047
		(*buf)[idx+0] = byte(0b11000000 | (c >> 6))
		(*buf)[idx+1] = byte(0b10000000 | (c & 0b00111111))
		return 2
	}

	if c <= 0xFFFF { // 65535
		(*buf)[idx+0] = byte(0b11100000 | (c >> 12))
		(*buf)[idx+1] = byte(0b10000000 | ((c >> 6) & 0b00111111))
		(*buf)[idx+2] = byte(0b10000000 | (c & 0b00111111))
		return 3
	}

	(*buf)[idx+0] = byte(0b11110000 | (c >> 18))
	(*buf)[idx+1] = byte(0b10000000 | ((c >> 12) & 0b00111111))
	(*buf)[idx+2] = byte(0b10000000 | ((c >> 6) & 0b00111111))
	(*buf)[idx+3] = byte(0b10000000 | (c & 0b00111111))
	return 4
}
