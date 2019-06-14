package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWrongCase(t *testing.T) {
	_, err := New(10, 10)
	assert.Equal(t, "objectList Must be a slice", err.Error())

	_, err = New(struct{}{}, 10)
	assert.Equal(t, "objectList Must be a slice", err.Error())

	_, err = New("", 10)
	assert.Equal(t, "objectList Must be a slice", err.Error())

	_, err = New(
		10,
		10,
		[]Option{WithCopySlice(true)}...,
	)
	assert.Equal(t, "s parameter must be a slice", err.Error())

	_, err = New(
		struct{}{},
		10,
		[]Option{WithCopySlice(true)}...,
	)
	assert.Equal(t, "s parameter must be a slice", err.Error())

	_, err = New(
		"",
		10,
		[]Option{WithCopySlice(true)}...,
	)
	assert.Equal(t, "s parameter must be a slice", err.Error())
}

func TestCount(t *testing.T) {
	s := []string{"one", "two", "three"}
	p, err := New(s, 10)

	assert.Nil(t, err)
	assert.Equal(t, 3, p.Count())

	i := []int{1, 2, 3}
	p, err = New(i, 10)

	assert.Nil(t, err)
	assert.Equal(t, 3, p.Count())

	es := []struct{}{struct{}{}, struct{}{}, struct{}{}, struct{}{}, struct{}{}}
	p, err = New(es, 10)

	assert.Nil(t, err)
	assert.Equal(t, 5, p.Count())

	p = &paginator{}
	assert.Equal(t, -1, p.Count())
}

func TestNumberOfPages(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	p, err := New(s, 2)
	assert.Nil(t, err)
	assert.Equal(t, 5, p.NumberOfPages())

	p, err = New(s, 3)
	assert.Nil(t, err)
	assert.Equal(t, 4, p.NumberOfPages())

	p, err = New(s, 7)
	assert.Nil(t, err)
	assert.Equal(t, 2, p.NumberOfPages())

	p, err = New(s, 20)
	assert.Nil(t, err)
	assert.Equal(t, 1, p.NumberOfPages())

	p, err = New(s, 1)
	assert.Nil(t, err)
	assert.Equal(t, 10, p.NumberOfPages())

	emptyS := []struct{}{}
	p, err = New(emptyS, 10)
	assert.Nil(t, err)
	assert.Equal(t, 1, p.NumberOfPages())

	p = &paginator{}
	assert.Equal(t, -1, p.NumberOfPages())
}

func TestPerPage(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	p, err := New(s, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, p.PerPage())

	p = &paginator{}
	assert.Equal(t, -1, p.PerPage())
}

func TestValidatePageNumber(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	p, err := New(s, 2)

	assert.Nil(t, err)
	assert.False(t, p.ValidatePageNumber(0))
	assert.False(t, p.ValidatePageNumber(-1))

	assert.True(t, p.ValidatePageNumber(1))
	assert.True(t, p.ValidatePageNumber(5))

	assert.False(t, p.ValidatePageNumber(6))
}

func TestIsVerify(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	p, err := New(s, 2)
	assert.Nil(t, err)
	assert.True(t, p.IsVerify())

	p = &paginator{}
	assert.False(t, p.IsVerify())
}

func TestPage(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	p, err := New(s, 2)

	assert.Nil(t, err)
	page, err := p.Page(1)
	assert.True(t, page.IsVerify())

	p = &paginator{}
	page, err = p.Page(1)
	assert.NotNil(t, err)
	assert.Equal(t, "Please use paginator.New() to create a paginator", err.Error())
}
