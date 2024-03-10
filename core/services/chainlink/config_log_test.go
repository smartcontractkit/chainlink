package chainlink

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestLogConfig(t *testing.T) {
	opts := GeneralConfigOpts{
		ConfigStrings: []string{fullTOML},
	}
	cfg, err := opts.New()
	require.NoError(t, err)

	log := cfg.Log()
	file := log.File()

	assert.Equal(t, "log/file/dir", file.Dir())
	assert.Equal(t, uint64(100*utils.GB), uint64(file.MaxSize()))
	assert.Equal(t, int64(17), file.MaxAgeDays())
	assert.Equal(t, int64(9), file.MaxBackups())
	assert.Equal(t, true, log.UnixTimestamps())
	assert.Equal(t, true, log.JSONConsole())
	assert.Equal(t, zapcore.Level(3), log.DefaultLevel())
	assert.Equal(t, zapcore.Level(3), log.Level())
}
