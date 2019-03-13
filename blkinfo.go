package blkinfo

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// BlkInfo shows block device information.
type BlkInfo struct {
	MajorMinor    string     `json:"major_minor"     yaml:"major_minor"    `
	Path          string     `json:"path"            yaml:"path"           `
	RealPath      string     `json:"real_path"       yaml:"real_path"      `
	ParentPath    string     `json:"parent_path"     yaml:"parent_path"    `
	ChildPaths    []string   `json:"child_paths"     yaml:"child_paths"    `
	SysPath       string     `json:"sys_path"        yaml:"sys_path"       `
	Sys           *Sys       `json:"sys"             yaml:"sys"            `
	UdevDataPath  string     `json:"udev_data_path"  yaml:"udev_data_path" `
	UdevData      []string   `json:"udev_data"       yaml:"udev_data"      `
	MountInfoPath string     `json:"mount_info_path" yaml:"mount_info_path"`
	MountInfo     *MountInfo `json:"mount_info"      yaml:"mount_info"     `
}

// Sys shows sys information.
type Sys struct {
	// See https://github.com/torvalds/linux/blob/d13937116f1e82bf508a6325111b322c30c85eb9/fs/block_dev.c#L1229-L1242
	// /sys/block/dm-0/slaves/sda  --> /sys/block/sda
	// /sys/block/sda/holders/dm-0 --> /sys/block/dm-0
	Uevent  []string `json:"uevent"  yaml:"uevent" `
	Slaves  []string `json:"slaves"  yaml:"slaves" `
	Holders []string `json:"holders" yaml:"holders"`
}

// MountInfo shows mount information.
type MountInfo struct {
	// See https://github.com/torvalds/linux/blob/d8372ba8ce288acdfce67cb873b2a741785c2e88/Documentation/filesystems/proc.txt#L1711
	MountID        string   `json:"mount_id"        yaml:"mount_id"       `
	ParentID       string   `json:"parent_id"       yaml:"parent_id"      `
	MajorMinor     string   `json:"major_minor"     yaml:"major_minor"    `
	Root           string   `json:"root"            yaml:"root"           `
	MountPoint     string   `json:"mount_point"     yaml:"mount_point"    `
	MountOptions   []string `json:"mount_options"   yaml:"mount_options"  `
	OptionalFields []string `json:"optional_fields" yaml:"optional_fields"`
	FilesystemType string   `json:"filesystem_type" yaml:"filesystem_type"`
	MountSource    string   `json:"mount_source"    yaml:"mount_source"   `
	SuperOptions   []string `json:"super_options"   yaml:"super_options"  `
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

	bi.Sys.Slaves, err = ls(filepath.Join(bi.SysPath, "slaves"))
	if err != nil {
		return nil, err
	}

	bi.Sys.Holders, err = ls(filepath.Join(bi.SysPath, "holders"))
	if err != nil {
		return nil, err
	}

	bi.MajorMinor, err = majorMinor(bi.SysPath)
	if err != nil {
		return nil, err
	}

	bi.UdevDataPath = filepath.Join("/", "run", "udev", "data", "b"+bi.MajorMinor)
	bi.UdevData, err = lines(bi.UdevDataPath)
	if err != nil {
		return nil, err
	}

	bi.MountInfoPath = filepath.Join("/", "proc", "self", "mountinfo")
	bi.MountInfo, err = newMountInfo(bi.MountInfoPath, bi.RealPath)
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

func newMountInfo(mountInfoPath string, realPath string) (*MountInfo, error) {
	lines, err := readFile(mountInfoPath)
	if err != nil {
		return nil, err
	}

	mountInfo := &MountInfo{
		MountOptions:   []string{},
		OptionalFields: []string{},
		SuperOptions:   []string{},
	}

	for _, line := range strings.Split(lines, "\n") {
		separated := strings.SplitN(line, " - ", 2)
		separatedFirst := strings.Fields(separated[0])
		separatedLast := strings.Fields(separated[1])

		if realPath == separatedLast[1] {
			mountInfo.MountID = separatedFirst[0]
			mountInfo.ParentID = separatedFirst[1]
			mountInfo.MajorMinor = separatedFirst[2]
			mountInfo.Root = separatedFirst[3]
			mountInfo.MountPoint = separatedFirst[4]
			mountInfo.MountOptions = strings.Split(separatedFirst[5], ",")
			mountInfo.OptionalFields = separatedFirst[6:]
			mountInfo.FilesystemType = separatedLast[0]
			mountInfo.MountSource = separatedLast[1]
			mountInfo.SuperOptions = strings.Split(separatedLast[2], ",")
			return mountInfo, nil
		}
	}

	return mountInfo, nil
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
