package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name SendOnlyNode --output ../mocks/ --case=underscore

// SendOnlyNode represents one ethereum node used as a sendonly
type SendOnlyNode interface {
	// Start may attempt to connect to the node, but should only return error for misconfiguration - never for temporary errors.
	Start(context.Context) error
	Close()

	ChainID() (chainID *big.Int)

	SendTransaction(ctx context.Context, tx *types.Transaction) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error

	String() string
}

//go:generate mockery --name TxSender --output ./mocks/ --case=underscore

type TxSender interface {
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	ChainID(context.Context) (*big.Int, error)
}

//go:generate mockery --name BatchSender --output ./mocks/ --case=underscore

type BatchSender interface {
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

// It only supports sending transactions
// It must a http(s) url
type sendOnlyNode struct {
	uri         url.URL
	batchSender BatchSender
	sender      TxSender
	log         logger.Logger
	dialed      bool
	name        string
	chainID     *big.Int
	chStop      chan struct{}
}

// NewSendOnlyNode returns a new sendonly node
func NewSendOnlyNode(lggr logger.Logger, httpuri url.URL, name string, chainID *big.Int) SendOnlyNode {
	s := new(sendOnlyNode)
	s.name = name
	s.log = lggr.Named("SendOnlyNode").Named(name).With(
		"nodeTier", "sendonly",
	)
	s.uri = httpuri
	s.chainID = chainID
	s.chStop = make(chan struct{})
	return s
}

// Start setups up and verifies the sendonly node
// Should only be called once in a node's lifecycle
// TODO: Failures to dial should put it into a retry loop
// https://app.shortcut.com/chainlinklabs/story/28182/eth-node-failover-consider-introducing-a-state-for-sendonly-nodes
func (s *sendOnlyNode) Start(startCtx context.Context) error {
	s.log.Debugw("evmclient.Client#Dial(...)")
	if s.dialed {
		panic("evmclient.Client.Dial(...) should only be called once during the node's lifetime.")
	}

	rpc, err := rpc.DialHTTP(s.uri.String())
	if err != nil {
		return errors.Wrapf(err, "failed to dial secondary client: %v", s.uri.Redacted())
	}
	s.dialed = true
	geth := ethclient.NewClient(rpc)
	s.SetEthClient(rpc, geth)

	if id, err := s.getChainID(startCtx); err != nil {
		s.log.Warn("sendonly rpc ChainID verification skipped", "err", err)
	} else if id.Cmp(s.chainID) != 0 {
		return errors.Errorf(
			"sendonly rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
			id.String(),
			s.chainID.String(),
			s.name,
		)
	}

	return nil
}

func (s *sendOnlyNode) SetEthClient(newBatchSender BatchSender, newSender TxSender) {
	if s.sender != nil {
		log.Panicf("sendOnlyNode.SetEthClient should only be called once!")
		return
	}
	s.batchSender = newBatchSender
	s.sender = newSender
}

func (s *sendOnlyNode) Close() {
	close(s.chStop)
}

func (s *sendOnlyNode) logTiming(lggr logger.Logger, duration time.Duration, err error, callName string) {
	promEVMPoolRPCCallTiming.
		WithLabelValues(
			s.chainID.String(),             // chain id
			s.name,                         // node name
			s.uri.Host,                     // rpc domain
			"true",                         // is send only
			strconv.FormatBool(err == nil), // is successful
			callName,                       // rpc call name
		).
		Observe(float64(duration))
	lggr.Debugw(fmt.Sprintf("SendOnly RPC call: evmclient.#%s", callName),
		"duration", duration,
		"rpcDomain", s.uri.Host,
		"name", s.name,
		"chainID", s.chainID,
		"sendOnly", true,
		"err", err,
	)
}

func (s *sendOnlyNode) SendTransaction(parentCtx context.Context, tx *types.Transaction) (err error) {
	defer func(start time.Time) {
		s.logTiming(s.log, time.Since(start), err, "SendTransaction")
	}(time.Now())

	ctx, cancel := s.makeQueryCtx(parentCtx)
	defer cancel()
	return s.wrap(s.sender.SendTransaction(ctx, tx))
}

func (s *sendOnlyNode) BatchCallContext(parentCtx context.Context, b []rpc.BatchElem) (err error) {
	defer func(start time.Time) {
		s.logTiming(s.log.With("nBatchElems", len(b)), time.Since(start), err, "BatchCallContext")
	}(time.Now())

	ctx, cancel := s.makeQueryCtx(parentCtx)
	defer cancel()
	return s.wrap(s.batchSender.BatchCallContext(ctx, b))
}

func (s *sendOnlyNode) ChainID() (chainID *big.Int) {
	return s.chainID
}

func (s *sendOnlyNode) wrap(err error) error {
	return wrap(err, fmt.Sprintf("sendonly http (%s)", s.uri.Redacted()))
}

func (s *sendOnlyNode) String() string {
	return fmt.Sprintf("(secondary)%s:%s", s.name, s.uri.Redacted())
}

// getChainID returns the chainID and converts zero/empty values to errors.
func (s *sendOnlyNode) getChainID(parentCtx context.Context) (*big.Int, error) {
	ctx, cancel := s.makeQueryCtx(parentCtx)
	defer cancel()

	chainID, err := s.sender.ChainID(ctx)
	if err != nil {
		return nil, err
	} else if chainID.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("zero/empty value")
	}
	return chainID, nil
}

// makeQueryCtx returns a context that cancels if:
// 1. Passed in ctx cancels
// 2. chStop is closed
// 3. Default timeout is reached (queryTimeout)
func (s *sendOnlyNode) makeQueryCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	var chCancel, timeoutCancel context.CancelFunc
	ctx, chCancel = utils.WithCloseChan(ctx, s.chStop)
	ctx, timeoutCancel = context.WithTimeout(ctx, queryTimeout)
	cancel := func() {
		chCancel()
		timeoutCancel()
	}
	return ctx, cancel
}
