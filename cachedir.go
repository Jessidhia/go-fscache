package fscache

import (
	"os"
	"path/filepath"
	"sync"
)

type CacheDir struct {
	mutex sync.RWMutex

	compressionLevel int
	cacheDir         string
}

// Creates (or opens) a CacheDir using the given path.
func NewCacheDir(path string) (cd *CacheDir, err error) {
	cd = &CacheDir{
		compressionLevel: DefaultCompressionLevel,
	}
	if err = cd.SetCacheDir(path); err != nil {
		return nil, err
	}
	return
}

// Sets the directory that will back this cache.
//
// Will try to os.MkdirAll the given path; if that fails,
// then the CacheDir is not modified.
func (cd *CacheDir) SetCacheDir(path string) (err error) {
	cd.mutex.Lock()
	defer cd.mutex.Unlock()

	path = filepath.Join(filterDotsAll(filepath.SplitList(path)...)...)

	if err = os.MkdirAll(path, 0777); err != nil {
		return
	}
	cd.cacheDir = path
	return
}

// Gets the path to the cache directory.
func (cd *CacheDir) GetCacheDir() string {
	cd.mutex.RLock()
	defer cd.mutex.RUnlock()

	return cd.cacheDir
}

// Opens the file that backs the specified key.
func (cd *CacheDir) Open(key ...CacheKey) (fh *os.File, err error) {
	return os.Open(cd.cachePath(key...))
}

// Opens the file that backs the specified key using os.OpenFile.
//
// The permission bits are always 0666, which then get filtered by umask.
func (cd *CacheDir) OpenFlags(flags int, key ...CacheKey) (fh *os.File, err error) {
	return os.OpenFile(cd.cachePath(key...), flags, 0666)
}

// Creates a new file to back the specified key.
func (cd *CacheDir) Create(key ...CacheKey) (fh *os.File, err error) {
	subItem := cd.cachePath(key...)
	subDir := filepath.Dir(subItem)

	if err = os.MkdirAll(subDir, 0777); err != nil {
		return nil, err
	}
	return os.Create(subItem)
}

// Deletes the file that backs the specified key.
func (cd *CacheDir) Delete(key ...CacheKey) (err error) {
	return os.Remove(cd.cachePath(key...))
}

// Deletes the specified key and all subkeys.
func (cd *CacheDir) DeleteAll(key ...CacheKey) (err error) {
	return os.RemoveAll(cd.cachePath(key...))
}
