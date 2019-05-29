package paginator

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPageLength(t *testing.T) {
	s := []string{"one", "two", "three"}
	page, err := NewPage(s, 1, &paginator{})
	assert.Nil(t, err)
	assert.Equal(t, 3, page.Length())
}

func TestPageNextPageNumber(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	pagin, err := New(s, 2)
	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, page.NextPageNumber())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.Equal(t, 5, page.NextPageNumber())

	page, err = pagin.Page(7)
	assert.Equal(t, "pageNumber is wrong", err.Error())
}

func TestPagePreviousPageNumber(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	pagin, err := New(s, 2)
	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, page.PreviousPageNumber())

	page, err = pagin.Page(2)
	assert.Nil(t, err)
	assert.Equal(t, 1, page.PreviousPageNumber())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.Equal(t, 4, page.PreviousPageNumber())

	page, err = pagin.Page(7)
	assert.Equal(t, "pageNumber is wrong", err.Error())
}

func TestPageBind(t *testing.T) {
	type CustomStructData struct {
		tI int
		tS string
	}

	type CustomStruct struct {
		tI   int
		tS   string
		Data CustomStructData
	}

	s, sPB := []string{"one", "two", "three"}, []string{}
	page, err := NewPage(s, 1, &paginator{})
	assert.Nil(t, err)
	assert.Nil(t, page.Bind(&sPB))
	assert.True(t, reflect.DeepEqual(s, sPB))

	i, iPB := []int{1, 2, 3}, []int{}
	page, err = NewPage(i, 1, &paginator{})
	assert.Nil(t, err)
	assert.Nil(t, page.Bind(&iPB))
	assert.True(t, reflect.DeepEqual(i, iPB))

	e, ePB := []struct{}{struct{}{}, struct{}{}, struct{}{}, struct{}{}, struct{}{}}, []struct{}{}
	page, err = NewPage(e, 1, &paginator{})
	assert.Nil(t, err)
	assert.Nil(t, page.Bind(&ePB))
	assert.True(t, reflect.DeepEqual(e, ePB))

	cs, csPB := []CustomStruct{CustomStruct{tI: 1}, CustomStruct{tS: "S"}, CustomStruct{tI: 1, tS: "S", Data: CustomStructData{tI: 1, tS: "S"}}}, []CustomStruct{}
	page, err = NewPage(cs, 1, &paginator{})
	assert.Nil(t, err)
	assert.Nil(t, page.Bind(&csPB))
	assert.True(t, reflect.DeepEqual(cs, csPB))

	f, fPB := []float64{0.1, 0.2, 0.3}, []string{}

	page, err = NewPage(f, 1, &paginator{})
	assert.Nil(t, err)
	err = page.Bind(fPB)
	assert.Equal(t, "destination not a pointer", err.Error())
	var nilPtr *struct{}
	err = page.Bind(nilPtr)
	assert.Equal(t, "error nil pointer", err.Error())
	err = page.Bind(&nilPtr)
	assert.Equal(t, "cannot assignable", err.Error())
}

func TestPageHasNext(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	pagin, err := New(s, 2)
	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.True(t, page.HasNext())

	page, err = pagin.Page(4)
	assert.Nil(t, err)
	assert.True(t, page.HasNext())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.False(t, page.HasNext())

	pagin, err = New(s, 20)
	assert.Nil(t, err)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.False(t, page.HasNext())
}

func TestPageHasPrevious(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	pagin, err := New(s, 2)
	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.False(t, page.HasPrevious())

	page, err = pagin.Page(4)
	assert.Nil(t, err)
	assert.True(t, page.HasPrevious())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.True(t, page.HasPrevious())

	pagin, err = New(s, 20)
	assert.Nil(t, err)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.False(t, page.HasPrevious())
}

func TestPageHasOtherPage(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	pagin, err := New(s, 2)

	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.True(t, page.HasOtherPages())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.True(t, page.HasOtherPages())

	pagin, err = New(s, 20)
	assert.Nil(t, err)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.False(t, page.HasOtherPages())
}

func TestPageStartIndex(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	pagin, err := New(s, 2)
	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, page.StartIndex())

	page, err = pagin.Page(3)
	assert.Nil(t, err)
	assert.Equal(t, 5, page.StartIndex())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.Equal(t, 9, page.StartIndex())

	pagin, err = New(s, 3)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, page.StartIndex())

	page, err = pagin.Page(3)
	assert.Nil(t, err)
	assert.Equal(t, 7, page.StartIndex())

	page, err = pagin.Page(4)
	assert.Nil(t, err)
	assert.Equal(t, 10, page.StartIndex())

	pagin, err = New(s, 20)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, page.StartIndex())

	pagin, err = New([]struct{}{}, 10)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 0, page.StartIndex())
}

func TestPageEndIndex(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	pagin, err := New(s, 2)
	page, err := pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, page.EndIndex())

	page, err = pagin.Page(3)
	assert.Nil(t, err)
	assert.Equal(t, 6, page.EndIndex())

	page, err = pagin.Page(5)
	assert.Nil(t, err)
	assert.Equal(t, 10, page.EndIndex())

	pagin, err = New(s, 3)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 3, page.EndIndex())

	page, err = pagin.Page(3)
	assert.Nil(t, err)
	assert.Equal(t, 9, page.EndIndex())

	page, err = pagin.Page(4)
	assert.Nil(t, err)
	assert.Equal(t, 10, page.EndIndex())

	pagin, err = New(s, 20)

	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 10, page.EndIndex())

	pagin, err = New([]struct{}{}, 10)
	page, err = pagin.Page(1)
	assert.Nil(t, err)
	assert.Equal(t, 0, page.EndIndex())
}

func TestNewPage(t *testing.T) {
	s := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

	_, err := NewPage(s, 10, &paginator{})
	assert.Nil(t, err)

	// Wrong case
	_, err = NewPage(10, 10, &paginator{})
	assert.NotNil(t, err)

	_, err = NewPage(struct{}{}, 10, &paginator{})
	assert.Equal(t, "objectList Must be a slice", err.Error())

	_, err = NewPage("", 10, &paginator{})
	assert.Equal(t, "objectList Must be a slice", err.Error())
}
