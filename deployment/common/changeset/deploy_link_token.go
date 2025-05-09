package changeset

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"golang.org/x/sync/errgroup"

	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ deployment.ChangeSet[[]uint64] = DeployLinkToken

const (
	TokenDecimalsSolana = 9
)

// DeployLinkToken deploys a link token contract to the chain identified by the ChainSelector.
func DeployLinkToken(e deployment.Environment, chains []uint64) (deployment.ChangesetOutput, error) {
	err := deployment.ValidateSelectorsInEnvironment(e, chains)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	newAddresses := deployment.NewMemoryAddressBook()
	deployGrp := errgroup.Group{}
	for _, chain := range chains {
		family, err := chainsel.GetSelectorFamily(chain)
		if err != nil {
			return deployment.ChangesetOutput{AddressBook: newAddresses}, err
		}
		var deployFn func() error
		switch family {
		case chainsel.FamilyEVM:
			// Deploy EVM LINK token
			deployFn = func() error {
				_, err := deployLinkTokenContractEVM(
					e.Logger, e.Chains[chain], newAddresses,
				)
				return err
			}
		case chainsel.FamilySolana:
			// Deploy Solana LINK token
			deployFn = func() error {
				err := deployLinkTokenContractSolana(
					e.Logger, e.SolChains[chain], newAddresses,
				)
				return err
			}
		}
		deployGrp.Go(func() error {
			err := deployFn()
			if err != nil {
				e.Logger.Errorw("Failed to deploy link token", "chain", chain, "err", err)
				return fmt.Errorf("failed to deploy link token for chain %d: %w", chain, err)
			}
			return nil
		})
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, deployGrp.Wait()
}

// DeployStaticLinkToken deploys a static link token contract to the chain identified by the ChainSelector.
func DeployStaticLinkToken(e deployment.Environment, chains []uint64) (deployment.ChangesetOutput, error) {
	err := deployment.ValidateSelectorsInEnvironment(e, chains)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	newAddresses := deployment.NewMemoryAddressBook()
	for _, chainSel := range chains {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment: %d", chainSel)
		}
		_, err := cldf.DeployContract[*link_token_interface.LinkToken](e.Logger, chain, newAddresses,
			func(chain deployment.Chain) cldf.ContractDeploy[*link_token_interface.LinkToken] {
				linkTokenAddr, tx, linkToken, err2 := link_token_interface.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
				return cldf.ContractDeploy[*link_token_interface.LinkToken]{
					Address:  linkTokenAddr,
					Contract: linkToken,
					Tx:       tx,
					Tv:       deployment.NewTypeAndVersion(types.StaticLinkToken, deployment.Version1_0_0),
					Err:      err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy static link token", "chain", chain.String(), "err", err)
			return deployment.ChangesetOutput{}, err
		}
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}

func deployLinkTokenContractEVM(
	lggr logger.Logger,
	chain deployment.Chain,
	ab deployment.AddressBook,
) (*cldf.ContractDeploy[*link_token.LinkToken], error) {
	linkToken, err := cldf.DeployContract[*link_token.LinkToken](lggr, chain, ab,
		func(chain deployment.Chain) cldf.ContractDeploy[*link_token.LinkToken] {
			linkTokenAddr, tx, linkToken, err2 := link_token.DeployLinkToken(
				chain.DeployerKey,
				chain.Client,
			)
			return cldf.ContractDeploy[*link_token.LinkToken]{
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
	tokenAdminPubKey := chain.DeployerKey.PublicKey()
	mint, _ := solana.NewRandomPrivateKey()
	mintPublicKey := mint.PublicKey() // this is the token address
	instructions, err := solTokenUtil.CreateToken(
		context.Background(),
		solana.Token2022ProgramID,
		mintPublicKey,
		tokenAdminPubKey,
		TokenDecimalsSolana,
		chain.Client,
		cldf.SolDefaultCommitment,
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
	lggr.Infow("Deployed contract", "Contract", tv.String(), "addr", mint.PublicKey().String(), "chain", chain.String())
	err = ab.Save(chain.Selector, mint.PublicKey().String(), tv)
	if err != nil {
		lggr.Errorw("Failed to save link token", "chain", chain.String(), "err", err)
		return err
	}
	return nil
}
