package vault

import (
	"errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func newUserError(msg string) *userError {
	return &userError{msg: msg}
}

type userError struct {
	msg string
}

func (u *userError) Error() string {
	return u.msg
}

func (u *userError) Is(target error) bool {
	_, ok := target.(*userError)
	return ok
}

func userFacingError(err error, fallback string) string {
	if errors.Is(err, &userError{}) {
		return err.Error()
	}

	return fallback
}

func logUserErrorAware(l logger.Logger, msg string, err error, keysAndValues ...interface{}) {
	keysAndValues = append(keysAndValues, "error", err)
	lggr := l.Helper(1)
	if errors.Is(err, &userError{}) {
		lggr.Debugw(msg, keysAndValues...)
		return
	}

	lggr.Errorw(msg, keysAndValues...)
}
