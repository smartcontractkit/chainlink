package service

import (
	"database/sql"
	"fmt"

	"chainlink/ingester/client"
	"chainlink/ingester/logger"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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
		"eth-chain-id", config.NetworkID,
		"db-host", config.DatabaseHost,
		"db-name", config.DatabaseName,
		"db-port", config.DatabasePort,
		"db-username", config.DatabaseUsername,
	)

	ec, err := client.NewClient(config.EthereumURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ETH client: %+v", err)
	}

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

	q := ethereum.FilterQuery{}
	logChan := make(chan types.Log)
	_, err = ec.SubscribeToLogs(logChan, q)
	if err != nil {
		return nil, err
	}

	go func() {
		logger.Debug("Listening for logs")
		for log := range logChan {
			address := make([]byte, 20)
			copy(address, log.Address[:])

			topics := make([]byte, len(log.Topics)*len(common.Hash{}))
			for index, topic := range log.Topics {
				copy(topics[index*len(common.Hash{}):], topic.Bytes())
			}

			logger.Debugw("Oberved new log", "blockHash", log.BlockHash, "index", log.Index, "removed", log.Removed)
			_, err := db.Exec(`INSERT INTO "ethereum_log" ("address", "topics", "data", "blockNumber", "txHash", "txIndex", "blockHash", "index", "removed") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
				address,
				topics,
				log.Data,
				log.BlockNumber,
				log.TxHash.Bytes(),
				log.TxIndex,
				log.BlockHash.Bytes(),
				log.Index,
				log.Removed)
			if err != nil {
				logger.Errorw("Insert failed", "error", err)
			}
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

			logger.Debugw("Observed new head", "blockHeight", head.Number, "blockHash", head.Hash())
			_, err := db.Exec(`INSERT INTO "ethereum_head" ("blockHash", "parentHash", "uncleHash", "coinbase", "root", "txHash", "receiptHash", "bloom", "difficulty", "number", "gasLimit", "gasUsed", "time", "extra", "mixDigest", "nonce") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`,
				head.Hash().Bytes(),
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
