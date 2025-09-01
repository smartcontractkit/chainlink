package sequence

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	fee_quoter "github.com/smartcontractkit/chainlink-aptos/bindings/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_token_pools/managed_token_pool"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// Deploy Token Pool sequence input
type DeployTokenPoolSeqInput struct {
	TokenCodeObjAddress aptos.AccountAddress
	TokenAddress        aptos.AccountAddress
	TokenOwnerAddress   aptos.AccountAddress
	PoolType            cldf.ContractType
}
type DeployTokenPoolSeqOutput struct {
	TokenPoolAddress aptos.AccountAddress
	MCMSOps          []mcmstypes.BatchOperation
}

// DeployAptosTokenPoolSequence deploys token pool to the same address as Token Object Address
var DeployAptosTokenPoolSequence = operations.NewSequence(
	"deploy-aptos-token-pool",
	operation.Version1_0_0,
	"Deploys token and token pool and configures",
	deployAptosTokenPoolSequence,
)

func deployAptosTokenPoolSequence(b operations.Bundle, deps operation.AptosDeps, in DeployTokenPoolSeqInput) (DeployTokenPoolSeqOutput, error) {
	var mcmsOperations []mcmstypes.BatchOperation
	var txs []mcmstypes.Transaction
	mcmsAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].MCMSAddress
	mcmsContract := mcmsbind.Bind(mcmsAddress, deps.AptosChain.Client)

	// 1 - Cleanup staging area
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, mcmsAddress)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// 2 - Deploy token pool package
	// Get a deterministic seed using token address and pool type
	tokenPoolSeed := fmt.Sprintf("%s::%s", in.TokenAddress.StringLong(), in.PoolType.String())
	deployTokenPoolPackageReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolPackageOp, deps, tokenPoolSeed)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	tokenPoolObjectAddress := deployTokenPoolPackageReport.Output.TokenPoolObjectAddress
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenPoolPackageReport.Output.MCMSOps)...)

	// 3 - Deploy token pool module
	deployTokenPoolModuleInput := operation.DeployTokenPoolModuleInput{
		TokenAddress:        in.TokenAddress,
		TokenCodeObjAddress: in.TokenCodeObjAddress,
		TokenPoolObjAddress: tokenPoolObjectAddress,
		PoolType:            in.PoolType,
	}
	deployTokenPoolModuleReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolModuleOp, deps, deployTokenPoolModuleInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenPoolModuleReport.Output)...)

	// 4 - ProposeAdministrator
	// The initial administrator of the token pool will be set to the MCMS resource account owning CCIP -
	// when calling admin function on the TAR, this signer will be used.
	initialAdministrator, err := mcmsContract.MCMSRegistry().GetRegisteredOwnerAddress(nil, deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, fmt.Errorf("failed to get CCIP owner address to be set as an initial administrator: %w", err)
	}
	proposeAdministratorIn := operation.ProposeAdministratorInput{
		TokenAddress:       in.TokenAddress,
		TokenAdministrator: initialAdministrator,
	}
	paReport, err := operations.ExecuteOperation(b, operation.ProposeAdministratorOp, deps, proposeAdministratorIn)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	txs = append(txs, paReport.Output)

	// 5 - AcceptAdminRole
	aaReport, err := operations.ExecuteOperation(b, operation.AcceptAdminRoleOp, deps, in.TokenAddress)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	txs = append(txs, aaReport.Output)

	// 6 - SetPool
	setPoolIn := operation.SetPoolInput{
		TokenAddress:     in.TokenAddress,
		TokenPoolAddress: tokenPoolObjectAddress,
	}
	spReport, err := operations.ExecuteOperation(b, operation.SetPoolOp, deps, setPoolIn)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	txs = append(txs, spReport.Output)

	// 7 - Grant BnM permission to the token pool
	// TODO: BnM Pool should also have this
	if in.PoolType == shared.AptosManagedTokenPoolType {
		// Get the token pool state address
		tokenPoolStateAddress := tokenPoolObjectAddress.ResourceAccount([]byte("CcipManagedTokenPool"))
		gmReport, err := operations.ExecuteOperation(b, operation.ApplyAllowedMintersOp, deps, operation.ApplyAllowedMintersInput{
			TokenCodeObjectAddress: in.TokenCodeObjAddress,
			MintersToAdd:           []aptos.AccountAddress{tokenPoolStateAddress},
		})
		if err != nil {
			return DeployTokenPoolSeqOutput{}, err
		}
		txs = append(txs, gmReport.Output)

		gbReport, err := operations.ExecuteOperation(b, operation.ApplyAllowedBurnersOp, deps, operation.ApplyAllowedBurnersInput{
			TokenCodeObjectAddress: in.TokenCodeObjAddress,
			BurnersToAdd:           []aptos.AccountAddress{tokenPoolStateAddress},
		})
		if err != nil {
			return DeployTokenPoolSeqOutput{}, err
		}
		txs = append(txs, gbReport.Output)

		mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
			Transactions:  txs,
		})
	}

	return DeployTokenPoolSeqOutput{
		TokenPoolAddress: tokenPoolObjectAddress,
		MCMSOps:          mcmsOperations,
	}, nil
}

// Connect Token Pool sequence input
type ConnectTokenPoolSeqInput struct {
	TokenPoolAddress                    aptos.AccountAddress
	RemotePools                         map[uint64]RemotePool
	RemotePoolsToRemove                 []uint64 // To re-set a pool also add its address on the removing list
	TokenAddress                        aptos.AccountAddress
	TokenTransferFeeByRemoteChainConfig map[uint64]fee_quoter.TokenTransferFeeConfig
}

type RemotePool struct {
	RemotePoolAddress  []byte
	RemoteTokenAddress []byte
	config.RateLimiterConfig
}

var ConnectTokenPoolSequence = operations.NewSequence(
	"connect-aptos-evm-token-pools",
	operation.Version1_0_0,
	"Connects EVM<>Aptos lanes token pools",
	connectTokenPoolSequence,
)

func connectTokenPoolSequence(b operations.Bundle, deps operation.AptosDeps, in ConnectTokenPoolSeqInput) (mcmstypes.BatchOperation, error) {
	var txs []mcmstypes.Transaction

	// Chain updates
	applyChainUpdatesInput := operation.ApplyChainUpdatesInput{
		TokenPoolAddress:             in.TokenPoolAddress,
		RemoteChainSelectorsToRemove: in.RemotePoolsToRemove,
		RemoteChainSelectorsToAdd:    nil,
		RemotePoolAddresses:          nil,
		RemoteTokenAddresses:         nil,
	}

	// Remote Pool Adds
	addRemotePoolsInput := operation.AddRemotePoolsInput{
		TokenPoolAddress:     in.TokenPoolAddress,
		RemoteChainSelectors: nil,
		RemotePoolAddresses:  nil,
	}

	// Update rate limits
	setChainRLConfigsInput := operation.SetChainRLConfigsInput{
		TokenPoolAddress:     in.TokenPoolAddress,
		RemoteChainSelectors: nil,
		OutboundIsEnableds:   nil,
		OutboundCapacities:   nil,
		OutboundRates:        nil,
		InboundIsEnableds:    nil,
		InboundCapacities:    nil,
		InboundRates:         nil,
	}

	tokenPool := managed_token_pool.Bind(in.TokenPoolAddress, deps.AptosChain.Client)
	supportedChains, err := tokenPool.ManagedTokenPool().GetSupportedChains(nil)
	if err != nil {
		b.Logger.Debugf("failed to get supported chains from token pool %s, likely because it isn't deployed yet: %v", in.TokenPoolAddress.StringLong(), err)
	}
	for remoteSel, remotePool := range in.RemotePools {
		// Always apply rate limits
		setChainRLConfigsInput.RemoteChainSelectors = append(setChainRLConfigsInput.RemoteChainSelectors, remoteSel)
		setChainRLConfigsInput.OutboundIsEnableds = append(setChainRLConfigsInput.OutboundIsEnableds, remotePool.OutboundIsEnabled)
		setChainRLConfigsInput.OutboundCapacities = append(setChainRLConfigsInput.OutboundCapacities, remotePool.OutboundCapacity)
		setChainRLConfigsInput.OutboundRates = append(setChainRLConfigsInput.OutboundRates, remotePool.OutboundRate)
		setChainRLConfigsInput.InboundIsEnableds = append(setChainRLConfigsInput.InboundIsEnableds, remotePool.InboundIsEnabled)
		setChainRLConfigsInput.InboundCapacities = append(setChainRLConfigsInput.InboundCapacities, remotePool.InboundCapacity)
		setChainRLConfigsInput.InboundRates = append(setChainRLConfigsInput.InboundRates, remotePool.InboundRate)

		isSupportedChain := slices.Contains(supportedChains, remoteSel)
		if !isSupportedChain {
			// Only add the remote chain if it isn't supported yet
			applyChainUpdatesInput.RemoteChainSelectorsToAdd = append(applyChainUpdatesInput.RemoteChainSelectorsToAdd, remoteSel)
			applyChainUpdatesInput.RemotePoolAddresses = append(applyChainUpdatesInput.RemotePoolAddresses, [][]byte{remotePool.RemotePoolAddress})
			applyChainUpdatesInput.RemoteTokenAddresses = append(applyChainUpdatesInput.RemoteTokenAddresses, remotePool.RemoteTokenAddress)
		} else {
			// If the chain is supported, check if there's an updated remote pool that hasn't been configured yet
			configuredRemotePools, err := tokenPool.ManagedTokenPool().GetRemotePools(nil, remoteSel)
			if err != nil {
				return mcmstypes.BatchOperation{}, fmt.Errorf("failed to get remote pools from token pool for selector %d: %w", remoteSel, err)
			}
			isRemotePoolSupported := false
			for _, configuredRemotePool := range configuredRemotePools {
				if bytes.Equal(configuredRemotePool, remotePool.RemotePoolAddress) {
					isRemotePoolSupported = true
					break
				}
			}
			if !isRemotePoolSupported {
				addRemotePoolsInput.RemoteChainSelectors = append(addRemotePoolsInput.RemoteChainSelectors, remoteSel)
				addRemotePoolsInput.RemotePoolAddresses = append(addRemotePoolsInput.RemotePoolAddresses, remotePool.RemotePoolAddress)
			}
		}
	}

	// Apply chain updates if there are any
	if (len(applyChainUpdatesInput.RemoteChainSelectorsToAdd) + len(applyChainUpdatesInput.RemoteChainSelectorsToRemove)) > 0 {
		applyChainUpdatesReport, err := operations.ExecuteOperation(b, operation.ApplyChainUpdatesOp, deps, applyChainUpdatesInput)
		if err != nil {
			return mcmstypes.BatchOperation{}, err
		}
		txs = append(txs, applyChainUpdatesReport.Output)
	}

	// Add remote pools if there are any to apply
	if len(addRemotePoolsInput.RemoteChainSelectors) > 0 {
		addRemotePoolsReport, err := operations.ExecuteOperation(b, operation.AddRemotePoolsOp, deps, addRemotePoolsInput)
		if err != nil {
			return mcmstypes.BatchOperation{}, err
		}
		txs = append(txs, addRemotePoolsReport.Output...)
	}

	// Set chain rate limiter configs
	if len(setChainRLConfigsInput.RemoteChainSelectors) > 0 {
		setChainRateLimiterReport, err := operations.ExecuteOperation(b, operation.SetChainRateLimiterConfigsOp, deps, setChainRLConfigsInput)
		if err != nil {
			return mcmstypes.BatchOperation{}, err
		}
		txs = append(txs, setChainRateLimiterReport.Output)
	}

	// Apply token transfer fee configuration updates
	for destSelector, feeConfig := range in.TokenTransferFeeByRemoteChainConfig {
		applyTokenTransferFeeCfgInput := operation.ApplyTokenTransferFeeCfgInput{
			DestChainSelector: destSelector,
			ConfigsByToken:    map[string]fee_quoter.TokenTransferFeeConfig{in.TokenAddress.StringLong(): feeConfig},
		}
		applyTokenTransferFeeCfgReport, err := operations.ExecuteOperation(b, operation.ApplyTokenTransferFeeCfgOp, deps, applyTokenTransferFeeCfgInput)
		if err != nil {
			return mcmstypes.BatchOperation{}, err
		}
		txs = append(txs, applyTokenTransferFeeCfgReport.Output...)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
