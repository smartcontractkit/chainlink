package svr

import (
	"context"
	"encoding/hex"
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

	"github.com/ethereum/go-ethereum/ethclient/simulated"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
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
	clientPubKey         credentials.StaticSizedPublicKey
	keyBundle            ocr2key.KeyBundle
	observedLogs         *observer.ObservedLogs
	effectiveTransmitter common.Address
}

func setupNodes(t *testing.T, nNodes int, transactOpts *bind.TransactOpts, backend *simulated.Backend, clientCSAKeys []csakey.KeyV2) (oracles []confighelper.OracleIdentityExtra, nodes []node) {
	ports := freeport.GetN(t, nNodes)
	for i := range nNodes {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("core_node_%d", i), backend, clientCSAKeys[i])
		t.Logf("Node %d with transmitter %#v (%s) and peer id %s started on port %d", i, transmitter, fmt.Sprintf("%x", transmitter[:]), peerID, ports[i])

		keys, err := app.GetKeyStore().Eth().GetAll(context.Background())
		require.NoErrorf(t, err, "failed to get keys")
		if len(keys) == 0 {
			t.Logf("No keys found")
		} else {
			for _, key := range keys {
				t.Logf("Existing key address: %s", key.Address.Hex())
			}
		}

		offchainPublicKey, err := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		require.NoError(t, err)
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})

		/** TODO(gg): to use:
				oracles = append(oracles, confighelper2.OracleIdentityExtra{
					OracleIdentity: confighelper2.OracleIdentity{
						OnchainPublicKey:  node.KeyBundle.PublicKey(),
						TransmitAccount:   ocrtypes2.Account(node.EffectiveTransmitter.String()),
						OffchainPublicKey: node.KeyBundle.OffchainPublicKey(),
						PeerID:            node.PeerID,
					},
					ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
				})
		**/

		// deploy a forwarder

		forwarderAddress, _, authorizedForwarder, err := authorized_forwarder.DeployAuthorizedForwarder(transactOpts, backend.Client(), linkTokenAddress, transactOpts.From, common.Address{}, []byte{})
		require.NoError(t, err)
		backend.Commit()

		// set EOA as an authorized sender for the forwarder
		_, err = authorizedForwarder.SetAuthorizedSenders(transactOpts, []common.Address{common.HexToAddress(fmt.Sprintf("%x", transmitter[:]))})
		require.NoError(t, err)
		backend.Commit()

		// add forwarder address to be tracked in db
		forwarderORM := forwarders.NewORM(app.GetDB())
		chainID, err := backend.Client().ChainID(testutils.Context(t))
		require.NoError(t, err)
		_, err = forwarderORM.CreateForwarder(testutils.Context(t), forwarderAddress, ubig.Big(*chainID))
		require.NoError(t, err)

		// }
		// return &Node{
		// 	App:                  app,
		// 	PeerID:               p2pKey.PeerID().Raw(),
		// 	Transmitter:          transmitter,
		// 	EffectiveTransmitter: effectiveTransmitter,
		// 	KeyBundle:            kb,
		// }

		nodes = append(nodes, node{
			app:                  app,
			clientPubKey:         credentials.StaticSizedPublicKey(transmitter),
			keyBundle:            kb,
			observedLogs:         observedLogs,
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

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		/**
		  [Log]
		  Level = 'debug'

		  [Pyroscope]
		  ServerAddress = 'http://host.docker.internal:4040'
		  Environment = 'local'

		  [WebServer]
		  HTTPWriteTimeout = '30s'
		  SecureCookies = false
		  HTTPPort = {{.HTTPPort}}

		  [WebServer.TLS]
		  HTTPSPort = 0

		  [JobPipeline]
		  [JobPipeline.HTTPRequest]
		  DefaultTimeout = '30s'
		*/

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
		c.P2P.TraceLogging = ptr(true)

		// [P2P.V2]
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)

		// [Mercury]
		c.Mercury.VerboseLogging = ptr(true)
		c.Mercury.Transmitter.ReaperFrequency = commonconfig.MustNewDuration(0 * time.Millisecond)

		// [Log]
		c.Log.Level = ptr(toml.LogLevel(zapcore.DebugLevel))
	})

	// Create file logger core
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
		app = cltest.NewApplicationWithConfigV2OnSimulatedBlockchain(t, config, backend, p2pKey, ocr2kb, csaKey, lggr.Named(nodeName))
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
