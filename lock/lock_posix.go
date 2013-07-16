// +build !windows

package lock

import "github.com/tgulacsi/go-locking"

type flockLock struct {
	locking.FLock
}

func LockFile(p string) (FileLock, error) {
	flock, err := locking.NewFLock(p)
	if err == nil {
		return &flockLock{FLock: flock}, nil
	}
	return nil, err
}

func (fl *flockLock) Lock() error {
	return fl.FLock.Lock()
}

func (fl *flockLock) Unlock() error {
	return fl.FLock.Unlock()
}
