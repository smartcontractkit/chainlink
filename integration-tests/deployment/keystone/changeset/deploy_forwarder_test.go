package changeset_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	kslb "github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment/memory"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.AllChainSelectors()[0]
	t.Run("err if no capabilities registry on registry chain", func(t *testing.T) {
		ab := deployment.NewMemoryAddressBook()
		m := make(map[uint64]map[string]deployment.TypeAndVersion)
		m[registrySel] = map[string]deployment.TypeAndVersion{
			"0x0000000000000000000000000000000000000002": kslb.OCR3CapabilityTypeVersion,
		}
		deployment.NewMemoryAddressBookFromMap(m)
		// capabilities registry and ocr3 must be deployed on registry chain
		_, err := changeset.DeployForwarder(lggr, env, ab, registrySel)
		require.Error(t, err)
	})

	t.Run("err if no ocr3 on registry chain", func(t *testing.T) {
		ab := deployment.NewMemoryAddressBook()
		m := make(map[uint64]map[string]deployment.TypeAndVersion)
		m[registrySel] = map[string]deployment.TypeAndVersion{
			"0x0000000000000000000000000000000000000001": kslb.CapabilityRegistryTypeVersion,
		}
		deployment.NewMemoryAddressBookFromMap(m)
		// capabilities registry and ocr3 must be deployed on registry chain
		_, err := changeset.DeployForwarder(lggr, env, ab, registrySel)
		require.Error(t, err)
	})

	t.Run("should deploy forwarder", func(t *testing.T) {
		ab := deployment.NewMemoryAddressBook()
		// fake capabilities registry
		err := ab.Save(registrySel, "0x0000000000000000000000000000000000000001", kslb.CapabilityRegistryTypeVersion)
		require.NoError(t, err)

		// fake ocr3
		err = ab.Save(registrySel, "0x0000000000000000000000000000000000000002", kslb.OCR3CapabilityTypeVersion)
		require.NoError(t, err)
		// deploy forwarder
		resp, err := changeset.DeployForwarder(lggr, env, ab, registrySel)
		require.NoError(t, err)
		require.NotNil(t, resp)
		// registry, ocr3, forwarder should be deployed on registry chain
		addrs, err := resp.AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 3)

		// only forwarder on chain 1
		require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
		oaddrs, err := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
		require.NoError(t, err)
		require.Len(t, oaddrs, 1)
	})
}
