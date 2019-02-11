package blkinfo

import (
	"reflect"
	"testing"
)

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
