package solana

import (
	"github.com/gagliardetto/solana-go"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type state struct {
	forwarderProgramID, forwarderState solana.PublicKey
}

func loadOnchainState(env cldf.Environment, chainSel uint64) (*state, error) {
	addresses, err := env.DataStore.Addresses().Fetch()
	if err != nil {
		return nil, err
	}
	state := &state{}
	for _, tvStr := range addresses {
		switch tvStr.Type {
		case ForwarderContract:
			state.forwarderProgramID = solana.MustPublicKeyFromBase58(tvStr.Address)
		case ForwarderState:
			state.forwarderState = solana.MustPublicKeyFromBase58(tvStr.Address)
		}
	}

	return state, nil
}
