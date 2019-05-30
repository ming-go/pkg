package mstring

import (
	"errors"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

const Digits string = "0123456789"
const HexDigits string = "0123456789abcdefABCDEF"
const OctDigits string = "01234567"
const AsciiLowercase string = "abcdefghijklmnopqrstuvwxyz"
const AsciiUppercase string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const AsciiLetters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const Punctuation string = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
const Printable string = "" // Ref Python
const Whitespace string = ""

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Random(count int, s string) string {
	if count <= 0 {
		return ""
	}

	sLen := len(s)

	b := make([]byte, count)

	for i := range b {
		b[i] = s[rand.Intn(sLen)]
	}

	return string(b)
}

func RandomNumeric(count int) string {
	return Random(count, Digits)
}

func RandomAlphanumericpunctuation(count int) string {
	return Random(count, AsciiLetters+Digits+Punctuation)
}

func RandomAlphanumeric(count int) string {
	return Random(count, AsciiLetters+Digits)
}

func RandomAlphabetic(count int) string {
	return Random(count, AsciiLetters)
}

func RandomAlphanumericLowercase(count int) string {
	return Random(count, AsciiLowercase+Digits)
}

func RandomAlphanumericUppercase(count int) string {
	return Random(count, AsciiUppercase+Digits)
}

func RandomAlphabeticLowercase(count int) string {
	return Random(count, AsciiLowercase)
}

func RandomAlphabeticUppercase(count int) string {
	return Random(count, AsciiUppercase)
}

// TODO
func SubString(s string, beginIndex int, endIndex int) {
	// https://github.com/apache/commons-lang/blob/ae6a24dd439a7b778e35b484a3a6eae1a8eb64d7/src/main/java/org/apache/commons/lang3/StringUtils.java#L2694:26
	// http://hg.openjdk.java.net/jdk9/jdk9/jdk/file/daefa1109859/src/java.base/share/classes/java/lang/String.java#l1907
	// Ref Python slice
}

func Rotate(s string, shift int) string {
	sLen := len(s)

	if shift == 0 || sLen == 0 || (shift%sLen) == 0 {
		return s
	}

	offset := (shift % sLen)

	if offset < 0 {
		return s[-offset:] + s[0:-offset]
	}

	return s[sLen-offset:] + s[0:sLen-offset]

}

func SwapCase(s string) string {
	rs := []rune(s)

	for i, r := range rs {
		if unicode.IsLetter(r) == false {
			continue
		}

		if unicode.IsLower(r) == true {
			rs[i] = unicode.ToUpper(r)
			continue
		}

		rs[i] = unicode.ToLower(r)
	}

	return string(rs)
}

func CountMatches(s string, sub string) int {
	return 0
}

func IsNotEmpty(s string) bool {
	return len(s) != 0
}

func IsAlpha(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) == false {
			return false
		}
	}

	return true
}

func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) == false && unicode.IsDigit(r) == false {
			return false
		}
	}

	return true
}

func IsAlphanumericSpace(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) == false && unicode.IsDigit(r) == false && unicode.IsSpace(r) == false {
			return false
		}
	}

	return true
}

func IsAlphaSpace(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) == false && unicode.IsSpace(r) == false {
			return false
		}
	}

	return true
}

func IsAsciiPrintable(s string) bool {
	for _, r := range s {
		if unicode.IsPrint(r) == false {
			return false
		}
	}

	return true
}

func IsNumeric(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) == false {
			return false
		}
	}

	return true
}

func IsNumericSpace(s string) bool {
	for _, r := range s {
		if unicode.IsNumber(r) == false && unicode.IsSpace(r) == false {
			return false
		}
	}

	return true
}

func IsWhiteSpace(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) == false {
			return false
		}
	}

	return true
}

func IsAllLowerCase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) == false {
			return false
		}
	}

	return true
}

var ErrSyntax = errors.New("invalid syntax")

func Abbreviate(s string, maxWidth int) (string, error) {
	strLen := len(s)

	if strLen <= maxWidth {
		return s, nil
	}

	if maxWidth < 4 {
		return "", ErrSyntax
	}

	return s[0:maxWidth-3] + "...", nil
}

func IsAllUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) == false {
			return false
		}
	}

	return true
}

func Reverse(s string) string {
	r := []rune(s)

	for i, j := 0, len(r)-1; i < (len(r) / 2); i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}

	return string(r)
}

func ReverseDelimited(s, sep string) string {
	splits := strings.Split(s, sep)

	for i, j := 0, len(splits)-1; i < (len(splits) / 2); i, j = i+1, j-1 {
		splits[i], splits[j] = splits[j], splits[i]
	}

	return strings.Join(splits, sep)
}

func ReverseWords(s string) string {
	w := strings.Fields(s)

	for i, j := 0, len(w)-1; i < (len(w) / 2); i, j = i+1, j-1 {
		w[i], w[j] = w[j], w[i]
	}

	return strings.Join(w, " ")
}

// TODO
func InsertByte(s string, offset int, v byte) string {
	return ""
}

// TODO
func InsertRune(s string, offset int, v rune) string {
	return ""
}

func InsertString(s string, offset int, v string) string {
	return s[0:offset] + v + s[offset:]
}

func StringToRune(s string) []rune {
	return []rune(s)
}

func Center(s string, width int, fillRune rune) string {
	sLen := len(s)

	offset := width - sLen

	if width <= sLen {
		return s
	}

	appendStr := strings.Repeat(string(fillRune), offset/2)

	if offset%2 == 0 {
		return appendStr + s + appendStr
	}

	if sLen%2 == 0 {
		return string(fillRune) + appendStr + s + appendStr
	}

	return appendStr + s + appendStr + string(fillRune)
}

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

// TODO
func Float64ToString() {

}

func IsBlank(s string) bool {
	return len(strings.Trim(s, " ")) == 0
}

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
