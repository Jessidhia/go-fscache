// Cache that uses the filesystem for indexing keys.
//
// Prefer running on case-sensitive filesystems; running on a
// case-insensitive one has the side-effect that keys are also
// case-insensitive.
package fscache
