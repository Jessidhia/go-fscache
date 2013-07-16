// Wrapper for github.com/tgulacsi/go-locking since it doesn't compile on windows.
//
// On windows, returns a dummy lock that always succeeds. On other OSes,
// returns a *locking.FLock.
//
// Windows also does file locking on its own, but with different
// semantics.
package lock

type FileLock interface {
	Lock() error
	Unlock() error
}
