package deployment_solana_common

import (
	"github.com/smartcontractkit/chainlink/deployment"
	deployment_solana "github.com/smartcontractkit/chainlink/deployment/solana/extension"
)

type DeployAndMintSolanaTokenConfig struct {
	ChainSelector    uint64
	TokenProgramName string
	TokenDecimals    uint8
	MintAmount       uint64
	AmountToAddress  map[string]uint64 // address -> amount
	ATAList          []string
}

func DeployAndMintSolanaToken(e deployment.Environment, cfg DeployAndMintSolanaTokenConfig) (deployment.ChangesetOutput, error) {
	// get chain
	chain := e.SolChains[cfg.ChainSelector]
	// sol deps
	deps := deployment_solana.SolanaDeps{
		Client:         chain.Client,
		Auth:           *chain.DeployerKey,
		SendAndConfirm: deployment_solana.SendAndConfirm,
	}
	// get addresses
	deployInput := DeploySolanaTokenInput{
		TokenProgramName: cfg.TokenProgramName,
		TokenDecimals:    cfg.TokenDecimals,
	}

	deployReport, err := deployment.ExecuteOp(e.OpEnv, DeploySolanaTokenOp, deps, deployInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// ata input
	ataInput := CreateSolanaTokenATAInput{
		TokenPubkey:  deployReport.Output.Address,
		TokenProgram: cfg.TokenProgramName,
		ATAList:      cfg.ATAList,
	}

	// create ATAs
	_, err = deployment.ExecuteOp(e.OpEnv, CreateSolanaTokenATAOp, deps, ataInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// mint input
	mintInput := MintSolanaTokenInput{
		TokenProgram:    cfg.TokenProgramName,
		TokenPubkey:     deployReport.Output.Address,
		AmountToAddress: cfg.AmountToAddress,
	}

	// mint tokens
	_, err = deployment.ExecuteOp(e.OpEnv, MintSolanaTokenOp, deps, mintInput)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	// We can compile every address deployed in the address book
	ab := deployment.NewMemoryAddressBook()
	ab.Save(cfg.ChainSelector, deployReport.Output.Address.String(), deployment.TypeAndVersion{
		Type: "Token",
	})

	return deployment.ChangesetOutput{
		Reports:     e.OpEnv.Reporter.GetReports(),
		AddressBook: ab,
	}, nil

}
