package utils

import "strconv"

func StrToByte(_s string) byte {
	v, _ := strconv.Atoi(_s)
	return byte(v)
}

func HexToInt(_s string) int {
	v, _ := strconv.ParseInt(_s, 16, 0)
	return int(v)
}

func IndexOf(_buf *[]byte, _offset, _len int, _key byte) int {
	idx := _offset
	slice := *_buf
	for idx < _len {
		if slice[idx] == _key {
			break
		} else {
			idx++
		}
	}
	return idx
}
