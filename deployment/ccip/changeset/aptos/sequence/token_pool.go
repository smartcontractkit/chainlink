package sequence

import (
	"github.com/aptos-labs/aptos-go-sdk"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/operation"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// Deploy Token Pool sequence input
type DeployTokenPoolSeqInput struct {
	TokenObjAddress aptos.AccountAddress
	PoolType        cldf.ContractType
}

// DeployAptosTokenPoolSequence deploys token pool to the same address as Token Object Address
var DeployAptosTokenPoolSequence = operations.NewSequence(
	"deploy-aptos-token-pool",
	operation.Version1_0_0,
	"Deploys token and token pool and configures",
	deployAptosTokenPoolSequence,
)

func deployAptosTokenPoolSequence(b operations.Bundle, deps operation.AptosDeps, in DeployTokenPoolSeqInput) ([]mcmstypes.BatchOperation, error) {
	mcmsOperations := []mcmstypes.BatchOperation{}
	txs := []mcmstypes.Transaction{}

	// Cleanup staging area
	mcmsAddress := deps.CCIPOnChainState.AptosChains[deps.AptosChain.Selector].MCMSAddress
	cleanupReport, err := operations.ExecuteOperation(b, operation.CleanupStagingAreaOp, deps, mcmsAddress)
	if err != nil {
		return nil, err
	}
	if len(cleanupReport.Output.Transactions) > 0 {
		mcmsOperations = append(mcmsOperations, cleanupReport.Output)
	}

	// Deploy token pool package
	deployTokenPoolPackageReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolPackageOp, deps, in.TokenObjAddress)
	if err != nil {
		return nil, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenPoolPackageReport.Output)...)

	// Deploy token pool module
	deployTokenPoolModuleInput := operation.DeployTokenPoolModuleInput{
		TokenObjAddress: in.TokenObjAddress,
		PoolType:        in.PoolType,
	}
	deployTokenPoolModuleReport, err := operations.ExecuteOperation(b, operation.DeployTokenPoolModuleOp, deps, deployTokenPoolModuleInput)
	if err != nil {
		return nil, err
	}
	mcmsOperations = append(mcmsOperations, utils.ToBatchOperations(deployTokenPoolModuleReport.Output)...)

	// Grant minter permission to the token pool
	gmReport, err := operations.ExecuteOperation(b, operation.GrantMinterPermissionsOp, deps, in.TokenObjAddress)
	if err != nil {
		return nil, err
	}
	txs = append(txs, gmReport.Output)

	// Grant burner permission to the token pool
	gbReport, err := operations.ExecuteOperation(b, operation.GrantBurnerPermissionsOp, deps, in.TokenObjAddress)
	if err != nil {
		return nil, err
	}
	txs = append(txs, gbReport.Output)

	return mcmsOperations, nil
}

// Connect Token Pool sequence input
type ConnectTokenPoolSeqInput struct {
	TokenPoolAddress    aptos.AccountAddress
	RemotePools         map[uint64]RemotePool
	RemotePoolsToRemove []uint64 // To re-set a pool also add it's address on the removing list
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
