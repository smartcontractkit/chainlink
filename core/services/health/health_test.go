package health_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/health"
	"github.com/stretchr/testify/assert"
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

func TestCheck(t *testing.T) {
	for i, test := range []struct {
		checks   []health.Checkable
		healthy  bool
		expected map[string]error
	}{
		{[]health.Checkable{}, true, map[string]error{}},

		{[]health.Checkable{boolCheck(true)}, true, map[string]error{"0": nil}},

		{[]health.Checkable{boolCheck(true), boolCheck(true)}, true, map[string]error{"0": nil, "1": nil}},

		{[]health.Checkable{boolCheck(true), boolCheck(false)}, false, map[string]error{"0": nil, "1": ErrUnhealthy}},

		{[]health.Checkable{boolCheck(true), boolCheck(false), boolCheck(false)}, false, map[string]error{
			"0": nil,
			"1": ErrUnhealthy,
			"2": ErrUnhealthy,
		}},
	} {
		c := health.NewChecker()
		for i, check := range test.checks {
			c.Register(fmt.Sprint(i), check)
		}

		c.Start()

		healthy, results := c.IsHealthy()

		assert.Equal(t, test.healthy, healthy, "case %d", i)
		assert.Equal(t, test.expected, results, "case %d", i)
	}
}
