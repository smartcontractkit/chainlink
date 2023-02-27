package services_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services"
)

var ErrUnhealthy = errors.New("Unhealthy")

type boolCheck bool

func (b boolCheck) Ready() error {
	if b {
		return nil
	}
	return errors.New("Not ready")
}

func (b boolCheck) Healthy() error {
	if b {
		return nil
	}
	return ErrUnhealthy
}

func (b boolCheck) HealthReport() map[string]error {
	if b {
		return map[string]error{"boolCheck": nil}
	}
	return map[string]error{"boolCheck": ErrUnhealthy}
}

func TestCheck(t *testing.T) {
	for i, test := range []struct {
		checks   []services.Checkable
		healthy  bool
		expected map[string]error
	}{
		{[]services.Checkable{}, true, map[string]error{}},

		{[]services.Checkable{boolCheck(true)}, true, map[string]error{"0": nil}},

		{[]services.Checkable{boolCheck(true), boolCheck(true)}, true, map[string]error{"0": nil, "1": nil}},

		{[]services.Checkable{boolCheck(true), boolCheck(false)}, false, map[string]error{"0": nil, "1": ErrUnhealthy}},

		{[]services.Checkable{boolCheck(true), boolCheck(false), boolCheck(false)}, false, map[string]error{
			"0": nil,
			"1": ErrUnhealthy,
			"2": ErrUnhealthy,
		}},
	} {
		c := services.NewChecker()
		for i, check := range test.checks {
			require.NoError(t, c.Register(fmt.Sprint(i), check))
		}

		require.NoError(t, c.Start())

		healthy, results := c.IsHealthy()

		assert.Equal(t, test.healthy, healthy, "case %d", i)
		assert.Equal(t, test.expected, results, "case %d", i)
	}
}
