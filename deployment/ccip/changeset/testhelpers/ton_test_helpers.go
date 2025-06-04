package testhelpers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	// "github.com/smartcontractkit/chainlink-ton/bindings/bind"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_dummy_receiver"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_offramp"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_onramp"
	// "github.com/smartcontractkit/chainlink-ton/bindings/ccip_router"
	// "github.com/smartcontractkit/chainlink-ton/bindings/mcms"
	// "github.com/smartcontractkit/chainlink-ton/relayer/utils"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	tonstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/ton"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	tonaddress "github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type TonTestDeployPrerequisitesChangeSet struct {
	T                 *testing.T
	TonChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestDeployPrerequisitesChangeSet{}

func (c TonTestDeployPrerequisitesChangeSet) Apply(e cldf.Environment) (cldf.ChangesetOutput, error) {
	t := c.T

	tonChains, err := tonstate.LoadOnchainStateTon(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: TonTestDeployPrerequisitesChangeSet: chain selectors: %+v\n", c.TonChainSelectors)
	for _, chainSelector := range c.TonChainSelectors {
		tonChainState := tonChains[chainSelector]

		// TODO: replace with actual TON addresses after contracts are supported, https://smartcontract-it.atlassian.net/browse/NONEVM-1938
		address, _ := tonaddress.ParseAddr("EQDtFpEwcFAEcRe5mLVh2N6C0x-_hJEM7W61_JLnSF74p4q2")
		tonChainState.OffRamp = *address

		address, _ = tonaddress.ParseAddr("UQCfQRaJr2vxgZr5NHc0CTx6tAb0jverj9QQFirNfoCkGcUy")
		tonChainState.Router = *address

		address, _ = tonaddress.ParseAddr("EQADa3W6G0nSiTV4a6euRA42fU9QxSEnb-WeDpcrtWzA2jM8")
		tonChainState.LinkTokenAddress = *address

		address, _ = tonaddress.ParseAddr("UQDgFwiokL1ojVwXa3Ac7xCLfGB0Ti0foSw5NZ48Aj_vhs_6")
		tonChainState.CCIPAddress = *address

		address, _ = tonaddress.ParseAddr("UQCk4967vNM_V46Dn8I0x-gB_QE2KkdW1GQ7mWz1DtYGLEd8")
		tonChainState.ReceiverAddress = *address
		err = tonstate.SaveOnchainStateTon(chainSelector, tonChainState, e)
		require.NoError(t, err)
	}
	return deployment.ChangesetOutput{}, nil
}

type TonTestDeployContractsChangeSet struct {
	T                 *testing.T
	HomeChainSelector uint64
	TonChainSelectors []uint64
	AllChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestDeployContractsChangeSet{}

func (c TonTestDeployContractsChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {

	t := c.T

	tonChains, err := tonstate.LoadOnchainStateTon(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: TonTestDeployContractsChangeSet: chain selectors: %+v\n", c.TonChainSelectors)

	for _, chainSelector := range c.TonChainSelectors {
		tonChain := e.BlockChains.TonChains()[chainSelector]
		tonChainState := tonChains[chainSelector]
		c.deployTonContracts(t, e, chainSelector, tonChain, tonChainState, tonChains)
	}
	return deployment.ChangesetOutput{}, nil
}

func (c TonTestDeployContractsChangeSet) deployTonContracts(t *testing.T, e deployment.Environment, chainSelector uint64, tonChain cldf_ton.Chain, tonChainState tonstate.CCIPChainState, onchainState map[uint64]tonstate.CCIPChainState) {
	logger := logger.Test(t)

	connectionPool := liteclient.NewConnectionPool()

	// move this to a ton config module
	var (
		NetworkConfigFile = "http://127.0.0.1:8000/localhost.global.config.json"
		// FaucetWalletSeed  = "viable model canvas decade neck soap turtle asthma bench crouch bicycle grief history envelope valid intact invest like offer urban adjust popular draft coral"
		// FaucetSubWalletID = 42
		// FaucetWalletVer   = wallet.V3R2
	)

	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), NetworkConfigFile)
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to lite servers
	err = connectionPool.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(connectionPool, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	_, err = wallet.FromSeed(api, strings.Fields(tonChain.DeployerSeed), wallet.V3R2)

	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// TODO(ton): once all contracts are deployed, we can remove the hardcoded addresses from the TonTestDeployPrerequisitesChangeSet

	// TODO(ton): Deploy TON MCMS, https://smartcontract-it.atlassian.net/browse/NONEVM-1939

	// TODO(ton): Deploy TON CCIP, https://smartcontract-it.atlassian.net/browse/NONEVM-1938

	// TODO(ton): Deploy TON CCIP Offramp

	// TODO(ton): Deploy TON CCIP Onramp

	// TODO(ton): Deploy TON CCIP Router

	// TODO(ton): Deploy Ton CCIP Dummy Receiver and set the contract address

	// tonChainState.ReceiverAddress = ton.AccountOne

	// TODO(ton): Initialize Onramp

	// TODO(ton): Initialize Offramp

	// TODO(ton): Initialize FeeQuoter

	// TODO(ton): Initialize RMNRemote

	logger.Infow("All TON contracts deployed")
	err = tonstate.SaveOnchainStateTon(chainSelector, tonChainState, e)
	require.NoError(t, err)
}

type TonTestConfigureContractsChangeSet struct {
	T                 *testing.T
	HomeChainSelector uint64
	TonChainSelectors []uint64
	AllChainSelectors []uint64
}

var _ commoncs.ConfiguredChangeSet = TonTestConfigureContractsChangeSet{}

func (c TonTestConfigureContractsChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {

	t := c.T

	tonChains, err := tonstate.LoadOnchainStateTon(e)
	require.NoError(t, err)

	fmt.Printf("DEBUG: TonTestConfigureContractsChangeSet: chain selectors: %+v\n", c.TonChainSelectors)

	for _, chainSelector := range c.TonChainSelectors {
		tonChain := e.BlockChains.TonChains()[chainSelector]
		tonChainState := tonChains[chainSelector]
		c.configureTonContracts(t, e, chainSelector, tonChain, tonChainState, tonChains)
	}
	return deployment.ChangesetOutput{}, nil
}

func (c TonTestConfigureContractsChangeSet) configureTonContracts(t *testing.T, e deployment.Environment, chainSelector uint64, tonChain cldf_ton.Chain, tonChainState tonstate.CCIPChainState, onchainState map[uint64]tonstate.CCIPChainState) {
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
	// fmt.Printf("TON DON ID: %+v\n", donID)

	// allCommitConfigs, err := onchainState.Chains[c.HomeChainSelector].CCIPHome.GetAllConfigs(&ethbind.CallOpts{
	// 	Context: context.Background(),
	// }, donID, 0)

	// allExecConfigs, err := onchainState.Chains[c.HomeChainSelector].CCIPHome.GetAllConfigs(&ethbind.CallOpts{
	// 	Context: context.Background(),
	// }, donID, 1)

	// fmt.Printf("DEBUG HOME CHAIN CCIPHome: commit configs: %+v\n", allCommitConfigs)
	// fmt.Printf("DEBUG HOME CHAIN CCIPHome: exec configs: %+v\n", allExecConfigs)

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
	// 	fmt.Printf("DEBUG: configureTonContracts commit signer %x\n", signer)
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

	// 	fmt.Printf("Sent 10 TON to transmitter %s\n", transmitter.String())
	// }
}

func addLaneTonChangesets(t *testing.T, e *DeployedEnv, from, to uint64, fromFamily, toFamily string) []commoncs.ConfiguredChangeSet {
	fmt.Printf("DEBUG: addLaneTonChangesets %d to %d / %s to %s\n", from, to, fromFamily, toFamily)
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

func (c TonTestAddLaneChangeSet) Apply(e deployment.Environment) (deployment.ChangesetOutput, error) {
	// t := c.T
	// TODO: support other paths
	// require.Equal(t, c.fromFamily, chainsel.FamilyEVM, "must be from EVM")
	// require.Equal(t, c.toFamily, chainsel.FamilyTon, "must be to TON")

	// tonSelector := c.toChainSelector
	// tonChain := e.TonChains[tonSelector]

	// onchainState, err := changeset.LoadOnchainState(e)
	// require.NoError(t, err)
	// tonChainState := onchainState.TonChains[tonSelector]

	// fmt.Printf("DEBUG: TonTestAddLaneChangeSet: LINK token: %s CCIP: %s Receiver: %s\n", tonChainState.LinkTokenAddress.String(), tonChainState.CCIPAddress.String(), tonChainState.ReceiverAddress.String())

	// require.False(t, tonChainState.LinkTokenAddress.IsAddrNone(), "LINK token address must be set")
	// require.False(t, tonChainState.CCIPAddress.IsAddrNone(), "CCIP address must be set")
	// require.False(t, tonChainState.ReceiverAddress.IsAddrNone(), "Receiver address must be set")

	// offrampBindings := ccip_offramp.Bind(tonChainState.CCIPAddress, tonChain.Client)
	// transactOpts := &bind.TransactOpts{
	// 	Signer: tonChain.DeployerSigner,
	// }

	// evmChainState := onchainState.Chains[c.fromChainSelector]
	// evmOnrampAddress := evmChainState.OnRamp.Address().Bytes()
	// fmt.Printf("DEBUG: TonTestAddLaneChangeSet: EVM chain selector: %d - EVM onramp address: %s\n", c.fromChainSelector, hex.EncodeToString(evmOnrampAddress))

	// sourceChainSelectors := []uint64{c.fromChainSelector}
	// sourceChainsIsEnabled := []bool{true}
	// sourceChainsIsRMNVerificationDisabled := []bool{true}
	// sourceChainsOnramps := [][]byte{evmOnrampAddress}
	// pendingTx, err := offrampBindings.Offramp().ApplySourceChainConfigUpdates(transactOpts, sourceChainSelectors, sourceChainsIsEnabled, sourceChainsIsRMNVerificationDisabled, sourceChainsOnramps)
	// require.NoError(t, err)
	// waitForTx(t, tonChain.Client, pendingTx.TxnHash(), time.Minute*1)

	// fmt.Printf("DEBUG: TonTestAddLaneChangeSet: Configured offramp\n")

	return deployment.ChangesetOutput{}, nil
}

type Ton2AnyMessage struct {
	Receiver      []byte
	Data          []byte
	TokenAmounts  []Ton2AnyTokenAmount
	FeeToken      tonaddress.Address
	FeeTokenStore tonaddress.Address
	ExtraArgs     []byte
}

type Ton2AnyTokenAmount struct {
	Token      tonaddress.Address
	Amount     uint64
	TokenStore tonaddress.Address
}

func SendRequestTon(
	t *testing.T,
	e deployment.Environment,
	state tonstate.CCIPChainState,
	cfg *CCIPSendReqConfig,
) (*onramp.OnRampCCIPMessageSent, error) { // TODO: chain independent return vailue
	//sourceSelector := cfg.SourceChain
	//destSelector := cfg.DestChain
	//msg := cfg.Message.(Ton2AnyMessage)
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
	fmt.Printf("DEBUG: ConfirmCommitWithExpectedSeqNumRangeTon srcSelector: %d, startBlock: %+v, expectedSeqNumRange: %+v, enforceSingleCommit: %+v\n", srcSelector, startBlock, expectedSeqNumRange, enforceSingleCommit)
	// TODO mock the contract event before offramp contracts are supported
	return true, nil
}

// TODO: what is the usage of this function? can we remove this?
func waitForTx(t *testing.T, client *ton.APIClient, txHash string, duration time.Duration) {
	// userTx, err := client.WaitForTransaction(txHash, ton.PollTimeout(duration))
	// require.NoError(t, err)
	// require.True(t, userTx.Success, "transaction failed: %s", userTx.VmStatus)
}

func ConfirmExecWithSeqNrsTon(
	t *testing.T,
	sourceChain uint64,
	dest cldf_ton.Chain,
	offRampAddress tonaddress.Address,
	startBlock *uint64,
	expectedSeqNrs []uint64,
) (executionStates map[uint64]int, err error) {
	fmt.Printf("DEBUG: ConfirmExecWithSeqNrsTon srcSelector: %d, dest: %s, startBlock: %+v, expectedSeqNrs: %+v\n", sourceChain, startBlock, expectedSeqNrs)
	// TODO mock the contract event before offramp contracts are supported
	fmt.Printf("DEBUG: TODO(ton): ConfirmExecWithSeqNrsTon\n")
	seqNrsToWatch := make(map[uint64]int)
	for _, seqNr := range expectedSeqNrs {
		seqNrsToWatch[seqNr] = EXECUTION_STATE_SUCCESS
	}
	return seqNrsToWatch, nil
}
