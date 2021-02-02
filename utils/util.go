package utils

import (
	"fmt"
	"os"
)

// ConvertSize is Converting bytes into individual units
func ConvertSize(n int64) (ans string) {
	if n < 1024 {
		// B
		ans = fmt.Sprintf("%dB", n)
	} else if n < 1024*1024 {
		// KB
		ans = fmt.Sprintf("%.2fKB", float64(n)/1024)
	} else if n < 1024*1024*1024 {
		// M
		ans = fmt.Sprintf("%.2fM", float64(n)/1024/1024)
	} else {
		// G
		ans = fmt.Sprintf("%.2fG", float64(n)/1024/1024/1024)
	}
	return
}

// PathExists is determine if a folder exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
