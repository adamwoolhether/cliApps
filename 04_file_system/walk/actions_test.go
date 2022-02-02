package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name      string
		file      string
		ext       string
		minSize   int64
		nameMatch string
		expected  bool
	}{
		{"FilterNoExtension", "testdata/dir.log", "", 0, "", false},
		{"FilterExtensionMatch", "testdata/dir.log", ".log", 0, "", false},
		{"FilterExtensionNoMatch", "testdata/dir.log", ".sh", 0, "", true},
		{"FilterExtensionSizeMatch", "testdata/dir.log", ".log", 10, "", false},
		{"FilterExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, "", true},
		{"FilterMultiExtensionMatch", "testdata/dir.log", ".log,.sh", 0, "", false},
		{"FilterMultiExtensionNoMatch", "testdata/dir.log", ".zip,.exe", 0, "", true},
		{"FilterNoNameMatch", "testdata/dir.log", "", 0, "", false},
		{"FilterNameNoMatch", "testdata/dir.log", "", 0, "go*.sh", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tc.file, tc.ext, tc.minSize, 0, tc.nameMatch, info)
			if f != tc.expected {
				t.Errorf("Expected '%t', got '%t'\n", tc.expected, f)
			}
		})
	}
}
