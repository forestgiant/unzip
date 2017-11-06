package unzip

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Unzip function to handle zipped image directories
func Unzip(src, dest string) ([]string, error) {
	var files []string

	// Try to open the zip
	r, err := zip.OpenReader(src)
	if err != nil {
		return files, err
	}
	defer r.Close()

	// Range over all files in zip to
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return files, err
		}

		path := filepath.Join(dest, f.Name)
		files = append(files, path)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
		} else {
			var dir string
			if lastIndex := strings.LastIndex(path, string(os.PathSeparator)); lastIndex > -1 {
				dir = path[:lastIndex]
			}

			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				rc.Close()
				return files, err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				rc.Close()
				return files, err
			}

			_, err = io.Copy(f, rc)
			if err != nil {
				rc.Close()
				f.Close()
				return files, err
			}

			f.Close()
		}

		rc.Close()
	}
	return files, nil
}

// DownloadFile takes a url and dest filepath and returns a complete filepath to the location of the file
func DownloadFile(rawurl string, dest string) (filepath string, err error) {
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Failed to download file")
	}

	// Save the file to the specified filename or create a tempfile to save to
	var f *os.File
	if dest != "" {
		f, err = os.Create(dest)
		if err != nil {
			return "", err
		}
	} else {
		f, err = ioutil.TempFile("", "unzip")
		if err != nil {
			return "", err
		}
		dest = f.Name()
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	return dest, nil
}
