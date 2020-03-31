package service

import (
	"database/sql"
	"fmt"

	"chainlink/ingester/client"
	"chainlink/ingester/logger"

	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
)

// Application is an instance of the aggregator monitor application containing
// all clients and services
type Application struct {
	Heads       HeadsTracker
	Feeds       FeedsTracker
	Submissions SubmissionsTracker
	Config      *Config

	ETHClient client.ETH
}

// InterruptHandler is a function that is called after application startup
// designed to wait based on a specified interrupt
type InterruptHandler func()

// NewApplication returns an instance of the Application with
// all clients connected and services instantiated
func NewApplication(config *Config) (*Application, error) {
	logger.SetLogger(logger.CreateTestLogger(-1))

	logger.Infow(
		"Starting the Chainlink Ingester",
		"eth-url", config.EthereumURL,
		"eth-chain-id", config.NetworkID,
		"db-host", config.DatabaseHost,
		"db-name", config.DatabaseName,
		"db-port", config.DatabasePort,
		"db-username", config.DatabaseUsername,
	)

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUsername,
		config.DatabasePassword,
		config.DatabaseName,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	ec, err := client.NewClient(config.EthereumURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ETH client: %+v", err)
	}

	ht := NewHeadsTracker(db, ec)
	ft := NewFeedsTracker(ec, client.NewFeedsUI(config.FeedsUIURL), config.NetworkID)
	st := NewSubmissionsTracker(db, ec, ft)

	return &Application{
		ETHClient:   ec,
		Config:      config,
		Heads:       ht,
		Feeds:       ft,
		Submissions: st,
	}, nil
}

// Start will start all the services within the application and call the interrupt handler
func (a *Application) Start(ih InterruptHandler) {
	a.Heads.Start()
	a.Submissions.Start()
	a.Feeds.Start()

	ih()
}

// Stop will call each services that requires a clean shutdown to stop
func (a *Application) Stop() {
	a.Heads.Stop()
	logger.Info("Shutting down")
}
