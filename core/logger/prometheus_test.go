package logger

import (
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusLogger_Counters(t *testing.T) {
	t.Parallel()

	l := newPrometheusLogger(NullLogger)
	repeat(l.Warn, 1)
	repeat(l.Error, 2)
	repeat(l.Critical, 3)
	repeat(l.Panic, 4)
	repeat(l.Fatal, 5)

	assertCounterValue(t, warnCounter, 1)
	assertCounterValue(t, errorCounter, 2)
	assertCounterValue(t, criticalCounter, 3)
	assertCounterValue(t, panicCounter, 4)
	assertCounterValue(t, fatalCounter, 5)

	l.Warnf("msg")
	l.Warnw("msg")
	assertCounterValue(t, warnCounter, 3)

	l.Errorf("msg")
	l.Errorf("msg")
	l.ErrorIfClosing(&errorCloser{}, "foo")
	assertCounterValue(t, errorCounter, 5)

	l.Criticalf("msg")
	l.Criticalw("msg")
	assertCounterValue(t, criticalCounter, 5)

	l.Panicf("msg")
	l.Panicw("msg")
	l.Recover(nil)
	assertCounterValue(t, panicCounter, 7)

	l.Fatalf("msg")
	l.Fatalw("msg")
	assertCounterValue(t, fatalCounter, 7)
}

func assertCounterValue(t *testing.T, c prometheus.Counter, v int) {
	var m io_prometheus_client.Metric
	err := c.Write(&m)
	assert.NoError(t, err)
	assert.Equal(t, v, int(m.GetCounter().GetValue()))
}

func repeat(f func(args ...interface{}), c int) {
	for c > 0 {
		f()
		c--
	}
}

type errorCloser struct{}

func (c errorCloser) Close() error {
	return errors.New("error")
}
