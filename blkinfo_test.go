package blkinfo

import "testing"

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
