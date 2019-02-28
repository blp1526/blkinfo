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
	RealPath     string   `json:"real_path"      yaml:"real_path"     `
	ParentPath   string   `json:"parent_path"    yaml:"parent_path"   `
	ChildPaths   []string `json:"child_paths"    yaml:"child_paths"   `
	SysfsPath    string   `json:"sysfs_path"     yaml:"sysfs_path"    `
	SysfsUevent  []string `json:"sysfs_uevent"   yaml:"sysfs_uevent"  `
	Holders      []string `json:"holders"        yaml:"holders"       `
	Slaves       []string `json:"slaves"         yaml:"slaves"        `
	UdevDataPath string   `json:"udev_data_path" yaml:"udev_data_path"`
	UdevData     []string `json:"udev_data"      yaml:"udev_data"     `
	Mountpoint   string   `json:"mountpoint"     yaml:"mountpoint"    `
}

// New initializes *BlkInfo.
func New(devPath string) (*BlkInfo, error) {
	var err error
	bi := &BlkInfo{}

	bi.RealPath, err = filepath.EvalSymlinks(devPath)
	if err != nil {
		return nil, err
	}

	sysfsPath, parentPath, childPaths, err := paths(bi.RealPath)
	if err != nil {
		return nil, err
	}

	bi.SysfsPath = sysfsPath
	bi.ParentPath = parentPath
	bi.ChildPaths = childPaths

	bi.SysfsUevent, err = sysfsUevent(bi.SysfsPath)
	if err != nil {
		return nil, err
	}

	// https://github.com/torvalds/linux/blob/d13937116f1e82bf508a6325111b322c30c85eb9/fs/block_dev.c#L1229-L1242
	// /sys/block/dm-0/slaves/sda  --> /sys/block/sda
	// /sys/block/sda/holders/dm-0 --> /sys/block/dm-0
	bi.Holders, err = ls(filepath.Join(bi.SysfsPath, "holders"))
	if err != nil {
		return nil, err
	}

	bi.Slaves, err = ls(filepath.Join(bi.SysfsPath, "slaves"))
	if err != nil {
		return nil, err
	}

	majorMinor, err := majorMinor(bi.SysfsPath)
	if err != nil {
		return nil, err
	}

	bi.UdevDataPath = udevDataPath(majorMinor)

	bi.UdevData, err = udevData(bi.UdevDataPath)
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

func paths(realPath string) (sysfsPath string, parentPath string, childPaths []string, err error) {
	devName := filepath.Base(realPath)
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
				sysfsPath = filepath.Join(blockPath, fileInfoName)
				fileInfoList, err := ioutil.ReadDir(sysfsPath)
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

				parentPath = ""
				return sysfsPath, parentPath, childPaths, nil
			default:
				// for example /sys/block/sda/sda1
				sysfsPath = filepath.Join(blockPath, fileInfoName, devName)
				parentPath = filepath.Join("/", "dev", fileInfoName)
				childPaths = []string{}
				return sysfsPath, parentPath, childPaths, nil
			}
		}
	}

	return "", "", []string{}, errors.New("sysfsPath, parentPath, and childPaths are not found")
}

func sysfsUevent(sysfsPath string) ([]string, error) {
	sysfsUevent, err := readFile(filepath.Join(sysfsPath, "uevent"))
	if err != nil {
		return []string{}, err
	}

	return strings.Split(sysfsUevent, "\n"), nil
}

func majorMinor(sysfsPath string) (string, error) {
	majorMinor, err := readFile(filepath.Join(sysfsPath, "dev"))
	if err != nil {
		return "", err
	}

	return majorMinor, nil
}

func udevDataPath(majorMinor string) string {
	return filepath.Join("/", "run", "udev", "data", "b"+majorMinor)
}

func udevData(udevDataPath string) ([]string, error) {
	rawUdevData, err := readFile(udevDataPath)
	if err != nil {
		return []string{}, err
	}

	udevData := strings.Split(rawUdevData, "\n")
	return udevData, nil
}

// OSRelease shows /etc/os-release.
func (bi *BlkInfo) OSRelease() ([]string, error) {
	osRelease := []string{}

	if bi.Mountpoint == "" {
		return osRelease, errors.New("this device is not mounted")
	}

	osReleasePath := filepath.Join(bi.Mountpoint, "etc", "os-release")
	rawOSRelease, err := readFile(osReleasePath)
	if err != nil {
		return osRelease, err
	}

	osRelease = strings.Split(rawOSRelease, "\n")
	return osRelease, nil
}
