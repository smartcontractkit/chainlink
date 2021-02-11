package log_test

import (
	"bytes"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func TestSubscriber(t *testing.T) {
	t.Parallel()

	g := NewGomegaWithT(t)

	const (
		backfillDepth uint64 = 5
	)

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	store.Config.Set("BLOCK_BACKFILL_DEPTH", backfillDepth)

	var (
		ethClient = new(mocks.Client)
		sub       = new(mocks.Subscription)
		addrs     = map[int]common.Address{
			0: cltest.NewAddress(),
			1: cltest.NewAddress(),
			2: cltest.NewAddress(),
			3: cltest.NewAddress(), // `addr3` exists to simulate geth misbehaving by sending some logs that weren't requested
		}
		currentBlockNumber uint64 = 12
		firstBlock         uint64 = currentBlockNumber - backfillDepth + 1
		lastBlock          uint64 = firstBlock * 6
		logs               []types.Log

		dependentAwaiter = utils.NewDependentAwaiter()
		orm              = log.NewORM(store.DB)
		relayer          = log.ExportedNewRelayer(orm, store.Config, dependentAwaiter)
		subscriber       = log.ExportedNewSubscriber(orm, ethClient, store.Config, relayer, dependentAwaiter)

		chRawLogs   chan<- types.Log
		chchRawLogs = make(chan chan<- types.Log, 1)
	)
	store.EthClient = ethClient

	for i := firstBlock; i <= lastBlock; i++ {
		logs = append(logs, types.Log{Address: addrs[int(i%3)+1], BlockNumber: i, BlockHash: cltest.NewHash(), Topics: []common.Hash{}, Data: []byte{}})
	}

	expectSubscribe := func(t *testing.T, expectedFilterAddrs []common.Address) {
		t.Helper()

		ethClient.On("SubscribeFilterLogs", mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				query := args.Get(1).(ethereum.FilterQuery)
				require.Nil(t, query.FromBlock)
				require.Len(t, query.Addresses, len(expectedFilterAddrs))
				for _, addr := range expectedFilterAddrs {
					g.Expect(query.Addresses).Should(ContainElement(addr))
				}
				select {
				case chchRawLogs <- args.Get(2).(chan<- types.Log):
				case <-time.After(5 * time.Second):
					t.Fatal("could not send log channel")
				}
			}).
			Return(sub, nil).
			Once()
		sub.On("Err").Return(nil)
		sub.On("Unsubscribe").Return()

		latestBlockCall := ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Once()
		latestBlockCall.Run(func(mock.Arguments) {
			latestBlockCall.ReturnArguments = mock.Arguments{&models.Head{Number: int64(currentBlockNumber)}, nil}
		})

		filterCall := ethClient.On("FilterLogs", mock.Anything, mock.Anything).Once()
		filterCall.Run(func(args mock.Arguments) {
			query := args.Get(1).(ethereum.FilterQuery)
			require.Equal(t, big.NewInt(int64(currentBlockNumber-backfillDepth)), query.FromBlock)
			require.Len(t, query.Addresses, len(expectedFilterAddrs))
			for _, addr := range expectedFilterAddrs {
				g.Expect(query.Addresses).Should(ContainElement(addr))
			}
			var filteredLogs []types.Log
			for _, log := range logs {
				if int64(log.BlockNumber) > query.FromBlock.Int64() && log.BlockNumber <= currentBlockNumber {
					filteredLogs = append(filteredLogs, log)
				}
			}
			filterCall.ReturnArguments = mock.Arguments{filteredLogs, nil}
		})
	}

	subscriber.Start()
	defer subscriber.Stop()

	t.Run("adds contracts", func(t *testing.T) {
		expectSubscribe(t, []common.Address{addrs[0]})

		subscriber.NotifyAddContract(addrs[0])

		g.Eventually(func() int { return len(subscriber.ExportedContracts()) }).Should(Equal(1))
		g.Consistently(func() int { return len(subscriber.ExportedContracts()) }).Should(Equal(1))
		g.Expect(subscriber.ExportedContracts()).Should(HaveLen(1))
		g.Expect(subscriber.ExportedContracts()).Should(HaveKeyWithValue(addrs[0], uint64(1)))

		select {
		case chRawLogs = <-chchRawLogs:
		case <-time.After(5 * time.Second):
			t.Fatal("did not subscribe")
		}

		cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)
	})

	t.Run("does not save logs to the DB for which there are no subscribers", func(t *testing.T) {
		defer store.DB.Exec(`DELETE FROM eth_logs`)

		require.NotNil(t, chRawLogs, "failed to subscribe in previous test")

		sendLogs(t, logs, chRawLogs)
		time.Sleep(5 * time.Second)

		var dbLogs []types.Log
		err := store.DB.Raw(`SELECT * FROM eth_logs`).Scan(&dbLogs).Error
		require.NoError(t, err)
		require.Len(t, dbLogs, 0)
	})

	t.Run("saves the correct logs to the DB when there are subscribers for them", func(t *testing.T) {
		expectSubscribe(t, []common.Address{addrs[0], addrs[1]})

		subscriber.NotifyAddContract(addrs[1])

		g.Eventually(func() int { return len(subscriber.ExportedContracts()) }).Should(Equal(2))
		g.Consistently(func() int { return len(subscriber.ExportedContracts()) }).Should(Equal(2))
		g.Expect(subscriber.ExportedContracts()).Should(HaveLen(2))
		g.Expect(subscriber.ExportedContracts()).Should(HaveKeyWithValue(addrs[0], uint64(1)))
		g.Expect(subscriber.ExportedContracts()).Should(HaveKeyWithValue(addrs[1], uint64(1)))

		select {
		case chRawLogs = <-chchRawLogs:
		case <-time.After(5 * time.Second):
			t.Fatal("did not subscribe")
		}

		cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)

		for currentBlockNumber <= lastBlock/2 {
			expected := filterLogs(logs, map[common.Address]fromTo{
				addrs[1]: {0, currentBlockNumber},
			})

			toSend := filterLogs(logs, map[common.Address]fromTo{
				addrs[1]: {currentBlockNumber, currentBlockNumber},
				addrs[2]: {currentBlockNumber, currentBlockNumber},
				addrs[3]: {currentBlockNumber, currentBlockNumber},
			})

			sendLogs(t, toSend, chRawLogs)

			var dbLogs []types.Log
			var err error
			g.Eventually(func() []types.Log {
				dbLogs, err = log.FetchLogs(store.DB, `SELECT eth_logs.block_hash, eth_logs.block_number, eth_logs.index, eth_logs.address, eth_logs.topics, eth_logs.data FROM eth_logs ORDER BY block_number, address ASC`)
				require.NoError(t, err)
				return dbLogs
			}).Should(HaveLen(len(expected)))

			sortLogs(dbLogs)
			sortLogs(expected)
			require.Equal(t, expected, dbLogs)

			currentBlockNumber++
		}

		//
		// Add the final subscriber, and also duplicate addr1's subscription
		//

		expectSubscribe(t, []common.Address{addrs[0], addrs[1], addrs[2]})

		subscriber.NotifyAddContract(addrs[1])
		subscriber.NotifyAddContract(addrs[2])

		g.Eventually(func() int { return len(subscriber.ExportedContracts()) }).Should(Equal(3))
		g.Consistently(func() int { return len(subscriber.ExportedContracts()) }).Should(Equal(3))
		g.Expect(subscriber.ExportedContracts()).Should(HaveLen(3))
		g.Expect(subscriber.ExportedContracts()).Should(HaveKeyWithValue(addrs[0], uint64(1)))
		g.Expect(subscriber.ExportedContracts()).Should(HaveKeyWithValue(addrs[1], uint64(2)))
		g.Expect(subscriber.ExportedContracts()).Should(HaveKeyWithValue(addrs[2], uint64(1)))

		select {
		case chRawLogs = <-chchRawLogs:
		case <-time.After(5 * time.Second):
			t.Fatal("did not subscribe")
		}

		cltest.EventuallyExpectationsMet(t, ethClient, 5*time.Second, 10*time.Millisecond)

		addr2FirstBlock := currentBlockNumber - backfillDepth

		for currentBlockNumber <= lastBlock {
			expected := filterLogs(logs, map[common.Address]fromTo{
				addrs[1]: {0, currentBlockNumber},
				addrs[2]: {addr2FirstBlock, currentBlockNumber},
			})

			toSend := filterLogs(logs, map[common.Address]fromTo{
				addrs[1]: {currentBlockNumber, currentBlockNumber},
				addrs[2]: {currentBlockNumber, currentBlockNumber},
				addrs[3]: {currentBlockNumber, currentBlockNumber},
			})

			sendLogs(t, toSend, chRawLogs)

			var dbLogs []types.Log
			var err error
			g.Eventually(func() []types.Log {
				dbLogs, err = log.FetchLogs(store.DB, `SELECT eth_logs.block_hash, eth_logs.block_number, eth_logs.index, eth_logs.address, eth_logs.topics, eth_logs.data FROM eth_logs ORDER BY block_number, address ASC`)
				require.NoError(t, err)
				return dbLogs
			}).Should(HaveLen(len(expected)))

			sortLogs(dbLogs)
			sortLogs(expected)
			require.Equal(t, expected, dbLogs)

			currentBlockNumber++
		}
	})
}

func sendLogs(t *testing.T, logs []types.Log, chRawLogs chan<- types.Log) {
	t.Helper()

	for _, log := range logs {
		select {
		case chRawLogs <- log:
		case <-time.After(5 * time.Second):
			t.Fatal("timed out while sending")
		}
	}
}

func sortLogs(logs []types.Log) {
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber < logs[j].BlockNumber {
			return true
		}
		return logs[i].BlockNumber == logs[j].BlockNumber && bytes.Compare(logs[i].Address[:], logs[j].Address[:]) < 0
	})
}

type fromTo struct {
	fromBlock, toBlock uint64
}

func filterLogs(logs []types.Log, addrs map[common.Address]fromTo) []types.Log {
	var filtered []types.Log
	for _, log := range logs {
		fromTo, exists := addrs[log.Address]
		if !exists {
			continue
		}
		if log.BlockNumber >= fromTo.fromBlock && log.BlockNumber <= fromTo.toBlock {
			filtered = append(filtered, log)
		}
	}
	return filtered
}
