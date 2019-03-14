package blkinfo

import (
	"fmt"
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
