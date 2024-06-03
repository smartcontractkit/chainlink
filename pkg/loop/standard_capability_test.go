package loop_test

import (
	"context"
	"testing"

	"github.com/hashicorp/go-plugin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	sctest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/capability/test"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
)

func TestPluginStandardCapability(t *testing.T) {
	t.Parallel()

	log := logger.Test(t)

	stopCh := newStopCh(t)
	test.PluginTest(t, loop.PluginStandardCapabilityName,
		&loop.StandardCapabilityLoop{
			Logger:       log,
			PluginServer: sctest.StandardCapabilityService{},
			BrokerConfig: loop.BrokerConfig{
				Logger: logger.Test(t),
				StopCh: stopCh}},
		func(t *testing.T, s loop.StandardCapability) {
			info, err := s.Info(context.Background())
			assert.NoError(t, err)
			assert.Equal(t, "1", info.ID)
			assert.Equal(t, capabilities.CapabilityTypeAction, info.CapabilityType)

			err = s.Initialise(context.Background(), "", nil, nil, nil, nil, nil, nil)
			assert.NoError(t, err)
		})
}

func TestRunningStandardCapabilityPluginOutOfProcess(t *testing.T) {
	t.Parallel()
	stopCh := newStopCh(t)

	scs := newOutOfProcessStandardCapabilityService(t, true, stopCh)

	info, err := scs.Info(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "1", info.ID)
	assert.Equal(t, capabilities.CapabilityTypeAction, info.CapabilityType)

	err = scs.Initialise(context.Background(), "", nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
}

func newOutOfProcessStandardCapabilityService(t *testing.T, staticChecks bool, stopCh <-chan struct{}) loop.StandardCapability {
	scl := loop.StandardCapabilityLoop{Logger: logger.Test(t), BrokerConfig: loop.BrokerConfig{Logger: logger.Test(t), StopCh: stopCh}}
	cc := scl.ClientConfig()
	cc.Cmd = NewHelperProcessCommand(loop.PluginStandardCapabilityName, staticChecks, 0)
	c := plugin.NewClient(cc)
	t.Cleanup(c.Kill)
	client, err := c.Client()
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })
	require.NoError(t, client.Ping())
	i, err := client.Dispense(loop.PluginStandardCapabilityName)
	require.NoError(t, err)
	return i.(loop.StandardCapability)
}
