package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum"

	"github.com/smartcontractkit/chainlink/core/utils"
)

var _ ethereum.Subscription = &chainIDSubForwarder{}

// chainIDSubForwarder wraps a head subscription in order to intercept and augment each head with chainID before forwarding.
type chainIDSubForwarder struct {
	chainID *big.Int
	destCh  chan<- *Head

	srcCh  chan *Head
	srcSub ethereum.Subscription

	err   chan error
	unSub chan chan struct{}
}

func newChainIDSubForwarder(chainID *big.Int, ch chan<- *Head) *chainIDSubForwarder {
	return &chainIDSubForwarder{
		chainID: chainID,
		destCh:  ch,
		srcCh:   make(chan *Head),
		err:     make(chan error),
		unSub:   make(chan chan struct{}),
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
	defer close(c.srcCh)
	for {
		select {
		case err := <-c.srcSub.Err():
			select {
			case c.err <- err:
			case done := <-c.unSub:
				c.srcSub.Unsubscribe()
				close(done)
			}
			return

		case h := <-c.srcCh:
			h.EVMChainID = utils.NewBig(c.chainID)
			select {
			case c.destCh <- h:
			case done := <-c.unSub:
				c.srcSub.Unsubscribe()
				close(done)
				return
			}

		case done := <-c.unSub:
			c.srcSub.Unsubscribe()
			close(done)
			return
		}
	}
}

func (c *chainIDSubForwarder) Unsubscribe() {
	// wait for forwardLoop to unsubscribe
	done := make(chan struct{})
	c.unSub <- done
	<-done

	close(c.err)
}

func (c *chainIDSubForwarder) Err() <-chan error {
	return c.err
}
