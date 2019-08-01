package orm

import (
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type stringUnmarshaler string

func (s *stringUnmarshaler) UnmarshalText(text []byte) error {
	*s = stringUnmarshaler(text)
	return nil
}

type durationUnmarshaler time.Duration

func (d *durationUnmarshaler) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	*d = durationUnmarshaler(duration)
	return err
}

type boolUnmarshaler bool

func (b *boolUnmarshaler) UnmarshalText(text []byte) error {
	bl, err := strconv.ParseBool(string(text))
	*b = boolUnmarshaler(bl)
	return err
}

type logLevelUnmarshaler LogLevel

func (l *logLevelUnmarshaler) UnmarshalText(text []byte) error {
	var lvl LogLevel
	err := lvl.Set(string(text))
	*l = logLevelUnmarshaler(lvl)
	return err
}

type uint16Unmarshaler uint16

func (u *uint16Unmarshaler) UnmarshalText(text []byte) error {
	d, err := strconv.ParseUint(string(text), 10, 16)
	*u = uint16Unmarshaler(d)
	return err
}

type uint64Unmarshaler uint16

func (u *uint64Unmarshaler) UnmarshalText(text []byte) error {
	d, err := strconv.ParseUint(string(text), 10, 64)
	*u = uint64Unmarshaler(d)
	return err
}

type urlUnmarshaler url.URL

func (u *urlUnmarshaler) UnmarshalText(text []byte) error {
	rl, err := url.Parse(string(text))
	*u = urlUnmarshaler(*rl)
	return err
}

type addressUnmarshaler common.Address

func (u *addressUnmarshaler) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" {
		return nil
	} else if common.IsHexAddress(str) {
		val := common.HexToAddress(str)
		*u = addressUnmarshaler(val)
		return nil
	} else if i, ok := new(big.Int).SetString(str, 10); ok {
		val := common.BigToAddress(i)
		*u = addressUnmarshaler(val)
		return nil
	}

	return fmt.Errorf("Unable to parse '%s' into EIP55-compliant address", str)
}

//func parseAddress(str string) (interface{}, error) {
//if str == "" {
//return nil, nil
//} else if common.IsHexAddress(str) {
//val := common.HexToAddress(str)
//return &val, nil
//} else if i, ok := new(big.Int).SetString(str, 10); ok {
//val := common.BigToAddress(i)
//return &val, nil
//}
//return nil, fmt.Errorf("Unable to parse '%s' into EIP55-compliant address", str)
//}

//func parseLink(str string) (interface{}, error) {
//i, ok := new(assets.Link).SetString(str, 10)
//if !ok {
//return i, fmt.Errorf("Unable to parse '%v' into *assets.Link(base 10)", str)
//}
//return i, nil
//}

//func parseLogLevel(str string) (interface{}, error) {
//var lvl LogLevel
//err := lvl.Set(str)
//return lvl, err
//}

//func parsePort(str string) (interface{}, error) {
//d, err := strconv.ParseUint(str, 10, 16)
//return uint16(d), err
//}

//func parseURL(s string) (interface{}, error) {
//return url.Parse(s)
//}

//func parseBigInt(str string) (interface{}, error) {
//i, ok := new(big.Int).SetString(str, 10)
//if !ok {
//return i, fmt.Errorf("Unable to parse %v into *big.Int(base 10)", str)
//}
//return i, nil
//}

//func parseHomeDir(str string) (interface{}, error) {
//exp, err := homedir.Expand(str)
//if err != nil {
//return nil, err
//}
//return filepath.ToSlash(exp), nil
//}
