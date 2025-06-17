package common_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/stretchr/testify/require"

	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	capcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
)

func Test_HashedCapabilityId(t *testing.T) {
	transactor := testutils.MustNewSimTransactor(t)
	sb := backends.NewSimulatedBackend(core.GenesisAlloc{
		transactor.From: {Balance: assets.Ether(1000).ToInt()},
	}, 30e6)

	crAddress, _, _, err := kcr.DeployCapabilitiesRegistry(transactor, sb)
	require.NoError(t, err)
	sb.Commit()

	cr, err := kcr.NewCapabilitiesRegistry(crAddress, sb)
	require.NoError(t, err)

	// add a capability, ignore cap config for simplicity.
	_, err = cr.AddCapabilities(transactor, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          "ccip",
			Version:               "v1.0.0",
			CapabilityType:        0,
			ResponseType:          0,
			ConfigurationContract: common.Address{},
		},
	})
	require.NoError(t, err)
	sb.Commit()

	hidExpected, err := cr.GetHashedCapabilityId(nil, "ccip", "v1.0.0")
	require.NoError(t, err)

	hid, err := capcommon.HashedCapabilityID("ccip", "v1.0.0")
	require.NoError(t, err)

	require.Equal(t, hidExpected, hid)
}
