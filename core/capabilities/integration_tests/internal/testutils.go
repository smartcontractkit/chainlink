package testutils

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/require"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/keystone_capability_registry"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

// var OWNER_ADDR = types.MustEIP55Address("0x0000000000000000000000000000000000000001").Address()

var DataStreamsReportCapability = kcr.CapabilityRegistryCapability{
	LabelledName: "data-streams-report",
	Version:      "1.0.0",
	ResponseType: uint8(0),
}

var WriteChainCapability = kcr.CapabilityRegistryCapability{
	LabelledName: "write-chain",
	Version:      "1.0.1",
	ResponseType: uint8(1),
}

func StartNewChain(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend) {
	owner := testutils.MustNewSimTransactor(t)

	oneEth, _ := new(big.Int).SetString("100000000000000000000", 10)
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2 // 60 M blocks

	simulatedBackend := backends.NewSimulatedBackend(core.GenesisAlloc{owner.From: {
		Balance: oneEth,
	}}, gasLimit)
	simulatedBackend.Commit()

	return owner, simulatedBackend
}

func DeployCapabilityRegistry(t *testing.T, owner *bind.TransactOpts, simulatedBackend *backends.SimulatedBackend) (capabilityRegistry *kcr.CapabilityRegistry) {
	capabilityRegistryAddress, _, capabilityRegistry, err := kcr.DeployCapabilityRegistry(owner, simulatedBackend)
	require.NoError(t, err, "DeployCapabilityRegistry failed")

	fmt.Println("Deployed CapabilityRegistry at", capabilityRegistryAddress.Hex())
	simulatedBackend.Commit()

	return capabilityRegistry
}

func AddCapability(t *testing.T, owner *bind.TransactOpts, simulatedBackend *backends.SimulatedBackend, capabilityRegistry *kcr.CapabilityRegistry, capability kcr.CapabilityRegistryCapability) {
	_, err := capabilityRegistry.AddCapability(owner, capability)
	require.NoError(t, err, "AddCapability failed for %s", capability.LabelledName)
	simulatedBackend.Commit()
}
