package kvstore

import (
	"sync"

	pbKV "github.com/ming-go/pkg/kvstore/proto/kvstore"
	"github.com/spf13/cast"
	errors "golang.org/x/xerrors"
)

var (
	ErrNotExists = errors.New("Key does not exist")
	ErrTypeError = errors.New("Type error")
)

type kvstore struct {
	sm *sync.Map
}

func NewKVStore() *kvstore {
	return &kvstore{
		sm: new(sync.Map),
	}
}

func (kvs *kvstore) get(k string) (interface{}, bool) {
	return kvs.sm.Load(k)
}

func (kvs *kvstore) set(k string, v interface{}) {
	kvs.sm.Store(k, v)
}

func (kvs *kvstore) Get(k string) (interface{}, bool) {
	return get(k)
}

func (kvs *kvstore) Set(k string, v interface{}) {
	set(k, v)
}

func (kvs *kvstore) del(k string) {
	kvs.sm.Delete(k)
}

func (kvs *kvstore) Del(k string) {
	kvs.del(k)
}

func (kvs *kvstore) setBytes(k string, v []byte) {
	kvs.sm.Store(k, v)
}

func (kvs *kvstore) SetBytes(k string, v []byte) {
	kvs.setBytes(k, v)
}

func (kvs *kvstore) getBytes(k string) ([]byte, error) {
	if v, ok := kvs.sm.Load(k); ok == true {
		switch v := v.(type) {
		case []byte:
			return v, nil
		default:
			return nil, ErrTypeError
		}
	}

	return nil, ErrNotExists
}

func (kvs *kvstore) GetBytes(k string) ([]byte, error) {
	return kvs.getBytes(k)
}

func (kvs *kvstore) SetString(k string, v string) {
	kvs.set(k, v)
}

func (kvs *kvstore) GetString(k string) (string, error) {
	if v, exists := kvs.get(k); exists {
		s, err := cast.ToStringE(v)
		if err != nil {
			return "", err
		}

		return s, nil
	}

	return "", ErrNotExists
}

func (kvs *kvstore) SetBoolWithPB(k string, v bool) error {
	p := pbKV.Bool{
		V: v,
	}

	b, err := p.Marshal()
	if err != nil {
		return err
	}

	kvs.set(k, b)

	return nil
}

func (kvs *kvstore) GetBoolWithPB(k string) (bool, error) {
	if b, exists := kvs.get(k); exists {
		pbB := pbKV.Bool{}
		switch b := b.(type) {
		case []byte:
			err := pbB.Unmarshal([]byte(b))
			if err != nil {
				return false, err
			}
			return pbB.V, nil
		default:
			return false, ErrTypeError
		}
	}

	return false, ErrNotExists
}

func (kvs *kvstore) SetBool(k string, v bool) {
	kvs.set(k, v)
}

func (kvs *kvstore) GetBool(k string) (bool, error) {
	if v, exists := kvs.get(k); exists {
		return cast.ToBoolE(v)
	}

	return false, ErrNotExists
}
