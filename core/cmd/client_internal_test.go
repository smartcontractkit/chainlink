package cmd

import (
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/stretchr/testify/assert"
)

func TestInitConfigAndLogger_ClosesPreviousLogger(t *testing.T) {
	c := Client{}
	var closeFnCalled bool
	c.CloseLogger = func() error {
		closeFnCalled = true
		return nil
	}

	var opts chainlink.GeneralConfigOpts
	err := c.initConfigAndLogger(&opts, []string{}, "")
	assert.NoError(t, err)
	assert.True(t, closeFnCalled)
}

func TestInitConfigAndLogger_ReturnsErrorIfOldLoggerCantBeClosed(t *testing.T) {
	c := Client{}
	c.CloseLogger = func() error {
		return errors.New("error closing logger")
	}

	var opts chainlink.GeneralConfigOpts
	err := c.initConfigAndLogger(&opts, []string{}, "")
	assert.ErrorContains(t, err, "failed to close initialized logger")
}
