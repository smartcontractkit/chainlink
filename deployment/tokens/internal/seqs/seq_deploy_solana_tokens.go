package seqs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

// SeqDeploySolTokensDeps contains the dependencies for the sequence to deploy LINK token contracts
// on Solana.
type SeqDeploySolTokensDeps struct {
	SolChains map[uint64]cldf_solana.Chain
	AddrBook  cldf.AddressBook
	Datastore datastore.MutableDataStore
}

// SeqDeploySolTokensInput is the input to the SeqDeploySolTokensInput sequence.
type SeqDeploySolTokensInput struct {
	// ChainSelectors are the chain selectors of the chains to which the Link Token contract
	ChainSelectors []uint64 `json:"chainSelectors"`

	// Qualifier is a string that will be used to tag the deployed contracts in the address book and datastore.
	Qualifier string `json:"qualifier"`

	// Labels are a list of labels that will be used to tag the deployed contracts in the address book and datastore.
	Labels []string `json:"labels"`
}

// SeqDeploySolTokensOutput is the output of the SeqDeploySolTokensOutput sequence.
type SeqDeploySolTokensOutput struct {
	// Addresses are the addresses of the deployed Link Token contracts.
	Addresses []string `json:"address"`
}

// SeqDeploySolTokens is a sequence that deploys LINK token contracts across multiple Solana
// chains.
//
// All provided chain selectors must reference Solana chains.
var SeqDeploySolTokens = operations.NewSequence(
	"seq-deploy-sol-tokens",
	semver.MustParse("1.0.0"),
	"Deploy Solana LINK token contracts across multiple chains",
	func(b operations.Bundle, deps SeqDeploySolTokensDeps, in SeqDeploySolTokensInput) (SeqDeploySolTokensOutput, error) {
		out := SeqDeploySolTokensOutput{
			Addresses: make([]string, 0),
		}

		for _, csel := range in.ChainSelectors {
			fam, err := chainsel.GetSelectorFamily(csel)
			if err != nil {
				return out, err
			}

			if fam != chainsel.FamilySolana {
				return out, fmt.Errorf("chain selector %d is not a Solana chain", csel)
			}

			chain := deps.SolChains[csel]

			deployReport, err := operations.ExecuteOperation(b, ops.OpSolDeployLinkToken,
				ops.OpSolDeployLinkTokenDeps{
					Client:      chain.Client,
					ConfirmFunc: chain.Confirm,
				},
				ops.OpSolDeployLinkTokenInput{
					ChainSelector:       csel,
					TokenAdminPublicKey: chain.DeployerKey.PublicKey(),
				},
			)
			if err != nil {
				return out, err
			}

			_, err = operations.ExecuteSequence(b, SeqPersistAddress,
				SeqPersistAddressDeps{
					AddrBook:  deps.AddrBook,
					Datastore: deps.Datastore,
				},
				SeqPersistAddressInput{
					ChainSelector: csel,
					Address:       deployReport.Output.MintPublicKey.String(),
					Type:          deployReport.Output.Type,
					Version:       deployReport.Output.Version,
					Qualifier:     in.Qualifier,
					Labels:        in.Labels,
				},
			)
			if err != nil {
				return out, err
			}

			out.Addresses = append(out.Addresses, deployReport.Output.MintPublicKey.String())
		}

		return out, nil
	},
)
