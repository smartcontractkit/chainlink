package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ ethereum.Subscription = &chainIDSubForwarder{}

// chainIDSubForwarder wraps a head subscription in order to intercept and augment each head with chainID before forwarding.
type chainIDSubForwarder struct {
	chainID *big.Int
	destCh  chan<- *evmtypes.Head

	srcCh  chan *evmtypes.Head
	srcSub ethereum.Subscription

	done  chan struct{}
	err   chan error
	unSub chan struct{}
}

func newChainIDSubForwarder(chainID *big.Int, ch chan<- *evmtypes.Head) *chainIDSubForwarder {
	return &chainIDSubForwarder{
		chainID: chainID,
		destCh:  ch,
		srcCh:   make(chan *evmtypes.Head),
		done:    make(chan struct{}),
		err:     make(chan error),
		unSub:   make(chan struct{}, 1),
	}
}

// start spawns the forwarding loop for sub.
func (c *chainIDSubForwarder) start(sub ethereum.Subscription, err error) error {
	if err != nil {
		close(c.srcCh)
		return err
	}
	c.srcSub = sub
	go c.forwardLoop()
	return nil
}

// forwardLoop receives from src, adds the chainID, and then sends to dest.
// It also handles Unsubscribing, which may interrupt either forwarding operation.
func (c *chainIDSubForwarder) forwardLoop() {
	defer close(c.done)
	for {
		select {
		case err := <-c.srcSub.Err():
			select {
			case c.err <- err:
			case <-c.unSub:
				c.srcSub.Unsubscribe()
			}
			return

		case h := <-c.srcCh:
			h.EVMChainID = utils.NewBig(c.chainID)
			select {
			case c.destCh <- h:
			case <-c.unSub:
				c.srcSub.Unsubscribe()
				return
			}

		case <-c.unSub:
			c.srcSub.Unsubscribe()
			return
		}
	}
}

func (c *chainIDSubForwarder) Unsubscribe() {
	// tell forwardLoop to unsubscribe
	select {
	case c.unSub <- struct{}{}:
	default:
		// already triggered
	}
	// wait for forwardLoop to complete
	<-c.done
}

func (c *chainIDSubForwarder) Err() <-chan error {
	return c.err
}
