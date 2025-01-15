package ccip

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"go.uber.org/atomic"
	"math/big"
	"math/rand"
	"time"
)

type SeqNumRange struct {
	Start *atomic.Uint64
	End   *atomic.Uint64
}

type DestinationGun struct {
	l             logger.Logger
	env           deployment.Environment
	seqNums       map[ccipchangeset.SourceDestPair]SeqNumRange
	roundNum      *atomic.Int32
	chainSelector uint64
	receiver      common.Address
	testConfig    *ccip.LoadConfig
	loki          *wasp.LokiClient
}

func NewDestinationGun(l logger.Logger, chainSelector uint64, env deployment.Environment, receiver common.Address, overrides *ccip.LoadConfig, loki *wasp.LokiClient) (*DestinationGun, error) {
	seqNums := make(map[ccipchangeset.SourceDestPair]SeqNumRange)
	for _, cs := range env.AllChainSelectorsExcluding([]uint64{chainSelector}) {

		// query for the actual sequence number
		seqNums[ccipchangeset.SourceDestPair{
			SourceChainSelector: cs,
			DestChainSelector:   chainSelector,
		}] = SeqNumRange{
			Start: atomic.NewUint64(0),
			End:   atomic.NewUint64(0),
		}
	}
	dg := DestinationGun{
		l:             l,
		env:           env,
		seqNums:       seqNums,
		roundNum:      &atomic.Int32{},
		chainSelector: chainSelector,
		receiver:      receiver,
		testConfig:    overrides,
		loki:          loki,
	}

	err := dg.Validate()
	if err != nil {
		return nil, err
	}

	return &dg, nil
}

func (m *DestinationGun) Validate() error {
	if len(*m.testConfig.MessageTypeWeights) != 3 {
		return fmt.Errorf(
			"message type must have 3 weights corresponding to message only, token only, token with message")
	}
	sum := 0
	for _, weight := range *m.testConfig.MessageTypeWeights {
		sum += weight
	}
	if sum != 100 {
		return fmt.Errorf("message type weights must sum to 100")
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

	lokiLabels, err := setLokiLabels(src, m.chainSelector)
	if err != nil {
		m.l.Errorw("Failed setting loki labels", "error", err)
	}

	csPair := ccipchangeset.SourceDestPair{
		SourceChainSelector: src,
		DestChainSelector:   m.chainSelector,
	}
	m.l.Infow("Starting transmit with ",
		"RoundNum", requestedRound,
		"Source ChainSelector", src,
		"Destination ChainSelector", m.chainSelector)

	r := state.Chains[src].Router

	msg, err := m.GetMessage()
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	fee, err := r.GetFee(
		&bind.CallOpts{Context: context.Background()}, m.chainSelector, msg)
	if err != nil {
		m.l.Errorw("could not get fee ", "dstChainSelector", m.chainSelector, "msg", msg, "fee", fee, "err", err)
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}
	m.l.Debugw("setting fee for ", "srcChain", src, "dstChain", m.chainSelector, "fee", fee, "msg", msg)
	if msg.FeeToken == common.HexToAddress("0x0") {
		m.env.Chains[src].DeployerKey.Value = fee
		defer func() { m.env.Chains[src].DeployerKey.Value = nil }()
	}
	tx, err := r.CcipSend(
		m.env.Chains[src].DeployerKey,
		m.chainSelector,
		msg)
	if err != nil {
		m.l.Errorw("execution reverted from ", "sourceChain", src, "destchain", m.chainSelector, "err", err, "tx", tx)
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	blockNum, err := m.env.Chains[src].Confirm(tx)
	if err != nil {
		m.l.Errorw("could not confirm tx on source", "tx", tx, "err", err)
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

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

	SendMetricsToLoki(m.l, m.loki, lokiLabels, &LokiMetric{
		EventType:      transmitted,
		Timestamp:      time.Now(),
		SequenceNumber: it.Event.SequenceNumber,
	})

	// if this is the first time we are sending a message, set the start sequence number
	// if we ran into a concurrency issue, store the lowest sequence number
	if it.Event.SequenceNumber < m.seqNums[csPair].Start.Load() || m.seqNums[csPair].End.Load() == 0 {
		m.seqNums[csPair].Start.Store(it.Event.SequenceNumber)
	}

	// only store the greatest sequence number we have seen as the maximum
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
		return 0, fmt.Errorf("no other chains to send from")
	}
	index := m.roundNum.Load() % int32(len(otherCS))
	return otherCS[index], nil
}

// GetMessage will return the message to be sent while considering expected load of different messages
func (m *DestinationGun) GetMessage() (router.ClientEVM2AnyMessage, error) {
	rcv, err := utils.ABIEncode(`[{"type":"address"}]`, m.receiver)
	if err != nil {
		m.l.Error("Error encoding receiver address")
		return router.ClientEVM2AnyMessage{}, err
	}

	messages := []router.ClientEVM2AnyMessage{
		{
			Receiver:     rcv,
			Data:         common.Hex2Bytes("message"),
			TokenAmounts: nil,
			FeeToken:     common.HexToAddress("0x0"),
			ExtraArgs:    nil,
		},
		{
			Receiver: rcv,
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0"),
					Amount: big.NewInt(100),
				},
			},
			Data:      common.Hex2Bytes("hello world"),
			FeeToken:  common.HexToAddress("0x0"),
			ExtraArgs: nil,
		},
		{
			Receiver: rcv,
			Data:     common.Hex2Bytes("message with token"),
			TokenAmounts: []router.ClientEVMTokenAmount{
				{
					Token:  common.HexToAddress("0x0"),
					Amount: big.NewInt(100),
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
