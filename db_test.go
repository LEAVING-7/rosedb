package rosedb

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	opts := DefaultOptions("/tmp/rosedb")
	db, err := Open(opts)
	if err != nil {
		t.Error("open db err ", err)
	}

	key := []byte("my_list")
	writeCount := 600000
	for i := 0; i <= writeCount; i++ {
		err := db.LPush(key, GetValue128B())
		assert.Nil(t, err)
	}
}

func TestLogFileGC(t *testing.T) {
	opts := DefaultOptions("/tmp/rosedb")
	opts.LogFileGCInterval = time.Second * 7
	opts.LogFileGCRatio = 0.00001

	db, err := Open(opts)
	if err != nil {
		t.Error("open db err ", err)
	}

	writeCount := 800000
	for i := 0; i < writeCount; i++ {
		err := db.Set(GetKey(i), GetValue128B())
		assert.Nil(t, err)
	}

	var deleted [][]byte
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100000; i++ {
		k := rand.Intn(writeCount)
		key := GetKey(k)
		err := db.Delete(key)
		assert.Nil(t, err)
		deleted = append(deleted, key)
	}

	time.Sleep(time.Second * 12)
	for _, key := range deleted {
		_, err := db.Get(key)
		assert.Equal(t, err, ErrKeyNotFound)
	}
}

func TestInMemoryDataDump(t *testing.T) {
	opts := DefaultOptions("/tmp/rosedb")
	opts.InMemoryDataDumpInterval = time.Second * 3

	db, err := Open(opts)
	if err != nil {
		t.Error("open db err ", err)
	}

	listKey := []byte("my_list")
	writeCount := 8
	for i := 0; i < writeCount; i++ {
		v := GetValue128B()
		err := db.LPush(listKey, v)
		assert.Nil(t, err)
		if i == 0 || i == writeCount-1 {
			t.Log(string(v))
		}
	}
	time.Sleep(time.Second * 4)
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions("/tmp/rosedb")
	opts.InMemoryDataDumpInterval = time.Second * 7

	db, err := Open(opts)
	if err != nil {
		t.Error("open db err ", err)
	}

	listKey := []byte("my_list")
	v, err := db.LPop(listKey)
	t.Log(string(v))
	t.Log(err)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func init() {
	rand.Seed(time.Now().Unix())
}

// GetKey length: 32 Bytes
func GetKey(n int) []byte {
	return []byte("kvstore-bench-key------" + fmt.Sprintf("%09d", n))
}

func GetValue16B() []byte {
	var str bytes.Buffer
	for i := 0; i < 16; i++ {
		str.WriteByte(alphabet[rand.Int()%36])
	}
	return []byte(str.String())
}

func GetValue128B() []byte {
	var str bytes.Buffer
	for i := 0; i < 128; i++ {
		str.WriteByte(alphabet[rand.Int()%36])
	}
	return []byte(str.String())
}

func GetValue4K() []byte {
	var str bytes.Buffer
	for i := 0; i < 4096; i++ {
		str.WriteByte(alphabet[rand.Int()%36])
	}
	return []byte(str.String())
}