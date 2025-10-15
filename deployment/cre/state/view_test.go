package state_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/deployment/cre/state"
	test2 "github.com/smartcontractkit/chainlink/deployment/cre/test"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"

	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var oracleConfig = changeset.OracleConfig{
	DeltaProgressMillis:               30000,
	DeltaResendMillis:                 5000,
	DeltaInitialMillis:                5000,
	DeltaRoundMillis:                  2000,
	DeltaGraceMillis:                  500,
	DeltaCertifiedCommitRequestMillis: 1000,
	DeltaStageMillis:                  30000,
	MaxRoundsPerEpoch:                 10,
	TransmissionSchedule:              []int{},
	MaxDurationQueryMillis:            1000,
	MaxDurationObservationMillis:      1000,
	MaxDurationShouldAcceptMillis:     1000,
	MaxDurationShouldTransmitMillis:   1000,
	MaxFaultyOracles:                  1,
	ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
		MaxQueryLengthBytes:       1000000,
		MaxObservationLengthBytes: 1000000,
		MaxReportLengthBytes:      1000000,
		MaxOutcomeLengthBytes:     1000000,
		MaxReportCount:            20,
		MaxBatchSize:              20,
		OutcomePruningThreshold:   3600,
		RequestTimeout:            30 * time.Second,
	},
	UniqueReports: true,
}

func TestCREView(t *testing.T) {
	t.Parallel()
	env := test2.SetupEnvV2(t, false)
	//originalAddressBook := env.Env.ExistingAddresses
	oracleConfig.TransmissionSchedule = []int{len(env.Env.NodeIDs)}

	addrs := env.Env.DataStore.Addresses().Filter(
		datastore.AddressRefByChainSelector(env.RegistrySelector),
	)

	var newCapabilityRegistryAddr string
	for _, addr := range addrs {
		if newCapabilityRegistryAddr != "" {
			break
		}
		switch addr.Type {
		case datastore.ContractType(contracts.CapabilitiesRegistry):
			newCapabilityRegistryAddr = addr.Address
			continue
		default:
			continue
		}
	}

	t.Run("successfully generates a view of the CRE state", func(t *testing.T) {
		var prevView json.RawMessage = []byte("{}")
		a, err := state.ViewCRE(*env.Env, prevView)
		require.NoError(t, err)
		b, err := a.MarshalJSON()
		require.NoError(t, err)
		require.NotEmpty(t, b)

		var outView state.CREView
		require.NoError(t, json.Unmarshal(b, &outView))

		chainName, err := chain_selectors.GetChainNameFromSelector(env.RegistrySelector)
		require.NoError(t, err)

		viewChain, ok := outView.Chains[chainName]
		require.True(t, ok)
		_, ok = viewChain.CapabilityRegistry[newCapabilityRegistryAddr]
		require.True(t, ok)
		require.Len(t, viewChain.CapabilityRegistry, 1)
	})
}
