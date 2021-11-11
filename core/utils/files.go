package utils

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/pkg/errors"
)

// FileExists returns true if a file at the passed string exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
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
		return os.Chmod(path, (stat.Mode() & perms))
	}
	return nil
}

// WriteFileWithMaxPerms writes `data` to `path` and ensures that
// the file has permissions that are no more permissive than the given ones.
func WriteFileWithMaxPerms(path string, data []byte, perms os.FileMode) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return err
	}
	defer logger.ErrorIfClosing(f, "file")
	err = EnsureFileMaxPerms(f, perms)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

// CopyFileWithMaxPerms copies the file at `srcPath` to `dstPath`
// and ensures that it has permissions that are no more permissive than the given ones.
func CopyFileWithMaxPerms(srcPath, dstPath string, perms os.FileMode) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return errors.Wrap(err, "could not open source file")
	}
	defer logger.ErrorIfClosing(src, "source")

	dst, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perms)
	if err != nil {
		return errors.Wrap(err, "could not open destination file")
	}
	defer logger.ErrorIfClosing(dst, "destination")

	err = EnsureFileMaxPerms(dst, perms)
	if err != nil {
		return errors.Wrap(err, "could not set file permissions")
	}

	_, err = io.Copy(dst, src)
	return errors.Wrap(err, "could not copy file contents")
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
func EnsureFilepathMaxPerms(filepath string, perms os.FileMode) error {
	dst, err := os.OpenFile(filepath, os.O_RDWR, perms)
	if err != nil {
		return err
	}
	defer logger.ErrorIfClosing(dst, "file")

	return EnsureFileMaxPerms(dst, perms)
}

type FileSize uint64

var fsregex = regexp.MustCompile(`(\d+\.?\d*)(tb|gb|mb|kb|b)?`)

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
)

func (s FileSize) MarshalText() ([]byte, error) {
	if s > TB {
		return []byte(fmt.Sprintf("%.2ftb", float64(s)/TB)), nil
	} else if s > GB {
		return []byte(fmt.Sprintf("%.2fgb", float64(s)/GB)), nil
	} else if s > MB {
		return []byte(fmt.Sprintf("%.2fmb", float64(s)/MB)), nil
	} else if s > KB {
		return []byte(fmt.Sprintf("%.2fkb", float64(s)/KB)), nil
	}
	return []byte(fmt.Sprintf("%db", s)), nil
}

func (s *FileSize) UnmarshalText(bs []byte) error {
	matches := fsregex.FindAllStringSubmatch(strings.ToLower(string(bs)), -1)
	if len(matches) != 1 {
		return errors.Errorf(`bad filesize: "%v"`, string(bs))
	} else if len(matches[0]) != 3 {
		return errors.Errorf(`bad filesize: "%v"`, string(bs))
	}
	var (
		num  = matches[0][1]
		unit = matches[0][2]
	)
	bytes, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return err
	}

	switch unit {
	case "", "b":
	case "kb":
		bytes *= KB
	case "mb":
		bytes *= MB
	case "gb":
		bytes *= GB
	case "tb":
		bytes *= TB
	default:
		return errors.Errorf(`bad filesize unit: "%v"`, unit)
	}
	*s = FileSize(bytes)
	return nil
}

func (s FileSize) String() string {
	str, _ := s.MarshalText()
	return string(str)
}
