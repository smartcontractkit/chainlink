package svr

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	"github.com/smartcontractkit/wsrpc/credentials"

	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"

	"github.com/ethereum/go-ethereum/ethclient/simulated"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	logtoml "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var linkTokenAddress = common.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB")

type node struct {
	app                  chainlink.Application
	keyBundle            ocr2key.KeyBundle
	observedLogs         *observer.ObservedLogs
	transmitter          common.Address // the node's primary EOA address
	secondaryTransmitter common.Address // the node's secondary EOA address
	effectiveTransmitter common.Address // the node's forwarder address
}

func setupNodes(t *testing.T, nNodes int, transactOpts *bind.TransactOpts, backend *simulated.Backend, clientCSAKeys []csakey.KeyV2) (oracles []confighelper.OracleIdentityExtra, nodes []node) {
	ports := freeport.GetN(t, nNodes)
	for i := range nNodes {
		app, peerID, _, _, observedLogs := setupNode(t, ports[i], fmt.Sprintf("core_node_%d", i), backend, clientCSAKeys[i])
		t.Logf("Node %d with peer id %s started on port %d", i, peerID, ports[i])

		sendingKeys, err := app.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Len(t, sendingKeys, 2)
		primaryTransmitterAddr := sendingKeys[0].Address
		err = fundAddress(primaryTransmitterAddr, transactOpts, backend)
		require.NoError(t, err)
		backend.Commit()
		t.Logf("Funded primary transmitter for node %d: %s", i, primaryTransmitterAddr.String())

		secondaryTransmitterAddr := sendingKeys[1].Address
		err = fundAddress(secondaryTransmitterAddr, transactOpts, backend)
		require.NoError(t, err)
		backend.Commit()
		t.Logf("Funded secondary transmitter for node %d: %s", i, secondaryTransmitterAddr.String())

		// set up the secondary transmitter key
		// secondaryTransmitter, err := app.GetKeyStore().Eth().Create(testutils.Context(t), testutils.SimulatedChainID)
		// require.NoErrorf(t, err, "could not create secondary transmitter key for node %d", i)

		// // Explicitly enable the secondary transmitter for the chain
		// err = app.GetKeyStore().Eth().Enable(testutils.Context(t), secondaryTransmitter.Address, testutils.SimulatedChainID)
		// require.NoErrorf(t, err, "could not enable secondary transmitter key for node %d", i)

		// chainService, err := app.GetRelayers().LegacyEVMChains().Get(testutils.SimulatedChainID.String())
		// require.NoError(t, err)
		// chain, ok := chainService.(legacyevm.Chain)
		// require.True(t, ok)

		// txm := chain.TxManager()
		// txm.GetForwarderForEOAOCR2Feeds(testutils.Context(t), secondaryTransmitter.Address, common.Address{})
		// t.Logf("Closing txm for node %d", i)
		// require.NoError(t, txm.Close())
		// t.Logf("Starting txm for node %d", i)
		// require.NoError(t, txm.Start(testutils.Context(t)))

		// err = app.Stop()
		// require.NoError(t, err)
		// app.Start(testutils.Context(t))
		// require.NoError(t, err)

		// err = fundAddress(secondaryTransmitter.Address, transactOpts, backend)
		// require.NoError(t, err, "Funding secondary transmitter shouldn't fail for node %d", i)
		// backend.Commit()
		// t.Logf("Funded secondary transmitter for node %d: %s", i, secondaryTransmitter.Address.String())

		addresses, err := app.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		t.Logf("Enabled addresses for node %d: %v", i, addresses)
		require.Len(t, addresses, 2)

		kb, err := app.GetKeyStore().OCR2().Create(testutils.Context(t), "evm")
		require.NoError(t, err)

		// deploy a forwarder
		forwarderAddress, _, authorizedForwarder, err := authorized_forwarder.DeployAuthorizedForwarder(transactOpts, backend.Client(), linkTokenAddress, transactOpts.From, common.Address{}, []byte{})
		require.NoError(t, err)
		backend.Commit()

		// set primary and secondary EOA as an authorized sender for the forwarder
		_, err = authorizedForwarder.SetAuthorizedSenders(transactOpts, []common.Address{primaryTransmitterAddr, secondaryTransmitterAddr})
		require.NoError(t, err)
		backend.Commit()

		// add forwarder address to be tracked in db
		forwarderORM := forwarders.NewORM(app.GetDB())
		chainID, err := backend.Client().ChainID(testutils.Context(t))
		require.NoError(t, err)
		_, err = forwarderORM.CreateForwarder(testutils.Context(t), forwarderAddress, ubig.Big(*chainID))
		require.NoError(t, err)

		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  kb.PublicKey(),
				TransmitAccount:   ocr2types.Account(forwarderAddress.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})

		nodes = append(nodes, node{
			app:                  app,
			keyBundle:            kb,
			observedLogs:         observedLogs,
			transmitter:          primaryTransmitterAddr,
			secondaryTransmitter: secondaryTransmitterAddr,
			effectiveTransmitter: forwarderAddress,
		})
	}
	return
}

func setupNode(
	t *testing.T,
	port int,
	nodeName string,
	backend *simulated.Backend,
	csaKey csakey.KeyV2,
) (app chainlink.Application, peerID string, clientPubKey credentials.StaticSizedPublicKey, ocr2kb ocr2key.KeyBundle, observedLogs *observer.ObservedLogs) {
	k := big.NewInt(int64(port)) // keys unique to port
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(k)
	t.Logf("GEERT p2pkey is %s", p2pKey.PeerID().Raw())
	rdr := keystest.NewRandReaderFromSeed(int64(port))
	ocr2kb = ocr2key.MustNewInsecure(rdr, chaintype.EVM)

	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}

	tomlNode := toml.Node{
		Name:              ptr(nodeName),
		HTTPURL:           ptr(commonconfig.URL{Scheme: "http", Host: fmt.Sprintf("127.0.0.1:%d", port)}),
		SendOnly:          ptr(false),
		Order:             ptr(int32(1)),
		IsLoadBalancedRPC: ptr(false),
	}

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {

		// [Insecure]
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		// [JobPipeline]
		c.JobPipeline.MaxSuccessfulRuns = ptr(uint64(0))
		c.JobPipeline.VerboseLogging = ptr(true)

		// [Feature]
		c.Feature.UICSAKeys = ptr(true)
		c.Feature.LogPoller = ptr(true)
		c.Feature.FeedsManager = ptr(false)

		// [OCR]
		c.OCR.Enabled = ptr(false)

		// [OCR2]
		c.OCR2.Enabled = ptr(true)
		c.OCR2.ContractPollInterval = commonconfig.MustNewDuration(1 * time.Second)
		c.OCR2.CaptureEATelemetry = ptr(false)

		// [P2P]
		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.TraceLogging = ptr(false)

		// [P2P.V2]
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		// [Mercury] // TODO(gg): remove
		c.Mercury.VerboseLogging = ptr(true)
		c.Mercury.Transmitter.ReaperFrequency = commonconfig.MustNewDuration(0 * time.Millisecond)

		// [EVM]
		c.EVM[0].Nodes = toml.EVMNodes{&tomlNode}
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(5 * time.Second)
		c.EVM[0].Transactions.ForwardersEnabled = ptr(true)

		// [EVM.Transactions]
		// ForwardersEnabled = true

		// [EVM.Transactions.TransactionManagerV2]
		c.EVM[0].Transactions.TransactionManagerV2.Enabled = ptr(true)
		c.EVM[0].Transactions.TransactionManagerV2.BlockTime = commonconfig.MustNewDuration(11 * time.Second)
		// c.EVM[0].Transactions.TransactionManagerV2.CustomURL = ptr(commonconfig.URL{Scheme: "https", Host: "rpc-sepolia.flashbots.net", Path: "/fast"}) // TODO(gg): use flashbots mock?
		c.EVM[0].Transactions.TransactionManagerV2.DualBroadcast = ptr(true)

		// [EVM.Transactions.AutoPurge]
		// c.EVM[0].Transactions.AutoPurge.Enabled = ptr(true)
		// c.EVM[0].Transactions.AutoPurge.Threshold = ptr(uint32(5))
		// c.EVM[0].Transactions.AutoPurge.MinAttempts = ptr(uint32(100))
		// // c.EVM[0].Transactions.AutoPurge.DetectionApiUrl = ptr(commonconfig.URL{Scheme: "https", Host: "protecc.flashbots.net", Path: "/tx/"}) // TODO(gg): use flashbots mock?
		// c.EVM[0].GasEstimator.BumpThreshold = ptr(uint32(6))

		// if cfg.Transactions().TransactionManagerV2().Enabled() {
		// 		c.EVM[0].Transactions.TransactionManagerV2.Enabled = ptr(true)
		// 		c.EVM[0].Transactions.TransactionManagerV2.ForwardersEnabled = ptr(true)
		// 	}

		// [Log]
		c.Log.Level = ptr(logtoml.LogLevel(zapcore.DebugLevel))
	})

	// Create file logger core // TODO(gg): maybe refactor this into something more general?
	fileCore, logFile, err := createFileLogger(t, nodeName)
	if err != nil {
		t.Fatalf("Failed to create file logger: %v", err)
	}
	// Ensure log file is closed on cleanup
	t.Cleanup(func() {
		if logFile != nil {
			_ = logFile.Sync()
			_ = logFile.Close()
		}
	})

	// Create observed logger (for test assertions)
	lggr, observedLogs := logger.TestLoggerObserved(t, config.Log().Level())

	// Start a goroutine to write observed logs to file in real-time
	// This ensures the log file captures all logs from the node
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		lastCount := 0
		ctx := testutils.Context(t)
		for {
			select {
			case <-ticker.C:
				logs := observedLogs.All()
				if len(logs) > lastCount {
					// Write new logs to file
					for i := lastCount; i < len(logs); i++ {
						log := logs[i]
						entry := zapcore.Entry{
							Time:    log.Time,
							Level:   log.Level,
							Message: log.Message,
							Caller:  log.Caller,
						}
						_ = fileCore.Write(entry, log.Context)
					}
					lastCount = len(logs)
					// Sync file periodically
					_ = logFile.Sync()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	if backend != nil {
		app = cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, backend, p2pKey)
	} else {
		app = cltest.NewApplicationWithConfig(t, config, p2pKey, ocr2kb, csaKey, lggr.Named(nodeName))
	}

	err = app.Start(testutils.Context(t))
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	t.Logf("p2pkey raw is %s", p2pKey.PeerID().Raw())
	return app, p2pKey.PeerID().Raw(), csaKey.StaticSizedPublicKey(), ocr2kb, observedLogs
}

// createFileLogger creates a logger that writes to a file in the logs/ directory
func createFileLogger(t *testing.T, nodeName string) (zapcore.Core, *os.File, error) {
	// Find project root by looking for go.mod file
	// Start from current working directory and walk up
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Try to find project root (where go.mod is)
	logsDir := "logs"
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found project root, use logs directory there
			logsDir = filepath.Join(dir, "logs")
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root, use logs in current directory
			logsDir = filepath.Join(cwd, "logs")
			break
		}
		dir = parent
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		return nil, nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with test name and node name
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("%s_%s_%s.log", timestamp, t.Name(), nodeName)
	// Sanitize filename (remove invalid characters)
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, " ", "_")

	logPath := filepath.Join(logsDir, filename)

	file, err := os.Create(logPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create log file: %w", err)
	}

	t.Logf("Writing logs for %s to %s", nodeName, logPath)

	// Create encoder config
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000000Z07:00")
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core that writes to file
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(file),
		zapcore.DebugLevel,
	)

	return core, file, nil
}

func ptr[T any](t T) *T { return &t }
