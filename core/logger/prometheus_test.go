package logger

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusLogger_Counters(t *testing.T) {
	t.Parallel()

	createRandomNameCounter := func() prometheus.Counter {
		name := fmt.Sprintf("test_counter_%d", rand.Int31n(999999))
		return promauto.NewCounter(prometheus.CounterOpts{Name: name})
	}

	warnCounter := createRandomNameCounter()
	errorCounter := createRandomNameCounter()
	criticalCounter := createRandomNameCounter()
	panicCounter := createRandomNameCounter()
	fatalCounter := createRandomNameCounter()

	l := newPrometheusLoggerWithCounters(NullLogger, warnCounter, errorCounter, criticalCounter, panicCounter, fatalCounter)
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

	nl := l.Named("foo") // reusing counters
	repeat(nl.Warn, 1)
	repeat(nl.Error, 1)
	repeat(nl.Critical, 1)
	repeat(nl.Panic, 1)
	repeat(nl.Fatal, 1)

	assertCounterValue(t, warnCounter, 2)
	assertCounterValue(t, errorCounter, 3)
	assertCounterValue(t, criticalCounter, 4)
	assertCounterValue(t, panicCounter, 5)
	assertCounterValue(t, fatalCounter, 6)

	wl := l.With("bar") // reusing counters
	repeat(wl.Warn, 1)
	repeat(wl.Error, 1)
	repeat(wl.Critical, 1)
	repeat(wl.Panic, 1)
	repeat(wl.Fatal, 1)

	assertCounterValue(t, warnCounter, 3)
	assertCounterValue(t, errorCounter, 4)
	assertCounterValue(t, criticalCounter, 5)
	assertCounterValue(t, panicCounter, 6)
	assertCounterValue(t, fatalCounter, 7)

	l.Warnf("msg")
	l.Warnw("msg")
	assertCounterValue(t, warnCounter, 5)

	l.Errorf("msg")
	l.Errorf("msg")
	assertCounterValue(t, errorCounter, 6)

	l.Criticalf("msg")
	l.Criticalw("msg")
	assertCounterValue(t, criticalCounter, 7)

	l.Panicf("msg")
	l.Panicw("msg")
	l.Recover(nil)
	assertCounterValue(t, panicCounter, 9)

	l.Fatalf("msg")
	l.Fatalw("msg")
	assertCounterValue(t, fatalCounter, 9)
}

func assertCounterValue(t *testing.T, c prometheus.Counter, v int) {
	var m io_prometheus_client.Metric
	err := c.Write(&m)
	assert.NoError(t, err)
	assert.Equal(t, v, int(m.GetCounter().GetValue()))
}

func repeat(f func(args ...interface{}), c int) {
	for ; c > 0; c-- {
		f()
	}
}

type errorCloser struct{}

func (c errorCloser) Close() error {
	return errors.New("error")
}
