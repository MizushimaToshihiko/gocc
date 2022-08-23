package main

// Encode a given character in UTF-8.
// Don't use it now
func encodeUft8(buf *[]byte, c int, idx int) int {
	if c <= 0x7F {
		(*buf)[idx+0] = byte(c)
		return 1
	}

	if c <= 0x7FF {
		(*buf)[idx+0] = byte(int8(0b11000000 | (c >> 6)))
		(*buf)[idx+1] = byte(int8(0b10000000 | (c & 0b00111111)))
		return 2
	}

	if c <= 0xFFFF {
		(*buf)[idx+0] = byte(int8(0b11100000 | (c >> 12)))
		(*buf)[idx+1] = byte(int8(0b10000000 | ((c >> 6) & 0b00111111)))
		(*buf)[idx+2] = byte(int8(0b10000000 | (c & 0b00111111)))
		return 3
	}

	(*buf)[idx+0] = byte(int8(0b11110000 | (c >> 18)))
	(*buf)[idx+1] = byte(int8(0b10000000 | ((c >> 12) & 0b00111111)))
	(*buf)[idx+2] = byte(int8(0b10000000 | ((c >> 6) & 0b00111111)))
	(*buf)[idx+3] = byte(int8(0b10000000 | (c & 0b00111111)))
	return 4
}
