package usdcreader

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	typepkgmock "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/types/ccipocr3"

	sel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/usdc_reader_tester"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

const ChainID = 1337

func Test_USDCReader_MessageHashes(t *testing.T) {
	finalityDepth := 5

	ctx := testutils.Context(t)
	ethereumChain := cciptypes.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)
	ethereumDomainCCTP := reader.CCTPDestDomains[uint64(ethereumChain)]
	avalancheChain := cciptypes.ChainSelector(sel.AVALANCHE_MAINNET.Selector)
	avalancheDomainCCTP := reader.CCTPDestDomains[uint64(avalancheChain)]
	polygonChain := cciptypes.ChainSelector(sel.POLYGON_MAINNET.Selector)
	polygonDomainCCTP := reader.CCTPDestDomains[uint64(polygonChain)]

	ts := testSetup(ctx, t, ethereumChain, evmconfig.USDCReaderConfig, finalityDepth, false)

	mokAddrCodec := typepkgmock.NewMockAddressCodec(t)
	mokAddrCodec.On("AddressBytesToString", mock.Anything, mock.Anything).
		Return(func(addr cciptypes.UnknownAddress, _ cciptypes.ChainSelector) string {
			return "0x" + hex.EncodeToString(addr)
		}, nil).Maybe()
	mokAddrCodec.On("AddressStringToBytes", mock.Anything, mock.Anything).
		Return(func(addr string, _ cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
			addrBytes, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(addr, "0x")))
			if err != nil {
				return nil, err
			}
			return addrBytes, nil
		}).Maybe()
	usdcReader, err := reader.NewUSDCMessageReader(
		ctx,
		logger.TestLogger(t),
		map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
			ethereumChain: {
				SourceMessageTransmitterAddr: ts.contractAddr.String(),
			},
		},
		map[cciptypes.ChainSelector]contractreader.ContractReaderFacade{
			ethereumChain: ts.reader,
		}, mokAddrCodec)
	require.NoError(t, err)

	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 11)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 21)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 31)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 41)
	emitMessageSent(t, ts, ethereumDomainCCTP, polygonDomainCCTP, 31)
	emitMessageSent(t, ts, ethereumDomainCCTP, polygonDomainCCTP, 41)
	// Finalize events
	for i := 0; i < finalityDepth; i++ {
		ts.sb.Commit()
	}
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 51)

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, ts.lp.Replay(ctx, 1))

	tt := []struct {
		name           string
		tokens         map[reader.MessageTokenID]cciptypes.RampTokenAmount
		sourceChain    cciptypes.ChainSelector
		destChain      cciptypes.ChainSelector
		expectedMsgIDs []reader.MessageTokenID
	}{
		{
			name:           "empty messages should return empty response",
			tokens:         map[reader.MessageTokenID]cciptypes.RampTokenAmount{},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "single token message",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{reader.NewMessageTokenID(1, 1)},
		},
		{
			name: "single token message but different chain",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(31, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      polygonChain,
			expectedMsgIDs: []reader.MessageTokenID{reader.NewMessageTokenID(1, 2)},
		},
		{
			name: "message without matching nonce",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(1234, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "message without matching source domain",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, avalancheDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "message with multiple tokens",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(21, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain: ethereumChain,
			destChain:   avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{
				reader.NewMessageTokenID(1, 1),
				reader.NewMessageTokenID(1, 2),
			},
		},
		{
			name: "message with multiple tokens, one without matching nonce",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(12, ethereumDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 3): {
					ExtraData: reader.NewSourceTokenDataPayload(31, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain: ethereumChain,
			destChain:   avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{
				reader.NewMessageTokenID(1, 1),
				reader.NewMessageTokenID(1, 3),
			},
		},
		{
			name: "not finalized events are not returned",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 5): {
					ExtraData: reader.NewSourceTokenDataPayload(51, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hashes, err1 := usdcReader.MessagesByTokenID(ctx, tc.sourceChain, tc.destChain, tc.tokens)
			require.NoError(t, err1)

			require.Equal(t, len(tc.expectedMsgIDs), len(hashes))

			for _, msgID := range tc.expectedMsgIDs {
				_, ok := hashes[msgID]
				require.True(t, ok)
			}
		})
	}
}

// Benchmark Results:
// Benchmark_MessageHashes/Small_Dataset-14       3723        272421 ns/op       126949 B/op      2508 allocs/op
// Benchmark_MessageHashes/Medium_Dataset-14       196       6164706 ns/op      1501435 B/op     20274 allocs/op
// Benchmark_MessageHashes/Large_Dataset-14          7     163930268 ns/op     37193160 B/op    463954 allocs/op
//
// Notes:
// - Small dataset processes 3,723 iterations with 126KB memory usage per iteration.
// - Medium dataset processes 196 iterations with 1.5MB memory usage per iteration.
// - Large dataset processes only 7 iterations with ~37MB memory usage per iteration.
func Benchmark_MessageHashes(b *testing.B) {
	finalityDepth := 5

	// Adding a new parameter: tokenCount
	testCases := []struct {
		name       string
		msgCount   int
		startNonce int64
		tokenCount int
	}{
		{"Small_Dataset", 100, 1, 5},
		{"Medium_Dataset", 10_000, 1, 10},
		{"Large_Dataset", 100_000, 1, 50},
	}

	mokAddrCodec := typepkgmock.NewMockAddressCodec(b)
	mokAddrCodec.On("AddressBytesToString", mock.Anything, mock.Anything).
		Return(func(addr cciptypes.UnknownAddress, _ cciptypes.ChainSelector) string {
			return "0x" + hex.EncodeToString(addr)
		}, nil).Maybe()
	mokAddrCodec.On("AddressStringToBytes", mock.Anything, mock.Anything).
		Return(func(addr string, _ cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
			addrBytes, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(addr, "0x")))
			if err != nil {
				return nil, err
			}
			return addrBytes, nil
		}).Maybe()

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			ctx := testutils.Context(b)
			sourceChain := cciptypes.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)
			sourceDomainCCTP := reader.CCTPDestDomains[uint64(sourceChain)]
			destChain := cciptypes.ChainSelector(sel.AVALANCHE_MAINNET.Selector)
			destDomainCCTP := reader.CCTPDestDomains[uint64(destChain)]

			ts := testSetup(ctx, b, sourceChain, evmconfig.USDCReaderConfig, finalityDepth, true)

			usdcReader, err := reader.NewUSDCMessageReader(
				ctx,
				logger.TestLogger(b),
				map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
					sourceChain: {
						SourceMessageTransmitterAddr: ts.contractAddr.String(),
					},
				},
				map[cciptypes.ChainSelector]contractreader.ContractReaderFacade{
					sourceChain: ts.reader,
				}, mokAddrCodec)
			require.NoError(b, err)

			// Populate the database with the specified number of logs
			populateDatabase(b, ts, sourceChain, sourceDomainCCTP, destDomainCCTP, tc.startNonce, tc.msgCount, finalityDepth)

			// Create a map of tokens to query for, with the specified tokenCount
			tokens := make(map[reader.MessageTokenID]cciptypes.RampTokenAmount)
			for i := 1; i <= tc.tokenCount; i++ {
				//nolint:gosec // disable G115
				tokens[reader.NewMessageTokenID(cciptypes.SeqNum(i), 1)] = cciptypes.RampTokenAmount{
					ExtraData: reader.NewSourceTokenDataPayload(uint64(tc.startNonce)+uint64(i), sourceDomainCCTP).ToBytes(),
				}
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				hashes, err := usdcReader.MessagesByTokenID(ctx, sourceChain, destChain, tokens)
				require.NoError(b, err)
				require.Len(b, hashes, tc.tokenCount) // Ensure the number of matches is as expected
			}
		})
	}
}

func populateDatabase(b *testing.B,
	testEnv *testSetupData,
	source cciptypes.ChainSelector,
	sourceDomainCCTP uint32,
	destDomainCCTP uint32,
	startNonce int64,
	numOfMessages int,
	finalityDepth int) {
	ctx := testutils.Context(b)

	abi, err := usdc_reader_tester.USDCReaderTesterMetaData.GetAbi()
	require.NoError(b, err)

	var logs []logpoller.Log
	messageSentEventSig := abi.Events["MessageSent"].ID
	require.NoError(b, err)
	messageTransmitterAddress := testEnv.contractAddr

	for i := 0; i < numOfMessages; i++ {
		// Create topics array with just the event signature
		topics := [][]byte{
			messageSentEventSig[:], // Topic[0] is event signature
		}

		// Create log entry
		logs = append(logs, logpoller.Log{
			EVMChainID:     ubig.New(new(big.Int).SetUint64(uint64(source))),
			LogIndex:       int64(i + 1),
			BlockHash:      utils.NewHash(),
			BlockNumber:    int64(i + 1),
			BlockTimestamp: time.Now(),
			EventSig:       messageSentEventSig,
			Topics:         topics,
			Address:        messageTransmitterAddress,
			TxHash:         utils.NewHash(),
			Data:           createMessageSentLogPollerData(startNonce, i, sourceDomainCCTP, destDomainCCTP),
			CreatedAt:      time.Now(),
		})
	}

	require.NoError(b, testEnv.orm.InsertLogs(ctx, logs))
	require.NoError(b, testEnv.orm.InsertBlock(ctx, utils.RandomHash(), int64(numOfMessages+finalityDepth), time.Now(), int64(numOfMessages+finalityDepth)))
}

func createMessageSentLogPollerData(startNonce int64, i int, sourceDomainCCTP uint32, destDomainCCTP uint32) []byte {
	nonce := int(startNonce) + i

	var buf []byte

	buf = binary.BigEndian.AppendUint32(buf, reader.CCTPMessageVersion)

	buf = binary.BigEndian.AppendUint32(buf, sourceDomainCCTP)

	buf = binary.BigEndian.AppendUint32(buf, destDomainCCTP)
	// #nosec G115
	buf = binary.BigEndian.AppendUint64(buf, uint64(nonce))

	senderBytes := [12]byte{}
	buf = append(buf, senderBytes[:]...)

	var message [32]byte
	copy(message[:], buf)

	data := make([]byte, 0)

	offsetBytes := make([]byte, 32)
	binary.BigEndian.PutUint64(offsetBytes[24:], 32)
	data = append(data, offsetBytes...)

	lengthBytes := make([]byte, 32)
	binary.BigEndian.PutUint64(lengthBytes[24:], uint64(len(message)))
	data = append(data, lengthBytes...)

	data = append(data, message[:]...)
	return data
}

// we might want to use batching (evm/batching or evm/batching) but might be slow
func emitMessageSent(t *testing.T, testEnv *testSetupData, source, dest uint32, nonce uint64) {
	payload := utils.RandomBytes32()
	_, err := testEnv.contract.EmitMessageSent(
		testEnv.auth,
		reader.CCTPMessageVersion,
		source,
		dest,
		utils.RandomBytes32(),
		utils.RandomBytes32(),
		[32]byte{},
		nonce,
		payload[:],
	)
	require.NoError(t, err)
	testEnv.sb.Commit()
}

func testSetup(ctx context.Context, t testing.TB, readerChain cciptypes.ChainSelector, cfg evmtypes.ChainReaderConfig, depth int, useHeavyDB bool) *testSetupData {
	// Generate a new key pair for the simulated account
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	// Set up the genesis account with balance
	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	assert.True(t, ok)
	alloc := map[common.Address]gethtypes.Account{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := simulated.NewBackend(alloc, simulated.WithBlockGasLimit(0))
	// Create a transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ChainID))
	require.NoError(t, err)
	auth.GasLimit = uint64(0)

	address, _, _, err := usdc_reader_tester.DeployUSDCReaderTester(
		auth,
		simulatedBackend.Client(),
	)
	require.NoError(t, err)
	simulatedBackend.Commit()

	contract, err := usdc_reader_tester.NewUSDCReaderTester(address, simulatedBackend.Client())
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)

	// Parameterize database selection
	var db *sqlx.DB
	if useHeavyDB {
		_, db = heavyweight.FullTestDBV2(t, nil) // Use heavyweight database for benchmarks
	} else {
		db = pgtest.NewSqlxDB(t) // Use simple in-memory DB for tests
	}

	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            int64(depth),
		BackfillBatchSize:        10,
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(0).SetUint64(uint64(readerChain)))
	headTracker := headstest.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	orm := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(readerChain)), db, lggr)

	lp := logpoller.NewLogPoller(
		orm,
		cl,
		lggr,
		headTracker,
		lpOpts,
	)
	require.NoError(t, lp.Start(ctx))

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, cfg)
	require.NoError(t, err)

	err = cr.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, cr.Close())
		require.NoError(t, lp.Close())
		require.NoError(t, db.Close())
	})

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         auth,
		cl:           cl,
		reader:       cr,
		orm:          orm,
		db:           db,
		lp:           lp,
	}
}

type testSetupData struct {
	contractAddr common.Address
	contract     *usdc_reader_tester.USDCReaderTester
	sb           *simulated.Backend
	auth         *bind.TransactOpts
	cl           client.Client
	reader       types.ContractReader
	orm          logpoller.ORM
	db           *sqlx.DB
	lp           logpoller.LogPoller
}
