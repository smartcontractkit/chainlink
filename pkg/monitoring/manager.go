package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/utils"
)

// Manager restarts the multi-feed monitor whenever the feed configuration list has changed.
// In order to not be coupled with the MultiFeedMonitor component, it simply runs a function
// every time the feed configuration has changed. This is hooked up to the MultiFeedMonitor.Run method in the Monitor.
type Manager interface {
	Run(backgroundCtx context.Context, managed ManagedFunc)
	HTTPHandler() http.Handler
}

type ManagedFunc func(localCtx context.Context, data RDDData)

func NewManager(
	log Logger,
	rddPoller Poller,
) Manager {
	return &managerImpl{
		log,
		rddPoller,
		RDDData{},
		sync.Mutex{},
	}
}

type managerImpl struct {
	log       Logger
	rddPoller Poller

	currentData   RDDData
	currentDataMu sync.Mutex
}

func (m *managerImpl) Run(backgroundCtx context.Context, managed ManagedFunc) {
	var localCtx context.Context
	var localCtxCancel context.CancelFunc
	var localSubs *utils.Subprocesses
	for {
		select {
		case rawData := <-m.rddPoller.Updates():
			updatedData, ok := rawData.(RDDData)
			if !ok {
				m.log.Errorw("unexpected type for rdd updates", "type", fmt.Sprintf("%T", updatedData))
				continue
			}
			shouldRestartMonitor := false
			func() {
				m.currentDataMu.Lock()
				defer m.currentDataMu.Unlock()
				shouldRestartMonitor = isDifferentData(m.currentData, updatedData)
				if shouldRestartMonitor {
					m.currentData = updatedData
				}
			}()
			if !shouldRestartMonitor {
				continue
			}
			m.log.Infow("change in feeds configuration detected", "feeds", fmt.Sprintf("%#v", updatedData))
			// Terminate previous managed function if not the first run.
			if localCtxCancel != nil && localSubs != nil {
				localCtxCancel()
				localSubs.Wait()
			}
			// Start new managed function
			localCtx, localCtxCancel = context.WithCancel(backgroundCtx)
			localSubs = &utils.Subprocesses{}
			localSubs.Go(func() {
				managed(localCtx, updatedData)
			})
		case <-backgroundCtx.Done():
			if localCtxCancel != nil {
				localCtxCancel()
			}
			if localSubs != nil {
				localSubs.Wait()
			}
			m.log.Infow("manager stopped")
			return
		}
	}
}

func (m *managerImpl) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var currentData RDDData
		func() { // take a snaphost of the current feeds
			m.currentDataMu.Lock()
			defer m.currentDataMu.Unlock()
			currentData = m.currentData
		}()
		writer.Header().Set("content-type", "application/json")
		encoder := json.NewEncoder(writer)
		if err := encoder.Encode(currentData); err != nil {
			m.log.Errorw("failed to write current feeds to the http handler", "error", err)
		}
	})
}

// isDifferentData checks whether there is a difference between the current list of feeds and the new feeds - Manager
func isDifferentData(current, updated RDDData) bool {
	return !assert.ObjectsAreEqual(current, updated)
}
