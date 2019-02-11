package blkinfo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// BlkInfo ...
type BlkInfo struct {
	RealPath     string
	Paths        []string
	Mountpoint   string
	MajorMinor   string
	RawUdevData  string
	FsUUID       string
	FsType       string
	PartEntry    *PartEntry
	RawOSRelease string
	OS           *OS
}

// PartEntry ...
type PartEntry struct {
	Scheme string
	Type   string
	Number string
}

// OS ...
type OS struct {
	Name       string
	Version    string
	ID         string
	VersionID  string
	IDLike     string
	PrettyName string
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

	b := &BlkInfo{
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

	return b, nil
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

		path := fields[0]
		mountpoint := fields[1]

		if strings.HasPrefix(path, "/dev") {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return "", err
			}

			if realPath == realPath {
				return mountpoint, nil
			}
		}
	}

	return "", nil
}

func majorMinor(realPath string) (string, error) {
	devName := filepath.Base(realPath)
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
	partEntry := &PartEntry{}
	for _, line := range strings.Split(rawUdevData, "\n") {
		prefix := "E:ID_PART_ENTRY_SCHEME="
		if strings.HasPrefix(line, prefix) {
			partEntry.Scheme = strings.TrimPrefix(line, prefix)
			continue
		}

		prefix = "E:ID_PART_ENTRY_TYPE="
		if strings.HasPrefix(line, prefix) {
			partEntry.Type = strings.TrimPrefix(line, prefix)
			continue
		}

		prefix = "E:ID_PART_ENTRY_NUMBER="
		if strings.HasPrefix(line, prefix) {
			partEntry.Number = strings.TrimPrefix(line, prefix)
		}
	}

	return partEntry
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
		return "", errors.Errorf("%s is not a file", osReleasePath)
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
		prefix := "NAME="
		if strings.HasPrefix(line, prefix) {
			os.Name = trimQuotationMarks(strings.TrimPrefix(line, prefix))
			continue
		}

		prefix = "VERSION="
		if strings.HasPrefix(line, prefix) {
			os.Version = trimQuotationMarks(strings.TrimPrefix(line, prefix))
			continue
		}

		prefix = "ID="
		if strings.HasPrefix(line, prefix) {
			os.ID = trimQuotationMarks(strings.TrimPrefix(line, prefix))
			continue
		}

		prefix = "VERSION_ID="
		if strings.HasPrefix(line, prefix) {
			os.VersionID = trimQuotationMarks(strings.TrimPrefix(line, prefix))
			continue
		}

		prefix = "ID_LIKE="
		if strings.HasPrefix(line, prefix) {
			os.IDLike = trimQuotationMarks(strings.TrimPrefix(line, prefix))
			continue
		}

		prefix = "PRETTY_NAME="
		if strings.HasPrefix(line, prefix) {
			os.PrettyName = trimQuotationMarks(strings.TrimPrefix(line, prefix))
			continue
		}
	}

	return os
}
