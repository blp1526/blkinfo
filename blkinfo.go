package blkinfo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"errors"
)

// BlkInfo shows a block device info.
type BlkInfo struct {
	Path         string   `json:"path"           yaml:"path"          `
	RealPath     string   `json:"real_path"      yaml:"real_path"     `
	Mountpoint   string   `json:"mountpoint"     yaml:"mountpoint"    `
	ParentPath   string   `json:"parent_path"    yaml:"parent_path"   `
	ChildPaths   []string `json:"child_paths"    yaml:"child_paths"   `
	SysPath      string   `json:"sys_path"       yaml:"sys_path"      `
	Sys          *Sys     `json:"sys"            yaml:"sys"           `
	UdevDataPath string   `json:"udev_data_path" yaml:"udev_data_path"`
	UdevData     []string `json:"udev_data"      yaml:"udev_data"     `
}

// Sys shows a sys info.
type Sys struct {
	Uevent  []string `json:"uevent"  yaml:"uevent" `
	Slaves  []string `json:"slaves"  yaml:"slaves" `
	Holders []string `json:"holders" yaml:"holders"`
}

// New initializes *BlkInfo.
func New(path string) (*BlkInfo, error) {
	var err error

	if path == "" {
		return nil, errors.New("a path is not given")
	}

	bi := &BlkInfo{
		Sys: &Sys{},
	}

	bi.Path = path
	bi.RealPath, err = filepath.EvalSymlinks(bi.Path)
	if err != nil {
		return nil, err
	}

	bi.SysPath, bi.ParentPath, bi.ChildPaths, err = relatedPaths(bi.RealPath)
	if err != nil {
		return nil, err
	}

	bi.Sys.Uevent, err = lines(filepath.Join(bi.SysPath, "uevent"))
	if err != nil {
		return nil, err
	}

	// https://github.com/torvalds/linux/blob/d13937116f1e82bf508a6325111b322c30c85eb9/fs/block_dev.c#L1229-L1242
	// /sys/block/dm-0/slaves/sda  --> /sys/block/sda
	// /sys/block/sda/holders/dm-0 --> /sys/block/dm-0
	bi.Sys.Slaves, err = ls(filepath.Join(bi.SysPath, "slaves"))
	if err != nil {
		return nil, err
	}

	bi.Sys.Holders, err = ls(filepath.Join(bi.SysPath, "holders"))
	if err != nil {
		return nil, err
	}

	majorMinor, err := majorMinor(bi.SysPath)
	if err != nil {
		return nil, err
	}

	bi.UdevDataPath = udevDataPath(majorMinor)

	bi.UdevData, err = lines(bi.UdevDataPath)
	if err != nil {
		return nil, err
	}

	mtab, err := mtab()
	if err != nil {
		return nil, err
	}

	bi.Mountpoint, err = mountpoint(mtab, bi.RealPath)
	if err != nil {
		return nil, err
	}

	return bi, nil
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(b)), nil
}

func lines(path string) ([]string, error) {
	text, err := readFile(path)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(text, "\n"), nil
}

func trimQuotationMarks(s string) string {
	for _, q := range []string{`"`, `'`} {
		if strings.HasPrefix(s, q) && strings.HasSuffix(s, q) {
			s = strings.TrimPrefix(s, q)
			s = strings.TrimSuffix(s, q)
			break
		}
	}

	return s
}

func ls(path string) ([]string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return []string{}, nil
	}

	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	names := []string{}
	for _, fileInfo := range fileInfoList {
		names = append(names, fileInfo.Name())
	}

	return names, nil
}

func mtab() (string, error) {
	mtab, err := readFile(filepath.Join("/", "etc", "mtab"))
	if err != nil {
		return "", err
	}

	return mtab, nil
}

func mountpoint(mtab string, realPath string) (string, error) {
	for _, line := range strings.Split(mtab, "\n") {
		fields := strings.Fields(line)

		pathField := fields[0]
		mountpointField := fields[1]

		if strings.HasPrefix(pathField, "/dev") {
			realPathField, err := filepath.EvalSymlinks(pathField)
			if err != nil {
				return "", err
			}

			if realPathField == realPath {
				return mountpointField, nil
			}
		}
	}

	return "", nil
}

func relatedPaths(path string) (sysPath string, parentPath string, childPaths []string, err error) {
	devName := filepath.Base(path)
	blockPath := filepath.Join("/", "sys", "block")
	fileInfoList, err := ioutil.ReadDir(blockPath)
	if err != nil {
		return "", "", []string{}, err
	}

	for _, fileInfo := range fileInfoList {
		fileInfoName := fileInfo.Name()
		if strings.HasPrefix(devName, fileInfoName) {
			switch devName {
			case fileInfoName:
				// for example /sys/block/sda
				sysPath = filepath.Join(blockPath, fileInfoName)
				fileInfoList, err = ioutil.ReadDir(sysPath)
				if err != nil {
					return "", "", []string{}, err
				}

				childPaths = []string{}
				for _, fileInfo := range fileInfoList {
					fileInfoName = fileInfo.Name()
					if strings.HasPrefix(fileInfoName, devName) {
						childPaths = append(childPaths, filepath.Join("/", "dev", fileInfoName))
					}
				}
			default:
				// for example /sys/block/sda/sda1
				sysPath = filepath.Join(blockPath, fileInfoName, devName)
				parentPath = filepath.Join("/", "dev", fileInfoName)
				childPaths = []string{}
			}

			sysPath, err = filepath.EvalSymlinks(sysPath)
			if err != nil {
				return "", "", []string{}, err
			}

			return sysPath, parentPath, childPaths, nil
		}
	}

	return "", "", []string{}, errors.New("sysPath, parentPath, and childPaths are not found")
}

func majorMinor(sysPath string) (string, error) {
	majorMinor, err := readFile(filepath.Join(sysPath, "dev"))
	if err != nil {
		return "", err
	}

	return majorMinor, nil
}

func udevDataPath(majorMinor string) string {
	return filepath.Join("/", "run", "udev", "data", "b"+majorMinor)
}

// OSRelease shows /etc/os-release.
func (bi *BlkInfo) OSRelease() ([]string, error) {
	if bi.Mountpoint == "" {
		return []string{}, errors.New("this device is not mounted")
	}

	osReleasePath := filepath.Join(bi.Mountpoint, "etc", "os-release")
	osRelease, err := lines(osReleasePath)
	if err != nil {
		return []string{}, err
	}

	return osRelease, nil
}
