package solana

import (
	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/shared"
)

type state struct {
	forwarderProgramID, forwarderState solana.PublicKey
}

func loadOnchainState(env cldf.Environment, chainSel uint64) (*state, error) {
	addresses, err := env.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}
	state := &state{}
	for address, tvStr := range addresses {
		switch tvStr.Type {
		case shared.Forwarder:
			state.forwarderProgramID = solana.MustPublicKeyFromBase58(address)
		case shared.ForwarderState:
			state.forwarderState = solana.MustPublicKeyFromBase58(address)
		}
	}

	return state, nil
}
