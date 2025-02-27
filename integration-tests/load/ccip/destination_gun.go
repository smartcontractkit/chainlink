package ccip

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"time"

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

	err := dg.Validate()
	if err != nil {
		return nil, err
	}

	return &dg, nil
}

func (m *DestinationGun) Validate() error {
	if len(*m.testConfig.MessageTypeWeights) != 3 {
		return errors.New(
			"message type must have 3 weights corresponding to message only, token only, token with message")
	}
	sum := 0
	for _, weight := range *m.testConfig.MessageTypeWeights {
		sum += weight
	}
	if sum != 100 {
		return errors.New("message type weights must sum to 100")
	}
	return nil
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

	msg, err := m.GetMessage(src)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
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
		defer func() { acc.Value = nil }()
	}
	// todo: remove once CCIP-5143 is implemented
	acc.GasLimit = *m.testConfig.GasLimit
	
	m.l.Debugw("sending message ",
		"srcChain", src,
		"dstChain", m.chainSelector,
		"fee", fee,
		"msg", msg)
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
			timestamp: uint64(time.Now().Unix()),
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
	// TODO: make this smarter by checking if this chain has sent a message recently, if so, switch to the next chain
	// Currently performing a round robin
	otherCS := m.env.AllChainSelectorsExcluding([]uint64{m.chainSelector})
	if len(otherCS) == 0 {
		return 0, errors.New("no other chains to send from")
	}
	index := (int(m.roundNum.Load()) + m.chainOffset) % len(otherCS)
	return otherCS[index], nil
}

// GetMessage will return the message to be sent while considering expected load of different messages
func (m *DestinationGun) GetMessage(src uint64) (router.ClientEVM2AnyMessage, error) {
	rcv, err := utils.ABIEncode(`[{"type":"address"}]`, m.receiver)
	if err != nil {
		m.l.Error("Error encoding receiver address")
		return router.ClientEVM2AnyMessage{}, err
	}

	messages := []router.ClientEVM2AnyMessage{
		{
			Receiver:     rcv,
			Data:         common.Hex2Bytes("0xabcdefabcdef"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		},
		{
			Receiver: rcv,
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  m.state.Chains[src].LinkToken.Address(),
					Amount: big.NewInt(1),
				},
			},
			Data:      common.Hex2Bytes("0xabcdefabcdef"),
			FeeToken:  common.HexToAddress("0x0"),
			ExtraArgs: nil,
		},
		{
			Receiver: rcv,
			Data:     common.Hex2Bytes("message with token"),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  m.state.Chains[src].LinkToken.Address(),
					Amount: big.NewInt(1),
				},
			},
			FeeToken:  common.HexToAddress("0x0"),
			ExtraArgs: nil,
		},
	}
	// Select a random message
	randomValue := rand.Intn(100)
	switch {
	case randomValue < (*m.testConfig.MessageTypeWeights)[0]:
		return messages[0], nil
	case randomValue < (*m.testConfig.MessageTypeWeights)[0]+(*m.testConfig.MessageTypeWeights)[1]:
		return messages[1], nil
	default:
		return messages[2], nil
	}
}
