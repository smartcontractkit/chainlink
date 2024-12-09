package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

var _ deployment.ChangeSet[[]uint64] = DeployLinkToken

// DeployLinkToken deploys a link token contract to the chain identified by the ChainSelector.
func DeployLinkToken(e deployment.Environment, chains []uint64) (deployment.ChangesetOutput, error) {
	for _, chain := range chains {
		_, ok := e.Chains[chain]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
		}
	}
	newAddresses := deployment.NewMemoryAddressBook()
	for _, chain := range chains {
		_, err := deployLinkTokenContract(
			e.Logger, e.Chains[chain], newAddresses,
		)
		if err != nil {
			return deployment.ChangesetOutput{AddressBook: newAddresses}, err
		}
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

func deployLinkTokenContract(
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (*deployment.ContractDeploy[*link_token.LinkToken], error) {
	linkToken, err := deployment.DeployContract[*link_token.LinkToken](lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*link_token.LinkToken] {
			linkTokenAddr, tx, linkToken, err2 := link_token.DeployLinkToken(
				chain.DeployerKey,
				chain.Client,
			)
			return deployment.ContractDeploy[*link_token.LinkToken]{
				Address:  linkTokenAddr,
				Contract: linkToken,
				Tx:       tx,
				Tv:       deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return linkToken, err
	}
	return linkToken, nil
}
