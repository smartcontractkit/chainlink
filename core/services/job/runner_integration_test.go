package job_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/pkg/errors"

	"github.com/pelletier/go-toml"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/telemetry"

	"github.com/smartcontractkit/chainlink/core/store/models"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

var monitoringEndpoint = ocrtypes.MonitoringEndpoint(&telemetry.NoopAgent{})

func TestRunner(t *testing.T) {
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "pipeline_runner", true, true)
	defer cleanupDB()
	config.Set("DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS", true)
	db := oldORM.DB
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	defer eventBroadcaster.Close()

	pipelineORM := pipeline.NewORM(db, config, eventBroadcaster)
	runner := pipeline.NewRunner(pipelineORM, config)
	jobORM := job.NewORM(db, config.Config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer jobORM.Close()

	runner.Start()
	defer runner.Close()

	key := cltest.MustInsertRandomKey(t, db, 0)
	transmitterAddress := key.Address.Address()

	ethClient, _, _ := cltest.NewEthMocks(t)
	ethClient.On("HeaderByNumber", mock.Anything, (*big.Int)(nil)).Return(cltest.Head(10), nil)
	ethClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil, nil)

	t.Run("gets the election result winner", func(t *testing.T) {
		var httpURL string
		{
			mockElectionWinner, cleanupElectionWinner := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `Hal Finney`,
				func(header http.Header, s string) {
					var md models.BridgeMetaDataJSON
					require.NoError(t, json.Unmarshal([]byte(s), &md))
					assert.Equal(t, big.NewInt(10), md.Meta.LatestAnswer)
					assert.Equal(t, big.NewInt(100), md.Meta.UpdatedAt)
				})
			defer cleanupElectionWinner()
			mockVoterTurnout, cleanupVoterTurnout := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"data": {"result": 62.57}}`,
				func(header http.Header, s string) {
					var md models.BridgeMetaDataJSON
					require.NoError(t, json.Unmarshal([]byte(s), &md))
					assert.Equal(t, big.NewInt(10), md.Meta.LatestAnswer)
					assert.Equal(t, big.NewInt(100), md.Meta.UpdatedAt)
				},
			)
			defer cleanupVoterTurnout()
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"turnout": 61.942}`)
			defer cleanupHTTP()

			_, bridgeER := cltest.NewBridgeType(t, "election_winner", mockElectionWinner.URL)
			err := db.Create(bridgeER).Error
			require.NoError(t, err)

			_, bridgeVT := cltest.NewBridgeType(t, "voter_turnout", mockVoterTurnout.URL)
			err = db.Create(bridgeVT).Error
			require.NoError(t, err)

			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		dbSpec := MakeVoterTurnoutOCRJobSpecWithHTTPURL(t, db, transmitterAddress, httpURL)
		err := jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		jb, err := jobORM.FindJob(dbSpec.ID)
		require.NoError(t, err)
		m, err := models.MarshalBridgeMetaData(big.NewInt(10), big.NewInt(100))
		require.NoError(t, err)
		runID, results, err := runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{Val: m}, *logger.Default, true)
		require.NoError(t, err)

		require.Len(t, results.Values, 2)
		assert.Nil(t, results.Errors[0])
		assert.Nil(t, results.Errors[1])
		assert.Equal(t, "6225.6", results.Values[0].(decimal.Decimal).String())
		assert.Equal(t, "Hal Finney", results.Values[1].(string))

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
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
		_, bridge := cltest.NewBridgeType(t, "testbridge", "http://blah.com")
		require.NoError(t, db.Create(bridge).Error)
		dbSpec := makeOCRJobSpecFromToml(t, db, `
			type               = "offchainreporting"
			schemaVersion      = 1
			observationSource = """
				ds1          [type=bridge name="testbridge"];
			"""
		`)
		require.NoError(t, jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline))
		// Should not be able to delete a bridge in use.
		jids, err := jobORM.FindJobIDsWithBridge(bridge.Name.String())
		require.NoError(t, err)
		require.Equal(t, 1, len(jids))

		// But if we delete the job, then we can.
		require.NoError(t, jobORM.DeleteJob(context.Background(), dbSpec.ID))
		jids, err = jobORM.FindJobIDsWithBridge(bridge.Name.String())
		require.NoError(t, err)
		require.Equal(t, 0, len(jids))
	})

	t.Run("referencing a non-existent bridge should error", func(t *testing.T) {
		_, bridge := cltest.NewBridgeType(t, "testbridge2", "http://blah.com")
		require.NoError(t, db.Create(bridge).Error)
		dbSpec := makeOCRJobSpecFromToml(t, db, `
			type               = "offchainreporting"
			schemaVersion      = 1
			observationSource = """
				ds1          [type=bridge name="testbridge2"];
			"""
		`)
		require.Error(t,
			pipeline.ErrNoSuchBridge,
			errors.Cause(jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)))
	})

	config.Set("DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS", false)

	t.Run("handles the case where the parsed value is literally null", func(t *testing.T) {
		var httpURL string
		resp := `{"USD": null}`
		{
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			defer cleanupHTTP()
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, transmitterAddress, httpURL, false)
		err := jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		jb, err := jobORM.FindJob(dbSpec.ID)
		require.NoError(t, err)
		runID, results, err := runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)

		assert.Len(t, results.Errors, 1)
		assert.Len(t, results.Values, 1)
		assert.Contains(t, results.Errors[0].Error(), "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results.Values[0])

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 3)

		for _, run := range runs {
			if run.GetDotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				require.NotNil(t, resp, run.Output)
				assert.Equal(t, resp, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				// FIXME: Shouldn't it be the Val that is null?
				assert.Nil(t, run.Output)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Contains(t, run.Error.ValueOrZero(), "type <nil> cannot be converted to decimal.Decimal")
				assert.Nil(t, run.Output)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("handles the case where the jsonparse lookup path is missing from the http response", func(t *testing.T) {
		var httpURL string
		resp := "{\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}"
		{
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			defer cleanupHTTP()
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, transmitterAddress, httpURL, false)
		err := jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		jb, err := jobORM.FindJob(dbSpec.ID)
		require.NoError(t, err)
		runID, results, err := runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)

		assert.Len(t, results.Values, 1)
		assert.Len(t, results.Errors, 1)
		assert.Contains(t, results.Errors[0].Error(), pipeline.ErrTooManyErrors.Error())
		assert.Nil(t, results.Values[0])

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 3)

		for _, run := range runs {
			if run.GetDotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.Contains(t, run.Error.ValueOrZero(), "could not resolve path [\"USD\"] in {\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}")
				assert.Nil(t, run.Output)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Contains(t, run.Error.ValueOrZero(), pipeline.ErrTooManyErrors.Error())
				assert.Nil(t, run.Output)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("handles the case where the jsonparse lookup path is missing from the http response and lax is enabled", func(t *testing.T) {
		var httpURL string
		resp := "{\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}"
		{
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			defer cleanupHTTP()
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, transmitterAddress, httpURL, true)
		err := jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)
		jb, err := jobORM.FindJob(dbSpec.ID)
		require.NoError(t, err)

		runID, results, err := runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)

		assert.Len(t, results.Values, 1)
		assert.Contains(t, results.Errors[0].Error(), "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results.Values[0])

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 3)

		for _, run := range runs {
			if run.GetDotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.GetDotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				assert.Nil(t, run.Output)
			} else if run.GetDotID() == "ds1_multiply" {
				assert.Contains(t, run.Error.ValueOrZero(), "type <nil> cannot be converted to decimal.Decimal")
				assert.Nil(t, run.Output)
			} else {
				t.Fatalf("unknown task '%v'", run.GetDotID())
			}
		}
	})

	t.Run("missing required env vars", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		var os = job.Job{
			Pipeline: *pipeline.NewTaskDAG(),
		}
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
		os, err := offchainreporting.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		os.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(context.Background(), &os, os.Pipeline)
		require.NoError(t, err)
		var jb job.Job
		err = db.Preload("PipelineSpec").
			Preload("OffchainreportingOracleSpec").Where("id = ?", os.ID).
			First(&jb).Error
		require.NoError(t, err)
		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			ethClient,
			nil,
			nil,
			monitoringEndpoint)
		_, err = sd.ServicesForSpec(jb)
		// We expect this to fail as neither the required vars are not set either via the env nor the job itself.
		require.Error(t, err)
	})

	t.Run("use env for minimal bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		var os = job.Job{
			Pipeline: *pipeline.NewTaskDAG(),
		}
		s := `
		type               = "offchainreporting"
		schemaVersion      = 1
		contractAddress    = "%s"
		isBootstrapPeer    = true
`
		config.Set("P2P_PEER_ID", ek.PeerID)
		s = fmt.Sprintf(s, cltest.NewEIP55Address())
		os, err = offchainreporting.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		os.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(context.Background(), &os, os.Pipeline)
		require.NoError(t, err)
		var jb job.Job
		err = db.Preload("PipelineSpec").
			Preload("OffchainreportingOracleSpec").
			Where("id = ?", os.ID).
			First(&jb).Error
		require.NoError(t, err)
		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.

		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, config.Config, db)
		require.NoError(t, pw.Start())
		sd := offchainreporting.NewDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			ethClient,
			nil,
			pw,
			monitoringEndpoint,
		)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("use env for minimal non-bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		kb, _, err := keyStore.GenerateEncryptedOCRKeyBundle()
		require.NoError(t, err)
		var os = job.Job{
			Pipeline: *pipeline.NewTaskDAG(),
		}
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
		config.Set("P2P_PEER_ID", ek.PeerID)
		config.Set("P2P_BOOTSTRAP_PEERS", []string{"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
			"/dns4/chain.link/tcp/1235/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"})
		config.Set("OCR_KEY_BUNDLE_ID", kb.ID.String())
		config.Set("OCR_TRANSMITTER_ADDRESS", transmitterAddress)
		os, err = offchainreporting.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		os.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(context.Background(), &os, os.Pipeline)
		require.NoError(t, err)
		var jb job.Job
		err = db.Preload("PipelineSpec").
			Preload("OffchainreportingOracleSpec").Where("id = ?", os.ID).
			First(&jb).Error
		require.NoError(t, err)
		// Assert the override
		assert.Equal(t, jb.OffchainreportingOracleSpec.ObservationTimeout, models.Interval(cltest.MustParseDuration(t, "15s")))
		// Assert that this is unset
		assert.Equal(t, jb.OffchainreportingOracleSpec.BlockchainTimeout, models.Interval(0))
		assert.Equal(t, jb.MaxTaskDuration, models.Interval(cltest.MustParseDuration(t, "1s")))

		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, config.Config, db)
		require.NoError(t, pw.Start())
		sd := offchainreporting.NewDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			ethClient,
			nil,
			pw,
			monitoringEndpoint)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("test min non-bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		kb, _, err := keyStore.GenerateEncryptedOCRKeyBundle()
		require.NoError(t, err)
		var os = job.Job{
			Pipeline: *pipeline.NewTaskDAG(),
		}

		s := fmt.Sprintf(minimalNonBootstrapTemplate, cltest.NewEIP55Address(), ek.PeerID, transmitterAddress.Hex(), kb.ID, "http://blah.com", "")
		os, err = offchainreporting.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)

		os.MaxTaskDuration = models.Interval(cltest.MustParseDuration(t, "1s"))
		err = jobORM.CreateJob(context.Background(), &os, os.Pipeline)
		require.NoError(t, err)
		var jb job.Job
		err = db.Preload("PipelineSpec").
			Preload("OffchainreportingOracleSpec").
			Where("id = ?", os.ID).
			First(&jb).Error
		require.NoError(t, err)
		assert.Equal(t, jb.MaxTaskDuration, models.Interval(cltest.MustParseDuration(t, "1s")))

		config.Config.Set("P2P_LISTEN_PORT", 2000)           // Required to create job spawner delegate.
		config.Config.Set("P2P_PEER_ID", ek.PeerID.String()) // Required to create job spawner delegate.
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, config.Config, db)
		require.NoError(t, pw.Start())
		sd := offchainreporting.NewDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			ethClient,
			nil,
			pw,
			monitoringEndpoint)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("test min bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		var os = job.Job{
			Pipeline: *pipeline.NewTaskDAG(),
		}
		s := fmt.Sprintf(minimalBootstrapTemplate, cltest.NewEIP55Address(), ek.PeerID)
		os, err = offchainreporting.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		err = jobORM.CreateJob(context.Background(), &os, os.Pipeline)
		require.NoError(t, err)
		var jb job.Job
		err = db.Preload("PipelineSpec").
			Preload("OffchainreportingOracleSpec").Where("id = ?", os.ID).
			First(&jb).Error
		require.NoError(t, err)

		config.Config.Set("P2P_LISTEN_PORT", 2000)           // Required to create job spawner delegate.
		config.Config.Set("P2P_PEER_ID", ek.PeerID.String()) // Required to create job spawner delegate.
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, config.Config, db)
		require.NoError(t, pw.Start())
		sd := offchainreporting.NewDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			ethClient,
			nil,
			pw,
			monitoringEndpoint)
		_, err = sd.ServicesForSpec(jb)
		require.NoError(t, err)
	})

	t.Run("test job spec error is created", func(t *testing.T) {
		// Create a keystore with an ocr key bundle and p2p key.
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		kb, _, err := keyStore.GenerateEncryptedOCRKeyBundle()
		require.NoError(t, err)
		spec := fmt.Sprintf(ocrJobSpecTemplate, cltest.NewAddress().Hex(), ek.PeerID, kb.ID, transmitterAddress.Hex(), fmt.Sprintf(simpleFetchDataSourceTemplate, "blah", true))
		dbSpec := makeOCRJobSpecFromToml(t, db, spec)

		// Create an OCR job
		err = jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)
		var jb job.Job
		err = db.Preload("PipelineSpec").
			Preload("OffchainreportingOracleSpec").Where("id = ?", dbSpec.ID).
			First(&jb).Error
		require.NoError(t, err)

		config.Config.Set("P2P_LISTEN_PORT", 2000)           // Required to create job spawner delegate.
		config.Config.Set("P2P_PEER_ID", ek.PeerID.String()) // Required to create job spawner delegate.
		pw := offchainreporting.NewSingletonPeerWrapper(keyStore, config.Config, db)
		require.NoError(t, pw.Start())

		sd := offchainreporting.NewDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			ethClient,
			log.NewBroadcaster(log.NewORM(db), ethClient, config, nil),
			pw,
			monitoringEndpoint)
		services, err := sd.ServicesForSpec(jb)
		require.NoError(t, err)

		// Return an error getting the contract code.
		ethClient.On("CodeAt", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("no such code"))
		for _, s := range services {
			err = s.Start()
			require.NoError(t, err)
		}
		var se []job.SpecError
		require.Eventually(t, func() bool {
			err = db.Find(&se).Error
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
		_, err = jobORM.ClaimUnclaimedJobs(context.Background())
		require.NoError(t, err)
		err = jobORM.DeleteJob(context.Background(), jb.ID)
		require.NoError(t, err)
		err = db.Find(&se).Error
		require.NoError(t, err)
		require.Len(t, se, 0)

		// Noop once the job is gone.
		jobORM.RecordError(context.Background(), jb.ID, "test")
		err = db.Find(&se).Error
		require.NoError(t, err)
		require.Len(t, se, 0)
	})

	t.Run("deleting jobs", func(t *testing.T) {
		var httpURL string
		{
			resp := `{"USD": 42.42}`
			mockHTTP, cleanupHTTP := cltest.NewHTTPMockServer(t, http.StatusOK, "GET", resp)
			defer cleanupHTTP()
			httpURL = mockHTTP.URL
		}

		// Need a job in order to create a run
		dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, transmitterAddress, httpURL, false)
		err := jobORM.CreateJob(context.Background(), dbSpec, dbSpec.Pipeline)
		require.NoError(t, err)

		jb, err := jobORM.FindJob(dbSpec.ID)
		require.NoError(t, err)
		_, results, err := runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)
		assert.Len(t, results.Values, 1)
		assert.Nil(t, results.Errors[0])
		assert.Equal(t, "4242", results.Values[0].(decimal.Decimal).String())

		// Delete the job
		err = jobORM.DeleteJob(context.Background(), dbSpec.ID)
		require.NoError(t, err)

		// Create another run
		_, _, err = runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.Error(t, err)
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

		jbs := makeMinimalHTTPOracleSpec(t, cltest.NewEIP55Address().String(), cltest.DefaultPeerID, transmitterAddress.Hex(), cltest.DefaultOCRKeyBundleID, serv.URL, `timeout="1ns"`)
		err := jobORM.CreateJob(context.Background(), jbs, jbs.Pipeline)
		require.NoError(t, err)

		jb, err := jobORM.FindJob(jbs.ID)
		require.NoError(t, err)
		_, results, err := runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)
		assert.Nil(t, results.Values[0])

		// No task timeout should succeed.
		jbs = makeMinimalHTTPOracleSpec(t, cltest.NewEIP55Address().String(), cltest.DefaultPeerID, transmitterAddress.Hex(), cltest.DefaultOCRKeyBundleID, serv.URL, "")
		err = jobORM.CreateJob(context.Background(), jbs, jbs.Pipeline)
		require.NoError(t, err)
		jb, err = jobORM.FindJob(jbs.ID)
		require.NoError(t, err)
		_, results, err = runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)
		assert.Equal(t, 10.1, results.Values[0])
		assert.Nil(t, results.Errors[0])

		// Job specified task timeout should fail.
		jbs = makeMinimalHTTPOracleSpec(t, cltest.NewEIP55Address().String(), cltest.DefaultPeerID, transmitterAddress.Hex(), cltest.DefaultOCRKeyBundleID, serv.URL, "")
		jbs.MaxTaskDuration = models.Interval(time.Duration(1))
		jbs.Name = null.NewString("a job 3", true)
		err = jobORM.CreateJob(context.Background(), jbs, jbs.Pipeline)
		require.NoError(t, err)
		jb, err = jobORM.FindJob(jbs.ID)
		require.NoError(t, err)

		_, results, err = runner.ExecuteAndInsertFinishedRun(context.Background(), *jb.PipelineSpec, nil, pipeline.JSONSerializable{}, *logger.Default, true)
		require.NoError(t, err)
		assert.NotNil(t, results.Errors[0])
	})
}
