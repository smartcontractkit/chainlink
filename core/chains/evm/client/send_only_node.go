package client

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --quiet --name SendOnlyNode --output ../mocks/ --case=underscore

// SendOnlyNode represents one ethereum node used as a sendonly
type SendOnlyNode interface {
	// Start may attempt to connect to the node, but should only return error for misconfiguration - never for temporary errors.
	Start(context.Context) error
	Close() error

	ChainID() (chainID *big.Int)

	SendTransaction(ctx context.Context, tx *types.Transaction) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error

	String() string
	// State returns NodeState
	State() NodeState
	// Name is a unique identifier for this node.
	Name() string
}

//go:generate mockery --quiet --name TxSender --output ./mocks/ --case=underscore

type TxSender interface {
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	ChainID(context.Context) (*big.Int, error)
}

//go:generate mockery --quiet --name BatchSender --output ./mocks/ --case=underscore

type BatchSender interface {
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

var _ SendOnlyNode = &sendOnlyNode{}

// It only supports sending transactions
// It must a http(s) url
type sendOnlyNode struct {
	utils.StartStopOnce

	stateMu sync.RWMutex // protects state* fields
	state   NodeState

	uri         url.URL
	batchSender BatchSender
	sender      TxSender
	log         logger.Logger
	dialed      bool
	name        string
	chainID     *big.Int
	chStop      chan struct{}
	wg          sync.WaitGroup
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

func (s *sendOnlyNode) Start(ctx context.Context) error {
	return s.StartOnce(s.name, func() error {
		s.start(ctx)
		return nil
	})
}

// Start setups up and verifies the sendonly node
// Should only be called once in a node's lifecycle
func (s *sendOnlyNode) start(startCtx context.Context) {
	if s.state != NodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", s.state))
	}

	s.log.Debugw("evmclient.Client#Dial(...)")
	if s.dialed {
		panic("evmclient.Client.Dial(...) should only be called once during the node's lifetime.")
	}

	// DialHTTP doesn't actually make any external HTTP calls
	// It can only return error if the URL is malformed. No amount of retries
	// will change this result.
	rpc, err := rpc.DialHTTP(s.uri.String())
	if err != nil {
		promEVMPoolRPCNodeTransitionsToUnusable.WithLabelValues(s.chainID.String(), s.name).Inc()
		s.log.Errorw("Dial failed: EVM SendOnly Node is unusable", "err", err)
		s.setState(NodeStateUnusable)
		return
	}
	s.dialed = true
	geth := ethclient.NewClient(rpc)
	s.SetEthClient(rpc, geth)

	if s.chainID.Cmp(big.NewInt(0)) == 0 {
		// Skip verification if chainID is zero
		s.log.Warn("sendonly rpc ChainID verification skipped")
	} else {
		verifyCtx, verifyCancel := s.makeQueryCtx(startCtx)
		defer verifyCancel()

		chainID, err := s.sender.ChainID(verifyCtx)
		if err != nil || chainID.Cmp(s.chainID) != 0 {
			promEVMPoolRPCNodeTransitionsToUnreachable.WithLabelValues(s.chainID.String(), s.name).Inc()
			if err != nil {
				promEVMPoolRPCNodeTransitionsToUnreachable.WithLabelValues(s.chainID.String(), s.name).Inc()
				s.log.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
				s.setState(NodeStateUnreachable)
			} else {
				promEVMPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(s.chainID.String(), s.name).Inc()
				s.log.Errorf(
					"sendonly rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
					chainID.String(),
					s.chainID.String(),
					s.name,
				)
				s.setState(NodeStateInvalidChainID)
			}
			// Since it has failed, spin up the verifyLoop that will keep
			// retrying until success
			s.wg.Add(1)
			go s.verifyLoop()
			return
		}
	}

	promEVMPoolRPCNodeTransitionsToAlive.WithLabelValues(s.chainID.String(), s.name).Inc()
	s.setState(NodeStateAlive)
	s.log.Infow("Sendonly RPC Node is online", "nodeState", s.state)
}

func (s *sendOnlyNode) SetEthClient(newBatchSender BatchSender, newSender TxSender) {
	if s.sender != nil {
		log.Panicf("sendOnlyNode.SetEthClient should only be called once!")
		return
	}
	s.batchSender = newBatchSender
	s.sender = newSender
}

func (s *sendOnlyNode) Close() error {
	return s.StopOnce(s.name, func() error {
		close(s.chStop)
		s.wg.Wait()
		s.setState(NodeStateClosed)
		return nil
	})
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

func (s *sendOnlyNode) setState(state NodeState) (changed bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state == state {
		return false
	}
	s.state = state
	return true
}

func (s *sendOnlyNode) State() NodeState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.state
}

func (s *sendOnlyNode) Name() string {
	return s.name
}
