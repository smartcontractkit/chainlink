// +build windows

package orm

import (
	"net/url"
	"path/filepath"
)

func getLockPath(dbpath string) (string, error) {
	uri, err := url.Parse(dbpath)
	if err != nil {
		return "", err
	}

	// Remove leading slash left by url.Parse() for windows
	dbpath = filepath.ToSlash(filepath.Clean(uri.Path))
	if leadingSlash(dbpath) {
		dbpath = uri.Path[1:]
	}
	directory := filepath.Dir(dbpath)
	return filepath.Join(directory, "chainlink.lock"), nil
}

func leadingSlash(path string) bool {
	return path[0] == '/' || path[0] == '\\'
}
