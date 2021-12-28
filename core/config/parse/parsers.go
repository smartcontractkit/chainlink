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

func String(str string) (interface{}, error) {
	return str, nil
}

func Link(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse '%v' into *assets.Link(base 10)", str)
	}
	return i, nil
}

func LogLevel(str string) (interface{}, error) {
	var lvl zapcore.Level
	err := lvl.Set(str)
	return lvl, err
}

func Uint16(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

func Uint32(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func Uint64(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err
}

func Int64(s string) (interface{}, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err
}

func F32(s string) (interface{}, error) {
	v, err := strconv.ParseFloat(s, 32)
	return float32(v), err
}

func URL(s string) (interface{}, error) {
	return url.Parse(s)
}

func IP(s string) (interface{}, error) {
	return net.ParseIP(s), nil
}

func Duration(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

func FileSize(s string) (interface{}, error) {
	var fs utils.FileSize
	err := fs.UnmarshalText([]byte(s))
	return fs, err
}

func Bool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

func BigInt(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse %v into *big.Int(base 10)", str)
	}
	return i, nil
}

func HomeDir(str string) (interface{}, error) {
	exp, err := homedir.Expand(str)
	if err != nil {
		return nil, err
	}
	return filepath.ToSlash(exp), nil
}
