package testhelpers

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	tonaddress "github.com/xssnick/tonutils-go/address"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	tonstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/ton"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

type TonTestDeployPrerequisitesChangeSet struct {
	T                 *testing.T
	TonChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestDeployPrerequisitesChangeSet{}

func (c TonTestDeployPrerequisitesChangeSet) Apply(e cldf.Environment) (cldf.ChangesetOutput, error) {
	t := c.T

	tonChains, err := tonstate.LoadOnchainState(e)
	require.NoError(t, err)

	t.Logf("DEBUG: TonTestDeployPrerequisitesChangeSet: chain selectors: %+v\n", c.TonChainSelectors)
	for _, chainSelector := range c.TonChainSelectors {
		tonChainState := tonChains[chainSelector]

		// TODO: replace with actual TON addresses after contracts are supported, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
		address := tonaddress.MustParseAddr("EQDtFpEwcFAEcRe5mLVh2N6C0x-_hJEM7W61_JLnSF74p4q2")
		tonChainState.OffRamp = *address
		address = tonaddress.MustParseAddr("UQCfQRaJr2vxgZr5NHc0CTx6tAb0jverj9QQFirNfoCkGcUy")
		tonChainState.Router = *address
		address = tonaddress.MustParseAddr("EQADa3W6G0nSiTV4a6euRA42fU9QxSEnb-WeDpcrtWzA2jM8")
		tonChainState.LinkTokenAddress = *address
		address = tonaddress.MustParseAddr("UQDgFwiokL1ojVwXa3Ac7xCLfGB0Ti0foSw5NZ48Aj_vhs_6")
		tonChainState.CCIPAddress = *address
		address = tonaddress.MustParseAddr("UQCk4967vNM_V46Dn8I0x-gB_QE2KkdW1GQ7mWz1DtYGLEd8")
		tonChainState.ReceiverAddress = *address

		err = tonstate.SaveOnchainState(chainSelector, tonChainState, e)
		require.NoError(t, err)
	}
	return cldf.ChangesetOutput{}, nil
}

type TonTestDeployContractsChangeSet struct {
	T                 *testing.T
	HomeChainSelector uint64
	TonChainSelectors []uint64
	AllChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestDeployContractsChangeSet{}

func (c TonTestDeployContractsChangeSet) Apply(e cldf.Environment) (cldf.ChangesetOutput, error) {
	t := c.T

	tonChains, err := tonstate.LoadOnchainState(e)
	require.NoError(t, err)

	t.Logf("DEBUG: TonTestDeployContractsChangeSet: chain selectors: %+v\n", c.TonChainSelectors)

	for _, chainSelector := range c.TonChainSelectors {
		tonChain := e.BlockChains.TonChains()[chainSelector]
		tonChainState := tonChains[chainSelector]
		c.deployTonContracts(t, e, chainSelector, tonChain, tonChainState, tonChains)
	}
	return cldf.ChangesetOutput{}, nil
}

func (c TonTestDeployContractsChangeSet) deployTonContracts(t *testing.T, e cldf.Environment, chainSelector uint64, tonChain cldf_ton.Chain, tonChainState tonstate.CCIPChainState, onchainState map[uint64]tonstate.CCIPChainState) {
	deployer := tonChain.Wallet
	_ = deployer // TODO: use deployer in the rest of the code, pre-funded
	_ = onchainState
	// TODO(ton): once all contracts are deployed, we can remove the hardcoded addresses from the TonTestDeployPrerequisitesChangeSet
	// TODO(ton): Deploy TON MCMS, https://smartcontract-it.atlassian.net/browse/NONEVM-1939
	// TODO(ton): Deploy TON CCIP, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
	// TODO(ton): Deploy TON CCIP Offramp, Router, Onramp, Dummy Receiver and set the contract address
	// TODO(ton): Initialize Onramp, Offramp, FeeQuoter, RMNRemote

	err := tonstate.SaveOnchainState(chainSelector, tonChainState, e)
	require.NoError(t, err)
}

type TonTestConfigureContractsChangeSet struct {
	T                 *testing.T
	HomeChainSelector uint64
	TonChainSelectors []uint64
	AllChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestConfigureContractsChangeSet{}

func (c TonTestConfigureContractsChangeSet) Apply(e cldf.Environment) (cldf.ChangesetOutput, error) {
	t := c.T

	tonChains, err := tonstate.LoadOnchainState(e)
	require.NoError(t, err)

	t.Logf("DEBUG: TonTestConfigureContractsChangeSet: chain selectors: %+v\n", c.TonChainSelectors)

	for _, chainSelector := range c.TonChainSelectors {
		tonChain := e.BlockChains.TonChains()[chainSelector]
		tonChainState := tonChains[chainSelector]
		c.configureTonContracts(t, e, chainSelector, tonChain, tonChainState, tonChains)
	}
	return cldf.ChangesetOutput{}, nil
}

func (c TonTestConfigureContractsChangeSet) configureTonContracts(t *testing.T, e cldf.Environment, chainSelector uint64, tonChain cldf_ton.Chain, tonChainState tonstate.CCIPChainState, onchainState map[uint64]tonstate.CCIPChainState) {
	t.Logf("DEBUG: TonTestConfigureContractsChangeSet: chain selector: %d", chainSelector)
	// logger := logger.Test(t)

	// offrampBindings := ccip_offramp.Bind(tonChainState.CCIPAddress, tonChain.Client)

	// transactOpts := &bind.TransactOpts{
	// 	Signer: tonChain.DeployerSigner,
	// }

	// donID, err := internal.DonIDForChain(
	// 	onchainState.Chains[c.HomeChainSelector].CapabilityRegistry,
	// 	onchainState.Chains[c.HomeChainSelector].CCIPHome,
	// 	chainSelector,
	// )
	// require.NoError(t, err)
	// t.Logf("TON DON ID: %+v\n", donID)

	// allCommitConfigs, err := onchainState.Chains[c.HomeChainSelector].CCIPHome.GetAllConfigs(&ethbind.CallOpts{
	// 	Context: context.Background(),
	// }, donID, 0)

	// allExecConfigs, err := onchainState.Chains[c.HomeChainSelector].CCIPHome.GetAllConfigs(&ethbind.CallOpts{
	// 	Context: context.Background(),
	// }, donID, 1)

	// t.Logf("DEBUG HOME CHAIN CCIPHome: commit configs: %+v\n", allCommitConfigs)
	// t.Logf("DEBUG HOME CHAIN CCIPHome: exec configs: %+v\n", allExecConfigs)

	// ocr3Args, err := internal.BuildSetOCR3ConfigArgsTon(
	// 	donID, onchainState.Chains[c.HomeChainSelector].CCIPHome, chainSelector, globals.ConfigTypeActive)
	// require.NoError(t, err)

	// var commitArgs *internal.MultiOCR3BaseOCRConfigArgsTon = nil
	// var execArgs *internal.MultiOCR3BaseOCRConfigArgsTon = nil
	// for _, ocr3Arg := range ocr3Args {
	// 	if ocr3Arg.OcrPluginType == uint8(types.PluginTypeCCIPCommit) {
	// 		commitArgs = &ocr3Arg
	// 	} else if ocr3Arg.OcrPluginType == uint8(types.PluginTypeCCIPExec) {
	// 		execArgs = &ocr3Arg
	// 	} else {
	// 		t.Fatalf("unexpected ocr3 plugin type %s", ocr3Arg.OcrPluginType)
	// 	}
	// }
	// require.NotNil(t, commitArgs)
	// require.NotNil(t, execArgs)

	// commitSigners := [][]byte{}
	// for _, signer := range commitArgs.Signers {
	// 	commitSigners = append(commitSigners, signer)
	// 	t.Logf("DEBUG: configureTonContracts commit signer %x\n", signer)
	// }
	// commitTransmitters := []tonaddress.Address{}
	// for _, transmitter := range commitArgs.Transmitters {
	// 	address, err := utils.PublicKeyBytesToAddress(transmitter)
	// 	require.NoError(t, err)
	// 	commitTransmitters = append(commitTransmitters, address)
	// }
	// pendingTx, err := offrampBindings.Offramp().SetOcr3Config(transactOpts, commitArgs.ConfigDigest[:], uint8(types.PluginTypeCCIPCommit), commitArgs.F, commitArgs.IsSignatureVerificationEnabled, commitSigners, commitTransmitters)
	// require.NoError(t, err)
	// waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	// execSigners := [][]byte{}
	// for _, signer := range execArgs.Signers {
	// 	execSigners = append(execSigners, signer)
	// }
	// execTransmitters := []tonaddress.Address{}
	// for _, transmitter := range execArgs.Transmitters {
	// 	address, err := utils.PublicKeyBytesToAddress(transmitter)
	// 	require.NoError(t, err)
	// 	execTransmitters = append(execTransmitters, address)
	// }
	// pendingTx, err = offrampBindings.Offramp().SetOcr3Config(transactOpts, execArgs.ConfigDigest[:], uint8(types.PluginTypeCCIPExec), execArgs.F, execArgs.IsSignatureVerificationEnabled, execSigners, execTransmitters)
	// require.NoError(t, err)
	// waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	// logger.Infow("TON contracts configured")

	// for _, transmitter := range append(commitTransmitters, execTransmitters...) {
	// 	// 10 TON
	// 	entryFunction, err := ton.CoinTransferPayload(nil, transmitter, 1000000000)
	// 	require.NoError(t, err)

	// 	rawTxn, err := tonChain.Client.BuildTransaction(tonChain.DeployerSigner.AccountAddress(), ton.TransactionPayload{Payload: entryFunction})
	// 	require.NoError(t, err)

	// 	signedTxn, err := rawTxn.SignedTransaction(tonChain.DeployerSigner)
	// 	require.NoError(t, err)

	// 	submitResult, err := tonChain.Client.SubmitTransaction(signedTxn)
	// 	require.NoError(t, err)

	// 	waitForTx(t, tonChain.Client, submitResult.Hash, time.Minute*1)

	// 	t.Logf("Sent 10 TON to transmitter %s\n", transmitter.String())
	// }
}

func addLaneTonChangesets(t *testing.T, e *DeployedEnv, from, to uint64, fromFamily, toFamily string) []commoncs.ConfiguredChangeSet {
	t.Logf("DEBUG: addLaneTonChangesets %d to %d / %s to %s\n", from, to, fromFamily, toFamily)
	return []commoncs.ConfiguredChangeSet{TonTestAddLaneChangeSet{
		T:                 t,
		fromChainSelector: from,
		toChainSelector:   to,
		fromFamily:        fromFamily,
		toFamily:          toFamily,
	}}
}

type TonTestAddLaneChangeSet struct {
	T                 *testing.T
	fromChainSelector uint64
	toChainSelector   uint64
	fromFamily        string
	toFamily          string
}

var _ commoncs.ConfiguredChangeSet = TonTestAddLaneChangeSet{}

func (c TonTestAddLaneChangeSet) Apply(e cldf.Environment) (cldf.ChangesetOutput, error) {
	// t := c.T
	// TODO: support other paths
	// require.Equal(t, c.fromFamily, chainsel.FamilyEVM, "must be from EVM")
	// require.Equal(t, c.toFamily, chainsel.FamilyTon, "must be to TON")

	// tonSelector := c.toChainSelector
	// tonChain := e.TonChains[tonSelector]

	// onchainState, err := changeset.LoadOnchainState(e)
	// require.NoError(t, err)
	// tonChainState := onchainState.TonChains[tonSelector]

	// t.Logf("DEBUG: TonTestAddLaneChangeSet: LINK token: %s CCIP: %s Receiver: %s\n", tonChainState.LinkTokenAddress.String(), tonChainState.CCIPAddress.String(), tonChainState.ReceiverAddress.String())

	// require.False(t, tonChainState.LinkTokenAddress.IsAddrNone(), "LINK token address must be set")
	// require.False(t, tonChainState.CCIPAddress.IsAddrNone(), "CCIP address must be set")
	// require.False(t, tonChainState.ReceiverAddress.IsAddrNone(), "Receiver address must be set")

	// offrampBindings := ccip_offramp.Bind(tonChainState.CCIPAddress, tonChain.Client)
	// transactOpts := &bind.TransactOpts{
	// 	Signer: tonChain.DeployerSigner,
	// }

	// evmChainState := onchainState.Chains[c.fromChainSelector]
	// evmOnrampAddress := evmChainState.OnRamp.Address().Bytes()
	// t.Logf("DEBUG: TonTestAddLaneChangeSet: EVM chain selector: %d - EVM onramp address: %s\n", c.fromChainSelector, hex.EncodeToString(evmOnrampAddress))

	// sourceChainSelectors := []uint64{c.fromChainSelector}
	// sourceChainsIsEnabled := []bool{true}
	// sourceChainsIsRMNVerificationDisabled := []bool{true}
	// sourceChainsOnramps := [][]byte{evmOnrampAddress}
	// pendingTx, err := offrampBindings.Offramp().ApplySourceChainConfigUpdates(transactOpts, sourceChainSelectors, sourceChainsIsEnabled, sourceChainsIsRMNVerificationDisabled, sourceChainsOnramps)
	// require.NoError(t, err)
	// waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	// t.Logf("DEBUG: TonTestAddLaneChangeSet: Configured offramp\n")

	return cldf.ChangesetOutput{}, nil
}

func SendRequestTon(
	t *testing.T,
	e deployment.Environment,
	state tonstate.CCIPChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) { // TODO: chain independent return vailue
	// sourceSelector := cfg.SourceChain
	// destSelector := cfg.DestChain
	// msg := cfg.Message.(Ton2AnyMessage)
	return nil, errors.New("TODO(ton): SendRequestTon")
}

func ConfirmCommitWithExpectedSeqNumRangeTon(
	t *testing.T,
	srcSelector uint64,
	dest cldf_ton.Chain,
	offrampAddr tonaddress.Address,
	startBlock uint64,
	expectedSeqNumRange ccipocr3.SeqNumRange,
	enforceSingleCommit bool,
) (any, error) {
	t.Logf("DEBUG: ConfirmCommitWithExpectedSeqNumRangeTon srcSelector: %d, startBlock: %+v, expectedSeqNumRange: %+v, enforceSingleCommit: %+v\n", srcSelector, startBlock, expectedSeqNumRange, enforceSingleCommit)
	// TODO once offramp contracts are supported, we can add the logic to confirm commit with expected sequence number range
	return true, nil
}

func ConfirmExecWithSeqNrsTon(
	t *testing.T,
	sourceChain uint64,
	dest cldf_ton.Chain,
	offRampAddress tonaddress.Address,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	t.Logf("DEBUG: ConfirmExecWithSeqNrsTon srcSelector: %d, dest: %s, startBlock: %+v, expectedSeqNrs: %+v\n", sourceChain, startBlock, expectedSeqNrs)
	// TODO once offramp contracts are supported, we can add the logic to confirm execution with sequence numbers
	t.Logf("DEBUG: TODO(ton): ConfirmExecWithSeqNrsTon\n")
	seqNrsToWatch := make(map[uint64]int)
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = EXECUTION_STATE_SUCCESS
	}
	return seqNrsToWatch, nil
}
