// Package blkinfo implements methods for block device information.
package blkinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// BlkInfo shows block device information.
type BlkInfo struct {
	Path             string            `json:"path"             `
	ResolvedPath     string            `json:"resolved_path"    `
	ParentPath       string            `json:"parent_path"      `
	ChildPaths       []string          `json:"child_paths"      `
	SysPath          string            `json:"sys_path"         `
	ResolevedSysPath string            `json:"resolved_sys_path"`
	Sys              *Sys              `json:"sys"              `
	MajorMinor       string            `json:"major_minor"      `
	UdevDataPath     string            `json:"udev_data_path"   `
	UdevData         []string          `json:"udev_data"        `
	MountInfoPath    string            `json:"mount_info_path"  `
	MountInfo        *MountInfo        `json:"mount_info"       `
	OSReleasePath    string            `json:"os_release_path"  `
	OSRelease        map[string]string `json:"os_release"       `
}

// Sys shows sys information.
type Sys struct {
	// See https://github.com/torvalds/linux/blob/d13937116f1e82bf508a6325111b322c30c85eb9/fs/block_dev.c#L1229-L1242
	// /sys/block/dm-0/slaves/sda  --> /sys/block/sda
	// /sys/block/sda/holders/dm-0 --> /sys/block/dm-0
	Uevent  []string `json:"uevent" `
	Slaves  []string `json:"slaves" `
	Holders []string `json:"holders"`
}

// MountInfo shows mount information.
type MountInfo struct {
	// See https://github.com/torvalds/linux/blob/d8372ba8ce288acdfce67cb873b2a741785c2e88/Documentation/filesystems/proc.txt#L1711
	MountID        string   `json:"mount_id"       `
	ParentID       string   `json:"parent_id"      `
	MajorMinor     string   `json:"major_minor"    `
	Root           string   `json:"root"           `
	MountPoint     string   `json:"mount_point"    `
	MountOptions   []string `json:"mount_options"  `
	OptionalFields []string `json:"optional_fields"`
	FilesystemType string   `json:"filesystem_type"`
	MountSource    string   `json:"mount_source"   `
	SuperOptions   []string `json:"super_options"  `
}

// New initializes *BlkInfo.
func New(path string) (*BlkInfo, error) { // nolint: funlen
	var err error

	if path == "" {
		return nil, fmt.Errorf("a path is not given")
	}

	bi := &BlkInfo{
		Sys: &Sys{},
	}

	bi.Path = path
	bi.ResolvedPath, err = filepath.EvalSymlinks(bi.Path)

	if err != nil {
		return nil, err
	}

	bi.SysPath, bi.ResolevedSysPath, bi.ParentPath, bi.ChildPaths, err = relatedPaths(bi.ResolvedPath)
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
	bi.MountInfo, err = newMountInfo(bi.MountInfoPath, bi.ResolvedPath)

	if err != nil {
		return nil, err
	}

	bi.OSReleasePath = osReleasePath(bi.MountInfo.MountPoint)

	bi.OSRelease, err = newOSRelease(bi.OSReleasePath)

	if err != nil {
		return nil, err
	}

	return bi, nil
}

func readFile(path string) (string, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
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

func newMountInfo(mountInfoPath string, path string) (*MountInfo, error) {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, err
	}

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

		mountSource := separatedLast[1]
		if !strings.HasPrefix(mountSource, "/dev") {
			continue
		}

		realMountSource, err := filepath.EvalSymlinks(mountSource)
		if err != nil {
			return nil, err
		}

		if resolvedPath == realMountSource {
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

func relatedPaths(path string) (sysPath string, resolvedSysPath string, parentPath string, childPaths []string, err error) {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", "", "", []string{}, err
	}

	devName := filepath.Base(resolvedPath)
	blockPath := filepath.Join("/", "sys", "block")
	fileInfoList, err := ioutil.ReadDir(blockPath)

	if err != nil {
		return "", "", "", []string{}, err
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
					return "", "", "", []string{}, err
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

			resolvedSysPath, err := filepath.EvalSymlinks(sysPath)
			if err != nil {
				return "", "", "", []string{}, err
			}

			return sysPath, resolvedSysPath, parentPath, childPaths, nil
		}
	}

	return "", "", "", []string{}, fmt.Errorf("sysPath, parentPath, and childPaths are not found")
}

func majorMinor(sysPath string) (string, error) {
	majorMinor, err := readFile(filepath.Join(sysPath, "dev"))
	if err != nil {
		return "", err
	}

	return majorMinor, nil
}

func osReleasePath(mountPoint string) (path string) {
	if mountPoint != "" {
		path = filepath.Join(mountPoint, "etc", "os-release")
	}

	return path
}

func newOSRelease(osReleasePath string) (osRelease map[string]string, err error) {
	osRelease = map[string]string{}

	if osReleasePath != "" {
		osReleaseLines, err := lines(osReleasePath)
		if err != nil {
			return map[string]string{}, err
		}

		for _, osReleaseLine := range osReleaseLines {
			line := strings.TrimSpace(osReleaseLine)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			kv := strings.SplitN(osReleaseLine, "=", 2)
			expectedKVSize := 2

			if len(kv) != expectedKVSize {
				return map[string]string{}, fmt.Errorf(`unexpected osReleaseLine, "%s"`, osReleaseLine)
			}

			key := kv[0]
			value := kv[1]

			osRelease[key] = trimQuotationMarks(value)
		}
	}

	return osRelease, nil
}
