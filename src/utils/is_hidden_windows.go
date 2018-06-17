// +build windows

package utils

import (
	"path/filepath"
	"strings"
	"syscall"
)

func IsHidden(fileName string) bool {

	fName := filepath.Base(fileName)

	if strings.HasPrefix(fName, ".") && len(fName) > 1 {
		return true
	}

	p, e := syscall.UTF16PtrFromString(fileName)
	if e != nil {
		return false
	}

	attrs, e := syscall.GetFileAttributes(p)
	if e != nil {
		return false
	}
	return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
