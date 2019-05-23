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

type KVStore struct {
	sm *sync.Map
}

func NewKVStore() *KVStore {
	return &KVStore{
		sm: new(sync.Map),
	}
}

func (kvs *KVStore) get(k string) (interface{}, bool) {
	return kvs.sm.Load(k)
}

func (kvs *KVStore) set(k string, v interface{}) {
	kvs.sm.Store(k, v)
}

func (kvs *KVStore) Get(k string) (interface{}, bool) {
	return kvs.get(k)
}

func (kvs *KVStore) Set(k string, v interface{}) {
	kvs.set(k, v)
}

func (kvs *KVStore) del(k string) {
	kvs.sm.Delete(k)
}

func (kvs *KVStore) Del(k string) {
	kvs.del(k)
}

func (kvs *KVStore) setBytes(k string, v []byte) {
	kvs.sm.Store(k, v)
}

func (kvs *KVStore) SetBytes(k string, v []byte) {
	kvs.setBytes(k, v)
}

func (kvs *KVStore) getBytes(k string) ([]byte, error) {
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

func (kvs *KVStore) GetBytes(k string) ([]byte, error) {
	return kvs.getBytes(k)
}

func (kvs *KVStore) SetString(k string, v string) {
	kvs.set(k, v)
}

func (kvs *KVStore) GetString(k string) (string, error) {
	if v, exists := kvs.get(k); exists {
		s, err := cast.ToStringE(v)
		if err != nil {
			return "", err
		}

		return s, nil
	}

	return "", ErrNotExists
}

func (kvs *KVStore) SetBoolWithPB(k string, v bool) error {
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

func (kvs *KVStore) GetBoolWithPB(k string) (bool, error) {
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

func (kvs *KVStore) SetBool(k string, v bool) {
	kvs.set(k, v)
}

func (kvs *KVStore) GetBool(k string) (bool, error) {
	if v, exists := kvs.get(k); exists {
		return cast.ToBoolE(v)
	}

	return false, ErrNotExists
}
