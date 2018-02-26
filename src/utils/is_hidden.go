// +build !windows

package utils

import (
	"path/filepath"
	"strings"
)

func IsHidden(fileName string) bool {

	fName := filepath.Base(fileName)

	if strings.HasPrefix(fName, ".") && len(fName) > 1 {
		return true
	}

	return false
}
