package chainlink

import (
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap/zapcore"
)

func (l *legacyGeneralConfig) AppID() uuid.UUID {
	l.appIDOnce.Do(func() {
		l.appID = uuid.NewV4()
	})
	return l.appID
}

func (l *legacyGeneralConfig) DefaultLogLevel() zapcore.Level {
	return l.logLevelDefault
}

func (l *legacyGeneralConfig) LogLevel() (ll zapcore.Level) {
	l.logMu.RLock()
	ll = l.logLevelDefault
	l.logMu.RUnlock()
	return
}

func (l *legacyGeneralConfig) SetLogLevel(lvl zapcore.Level) error {
	l.logMu.Lock()
	l.logLevel = lvl
	l.logMu.Unlock()
	return nil
}

func (l *legacyGeneralConfig) LogSQL() (sql bool) {
	l.logMu.RLock()
	sql = l.logSQL
	l.logMu.RUnlock()
	return
}

func (l *legacyGeneralConfig) SetLogSQL(logSQL bool) {
	l.logMu.Lock()
	l.logSQL = logSQL
	l.logMu.Unlock()
}
