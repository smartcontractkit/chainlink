package changeset

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	solConfig "github.com/smartcontractkit/chainlink-ccip/chains/solana/contracts/tests/config"
	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

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

// Solana variant of DeployLinkToken -> do we want separate variants ? or we want to combine them ?
func DeployLinkTokenSolana(e deployment.Environment, chains []uint64) (deployment.ChangesetOutput, error) {
	fmt.Println("DeployLinkTokenSolana", e.SolChains)
	for _, chain := range chains {
		_, ok := e.SolChains[chain]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("sol chain not found in environment")
		}
	}
	newAddresses := deployment.NewMemoryAddressBook()
	for _, chain := range chains {
		err := deployLinkTokenContractSolana(
			e.Logger, e.SolChains[chain], newAddresses,
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

// for ethereum -> DeployContract -> calls
// 1. gethwrapper DeployContract
// 2. chain.Confirm
// 3. addressBook.Save
func deployLinkTokenContractSolana(
	lggr logger.Logger,
	chain deployment.SolChain,
	ab deployment.AddressBook,
) error {
	decimals := uint8(0)
	adminPublicKey := chain.DeployerKey.PublicKey()
	mint, _ := solana.NewRandomPrivateKey()
	mintPublicKey := mint.PublicKey()
	instructions, err := solTokenUtil.CreateToken(
		context.Background(), solConfig.Token2022Program, mintPublicKey, adminPublicKey, decimals, chain.Client, solRpc.CommitmentConfirmed,
	)
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return err
	}
	err = chain.Confirm(instructions, solCommomUtil.AddSigners(mint))
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return err
	}
	tv := deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)
	lggr.Infow("Deployed contract", "Contract", tv.String(), "addr", mintPublicKey.String(), "chain", chain.String())
	err = ab.Save(chain.Selector, mintPublicKey.String(), tv)
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return err
	}

	return nil
}
