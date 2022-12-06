package job_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"

	evmconfigmocks "github.com/smartcontractkit/chainlink/core/chains/evm/config/mocks"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	configtest2 "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	ocr2mocks "github.com/smartcontractkit/chainlink/core/services/ocr2/mocks"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	pkgconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	clhttptest "github.com/smartcontractkit/chainlink/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

var monitoringEndpoint = telemetry.MonitoringEndpointGenerator(&telemetry.NoopAgent{})

func TestRunner(t *testing.T) {
	config := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	keyStore := cltest.NewKeyStore(t, db, config)
	ethKeyStore := keyStore.Eth()

	ethClient := cltest.NewEthMocksWithDefaultChain(t)
	ethClient.On("HeadByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(10), nil)
	ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, nil)

	pipelineORM := pipeline.NewORM(db, logger.TestLogger(t), config)
	btORM := bridges.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: ethClient, GeneralConfig: config})
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	runner := pipeline.NewRunner(pipelineORM, btORM, config, cc, nil, nil, logger.TestLogger(t), c, c)
	jobORM := NewTestORM(t, db, cc, pipelineORM, btORM, keyStore, config)

	require.NoError(t, runner.Start(testutils.Context(t)))
	t.Cleanup(func() { assert.NoError(t, runner.Close()) })

	_, transmitterAddress := cltest.MustInsertRandomKey(t, ethKeyStore, 0)
	require.NoError(t, keyStore.OCR().Add(cltest.DefaultOCRKey))

	t.Run("gets the election result winner", func(t *testing.T) {
		var httpURL string
		mockElectionWinner := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `Hal Finney`,
			func(header http.Header, s string) {
				var md bridges.BridgeMetaDataJSON
				require.NoError(t, json.Unmarshal([]byte(s), &md))
				assert.Equal(t, big.NewInt(10), md.Meta.LatestAnswer)
				assert.Equal(t, big.NewInt(100), md.Meta.UpdatedAt)
			})
		mockVoterTurnout := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"data": {"result": 62.57}}`,
			func(header http.Header, s string) {
				var md bridges.BridgeMetaDataJSON
				require.NoError(t, json.Unmarshal([]byte(s), &md))
				assert.Equal(t, big.NewInt(10), md.Meta.LatestAnswer)
				assert.Equal(t, big.NewInt(100), md.Meta.UpdatedAt)
			},
		)
		mockHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"turnout": 61.942}`)

		httpURL = mockHTTP.URL
		_, bridgeER := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: mockElectionWinner.URL}, config)
		_, bridgeVT := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: mockVoterTurnout.URL}, config)

		// Need a job in order to create a run
		jb := MakeVoterTurnoutOCRJobSpecWithHTTPURL(t, transmitterAddress, httpURL, bridgeVT.Name.String(), bridgeER.Name.String())
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)
		require.NotNil(t, jb.PipelineSpec)

		m, err := bridges.MarshalBridgeMetaData(big.NewInt(10), big.NewInt(100))
		require.NoError(t, err)
		runID, results, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(map[string]interface{}{"jobRun": map[string]interface{}{"meta": m}}), logger.TestLogger(t), true)
		require.NoError(t, err)

		require.Len(t, results.Values, 2)
		require.GreaterOrEqual(t, len(results.FatalErrors), 2)
		assert.Nil(t, results.FatalErrors[0])
		assert.Nil(t, results.FatalErrors[1])
		require.GreaterOrEqual(t, len(results.AllErrors), 2)
		assert.Equal(t, "6225.6", results.Values[0].(decimal.Decimal).String())
		assert.Equal(t, "Hal Finney", results.Values[1].(string))

		// Verify individual task results
		var runs []pipeline.TaskRun
		sql := `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`
		err = db.Select(&runs, sql, runID)
		assert.NoError(t, err)
		assert.Len(t, runs, 8)

		for _, run := range runs {
			if run.GetDotID() == "answer2" {
				assert.Equal(t, "Hal Finney", run.Output.Val)
			} else if run.GetDotID() == "ds2" {
				assert.Equal(t, `{"turnout": 61.942}`, run.Output.Val)
			} else if run.GetDotID() == "ds2_parse" {
				assert.Equal(t, float64(61.942), run.Output.Val)
			} else if run.GetDotID() == "ds2_multiply" {
				assert.Equal(t, "6194.2", run.Output.Val)
			} else if run.GetDotID() == "ds1" {
				assert.Equal(t, `{"data": {"result": 62.57}}`, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.Equal(t, float64(62.57), run.Output.Val)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Equal(t, "6257", run.Output.Val)
			} else if run.GetDotID() == "answer1" {
				assert.Equal(t, "6225.6", run.Output.Val)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("must delete job before deleting bridge", func(t *testing.T) {
		_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)
		jb := makeOCRJobSpecFromToml(t, fmt.Sprintf(`
			type               = "offchainreporting"
			schemaVersion      = 1
			observationSource = """
				ds1          [type=bridge name="%s"];
			"""
		`, bridge.Name.String()))
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)
		// Should not be able to delete a bridge in use.
		jids, err := jobORM.FindJobIDsWithBridge(bridge.Name.String())
		require.NoError(t, err)
		require.Equal(t, 1, len(jids))

		// But if we delete the job, then we can.
		require.NoError(t, jobORM.DeleteJob(jb.ID))
		jids, err = jobORM.FindJobIDsWithBridge(bridge.Name.String())
		require.NoError(t, err)
		require.Equal(t, 0, len(jids))
	})

	t.Run("referencing a non-existent bridge should error", func(t *testing.T) {
		// Create a random bridge name
		_, b := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{}, config)

		// Reference a different one
		cfg := new(evmconfigmocks.ChainScopedConfig)
		cfg.On("Dev").Return(true)
		cfg.On("ChainType").Return(pkgconfig.ChainType(""))
		c := new(evmmocks.Chain)
		c.On("Config").Return(cfg)
		cs := new(evmmocks.ChainSet)
		cs.On("Get", mock.Anything).Return(c, nil)

		jb, err := ocr.ValidatedOracleSpecToml(cs, `
			type               = "offchainreporting"
			schemaVersion      = 1
			evmChainID         = 1
			contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
			isBootstrapPeer    = false
			blockchainTimeout  = "1s"
			observationTimeout = "10s"
			databaseTimeout    = "2s"
			contractConfigTrackerPollInterval="1s"
			contractConfigConfirmations=1
			observationGracePeriod = "2s"
			contractTransmitterTransmitTimeout = "500ms"
			contractConfigTrackerSubscribeInterval="1s"
			observationSource = """
			ds1          [type=bridge name=blah];
			ds1_parse    [type=jsonparse path="one,two"];
			ds1_multiply [type=multiply times=1.23];
			ds1 -> ds1_parse -> ds1_multiply -> answer1;
			answer1      [type=median index=0];
			"""
		`)
		require.NoError(t, err)
		// Should error creating it
		err = jobORM.CreateJob(&jb)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not all bridges exist")

		// Same for ocr2
		cfg2 := new(ocr2mocks.Config)
		cfg2.On("OCR2ContractTransmitterTransmitTimeout").Return(time.Second)
		cfg2.On("OCR2DatabaseTimeout").Return(time.Second)
		cfg2.On("Dev").Return(true)
		jb2, err := validate.ValidatedOracleSpecToml(cfg2, fmt.Sprintf(`
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
blockchainTimeout = "1s"
contractConfigTrackerPollInterval = "2s"
contractConfigConfirmations = 1
observationSource  = """
ds1          [type=bridge name="%s"];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name=blah];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`, b.Name.String()))
		require.NoError(t, err)
		// Should error creating it because of the juels per fee coin non-existent bridge
		err = jobORM.CreateJob(&jb2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not all bridges exist")

		// Duplicate bridge names that exist is ok
		cfg2.On("OCR2ContractTransmitterTransmitTimeout").Return(time.Second)
		cfg2.On("OCR2DatabaseTimeout").Return(time.Second)
		cfg2.On("Dev").Return(true)
		jb3, err := validate.ValidatedOracleSpecToml(cfg2, fmt.Sprintf(`
type               = "offchainreporting2"
pluginType         = "median"
schemaVersion      = 1
relay              = "evm"
contractID         = "0x613a38AC1659769640aaE063C651F48E0250454C"
blockchainTimeout = "1s"
contractConfigTrackerPollInterval = "2s"
contractConfigConfirmations = 1
observationSource  = """
ds1          [type=bridge name="%s"];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
ds1          [type=bridge name="%s"];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
ds2          [type=bridge name="%s"];
ds2_parse    [type=jsonparse path="one,two"];
ds2_multiply [type=multiply times=1.23];
ds2 -> ds2_parse -> ds2_multiply -> answer1;
answer1      [type=median index=0];
"""
`, b.Name.String(), b.Name.String(), b.Name.String()))
		require.NoError(t, err)
		// Should not error with duplicate bridges
		err = jobORM.CreateJob(&jb3)
		require.NoError(t, err)
	})

	t.Run("handles the case where the parsed value is literally null", func(t *testing.T) {
		var httpURL string
		resp := `{"USD": null}`
		{
			mockHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		jb := makeSimpleFetchOCRJobSpecWithHTTPURL(t, transmitterAddress, httpURL, false)
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)

		runID, results, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)

		assert.Len(t, results.FatalErrors, 1)
		assert.Len(t, results.Values, 1)
		assert.Contains(t, results.FatalErrors[0].Error(), "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results.Values[0])

		// Verify individual task results
		var runs []pipeline.TaskRun
		sql := `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`
		err = db.Select(&runs, sql, runID)
		assert.NoError(t, err)
		require.Len(t, runs, 3)

		for _, run := range runs {
			if run.GetDotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				require.NotNil(t, resp, run.Output)
				assert.Equal(t, resp, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				assert.False(t, run.Output.Valid)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Contains(t, run.Error.ValueOrZero(), "type <nil> cannot be converted to decimal.Decimal")
				assert.False(t, run.Output.Valid)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("handles the case where the jsonparse lookup path is missing from the http response", func(t *testing.T) {
		var httpURL string
		resp := "{\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}"
		{
			mockHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		jb := makeSimpleFetchOCRJobSpecWithHTTPURL(t, transmitterAddress, httpURL, false)
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)

		runID, results, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)

		assert.Len(t, results.Values, 1)
		assert.Len(t, results.FatalErrors, 1)
		assert.Contains(t, results.FatalErrors[0].Error(), pipeline.ErrTooManyErrors.Error())
		assert.Nil(t, results.Values[0])

		// Verify individual task results
		var runs []pipeline.TaskRun
		sql := `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`
		err = db.Select(&runs, sql, runID)
		assert.NoError(t, err)
		require.Len(t, runs, 3)

		for _, run := range runs {
			if run.GetDotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.Contains(t, run.Error.ValueOrZero(), "could not resolve path [\"USD\"] in {\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}")
				assert.False(t, run.Output.Valid)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Contains(t, run.Error.ValueOrZero(), pipeline.ErrTooManyErrors.Error())
				assert.False(t, run.Output.Valid)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("handles the case where the jsonparse lookup path is missing from the http response and lax is enabled", func(t *testing.T) {
		var httpURL string
		resp := "{\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}"
		{
			mockHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		jb := makeSimpleFetchOCRJobSpecWithHTTPURL(t, transmitterAddress, httpURL, true)
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)

		runID, results, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)

		assert.Len(t, results.Values, 1)
		assert.Contains(t, results.FatalErrors[0].Error(), "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results.Values[0])

		// Verify individual task results
		var runs []pipeline.TaskRun
		sql := `SELECT * FROM pipeline_task_runs WHERE pipeline_run_id = $1`
		err = db.Select(&runs, sql, runID)
		assert.NoError(t, err)
		require.Len(t, runs, 3)

		for _, run := range runs {
			if run.GetDotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				assert.False(t, run.Output.Valid)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Contains(t, run.Error.ValueOrZero(), "type <nil> cannot be converted to decimal.Decimal")
				assert.False(t, run.Output.Valid)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("missing required env vars", func(t *testing.T) {
		s := `
		type               = "offchainreporting"
		schemaVersion      = 1
		contractAddress    = "%s"
		isBootstrapPeer    = false
		observationSource = """
ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true" %s];
ds1_parse    [type=jsonparse path="USD" lax=true];
ds1 -> ds1_parse;
"""
`
		s = fmt.Sprintf(s, cltest.NewEIP55Address(), "http://blah.com", "")
		jb, err := ocr.ValidatedOracleSpecToml(cc, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &jb)
		require.NoError(t, err)
		jb.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)
		// Required to create job spawner delegate.
		config.Overrides.P2PListenPort = null.IntFrom(2000)
		sd := ocr.NewDelegate(
			db,
			jobORM,
			keyStore,
			nil,
			nil,
			nil,
			cc,
			logger.TestLogger(t),
			config,
			srvctest.Start(t, utils.NewMailboxMonitor(t.Name())),
		)
		_, err = sd.ServicesForSpec(jb)
		// We expect this to fail as neither the required vars are not set either via the env nor the job itself.
		require.Error(t, err)
	})

	t.Run("use env for minimal bootstrap", func(t *testing.T) {
		s := `
		type               = "offchainreporting"
		schemaVersion      = 1
		contractAddress    = "%s"
		isBootstrapPeer    = true
`
		s = fmt.Sprintf(s, cltest.NewEIP55Address())
		jb, err := ocr.ValidatedOracleSpecToml(cc, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &jb)
		require.NoError(t, err)
		jb.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)
		// Required to create job spawner delegate.
		config.Overrides.P2PListenPort = null.IntFrom(2000)

		lggr := logger.TestLogger(t)
		_, err = keyStore.P2P().Create()
		assert.NoError(t, err)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, config, db, lggr)
		require.NoError(t, pw.Start(testutils.Context(t)))
		sd := ocr.NewDelegate(
			db,
			jobORM,
			keyStore,
			nil,
			pw,
			monitoringEndpoint,
			cc,
			lggr,
			config,
			srvctest.Start(t, utils.NewMailboxMonitor(t.Name())),
		)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("use env for minimal non-bootstrap", func(t *testing.T) {
		kb, err := keyStore.OCR().Create()
		require.NoError(t, err)
		s := `
		type               = "offchainreporting"
		schemaVersion      = 1
		contractAddress    = "%s"
		isBootstrapPeer    = false
		observationTimeout = "15s"
		observationSource = """
ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true" %s];
ds1_parse    [type=jsonparse path="USD" lax=true];
ds1 -> ds1_parse;
"""
`
		s = fmt.Sprintf(s, cltest.NewEIP55Address(), "http://blah.com", "")
		tAddress := ethkey.EIP55AddressFromAddress(transmitterAddress)
		config.Overrides.P2PBootstrapPeers = []string{"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju", "/dns4/chain.link/tcp/1235/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"}
		config.Overrides.P2PV2Bootstrappers = []commontypes.BootstrapperLocator{}
		config.Overrides.OCRKeyBundleID = null.NewString(kb.ID(), true)
		config.Overrides.OCRTransmitterAddress = &tAddress
		jb, err := ocr.ValidatedOracleSpecToml(cc, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &jb)
		require.NoError(t, err)
		jb.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)
		// Assert the override
		assert.Equal(t, jb.OCROracleSpec.ObservationTimeout, models.Interval(cltest.MustParseDuration(t, "15s")))
		// Assert that this is default
		assert.Equal(t, models.Interval(20000000000), jb.OCROracleSpec.BlockchainTimeout)
		assert.Equal(t, models.Interval(cltest.MustParseDuration(t, "1s")), jb.MaxTaskDuration)

		// Required to create job spawner delegate.
		config.Overrides.P2PListenPort = null.IntFrom(2000)
		lggr := logger.TestLogger(t)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, config, db, lggr)
		require.NoError(t, pw.Start(testutils.Context(t)))
		sd := ocr.NewDelegate(
			db,
			jobORM,
			keyStore,
			nil,
			pw,
			monitoringEndpoint,
			cc,
			lggr,
			config,
			srvctest.Start(t, utils.NewMailboxMonitor(t.Name())),
		)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("test min non-bootstrap", func(t *testing.T) {
		kb, err := keyStore.OCR().Create()
		require.NoError(t, err)

		s := fmt.Sprintf(minimalNonBootstrapTemplate, cltest.NewEIP55Address(), transmitterAddress.Hex(), kb.ID(), "http://blah.com", "")
		jb, err := ocr.ValidatedOracleSpecToml(cc, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &jb)
		require.NoError(t, err)

		jb.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)
		assert.Equal(t, jb.MaxTaskDuration, models.Interval(cltest.MustParseDuration(t, "1s")))

		// Required to create job spawner delegate.
		config.Overrides.P2PListenPort = null.IntFrom(2000)
		lggr := logger.TestLogger(t)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, config, db, lggr)
		require.NoError(t, pw.Start(testutils.Context(t)))
		sd := ocr.NewDelegate(
			db,
			jobORM,
			keyStore,
			nil,
			pw,
			monitoringEndpoint,
			cc,
			lggr,
			config,
			srvctest.Start(t, utils.NewMailboxMonitor(t.Name())),
		)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("test min bootstrap", func(t *testing.T) {
		s := fmt.Sprintf(minimalBootstrapTemplate, cltest.NewEIP55Address())
		jb, err := ocr.ValidatedOracleSpecToml(cc, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &jb)
		require.NoError(t, err)
		err = jobORM.CreateJob(&jb)
		require.NoError(t, err)

		// Required to create job spawner delegate.
		config.Overrides.P2PListenPort = null.IntFrom(2000)
		lggr := logger.TestLogger(t)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, config, db, lggr)
		require.NoError(t, pw.Start(testutils.Context(t)))
		sd := ocr.NewDelegate(
			db,
			jobORM,
			keyStore,
			nil,
			pw,
			monitoringEndpoint,
			cc,
			lggr,
			config,
			srvctest.Start(t, utils.NewMailboxMonitor(t.Name())),
		)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("test job spec error is created", func(t *testing.T) {
		// Create a keystore with an ocr key bundle and p2p key.
		kb, err := keyStore.OCR().Create()
		require.NoError(t, err)
		spec := fmt.Sprintf(ocrJobSpecTemplate, testutils.NewAddress().Hex(), kb.ID(), transmitterAddress.Hex(), fmt.Sprintf(simpleFetchDataSourceTemplate, "blah", true))
		jb := makeOCRJobSpecFromToml(t, spec)

		// Create an OCR job
		err = jobORM.CreateJob(jb)
		require.NoError(t, err)

		// Required to create job spawner delegate.
		config.Overrides.P2PListenPort = null.IntFrom(2000)
		config.Overrides.P2PV2Bootstrappers = []commontypes.BootstrapperLocator{}
		lggr := logger.TestLogger(t)
		pw := ocrcommon.NewSingletonPeerWrapper(keyStore, config, db, lggr)
		require.NoError(t, pw.Start(testutils.Context(t)))

		sd := ocr.NewDelegate(
			db,
			jobORM,
			keyStore,
			nil,
			pw,
			monitoringEndpoint,
			cc,
			lggr,
			config,
			srvctest.Start(t, utils.NewMailboxMonitor(t.Name())),
		)
		services, err := sd.ServicesForSpec(*jb)
		require.NoError(t, err)

		// Return an error getting the contract code.
		ethClient.On("CodeAt", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("no such code"))
		ctx := testutils.Context(t)
		for _, s := range services {
			err = s.Start(ctx)
			require.NoError(t, err)
		}
		var se []job.SpecError
		require.Eventually(t, func() bool {
			err = db.Select(&se, `SELECT * FROM job_spec_errors`)
			require.NoError(t, err)
			return len(se) == 1
		}, time.Second, 100*time.Millisecond)
		require.Len(t, se, 1)
		assert.Equal(t, uint(1), se[0].Occurrences)

		for _, s := range services {
			err = s.Close()
			require.NoError(t, err)
		}

		// Ensure we can delete an errored
		err = jobORM.DeleteJob(jb.ID)
		require.NoError(t, err)
		se = []job.SpecError{}
		err = db.Select(&se, `SELECT * FROM job_spec_errors`)
		require.NoError(t, err)
		require.Len(t, se, 0)

		// TODO: This breaks the txdb connection, failing subsequent tests. Resolve in the future
		// Noop once the job is gone.
		// jobORM.RecordError(testutils.Context(t), jb.ID, "test")
		// err = db.Find(&se).Error
		// require.NoError(t, err)
		// require.Len(t, se, 0)
	})

	t.Run("timeouts", func(t *testing.T) {
		// There are 4 timeouts:
		// - ObservationTimeout = how long the whole OCR time needs to run, or it fails (default 10 seconds)
		// - config.JobPipelineMaxTaskDuration() = node level maximum time for a pipeline task (default 10 minutes)
		// - config.transmitterAddress, http specific timeouts (default 15s * 5 retries = 75s)
		// - "d1 [.... timeout="2s"]" = per task level timeout (should override the global config)
		serv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(1 * time.Millisecond)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"USD":10.1}`))
		}))
		defer serv.Close()

		jb := makeMinimalHTTPOracleSpec(t, db, config, cltest.NewEIP55Address().String(), transmitterAddress.Hex(), cltest.DefaultOCRKeyBundleID, serv.URL, `timeout="1ns"`)
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)

		_, results, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)
		assert.Nil(t, results.Values[0])

		// No task timeout should succeed.
		jb = makeMinimalHTTPOracleSpec(t, db, config, cltest.NewEIP55Address().String(), transmitterAddress.Hex(), cltest.DefaultOCRKeyBundleID, serv.URL, "")
		jb.Name = null.NewString("a job 2", true)
		err = jobORM.CreateJob(jb)
		require.NoError(t, err)
		_, results, err = runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)
		assert.Equal(t, 10.1, results.Values[0])
		assert.Nil(t, results.FatalErrors[0])

		// Job specified task timeout should fail.
		jb = makeMinimalHTTPOracleSpec(t, db, config, cltest.NewEIP55Address().String(), transmitterAddress.Hex(), cltest.DefaultOCRKeyBundleID, serv.URL, "")
		jb.MaxTaskDuration = models.Interval(time.Duration(1))
		jb.Name = null.NewString("a job 3", true)
		err = jobORM.CreateJob(jb)
		require.NoError(t, err)

		_, results, err = runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)
		assert.NotNil(t, results.FatalErrors[0])
	})

	t.Run("deleting jobs", func(t *testing.T) {
		var httpURL string
		{
			resp := `{"USD": 42.42}`
			mockHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		jb := makeSimpleFetchOCRJobSpecWithHTTPURL(t, transmitterAddress, httpURL, false)
		err := jobORM.CreateJob(jb)
		require.NoError(t, err)

		_, results, err := runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.NoError(t, err)
		assert.Len(t, results.Values, 1)
		assert.Nil(t, results.FatalErrors[0])
		assert.Equal(t, "4242", results.Values[0].(decimal.Decimal).String())

		// Delete the job
		err = jobORM.DeleteJob(jb.ID)
		require.NoError(t, err)

		// Create another run, it should fail
		_, _, err = runner.ExecuteAndInsertFinishedRun(testutils.Context(t), *jb.PipelineSpec, pipeline.NewVarsFrom(nil), logger.TestLogger(t), true)
		require.Error(t, err)
	})
}

func TestRunner_Success_Callback_AsyncJob(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)

	cfg := configtest2.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		t := true
		c.JobPipeline.ExternalInitiatorsEnabled = &t
		c.Database.Listener.FallbackPollInterval = models.MustNewDuration(10 * time.Millisecond)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient, cltest.UseRealExternalInitiatorManager)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: app.GetSqlxDB(), Client: ethClient, GeneralConfig: cfg})
	require.NoError(t, app.Start(testutils.Context(t)))

	var (
		eiName    = "substrate-ei"
		eiSpec    = map[string]interface{}{"foo": "bar"}
		eiRequest = map[string]interface{}{"result": 42}

		jobUUID = uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46")

		expectedCreateJobRequest = map[string]interface{}{
			"jobId":  jobUUID.String(),
			"type":   eiName,
			"params": eiSpec,
		}
	)

	// Setup EI
	var eiURL string
	var eiNotifiedOfCreate bool
	var eiNotifiedOfDelete bool
	{
		mockEI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !eiNotifiedOfCreate {
				require.Equal(t, http.MethodPost, r.Method)

				eiNotifiedOfCreate = true
				defer r.Body.Close()

				var gotCreateJobRequest map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&gotCreateJobRequest)
				require.NoError(t, err)

				require.Equal(t, expectedCreateJobRequest, gotCreateJobRequest)
				w.WriteHeader(http.StatusOK)
			} else {
				require.Equal(t, http.MethodDelete, r.Method)

				eiNotifiedOfDelete = true
				defer r.Body.Close()

				require.Equal(t, fmt.Sprintf("/%v", jobUUID.String()), r.URL.Path)
			}
		}))
		defer mockEI.Close()
		eiURL = mockEI.URL
	}

	// Create the EI record on the Core node
	var eia *auth.Token
	{
		eiCreate := map[string]string{
			"name": eiName,
			"url":  eiURL,
		}
		eiCreateJSON, err := json.Marshal(eiCreate)
		require.NoError(t, err)
		eip := cltest.CreateExternalInitiatorViaWeb(t, app, string(eiCreateJSON))
		eia = &auth.Token{
			AccessKey: eip.AccessKey,
			Secret:    eip.Secret,
		}
	}

	var responseURL string

	// Create the bridge on the Core node
	bridgeCalled := make(chan struct{}, 1)
	var bridgeName string
	{
		bridgeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			var bridgeRequest map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&bridgeRequest)
			require.NoError(t, err)

			require.Equal(t, float64(42), bridgeRequest["value"])

			responseURL = bridgeRequest["responseURL"].(string)

			w.WriteHeader(http.StatusOK)
			require.NoError(t, err)
			io.WriteString(w, `{"pending": true}`)
			bridgeCalled <- struct{}{}
		}))
		_, bridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{URL: bridgeServer.URL}, app.GetConfig())
		bridgeName = bridge.Name.String()
		defer bridgeServer.Close()
	}

	// Create the job spec on the Core node
	var jobID int32
	{
		tomlSpec := fmt.Sprintf(`
			type            = "webhook"
			schemaVersion   = 1
			externalJobID           = "%v"
			externalInitiators = [
				{
					name = "%s",
					spec = """
				%s
			"""
				}
			]
			observationSource   = """
				parse  [type=jsonparse path="result" data="$(jobRun.requestBody)"]
				ds1 [type=bridge async=true name="%s" timeout=0 requestData=<{"value": $(parse)}>]
				ds1_parse [type=jsonparse lax=false  path="data,result"]
				ds1_multiply [type=multiply times=1000000000000000000 index=0]
			
				parse->ds1->ds1_parse->ds1_multiply;
			"""
    `, jobUUID, eiName, cltest.MustJSONMarshal(t, eiSpec), bridgeName)

		_, err := webhook.ValidatedWebhookSpec(tomlSpec, app.GetExternalInitiatorManager())
		require.NoError(t, err)
		job := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: tomlSpec})))
		jobID = job.ID

		require.Eventually(t, func() bool { return eiNotifiedOfCreate }, 5*time.Second, 10*time.Millisecond, "expected external initiator to be notified of new job")
	}
	t.Run("simulate request from EI -> Core node with successful callback", func(t *testing.T) {
		cltest.AwaitJobActive(t, app.JobSpawner(), jobID, 3*time.Second)

		_ = cltest.CreateJobRunViaExternalInitiatorV2(t, app, jobUUID, *eia, cltest.MustJSONMarshal(t, eiRequest))

		pipelineORM := pipeline.NewORM(app.GetSqlxDB(), logger.TestLogger(t), cfg)
		bridgesORM := bridges.NewORM(app.GetSqlxDB(), logger.TestLogger(t), cfg)
		jobORM := NewTestORM(t, app.GetSqlxDB(), cc, pipelineORM, bridgesORM, app.KeyStore, cfg)

		// Trigger v2/resume
		select {
		case <-bridgeCalled:
		case <-time.After(time.Second):
			t.Fatal("expected bridge server to be called")
		}
		// Make the request
		{
			url, err := url.Parse(responseURL)
			require.NoError(t, err)
			client := app.NewHTTPClient(cltest.APIEmailAdmin)
			body := strings.NewReader(`{"value": {"data":{"result":"123.45"}}}`)
			response, cleanup := client.Patch(url.Path, body)
			defer cleanup()
			cltest.AssertServerResponse(t, response, http.StatusOK)
		}

		runs := cltest.WaitForPipelineComplete(t, 0, jobID, 1, 4, jobORM, 5*time.Second, 300*time.Millisecond)
		require.Len(t, runs, 1)
		run := runs[0]
		require.Len(t, run.PipelineTaskRuns, 4)
		require.Empty(t, run.PipelineTaskRuns[0].Error)
		require.Empty(t, run.PipelineTaskRuns[1].Error)
		require.Empty(t, run.PipelineTaskRuns[2].Error)
		require.Empty(t, run.PipelineTaskRuns[3].Error)
		require.Equal(t, pipeline.JSONSerializable{Val: []interface{}{"123450000000000000000"}, Valid: true}, run.Outputs)
		require.Equal(t, pipeline.RunErrors{null.String{NullString: sql.NullString{String: "", Valid: false}}}, run.FatalErrors)
	})
	// Delete the job
	{
		cltest.DeleteJobViaWeb(t, app, jobID)
		require.Eventually(t, func() bool { return eiNotifiedOfDelete }, 5*time.Second, 10*time.Millisecond, "expected external initiator to be notified of deleted job")
	}
}

func TestRunner_Error_Callback_AsyncJob(t *testing.T) {
	t.Parallel()

	ethClient := cltest.NewEthMocksWithStartupAssertions(t)

	cfg := configtest2.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		t := true
		c.JobPipeline.ExternalInitiatorsEnabled = &t
		c.Database.Listener.FallbackPollInterval = models.MustNewDuration(10 * time.Millisecond)
	})

	app := cltest.NewApplicationWithConfig(t, cfg, ethClient, cltest.UseRealExternalInitiatorManager)

	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: app.GetSqlxDB(), Client: ethClient, GeneralConfig: cfg})
	require.NoError(t, app.Start(testutils.Context(t)))

	var (
		eiName    = "substrate-ei"
		eiSpec    = map[string]interface{}{"foo": "bar"}
		eiRequest = map[string]interface{}{"result": 42}

		jobUUID = uuid.FromStringOrNil("0EEC7E1D-D0D2-476C-A1A8-72DFB6633F47")

		expectedCreateJobRequest = map[string]interface{}{
			"jobId":  jobUUID.String(),
			"type":   eiName,
			"params": eiSpec,
		}
	)

	// Setup EI
	var eiURL string
	var eiNotifiedOfCreate bool
	var eiNotifiedOfDelete bool
	{
		mockEI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !eiNotifiedOfCreate {
				require.Equal(t, http.MethodPost, r.Method)

				eiNotifiedOfCreate = true
				defer r.Body.Close()

				var gotCreateJobRequest map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&gotCreateJobRequest)
				require.NoError(t, err)

				require.Equal(t, expectedCreateJobRequest, gotCreateJobRequest)
				w.WriteHeader(http.StatusOK)
			} else {
				require.Equal(t, http.MethodDelete, r.Method)

				eiNotifiedOfDelete = true
				defer r.Body.Close()

				require.Equal(t, fmt.Sprintf("/%v", jobUUID.String()), r.URL.Path)
			}
		}))
		defer mockEI.Close()
		eiURL = mockEI.URL
	}

	// Create the EI record on the Core node
	var eia *auth.Token
	{
		eiCreate := map[string]string{
			"name": eiName,
			"url":  eiURL,
		}
		eiCreateJSON, err := json.Marshal(eiCreate)
		require.NoError(t, err)
		eip := cltest.CreateExternalInitiatorViaWeb(t, app, string(eiCreateJSON))
		eia = &auth.Token{
			AccessKey: eip.AccessKey,
			Secret:    eip.Secret,
		}
	}

	var responseURL string

	// Create the bridge on the Core node
	bridgeCalled := make(chan struct{}, 1)
	var bridgeName string
	{
		bridgeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()

			var bridgeRequest map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&bridgeRequest)
			require.NoError(t, err)

			require.Equal(t, float64(42), bridgeRequest["value"])

			responseURL = bridgeRequest["responseURL"].(string)

			w.WriteHeader(http.StatusOK)
			require.NoError(t, err)
			io.WriteString(w, `{"pending": true}`)
			bridgeCalled <- struct{}{}
		}))
		_, bridge := cltest.MustCreateBridge(t, app.GetSqlxDB(), cltest.BridgeOpts{URL: bridgeServer.URL}, app.GetConfig())
		bridgeName = bridge.Name.String()
		defer bridgeServer.Close()
	}

	// Create the job spec on the Core node
	var jobID int32
	{
		tomlSpec := fmt.Sprintf(`
			type            = "webhook"
			schemaVersion   = 1
			externalJobID           = "%v"
			externalInitiators = [
				{
					name = "%s",
					spec = """
				%s
			"""
				}
			]
			observationSource   = """
				parse  [type=jsonparse path="result" data="$(jobRun.requestBody)"]
				ds1 [type=bridge async=true name="%s" timeout=0 requestData=<{"value": $(parse)}>]
				ds1_parse [type=jsonparse lax=false  path="data,result"]
				ds1_multiply [type=multiply times=1000000000000000000 index=0]
			
				parse->ds1->ds1_parse->ds1_multiply;
			"""
    `, jobUUID, eiName, cltest.MustJSONMarshal(t, eiSpec), bridgeName)

		_, err := webhook.ValidatedWebhookSpec(tomlSpec, app.GetExternalInitiatorManager())
		require.NoError(t, err)
		job := cltest.CreateJobViaWeb(t, app, []byte(cltest.MustJSONMarshal(t, web.CreateJobRequest{TOML: tomlSpec})))
		jobID = job.ID

		require.Eventually(t, func() bool { return eiNotifiedOfCreate }, 5*time.Second, 10*time.Millisecond, "expected external initiator to be notified of new job")
	}
	t.Run("simulate request from EI -> Core node with erroring callback", func(t *testing.T) {
		_ = cltest.CreateJobRunViaExternalInitiatorV2(t, app, jobUUID, *eia, cltest.MustJSONMarshal(t, eiRequest))

		pipelineORM := pipeline.NewORM(app.GetSqlxDB(), logger.TestLogger(t), cfg)
		bridgesORM := bridges.NewORM(app.GetSqlxDB(), logger.TestLogger(t), cfg)
		jobORM := NewTestORM(t, app.GetSqlxDB(), cc, pipelineORM, bridgesORM, app.KeyStore, cfg)

		// Trigger v2/resume
		select {
		case <-bridgeCalled:
		case <-time.After(time.Second):
			t.Fatal("expected bridge server to be called")
		}
		// Make the request
		{
			url, err := url.Parse(responseURL)
			require.NoError(t, err)
			client := app.NewHTTPClient(cltest.APIEmailAdmin)
			body := strings.NewReader(`{"error": "something exploded in EA"}`)
			response, cleanup := client.Patch(url.Path, body)
			defer cleanup()
			cltest.AssertServerResponse(t, response, http.StatusOK)
		}

		runs := cltest.WaitForPipelineError(t, 0, jobID, 1, 4, jobORM, 5*time.Second, 300*time.Millisecond)
		require.Len(t, runs, 1)
		run := runs[0]
		require.Len(t, run.PipelineTaskRuns, 4)
		require.Empty(t, run.PipelineTaskRuns[0].Error)
		assert.True(t, run.PipelineTaskRuns[1].Error.Valid)
		assert.Equal(t, "something exploded in EA", run.PipelineTaskRuns[1].Error.String)
		assert.True(t, run.PipelineTaskRuns[2].Error.Valid)
		assert.True(t, run.PipelineTaskRuns[3].Error.Valid)
		require.Equal(t, pipeline.JSONSerializable{Val: []interface{}{interface{}(nil)}, Valid: true}, run.Outputs)
		require.Equal(t, pipeline.RunErrors{null.String{NullString: sql.NullString{String: "task inputs: too many errors", Valid: true}}}, run.FatalErrors)
	})
	// Delete the job
	{
		cltest.DeleteJobViaWeb(t, app, jobID)
		require.Eventually(t, func() bool { return eiNotifiedOfDelete }, 5*time.Second, 10*time.Millisecond, "expected external initiator to be notified of deleted job")
	}
}
