package util

import (
	"fmt"
	"strings"
)

func WrapErrorf(err error, fmtString string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(fmtString+": %s", append(args, err)...)
}

func WrapError(err error, msg string) error {
	return WrapErrorf(err, strings.ReplaceAll(msg, "%", "%%"))
}
