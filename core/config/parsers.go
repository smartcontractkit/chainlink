package config

import (
	"fmt"
	"math/big"
	"net"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/smartcontractkit/chainlink/core/assets"
)

func ParseString(str string) (interface{}, error) {
	return str, nil
}

func ParseAddress(str string) (interface{}, error) {
	if str == "" {
		return nil, nil
	} else if common.IsHexAddress(str) {
		val := common.HexToAddress(str)
		return &val, nil
	} else if i, ok := new(big.Int).SetString(str, 10); ok {
		val := common.BigToAddress(i)
		return &val, nil
	}
	return nil, fmt.Errorf("unable to parse '%s' into EIP55-compliant address", str)
}

func ParseLink(str string) (interface{}, error) {
	i, ok := new(assets.Link).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse '%v' into *assets.Link(base 10)", str)
	}
	return i, nil
}

func ParseLogLevel(str string) (interface{}, error) {
	var lvl LogLevel
	err := lvl.Set(str)
	return lvl, err
}

func ParseUint16(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 16)
	return uint16(v), err
}

func ParseUint32(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 32)
	return uint32(v), err
}

func ParseUint64(s string) (interface{}, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err
}

func ParseF32(s string) (interface{}, error) {
	v, err := strconv.ParseFloat(s, 32)
	return v, err
}

func ParseURL(s string) (interface{}, error) {
	return url.Parse(s)
}

func ParseIP(s string) (interface{}, error) {
	return net.ParseIP(s), nil
}

func ParseDuration(s string) (interface{}, error) {
	return time.ParseDuration(s)
}

func ParseBool(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

func ParseBigInt(str string) (interface{}, error) {
	i, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return i, fmt.Errorf("unable to parse %v into *big.Int(base 10)", str)
	}
	return i, nil
}

func ParseHomeDir(str string) (interface{}, error) {
	exp, err := homedir.Expand(str)
	if err != nil {
		return nil, err
	}
	return filepath.ToSlash(exp), nil
}
