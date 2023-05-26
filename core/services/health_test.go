package services_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services"
)

var ErrUnhealthy = errors.New("Unhealthy")

type boolCheck struct {
	name    string
	healthy bool
}

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
		for i, check := range test.checks {
			require.NoError(t, c.Register(fmt.Sprint(i), check))
		}

		require.NoError(t, c.Start())

		healthy, results := c.IsHealthy()

		assert.Equal(t, test.healthy, healthy, "case %d", i)
		assert.Equal(t, test.expected, results, "case %d", i)
	}
}
