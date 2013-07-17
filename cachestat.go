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

	if err = cd.ChtimeNoLock(time.Now(), key...); err == nil {
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

// Same as Chtime, but does not try to acquire a file lock on the file
// before. Only call this if you already hold a lock to the file.
func (cd *CacheDir) ChtimeNoLock(t time.Time, key ...CacheKey) (err error) {
	return os.Chtimes(cd.cachePath(key...), time.Now(), t)
}

// Sets the mtime of the file backing the given key to the specified time.
func (cd *CacheDir) Chtime(t time.Time, key ...CacheKey) (err error) {
	lock, err := cd.Lock(key...)
	if err != nil {
		return
	}
	defer lock.Unlock()

	return cd.ChtimeNoLock(t, key...)
}
