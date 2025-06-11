package seqs

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

// SeqDeployEVMTokensDeps contains the dependencies for the SeqDeployEVMTokensDeps sequence.
type SeqDeployEVMTokensDeps struct {
	EVMChains map[uint64]cldf_evm.Chain
	AddrBook  cldf.AddressBook
	Datastore datastore.MutableDataStore
}

// SeqDeployEVMTokensInput is the input to the SeqDeployEVMTokensInput sequence.
type SeqDeployEVMTokensInput struct {
	// ChainSelectors are the chain selectors of the chains to which the Link Token contract
	ChainSelectors []uint64 `json:"chainSelectors"`

	// Qualifier is a string that will be used to tag the deployed contracts in the address book and datastore.
	Qualifier string `json:"qualifier"`

	// Labels are a list of labels that will be used to tag the deployed contracts in the address book and datastore.
	Labels []string `json:"labels"`
}

// SeqDeployEVMTokensOutput is the output of the SeqDeployTokens sequence.
type SeqDeployEVMTokensOutput struct {
	// Addresses are the addresses of the deployed Link Token contracts.
	Addresses []string `json:"address"`
}

// SeqDeployEVMTokens is a sequence that deploys LINK token contracts across EVM chains.
var SeqDeployEVMTokens = operations.NewSequence(
	"seq-deploy-tokens",
	semver.MustParse("1.0.0"),
	"Deploy LINK token contracts across multiple chains",
	func(b operations.Bundle, deps SeqDeployEVMTokensDeps, in SeqDeployEVMTokensInput) (SeqDeployEVMTokensOutput, error) {
		out := SeqDeployEVMTokensOutput{
			Addresses: make([]string, 0),
		}

		for _, csel := range in.ChainSelectors {
			fam, err := chainsel.GetSelectorFamily(csel)
			if err != nil {
				return out, err
			}

			if fam != chainsel.FamilyEVM {
				return out, fmt.Errorf("chain selector %d is not an evm chain", csel)
			}

			chain := deps.EVMChains[csel]

			// Deploy the link token
			deployReport, err := operations.ExecuteOperation(b, ops.OpEVMDeployLinkToken,
				ops.OpEVMDeployLinkTokenDeps{
					Auth:        chain.DeployerKey,
					Backend:     chain.Client,
					ConfirmFunc: chain.Confirm,
				},
				ops.OpEVMDeployLinkTokenInput{
					ChainSelector: csel,
				},
			)
			if err != nil {
				return out, err
			}

			out.Addresses = append(out.Addresses, deployReport.Output.Address.String())

			_, err = operations.ExecuteSequence(b, SeqPersistAddress,
				SeqPersistAddressDeps{
					AddrBook:  deps.AddrBook,
					Datastore: deps.Datastore,
				},
				SeqPersistAddressInput{
					ChainSelector: csel,
					Address:       deployReport.Output.Address.String(),
					Type:          deployReport.Output.Type,
					Version:       deployReport.Output.Version,
					Qualifier:     in.Qualifier,
					Labels:        in.Labels,
				},
			)
			if err != nil {
				return out, err
			}
		}

		return out, nil
	},
)
