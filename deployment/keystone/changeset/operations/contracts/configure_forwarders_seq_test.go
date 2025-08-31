package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func Test_ConfigureForwardersSeq(t *testing.T) {
	te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
		WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
		AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
		WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
		NumChains:       2,
	})

	var wfNodes []string
	for _, id := range te.GetP2PIDs("wfDon") {
		wfNodes = append(wfNodes, id.String())
	}
	configureForwardersDeps := contracts.ConfigureForwardersSeqDeps{
		Env:      &te.Env,
		Registry: te.CapabilitiesRegistry(),
	}
	configureForwardersInput := contracts.ConfigureForwardersSeqInput{
		RegistryChainSel: te.RegistrySelector,
		DONs: []contracts.ConfigureKeystoneDON{
			{
				Name:    "wfDon",
				NodeIDs: wfNodes,
			},
		},
	}
	b := optest.NewBundle(t)
	_, err := operations.ExecuteSequence(b, contracts.ConfigureForwardersSeq, configureForwardersDeps, configureForwardersInput)
	require.NoError(t, err)
}
