package dmutex

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	errors "golang.org/x/xerrors"
)

var (
	TestDatabase   = os.Getenv("DMUTEX_MGO_DATABASE")
	TestCollection = os.Getenv("DMUTEX_MGO_COLLECTION")
	TestHost       = os.Getenv("DMUTEX_MGO_ADDRS")
	TestUserName   = os.Getenv("DMUTEX_MGO_USERNAME")
	TestPassword   = os.Getenv("DMUTEX_MGO_PASSWORD")
)

var session *mgo.Session

func init() {
	var err error
	session, err = mgo.DialWithInfo(
		&mgo.DialInfo{
			Addrs:    []string{TestHost},
			Database: TestDatabase,
			Username: TestUserName,
			Password: TestPassword,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}

func loadTest() {
	session = session.Copy()
	defer session.Close()
	err := session.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func TestEnsureMongoDBMutexIndex(t *testing.T) {
	loadTest()
	session := session.Copy()
	defer session.Close()

	err := EnsureMongoDBMutexIndex(nil, "", "", time.Second)
	assert.True(t, errors.Is(ErrParamsIsNil, err))

	c := session.DB(TestDatabase).C(TestCollection)
	err = EnsureMongoDBMutexIndex(c, DefaultUniqueField, DefaultTTLField, 60*time.Second)
	assert.Nil(t, err)

	c = session.DB(TestDatabase).C("not_lock")
	err = EnsureMongoDBMutexIndex(c, DefaultUniqueField, DefaultTTLField, 60*time.Second)
	assert.True(t, errors.Is(ErrMutexIndexNotExist, err))

	c = session.DB(TestDatabase).C("mtcode.lock.not.exist")
	err = EnsureMongoDBMutexIndex(c, DefaultUniqueField, DefaultTTLField, 60*time.Second)
	assert.True(t, errors.Is(ErrMutexIndexNotExist, err))
}

func TestNewMongoDBImplCustomField(t *testing.T) {
	testCustomTTLField := "testttlfield"
	testCustomUniqueField := "testuniquefield"

	m, _ := NewMongoDBImpl(&mgo.Collection{}, &MutexMongoDBImplConfig{
		TTLField:    testCustomTTLField,
		UniqueField: testCustomUniqueField,
	})

	assert.Equal(t, m.ttlField, testCustomTTLField)
	assert.Equal(t, m.uniqueField, testCustomUniqueField)

}

func TestNewMongoDBImplWrongCase(t *testing.T) {
	_, err := NewMongoDBImpl(nil, &MutexMongoDBImplConfig{})
	assert.True(t, errors.Is(ErrParamsIsNil, err))

	_, err = NewMongoDBImpl(nil, nil)
	assert.True(t, errors.Is(ErrParamsIsNil, err))

	_, err = NewMongoDBImpl(&mgo.Collection{}, nil)
	assert.True(t, errors.Is(ErrParamsIsNil, err))
}

func TestMutexMongoDBImplIsLockerImplements(t *testing.T) {
	m, err := NewMongoDBImpl(&mgo.Collection{}, &MutexMongoDBImplConfig{})
	assert.Nil(t, err)
	var im interface{} = m
	_, ok := im.(Locker)
	assert.True(t, ok)
}

func TestMutexMongoDBImplWrongLockCase(t *testing.T) {
	loadTest()
	session := session.Copy()
	defer session.Close()

	c := session.DB(TestDatabase).C(TestCollection)
	testKey := "test-key-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	count, err := c.Find(bson.M{"key": testKey}).Count()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	m, err := NewMongoDBImpl(c, &MutexMongoDBImplConfig{})
	assert.Nil(t, err)
	assert.NotNil(t, m.c)
	assert.Equal(t, m.c.FullName, TestDatabase+"."+TestCollection)

	err = m.Lock(testKey)
	assert.Nil(t, err)

	count, err = c.Find(bson.M{"key": testKey}).Count()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	assert.True(t, errors.Is(ErrDuplicateKey, m.Lock(testKey)))

	err = m.Unlock(testKey)
	assert.Nil(t, err)

	count, err = c.Find(bson.M{"key": testKey}).Count()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	assert.True(t, errors.Is(ErrNotFoundKey, m.Unlock(testKey)))
}

func TestMutexMongoDBImpl(t *testing.T) {
	loadTest()
	session := session.Copy()
	defer session.Close()

	c := session.DB(TestDatabase).C(TestCollection)
	testKey := "test-key-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	count, err := c.Find(bson.M{"key": testKey}).Count()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	m, err := NewMongoDBImpl(c, &MutexMongoDBImplConfig{})
	assert.Nil(t, err)
	assert.NotNil(t, m.c)
	assert.Equal(t, m.c.FullName, TestDatabase+"."+TestCollection)

	err = m.Lock(testKey)
	assert.Nil(t, err)

	count, err = c.Find(bson.M{"key": testKey}).Count()
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	err = m.Unlock(testKey)
	assert.Nil(t, err)

	count, err = c.Find(bson.M{"key": testKey}).Count()
	assert.Nil(t, err)
	assert.Equal(t, 0, count)
}
