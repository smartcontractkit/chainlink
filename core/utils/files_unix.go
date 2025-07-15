//go:build unix

package utils

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

// IsFileOwnedByChainlink attempts to read fileInfo to verify file owner
func IsFileOwnedByChainlink(fileInfo os.FileInfo) (bool, error) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return false, errors.Errorf("Unable to determine file owner of %s", fileInfo.Name())
	}
	return int(stat.Uid) == os.Getuid(), nil
}
