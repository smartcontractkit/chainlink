package changeset

import (
	"fmt"
	"math/big"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// Contract-type labels used to resolve bridge dependencies from the datastore. These should match the
// `type`  of the address in the datastore and the (type, version) keys registered in the proposal ABI
// registry (treasury-cld-changesets VaultEVMAbiRegistry), otherwise lookups or proposal decoding fail.
const (
	WETHContractType                cldf.ContractType = "WETH"
	RouterContractType              cldf.ContractType = "Router"
	EtherSenderReceiverContractType cldf.ContractType = "EtherSenderReceiver"
)

// Changeset-level tunables. We can change them here to customize; they are intentionally not exposed in the YAML payload.
var (
	// TestBridgeAmount is the small amount (wei) bridged by the test proposal when TestNeeded is true.
	TestBridgeAmount = big.NewInt(1_000_000_000_000_000) // 0.001 ETH

	// CCIPFeeBufferBps is applied to the live CCIP fee quote to compute the native value sent with
	// ccipSend (15000 = 150%). NOTE: for native fee the Router consumes the whole msg.value and does
	// not refund the excess, so the buffer above the actual fee is knowingly donated.
	// This payment of additional fees is fine as indirectly its returned to Finance via the CCIP onRamps
	CCIPFeeBufferBps = big.NewInt(15_000)
)

const (
	bpsDenominator        = int64(10_000)
	defaultBridgeGasLimit = uint64(100_000)
)

// BridgeNativeChangeset bridges native/wrapped-native value held by a source chain's Timelock to one
// or more destination chains via CCIP, delivering native value to a destination address through an EtherSenderReceiver.
// It only ever produces MCMS timelock proposals.
var BridgeNativeChangeset cldf.ChangeSetV2[types.CCIPBridgeConfig] = bridgeNativeChangeset{}

type bridgeNativeChangeset struct{}

func (b bridgeNativeChangeset) VerifyPreconditions(e cldf.Environment, cfg types.CCIPBridgeConfig) error {
	return ValidateCCIPBridgeConfig(e.GetContext(), e, cfg)
}

func (b bridgeNativeChangeset) Apply(e cldf.Environment, cfg types.CCIPBridgeConfig) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	evmChains := e.BlockChains.EVMChains()
	sourceChain, ok := evmChains[cfg.FromChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("source chain %d not found in environment", cfg.FromChainSelector)
	}

	lggr.Infow("Starting CCIP bridge",
		"from", cfg.FromChainSelector,
		"token_type", cfg.TokenType,
		"destinations", len(cfg.Transfers),
		"test_mode", cfg.TestNeeded)

	deps := VaultDeps{
		Chain:       sourceChain,
		Auth:        sourceChain.DeployerKey,
		DataStore:   e.DataStore,
		Environment: e,
	}

	seqReport, err := operations.ExecuteSequence(
		e.OperationsBundle,
		BridgeCCIPSequence,
		deps,
		BridgeCCIPSequenceInput{Config: cfg},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute CCIP bridge sequence: %w", err)
	}

	lggr.Infow("CCIP bridge proposal generation completed",
		"from", cfg.FromChainSelector,
		"proposals", len(seqReport.Output.MCMSTimelockProposals))

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: seqReport.Output.MCMSTimelockProposals,
		Reports:               seqReport.ExecutionReports,
	}, nil
}
