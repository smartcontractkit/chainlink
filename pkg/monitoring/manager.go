package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/stretchr/testify/assert"
)

// Manager restarts the managed function with a new list of updates whenever something changed.
type Manager interface {
	Start(backgroundCtx context.Context, backgroundWg *sync.WaitGroup, managed ManagedFunc)
	HTTPHandler() http.Handler
}

type ManagedFunc func(localCtx context.Context, localWg *sync.WaitGroup, feeds []FeedConfig)

func NewManager(
	log Logger,
	rddPoller Poller,
) Manager {
	return &managerImpl{
		log,
		rddPoller,
		[]FeedConfig{},
		sync.Mutex{},
	}
}

type managerImpl struct {
	log       Logger
	rddPoller Poller

	currentFeeds   []FeedConfig
	currentFeedsMu sync.Mutex
}

func (m *managerImpl) Start(backgroundCtx context.Context, backgroundWg *sync.WaitGroup, managed ManagedFunc) {
	var localCtx context.Context
	var localCtxCancel context.CancelFunc
	var localWg *sync.WaitGroup
	for {
		select {
		case rawUpdatedFeeds := <-m.rddPoller.Updates():
			updatedFeeds, ok := rawUpdatedFeeds.([]FeedConfig)
			if !ok {
				m.log.Errorw("unexpected type for rdd updates", "type", fmt.Sprintf("%T", updatedFeeds))
				continue
			}
			shouldRestartMonitor := false
			func() {
				m.currentFeedsMu.Lock()
				defer m.currentFeedsMu.Unlock()
				shouldRestartMonitor = isDifferentFeeds(m.currentFeeds, updatedFeeds)
				if shouldRestartMonitor {
					m.currentFeeds = updatedFeeds
				}
			}()
			if !shouldRestartMonitor {
				continue
			}
			m.log.Infow("change in feeds configuration detected", "feeds", updatedFeeds)
			// Terminate previous managed function if not the first run.
			if localCtxCancel != nil && localWg != nil {
				localCtxCancel()
				localWg.Wait()
			}
			// Start new managed function
			localCtx, localCtxCancel = context.WithCancel(backgroundCtx)
			localWg = &sync.WaitGroup{}
			localWg.Add(1)
			go func() {
				defer localWg.Done()
				managed(localCtx, localWg, updatedFeeds)
			}()
		case <-backgroundCtx.Done():
			if localCtxCancel != nil {
				localCtxCancel()
			}
			if localWg != nil {
				localWg.Wait()
			}
			m.log.Infow("manager stopped")
			return
		}
	}
}

func (m *managerImpl) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var currentFeeds []FeedConfig
		func() { // take a snaphost of the current feeds
			m.currentFeedsMu.Lock()
			defer m.currentFeedsMu.Unlock()
			currentFeeds = m.currentFeeds
		}()
		writer.Header().Set("content-type", "application/json")
		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(currentFeeds); err != nil {
			m.log.Errorw("failed to write current feeds to the http handler", "error", err)
		}
	})
}

// isDifferentFeeds checks whether there is a difference between the current list of feeds and the new feeds - Manager
func isDifferentFeeds(current, updated []FeedConfig) bool {
	return !assert.ObjectsAreEqual(current, updated)
}
