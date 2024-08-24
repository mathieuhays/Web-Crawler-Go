package main

import (
	"errors"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expected    string
		expectedErr error
	}{
		{
			name:        "remove scheme",
			inputURL:    "https://blog.boot.dev/path",
			expected:    "blog.boot.dev/path",
			expectedErr: nil,
		},
		{
			name:        "remove trailing slash",
			inputURL:    "http://blog.boot.dev/",
			expected:    "blog.boot.dev",
			expectedErr: nil,
		},
		{
			name:        "consistent case",
			inputURL:    "http://blog.Boot.dev/",
			expected:    "blog.boot.dev",
			expectedErr: nil,
		},
		{
			name:        "invalid url",
			inputURL:    "hey",
			expected:    "",
			expectedErr: errInvalidURL,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
