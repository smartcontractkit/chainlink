package logger

import (
	"log"
	"os"
)

type PrettyConsole struct {
	io *os.File
}

func (PrettyConsole) Sync() error  { return nil }
func (PrettyConsole) Close() error { return nil }

func (pc PrettyConsole) Write(p []byte) (n int, err error) {
	log.Print(string(p))
	return len(p), nil
}
