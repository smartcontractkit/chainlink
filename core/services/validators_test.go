package services_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateJob(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	// Create a funding key.
	require.NoError(t, keyStore.Eth.Unlock(cltest.Password))
	fundingKey, _, err := keyStore.Eth.EnsureFundingKey()
	require.NoError(t, err)
	tests := []struct {
		name  string
		input []byte
		want  error
	}{
		{"base case", cltest.MustReadFile(t, "../testdata/jsonspecs/hello_world_job.json"), nil},
		{
			"error in job",
			cltest.MustReadFile(t, "../testdata/jsonspecs/invalid_endat_job.json"),
			models.NewJSONAPIErrorsWith("StartAt cannot be before EndAt"),
		},
		{
			"error in runat initr",
			cltest.MustReadFile(t, "../testdata/jsonspecs/run_at_wo_time_job.json"),
			models.NewJSONAPIErrorsWith("RunAt must have a time"),
		},
		{
			"error in task",
			cltest.MustReadFile(t, "../testdata/jsonspecs/nonexistent_task_job.json"),
			models.NewJSONAPIErrorsWith("idonotexist is not a supported adapter type"),
		},
		{
			"zero initiators",
			cltest.MustReadFile(t, "../testdata/jsonspecs/zero_initiators.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one initiator only",
			cltest.MustReadFile(t, "../testdata/jsonspecs/initiator_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one task only",
			cltest.MustReadFile(t, "../testdata/jsonspecs/task_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"runlog and ethtx with an address",
			cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_ethtx_w_address_job.json"),
			models.NewJSONAPIErrorsWith("Cannot set EthTx Task's address parameter with a RunLog Initiator"),
		},
		{
			"runlog and ethtx with a function selector",
			cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_ethtx_w_funcselector_job.json"),
			models.NewJSONAPIErrorsWith("Cannot set EthTx Task's function selector parameter with a RunLog Initiator"),
		},
		{
			"runlog and ethtx with a fromAddress that doesn't match one of our keys",
			cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_ethtx_w_missing_fromAddress_job.json"),
			models.NewJSONAPIErrorsWith("error address 0x0f416A5a298F05d386CfE8164f342Bec5b5E10D7 not in keystore finding key for address 0x0f416a5a298f05d386cfe8164f342bec5b5e10d7"),
		},
		{
			"runlog with two ethtx tasks",
			cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_2_ethlogs_job.json"),
			models.NewJSONAPIErrorsWith("Cannot RunLog initiated jobs cannot have more than one EthTx Task"),
		},
		{
			"cannot use funding key",
			[]byte(fmt.Sprintf(string(cltest.MustReadFile(t, "../testdata/jsonspecs/runlog_ethtx_template_fromAddress_job.json")), fundingKey.Address.String())),
			models.NewJSONAPIErrorsWith(fmt.Sprintf("address %v is a funding address, cannot use it to send transactions", fundingKey.Address.String())),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var j models.JobSpec
			assert.NoError(t, json.Unmarshal(test.input, &j))
			result := services.ValidateJob(j, store, keyStore)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateJob_RejectsSleepAdapterWhenExperimentalAdaptersAreDisabled(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	sleepingJob := cltest.NewJobWithWebInitiator()
	sleepingJob.Tasks[0].Type = adapters.TaskTypeSleep

	store.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", true)
	assert.NoError(t, services.ValidateJob(sleepingJob, store, keyStore))

	store.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", false)
	assert.Error(t, services.ValidateJob(sleepingJob, store, keyStore))
}

func TestValidateBridgeType(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	tests := []struct {
		description string
		request     models.BridgeTypeRequest
		want        error
	}{
		{
			"no adapter name",
			models.BridgeTypeRequest{
				URL: cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("No name specified"),
		},
		{
			"invalid adapter name",
			models.BridgeTypeRequest{
				Name: "invalid/adapter",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("task type validation: name invalid/adapter contains invalid characters"),
		},
		{
			"invalid with blank url",
			models.BridgeTypeRequest{
				Name: "validadaptername",
				URL:  cltest.WebURL(t, ""),
			},
			models.NewJSONAPIErrorsWith("URL must be present"),
		},
		{
			"valid url",
			models.BridgeTypeRequest{
				Name: "adapterwithvalidurl",
				URL:  cltest.WebURL(t, "//denergy"),
			},
			nil,
		},
		{
			"valid docker url",
			models.BridgeTypeRequest{
				Name: "adapterwithdockerurl",
				URL:  cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
			},
			nil,
		},
		{
			"valid MinimumContractPayment positive",
			models.BridgeTypeRequest{
				Name:                   "adapterwithdockerurl",
				URL:                    cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
				MinimumContractPayment: assets.NewLink(1),
			},
			nil,
		},
		{
			"invalid MinimumContractPayment negative",
			models.BridgeTypeRequest{
				Name:                   "adapterwithdockerurl",
				URL:                    cltest.WebURL(t, "http://chainlink_cmc-adapter_1:8080"),
				MinimumContractPayment: assets.NewLink(-1),
			},
			models.NewJSONAPIErrorsWith("MinimumContractPayment must be positive"),
		},
		{
			"existing core adapter",
			models.BridgeTypeRequest{
				Name: "ethtx",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			models.NewJSONAPIErrorsWith("Bridge Type ethtx is a native adapter"),
		},
		{
			"new external adapter",
			models.BridgeTypeRequest{
				Name: "gdaxprice",
				URL:  cltest.WebURL(t, "https://denergy.eth"),
			},
			nil,
		}}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := services.ValidateBridgeType(&test.request, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateBridgeNotExist(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// Create a duplicate
	bt := models.BridgeType{}
	bt.Name = models.MustNewTaskType("solargridreporting")
	bt.URL = cltest.WebURL(t, "https://denergy.eth")
	assert.NoError(t, store.CreateBridgeType(&bt))

	newBridge := models.BridgeTypeRequest{
		Name: "solargridreporting",
	}
	expected := models.NewJSONAPIErrorsWith("Bridge Type solargridreporting already exists")
	result := services.ValidateBridgeTypeNotExist(&newBridge, store)
	assert.Equal(t, expected, result)
}

func TestValidateExternalInitiator(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	url := cltest.WebURL(t, "https://a.web.url")

	//  Add duplicate
	exi := models.ExternalInitiator{
		Name: "duplicate",
		URL:  &url,
	}

	assert.NoError(t, store.CreateExternalInitiator(&exi))

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"basic", `{"name":"bitcoin","url":"https://test.url"}`, false},
		{"basic w/ underscore", `{"name":"bit_coin","url":"https://test.url"}`, false},
		{"basic w/ underscore in url", `{"name":"bitcoin","url":"https://chainlink_bit-coin_1.url"}`, false},
		{"missing url", `{"name":"missing_url"}`, false},
		{"duplicate name", `{"name":"duplicate","url":"https://test.url"}`, true},
		{"invalid name characters", `{"name":"<invalid>","url":"https://test.url"}`, true},
		{"missing name", `{"url":"https://test.url"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var exr models.ExternalInitiatorRequest

			assert.NoError(t, json.Unmarshal([]byte(test.input), &exr))
			result := services.ValidateExternalInitiator(&exr, store)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}

func TestValidateInitiator(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	startAt := time.Now()
	endAt := startAt.Add(time.Second)
	job := cltest.NewJob()
	job.StartAt = cltest.NullableTime(startAt)
	job.EndAt = cltest.NullableTime(endAt)
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"web", `{"type":"web"}`, false},
		{"ethlog", `{"type":"ethlog"}`, false},
		{"external", `{"type":"external","params":{"name":"bitcoin"}}`, false},
		{"runlog", `{"type":"runlog"}`, false},
		{"runat", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, utils.ISO8601UTC(startAt)), false},
		{"runat w/o time", `{"type":"runat"}`, true},
		{"runat w time before start at", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, startAt.Add(-1*time.Second).Unix()), true},
		{"runat w time after end at", fmt.Sprintf(`{"type":"runat","params": {"time":"%v"}}`, endAt.Add(time.Second).Unix()), true},
		{"cron standard", `{"type":"cron","params": {"schedule":"CRON_TZ=UTC * * * * *"}}`, false},
		{"cron with 6 fields", `{"type":"cron","params": {"schedule":"CRON_TZ=UTC * * * * * *"}}`, false},
		{"cron w/o schedule", `{"type":"cron"}`, true},
		{"external w/o name", `{"type":"external"}`, true},
		{"non-existent initiator", `{"type":"doesntExist"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var initr models.Initiator
			assert.NoError(t, json.Unmarshal([]byte(test.input), &initr))
			result := services.ValidateInitiator(initr, job, store)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}

func TestValidateServiceAgreement(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
	err := keyStore.Eth.Unlock(cltest.Password)
	_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, keyStore.Eth, 0)
	assert.NoError(t, err)

	oracles := []string{fromAddress.Hex()}

	basic := string(cltest.MustReadFile(t, "../testdata/jsonspecs/hello_world_agreement.json"))
	basic = cltest.MustJSONSet(t, basic, "oracles", oracles)
	threeDays, _ := time.ParseDuration("72h")
	basic = cltest.MustJSONSet(t, basic, "endAt", time.Now().Add(threeDays))

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"basic", basic, false},
		{"no payment", cltest.MustJSONDel(t, basic, "payment"), true},
		{"less than minimum payment", cltest.MustJSONSet(t, basic, "payment", "1"), true},
		{"less than minimum expiration", cltest.MustJSONSet(t, basic, "expiration", 1), true},
		{"without being listed as an oracle", cltest.MustJSONSet(t, basic, "oracles", []string{}), true},
		{"past allowed end at", cltest.MustJSONSet(t, basic, "endAt", "3000-06-19T22:17:19Z"), true},
		{"before allowed end at", cltest.MustJSONSet(t, basic, "endAt", "2018-06-19T22:17:19Z"), true},
		{"more than one initiator should fail",
			cltest.MustJSONSet(t, basic, "initiators",
				[]models.Initiator{{
					JobSpecID:       models.NewJobID(),
					Type:            models.InitiatorServiceAgreementExecutionLog,
					InitiatorParams: models.InitiatorParams{},
				}, {
					JobSpecID:       models.NewJobID(),
					Type:            models.InitiatorWeb,
					InitiatorParams: models.InitiatorParams{},
				},
				}),
			true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sa, err := cltest.ServiceAgreementFromString(test.input)
			require.NoError(t, err)

			result := services.ValidateServiceAgreement(sa, store, keyStore)

			cltest.AssertError(t, test.wantError, result)
		})
	}
}

const validInitiator = `{
	"type": "fluxmonitor",
	"params": {
		"address": "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42",
		"requestdata": {
			"data":{"coin":"ETH","market":"USD"}
		},
		"feeds": [
			"https://lambda.staging.devnet.tools/bnc/call",
			"https://lambda.staging.devnet.tools/cc/call",
			"https://lambda.staging.devnet.tools/cmc/call"
		],
		"idleTimer": {
			"duration": "1m"
		},
		"pollTimer": {
			"period": "1m"
		},
		"threshold": 0.5,
		"precision": 2
	}
}`

func TestValidateInitiator_FluxMonitorHappy(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJob()
	var initr models.Initiator
	require.NoError(t, json.Unmarshal([]byte(validInitiator), &initr))
	err := services.ValidateInitiator(initr, job, store)
	require.NoError(t, err)
}

func TestValidateInitiator_FluxMonitorErrors(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJob()
	tests := []struct {
		Field   string
		JSONStr string
	}{
		{"address", cltest.MustJSONDel(t, validInitiator, "params.address")},
		{"feeds", cltest.MustJSONSet(t, validInitiator, "params.feeds", []string{})},
		{"threshold", cltest.MustJSONDel(t, validInitiator, "params.threshold")},
		{"must be positive", cltest.MustJSONSet(t, validInitiator, "params.threshold", -5)},
		{"requestdata", cltest.MustJSONDel(t, validInitiator, "params.requestdata")},
		{"pollTimer enabled, but no period specified", cltest.MustJSONDel(t, validInitiator, "params.pollTimer.period")},
		{"period must be equal or greater than 15s", cltest.MustJSONSet(t, validInitiator, "params.pollTimer.period", "1s")},
		{"idleTimer.duration must be >= than pollTimer.period", cltest.MustJSONSet(t, validInitiator, "params.idleTimer.duration", "30s")},
	}
	for _, test := range tests {
		t.Run("bad "+test.Field, func(t *testing.T) {
			var initr models.Initiator
			require.NoError(t, json.Unmarshal([]byte(test.JSONStr), &initr))
			err := services.ValidateInitiator(initr, job, store)
			require.Error(t, err)
			assert.Contains(t, err.Error(), test.Field)
		})
	}
}

func TestValidateInitiator_FeedsHappy(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	bridge := &models.BridgeType{
		Name: models.MustNewTaskType("testbridge"),
		URL:  cltest.WebURL(t, "https://testing.com/bridges"),
	}
	require.NoError(t, store.CreateBridgeType(bridge))

	job := cltest.NewJob()
	var initr models.Initiator
	require.NoError(t, json.Unmarshal([]byte(validInitiator), &initr))
	initr.Feeds = cltest.JSONFromString(t, `["https://lambda.staging.devnet.tools/bnc/call", {"bridge": "testbridge"}]`)
	err := services.ValidateInitiator(initr, job, store)
	require.NoError(t, err)
}

func TestValidateInitiator_FeedsErrors(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	bridge := &models.BridgeType{
		Name: models.MustNewTaskType("testbridge"),
		URL:  cltest.WebURL(t, "https://testing.com/bridges"),
	}
	require.NoError(t, store.CreateBridgeType(bridge))

	job := cltest.NewJob()
	tests := []struct {
		description string
		FeedsJSON   string
	}{
		{"invalid url", `["invalid/url"]`},
		{"invalid bridge name", `[{"bridge": "doesnotexist"}]`},
		{"invalid url type", `[1]`},
		{"invalid bridge type", `[{"bridge": 1}]`},
		{"valid url, invalid bridge", `["http://example.com", {"bridge": "doesnotexist"}]`},
		{"invalid url, valid bridge", `["invalid/url", {"bridge": "testbridge"}]`},
		{"missing bridge", `[{"bridgeName": "doesnotexist"}]`},
		{"unsupported bridge properties", `[{"bridge": "testbridge", "foo": "bar"}]`},
		{"invalid entry", `["http://example.com", {"bridge": "testbridge"}, 1]`},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var initr models.Initiator
			require.NoError(t, json.Unmarshal([]byte(validInitiator), &initr))
			initr.Feeds = cltest.JSONFromString(t, test.FeedsJSON)
			err := services.ValidateInitiator(initr, job, store)
			require.Error(t, err)
		})
	}
}

func TestValidateJob_VRF_Happy(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	input := cltest.MustReadFile(t, "../testdata/jsonspecs/randomness_job.json")

	var j models.JobSpec
	assert.NoError(t, json.Unmarshal(input, &j))
	err := services.ValidateJob(j, store, keyStore)
	assert.NoError(t, err)
}

func TestValidateJob_VRF_Error(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	input := cltest.MustReadFile(t, "../testdata/jsonspecs/randomness_job.json")

	makeVRFJob := func() models.JobSpec {
		var job models.JobSpec
		assert.NoError(t, json.Unmarshal(input, &job))
		return job
	}

	missingPubKeyJson := models.JSON{
		Result: gjson.ParseBytes([]byte(`{}`)),
	}

	job1 := makeVRFJob()
	job2 := makeVRFJob()
	job3 := makeVRFJob()
	job4 := makeVRFJob()

	job1.Tasks[0].Params = missingPubKeyJson
	job2.Tasks[0].MinRequiredIncomingConfirmations.Uint32 = 0
	job3.Initiators[0].Address = utils.ZeroAddress
	job4.Initiators = append(job2.Initiators, models.Initiator{Type: models.InitiatorWeb})

	for _, test := range []struct {
		name string
		job  models.JobSpec
	}{
		{"mising public key", job1},
		{"mising min confirmations", job2},
		{"mising contract address", job3},
		{"single initiator", job4},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := services.ValidateJob(test.job, store, keyStore)
			assert.Error(t, err)
		})
	}
}
