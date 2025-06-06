package changeset

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	eth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	"golang.org/x/sync/errgroup"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solCommonUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	solTokenUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var _ cldf.ChangeSet[[]uint64] = DeployLinkToken

// DeployLinkToken deploys a link token contract to the chain identified by the ChainSelector.
func DeployLinkToken(e cldf.Environment, chains []uint64) (cldf.ChangesetOutput, error) {
	err := deployment.ValidateSelectorsInEnvironment(e, chains)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	newAddresses := cldf.NewMemoryAddressBook()
	deployGrp := errgroup.Group{}
	for _, chain := range chains {
		family, err := chainsel.GetSelectorFamily(chain)
		if err != nil {
			return cldf.ChangesetOutput{AddressBook: newAddresses}, err
		}
		var deployFn func() error
		switch family {
		case chainsel.FamilyEVM:
			// Deploy EVM LINK token
			deployFn = func() error {
				_, err := deployLinkTokenContractEVM(
					e.Logger, e.BlockChains.EVMChains()[chain], newAddresses,
				)
				return err
			}
		default:
			return cldf.ChangesetOutput{}, fmt.Errorf("unsupported chain family %s", family)
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
	return cldf.ChangesetOutput{AddressBook: newAddresses}, deployGrp.Wait()
}

// DeployStaticLinkToken deploys a static link token contract to the chain identified by the ChainSelector.
func DeployStaticLinkToken(e cldf.Environment, chains []uint64) (cldf.ChangesetOutput, error) {
	err := deployment.ValidateSelectorsInEnvironment(e, chains)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	newAddresses := cldf.NewMemoryAddressBook()
	for _, chainSel := range chains {
		chain, ok := e.BlockChains.EVMChains()[chainSel]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain not found in environment: %d", chainSel)
		}
		_, err := cldf.DeployContract[*link_token_interface.LinkToken](e.Logger, chain, newAddresses,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*link_token_interface.LinkToken] {
				linkTokenAddr, tx, linkToken, err2 := link_token_interface.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
				return cldf.ContractDeploy[*link_token_interface.LinkToken]{
					Address:  linkTokenAddr,
					Contract: linkToken,
					Tx:       tx,
					Tv:       cldf.NewTypeAndVersion(types.StaticLinkToken, deployment.Version1_0_0),
					Err:      err2,
				}
			})
		if err != nil {
			e.Logger.Errorw("Failed to deploy static link token", "chain", chain.String(), "err", err)
			return cldf.ChangesetOutput{}, err
		}
	}
	return cldf.ChangesetOutput{AddressBook: newAddresses}, nil
}

func deployLinkTokenContractEVM(
	lggr logger.Logger,
	chain cldf_evm.Chain,
	ab cldf.AddressBook,
) (*cldf.ContractDeploy[*link_token.LinkToken], error) {
	linkToken, err := cldf.DeployContract[*link_token.LinkToken](lggr, chain, ab,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*link_token.LinkToken] {
			var (
				linkTokenAddr common.Address
				tx            *eth_types.Transaction
				linkToken     *link_token.LinkToken
				err2          error
			)
			if !chain.IsZkSyncVM {
				linkTokenAddr, tx, linkToken, err2 = link_token.DeployLinkToken(
					chain.DeployerKey,
					chain.Client,
				)
			} else {
				linkTokenAddr, _, linkToken, err2 = link_token.DeployLinkTokenZk(
					nil,
					chain.ClientZkSyncVM,
					chain.DeployerKeyZkSyncVM,
					chain.Client,
				)
			}
			return cldf.ContractDeploy[*link_token.LinkToken]{
				Address:  linkTokenAddr,
				Contract: linkToken,
				Tx:       tx,
				Tv:       cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0),
				Err:      err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy link token", "chain", chain.String(), "err", err)
		return linkToken, err
	}
	return linkToken, nil
}

type DeploySolanaLinkTokenConfig struct {
	ChainSelector uint64
	TokenPrivKey  solana.PrivateKey
	TokenDecimals uint8
}

func DeploySolanaLinkToken(e cldf.Environment, cfg DeploySolanaLinkTokenConfig) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.SolanaChains()[cfg.ChainSelector]
	mint := cfg.TokenPrivKey
	instructions, err := solTokenUtil.CreateToken(
		context.Background(),
		solana.TokenProgramID,
		mint.PublicKey(),
		chain.DeployerKey.PublicKey(),
		cfg.TokenDecimals,
		chain.Client,
		cldf_solana.SolDefaultCommitment,
	)
	if err != nil {
		e.Logger.Errorw("Failed to generate instructions for link token deployment", "chain", chain.String(), "err", err)
		return cldf.ChangesetOutput{}, err
	}
	err = chain.Confirm(instructions, solCommonUtil.AddSigners(mint))
	if err != nil {
		e.Logger.Errorw("Failed to confirm instructions for link token deployment", "chain", chain.String(), "err", err)
		return cldf.ChangesetOutput{}, err
	}
	tv := cldf.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0)
	e.Logger.Infow("Deployed contract", "Contract", tv.String(), "addr", mint.PublicKey().String(), "chain", chain.String())
	newAddresses := cldf.NewMemoryAddressBook()
	err = newAddresses.Save(chain.Selector, mint.PublicKey().String(), tv)
	if err != nil {
		e.Logger.Errorw("Failed to save link token", "chain", chain.String(), "err", err)
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}
