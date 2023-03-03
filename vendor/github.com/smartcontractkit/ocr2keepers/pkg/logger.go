package ocr2keepers

import "github.com/smartcontractkit/libocr/commontypes"

type logWriter struct {
	l commontypes.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.l.Debug(string(p), nil)
	n = len(p)
	return
}
