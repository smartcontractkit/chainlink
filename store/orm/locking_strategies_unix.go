// +build !windows

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
	dbpath = uri.Path
	directory := filepath.Dir(dbpath)
	return filepath.Join(directory, "chainlink.lock"), nil
}
