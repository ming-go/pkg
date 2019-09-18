/*
	Please make sure the key field has a Unique Indexes and that
	created_at field a has TTL Indexes.
*/

package dmutex

import (
	"strings"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	errors "golang.org/x/xerrors"
)

const (
	DefaultUniqueField = "key"
	DefaultTTLField    = "created_at"
)

var (
	ErrMutexIndexNotExist = errors.New("unique index or ttl index not exists")
	ErrParamsIsNil        = errors.New("params cannot be nil")
	ErrDuplicateKey       = errors.New("duplicate key error")
	ErrNotFoundKey        = errors.New("key does not exist")
)

// EnsureMongoDBMutexIndex Ensure the index required by mutex exists
func EnsureMongoDBMutexIndex(
	c *mgo.Collection,
	uniqueField string,
	ttlField string,
	expireAfter time.Duration,
) error {
	if c == nil {
		return ErrParamsIsNil
	}

	indexs, err := c.Indexes()
	if err != nil {
		if err.Error() == "no collection" {
			return ErrMutexIndexNotExist
		}

		return err
	}

	ttlFieldIsExist := false
	uniqueFieldIsExist := false

	for _, v := range indexs {
		if len(v.Key) == 1 {
			if v.Key[0] == uniqueField && v.Unique {
				uniqueFieldIsExist = true
			}

			if v.Key[0] == ttlField && v.ExpireAfter == expireAfter {
				ttlFieldIsExist = true
			}
		}
	}

	if !ttlFieldIsExist && !uniqueFieldIsExist {
		return ErrMutexIndexNotExist
	}

	return nil
}

type MutexMongoDBImplConfig struct {
	TTLField    string
	UniqueField string
}

// A mutexMongoDBImpl is a mutual exclusion lock.
// This mutex is cross-service and implemented by mongodb.
type mutexMongoDBImpl struct {
	ttlField    string
	uniqueField string
	c           *mgo.Collection
}

// NewMongoDBImpl creates a new mutex, Implemented by mongodb.
func NewMongoDBImpl(c *mgo.Collection, cfg *MutexMongoDBImplConfig) (*mutexMongoDBImpl, error) {
	if c == nil || cfg == nil {
		return nil, ErrParamsIsNil
	}

	m := &mutexMongoDBImpl{
		uniqueField: DefaultUniqueField,
		ttlField:    DefaultTTLField,
		c:           c,
	}

	if strings.Trim(cfg.UniqueField, " ") != "" {
		m.uniqueField = cfg.UniqueField
	}

	if strings.Trim(cfg.TTLField, " ") != "" {
		m.ttlField = cfg.TTLField
	}

	return m, nil
}

// Lock
// If the lock key is already in use
// will return an ErrDuplicateKey.
func (ml *mutexMongoDBImpl) Lock(key string) error {
	err := ml.c.Insert(bson.M{
		ml.uniqueField: key,
		ml.ttlField:    time.Now(),
	})

	if err != nil && mgo.IsDup(err) {
		return ErrDuplicateKey
	}

	return err
}

// Unlock
// If the lock key does not exist
// will return an ErrNotFoundKey.
func (ml *mutexMongoDBImpl) Unlock(key string) error {
	err := ml.c.Remove(bson.M{
		"key": key,
	})

	if err != nil && err == mgo.ErrNotFound {
		return ErrNotFoundKey
	}

	return err
}
