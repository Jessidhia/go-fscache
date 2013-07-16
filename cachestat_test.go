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

func TestCache_Stat(T *testing.T) {
	cd, err := fscache.NewCacheDir(".testdir")
	if err != nil {
		T.Fatal(err)
		return
	}

	err = cd.Touch("test", "stat")
	if err != nil {
		T.Fatal(err)
		return
	}

	stat, err := cd.Stat("test", "stat")
	if err != nil {
		T.Fatal(err)
		return
	}

	T.Log("IsDir:", stat.IsDir())
	T.Log("ModTime:", stat.ModTime())
	T.Log("Name: " + stat.Name())
	T.Log("Size:", stat.Size(), "bytes")

	if stat.IsDir() || stat.ModTime().IsZero() || stat.Name() != "stat" || stat.Size() != 0 {
		T.Error("Stat returned unexpected data")
	}
}

func TestCache_Touch(T *testing.T) {
	cd, err := fscache.NewCacheDir(".testdir")
	if err != nil {
		T.Fatal(err)
		return
	}

	err = cd.Touch("test", "touch", "file")
	if err != nil {
		T.Fatal(err)
		return
	}

	stat, err := cd.Stat("test", "touch")
	if err != nil {
		T.Fatal(err)
		return
	}

	if !stat.IsDir() {
		T.Error("Expected touch to be a dir, file found")
	}

	// can be anything -- just make a nonblank file
	size, err := cd.Set(stat.ModTime(), "test", "touch", "subdir", "file")
	if err != nil {
		T.Fatal(err)
		return
	}

	stat, err = cd.Stat("test", "touch", "subdir", "file")
	if err != nil {
		T.Fatal(err)
		return
	}
	if stat.Size() != size {
		panic("File modified outside of test")
	}

	T.Log("Waiting 5ms")
	time.Sleep(5 * time.Millisecond)

	err = cd.Touch("test", "touch", "subdir", "file")
	if err != nil {
		T.Fatal(err)
	}
	stat2, err := cd.Stat("test", "touch", "subdir", "file")

	if stat.Size() != stat2.Size() {
		T.Errorf("Touch modified the size")
	}
	if !stat2.ModTime().After(stat.ModTime()) {
		T.Errorf("Touch did not update timestamp (FAT filesystem?)")
	}
}
