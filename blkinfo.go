package blkinfo

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// BlkInfo ...
type BlkInfo struct {
	realDevPath string
	mountpoint  string
	number      string
}

// New ...
func New(devPath string) (*BlkInfo, error) {
	realDevPath, err := filepath.EvalSymlinks(devPath)
	if err != nil {
		return nil, err
	}

	mtab, err := mtab()
	if err != nil {
		return nil, err
	}

	mountpoint, err := mountpoint(mtab, realDevPath)
	if err != nil {
		return nil, err
	}

	number, err := number(realDevPath)
	if err != nil {
		return nil, err
	}

	bi := &BlkInfo{
		realDevPath: realDevPath,
		mountpoint:  mountpoint,
		number:      number,
	}

	return bi, nil
}

// UdevInfo ...
func (bi *BlkInfo) UdevInfo() (string, error) {
	udevInfo, err := readFile(filepath.Join("/", "run", "udev", "data", "b"+bi.number))
	if err != nil {
		return "", err
	}

	return udevInfo, nil
}

// OsInfo ...
func (bi *BlkInfo) OsInfo() (string, error) {
	if bi.mountpoint == "" {
		return "", errors.New("no mountpoint")
	}

	osReleasePath := filepath.Join(bi.mountpoint, "etc", "os-release")
	osInfo, err := readFile(osReleasePath)
	if err != nil {
		return "", err
	}

	return osInfo, nil
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

func mtab() (string, error) {
	mtab, err := readFile(filepath.Join("/", "etc", "mtab"))
	if err != nil {
		return "", err
	}

	return mtab, nil
}

func mountpoint(mtab string, realDevPath string) (string, error) {
	for _, line := range strings.Split(mtab, "\n") {
		fields := strings.Fields(line)

		path := fields[0]
		mountpoint := fields[1]

		if strings.HasPrefix(path, "/dev") {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return "", err
			}

			if realPath == realDevPath {
				return mountpoint, nil
			}
		}
	}

	return "", nil
}

func number(realDevPath string) (string, error) {
	devName := filepath.Base(realDevPath)
	sysBlockPath := filepath.Join("/", "sys", "block")
	fileInfos, err := ioutil.ReadDir(sysBlockPath)
	if err != nil {
		return "", err
	}

	numberPath := ""
	for _, fileInfo := range fileInfos {
		fileInfoName := fileInfo.Name()
		if strings.HasPrefix(devName, fileInfoName) {
			numberPath = filepath.Join(sysBlockPath, fileInfoName)
			if devName != fileInfoName {
				// name is a partition.
				numberPath = filepath.Join(numberPath, devName)
			}

			numberPath = filepath.Join(numberPath, "dev")
			break
		}
	}

	if numberPath == "" {
		return "", errors.New("not found")
	}

	number, err := readFile(numberPath)
	if err != nil {
		return "", err
	}

	return number, nil
}
