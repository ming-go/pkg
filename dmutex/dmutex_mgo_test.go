package dmutex

import (
	"strconv"
	"testing"
	"time"

	"git.cchntek.com/CypressModule/dbpool"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/assert"
	errors "golang.org/x/xerrors"
)

const (
	TestDatabase   = "transaction"
	TestCollection = "mtcode.lock"
	TestConnName   = "test-dmutex-mgo"
)

func loadTest() error {
	config := dbpool.Config{
		ConnName: TestConnName,
	}
	config.Host = "localhost:27017"
	config.UserName = "mongoDBUsername"
	config.Password = "mongoDBPassword"
	config.DBName = "admin"

	return dbpool.GetPool().NewDB(dbpool.Mgo, config)
}

func TestEnsureMongoDBMutexIndex(t *testing.T) {
	loadTest()
	session, err := dbpool.GetPool().GetMgoDB(TestConnName)
	defer session.Close()
	assert.NoError(t, err)

	err = EnsureMongoDBMutexIndex(nil, "", "", time.Second)
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
	session, err := dbpool.GetPool().GetMgoDB(TestConnName)
	defer session.Close()
	assert.NoError(t, err)

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
	session, err := dbpool.GetPool().GetMgoDB(TestConnName)
	defer session.Close()
	assert.NoError(t, err)

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
