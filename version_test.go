package blkinfo

import "testing"

func TestVersion(t *testing.T) {
	tests := []struct {
		want string
	}{
		{
			want: "dev",
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
			want: "none",
		},
	}

	for _, tt := range tests {
		got := Revision()
		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}

func TestBuiltAt(t *testing.T) {
	tests := []struct {
		want string
	}{
		{
			want: "unknown",
		},
	}

	for _, tt := range tests {
		got := BuiltAt()
		if got != tt.want {
			t.Errorf("got: %v, tt.want: %v", got, tt.want)
		}
	}
}
