package logger

import (
	"go.uber.org/zap"
	"log"
)

var logger *Logger

type Logger struct {
	*zap.Logger
}

func init() {
	config := zap.NewProductionConfig()
	zap, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger = &Logger{zap}
}

func ForGin() *Logger {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	zap, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}
	return &Logger{zap}
}

func (self *Logger) Write(b []byte) (n int, err error) {
	self.Info(string(b))
	return len(b), nil
}

func Get() *Logger {
	return logger
}
