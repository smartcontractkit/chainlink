//go:build windows

package utils

import (
	"os"

	"github.com/pkg/errors"
)

func IsFileOwnedByChainlink(fileInfo os.FileInfo) (bool, error) {
	// On Windows, we cannot reliably determine the owner of a file
	// using the same method as Unix. Instead, we assume that if the
	// file exists and is accessible, it is owned by the current user.
	if fileInfo == nil {
		return false, errors.New("fileInfo is nil")
	}
	return true, nil // Assume ownership for simplicity
}
