package solana

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfdeployment "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	ks_forwarder "github.com/smartcontractkit/chainlink-solana/contracts/generated/keystone_forwarder"
	cdeployment "github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

var _ cldf.ChangeSet[*DeployRequest] = DeployForwarder

// move
var Forwarder cldfdeployment.ContractType = "Forwarder"

func DeployForwarder(env cldf.Environment, req *DeployRequest) (cldf.ChangesetOutput, error) {
	if req.BuildConfig != nil {
		err := helpers.BuildSolana(env, req.BuildConfig, keystoneBuildParams)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
	}
	chain := env.SolChains[req.ChainSel]
	ab := cldf.NewMemoryAddressBook()

	address, err := helpers.DeployAndMaybeSaveToAddressBook(env, chain, ab, Forwarder, cdeployment.Version1_0_0, false, "")
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// initialize
	ks_forwarder.SetProgramID(address)

	key, err := solana.NewRandomPrivateKey()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create random keys: %w", err)
	}

	instruction, err := ks_forwarder.NewInitializeInstruction(key.PublicKey(), chain.DeployerKey.PublicKey(), solana.SystemProgramID).ValidateAndBuild()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	instructions := []solana.Instruction{instruction}
	if err = chain.Confirm(instructions, solanaUtils.AddSigners(key)); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm ")
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
	}, nil
}
