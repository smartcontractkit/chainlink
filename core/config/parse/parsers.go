package parse

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// String parser
func String(str string) (interface{}, error) {
	return str, nil
}

// Link parser
func Link(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse '%v' into *assets.Link(base 10)", str)
	}
	return i, nil
}

// LogLevel sets log level
func LogLevel(str string) (interface{}, error) {
	var lvl zapcore.Level
	err := lvl.Set(str)
	return lvl, err
}

// Uint16 converts string to uint16
func Uint16(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

// Uint32 converts string to uint32
func Uint32(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

// Uint64 converts string to uint64
func Uint64(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err
}

// Int64 converts string to int64
func Int64(s string) (interface{}, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err
}

// F32 converts string to float32
func F32(s string) (interface{}, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

// URL converts string to parsed URL type
func URL(s string) (interface{}, error) {
	return url.Parse(s)
}

// IP converts string to parsed IP type
func IP(s string) (interface{}, error) {
	return net.ParseIP(s), nil
}

// Duration converts string to parsed Duratin type
func Duration(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

// FileSize parses string as FileSize type
func FileSize(s string) (interface{}, error) {
	var fs utils.FileSize
	err := fs.UnmarshalText([]byte(s))
	return fs, err
}

// Bool parses string as a bool type
func Bool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

// BigInt parses string into a big int
func BigInt(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse %v into *big.Int(base 10)", str)
	}
	return i, nil
}

// HomeDir parses string as a file path
func HomeDir(str string) (interface{}, error) {
	exp, err := homedir.Expand(str)
	if err != nil {
		return nil, err
	}
	return filepath.ToSlash(exp), nil
}
