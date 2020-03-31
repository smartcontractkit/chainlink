package service

import (
	"database/sql"

	"chainlink/ingester/client"
	"chainlink/ingester/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"go.uber.org/multierr"

	"fmt"
	"sync"
	"time"
)

// SubmissionsTracker is an interface for subscribing to new Chainlink aggregator feeds
type SubmissionsTracker interface {
	Start() error
	Stop()
}

type submissionsTracker struct {
	db           *sql.DB
	eth          client.ETH
	mutex        sync.Mutex
	subs         map[string]client.Subscription
	feedsTracker FeedsTracker
}

// NewSubmissionsTracker
func NewSubmissionsTracker(db *sql.DB, eth client.ETH, feedsTracker FeedsTracker) SubmissionsTracker {
	ht := &submissionsTracker{
		db:           db,
		eth:          eth,
		feedsTracker: feedsTracker,
		mutex:        sync.Mutex{},
		subs:         map[string]client.Subscription{},
	}
	return ht
}

// Start subscribes to new submissions when a new feed is added
func (st *submissionsTracker) Start() error {
	newFeeds := st.feedsTracker.Subscribe()

	go func() {
		for {
			select {
			case aggregator, open := <-newFeeds:
				if !open {
					return
				}
				go st.listenForSubmissions(aggregator)
			}
		}
	}()

	return nil
}

// Stop unsubscribes from all submissions for the tracked feeds
func (st *submissionsTracker) Stop() {
	st.mutex.Lock()
	defer st.mutex.Unlock()

	for _, sub := range st.subs {
		sub.Unsubscribe()
	}
}

// SubscribeToSubmissionReceived takes the aggregator and starts listening new submissions
func (st *submissionsTracker) SubscribeToSubmissionReceived(agg client.Aggregator) error {
	logChan := make(chan types.Log)
	sub, err := agg.SubscribeToSubmissionReceived(logChan)
	if err != nil {
		return err
	}

	st.addSubscription(agg, sub)
	defer st.deleteSubscription(agg, sub)
	logger.Info("Listening for new submission received logs")

	select {
	case srLog := <-logChan:
		defer sub.Unsubscribe()
		logger.Info("New submission log received")
		if sr, err := agg.UnmarshalSubmissionReceivedEvent(srLog); err != nil {
			return fmt.Errorf("failed to unmarshal submission received log: %+v", err)
		} else {
			logger.Infow(
				"Submission received",
				"RoundID",
				sr.RoundID,
			)

			address := make([]byte, 20)
			copy(address, srLog.Address[:])

			topics := make([]byte, len(srLog.Topics)*len(common.Hash{}))
			for index, topic := range srLog.Topics {
				copy(topics[index*len(common.Hash{}):], topic.Bytes())
			}
			_, err := st.db.Exec(
				`INSERT INTO "ethereum_log" ("address", "topics", "data", "blockNumber", "txHash", "txIndex", "blockHash", "index", "removed", "type")
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
				address,
				topics,
				srLog.Data,
				srLog.BlockNumber,
				srLog.TxHash.Bytes(),
				srLog.TxIndex,
				srLog.BlockHash.Bytes(),
				srLog.Index,
				srLog.Removed,
				"SubmissionReceived",
			)
			if err != nil {
				logger.Errorw("Insert failed", "error", err)
			}
		}
	case err, open := <-sub.Err():
		if err != nil {
			logger.Errorf("Error subscribing to submissions received: %+v", err)
		}
		if open {
			sub.Unsubscribe()
		}
		return err
	}
	return nil
}

func (st *submissionsTracker) listenForSubmissions(agg client.Aggregator) {
	logger.Infow(
		"New feed detected, listen for SubmissionReceived events",
		"Name",
		agg.Name(),
	)
	for {
		if err := st.SubscribeToSubmissionReceived(agg); err != nil {
			for _, e := range multierr.Errors(err) {
				logger.Errorf("Failed to subscribe to submission received: %+v", e)
			}
		}
		time.Sleep(1000)
	}
}

func (st *submissionsTracker) addSubscription(agg client.Aggregator, sub client.Subscription) {
	st.mutex.Lock()
	st.subs[agg.Address().String()] = sub
	st.mutex.Unlock()
}

func (st *submissionsTracker) deleteSubscription(agg client.Aggregator, sub client.Subscription) {
	st.mutex.Lock()
	delete(st.subs, agg.Address().String())
	st.mutex.Unlock()
}
