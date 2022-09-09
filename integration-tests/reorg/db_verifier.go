package reorg

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"go.uber.org/atomic"
	"go.uber.org/ratelimit"
)

type LogData struct {
	EVMChainID  int       `db:"evm_chain_id"`
	LogIndex    int64     `db:"log_index"`
	BlockHash   []byte    `db:"block_hash"`
	BlockNumber int64     `db:"block_number"`
	Address     []byte    `db:"address"`
	EventSig    []byte    `db:"event_sig"`
	Topics      []byte    `db:"topics"`
	TXHash      []byte    `db:"tx_hash"`
	Data        []byte    `db:"data"`
	CreatedAt   time.Time `db:"created_at"`
}

type DBVerifier struct {
	db              *ctfClient.PostgresConnector
	txns            []*types.Transaction
	logsPerTX       int
	fromBlockNumber uint64
}

func NewDBVerifier(db *ctfClient.PostgresConnector, txns []*types.Transaction, logsPerTX int, fromBlockNumber uint64) *DBVerifier {
	return &DBVerifier{
		db:              db,
		txns:            txns,
		logsPerTX:       logsPerTX,
		fromBlockNumber: fromBlockNumber,
	}
}

func (v *DBVerifier) verifyLogsReceived() {
	var logsCount []int
	for {
		time.Sleep(3 * time.Second)
		err := v.db.Select(&logsCount, fmt.Sprintf("select count(*) from logs where block_number > %d", v.fromBlockNumber))
		Expect(err).ShouldNot(HaveOccurred())
		log.Info().
			Int("Logs", logsCount[0]).
			Int("RequiredLogs", len(v.txns)*v.logsPerTX).
			Uint64("Block", v.fromBlockNumber).
			Msg("Checking received logs after block")
		if logsCount[0] >= len(v.txns)*v.logsPerTX {
			break
		}
	}
}

func (v *DBVerifier) fetchLogs() []LogData {
	log.Info().Uint64("Block", v.fromBlockNumber).Msg("Fetching logs from block")
	rl := ratelimit.New(10)
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	totalLogs := len(v.txns) * v.logsPerTX
	allLogs := make([]LogData, 0)
	chunkSize := 5000
	for i := 0; i <= totalLogs; i += chunkSize {
		i := i
		wg.Add(1)
		rl.Take()
		go func() {
			query := fmt.Sprintf("select * from logs where block_number > %d limit %d offset %d;", v.fromBlockNumber, chunkSize, i)
			log.Info().Int("Offset", i).Msg("Fetching logs")
			var d []LogData
			err := v.db.Select(&d, query)
			Expect(err).ShouldNot(HaveOccurred())
			mu.Lock()
			defer mu.Unlock()
			allLogs = append(allLogs, d...)
			wg.Done()
		}()
	}
	wg.Wait()
	return allLogs
}

func (v *DBVerifier) verifyTransactionsStored(logs []LogData) {
	log.Info().Msg("Verifying transactions")
	success, fail := atomic.NewInt64(0), atomic.NewInt64(0)
	for _, txn := range v.txns {
		found := false
		for _, l := range logs {
			if bytes.Equal(txn.Hash().Bytes(), l.TXHash) {
				found = true
				success.Add(1)
				break
			}
		}
		if !found {
			log.Warn().Str("Hash", txn.Hash().Hex()).Msg("Transaction was not found in DB")
			fail.Add(1)
		}
	}
	log.Warn().Int64("Success", success.Load()).Int64("Fail", fail.Load()).Msg("Verification results")
	Expect(fail.Load()).To(BeEquivalentTo(int64(0)))
	Expect(success.Load()).To(BeEquivalentTo(len(v.txns)))
}

func (v *DBVerifier) VerifyAllTransactionsStored() {
	v.verifyLogsReceived()
}
