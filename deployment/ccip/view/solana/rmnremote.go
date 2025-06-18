package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	solRmnRemote "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/rmn_remote"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"
)

type RMNRemoteView struct {
	ConfigPDA          string   `json:"configPDA,omitempty"`
	CursePDA           string   `json:"cursePDA,omitempty"`
	Version            uint8    `json:"version,omitempty"`
	Owner              string   `json:"owner,omitempty"`
	ProposedOwner      string   `json:"proposedOwner,omitempty"`
	DefaultCodeVersion string   `json:"defaultCodeVersion,omitempty"`
	CurseSubjects      []string `json:"curses,omitempty"`
}

func GenerateRMNRemoteView(chain cldf_solana.Chain, program solana.PublicKey, remoteChains []uint64, tokens []solana.PublicKey) (RMNRemoteView, error) {
	view := RMNRemoteView{}
	var config solRmnRemote.Config
	configPDA, _, _ := solState.FindRMNRemoteConfigPDA(program)
	err := chain.GetAccountDataBorshInto(context.Background(), configPDA, &config)
	if err != nil {
		return view, fmt.Errorf("config not found in existing state, initialize rmn first %d", chain.Selector)
	}
	view.ConfigPDA = configPDA.String()
	view.DefaultCodeVersion = config.DefaultCodeVersion.String()
	view.Owner = config.Owner.String()
	view.ProposedOwner = config.ProposedOwner.String()

	var curseAccount solRmnRemote.Curses
	cursePDA, _, _ := solState.FindRMNRemoteCursesPDA(program)
	if err = chain.GetAccountDataBorshInto(context.Background(), cursePDA, &curseAccount); err != nil {
		return view, fmt.Errorf("failed to get curse pda: %w", err)
	}
	view.CursePDA = cursePDA.String()
	view.CurseSubjects = make([]string, len(curseAccount.CursedSubjects))
	for i, curse := range curseAccount.CursedSubjects {
		view.CurseSubjects[i] = base58.Encode(curse.Value[:])
	}
	return view, nil
}
