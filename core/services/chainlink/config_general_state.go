package chainlink

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

func (g *generalConfig) AppID() uuid.UUID {
	g.appIDOnce.Do(func() {
		g.appID = uuid.NewV4()
	})
	return g.appID
}

func (g *generalConfig) DefaultLogLevel() zapcore.Level {
	return g.logLevelDefault
}

func (g *generalConfig) LogLevel() (ll zapcore.Level) {
	g.logMu.RLock()
	ll = g.logLevelDefault
	g.logMu.RUnlock()
	return
}

func (g *generalConfig) SetLogLevel(lvl zapcore.Level) error {
	g.logMu.Lock()
	g.logLevel = lvl
	g.logMu.Unlock()
	return nil
}

func (g *generalConfig) LogSQL() (sql bool) {
	g.logMu.RLock()
	sql = g.logSQL
	g.logMu.RUnlock()
	return
}

func (g *generalConfig) SetLogSQL(logSQL bool) {
	g.logMu.Lock()
	g.logSQL = logSQL
	g.logMu.Unlock()
}
