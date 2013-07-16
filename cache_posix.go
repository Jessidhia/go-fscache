// +build !windows

package fscache

import (
	"strings"
)

func filterDots(parts ...string) []string {
	return parts
}

func filterDotsAll(parts ...string) []string {
	return parts
}
