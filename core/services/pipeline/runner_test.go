package pipeline_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/pelletier/go-toml"

	"github.com/smartcontractkit/chainlink/core/services"

	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
)

func TestRunner(t *testing.T) {
	config, oldORM, cleanupDB := cltest.BootstrapThrowawayORM(t, "pipeline_runner", true, true)
	defer cleanupDB()
	config.Set("DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS", true)
	db := oldORM.DB
	eventBroadcaster := postgres.NewEventBroadcaster(config.DatabaseURL(), 0, 0)
	eventBroadcaster.Start()
	defer eventBroadcaster.Stop()

	pipelineORM := pipeline.NewORM(db, config, eventBroadcaster)
	runner := pipeline.NewRunner(pipelineORM, config)
	jobORM := job.NewORM(db, config, pipelineORM, eventBroadcaster, &postgres.NullAdvisoryLocker{})
	defer jobORM.Close()

	runner.Start()
	defer runner.Stop()

	t.Run("gets the election result winner", func(t *testing.T) {
		var httpURL string
		{
			mockElectionWinner, cleanupElectionWinner := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `Hal Finney`)
			defer cleanupElectionWinner()
			mockVoterTurnout, cleanupVoterTurnout := cltest.NewHTTPMockServer(t, http.StatusOK, "POST", `{"data": {"result": 62.57}}`)
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
		ocrSpec, dbSpec := makeVoterTurnoutOCRJobSpecWithHTTPURL(t, db, httpURL)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 2)
		assert.NoError(t, results[0].Error)
		assert.NoError(t, results[1].Error)
		assert.Equal(t, "6225.6", results[0].Value)
		assert.Equal(t, "Hal Finney", results[1].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		assert.Len(t, runs, 9)

		for _, run := range runs {
			if run.DotID() == "answer2" {
				assert.Equal(t, "Hal Finney", run.Output.Val)
			} else if run.DotID() == "ds2" {
				assert.Equal(t, `{"turnout": 61.942}`, run.Output.Val)
			} else if run.DotID() == "ds2_parse" {
				assert.Equal(t, float64(61.942), run.Output.Val)
			} else if run.DotID() == "ds2_multiply" {
				assert.Equal(t, "6194.2", run.Output.Val)
			} else if run.DotID() == "ds1" {
				assert.Equal(t, `{"data": {"result": 62.57}}`, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.Equal(t, float64(62.57), run.Output.Val)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "6257", run.Output.Val)
			} else if run.DotID() == "answer1" {
				assert.Equal(t, "6225.6", run.Output.Val)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{"6225.6", "Hal Finney"}, run.Output.Val)
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
			}
		}
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
		ocrSpec, dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, httpURL, false)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.EqualError(t, results[0].Error, "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results[0].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 4)

		for _, run := range runs {
			if run.DotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				require.NotNil(t, resp, run.Output)
				assert.Equal(t, resp, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				// FIXME: Shouldn't it be the Val that is null?
				assert.Nil(t, run.Output)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "type <nil> cannot be converted to decimal.Decimal", run.Error.ValueOrZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{nil}, run.Output.Val)
				assert.Equal(t, "[\"type \\u003cnil\\u003e cannot be converted to decimal.Decimal\"]", run.Error.ValueOrZero())
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
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
		ocrSpec, dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, httpURL, false)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.EqualError(t, results[0].Error, "could not resolve path [\"USD\"] in {\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}")
		assert.Nil(t, results[0].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 4)

		for _, run := range runs {
			if run.DotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.Equal(t, "could not resolve path [\"USD\"] in {\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}", run.Error.ValueOrZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "could not resolve path [\"USD\"] in {\"Response\":\"Error\",\"Message\":\"You are over your rate limit please upgrade your account!\",\"HasWarning\":false,\"Type\":99,\"RateLimit\":{\"calls_made\":{\"second\":5,\"minute\":5,\"hour\":955,\"day\":10004,\"month\":15146,\"total_calls\":15152},\"max_calls\":{\"second\":20,\"minute\":300,\"hour\":3000,\"day\":10000,\"month\":75000}},\"Data\":{}}", run.Error.ValueOrZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{nil}, run.Output.Val)
				assert.Equal(t, "[\"could not resolve path [\\\"USD\\\"] in {\\\"Response\\\":\\\"Error\\\",\\\"Message\\\":\\\"You are over your rate limit please upgrade your account!\\\",\\\"HasWarning\\\":false,\\\"Type\\\":99,\\\"RateLimit\\\":{\\\"calls_made\\\":{\\\"second\\\":5,\\\"minute\\\":5,\\\"hour\\\":955,\\\"day\\\":10004,\\\"month\\\":15146,\\\"total_calls\\\":15152},\\\"max_calls\\\":{\\\"second\\\":20,\\\"minute\\\":300,\\\"hour\\\":3000,\\\"day\\\":10000,\\\"month\\\":75000}},\\\"Data\\\":{}}\"]", run.Error.ValueOrZero())
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
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
		ocrSpec, dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, httpURL, true)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the final pipeline results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.EqualError(t, results[0].Error, "type <nil> cannot be converted to decimal.Decimal")
		assert.Nil(t, results[0].Value)

		// Verify individual task results
		var runs []pipeline.TaskRun
		err = db.
			Preload("PipelineTaskSpec").
			Where("pipeline_run_id = ?", runID).
			Find(&runs).Error
		assert.NoError(t, err)
		require.Len(t, runs, 4)

		for _, run := range runs {
			if run.DotID() == "ds1" {
				assert.True(t, run.Error.IsZero())
				assert.Equal(t, resp, run.Output.Val)
			} else if run.DotID() == "ds1_parse" {
				assert.True(t, run.Error.IsZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "ds1_multiply" {
				assert.Equal(t, "type <nil> cannot be converted to decimal.Decimal", run.Error.ValueOrZero())
				assert.Nil(t, run.Output)
			} else if run.DotID() == "__result__" {
				assert.Equal(t, []interface{}{nil}, run.Output.Val)
				assert.Equal(t, "[\"type \\u003cnil\\u003e cannot be converted to decimal.Decimal\"]", run.Error.ValueOrZero())
			} else {
				t.Fatalf("unknown task '%v'", run.DotID())
			}
		}
	})

	t.Run("missing required env vars", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		var os = offchainreporting.OracleSpec{
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
		_, err := services.ValidatedOracleSpecToml(config.Config, fmt.Sprintf(s, cltest.NewEIP55Address(), "http://blah.com", ""))
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		js := models.JobSpecV2{
			MaxTaskDuration:             models.Interval(cltest.MustParseDuration(t, "1s")),
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               os.SchemaVersion,
		}
		err = jobORM.CreateJob(context.Background(), &js, os.TaskDAG())
		require.NoError(t, err)
		var jb models.JobSpecV2
		err = db.Preload("OffchainreportingOracleSpec", "id = ?", js.ID).
			Find(&jb).Error
		require.NoError(t, err)
		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewJobSpawnerDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			nil,
			nil)
		_, err = sd.ServicesForSpec(sd.FromDBRow(jb))
		// We expect this to fail as neither the required vars are not set either via the env nor the job itself.
		require.Error(t, err)
	})

	t.Run("use env for minimal bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		var os = offchainreporting.OracleSpec{
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
		_, err = services.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		js := models.JobSpecV2{
			MaxTaskDuration:             models.Interval(cltest.MustParseDuration(t, "1s")),
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               os.SchemaVersion,
		}
		err = jobORM.CreateJob(context.Background(), &js, os.TaskDAG())
		require.NoError(t, err)
		var jb models.JobSpecV2
		err = db.Preload("OffchainreportingOracleSpec", "id = ?", js.ID).
			Find(&jb).Error
		require.NoError(t, err)
		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewJobSpawnerDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			nil,
			nil)
		_, err = sd.ServicesForSpec(sd.FromDBRow(jb))
		require.NoError(t, err)
	})

	t.Run("use env for minimal non-bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		kb, _, err := keyStore.GenerateEncryptedOCRKeyBundle()
		require.NoError(t, err)
		var os = offchainreporting.OracleSpec{
			Pipeline: *pipeline.NewTaskDAG(),
		}
		s := `
		type               = "offchainreporting"
		schemaVersion      = 1
		contractAddress    = "%s"
		isBootstrapPeer    = false 
		observationTimeout = "10s"
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
		config.Set("OCR_TRANSMITTER_ADDRESS", cltest.DefaultKey)
		_, err = services.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		js := models.JobSpecV2{
			MaxTaskDuration:             models.Interval(cltest.MustParseDuration(t, "1s")),
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               os.SchemaVersion,
		}
		err = jobORM.CreateJob(context.Background(), &js, os.TaskDAG())
		require.NoError(t, err)
		var jb models.JobSpecV2
		err = db.Preload("OffchainreportingOracleSpec", "id = ?", js.ID).
			Find(&jb).Error
		require.NoError(t, err)
		assert.Equal(t, jb.MaxTaskDuration, models.Interval(cltest.MustParseDuration(t, "1s")))

		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewJobSpawnerDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			nil,
			nil)
		_, err = sd.ServicesForSpec(sd.FromDBRow(jb))
		require.NoError(t, err)
	})

	t.Run("test min non-bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		kb, _, err := keyStore.GenerateEncryptedOCRKeyBundle()
		require.NoError(t, err)
		var os = offchainreporting.OracleSpec{
			Pipeline: *pipeline.NewTaskDAG(),
		}

		s := fmt.Sprintf(minimalNonBootstrapTemplate, cltest.NewEIP55Address(), ek.PeerID, cltest.DefaultKey, kb.ID, "http://blah.com", "")
		_, err = services.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)

		err = jobORM.CreateJob(context.Background(), &models.JobSpecV2{
			MaxTaskDuration:             models.Interval(cltest.MustParseDuration(t, "1s")),
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               os.SchemaVersion,
		}, os.TaskDAG())
		require.NoError(t, err)
		var jb models.JobSpecV2
		err = db.Preload("OffchainreportingOracleSpec", "p2p_peer_id = ?", ek.PeerID).
			Find(&jb).Error
		require.NoError(t, err)
		assert.Equal(t, jb.MaxTaskDuration, models.Interval(cltest.MustParseDuration(t, "1s")))

		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewJobSpawnerDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			nil,
			nil)
		_, err = sd.ServicesForSpec(sd.FromDBRow(jb))
		require.NoError(t, err)
	})

	t.Run("test min bootstrap", func(t *testing.T) {
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		var os = offchainreporting.OracleSpec{
			Pipeline: *pipeline.NewTaskDAG(),
		}
		s := fmt.Sprintf(minimalBootstrapTemplate, cltest.NewEIP55Address(), ek.PeerID)
		_, err = services.ValidatedOracleSpecToml(config.Config, s)
		require.NoError(t, err)
		err = toml.Unmarshal([]byte(s), &os)
		require.NoError(t, err)
		err = jobORM.CreateJob(context.Background(), &models.JobSpecV2{
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               os.SchemaVersion,
		}, os.TaskDAG())
		require.NoError(t, err)
		var jb models.JobSpecV2
		err = db.Preload("OffchainreportingOracleSpec", "p2p_peer_id = ?", ek.PeerID).
			Find(&jb).Error
		require.NoError(t, err)

		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewJobSpawnerDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			nil,
			nil)
		_, err = sd.ServicesForSpec(sd.FromDBRow(jb))
		require.NoError(t, err)
	})

	t.Run("test job spec error is created", func(t *testing.T) {
		// Create a keystore with an ocr key bundle and p2p key.
		keyStore := offchainreporting.NewKeyStore(db, utils.GetScryptParams(config.Config))
		_, ek, err := keyStore.GenerateEncryptedP2PKey()
		require.NoError(t, err)
		kb, _, err := keyStore.GenerateEncryptedOCRKeyBundle()
		require.NoError(t, err)
		spec := fmt.Sprintf(ocrJobSpecTemplate, cltest.NewAddress().Hex(), ek.PeerID, kb.ID, cltest.DefaultKey, fmt.Sprintf(simpleFetchDataSourceTemplate, "blah", true))
		ocrspec, dbSpec := makeOCRJobSpecWithHTTPURL(t, db, spec)

		// Create an OCR job
		err = jobORM.CreateJob(context.Background(), dbSpec, ocrspec.TaskDAG())
		require.NoError(t, err)
		var jb models.JobSpecV2
		err = db.Preload("OffchainreportingOracleSpec", "p2p_peer_id = ?", ek.PeerID).
			Find(&jb).Error
		require.NoError(t, err)

		config.Config.Set("P2P_LISTEN_PORT", 2000) // Required to create job spawner delegate.
		sd := offchainreporting.NewJobSpawnerDelegate(
			db,
			jobORM,
			config.Config,
			keyStore,
			nil,
			nil,
			nil)
		services, err := sd.ServicesForSpec(sd.FromDBRow(jb))
		require.NoError(t, err)

		// Start and stop the service to generate errors.
		// We expect a database timeout and a context cancellation
		// error to show up as pipeline_spec_errors.
		for _, s := range services {
			err = s.Start()
			require.NoError(t, err)
			err = s.Close()
			require.NoError(t, err)
		}

		var se []models.JobSpecErrorV2
		err = db.Find(&se).Error
		require.NoError(t, err)
		require.Len(t, se, 2)
		assert.Equal(t, uint(1), se[0].Occurrences)
		assert.Equal(t, uint(1), se[1].Occurrences)

		// Ensure we can delete an errored job.
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
		ocrSpec, dbSpec := makeSimpleFetchOCRJobSpecWithHTTPURL(t, db, httpURL, false)
		err := jobORM.CreateJob(context.Background(), dbSpec, ocrSpec.TaskDAG())
		require.NoError(t, err)

		runID, err := runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.NoError(t, err)

		// Verify the results
		results, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)

		assert.Len(t, results, 1)
		assert.Nil(t, results[0].Error)
		assert.Equal(t, "4242", results[0].Value)

		// Delete the job
		err = jobORM.DeleteJob(ctx, dbSpec.ID)
		require.NoError(t, err)

		// Create another run
		_, err = runner.CreateRun(context.Background(), dbSpec.ID, nil)
		require.EqualError(t, err, fmt.Sprintf("no job found with id %v (most likely it was deleted)", dbSpec.ID))

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = runner.AwaitRun(ctx, runID)
		require.EqualError(t, err, fmt.Sprintf("could not determine if run is finished (run ID: %v): record not found", runID))
	})

	t.Run("timeouts", func(t *testing.T) {
		// There are 4 timeouts:
		// - ObservationTimeout = how long the whole OCR time needs to run, or it fails (default 10 seconds)
		// - config.JobPipelineMaxTaskDuration() = node level maximum time for a pipeline task (default 10 minutes)
		// - config.DefaultHTTPTimeout() * config.DefaultMaxHTTPAttempts() = global, http specific timeouts (default 15s * 5 retries = 75s)
		// - "d1 [.... timeout="2s"]" = per task level timeout (should override the global config)
		serv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(1 * time.Millisecond)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"USD":10.1}`))
		}))
		defer serv.Close()

		os := makeMinimalHTTPOracleSpec(t, cltest.NewEIP55Address().String(), cltest.DefaultPeerID, cltest.DefaultKey, cltest.DefaultOCRKeyBundleID, serv.URL, `timeout="1ns"`)
		jb := &models.JobSpecV2{
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Name:                        null.NewString("a job", true),
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               1,
		}
		err := jobORM.CreateJob(context.Background(), jb, os.TaskDAG())
		require.NoError(t, err)
		runID, err := runner.CreateRun(context.Background(), jb.ID, nil)
		require.NoError(t, err)
		err = runner.AwaitRun(context.Background(), runID)
		require.NoError(t, err)
		r, err := runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)
		assert.Error(t, r[0].Error)

		// No task timeout should succeed.
		os = makeMinimalHTTPOracleSpec(t, cltest.NewEIP55Address().String(), cltest.DefaultPeerID, cltest.DefaultKey, cltest.DefaultOCRKeyBundleID, serv.URL, "")
		jb = &models.JobSpecV2{
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Name:                        null.NewString("a job 2", true),
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               1,
		}
		err = jobORM.CreateJob(context.Background(), jb, os.TaskDAG())
		require.NoError(t, err)
		runID, err = runner.CreateRun(context.Background(), jb.ID, nil)
		require.NoError(t, err)
		err = runner.AwaitRun(context.Background(), runID)
		require.NoError(t, err)
		r, err = runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)
		assert.Equal(t, 10.1, r[0].Value)
		assert.NoError(t, r[0].Error)

		// Job specified task timeout should fail.
		os = makeMinimalHTTPOracleSpec(t, cltest.NewEIP55Address().String(), cltest.DefaultPeerID, cltest.DefaultKey, cltest.DefaultOCRKeyBundleID, serv.URL, "")
		jb = &models.JobSpecV2{
			MaxTaskDuration:             models.Interval(time.Duration(1)),
			OffchainreportingOracleSpec: &os.OffchainReportingOracleSpec,
			Name:                        null.NewString("a job 3", true),
			Type:                        string(offchainreporting.JobType),
			SchemaVersion:               1,
		}
		err = jobORM.CreateJob(context.Background(), jb, os.TaskDAG())
		require.NoError(t, err)
		runID, err = runner.CreateRun(context.Background(), jb.ID, nil)
		require.NoError(t, err)
		err = runner.AwaitRun(context.Background(), runID)
		require.NoError(t, err)
		r, err = runner.ResultsForRun(context.Background(), runID)
		require.NoError(t, err)
		assert.Error(t, r[0].Error)
	})
}
