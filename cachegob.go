package fscache

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
)

// The default compression level of new CacheDir objects.
const DefaultCompressionLevel = gzip.BestCompression

func (cd *CacheDir) SetCompressionLevel(level int) {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	cd.compressionLevel = level
}

// Retrieves the current gzip compression level.
func (cd *CacheDir) GetCompressionLevel() int {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	return cd.compressionLevel
}

// Calls Get to retrieve the requested key from the cache.
//
// If the key is expired, then it is removed from the cache.
func (cd *CacheDir) GetAndExpire(v interface{}, max time.Duration, key ...CacheKey) (mtime time.Time, expired bool, err error) {
	mtime, err = cd.Get(v, key...)

	if err != nil && time.Now().Sub(mtime) > max {
		expired = true
		err = cd.Delete(key...)
	}
	return
}

// Gets the requested key from the cache. The given interface{} must be a pointer
// or otherwise be modifiable; otherwise Get will panic.
func (cd *CacheDir) Get(v interface{}, key ...CacheKey) (mtime time.Time, err error) {
	val := reflect.ValueOf(v)
	if k := val.Kind(); k == reflect.Ptr || k == reflect.Interface {
		val = val.Elem()
	}
	if !val.CanSet() {
		// API caller error
		panic("(*CacheDir).Get(): given interface{} is not setable")
	}

	lock, err := cd.Lock(key...)
	if err != nil {
		return
	}
	defer func() {
		// We may unlock it early
		if lock != nil {
			lock.Unlock()
		}
	}()

	fh, err := cd.Open(key...)
	if err != nil {
		return
	}
	stat, err := fh.Stat()
	if err != nil {
		return
	}
	mtime = stat.ModTime()

	buf := bytes.Buffer{}
	if _, err = io.Copy(&buf, fh); err != nil {
		fh.Close()
		return
	}
	if err = fh.Close(); err != nil {
		return
	}

	if lock != nil {
		// early unlock
		lock.Unlock()
		lock = nil
	}

	gz, err := gzip.NewReader(&buf)
	if err != nil {
		return
	}
	defer func() {
		if e := gz.Close(); err == nil {
			err = e
		}
	}()

	switch f := gz.Header.Comment; f {
	case "encoding/gob":
		dec := gob.NewDecoder(gz)
		err = dec.Decode(v)
	default:
		err = errors.New(fmt.Sprintf("Cached data (format %q) is not in a known format", f))
	}

	return
}

// Stores the given interface{} in the cache. Returns the size of the resulting file and the error, if any.
//
// Compresses the resulting data using gzip with the compression level set by SetCompressionLevel().
func (cd *CacheDir) Set(v interface{}, key ...CacheKey) (n int64, err error) {
	if v := reflect.ValueOf(v); !v.IsValid() {
		panic("reflect.ValueOf() returned invaled value")
	} else if k := v.Kind(); k == reflect.Ptr || k == reflect.Interface {
		if v.IsNil() {
			return // no point in saving nil
		}
	}

	// First we encode to memory -- we don't want to create/truncate a file and put bad data in it.
	buf := bytes.Buffer{}
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return 0, err
	}
	gz.Header.Comment = "encoding/gob"

	enc := gob.NewEncoder(gz)
	err = enc.Encode(v)

	if e := gz.Close(); err == nil {
		err = e
	}

	if err != nil {
		return 0, err
	}

	// We have good data, time to actually put it in the cache
	lock, err := cd.Lock(key...)
	switch {
	case err == nil:
		// AOK
		defer lock.Unlock()
	case os.IsNotExist(err):
		// new file
	default:
		return 0, err
	}

	fh, err := cd.Create(key...)
	if err != nil {
		return 0, err
	}
	if lock == nil {
		// the file didn't exist before, but it does now
		lock, err = cd.Lock(key...)
		if err != nil {
			return 0, err
		}
		defer lock.Unlock()
	}

	defer func() {
		if e := fh.Close(); err == nil {
			err = e
		}
	}()
	n, err = io.Copy(fh, &buf)
	return
}
