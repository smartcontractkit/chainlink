package logger

import (
	"time"
)

// LogConfig stores key value pairs for configuring package specific logging
type LogConfig struct {
	ID          int64
	ServiceName string
	LogLevel    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
