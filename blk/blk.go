package blk

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func mtab() (string, error) {
	mtabPath := filepath.Join("/", "etc", "mtab")
	b, err := ioutil.ReadFile(mtabPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}
