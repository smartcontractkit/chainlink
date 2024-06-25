//go:build playground
// +build playground

package chainreader

import (
	"context"
	_ "embed"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	types2 "github.com/smartcontractkit/chainlink-common/pkg/types"
	query2 "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	logger2 "github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/stretchr/testify/assert"
)

const chainID = 1337

type testSetupData struct {
	contractAddr common.Address
	contract     *Chainreader
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
}

func TestChainReader(t *testing.T) {
	ctx := testutils.Context(t)
	lggr := logger2.NullLogger
	d := testSetup(t, ctx)

	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            0,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, d.sb, big.NewInt(chainID))
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(chainID), db, lggr), cl, lggr, lpOpts)
	assert.NoError(t, lp.Start(ctx))

	const (
		ContractNameAlias = "myCoolContract"

		FnAliasGetCount = "myCoolFunction"
		FnGetCount      = "getEventCount"

		FnAliasGetNumbers = "GetNumbers"
		FnGetNumbers      = "getNumbers"

		FnAliasGetPerson = "GetPerson"
		FnGetPerson      = "getPerson"

		EventNameAlias = "myCoolEvent"
		EventName      = "SimpleEvent"
	)

	// Initialize chainReader
	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			ContractNameAlias: {
				ContractABI: ChainreaderMetaData.ABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					EventNameAlias: {
						ChainSpecificName:       EventName,
						ReadType:                evmtypes.Event,
						ConfidenceConfirmations: map[string]int{"0.0": 0, "1.0": 0},
					},
					FnAliasGetCount: {
						ChainSpecificName: FnGetCount,
					},
					FnAliasGetNumbers: {
						ChainSpecificName:   FnGetNumbers,
						OutputModifications: codec.ModifiersConfig{},
					},
					FnAliasGetPerson: {
						ChainSpecificName: FnGetPerson,
						OutputModifications: codec.ModifiersConfig{
							&codec.RenameModifierConfig{
								Fields: map[string]string{"Name": "NameField"}, // solidity name -> go struct name
							},
						},
					},
				},
			},
		},
	}

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, cl, cfg)
	assert.NoError(t, err)
	err = cr.Bind(ctx, []types2.BoundContract{
		{
			Address: d.contractAddr.String(),
			Name:    ContractNameAlias,
			Pending: false,
		},
	})
	assert.NoError(t, err)

	err = cr.Start(ctx)
	assert.NoError(t, err)
	for {
		if err := cr.Ready(); err == nil {
			break
		}
	}

	emitEvents(t, d, ctx) // Calls the contract to emit events

	// (hack) Sometimes LP logs are missing, commit several times and wait few seconds to make it work.
	for i := 0; i < 100; i++ {
		d.sb.Commit()
	}
	time.Sleep(5 * time.Second)

	t.Run("simple contract read", func(t *testing.T) {
		var cnt big.Int
		err = cr.GetLatestValue(ctx, ContractNameAlias, FnAliasGetCount, map[string]interface{}{}, &cnt)
		assert.NoError(t, err)
		assert.Equal(t, int64(10), cnt.Int64())
	})

	t.Run("read array", func(t *testing.T) {
		var nums []big.Int
		err = cr.GetLatestValue(ctx, ContractNameAlias, FnAliasGetNumbers, map[string]interface{}{}, &nums)
		assert.NoError(t, err)
		assert.Len(t, nums, 10)
		for i := 1; i <= 10; i++ {
			assert.Equal(t, int64(i), nums[i-1].Int64())
		}
	})

	t.Run("read struct", func(t *testing.T) {
		person := struct {
			NameField string
			Age       *big.Int // WARN: specifying a wrong data type e.g. int instead of *big.Int fails silently with a default value of 0
		}{}
		err = cr.GetLatestValue(ctx, ContractNameAlias, FnAliasGetPerson, map[string]interface{}{}, &person)
		assert.Equal(t, "Dim", person.NameField)
		assert.Equal(t, int64(18), person.Age.Int64())
	})

	t.Run("read events", func(t *testing.T) {
		var myDataType *big.Int
		seq, err := cr.QueryKey(
			ctx,
			ContractNameAlias,
			query2.KeyFilter{
				Key:         EventNameAlias,
				Expressions: []query2.Expression{},
			},
			query2.LimitAndSort{},
			myDataType,
		)
		assert.NoError(t, err)
		assert.Equal(t, 10, len(seq), "expected 10 events from chain reader")
		for _, v := range seq {
			// TODO: for some reason log poller does not populate event data
			t.Logf("(chain reader) got event: (data=%v) (hash=%x)", v.Data, v.Hash)
		}
	})
}

func testSetup(t *testing.T, ctx context.Context) *testSetupData {
	// Generate a new key pair for the simulated account
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	// Set up the genesis account with balance
	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	assert.True(t, ok)
	alloc := map[common.Address]core.GenesisAccount{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := backends.NewSimulatedBackend(alloc, 0)
	// Create a transactor

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	assert.NoError(t, err)
	auth.GasLimit = uint64(0)

	// Deploy the contract
	address, tx, _, err := DeployChainreader(auth, simulatedBackend)
	assert.NoError(t, err)
	simulatedBackend.Commit()
	t.Logf("contract deployed: addr=%s tx=%s", address.Hex(), tx.Hash())

	// Setup contract client
	contract, err := NewChainreader(address, simulatedBackend)
	assert.NoError(t, err)

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         auth,
	}
}

func emitEvents(t *testing.T, d *testSetupData, ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Start emitting events
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			_, err := d.contract.EmitEvent(d.auth)
			assert.NoError(t, err)
			d.sb.Commit()
		}
	}()

	// Listen events using go-ethereum lib
	go func() {
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(0),
			Addresses: []common.Address{d.contractAddr},
		}
		logs := make(chan types.Log)
		sub, err := d.sb.SubscribeFilterLogs(ctx, query, logs)
		assert.NoError(t, err)

		numLogs := 0
		defer wg.Done()
		for {
			// Wait for the events
			select {
			case err := <-sub.Err():
				assert.NoError(t, err, "got an unexpected error")
			case vLog := <-logs:
				assert.Equal(t, d.contractAddr, vLog.Address, "got an unexpected address")
				t.Logf("(geth) got new log (cnt=%d) (data=%x) (topics=%s)", numLogs, vLog.Data, vLog.Topics)
				numLogs++
				if numLogs == 10 {
					return
				}
			}
		}
	}()

	wg.Wait() // wait for all the events to be consumed
}
