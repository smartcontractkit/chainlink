package changeset_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslb "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
)

func TestDeployOCR3(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	ab := deployment.NewMemoryAddressBook()
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.AllChainSelectors()[0]
	// err if no capabilities registry on chain 0
	_, err := changeset.DeployOCR3(lggr, env, ab, registrySel)
	require.Error(t, err)

	// fake capabilities registry
	err = ab.Save(registrySel, "0x0000000000000000000000000000000000000001", kslb.CapabilityRegistryTypeVersion)
	require.NoError(t, err)
	resp, err := changeset.DeployOCR3(lggr, env, ab, registrySel)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// OCR3 should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 2)

	// nothing on chain 1
	require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
	assert.Len(t, oaddrs, 0)
}
