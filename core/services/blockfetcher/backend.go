package blockfetcher

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// HeadSubscription - Managing fetching of data for a single head
type HeadSubscription interface {
	Head() models.Head

	// Block - will deliver the block data
	Block() <-chan types.Block

	// Receipts - will deliver receipt data
	Receipts() <-chan []types.Receipt

	// Unsubscribe - stops inflow of data for this head
	Unsubscribe()
}

// Backend - Managing fetching of heads and block data
type Backend interface {

	// Subscribe - Start receiving newest heads and data for them
	Subscribe() <-chan HeadSubscription

	HeadByNumber(n *big.Int) HeadSubscription
}
