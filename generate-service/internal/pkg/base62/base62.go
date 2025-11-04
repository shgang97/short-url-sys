package base62

import (
	"fmt"
	"math"
	"strings"
)

const (
	base         = 62
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Encode 将数字编码为Base62字符串
func Encode(num uint64) string {
	if num == 0 {
		return "0"
	}

	var encoded []byte
	for num > 0 {
		remainder := num % base
		num /= base
		encoded = append([]byte{characterSet[remainder]}, encoded...)
	}

	return string(encoded)
}

// Decode 将数字编码为Base62字符串
func Decode(encoded string) (uint64, error) {
	var num uint64
	length := len(encoded)
	for i, char := range encoded {
		pos := strings.IndexRune(characterSet, char)
		if pos == -1 {
			return 0, fmt.Errorf("invalid character '%c'", char)
		}
		num += uint64(pos) * uint64(math.Pow(float64(base), float64(length-i-1)))
	}
	return num, nil
}

// PadLeft 在字符串左侧填充字符到指定长度
func PadLeft(str string, length int, pacChar byte) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(string(pacChar), length-len(str)) + str
}
