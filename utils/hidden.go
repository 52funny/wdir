//go:build !windows

package utils

import (
	"strings"
)

// Check whether the file is hidden
func FileHidden(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}
