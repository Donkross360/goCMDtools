package main

import (
	"bytes"
	"strings"
	"testing"
)

// to normalize paths before comparison for windows
func normalizePath(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}
func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},

		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},

		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},

		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 20, list: true}, expected: ""},

		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".gz", size: 0, list: true}, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer

			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			res := buffer.String()

			if tc.expected != normalizePath(res) {
				t.Errorf("expected %q, got %q instaed\n", tc.expected, res)
			}
		})
	}
}
