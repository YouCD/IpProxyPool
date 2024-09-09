package fileutil

import (
	"os"
)

// PathExists return true if given path exist.
func PathExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
