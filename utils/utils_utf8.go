package utils

import "unicode/utf8"

func EncodeUTF8(s string) []byte {
	rs := []rune(s)
	bs := make([]byte, len(rs)*utf8.UTFMax)

	count := 0
	for _, r := range rs {
		count += utf8.EncodeRune(bs[count:], r)
	}
	return bs[:count]
}

func EncodeUTF8ToBuf(s string, buf []byte) (n int) {
	rs := []rune(s)
	for _, r := range rs {
		n += utf8.EncodeRune(buf[n:], r)
	}
	return
}

func DecodeUTF8(buf []byte) string {
	var l = len(buf)
	rs := make([]rune, 0, l/2)
	for off := 0; off < len(buf); {
		r, n := utf8.DecodeRune(buf[off:])
		rs = append(rs, r)
		off += n
	}
	return string(rs)
}
