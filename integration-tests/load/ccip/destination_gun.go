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
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"go.uber.org/atomic"
	"time"
)

type ChainSelectorPair struct {
	src uint64
	dst uint64
}

type DestinationGun struct {
	l             logger.Logger
	env           deployment.Environment
	seqNums       map[ChainSelectorPair]*atomic.Uint64
	roundNum      *atomic.Int32
	chainSelector uint64
	receiver      common.Address
	loki          *wasp.LokiClient
}

func NewDestinationGun(l logger.Logger, chainSelector uint64, env deployment.Environment, receiver common.Address, loki *wasp.LokiClient) *DestinationGun {
	seqNums := make(map[ChainSelectorPair]*atomic.Uint64)
	for _, cs := range env.AllChainSelectorsExcluding([]uint64{chainSelector}) {

		seqNums[ChainSelectorPair{
			src: cs,
			dst: chainSelector,
		}] = atomic.NewUint64(1)
	}
	return &DestinationGun{
		l:             l,
		env:           env,
		seqNums:       seqNums,
		roundNum:      &atomic.Int32{},
		chainSelector: chainSelector,
		receiver:      receiver,
		loki:          loki,
	}
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

	csPair := ChainSelectorPair{
		src: src,
		dst: m.chainSelector,
	}
	m.seqNums[csPair].Add(1)
	m.l.Infow("Starting transmit with ", "RoundNum", requestedRound, "Destination ChainSelector", m.chainSelector, "Source ChainSelector", src, "SequenceNumber", m.seqNums[csPair].Load())

	r := state.Chains[src].Router

	msg, err := m.GetMessage()
	if err != nil {
		return &wasp.Response{Error: err.Error(), Group: waspGroup, Failed: true}
	}

	fee, err := r.GetFee(
		&bind.CallOpts{Context: context.Background()}, m.chainSelector, msg)
	if err != nil {
		m.l.Errorw("could not get fee ", "dstChainSelector", m.chainSelector, "msg", msg, "fee", fee)
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

	lokiLabels, err := setLokiLabels(src, m.chainSelector)
	if err != nil {
		m.l.Errorw("Failed setting loki labels", "error", err)
	}
	SendMetricsToLoki(m.l, m.loki, lokiLabels, &LokiMetric{
		EventType:      transmitted,
		Timestamp:      time.Now(),
		SequenceNumber: m.seqNums[csPair].Load(),
	})

	return &wasp.Response{Failed: false, Group: waspGroup}
}

// MustSourceChain will return a chain selector to send a message from
func (m *DestinationGun) MustSourceChain() (uint64, error) {

	// TODO: make this smarter by checking if this chain has sent a message recently, if so, switch to the next chain
	otherCS := m.env.AllChainSelectorsExcluding([]uint64{m.chainSelector})
	if len(otherCS) == 0 {
		return 0, fmt.Errorf("no other chains to send from")
	}
	index := m.roundNum.Load() % int32(len(otherCS))
	return otherCS[index], nil
}

// GetMessage will return the message to be sent while considering expected load of different messages
// TODO: implement randomness and different types of messages
func (m *DestinationGun) GetMessage() (router.ClientEVM2AnyMessage, error) {
	rcv, err := utils.ABIEncode(`[{"type":"address"}]`, m.receiver)
	if err != nil {
		m.l.Error("Error encoding receiver address")
		return router.ClientEVM2AnyMessage{}, err
	}

	return router.ClientEVM2AnyMessage{
		Receiver:     rcv,
		Data:         common.Hex2Bytes("hello world"),
		TokenAmounts: nil,
		FeeToken:     common.HexToAddress("0x0"),
		ExtraArgs:    nil,
	}, nil
}
