package mstring

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringIsEmpty(t *testing.T) {
	assert.True(t, IsEmpty(""))
	assert.False(t, IsEmpty(" "))
	assert.False(t, IsEmpty("   "))
	assert.False(t, IsEmpty("%"))
	assert.False(t, IsEmpty(" % "))
	assert.False(t, IsEmpty("   %   "))
}

func TestStringIsBlank(t *testing.T) {
	assert.True(t, IsBlank(""))
	assert.True(t, IsBlank(" "))
	assert.True(t, IsBlank("   "))
	assert.False(t, IsBlank("%"))
	assert.False(t, IsBlank(" % "))
	assert.False(t, IsBlank("   %   "))
}

func TestIsURI(t *testing.T) {
	assert.True(t, IsURI("https://www.google.com.tw/"))
	assert.True(t, IsURI("http://www.google.com.tw/"))
	assert.True(t, IsURI("https://alone.tw/"))
	assert.True(t, IsURI("http://alone.tw/"))
	assert.False(t, IsURI("alone.tw/"))
	assert.False(t, IsURI("httpalonetw"))
}

func TestSwapCase(t *testing.T) {
	sliceTestCases := []struct {
		input  string
		output string
	}{
		{"", ""},
		{"The dog has a BONE", "tHE DOG HAS A bone"},
		{"!@#&*$^)a$_!@#_@+", "!@#&*$^)A$_!@#_@+"},
	}

	for _, testCase := range sliceTestCases {
		output := SwapCase(testCase.input)

		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Input: %v, Output: %v, Expected: %v", testCase.input, output, testCase.output)
		}
	}

}

func TestAbbreviate(t *testing.T) {
	sliceTestCases := []struct {
		input    string
		maxWidth int
		output   string
		err      error
	}{
		{"", 4, "", nil},
		{"abcdefg", 6, "abc...", nil},
		{"abcdefg", 7, "abcdefg", nil},
		{"abcdefg", 8, "abcdefg", nil},
		{"abcdefg", 4, "a...", nil},
		{"abcdefg", 3, "", ErrSyntax},
	}

	for _, testCase := range sliceTestCases {
		output, err := Abbreviate(testCase.input, testCase.maxWidth)

		if testCase.output != output || !reflect.DeepEqual(testCase.err, err) {
			t.Errorf("Input: (%v, %v), Output: (%v, %v), Expected: %v, %v", testCase.input, testCase.maxWidth, output, err, testCase.output, testCase.err)
		}
	}
}

func TestRotate(t *testing.T) {
	sliceTestCases := []struct {
		input  string
		shift  int
		output string
	}{
		{"", 4, ""},
		{"abcdefg", 0, "abcdefg"},
		{"abcdefg", 2, "fgabcde"},
		{"abcdefg", -2, "cdefgab"},
		{"abcdefg", 7, "abcdefg"},
		{"abcdefg", -7, "abcdefg"},
		{"abcdefg", 9, "fgabcde"},
		{"abcdefg", -9, "cdefgab"},
	}

	for _, testCase := range sliceTestCases {
		output := Rotate(testCase.input, testCase.shift)

		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Input: %v, Output: %v, Expected: %v", testCase.input, output, testCase.output)
		}
	}
}

func TestInsertString(t *testing.T) {
	sliceTestCases := []struct {
		input  string
		offset int
		v      string
		output string
	}{
		{"", 0, "", ""},
		{"abcdefg", 2, "INSERT", "abINSERTcdefg"},
		{"abcdefg", 6, "INSERT", "abcdefINSERTg"},
		{"abcdefg", 7, "INSERT", "abcdefgINSERT"},
	}

	for _, testCase := range sliceTestCases {
		output := InsertString(testCase.input, testCase.offset, testCase.v)

		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Input: (%v, %v, %v), Output: %v, Expected: %v", testCase.input, testCase.offset, testCase.v, output, testCase.output)
		}
	}
}

func TestCenter(t *testing.T) {
	sliceTestCases := []struct {
		input    string
		width    int
		fillRune rune
		output   string
	}{
		{"ab", 0, '*', "ab"},
		{"ab", 2, '*', "ab"},
		{"ab", 3, '*', "*ab"},
		{"abc", 4, '*', "abc*"},
		{"abc", 17, '*', "*******abc*******"},
		{"abc", 22, '*', "*********abc**********"},
		{"ab", 17, '*', "********ab*******"},
		{"ab", 22, '*', "**********ab**********"},
	}

	for _, testCase := range sliceTestCases {
		output := Center(testCase.input, testCase.width, testCase.fillRune)

		if !reflect.DeepEqual(output, testCase.output) {
			t.Errorf("Input: (%v, %v, %q), Output: %v, Expected: %v", testCase.input, testCase.width, testCase.fillRune, output, testCase.output)
		}
	}
}
