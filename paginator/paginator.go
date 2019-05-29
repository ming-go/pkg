package paginator

import (
	"errors"
	"math"
	"reflect"
	//"sync"
)

/*
// Generate Page Cache?
type memCache struct {
	memCache sync.Map
}
*/

type paginator struct {
	rVOL       reflect.Value
	perPage    int
	count      int
	countIsSet bool
	copySlice  bool
	isVerify   bool
}

type Paginator interface {
	Page(number int) (Page, error)
	Count() int
	NumberOfPages() int
	ValidatePageNumber(int) bool
	PerPage() int
	IsVerify() bool
}

func New(objectList interface{}, perPage int, options ...Option) (Paginator, error) {
	p := &paginator{
		perPage:    perPage,
		copySlice:  false,
		isVerify:   true,
		count:      -1,
		countIsSet: false,
	}

	for _, option := range options {
		option.Apply(p)
	}

	if p.copySlice {
		cOL, err := copySlice(objectList)
		if err != nil {
			return nil, err
		}

		p.rVOL = reflect.ValueOf(cOL)
	} else {
		if !isSlice(objectList) {
			return nil, errors.New("objectList Must be a slice")
		}

		p.rVOL = reflect.ValueOf(objectList)
	}

	return p, nil
}

func (p *paginator) ValidatePageNumber(pageNumber int) bool {
	if pageNumber < 1 {
		return false
	}

	if pageNumber > p.NumberOfPages() {
		return false
	}

	return true
}

func (p *paginator) Count() int {
	if !p.isVerify {
		return -1
	}

	if !p.countIsSet {
		p.count = p.rVOL.Len()
		p.countIsSet = true
	}

	return p.count
}

func (p *paginator) Page(pageNumber int) (Page, error) {
	if !p.isVerify {
		return nil, errors.New("Please use paginator.New() to create a paginator")
	}

	if !p.ValidatePageNumber(pageNumber) {
		return nil, errors.New("pageNumber is wrong")
	}

	bottom := (pageNumber - 1) * p.perPage
	top := bottom + p.perPage

	if top >= p.Count() {
		top = p.Count()
	}

	return NewPage(p.rVOL.Slice(bottom, top).Interface(), pageNumber, p)
}

func (p *paginator) NumberOfPages() int {
	if !p.isVerify {
		return -1
	}

	/*
		count := p.Count()
		if count == 0 {
			return 0
		}
	*/

	hits := maxInt(1, p.Count())
	return int(math.Ceil(float64(hits) / float64(p.perPage)))
}

func (p *paginator) PerPage() int {
	if !p.isVerify {
		return -1
	}

	return p.perPage
}

func (p *paginator) IsVerify() bool {
	return p.isVerify
}
