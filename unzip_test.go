package unzip

import (
	"os"
	"testing"
)

func Test_unzip(t *testing.T) {
	tests := []struct {
		src        string
		dest       string
		files      []string
		shouldPass bool
	}{
		{"", "", []string{}, false},
		{"testdata/test.zip", "testdata/temp", []string{"testdata/temp/test.txt"}, true},
	}
	for _, test := range tests {
		files, err := Unzip(test.src, test.dest)
		if test.shouldPass == (err != nil) {
			t.Fatal(err)
		}
		defer os.RemoveAll(test.dest)

		for _, f := range files {
			for _, expected := range test.files {
				if f != expected {
					t.Fatal("Files did not match expected")
				}
			}
		}
	}
}

func Test_downloadFile(t *testing.T) {
	tests := []struct {
		url        string
		filepath   string
		shouldPass bool
	}{
		{"", "", false},
		{"http://google.com/", "", true},
		{"http://test.fg", "somefile.xlsx", false},
		{"", "somefile.txt", false},
	}
	for i, test := range tests {
		_, err := DownloadFile(test.url, test.filepath)
		if test.shouldPass == (err != nil) {
			t.Fatalf("Test %d Error: %v. Should Pass: %t, Passed: %t", i, err, test.shouldPass, err == nil)
		}
	}
}
