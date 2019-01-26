package blk

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Errors
var (
	// ErrNotFound ...
	ErrNotFound = errors.New("Not Found")
)

// Paths
var (
	// MtabPath ...
	MtabPath = filepath.Join("/", "etc", "mtab")
	// SysBlockPath ...
	SysBlockPath = filepath.Join("/", "sys", "block")
)

func mtab() (string, error) {
	b, err := ioutil.ReadFile(MtabPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

// GetPath ...
func GetPath(mountpoint string) (string, error) {
	mtab, err := mtab()
	if err != nil {
		return "", err
	}

	lines := strings.Split(mtab, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[1] == mountpoint {
			return fields[0], nil
		}
	}

	return "", ErrNotFound
}
