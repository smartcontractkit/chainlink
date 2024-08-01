package telemetry

import (
	"fmt"
	"log"
)

const (
	ServiceName    = "automation-ocr3"
	LogPkgStdFlags = log.Lshortfile
)

func WrapLogger(logger *log.Logger, ns string) *log.Logger {
	return log.New(logger.Writer(), fmt.Sprintf("[%s | %s]", ServiceName, ns), LogPkgStdFlags)
}
