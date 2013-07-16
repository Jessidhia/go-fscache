package fscache

import (
	"github.com/Kovensky/go-fscache/lock"
)

// Locks the file that backs the given key.
//
// If the call is successful, it's the caller's responsibility to call Unlock on the returned lock.
func (cd *CacheDir) Lock(key ...CacheKey) (lock.FileLock, error) {
	l, err := lock.LockFile(cd.cachePath(key...))
	if l != nil {
		l.Lock()
	}
	return l, err
}
