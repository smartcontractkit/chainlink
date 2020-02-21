package service

import (
	"fmt"

	"ingester/client"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
)

// Application is an instance of the aggregator monitor application containing
// all clients and services
type Application struct {
	Config *Config

	ETHClient client.ETH
}

// InterruptHandler is a function that is called after application startup
// designed to wait based on a specified interrupt
type InterruptHandler func()

// NewApplication returns an instance of the Application with
// all clients connected and services instantiated
func NewApplication(config *Config) (*Application, error) {
	log.Info("Starting the ingester")

	ec, err := client.NewClient(config.EthereumURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ETH client: %+v", err)
	}

	q := ethereum.FilterQuery{}
	logChan := make(chan types.Log)
	_, err = ec.SubscribeToLogs(logChan, q)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			log := <-logChan
			fmt.Println("got log", log.Address.Hex(), log.Topics, log.Data, log.BlockNumber, log.TxHash, log.TxIndex, log.BlockHash, log.Index, log.Removed)
		}
	}()

	headChan := make(chan client.BlockHeader)
	_, err = ec.SubscribeToNewHeads(headChan)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			head := <-headChan
			fmt.Println("got head", head)
		}
	}()

	return &Application{
		ETHClient: ec,
		Config:    config,
	}, nil
}

// Start will start all the services within the application and call the interrupt handler
func (a *Application) Start(ih InterruptHandler) {
	ih()
}

// Stop will call each services that requires a clean shutdown to stop
func (a *Application) Stop() {
	log.Info("Stopping the ingester")
}
