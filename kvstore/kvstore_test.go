package kvstore

import (
	"strconv"
	"testing"
)

func BenchmarkSet(b *testing.B) {
	kvs := NewKVStore()

	for n := 0; n < b.N; n++ {
		kvs.sm.Store(strconv.Itoa(n), true)
	}
}

func BenchmarkSetBool(b *testing.B) {
	kvs := NewKVStore()

	for n := 0; n < b.N; n++ {
		kvs.SetBool(strconv.Itoa(n), true)
		kvs.GetBool(strconv.Itoa(n))
	}
}

func BenchmarkSetBoolWithPB(b *testing.B) {
	kvs := NewKVStore()

	for n := 0; n < b.N; n++ {
		kvs.SetBoolWithPB(strconv.Itoa(n), true)
		kvs.GetBoolWithPB(strconv.Itoa(n))
	}
}
