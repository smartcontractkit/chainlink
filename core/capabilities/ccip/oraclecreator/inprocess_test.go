package oraclecreator_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/oraclecreator"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	ocr2validate "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestOracleCreator_CreateBootstrap(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	keyStore := keystore.New(db, utils.DefaultScryptParams, logger.NullLogger)
	require.NoError(t, keyStore.Unlock(testutils.Context(t), cltest.Password), "unable to unlock keystore")
	p2pKey, err := keyStore.P2P().Create(testutils.Context(t))
	require.NoError(t, err)
	peerID := p2pKey.PeerID()
	listenPort := freeport.GetOne(t)
	generalConfig := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.P2P.PeerID = ptr(peerID)
		c.P2P.TraceLogging = ptr(false)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.ListenAddresses = ptr([]string{fmt.Sprintf("127.0.0.1:%d", listenPort)})

		c.OCR2.Enabled = ptr(true)
	})
	peerWrapper := ocrcommon.NewSingletonPeerWrapper(keyStore, generalConfig.P2P(), generalConfig.OCR(), db, logger.NullLogger)
	require.NoError(t, peerWrapper.Start(testutils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, peerWrapper.Close()) })

	// NOTE: this is a bit of a hack to get the OCR2 job created in order to use the ocr db
	// the ocr2_contract_configs table has a foreign key constraint on ocr2_oracle_spec_id
	// which is passed into ocr2.NewDB.
	pipelineORM := pipeline.NewORM(db,
		logger.NullLogger, generalConfig.JobPipeline().MaxSuccessfulRuns())
	bridgesORM := bridges.NewORM(db)

	jobORM := job.NewORM(db, pipelineORM, bridgesORM, keyStore, logger.TestLogger(t))
	t.Cleanup(func() { assert.NoError(t, jobORM.Close()) })

	jb, err := ocr2validate.ValidatedOracleSpecToml(testutils.Context(t), generalConfig.OCR2(), generalConfig.Insecure(), testspecs.GetOCR2EVMSpecMinimal(), nil)
	require.NoError(t, err)
	const juelsPerFeeCoinSource = `
	ds          [type=http method=GET url="https://chain.link/ETH-USD"];
	ds_parse    [type=jsonparse path="data.price" separator="."];
	ds_multiply [type=multiply times=100];
	ds -> ds_parse -> ds_multiply;`

	_, address := cltest.MustInsertRandomKey(t, keyStore.Eth())
	jb.Name = null.StringFrom("Job 1")
	jb.OCR2OracleSpec.TransmitterID = null.StringFrom(address.String())
	jb.OCR2OracleSpec.PluginConfig["juelsPerFeeCoinSource"] = juelsPerFeeCoinSource

	err = jobORM.CreateJob(testutils.Context(t), &jb)
	require.NoError(t, err)

	cltest.AssertCount(t, db, "ocr2_oracle_specs", 1)
	cltest.AssertCount(t, db, "jobs", 1)

	var oracleSpecID int32
	err = db.Get(&oracleSpecID, "SELECT id FROM ocr2_oracle_specs LIMIT 1")
	require.NoError(t, err)

	ocrdb := ocr2.NewDB(db, oracleSpecID, 0, logger.NullLogger)

	oc := oraclecreator.New(
		nil,
		nil,
		nil,
		peerWrapper,
		uuid.Max,
		0,
		false,
		nil,
		ocrdb,
		logger.TestLogger(t),
		&mockEndpointGen{},
		[]commontypes.BootstrapperLocator{},
		nil,
	)

	chainSelector := chainsel.GETH_TESTNET.Selector
	oracles, offchainConfig := ocrOffchainConfig(t, keyStore)
	bootstrapP2PID, err := p2pkey.MakePeerID(oracles[0].PeerID)
	require.NoError(t, err)
	transmitters := func() [][]byte {
		var transmitters [][]byte
		for _, o := range oracles {
			transmitters = append(transmitters, hexutil.MustDecode(string(o.TransmitAccount)))
		}
		return transmitters
	}()
	configDigest := ccipConfigDigest()
	bootstrap, err := oc.CreateBootstrapOracle(cctypes.OCR3ConfigWithMeta{
		ConfigDigest: configDigest,
		ConfigCount:  1,
		Config: reader.OCR3Config{
			ChainSelector:         ccipocr3.ChainSelector(chainSelector),
			OfframpAddress:        testutils.NewAddress().Bytes(),
			PluginType:            uint8(cctypes.PluginTypeCCIPCommit),
			F:                     1,
			OffchainConfigVersion: 30,
			BootstrapP2PIds:       [][32]byte{bootstrapP2PID},
			P2PIds: func() [][32]byte {
				var ids [][32]byte
				for _, o := range oracles {
					id, err2 := p2pkey.MakePeerID(o.PeerID)
					require.NoError(t, err2)
					ids = append(ids, id)
				}
				return ids
			}(),
			Signers: func() [][]byte {
				var signers [][]byte
				for _, o := range oracles {
					signers = append(signers, o.OnchainPublicKey)
				}
				return signers
			}(),
			Transmitters:   transmitters,
			OffchainConfig: offchainConfig,
		},
	})
	require.NoError(t, err)
	require.NoError(t, bootstrap.Start())
	t.Cleanup(func() { assert.NoError(t, bootstrap.Close()) })

	tests.AssertEventually(t, func() bool {
		c, err := ocrdb.ReadConfig(testutils.Context(t))
		require.NoError(t, err)
		return c.ConfigDigest == configDigest
	})
}

func ccipConfigDigest() [32]byte {
	rand32Bytes := testutils.Random32Byte()
	// overwrite first four bytes to be 0x000a, to match the prefix in libocr.
	rand32Bytes[0] = 0x00
	rand32Bytes[1] = 0x0a
	return rand32Bytes
}

type mockEndpointGen struct{}

func (m *mockEndpointGen) GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) commontypes.MonitoringEndpoint {
	return &telemetry.NoopAgent{}
}

func ptr[T any](b T) *T {
	return &b
}

func ocrOffchainConfig(t *testing.T, ks keystore.Master) (oracles []confighelper2.OracleIdentityExtra, offchainConfig []byte) {
	for i := 0; i < 4; i++ {
		kb, err := ks.OCR2().Create(testutils.Context(t), chaintype.EVM)
		require.NoError(t, err)
		p2pKey, err := ks.P2P().Create(testutils.Context(t))
		require.NoError(t, err)
		ethKey, err := ks.Eth().Create(testutils.Context(t))
		require.NoError(t, err)
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OffchainPublicKey: kb.OffchainPublicKey(),
				OnchainPublicKey:  types.OnchainPublicKey(kb.OnChainPublicKey()),
				PeerID:            p2pKey.ID(),
				TransmitAccount:   types.Account(ethKey.Address.Hex()),
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}
	var schedule []int
	for range oracles {
		schedule = append(schedule, 1)
	}
	offchainConfig, onchainConfig := []byte{}, []byte{}
	f := uint8(1)

	_, _, _, _, _, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		30*time.Second, // deltaProgress
		10*time.Second, // deltaResend
		20*time.Second, // deltaInitial
		2*time.Second,  // deltaRound
		20*time.Second, // deltaGrace
		10*time.Second, // deltaCertifiedCommitRequest
		10*time.Second, // deltaStage
		3,              // rmax
		schedule,
		oracles,
		offchainConfig,
		50*time.Millisecond, // maxDurationQuery
		5*time.Second,       // maxDurationObservation
		10*time.Second,      // maxDurationShouldAcceptAttestedReport
		10*time.Second,      // maxDurationShouldTransmitAcceptedReport
		int(f),
		onchainConfig)
	require.NoError(t, err, "failed to create contract config")

	return oracles, offchainConfig
}
