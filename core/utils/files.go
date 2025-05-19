package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// FileExists returns true if a file at the passed string exists.
func FileExists(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed to check if file exists %q", name)
	}
	return true, nil
}

// TooPermissive checks if the file has more than the allowed permissions
func TooPermissive(fileMode, maxAllowedPerms os.FileMode) bool {
	return fileMode&^maxAllowedPerms != 0
}

// IsFileOwnedByChainlink attempts to read fileInfo to verify file owner
func IsFileOwnedByChainlink(fileInfo os.FileInfo) (bool, error) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return false, errors.Errorf("Unable to determine file owner of %s", fileInfo.Name())
	}
	return int(stat.Uid) == os.Getuid(), nil
}

// EnsureDirAndMaxPerms ensures that the given path exists, that it's a directory,
// and that it has permissions that are no more permissive than the given ones.
//
// - If the path does not exist, it is created
// - If the path exists, but is not a directory, an error is returned
// - If the path exists, and is a directory, but has the wrong perms, it is chmod'ed
func EnsureDirAndMaxPerms(path string, perms os.FileMode) error {
	stat, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		// Regular error
		return err
	} else if os.IsNotExist(err) {
		// Dir doesn't exist, create it with desired perms
		return os.MkdirAll(path, perms)
	} else if !stat.IsDir() {
		// Path exists, but it's a file, so don't clobber
		return errors.Errorf("%v already exists and is not a directory", path)
	} else if stat.Mode() != perms {
		// Dir exists, but wrong perms, so chmod
		return os.Chmod(path, stat.Mode()&perms)
	}
	return nil
}

// WriteFileWithMaxPerms writes `data` to `path` and ensures that
// the file has permissions that are no more permissive than the given ones.
func WriteFileWithMaxPerms(path string, data []byte, perms os.FileMode) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return err
	}
	defer func() { err = multierr.Combine(err, f.Close()) }()
	err = EnsureFileMaxPerms(f, perms)
	if err != nil {
		return
	}
	_, err = f.Write(data)
	return
}

// EnsureFileMaxPerms ensures that the given file has permissions
// that are no more permissive than the given ones.
func EnsureFileMaxPerms(file *os.File, perms os.FileMode) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Mode() == perms {
		return nil
	}
	return file.Chmod(stat.Mode() & perms)
}

// EnsureFilepathMaxPerms ensures that the file at the given filepath
// has permissions that are no more permissive than the given ones.
func EnsureFilepathMaxPerms(filepath string, perms os.FileMode) (err error) {
	dst, err := os.OpenFile(filepath, os.O_RDWR, perms)
	if err != nil {
		return err
	}
	defer func() { err = multierr.Combine(err, dst.Close()) }()
	return EnsureFileMaxPerms(dst, perms)
}

// FileSize repesents a file size in bytes.
type FileSize uint64

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
)

var (
	fsregex = regexp.MustCompile(`(\d+\.?\d*)(tb|gb|mb|kb|b)?`)

	fsUnitMap = map[string]int{
		"tb": TB,
		"gb": GB,
		"mb": MB,
		"kb": KB,
		"b":  1,
		"":   1,
	}
)

// MarshalText encodes s as a human readable string.
func (s FileSize) MarshalText() ([]byte, error) {
	if s >= TB {
		return []byte(fmt.Sprintf("%.2ftb", float64(s)/TB)), nil
	} else if s >= GB {
		return []byte(fmt.Sprintf("%.2fgb", float64(s)/GB)), nil
	} else if s >= MB {
		return []byte(fmt.Sprintf("%.2fmb", float64(s)/MB)), nil
	} else if s >= KB {
		return []byte(fmt.Sprintf("%.2fkb", float64(s)/KB)), nil
	}
	return []byte(fmt.Sprintf("%db", s)), nil
}

// UnmarshalText parses a file size from bs in to s.
func (s *FileSize) UnmarshalText(bs []byte) error {
	lc := strings.ToLower(strings.TrimSpace(string(bs)))
	matches := fsregex.FindAllStringSubmatch(lc, -1)
	if len(matches) != 1 || len(matches[0]) != 3 || fmt.Sprintf("%s%s", matches[0][1], matches[0][2]) != lc {
		return errors.Errorf(`bad filesize expression: "%v"`, string(bs))
	}

	var (
		num  = matches[0][1]
		unit = matches[0][2]
	)

	value, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return errors.Errorf(`bad filesize value: "%v"`, string(bs))
	}

	u, ok := fsUnitMap[unit]
	if !ok {
		return errors.Errorf(`bad filesize unit: "%v"`, unit)
	}

	*s = FileSize(value * float64(u))
	return nil
}

func (s FileSize) String() string {
	str, _ := s.MarshalText()
	return string(str)
}
