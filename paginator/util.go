package paginator

import (
	"errors"
	"reflect"
)

func isSlice(s interface{}) bool {
	return reflect.TypeOf(s).Kind() == reflect.Slice
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func copySlice(s interface{}) (interface{}, error) {
	if !isSlice(s) {
		return nil, errors.New("s parameter must be a slice")
	}

	sT, sV := reflect.TypeOf(s), reflect.ValueOf(s)
	sC := reflect.MakeSlice(sT, sV.Len(), sV.Len())
	reflect.Copy(sC, sV)
	return sC.Interface(), nil
}
