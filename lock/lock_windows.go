package lock

func LockFile(p string) (FileLock, error) {
	return &winLock{}, nil
}

type winLock struct{}

func (_ *winLock) Lock() error   { return nil }
func (_ *winLock) Unlock() error { return nil }
