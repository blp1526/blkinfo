package blkinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"errors"
)

// BlkInfo ...
type BlkInfo struct {
	RealPath     string     `json:"real_path"      yaml:"real_path"     `
	Paths        []string   `json:"paths"          yaml:"paths"         `
	Mountpoint   string     `json:"mountpoint"     yaml:"mountpoint"    `
	MajorMinor   string     `json:"major_minor"    yaml:"major_minor"   `
	RawUdevData  string     `json:"raw_udev_data"  yaml:"raw_udev_data" `
	FsUUID       string     `json:"fs_uuid"        yaml:"fs_uuid"       `
	FsType       string     `json:"fs_type"        yaml:"fs_type"       `
	PartEntry    *PartEntry `json:"part_entry"     yaml:"part_entry"    `
	RawOSRelease string     `json:"raw_os_release" yaml:"raw_os_release"`
	OS           *OS        `json:"os"             yaml:"os"            `
}

// PartEntry ...
type PartEntry struct {
	Scheme string `json:"scheme" yaml:"scheme"`
	Type   string `json:"type"   yaml:"type"  `
	Number string `json:"number" yaml:"number"`
}

// OS ...
type OS struct {
	Name       string `json:"name"        yaml:"name"       `
	Version    string `json:"version"     yaml:"version"    `
	ID         string `json:"id"          yaml:"id"         `
	VersionID  string `json:"version_id"  yaml:"version_id" `
	IDLike     string `json:"id_like"     yaml:"id_like"    `
	PrettyName string `json:"pretty_name" yaml:"pretty_name"`
}

// New ...
func New(devPath string) (*BlkInfo, error) {
	realPath, err := filepath.EvalSymlinks(devPath)
	if err != nil {
		return nil, err
	}

	mtab, err := mtab()
	if err != nil {
		return nil, err
	}

	mountpoint, err := mountpoint(mtab, realPath)
	if err != nil {
		return nil, err
	}

	majorMinor, err := majorMinor(realPath)
	if err != nil {
		return nil, err
	}

	rawUdevData, err := rawUdevData(majorMinor)
	if err != nil {
		return nil, err
	}

	rawOSRelease, err := rawOSRelease(mountpoint)
	if err != nil {
		return nil, err
	}

	bi := &BlkInfo{
		RealPath:     realPath,
		Mountpoint:   mountpoint,
		MajorMinor:   majorMinor,
		RawUdevData:  rawUdevData,
		Paths:        paths(rawUdevData),
		FsUUID:       fsUUID(rawUdevData),
		FsType:       fsType(rawUdevData),
		PartEntry:    newPartEntry(rawUdevData),
		RawOSRelease: rawOSRelease,
		OS:           newOS(rawOSRelease),
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
		s = strings.TrimPrefix(s, q)
		s = strings.TrimSuffix(s, q)
	}

	return s
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

func majorMinor(realPath string) (string, error) {
	baseName := filepath.Base(realPath)
	sysBlockPath := filepath.Join("/", "sys", "block")
	fileInfos, err := ioutil.ReadDir(sysBlockPath)
	if err != nil {
		return "", err
	}

	numberPath := ""
	for _, fileInfo := range fileInfos {
		fileInfoName := fileInfo.Name()
		if strings.HasPrefix(baseName, fileInfoName) {
			numberPath = filepath.Join(sysBlockPath, fileInfoName)
			if baseName != fileInfoName {
				// name is a partition.
				numberPath = filepath.Join(numberPath, baseName)
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

func rawUdevData(majorMinor string) (string, error) {
	rawUdevData, err := readFile(filepath.Join("/", "run", "udev", "data", "b"+majorMinor))
	if err != nil {
		return "", err
	}

	return rawUdevData, nil
}

func paths(rawUdevData string) []string {
	paths := []string{}

	for _, line := range strings.Split(rawUdevData, "\n") {
		prefix := "S:"
		if strings.HasPrefix(line, prefix) {
			paths = append(paths, filepath.Join("/dev", strings.TrimPrefix(line, prefix)))
		}
	}

	return paths
}

func fsUUID(rawUdevData string) string {
	for _, line := range strings.Split(rawUdevData, "\n") {
		prefix := "E:ID_FS_UUID="
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix)
		}
	}

	return ""
}

func fsType(rawUdevData string) string {
	for _, line := range strings.Split(rawUdevData, "\n") {
		prefix := "E:ID_FS_TYPE="
		if strings.HasPrefix(line, prefix) {
			return strings.TrimPrefix(line, prefix)
		}
	}

	return ""
}

func newPartEntry(rawUdevData string) *PartEntry {
	pe := &PartEntry{}
	for _, line := range strings.Split(rawUdevData, "\n") {
		if strings.HasPrefix(line, "E:ID_PART_ENTRY") {
			s := strings.SplitN(line, "=", 2)
			key := s[0]
			value := trimQuotationMarks(s[1])

			switch key {
			case "E:ID_PART_ENTRY_SCHEME":
				pe.Scheme = value
			case "E:ID_PART_ENTRY_TYPE":
				pe.Type = value
			case "E:ID_PART_ENTRY_NUMBER":
				pe.Number = value
			}
		}
	}

	return pe
}

func rawOSRelease(mountpoint string) (string, error) {
	if mountpoint == "" {
		return "", nil
	}

	osReleasePath := filepath.Join(mountpoint, "etc", "os-release")
	fileInfo, err := os.Stat(osReleasePath)
	if err == os.ErrNotExist {
		return "", nil
	}

	if err != nil && err != os.ErrNotExist {
		return "", err
	}

	if fileInfo.IsDir() {
		return "", fmt.Errorf("%s is not a file", osReleasePath)
	}

	rawOSRelease, err := readFile(osReleasePath)
	if err != nil {
		return "", err
	}

	return rawOSRelease, nil
}

func newOS(rawOSRelease string) *OS {
	os := &OS{}
	for _, line := range strings.Split(rawOSRelease, "\n") {
		s := strings.SplitN(line, "=", 2)
		key := s[0]
		value := trimQuotationMarks(s[1])

		switch key {
		case "NAME":
			os.Name = value
		case "VERSION":
			os.Version = value
		case "ID":
			os.ID = value
		case "VERSION_ID":
			os.VersionID = value
		case "ID_LIKE":
			os.IDLike = value
		case "PRETTY_NAME":
			os.PrettyName = value
		}
	}

	return os
}
