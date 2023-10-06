package utils

import (
	"path/filepath"
	"testing"
)

func TestGetCWD(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "cwd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("cwd: %+v", GetCWD())
		})
	}
}

func TestAbsPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "abs path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := filepath.Abs("test")
			t.Logf("abs path: %+v, err: %+v", path, err)
		})
	}
}
