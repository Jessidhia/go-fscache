package fscache

import (
	"os"
	"time"
)

// Calls os.Stat() on the file or folder that backs the given key.
func (cd *CacheDir) Stat(key ...CacheKey) (stat os.FileInfo, err error) {
	lock, err := cd.Lock(key...)
	if err != nil {
		return
	}
	defer lock.Unlock()

	fh, err := cd.Open(key...)
	if err != nil {
		return
	}
	defer fh.Close()

	return fh.Stat()
}

// Updates the mtime of the file backing the given key.
//
// Creates an empty file if it doesn't exist.
func (cd *CacheDir) Touch(key ...CacheKey) (err error) {
	lock, err := cd.Lock(key...)
	switch {
	case err == nil:
		// AOK
		defer lock.Unlock()
	case os.IsNotExist(err):
		// new file
	default:
		return
	}

	if err = os.Chtimes(cd.cachePath(key...), time.Now(), time.Now()); err == nil {
		return
	}

	fh, err := cd.OpenFlags(os.O_APPEND|os.O_CREATE, key...)
	switch {
	case err == nil:
		// AOK
	case os.IsNotExist(err): // directory path not created yet
		fh, err = cd.Create(key...)
		if err != nil {
			return
		}
	default:
		return
	}
	return fh.Close()
}
