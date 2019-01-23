package blk

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ErrNotFound ...
var ErrNotFound = errors.New("Not Found")

func mtab() (string, error) {
	mtabPath := filepath.Join("/", "etc", "mtab")
	b, err := ioutil.ReadFile(mtabPath)
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
