package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
)

var ErrUnhealthy = errors.New("Unhealthy")

type boolCheck struct {
	name    string
	healthy bool
}

func (b boolCheck) Name() string { return b.name }

func (b boolCheck) Ready() error {
	if b.healthy {
		return nil
	}
	return errors.New("Not ready")
}

func (b boolCheck) HealthReport() map[string]error {
	if b.healthy {
		return map[string]error{b.name: nil}
	}
	return map[string]error{b.name: ErrUnhealthy}
}

func TestCheck(t *testing.T) {
	for i, test := range []struct {
		checks   []services.Checkable
		healthy  bool
		expected map[string]error
	}{
		{[]services.Checkable{}, true, map[string]error{}},

		{[]services.Checkable{boolCheck{"0", true}}, true, map[string]error{"0": nil}},

		{[]services.Checkable{boolCheck{"0", true}, boolCheck{"1", true}}, true, map[string]error{"0": nil, "1": nil}},

		{[]services.Checkable{boolCheck{"0", true}, boolCheck{"1", false}}, false, map[string]error{"0": nil, "1": ErrUnhealthy}},

		{[]services.Checkable{boolCheck{"0", true}, boolCheck{"1", false}, boolCheck{"2", false}}, false, map[string]error{
			"0": nil,
			"1": ErrUnhealthy,
			"2": ErrUnhealthy,
		}},
	} {
		c := services.NewChecker()
		for _, check := range test.checks {
			require.NoError(t, c.Register(check))
		}

		require.NoError(t, c.Start())

		healthy, results := c.IsHealthy()

		assert.Equal(t, test.healthy, healthy, "case %d", i)
		assert.Equal(t, test.expected, results, "case %d", i)
	}
}

func TestNewInBackupHealthReport(t *testing.T) {
	lggr, observed := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	ibhr := services.NewInBackupHealthReport(1234, lggr)

	ibhr.Start()
	require.Eventually(t, func() bool { return observed.Len() >= 1 }, time.Second*5, time.Millisecond*100)
	require.Equal(t, "Starting InBackupHealthReport", observed.TakeAll()[0].Message)

	res, err := http.Get("http://localhost:1234/health")
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	ibhr.Stop()
	require.Eventually(t, func() bool { return observed.Len() >= 1 }, time.Second*5, time.Millisecond*100)
	require.Equal(t, "InBackupHealthReport shutdown complete", observed.TakeAll()[0].Message)
}
