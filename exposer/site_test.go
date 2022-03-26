package exposer

import (
	"testing"
)

func TestNewFileName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "aaa/bbb", want: "aaa_bbb"},
		{input: "aaa/bbb/ccc", want: "aaa_bbb_ccc"},
		{input: "ああ/", want: "ああ_"},
	}

	for _, tt := range tests {
		got := newFileName(tt.input)
		if got != tt.want {
			t.Errorf("want: %s, got: %s", tt.want, got)
		}
	}
}
