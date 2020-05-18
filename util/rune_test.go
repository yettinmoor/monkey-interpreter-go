package util

import "testing"

func TestRuneFuncs(t *testing.T) {
	tests := []string{
		"abc123",
		"héllo",
		"çëłłß",
		"日本語",
	}
	for _, tt := range tests {
		if !IsValidIdentifier(tt) {
			t.Errorf("%s failed", tt)
		}
	}
}
