package chainlink

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

func (g *generalConfig) AppID() uuid.UUID {
	g.appIDOnce.Do(func() {
		if g.c.AppID != (uuid.UUID{}) {
			return // already set (e.g. test override)
		}
		g.c.AppID = uuid.NewV4() // randomize
	})
	return g.c.AppID
}

func (g *generalConfig) DefaultLogLevel() zapcore.Level {
	return g.logLevelDefault
}

func (g *generalConfig) LogLevel() (ll zapcore.Level) {
	g.logMu.RLock()
	ll = g.c.Log.Level
	g.logMu.RUnlock()
	return
}

func (g *generalConfig) SetLogLevel(lvl zapcore.Level) error {
	g.logMu.Lock()
	g.c.Log.Level = lvl
	g.logMu.Unlock()
	return nil
}

func (g *generalConfig) LogSQL() (sql bool) {
	g.logMu.RLock()
	sql = g.c.Log.SQL
	g.logMu.RUnlock()
	return
}

func (g *generalConfig) SetLogSQL(logSQL bool) {
	g.logMu.Lock()
	g.c.Log.SQL = logSQL
	g.logMu.Unlock()
}
