package blkinfo

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"errors"
)

// BlkInfo shows a block device info.
type BlkInfo struct {
	RealPath   string   `json:"real_path"   yaml:"real_path"  `
	Mountpoint string   `json:"mountpoint"  yaml:"mountpoint" `
	SysfsPath  string   `json:"sysfs_path"  yaml:"sysfs_path" `
	MajorMinor string   `json:"major_minor" yaml:"major_minor"`
	UdevData   []string `json:"udev_data"   yaml:"udev_data"  `
}

// New initializes *BlkInfo.
func New(devPath string) (*BlkInfo, error) {
	var err error
	bi := &BlkInfo{}

	bi.RealPath, err = filepath.EvalSymlinks(devPath)
	if err != nil {
		return nil, err
	}

	bi.SysfsPath, err = getSysfsPath(bi.RealPath)
	if err != nil {
		return nil, err
	}

	bi.MajorMinor, err = getMajorMinor(bi.SysfsPath)
	if err != nil {
		return nil, err
	}

	bi.UdevData, err = getUdevData(bi.MajorMinor)
	if err != nil {
		return nil, err
	}

	mtab, err := getMtab()
	if err != nil {
		return nil, err
	}

	bi.Mountpoint, err = getMountpoint(mtab, bi.RealPath)
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

func getMtab() (string, error) {
	mtab, err := readFile(filepath.Join("/", "etc", "mtab"))
	if err != nil {
		return "", err
	}

	return mtab, nil
}

func getMountpoint(mtab string, realPath string) (string, error) {
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

func getSysfsPath(realPath string) (string, error) {
	// https://github.com/torvalds/linux/blob/d13937116f1e82bf508a6325111b322c30c85eb9/fs/block_dev.c#L1229-L1242
	// /sys/block/dm-0/slaves/sda  --> /sys/block/sda
	// /sys/block/sda/holders/dm-0 --> /sys/block/dm-0
	devName := filepath.Base(realPath)
	blockPath := filepath.Join("/", "sys", "block")
	fileInfos, err := ioutil.ReadDir(blockPath)
	if err != nil {
		return "", err
	}

	sysfsPath := ""
	for _, fileInfo := range fileInfos {
		fileInfoName := fileInfo.Name()
		if strings.HasPrefix(devName, fileInfoName) {
			// for example /sys/block/sda
			sysfsPath = filepath.Join(blockPath, fileInfoName)
			if devName != fileInfoName {
				// for example /sys/block/sda/sda1
				sysfsPath = filepath.Join(sysfsPath, devName)
			}

			break
		}
	}

	if sysfsPath == "" {
		return "", errors.New("not found")
	}

	return sysfsPath, nil
}

func getMajorMinor(sysfsPath string) (string, error) {
	majorMinor, err := readFile(filepath.Join(sysfsPath, "dev"))
	if err != nil {
		return "", err
	}

	return majorMinor, nil
}

func getUdevData(majorMinor string) ([]string, error) {
	rawUdevData, err := readFile(filepath.Join("/", "run", "udev", "data", "b"+majorMinor))
	if err != nil {
		return []string{}, err
	}

	udevData := strings.Split(rawUdevData, "\n")
	return udevData, nil
}

// GetOSRelease gets /etc/os-release.
func (bi *BlkInfo) GetOSRelease() ([]string, error) {
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
