package fscache

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"
)

// An arbitrary object that can be stringified by fmt.Sprint().
//
// The stringification is filtered to ensure it doesn't contain characters
// that are invalid on Windows, which has the most restrictive filesystem.
// The "bad" characters (\, /, :, *, ?, ", <, >, |) are replaced with _.
//
// On a list of CacheKeys, the last component is taken to represent a file
// and all the other components represent the intermediary directories.
// This means that it's not possible to have subkeys of an existing file key.
//
// NOTE: when running on Windows, directories that start with a '.' get the
// '.' replaced by a '_'. This is because regular Windows tools can't deal
// with directories starting with a dot.
type CacheKey interface{}

// All "bad characters" that can't go in Windows paths.
// It's a superset of the "bad characters" on other OSes, so this works.
var badPath = regexp.MustCompile(`[\\/:\*\?\"<>\|]`)

func stringify(stuff ...CacheKey) []string {
	ret := make([]string, len(stuff))
	for i := range stuff {
		s := fmt.Sprint(stuff[i])
		ret[i] = badPath.ReplaceAllLiteralString(s, "_")
	}
	return ret
}

// Each key but the last is treated as a directory.
// The last key is treated as a regular file.
//
// This also means that cache keys that are file-backed
// cannot have subkeys.
func (cd *CacheDir) cachePath(key ...CacheKey) string {
	parts := append([]string{cd.GetCacheDir()}, stringify(key...)...)
	p := filepath.Join(filterDots(parts...)...)
	return p
}

var invalidPath = []CacheKey{".invalid"}

// Returns the time the given key was marked as invalid.
// If the key is valid, then calling IsZero() on the returned
// time will return true.
func (cd *CacheDir) GetInvalid(key ...CacheKey) (ts time.Time) {
	invKey := append(invalidPath, key...)

	stat, _ := cd.Stat(invKey...)
	return stat.ModTime()
}

// Checks if the given key is not marked as invalid, or if it is,
// checks if it was marked more than maxDuration time ago.
//
// Calls UnsetInvalid if the keys are valid.
func (cd *CacheDir) IsValid(maxDuration time.Duration, key ...CacheKey) bool {
	ts := cd.GetInvalid(key...)

	switch {
	case ts.IsZero():
		return true
	case time.Now().Sub(ts) > maxDuration:
		cd.UnsetInvalid(key...)
		return true
	default:
		return false
	}
}

// Deletes the given key and caches it as invalid.
func (cd *CacheDir) SetInvalid(key ...CacheKey) error {
	invKey := append(invalidPath, key...)

	cd.Delete(key...)
	return cd.Touch(invKey...)
}

// Removes the given key from the invalid key cache.
func (cd *CacheDir) UnsetInvalid(key ...CacheKey) error {
	invKey := append(invalidPath, key...)

	return cd.Delete(invKey...)
}
