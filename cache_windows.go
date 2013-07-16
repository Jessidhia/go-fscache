package fscache

import (
	"strings"
)

func filterDots(parts ...string) []string {
	// no subdir
	if len(parts) < 2 {
		return parts
	}
	filterDotsAll(parts[:len(parts)-1]...)
	return parts
}

func filterDotsAll(parts ...string) []string {
	for i := range parts {
		if parts[i][0] == '.' {
			parts[i] = strings.Replace(parts[i], ".", "_", 1)
		}
	}
	return parts
}
