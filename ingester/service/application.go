package service

import (
	"database/sql"
	"fmt"

	"ingester/client"
	"ingester/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
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
	logger.SetLogger(logger.CreateTestLogger(-1))

	logger.Infow(
		"Starting the Chainlink Ingester",
		"eth-url", config.EthereumURL,
		"db-url", config.DatabaseURL,
		"eth-chain-id", config.NetworkID)

	ec, err := client.NewClient(config.EthereumURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ETH client: %+v", err)
	}

	pool, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	q := ethereum.FilterQuery{}
	logChan := make(chan types.Log)
	_, err = ec.SubscribeToLogs(logChan, q)
	if err != nil {
		return nil, err
	}

	go func() {
		logger.Debug("Listening for logs")
		for log := range logChan {
			logger.Debugw("Got Log",
				"address", log.Address.Hex(),
				"topics", log.Topics,
				"data", log.Data,
				"blockNumber", log.BlockNumber,
				"txHash", log.TxHash,
				"txIndex", log.TxIndex,
				"blockHash", log.BlockHash,
				"index", log.Index,
				"removed", log.Removed)
		}
	}()

	headChan := make(chan types.Header)
	_, err = ec.SubscribeToNewHeads(headChan)
	if err != nil {
		return nil, err
	}

	go func() {
		logger.Debug("Listening for heads")
		for head := range headChan {
			nonce := make([]byte, 8)
			copy(nonce, head.Nonce[:])

			logger.Debugw("Got head", "head", head)
			_, err := pool.Exec(`INSERT INTO "ethereum_head" (parent_hash, uncle_hash, coinbase, root, tx_hash, receipt_hash, bloom, difficulty, number, gas_limit, gas_used, time, extra, mix_digest, nonce) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15);`,
				head.ParentHash,
				head.UncleHash,
				head.Coinbase,
				head.Root,
				head.TxHash,
				head.ReceiptHash,
				head.Bloom.Bytes(),
				head.Difficulty.String(),
				head.Number.String(),
				head.GasLimit,
				head.GasUsed,
				head.Time,
				head.Extra,
				head.MixDigest,
				nonce)
			if err != nil {
				logger.Errorw("Insert failed", "error", err)
			}
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
	logger.Info("Shutting down")
}
