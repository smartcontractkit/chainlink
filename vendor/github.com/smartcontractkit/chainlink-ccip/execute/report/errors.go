package report

import "errors"

var ErrEmptyReport = errors.New("no messages can fit in the report")
var ErrNotReady = errors.New("token data not ready")
