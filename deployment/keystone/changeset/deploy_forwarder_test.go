package changeset_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	kslb "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	var (
		ocrTV = deployment.NewTypeAndVersion(kslb.OCR3Capability, deployment.Version1_0_0)
		crTV  = deployment.NewTypeAndVersion(kslb.CapabilitiesRegistry, deployment.Version1_0_0)
	)

	registrySel := env.AllChainSelectors()[0]
	t.Run("err if no capabilities registry on registry chain", func(t *testing.T) {
		m := make(map[uint64]map[string]deployment.TypeAndVersion)
		m[registrySel] = map[string]deployment.TypeAndVersion{
			"0x0000000000000000000000000000000000000002": ocrTV,
		}
		env.ExistingAddresses = deployment.NewMemoryAddressBookFromMap(m)
		// capabilities registry and ocr3 must be deployed on registry chain
		_, err := changeset.DeployForwarder(env, registrySel)
		require.Error(t, err)
	})

	t.Run("err if no ocr3 on registry chain", func(t *testing.T) {
		m := make(map[uint64]map[string]deployment.TypeAndVersion)
		m[registrySel] = map[string]deployment.TypeAndVersion{
			"0x0000000000000000000000000000000000000001": crTV,
		}
		env.ExistingAddresses = deployment.NewMemoryAddressBookFromMap(m)
		// capabilities registry and ocr3 must be deployed on registry chain
		_, err := changeset.DeployForwarder(env, registrySel)
		require.Error(t, err)
	})

	t.Run("should deploy forwarder", func(t *testing.T) {
		ab := deployment.NewMemoryAddressBook()
		// fake capabilities registry
		err := ab.Save(registrySel, "0x0000000000000000000000000000000000000001", crTV)
		require.NoError(t, err)

		// fake ocr3
		err = ab.Save(registrySel, "0x0000000000000000000000000000000000000002", ocrTV)
		require.NoError(t, err)
		// deploy forwarder
		env.ExistingAddresses = ab
		resp, err := changeset.DeployForwarder(env, registrySel)
		require.NoError(t, err)
		require.NotNil(t, resp)
		// registry, ocr3, forwarder should be deployed on registry chain
		addrs, err := resp.AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// only forwarder on chain 1
		require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
		oaddrs, err := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
		require.NoError(t, err)
		require.Len(t, oaddrs, 1)
	})
}
