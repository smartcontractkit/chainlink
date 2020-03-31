package service

import (
	"fmt"
	"time"

	"chainlink/ingester/client"
	"chainlink/ingester/logger"
	"chainlink/ingester/util"

	"github.com/ethereum/go-ethereum/common"
)

// FeedsTracker is an interface for subscribing to new Chainlink aggregator feeds
type FeedsTracker interface {
	util.TickerService
	Subscribe() <-chan client.Aggregator
}

type feedsTracker struct {
	util.Ticker

	feedsByAddress map[common.Address]client.Aggregator
	feedsChannel   chan client.Aggregator

	eth       client.ETH
	feedsUI   client.FeedsUI
	networkId int
}

// NewFeedsTracker returns an instantiated instance of a FeedsTracker implementation
func NewFeedsTracker(eth client.ETH, feedsUI client.FeedsUI, networkId int) FeedsTracker {
	ag := &feedsTracker{
		feedsByAddress: map[common.Address]client.Aggregator{},
		feedsChannel:   make(chan client.Aggregator),
		eth:            eth,
		feedsUI:        feedsUI,
		networkId:      networkId,
	}
	ag.Ticker = util.Ticker{
		Ticker: time.NewTicker(time.Minute),
		Name:   "feedsTracker",
		Impl:   ag,
		Done:   make(chan struct{}),
		Exited: make(chan struct{}),
	}
	return ag
}

// Tick implements the ServiceTicker Tick interface which fetches the list of feeds
// via an API call, checking to see if there's any new feeds that aren't yet subscribed
func (ft *feedsTracker) Tick() {
	feeds, err := ft.feedsUI.Feeds()
	if err != nil {
		logger.Errorf("Error while fetching feeds: %v", err)
	}

	for _, feed := range feeds {
		if !ft.supported(feed) {
			continue
		}

		if agg, err := client.NewAggregator(ft.eth, ft.feedsUI, feed.Name, feed.ContractAddress); err != nil {
			logger.Errorf("Error while creating new aggregator: %v", err)
		} else if err := ft.add(feed.ContractAddress, agg); err != nil {
			logger.Warnf(
				fmt.Sprintf("Ignoring aggregator contract due to latestRound error: %v", err),
				"name", feed.Name,
				"address", feed.ContractAddress.String(),
				"networkId", feed.NetworkID,
			)
		}
	}
}

// Subscribe returns the feeds channel so that any subscriber can receive new feeds
func (ft *feedsTracker) Subscribe() <-chan client.Aggregator {
	return ft.feedsChannel
}

// Currently only supports Aggregator v2
func (ft *feedsTracker) supported(feed *client.UIFeed) bool {
	return feed.NetworkID == ft.networkId && feed.ContractVersion == 2
}

func (ft *feedsTracker) add(address common.Address, agg client.Aggregator) error {
	if _, ok := ft.feedsByAddress[address]; ok {
		return nil
	}
	round, err := agg.LatestRound()
	if err != nil {
		return fmt.Errorf("invalid aggregator address: %+v", err)
	}

	logger.Infow(
		"New feed found",
		"address", agg.Address().String(),
		"name", agg.Name(),
		"round", round,
	)
	ft.feedsByAddress[address] = agg
	select {
	case ft.feedsChannel <- agg:
	default:
		logger.Warn("Start subscribers before FeedsTracker")
	}
	return nil
}
