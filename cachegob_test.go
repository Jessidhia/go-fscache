package fscache_test

import (
	"github.com/Kovensky/go-fscache"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestCache_Set(T *testing.T) {
	cd, err := fscache.NewCacheDir(".testdir")
	if err != nil {
		T.Fatal(err)
		return
	}

	val := rand.Int63()
	val2 := ^val // make it different

	_, err = cd.Set(val, "test", "set")
	if err != nil {
		T.Fatal(err)
		return
	}

	_, err = cd.Get(&val2, "test", "set")
	if err != nil {
		T.Fatal(err)
		return
	}

	if val != val2 {
		T.Errorf("Expected %d, got %d", val, val2)
	}
}
