package ccip

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	selectors "github.com/smartcontractkit/chain-selectors"
	"go.uber.org/atomic"

	solccip "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/ccip"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/message_hasher"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
)

type SeqNumRange struct {
	Start *atomic.Uint64
	End   *atomic.Uint64
}

type DestinationGun struct {
	l                logger.Logger
	env              cldf.Environment
	state            *stateview.CCIPOnChainState
	roundNum         *atomic.Int32
	chainSelector    uint64
	receiver         []byte
	testConfig       *ccip.LoadConfig
	evmSourceKeys    map[uint64]*bind.TransactOpts
	solanaSourceKeys map[uint64]*solana.PrivateKey
	chainOffset      int
	metricPipe       chan messageData
}

func NewDestinationGun(
	l logger.Logger,
	chainSelector uint64,
	env cldf.Environment,
	state *stateview.CCIPOnChainState,
	receiver []byte,
	overrides *ccip.LoadConfig,
	evmSourceKeys map[uint64]*bind.TransactOpts,
	solanaSourceKeys map[uint64]*solana.PrivateKey,
	chainOffset int,
	metricPipe chan messageData,
) (*DestinationGun, error) {
	dg := DestinationGun{
		l:                l,
		env:              env,
		state:            state,
		roundNum:         &atomic.Int32{},
		chainSelector:    chainSelector,
		receiver:         receiver,
		testConfig:       overrides,
		evmSourceKeys:    evmSourceKeys,
		solanaSourceKeys: solanaSourceKeys,
		chainOffset:      chainOffset,
		metricPipe:       metricPipe,
	}

	return &dg, nil
}

func (m *DestinationGun) Call(_ *wasp.Generator) *wasp.Response {
	m.roundNum.Add(1)
	src, err := m.mustSourceChain()
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: "", Failed: true}
	}
	waspGroup := fmt.Sprintf("%d->%d", src, m.chainSelector)

	selectorFamily, err := selectors.GetSelectorFamily(src)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	switch selectorFamily {
	case selectors.FamilyEVM:
		err = m.sendEVMSourceMessage(src)
	case selectors.FamilySolana:
		err = m.sendSOLSourceMessage(src)
	}
	if err != nil {
		m.l.Errorw("Failed to transmit message",
			"gun", waspGroup,
			"sourceChainFamily", selectorFamily,
			err, cldf.MaybeDataErr(err))
		if m.metricPipe != nil {
			// in the event of an error, still push a metric
			// sequence numbers start at 1 so using 0 as a sentinel value
			data := messageData{
				eventType: transmitted,
				srcDstSeqNum: srcDstSeqNum{
					src:    src,
					dst:    m.chainSelector,
					seqNum: 0,
				},
				timestamp: uint64(time.Now().Unix()), //nolint:gosec // G115
			}
			m.metricPipe <- data
		}

		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	return &wasp.Response{Failed: false, Group: waspGroup}
}

// mustSourceChain will return a chain selector to send a message from
func (m *DestinationGun) mustSourceChain() (uint64, error) {
	otherCS := m.env.BlockChains.ListChainSelectors(cldf_chain.WithChainSelectorsExclusion([]uint64{m.chainSelector}))

	if len(otherCS) == 0 {
		return 0, errors.New("no other chains to send from")
	}
	index := (int(m.roundNum.Load()) + m.chainOffset) % len(otherCS)
	return otherCS[index], nil
}

func (m *DestinationGun) sendEVMSourceMessage(src uint64) error {
	acc := m.evmSourceKeys[src]
	r := m.state.Chains[src].Router

	msg, gasLimit, err := m.GetEVMMessage(src)
	if err != nil {
		return err
	}
	// Set the gas limit for this tx
	if gasLimit != 0 {
		//nolint:gosec // it's okay here
		acc.GasLimit = uint64(gasLimit)
	}

	fee, err := r.GetFee(
		&bind.CallOpts{Context: context.Background()}, m.chainSelector, msg)
	if err != nil {
		m.l.Errorw("could not get fee ",
			"dstChainSelector", m.chainSelector,
			"fee", fee,
			"err", cldf.MaybeDataErr(err))
		return err
	}
	if msg.FeeToken == common.HexToAddress("0x0") {
		acc.Value = fee
	}
	msgWithoutData := msg
	msgWithoutData.Data = nil
	m.l.Debugw("sending message ",
		"srcChain", src,
		"dstChain", m.chainSelector,
		"fee", fee,
		"msg size", len(msg.Data),
		"msgWithoutData", msgWithoutData)
	tx, err := r.CcipSend(
		acc,
		m.chainSelector,
		msg)
	if err != nil {
		m.l.Errorw("execution reverted from ",
			"sourceChain", src,
			"destchain", m.chainSelector,
			"err", cldf.MaybeDataErr(err))
		return err
	}

	_, err = m.env.BlockChains.EVMChains()[src].Confirm(tx)
	if err != nil {
		m.l.Errorw("could not confirm tx on source", "tx", tx, "err", cldf.MaybeDataErr(err))
		return err
	}

	return nil
}

// GetEVMMessage will return the message to be sent while considering expected load of different messages
// returns the message, gas limit
func (m *DestinationGun) GetEVMMessage(src uint64) (router.ClientEVM2AnyMessage, int64, error) {
	dstSelFamily, err := selectors.GetSelectorFamily(m.chainSelector)
	if err != nil {
		return router.ClientEVM2AnyMessage{}, 0, fmt.Errorf("destination chain family for %d is not supported ", m.chainSelector)
	}
	rcv, extraArgs := []byte{}, []byte{}
	svmExtraArgs := message_hasher.ClientSVMExtraArgsV1{}
	var tokenReceiver solana.PublicKey

	// Select a message type based on ratio
	randomValue := mathrand.Intn(100)
	accumulatedRatio := 0
	var selectedMsgDetails *ccip.MsgDetails

	for _, msg := range *m.testConfig.MessageDetails {
		accumulatedRatio += *msg.Ratio
		if randomValue < accumulatedRatio {
			selectedMsgDetails = &msg
			break
		}
	}

	if selectedMsgDetails == nil {
		return router.ClientEVM2AnyMessage{}, 0, errors.New("failed to select message type")
	}

	m.l.Infow("Selected message type", "msgType", *selectedMsgDetails.MsgType)

	switch dstSelFamily {
	case selectors.FamilyEVM:
		rcv, err = utils.ABIEncode(`[{"type":"address"}]`, common.BytesToAddress(m.receiver))
		if err != nil {
			m.l.Error("Error encoding receiver address")
			return router.ClientEVM2AnyMessage{}, 0, err
		}
		extraArgs, err = GetEVMExtraArgsV2(big.NewInt(0), *m.testConfig.OOOExecution)
		if err != nil {
			m.l.Error("Error encoding extra args for evm dest")
			return router.ClientEVM2AnyMessage{}, 0, err
		}
	case selectors.FamilySolana:
		receiverTargetAccountPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, solana.PublicKeyFromBytes(m.receiver))
		receiverExternalExecutionConfigPDA, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, solana.PublicKeyFromBytes(m.receiver))
		rcv = common.LeftPadBytes(m.receiver, 32)

		accounts := [][32]byte{
			receiverExternalExecutionConfigPDA,
			receiverTargetAccountPDA,
			solana.SystemProgramID,
		}

		svmExtraArgs = message_hasher.ClientSVMExtraArgsV1{
			AccountIsWritableBitmap:  solccip.GenerateBitMapForIndexes([]int{0, 1}),
			Accounts:                 accounts,
			AllowOutOfOrderExecution: *m.testConfig.OOOExecution,
			ComputeUnits:             150000,
		}
	}
	message := router.ClientEVM2AnyMessage{
		Receiver:  rcv,
		FeeToken:  common.HexToAddress("0x0"),
		ExtraArgs: extraArgs,
	}

	// Set data length if it's a data transfer
	if selectedMsgDetails.IsDataTransfer() {
		dataLength := *selectedMsgDetails.DataLengthBytes
		switch dstSelFamily {
		case selectors.FamilyEVM:
			dataLength = *selectedMsgDetails.DataLengthBytes
		case selectors.FamilySolana:
			dataLength = *m.testConfig.SolanaDataSize
		}
		data := make([]byte, dataLength)
		_, err2 := rand.Read(data)
		if err2 != nil {
			return router.ClientEVM2AnyMessage{}, 0, err2
		}
		message.Data = data
	}

	// When it's not a programmable token transfer the receiver can be an EOA, we use a random address to denote that
	if selectedMsgDetails.IsTokenOnlyTransfer() {
		if dstSelFamily == selectors.FamilyEVM {
			receiver, err := utils.ABIEncode(`[{"type":"address"}]`, common.HexToAddress(utils.RandomAddress().Hex()))
			if err != nil {
				m.l.Error("Error encoding receiver address")
				return router.ClientEVM2AnyMessage{}, 0, err
			}
			message.Receiver = receiver
		}
	}

	// Set token amounts if it's a token transfer
	if selectedMsgDetails.IsTokenTransfer() {
		message.TokenAmounts = []router.ClientEVMTokenAmount{
			{
				Token:  m.state.Chains[src].LinkToken.Address(),
				Amount: big.NewInt(1),
			},
		}
		if dstSelFamily == selectors.FamilySolana {
			tokenReceiver, _, err = soltokens.FindAssociatedTokenAddress(
				solana.Token2022ProgramID,
				m.state.SolChains[m.chainSelector].LinkToken,
				m.state.SolChains[m.chainSelector].Receiver)
			if err != nil {
				m.l.Errorw("Error getting token receiver address")
				return router.ClientEVM2AnyMessage{}, 0, err
			}
			svmExtraArgs.TokenReceiver = tokenReceiver
		}
	}

	gasLimit := int64(0)
	if selectedMsgDetails.DestGasLimit != nil {
		gasLimit = *selectedMsgDetails.DestGasLimit
	}

	if dstSelFamily == selectors.FamilySolana {
		extraArgs, err = ccipevm.SerializeClientSVMExtraArgsV1(svmExtraArgs)
		if err != nil {
			m.l.Errorw("Error encoding extra args for sol dest")
			return router.ClientEVM2AnyMessage{}, 0, err
		}
		message.ExtraArgs = extraArgs
	}

	return message, gasLimit, nil
}

func GetEVMExtraArgsV2(gasLimit *big.Int, allowOutOfOrder bool) ([]byte, error) {
	EVMV2Tag := hexutil.MustDecode("0x181dcf10")

	encodedArgs, err := utils.ABIEncode(`[{"type":"uint256"},{"type":"bool"}]`, gasLimit, allowOutOfOrder)
	if err != nil {
		return nil, err
	}

	return append(EVMV2Tag, encodedArgs...), nil
}

func (m *DestinationGun) sendSOLSourceMessage(src uint64) error {
	msg, err := m.getSolanaMessage(src)
	if err != nil {
		return err
	}

	sendRequestCfg := testhelpers.CCIPSendReqConfig{
		SourceChain:  src,
		DestChain:    m.chainSelector,
		IsTestRouter: false,
		Message:      msg,
		MaxRetries:   1,
	}
	_, err = testhelpers.SendRequestSol(m.env, *m.state, &sendRequestCfg)
	if err != nil {
		m.l.Errorw("execution reverted from ",
			"sourceChain", src,
			"destchain", m.chainSelector,
			"err", cldf.MaybeDataErr(err))
	}
	return err
}

func (m *DestinationGun) getSolanaMessage(src uint64) (ccip_router.SVM2AnyMessage, error) {
	// Select a message type based on ratio
	randomValue := mathrand.Intn(100)
	accumulatedRatio := 0
	var selectedMsgDetails *ccip.MsgDetails

	for _, msg := range *m.testConfig.MessageDetails {
		accumulatedRatio += *msg.Ratio
		if randomValue < accumulatedRatio {
			selectedMsgDetails = &msg
			break
		}
	}

	if selectedMsgDetails == nil {
		return ccip_router.SVM2AnyMessage{}, errors.New("failed to select message type")
	}

	m.l.Infow("Selected message type", "msgType", *selectedMsgDetails.MsgType)
	message := ccip_router.SVM2AnyMessage{
		Receiver:  common.LeftPadBytes(m.receiver, 32),
		ExtraArgs: []byte{},
	}
	switch {
	case selectedMsgDetails.IsDataTransfer():
		data := make([]byte, *m.testConfig.SolanaDataSize)
		_, err := rand.Read(data)
		if err != nil {
			return ccip_router.SVM2AnyMessage{}, err
		}
		message.Data = data
	case selectedMsgDetails.IsTokenTransfer():
		message.TokenAmounts = []ccip_router.SVMTokenAmount{
			{
				Token:  m.state.SolChains[src].LinkToken,
				Amount: 1,
			},
		}
	}

	return message, nil
}
