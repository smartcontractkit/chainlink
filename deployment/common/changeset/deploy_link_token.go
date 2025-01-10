package changeset

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

var _ deployment.ChangeSet[[]uint64] = DeployLinkToken

const (
	TokenDecimalsSolana = 9
)

// DeployLinkToken deploys a link token contract to the chain identified by the ChainSelector.
func DeployLinkToken(e deployment.Environment, chains []uint64) (deployment.ChangesetOutput, error) {
	for _, chain := range chains {
		_, evmOk := e.Chains[chain]
		_, solOk := e.SolChains[chain]
		if !evmOk && !solOk {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chain)
		}
	}
	newAddresses := deployment.NewMemoryAddressBook()
	for _, chain := range chains {
		family, err := chainsel.GetSelectorFamily(chain)
		if err != nil {
			return deployment.ChangesetOutput{AddressBook: newAddresses}, err
		}
		switch family {
		case chainsel.FamilyEVM:
			// Deploy EVM LINK token
			_, err := deployLinkTokenContractEVM(
				e.Logger, e.Chains[chain], newAddresses,
			)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}
		case chainsel.FamilySolana:
			// Deploy Solana LINK token
			err := deployLinkTokenContractSolana(
				e.Logger, e.SolChains[chain], newAddresses,
			)
			if err != nil {
				return deployment.ChangesetOutput{AddressBook: newAddresses}, err
			}
		}
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

func deployLinkTokenContractEVM(
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

func deployLinkTokenContractSolana(
	lggr logger.Logger,
	chain deployment.SolChain,
	ab deployment.AddressBook,
) error {
	adminPublicKey := chain.DeployerKey.PublicKey()
	mint, _ := solana.NewRandomPrivateKey()
	// this is the token address
	mintPublicKey := mint.PublicKey()
	instructions, err := solTokenUtil.CreateToken(
		context.Background(), solana.Token2022ProgramID, mintPublicKey, adminPublicKey, TokenDecimalsSolana, chain.Client, solRpc.CommitmentConfirmed,
	)
	if err != nil {
		lggr.Errorw("Failed to generate instructions for link token deployment", "chain", chain.String(), "err", err)
		return err
	}
	err = chain.Confirm(instructions, solCommomUtil.AddSigners(mint))
	if err != nil {
		lggr.Errorw("Failed to confirm instructions for link token deployment", "chain", chain.String(), "err", err)
		return err
	}
	tv := deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)
	lggr.Infow("Deployed contract", "Contract", tv.String(), "addr", mintPublicKey.String(), "chain", chain.String())
	err = ab.Save(chain.Selector, mintPublicKey.String(), tv)
	if err != nil {
		lggr.Errorw("Failed to save link token", "chain", chain.String(), "err", err)
		return err
	}

	return nil
}
