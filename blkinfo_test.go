package blkinfo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		path string
		want string
		err  bool
	}{
		{
			path: "/no/such/file",
			want: "",
			err:  true,
		},
		{
			path: "LICENSE",
			want: `Copyright (c) 2019 Shingo Kawamura

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.`,
			err: false,
		},
	}

	for _, tt := range tests {
		got, err := readFile(tt.path)
		errMsg := fmt.Sprintf("tt: %+v, got: %v, err: %v", tt, got, err)
		if tt.err && err == nil {
			t.Errorf(errMsg)
		}

		if !tt.err && err != nil {
			t.Errorf(errMsg)
		}

		if got != tt.want {
			t.Errorf(errMsg)
		}
	}
}

func TestTrimQuottionMarks(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{
			s:    `foo`,
			want: `foo`,
		},
		{
			s:    `"foo"`,
			want: `foo`,
		},
		{
			s:    `'foo'`,
			want: `foo`,
		},
		{
			s:    `"'foo'"`,
			want: `'foo'`,
		},
		{
			s:    `'"foo"'`,
			want: `"foo"`,
		},
		{
			s:    `"foo`,
			want: `"foo`,
		},
		{
			s:    `'foo`,
			want: `'foo`,
		},
		{
			s:    `foo"`,
			want: `foo"`,
		},
		{
			s:    `foo'`,
			want: `foo'`,
		},
	}

	for _, tt := range tests {
		got := trimQuotationMarks(tt.s)
		if got != tt.want {
			t.Errorf("tt: %+v, got: %v", tt, got)
		}
	}
}

func TestFsUUID(t *testing.T) {
	tests := []struct {
		rawUdevData string
		want        string
	}{
		{
			rawUdevData: "",
			want:        "",
		},
		{
			rawUdevData: "E:ID_FS_UUID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			want:        "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
		},
	}

	for _, tt := range tests {
		got := fsUUID(tt.rawUdevData)
		if got != tt.want {
			t.Errorf("tt: %#v, got: %v", tt, got)
		}
	}
}

func TestFsType(t *testing.T) {
	tests := []struct {
		rawUdevData string
		want        string
	}{
		{
			rawUdevData: "",
			want:        "",
		},
		{
			rawUdevData: "E:ID_FS_TYPE=ext4",
			want:        "ext4",
		},
	}

	for _, tt := range tests {
		got := fsType(tt.rawUdevData)
		if got != tt.want {
			t.Errorf("tt: %#v, got: %v", tt, got)
		}
	}
}

func TestPaths(t *testing.T) {
	tests := []struct {
		rawUdevData string
		want        []string
	}{
		{
			rawUdevData: ``,
			want:        []string{},
		},
		{
			rawUdevData: `S:disk/by-path/pci-0000:00:00.0-scsi-0:0:0:0-part1
S:disk/by-partuuid/xxxxxxxx-xx
S:disk/by-uuid/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`,
			want: []string{
				"/dev/disk/by-path/pci-0000:00:00.0-scsi-0:0:0:0-part1",
				"/dev/disk/by-partuuid/xxxxxxxx-xx",
				"/dev/disk/by-uuid/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			},
		},
	}

	for _, tt := range tests {
		got := paths(tt.rawUdevData)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("tt: %#v, got: %v", tt, got)
		}
	}
}

func TestNewPartTable(t *testing.T) {
	tests := []struct {
		rawUdevData string
		want        *PartTable
	}{
		{
			rawUdevData: ``,
			want:        &PartTable{},
		},
		{
			rawUdevData: `
E:ID_PART_TABLE_TYPE=dos
E:ID_PART_TABLE_UUID=xxxxxxxx
`,
			want: &PartTable{
				Type: "dos",
				UUID: "xxxxxxxx",
			},
		},
	}
	for _, tt := range tests {
		got := newPartTable(tt.rawUdevData)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %+v, tt.want: %+v", got, tt.want)
		}
	}
}

func TestNewPartEntry(t *testing.T) {
	tests := []struct {
		rawUdevData string
		want        *PartEntry
	}{
		{
			rawUdevData: ``,
			want:        &PartEntry{},
		},
		{
			rawUdevData: `
E:ID_PART_ENTRY_SCHEME=dos
E:ID_PART_ENTRY_UUID=xxxxxxxx-xx
E:ID_PART_ENTRY_TYPE=0x83
E:ID_PART_ENTRY_FLAGS=0x80
E:ID_PART_ENTRY_NUMBER=1
E:ID_PART_ENTRY_OFFSET=2048
E:ID_PART_ENTRY_SIZE=209711104
E:ID_PART_ENTRY_DISK=8:0
`,
			want: &PartEntry{
				Scheme: "dos",
				Type:   "0x83",
				Number: "1",
			},
		},
	}
	for _, tt := range tests {
		got := newPartEntry(tt.rawUdevData)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %+v, tt.want: %+v", got, tt.want)
		}
	}
}

func TestNewOS(t *testing.T) {
	tests := []struct {
		rawOSRelease string
		want         *OS
	}{
		{
			rawOSRelease: ``,
			want:         &OS{},
		},
		{
			rawOSRelease: `
NAME="Ubuntu"
VERSION="18.04.2 LTS (Bionic Beaver)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 18.04.2 LTS"
VERSION_ID="18.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
VERSION_CODENAME=bionic
UBUNTU_CODENAME=bionic
`,
			want: &OS{
				Name:       "Ubuntu",
				Version:    "18.04.2 LTS (Bionic Beaver)",
				ID:         "ubuntu",
				VersionID:  "18.04",
				IDLike:     "debian",
				PrettyName: "Ubuntu 18.04.2 LTS",
			},
		},
	}
	for _, tt := range tests {
		got := newOS(tt.rawOSRelease)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got: %+v, tt.want: %+v", got, tt.want)
		}
	}
}
