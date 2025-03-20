package ccip

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_2_0/router"
)

type SeqNumRange struct {
	Start *atomic.Uint64
	End   *atomic.Uint64
}

type DestinationGun struct {
	l             logger.Logger
	env           deployment.Environment
	state         *ccipchangeset.CCIPOnChainState
	roundNum      *atomic.Int32
	chainSelector uint64
	receiver      common.Address
	testConfig    *ccip.LoadConfig
	messageKeys   map[uint64]*bind.TransactOpts
	chainOffset   int
	metricPipe    chan messageData
}

func NewDestinationGun(
	l logger.Logger,
	chainSelector uint64,
	env deployment.Environment,
	state *ccipchangeset.CCIPOnChainState,
	receiver common.Address,
	overrides *ccip.LoadConfig,
	messageKeys map[uint64]*bind.TransactOpts,
	chainOffset int,
	metricPipe chan messageData,
) (*DestinationGun, error) {
	dg := DestinationGun{
		l:             l,
		env:           env,
		state:         state,
		roundNum:      &atomic.Int32{},
		chainSelector: chainSelector,
		receiver:      receiver,
		testConfig:    overrides,
		messageKeys:   messageKeys,
		chainOffset:   chainOffset,
		metricPipe:    metricPipe,
	}

	return &dg, nil
}

func (m *DestinationGun) Call(_ *wasp.Generator) *wasp.Response {
	m.roundNum.Add(1)

	waspGroup := fmt.Sprintf("%d-%s", m.chainSelector, "messageOnly")

	state, err := ccipchangeset.LoadOnchainState(m.env)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	src, err := m.MustSourceChain()
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	acc := m.messageKeys[src]

	r := state.Chains[src].Router

	msg, gasLimit, err := m.GetMessage(src)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
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
			"msg", msg,
			"fee", fee,
			"err", deployment.MaybeDataErr(err))
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
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
			"err", deployment.MaybeDataErr(err))

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

		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	_, err = m.env.Chains[src].Confirm(tx)
	if err != nil {
		m.l.Errorw("could not confirm tx on source", "tx", tx, "err", deployment.MaybeDataErr(err))
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	return &wasp.Response{Failed: false, Group: waspGroup}
}

// MustSourceChain will return a chain selector to send a message from
func (m *DestinationGun) MustSourceChain() (uint64, error) {
	otherCS := m.env.AllChainSelectorsExcluding([]uint64{m.chainSelector})
	if len(otherCS) == 0 {
		return 0, errors.New("no other chains to send from")
	}
	index := mathrand.Intn(len(otherCS))
	return otherCS[index], nil
}

// GetMessage will return the message to be sent while considering expected load of different messages
// returns the message, gas limit
func (m *DestinationGun) GetMessage(src uint64) (router.ClientEVM2AnyMessage, int64, error) {
	rcv, err := utils.ABIEncode(`[{"type":"address"}]`, m.receiver)
	if err != nil {
		m.l.Error("Error encoding receiver address")
		return router.ClientEVM2AnyMessage{}, 0, err
	}
	extraArgs, err := GetEVMExtraArgsV2(big.NewInt(0), *m.testConfig.OOOExecution)
	if err != nil {
		m.l.Error("Error encoding extra args")
		return router.ClientEVM2AnyMessage{}, 0, err
	}

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

	message := router.ClientEVM2AnyMessage{
		Receiver:  rcv,
		FeeToken:  common.HexToAddress("0x0"),
		ExtraArgs: extraArgs,
	}

	// Set data length if it's a data transfer
	if selectedMsgDetails.IsDataTransfer() {
		dataLength := *selectedMsgDetails.DataLengthBytes
		data := make([]byte, dataLength)
		_, err2 := rand.Read(data)
		if err2 != nil {
			return router.ClientEVM2AnyMessage{}, 0, err2
		}
		message.Data = data
	}

	// When it's not a programmable token transfer the receiver can be an EOA, we use a random address to denote that
	if selectedMsgDetails.IsTokenOnlyTransfer() {
		receiver, err := utils.ABIEncode(`[{"type":"address"}]`, common.HexToAddress(utils.RandomAddress().Hex()))
		if err != nil {
			m.l.Error("Error encoding receiver address")
			return router.ClientEVM2AnyMessage{}, 0, err
		}
		message.Receiver = receiver
	}

	// Set token amounts if it's a token transfer
	if selectedMsgDetails.IsTokenTransfer() {
		message.TokenAmounts = []router.ClientEVMTokenAmount{
			{
				Token:  m.state.Chains[src].LinkToken.Address(),
				Amount: big.NewInt(1),
			},
		}
	}

	gasLimit := int64(0)
	if selectedMsgDetails.DestGasLimit != nil {
		gasLimit = *selectedMsgDetails.DestGasLimit
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
