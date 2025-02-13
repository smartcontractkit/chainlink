package ccip

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink/deployment"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
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
	seqNums       map[testhelpers.SourceDestPair]SeqNumRange
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
	seqNums := make(map[testhelpers.SourceDestPair]SeqNumRange)
	for _, cs := range env.AllChainSelectorsExcluding([]uint64{chainSelector}) {
		seqNums[testhelpers.SourceDestPair{
			SourceChainSelector: cs,
			DestChainSelector:   chainSelector,
		}] = SeqNumRange{
			Start: atomic.NewUint64(math.MaxUint64),
			End:   atomic.NewUint64(0),
		}
	}
	dg := DestinationGun{
		l:             l,
		env:           env,
		state:         state,
		seqNums:       seqNums,
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
	requestedRound := m.roundNum.Load()

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

	csPair := testhelpers.SourceDestPair{
		SourceChainSelector: src,
		DestChainSelector:   m.chainSelector,
	}

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
	m.l.Debugw("sending message ",
		"srcChain", src,
		"dstChain", m.chainSelector,
		"round", requestedRound,
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
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	blockNum, err := m.env.Chains[src].Confirm(tx)
	if err != nil {
		m.l.Errorw("could not confirm tx on source", "tx", tx, "err", deployment.MaybeDataErr(err))
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	// todo: wasp should not manage confirming the message
	// instead, we should manage the sequence number atomically (at a higher level)
	it, err := state.Chains[src].OnRamp.FilterCCIPMessageSent(&bind.FilterOpts{
		Start:   blockNum,
		End:     &blockNum,
		Context: context.Background(),
	}, []uint64{m.chainSelector}, []uint64{})
	if err != nil {
		m.l.Errorw("could not find sent message event on src chain", "src", src, "dst", m.chainSelector, "err", err)
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}
	if !it.Next() {
		m.l.Errorw("Could not find event")
		return &wasp.Response{Error: "Could not iterate", Group: waspGroup, Failed: true}
	}

	m.l.Infow("Transmitted message with",
		"sourceChain", src,
		"destChain", m.chainSelector,
		"sequence number", it.Event.SequenceNumber)

	// push metric to metric manager for eventual distribution to loki
	blockNum = it.Event.Raw.BlockNumber
	header, err := m.env.Chains[src].Client.HeaderByNumber(m.env.GetContext(), new(big.Int).SetUint64(blockNum))
	if err != nil {
		return &wasp.Response{Error: "Could not get timestamp of block number", Group: waspGroup, Failed: true}
	}
	m.metricPipe <- messageData{
		eventType: transmitted,
		srcDstSeqNum: srcDstSeqNum{
			src:    src,
			dst:    m.chainSelector,
			seqNum: it.Event.SequenceNumber,
		},
		timestamp: header.Time,
		round:     int(requestedRound),
	}

	// always store the lowest seen number as the start seq num
	if it.Event.SequenceNumber < m.seqNums[csPair].Start.Load() {
		m.seqNums[csPair].Start.Store(it.Event.SequenceNumber)
	}

	// always store the greatest sequence number we have seen as the maximum
	if it.Event.SequenceNumber > m.seqNums[csPair].End.Load() {
		m.seqNums[csPair].End.Store(it.Event.SequenceNumber)
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
