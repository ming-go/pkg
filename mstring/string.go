package mstring

import (
	"net/url"
	"strings"
	"unsafe"
)

func IsEmpty(s string) bool {
	return len(s) == 0
}

func IsURI(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	return true
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Join(strs ...string) string {
	sb := &strings.Builder{}

	for _, str := range strs {
		sb.WriteString(str)
	}

	return sb.String()
}
