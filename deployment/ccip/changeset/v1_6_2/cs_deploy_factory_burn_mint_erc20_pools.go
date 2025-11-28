package v1_6_2

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6_2"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

var _ cldf.ChangeSet[ccipseq.DeployFactoryBurnMintERC20PoolsConfig] = DeployFactoryBurnMintERC20PoolsChangeset

// DeployPrerequisitesChangeset deploys a FactoryBurnMintERC20Token on a chain.
// It also deploys a BurnMintTokenPool, BurnFromMintTokenPool, BurnWithFromMintTokenPool,
// and a LockReleaseTokenPool for the newly deployed token. If any of these contracts are
// already deployed on the chain, they are not deployed again here.
func DeployFactoryBurnMintERC20PoolsChangeset(env cldf.Environment, c ccipseq.DeployFactoryBurnMintERC20PoolsConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid DeployFactoryBurnMintERC20PoolsConfig: %w", err)
	}
	report, err := deployFactoryBurnMintERC20PoolsForChains(env, c)
	if err != nil {
		return cldf.ChangesetOutput{
			Reports: report.ExecutionReports,
		}, fmt.Errorf("failed to deploy FactoryBurnMintERC20 and pool contracts: %w", err)
	}

	addressBook := cldf.NewMemoryAddressBook()
	for chainSel, addresses := range report.Output {
		for address, typeAndVersion := range addresses {
			err := addressBook.Save(chainSel, address, cldf.MustTypeAndVersionFromString(typeAndVersion))
			if err != nil {
				return cldf.ChangesetOutput{
					Reports: report.ExecutionReports,
				}, fmt.Errorf("failed to save address %s for chain %d: %w", address, chainSel, err)
			}
		}
	}
	ds, err := shared.PopulateDataStore(addressBook)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate in-memory DataStore: %w", err)
	}
	return cldf.ChangesetOutput{
		Reports:     report.ExecutionReports,
		AddressBook: addressBook,
		DataStore:   ds,
	}, nil
}

func deployFactoryBurnMintERC20PoolsForChains(
	e cldf.Environment,
	c ccipseq.DeployFactoryBurnMintERC20PoolsConfig,
) (operations.SequenceReport[ccipseq.DeployFactoryBurnMintERC20PoolsSeqConfig, map[uint64]map[string]string], error) {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err", err)
		return operations.SequenceReport[ccipseq.DeployFactoryBurnMintERC20PoolsSeqConfig, map[uint64]map[string]string]{}, err
	}

	addresses := make(map[uint64]ccipseq.CCIPContractAddresses)
	for chainSel, _ := range c.TokenAndPoolContractParamsPerChain {
		// get addresses
		rmnProxyAddress := common.Address{}
		routerAddress := common.Address{}
		factoryBurnMintERC20Address := common.Address{}
		burnMintTokenPoolAddress := common.Address{}
		burnFromMintTokenPoolAddress := common.Address{}
		burnWithFromMintTokenPoolAddress := common.Address{}
		lockReleaseTokenPoolAddress := common.Address{}

		chainState, chainExists := existingState.EVMChainState(chainSel)
		if !chainExists {
			return operations.SequenceReport[ccipseq.DeployFactoryBurnMintERC20PoolsSeqConfig, map[uint64]map[string]string]{}, fmt.Errorf("chain %d does not exist in state", chainSel)
		}
		if chainState.RMNProxy != nil {
			rmnProxyAddress = chainState.RMNProxy.Address()
		}
		if chainState.Router != nil {
			routerAddress = chainState.Router.Address()
		}
		if chainState.FactoryBurnMintERC20Token != nil {
			factoryBurnMintERC20Address = chainState.FactoryBurnMintERC20Token.Address()
		}
		if tokenPools, ok := chainState.BurnMintTokenPools[shared.FactoryBurnMintERC20Symbol]; ok {
			if tokenPool, ok := tokenPools[shared.CurrentTokenPoolVersion]; ok {
				burnMintTokenPoolAddress = tokenPool.Address()
			}
		}
		if tokenPools, ok := chainState.BurnFromMintTokenPools[shared.FactoryBurnMintERC20Symbol]; ok {
			if tokenPool, ok := tokenPools[shared.CurrentTokenPoolVersion]; ok {
				burnFromMintTokenPoolAddress = tokenPool.Address()
			}
		}
		if tokenPools, ok := chainState.BurnWithFromMintTokenPools[shared.FactoryBurnMintERC20Symbol]; ok {
			if tokenPool, ok := tokenPools[shared.CurrentTokenPoolVersion]; ok {
				burnWithFromMintTokenPoolAddress = tokenPool.Address()
			}
		}
		if tokenPools, ok := chainState.LockReleaseTokenPools[shared.FactoryBurnMintERC20Symbol]; ok {
			if tokenPool, ok := tokenPools[shared.CurrentTokenPoolVersion]; ok {
				lockReleaseTokenPoolAddress = tokenPool.Address()
			}
		}

		addresses[chainSel] = ccipseq.CCIPContractAddresses{
			RMNProxyAddress:                  rmnProxyAddress,
			RouterAddress:                    routerAddress,
			FactoryBurnMintERC20Address:      factoryBurnMintERC20Address,
			BurnMintTokenPoolAddress:         burnMintTokenPoolAddress,
			BurnFromMintTokenPoolAddress:     burnFromMintTokenPoolAddress,
			BurnWithFromMintTokenPoolAddress: burnWithFromMintTokenPoolAddress,
			LockReleaseTokenPoolAddress:      lockReleaseTokenPoolAddress,
		}
	}

	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipseq.DeployFactoryBurnMintERC20PoolsSeq,
		e.BlockChains.EVMChains(),
		ccipseq.DeployFactoryBurnMintERC20PoolsSeqConfig{
			AddressesPerChain:                     addresses,
			DeployFactoryBurnMintERC20PoolsConfig: c,
			GasBoostConfigPerChain:                c.GasBoostConfigPerChain,
		},
	)
	if err != nil {
		return report, fmt.Errorf("failed to deploy FactoryBurnMintERC20 and pool contracts: %w", err)
	}

	return report, nil
}
