package chainlink

import (
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var _ config.Log = (*logConfig)(nil)

type logConfig struct {
	c            toml.Log
	rootDir      func() string
	defaultLevel zapcore.Level
	level        func() zapcore.Level
}

type fileConfig struct {
	c       toml.LogFile
	rootDir func() string
}

func (f *fileConfig) Dir() string {
	s := *f.c.Dir
	if s == "" {
		s = f.rootDir()
	}
	return s
}

func (f *fileConfig) MaxSize() utils.FileSize {
	return *f.c.MaxSize
}

func (f *fileConfig) MaxAgeDays() int64 {
	return *f.c.MaxAgeDays
}

func (f *fileConfig) MaxBackups() int64 {
	return *f.c.MaxBackups
}

func (l *logConfig) File() config.File {
	return &fileConfig{c: l.c.File, rootDir: l.rootDir}
}

func (l *logConfig) UnixTimestamps() bool {
	return *l.c.UnixTS
}

func (l *logConfig) JSONConsole() bool {
	return *l.c.JSONConsole
}

func (l *logConfig) DefaultLevel() zapcore.Level {
	return l.defaultLevel
}

func (l *logConfig) Level() zapcore.Level {
	return l.level()
}
