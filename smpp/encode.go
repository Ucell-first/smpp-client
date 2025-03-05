package smpp

import (
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func isUnicode(text string) bool {
	for _, r := range text {
		if r > 127 {
			return true
		}
	}
	return false
}

func encodeUCS2(text string) []byte {
	encoder := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()
	result, _, _ := transform.Bytes(encoder, []byte(text))
	return result
}
