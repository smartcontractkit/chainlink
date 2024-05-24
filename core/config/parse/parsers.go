package parse

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/go-homedir"
	pkgerrors "github.com/pkg/errors"
	"go.uber.org/zap/zapcore"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func String(str string) (string, error) {
	return str, nil
}

func Link(str string) (*commonassets.Link, error) {
	i, ok := new(commonassets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse '%s'", str)
	}
	return i, nil
}

func LogLevel(str string) (zapcore.Level, error) {
	var lvl zapcore.Level
	err := lvl.Set(str)
	return lvl, err
}

func Uint16(s string) (uint16, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

func Uint32(s string) (uint32, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func Uint64(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err
}

func Int64(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err
}

func F32(s string) (float32, error) {
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

func FileSize(s string) (utils.FileSize, error) {
	var fs utils.FileSize
	err := fs.UnmarshalText([]byte(s))
	return fs, err
}

// Bool parses string as a bool type
func Bool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

func BigInt(str string) (*big.Int, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse %v into *big.Int(base 10)", str)
	}
	return i, nil
}

func Wei(str string) (w *assets.Wei, err error) {
	w = new(assets.Wei)
	err = w.UnmarshalText([]byte(str))
	return w, err
}

func HomeDir(str string) (string, error) {
	exp, err := homedir.Expand(str)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(exp), nil
}

func DatabaseURL(s string) (url.URL, error) {
	uri, err := url.Parse(s)
	if err != nil {
		return url.URL{}, pkgerrors.Wrapf(err, "invalid database url %s", s)
	}
	if uri.String() == "" {
		return *uri, nil
	}
	static.SetConsumerName(uri, "Default", nil)
	return *uri, nil
}
