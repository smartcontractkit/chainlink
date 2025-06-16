package sequence

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

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
	TokenObjAddress   aptos.AccountAddress
	TokenAddress      aptos.AccountAddress
	TokenOwnerAddress aptos.AccountAddress
	PoolType          cldf.ContractType
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

	// 1 - Cleanup staging area
	mcmsAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].MCMSAddress
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, mcmsAddress)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// 2 - Set token Registrar
	// Get a deterministic seed using token address and pool type
	tokenPoolSeed := fmt.Sprintf("%s::%s", in.TokenAddress.String(), in.PoolType.String())
	// Calculate token pool owner address and set token registrar
	mcmsContract := mcmsbind.Bind(mcmsAddress, deps.AptosChain.Client)
	tokenPoolOwnerAddress, err := mcmsContract.MCMSRegistry().GetNewCodeObjectOwnerAddress(nil, []byte(tokenPoolSeed))
	if err != nil {
		return DeployTokenPoolSeqOutput{}, fmt.Errorf("failed to get new code object owner address: %w", err)
	}
	setTokenRegistrarInput := operation.SetTokenRegistrarInput{
		TokenAddress:          in.TokenAddress,
		TokenPoolOwnerAddress: tokenPoolOwnerAddress,
	}
	setRegReport, err := operations.ExecuteOperation(b, operation.SetTokenRegistrarOp, deps, setTokenRegistrarInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  []mcmstypes.Transaction{setRegReport.Output},
	})

	// 3 - Deploy token pool package
	deployTokenPoolPackageReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolPackageOp, deps, tokenPoolSeed)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	tokenPoolObjectAddress := deployTokenPoolPackageReport.Output.TokenPoolObjectAddress
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenPoolPackageReport.Output.MCMSOps)...)

	// 4 - Deploy token pool module
	// The initial administrator of the token pool will be set to the MCMS resource account owning CCIP -
	// when calling admin function on the TAR, this signer will be used.
	initialAdministrator, err := mcmsContract.MCMSRegistry().GetRegisteredOwnerAddress(nil, deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].CCIPAddress)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, fmt.Errorf("failed to get CCIP owner address to be set as an initial administrator: %w", err)
	}
	deployTokenPoolModuleInput := operation.DeployTokenPoolModuleInput{
		TokenObjAddress:      in.TokenObjAddress,
		TokenPoolObjAddress:  tokenPoolObjectAddress,
		InitialAdministrator: initialAdministrator,
		PoolType:             in.PoolType,
	}
	deployTokenPoolModuleReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolModuleOp, deps, deployTokenPoolModuleInput)
	if err != nil {
		return DeployTokenPoolSeqOutput{}, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenPoolModuleReport.Output)...)
	// 5 - Grant BnM permission to the token pool
	// TODO: BnM Pool should also have this
	if in.PoolType == shared.AptosManagedTokenPoolType {
		// Get the token pool state address
		tokenPoolStateAddress := tokenPoolObjectAddress.ResourceAccount([]byte("CcipManagedTokenPool"))
		var txs []mcmstypes.Transaction
		gmReport, err := operations.ExecuteOperation(b, operation.GrantMinterPermissionsOp, deps, operation.GrantRolePermissionsInput{
			TokenObjAddress:       in.TokenObjAddress,
			TokenPoolStateAddress: tokenPoolStateAddress,
		})
		if err != nil {
			return DeployTokenPoolSeqOutput{}, err
		}
		txs = append(txs, gmReport.Output)

		gbReport, err := operations.ExecuteOperation(b, operation.GrantBurnerPermissionsOp, deps, operation.GrantRolePermissionsInput{
			TokenObjAddress:       in.TokenObjAddress,
			TokenPoolStateAddress: tokenPoolStateAddress,
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
	TokenPoolAddress    aptos.AccountAddress
	RemotePools         map[uint64]RemotePool
	RemotePoolsToRemove []uint64 // To re-set a pool also add its address on the removing list
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

	// Re-organize remote pool variables into contract input format
	var remoteChainSelectors []uint64
	var remotePoolAddresses [][][]byte
	var remoteTokenAddresses [][]byte
	var outboundIsEnableds []bool
	var outboundCapacities []uint64
	var outboundRates []uint64
	var inboundIsEnableds []bool
	var inboundCapacities []uint64
	var inboundRates []uint64

	for remoteSel, remotePool := range in.RemotePools {
		remoteChainSelectors = append(remoteChainSelectors, remoteSel)
		remotePoolAddresses = append(remotePoolAddresses, [][]byte{remotePool.RemotePoolAddress})
		remoteTokenAddresses = append(remoteTokenAddresses, remotePool.RemoteTokenAddress)
		outboundIsEnableds = append(outboundIsEnableds, remotePool.OutboundIsEnabled)
		outboundCapacities = append(outboundCapacities, remotePool.OutboundCapacity)
		outboundRates = append(outboundRates, remotePool.OutboundRate)
		inboundIsEnableds = append(inboundIsEnableds, remotePool.InboundIsEnabled)
		inboundCapacities = append(inboundCapacities, remotePool.InboundCapacity)
		inboundRates = append(inboundRates, remotePool.InboundRate)
	}

	// Apply chain updates
	applyChainUpdatesInput := operation.ApplyChainUpdatesInput{
		RemoteChainSelectorsToRemove: in.RemotePoolsToRemove,
		RemoteChainSelectorsToAdd:    remoteChainSelectors,
		RemotePoolAddresses:          remotePoolAddresses,
		RemoteTokenAddresses:         remoteTokenAddresses,
		TokenPoolAddress:             in.TokenPoolAddress,
	}
	applyChainUpdatesReport, err := operations.ExecuteOperation(b, operation.ApplyChainUpdatesOp, deps, applyChainUpdatesInput)
	if err != nil {
		return mcmstypes.BatchOperation{}, err
	}
	txs = append(txs, applyChainUpdatesReport.Output)

	// Set chain rate limiter configs
	if len(remoteChainSelectors) > 0 {
		setChainRateLimiterInput := operation.SetChainRLConfigsInput{
			RemoteChainSelectors: remoteChainSelectors,
			OutboundIsEnableds:   outboundIsEnableds,
			OutboundCapacities:   outboundCapacities,
			OutboundRates:        outboundRates,
			InboundIsEnableds:    inboundIsEnableds,
			InboundCapacities:    inboundCapacities,
			InboundRates:         inboundRates,
			TokenPoolAddress:     in.TokenPoolAddress,
		}
		setChainRateLimiterReport, err := operations.ExecuteOperation(b, operation.SetChainRateLimiterConfigsOp, deps, setChainRateLimiterInput)
		if err != nil {
			return mcmstypes.BatchOperation{}, err
		}
		txs = append(txs, setChainRateLimiterReport.Output)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(deps.AptosChain.Selector),
		Transactions:  txs,
	}, nil
}
