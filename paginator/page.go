package paginator

import (
	"errors"
	"reflect"
)

type Page interface {
	Bind(interface{}) error
	Number() int       //
	List() interface{} //
	HasPrevious() bool
	HasNext() bool
	HasOtherPages() bool
	NextPageNumber() int
	PreviousPageNumber() int
	StartIndex() int
	EndIndex() int
	Length() int
	IsVerify() bool //
}

type page struct {
	objectList  interface{}
	rfList      reflect.Value
	number      int
	length      int
	lengthIsSet bool
	paginator   Paginator
	isVerify    bool
}

func NewPage(objectList interface{}, number int, paginator Paginator) (Page, error) {
	//paginator *Paginator)
	/*
		cOL, err := copySlice(objectList)

		if err != nil {
			return &page{isVerify: false}, err
		}
	*/

	if !isSlice(objectList) {
		return nil, errors.New("objectList Must be a slice")
	}

	return &page{
		objectList: objectList,
		rfList:     reflect.ValueOf(objectList),
		number:     number,
		paginator:  paginator,
		isVerify:   true,
	}, nil

}

func (p *page) Number() int {
	return (*p).number
}

func (p *page) List() interface{} {
	return p.objectList
}

func (p *page) Bind(dest interface{}) error {
	pCOL, err := copySlice(p.objectList)
	if err != nil {
		return err
	}

	var pV reflect.Value = reflect.ValueOf(pCOL)
	var dpV reflect.Value = reflect.ValueOf(dest)

	if dpV.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}

	if dpV.IsNil() {
		return errors.New("error nil pointer")
	}

	if !pV.IsValid() {
		pV = reflect.ValueOf(pV)
	}

	var dV reflect.Value = reflect.Indirect(dpV)

	if !pV.IsValid() || !pV.Type().AssignableTo(dV.Type()) {
		return errors.New("cannot assignable")
	}

	dV.Set(pV)
	return nil
}

func (p *page) HasNext() bool {
	return p.number < p.paginator.NumberOfPages()
}

func (p *page) HasPrevious() bool {
	return p.number > 1
}

func (p *page) HasOtherPages() bool {
	return p.HasNext() || p.HasPrevious()
}

func (p *page) NextPageNumber() int {
	if !p.paginator.ValidatePageNumber(p.number + 1) {
		return p.number
	}

	return p.number + 1
}

func (p *page) PreviousPageNumber() int {
	if !p.paginator.ValidatePageNumber(p.number - 1) {
		return p.number
	}

	return p.number - 1
}

func (p *page) Length() int {
	if !(*p).lengthIsSet {
		(*p).length = (*p).rfList.Len()
		(*p).lengthIsSet = true
	}

	return (*p).length
}

func (p *page) StartIndex() int {
	if p.paginator.Count() == 0 {
		return 0
	}

	return (p.paginator.PerPage() * (p.number - 1)) + 1
}

func (p *page) EndIndex() int {
	if p.number == p.paginator.NumberOfPages() {
		return p.paginator.Count()
	}

	return p.number * p.paginator.PerPage()
}

func (p *page) IsVerify() bool {
	return p.isVerify
}
