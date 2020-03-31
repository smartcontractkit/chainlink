package util

import (
	"time"

	"chainlink/ingester/logger"
)

// TickerService is an interface that provides
// basic ticking functionality to a service
type TickerService interface {
	Start()
	Stop()
	Tick()
}

// Ticker is the implementation of the TickerService
// interface, managing the lifecycle of a ticking service
type Ticker struct {
	Ticker *time.Ticker
	Name   string
	Impl   TickerService
	Done   chan struct{}
	Exited chan struct{}
}

// Start will start the ticker service
func (t *Ticker) Start() {
	go t.tick()
}

// Stop will stop the ticket service, waiting
// for any any tick to first complete
func (t *Ticker) Stop() {
	close(t.Done)
	logger.Info("Waiting for service to exit")
	<-t.Exited
	logger.Info("Safely exited")
}

func (t *Ticker) tick() {
	logger.Info("Running initial service tick")
	t.Impl.Tick()
	logger.Info("Initial tick complete")

	defer close(t.Exited)
	for {
		select {
		case <-t.Done:
			return
		case <-t.Ticker.C:
			logger.Debug("Service tick")
			t.Impl.Tick()
			logger.Debug("Tick complete")
		}
	}
}
