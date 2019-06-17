package blkinfo

import "testing"

func TestVersion(t *testing.T) {
	tests := []struct {
		want string
	}{
		{
			want: "0.0.1",
		},
	}

	for _, tt := range tests {
		got := Version()
		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}

func TestRevision(t *testing.T) {
	tests := []struct {
		want string
	}{
		{
			want: "devel",
		},
	}

	for _, tt := range tests {
		got := Revision()
		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}
