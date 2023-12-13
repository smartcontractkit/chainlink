package client

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/common/types"
)

//go:generate mockery --quiet --name sendOnlyClient --structname mockSendOnlyClient --filename "mock_send_only_client_test.go" --inpackage --case=underscore
type sendOnlyClient[
	CHAIN_ID types.ID,
] interface {
	Close()
	ChainID(context.Context) (CHAIN_ID, error)
	DialHTTP() error
}

// SendOnlyNode represents one node used as a sendonly
//
//go:generate mockery --quiet --name SendOnlyNode --structname mockSendOnlyNode --filename "mock_send_only_node_test.go"  --inpackage --case=underscore
type SendOnlyNode[
	CHAIN_ID types.ID,
	RPC sendOnlyClient[CHAIN_ID],
] interface {
	// Start may attempt to connect to the node, but should only return error for misconfiguration - never for temporary errors.
	Start(context.Context) error
	Close() error

	ConfiguredChainID() CHAIN_ID
	RPC() RPC

	String() string
	// State returns nodeState
	State() nodeState
	// Name is a unique identifier for this node.
	Name() string
}

// It only supports sending transactions
// It must use an http(s) url
type sendOnlyNode[
	CHAIN_ID types.ID,
	RPC sendOnlyClient[CHAIN_ID],
] struct {
	services.StateMachine

	stateMu sync.RWMutex // protects state* fields
	state   nodeState

	rpc     RPC
	uri     url.URL
	log     logger.Logger
	name    string
	chainID CHAIN_ID
	chStop  services.StopChan
	wg      sync.WaitGroup
}

// NewSendOnlyNode returns a new sendonly node
func NewSendOnlyNode[
	CHAIN_ID types.ID,
	RPC sendOnlyClient[CHAIN_ID],
](
	lggr logger.Logger,
	httpuri url.URL,
	name string,
	chainID CHAIN_ID,
	rpc RPC,
) SendOnlyNode[CHAIN_ID, RPC] {
	s := new(sendOnlyNode[CHAIN_ID, RPC])
	s.name = name
	s.log = logger.Named(logger.Named(lggr, "SendOnlyNode"), name)
	s.log = logger.With(s.log,
		"nodeTier", "sendonly",
	)
	s.rpc = rpc
	s.uri = httpuri
	s.chainID = chainID
	s.chStop = make(chan struct{})
	return s
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) Start(ctx context.Context) error {
	return s.StartOnce(s.name, func() error {
		s.start(ctx)
		return nil
	})
}

// Start setups up and verifies the sendonly node
// Should only be called once in a node's lifecycle
func (s *sendOnlyNode[CHAIN_ID, RPC]) start(startCtx context.Context) {
	if s.State() != nodeStateUndialed {
		panic(fmt.Sprintf("cannot dial node with state %v", s.state))
	}

	err := s.rpc.DialHTTP()
	if err != nil {
		promPoolRPCNodeTransitionsToUnusable.WithLabelValues(s.chainID.String(), s.name).Inc()
		s.log.Errorw("Dial failed: SendOnly Node is unusable", "err", err)
		s.setState(nodeStateUnusable)
		return
	}
	s.setState(nodeStateDialed)

	if s.chainID.String() == "0" {
		// Skip verification if chainID is zero
		s.log.Warn("sendonly rpc ChainID verification skipped")
	} else {
		chainID, err := s.rpc.ChainID(startCtx)
		if err != nil || chainID.String() != s.chainID.String() {
			promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(s.chainID.String(), s.name).Inc()
			if err != nil {
				promPoolRPCNodeTransitionsToUnreachable.WithLabelValues(s.chainID.String(), s.name).Inc()
				s.log.Errorw(fmt.Sprintf("Verify failed: %v", err), "err", err)
				s.setState(nodeStateUnreachable)
			} else {
				promPoolRPCNodeTransitionsToInvalidChainID.WithLabelValues(s.chainID.String(), s.name).Inc()
				s.log.Errorf(
					"sendonly rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
					chainID.String(),
					s.chainID.String(),
					s.name,
				)
				s.setState(nodeStateInvalidChainID)
			}
			// Since it has failed, spin up the verifyLoop that will keep
			// retrying until success
			s.wg.Add(1)
			go s.verifyLoop()
			return
		}
	}

	promPoolRPCNodeTransitionsToAlive.WithLabelValues(s.chainID.String(), s.name).Inc()
	s.setState(nodeStateAlive)
	s.log.Infow("Sendonly RPC Node is online", "nodeState", s.state)
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) Close() error {
	return s.StopOnce(s.name, func() error {
		s.rpc.Close()
		close(s.chStop)
		s.wg.Wait()
		s.setState(nodeStateClosed)
		return nil
	})
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) ConfiguredChainID() CHAIN_ID {
	return s.chainID
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) RPC() RPC {
	return s.rpc
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) String() string {
	return fmt.Sprintf("(%s)%s:%s", Secondary.String(), s.name, s.uri.Redacted())
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) setState(state nodeState) (changed bool) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	if s.state == state {
		return false
	}
	s.state = state
	return true
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) State() nodeState {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.state
}

func (s *sendOnlyNode[CHAIN_ID, RPC]) Name() string {
	return s.name
}
