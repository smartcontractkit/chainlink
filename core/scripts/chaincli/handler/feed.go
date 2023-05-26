package handler

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
	feed "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/derived_price_feed_wrapper"
)

// Feed is the price feeds commands handler
type Feed struct {
	*baseHandler

	baseAddress  common.Address
	quoteAddress common.Address
	decimals     uint8
}

// NewFeed is the constructor of Feed
func NewFeed(cfg *config.Config) *Feed {
	return &Feed{
		baseHandler:  NewBaseHandler(cfg),
		baseAddress:  common.HexToAddress(cfg.FeedBaseAddr),
		quoteAddress: common.HexToAddress(cfg.FeedQuoteAddr),
		decimals:     cfg.FeedDecimals,
	}
}

// DeployDerivedPriceFeed deploys and approves the derived price feed.
func (h *Feed) DeployDerivedPriceFeed(ctx context.Context) {
	// Deploy derived price feed
	log.Println("Deploying derived price feed...")
	feedAddr, deployFeedTx, _, err := feed.DeployDerivedPriceFeed(h.buildTxOpts(ctx), h.client,
		h.baseAddress,
		h.quoteAddress,
		h.decimals,
	)
	if err != nil {
		log.Fatal("DeployDerivedPriceFeed failed: ", err)
	}
	log.Println("Waiting for derived price feed contract deployment confirmation...", deployFeedTx.Hash().Hex())
	h.waitDeployment(ctx, deployFeedTx)
	log.Println(feedAddr.Hex(), ": Derived price feed successfully deployed - ", deployFeedTx.Hash().Hex())

	// Approve derived price feed
	approveRegistryTx, err := h.linkToken.Approve(h.buildTxOpts(ctx), feedAddr, h.approveAmount)
	if err != nil {
		log.Fatal(feedAddr.Hex(), ": Approve failed - ", err)
	}
	h.waitTx(ctx, approveRegistryTx)
	log.Println(feedAddr.Hex(), ": Derived price feed successfully approved - ", approveRegistryTx.Hash().Hex())
}
