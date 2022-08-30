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

// Read a UTF-8-encoded Unicode code point from a source file.
// We assume that source files are always in UTF-8.
//
// UTF-8 is as variable-width encoding in which one code point is
// encoed in one to four bytes. One byte UTF-8 code points are
// identical to ASCII. Non-ASCII characters are encoded using more
// than on bytes.
// =>
// Actually, it is decoded when converted to []byte, so do nothing.
// It's just looking for the following single quotes.
// func decodeUtf8(idx *int, p *[]byte) int64 {
// 	// if (*p)[*idx] < 128 {
// 	// 	*idx += 1
// 	// 	return int64((*p)[*idx])
// 	// }

// 	var i int
// 	for i = *idx; i < 4 && (*p)[i] != '\''; i++ {
// 	}
// 	*idx += *idx + i-1

// 	return i

// }
