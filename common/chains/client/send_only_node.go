package client

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	nodetypes "github.com/smartcontractkit/chainlink/v2/common/chains/client/types"
	"github.com/smartcontractkit/chainlink/v2/common/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// SendOnlyNode represents one node used as a sendonly
type SendOnlyNode[
	CHAIN_ID types.ID,
	RPC_CLIENT nodetypes.SendOnlyClientAPI[CHAIN_ID],
] interface {
	// Start may attempt to connect to the node, but should only return error for misconfiguration - never for temporary errors.
	Start(context.Context) error
	Close() error

	ConfiguredChainID() CHAIN_ID
	RPCClient() RPC_CLIENT

	String() string
	// State returns NodeState
	State() NodeState
	// Name is a unique identifier for this node.
	Name() string
}

// It only supports sending transactions
// It must use an http(s) url
type sendOnlyNode[
	CHAIN_ID types.ID,
	RPC_CLIENT nodetypes.SendOnlyClientAPI[CHAIN_ID],
] struct {
	utils.StartStopOnce

	stateMu sync.RWMutex // protects state* fields
	state   NodeState

	rpcClient RPC_CLIENT
	uri       url.URL
	log       logger.Logger
	name      string
	chainID   CHAIN_ID
	chStop    utils.StopChan
	wg        sync.WaitGroup
}

// NewSendOnlyNode returns a new sendonly node
func NewSendOnlyNode[
	CHAIN_ID types.ID,
	RPC_CLIENT nodetypes.SendOnlyClientAPI[CHAIN_ID],
](
	lggr logger.Logger,
	httpuri url.URL,
	name string,
	chainID CHAIN_ID,
	rpcClient RPC_CLIENT,
) SendOnlyNode[CHAIN_ID, RPC_CLIENT] {
	s := new(sendOnlyNode[CHAIN_ID, RPC_CLIENT])
	s.name = name
	s.log = lggr.Named("SendOnlyNode").Named(name).With(
		"nodeTier", "sendonly",
	)
	s.rpcClient = rpcClient
	s.uri = httpuri
	s.chainID = chainID
	s.chStop = make(chan struct{})
	return s
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) Start(ctx context.Context) error {
	return s.StartOnce(s.name, func() error {
		s.start(ctx)
		return nil
	})
}

// Start setups up and verifies the sendonly node
// Should only be called once in a node's lifecycle
func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) start(startCtx context.Context) {
	if s.State() != NodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", s.state))
	}

	err := s.rpcClient.DialHTTP()
	if err != nil {
		promPoolRPCNodeTransitionsToUnusable.WithLabelValues(s.chainID.String(), s.name).Inc()
		s.log.Errorw("Dial failed: SendOnly Node is unusable", "err", err)
		s.setState(NodeStateUnusable)
		return
	}
	s.setState(NodeStateDialed)

	if s.chainID.String() == "0" {
		// Skip verification if chainID is zero
		s.log.Warn("sendonly rpc ChainID verification skipped")
	} else {
		chainID, err := s.rpcClient.ChainID(startCtx)
		if err != nil || chainID.String() != s.chainID.String() {
			promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(s.chainID.String(), s.name).Inc()
			if err != nil {
				promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(s.chainID.String(), s.name).Inc()
				s.log.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
				s.setState(NodeStateUnreachable)
			} else {
				promPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(s.chainID.String(), s.name).Inc()
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

	promPoolRPCNodeTransitionsToAlive.WithLabelValues(s.chainID.String(), s.name).Inc()
	s.setState(NodeStateAlive)
	s.log.Infow("Sendonly RPC Node is online", "nodeState", s.state)
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) Close() error {
	return s.StopOnce(s.name, func() error {
		s.rpcClient.Close()
		s.wg.Wait()
		s.setState(NodeStateClosed)
		return nil
	})
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) ConfiguredChainID() CHAIN_ID {
	return s.chainID
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) RPCClient() RPC_CLIENT {
	return s.rpcClient
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) String() string {
	return fmt.Sprintf("(secondary)%s:%s", s.name, s.uri.Redacted())
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) setState(state NodeState) (changed bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state == state {
		return false
	}
	s.state = state
	return true
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) State() NodeState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.state
}

func (s *sendOnlyNode[CHAIN_ID, RPC_CLIENT]) Name() string {
	return s.name
}
