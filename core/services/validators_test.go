package services_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"

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
	tests := []struct {
		name  string
		input []byte
		want  error
	}{
		{"base case", cltest.MustReadFile(t, "testdata/hello_world_job.json"), nil},
		{
			"error in job",
			cltest.MustReadFile(t, "testdata/invalid_endat_job.json"),
			models.NewJSONAPIErrorsWith("StartAt cannot be before EndAt"),
		},
		{
			"error in runat initr",
			cltest.MustReadFile(t, "testdata/run_at_wo_time_job.json"),
			models.NewJSONAPIErrorsWith("RunAt must have a time"),
		},
		{
			"error in task",
			cltest.MustReadFile(t, "testdata/nonexistent_task_job.json"),
			models.NewJSONAPIErrorsWith("idonotexist is not a supported adapter type"),
		},
		{
			"zero initiators",
			cltest.MustReadFile(t, "testdata/zero_initiators.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one initiator only",
			cltest.MustReadFile(t, "testdata/initiator_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"one task only",
			cltest.MustReadFile(t, "testdata/task_only_job.json"),
			models.NewJSONAPIErrorsWith("Must have at least one Initiator and one Task"),
		},
		{
			"runlog and ethtx with an address",
			cltest.MustReadFile(t, "testdata/runlog_ethtx_w_address_job.json"),
			models.NewJSONAPIErrorsWith("Cannot set EthTx Task's address parameter with a RunLog Initiator"),
		},
		{
			"runlog and ethtx with a function selector",
			cltest.MustReadFile(t, "testdata/runlog_ethtx_w_funcselector_job.json"),
			models.NewJSONAPIErrorsWith("Cannot set EthTx Task's function selector parameter with a RunLog Initiator"),
		},
		{
			"runlog and ethtx with a fromAddress that doesn't match one of our keys",
			cltest.MustReadFile(t, "testdata/runlog_ethtx_w_missing_fromAddress_job.json"),
			models.NewJSONAPIErrorsWith("Cannot set EthTx Task's fromAddress parameter: the node does not have this private key in the database"),
		},
		{
			"runlog with two ethtx tasks",
			cltest.MustReadFile(t, "testdata/runlog_2_ethlogs_job.json"),
			models.NewJSONAPIErrorsWith("Cannot RunLog initiated jobs cannot have more than one EthTx Task"),
		},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var j models.JobSpec
			assert.NoError(t, json.Unmarshal(test.input, &j))
			result := services.ValidateJob(j, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateJob_RejectsSleepAdapterWhenExperimentalAdaptersAreDisabled(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	sleepingJob := cltest.NewJobWithWebInitiator()
	sleepingJob.Tasks[0].Type = adapters.TaskTypeSleep

	store.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", true)
	assert.NoError(t, services.ValidateJob(sleepingJob, store))

	store.Config.Set("ENABLE_EXPERIMENTAL_ADAPTERS", false)
	assert.Error(t, services.ValidateJob(sleepingJob, store))
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
	err := store.KeyStore.Unlock("password")
	assert.NoError(t, err)
	_, err = store.KeyStore.NewAccount() // matches correct_password.txt
	assert.NoError(t, err)

	oracles := []string{cltest.DefaultKeyAddress.Hex()}

	basic := string(cltest.MustReadFile(t, "testdata/hello_world_agreement.json"))
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
					JobSpecID:       models.NewID(),
					Type:            models.InitiatorServiceAgreementExecutionLog,
					InitiatorParams: models.InitiatorParams{},
				}, {
					JobSpecID:       models.NewID(),
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

			result := services.ValidateServiceAgreement(sa, store)

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

func TestValidateOracleSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os offchainreporting.OracleSpec, err error)
	}{
		{
			name: "decodes valid oracle spec toml",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xaA07d525B4006a2f927D79CA78a23A8ee680A32A"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1, int(os.SchemaVersion))
				assert.False(t, os.IsBootstrapPeer)
			},
		},
		{
			name: "decodes bootstrap toml",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = true
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1, int(os.SchemaVersion))
				assert.True(t, os.IsBootstrapPeer)
			},
		},
		{
			name: "raises error on extra keys",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = true
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xaA07d525B4006a2f927D79CA78a23A8ee680A32A"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="one,two"];
ds1_multiply [type=multiply times=1.23];
ds1 -> ds1_parse -> ds1_multiply -> answer1;
answer1      [type=median index=0];
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "unrecognised key for bootstrap peer: keyBundleID")
				assert.Contains(t, err.Error(), "unrecognised key for bootstrap peer: transmitterAddress")
				assert.Contains(t, err.Error(), "unrecognised key for bootstrap peer: observationTimeout")
				assert.Contains(t, err.Error(), "unrecognised key for bootstrap peer: observationSource")
			},
		},
		{
			name: "empty pipeline string non-bootstrap node",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = false
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid dot",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = false
observationSource = """
->
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "invalid peer address",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/invalid/peer/address"]
isBootstrapPeer    = false
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero timeouts",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
blockchainTimeout  = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "non-zero intervals",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
isBootstrapPeer    = false
contractConfigTrackerSubscribeInterval = "0s"
observationSource = """
blah
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "broken monitoring endpoint",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = true
monitoringEndpoint = "\t/fd\2ff )(*&^%$#@"
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.EqualError(t, err, "(8, 23): invalid escape sequence: \\2")
			},
		},
		{
			name: "max task duration > observation timeout should error",
			toml: `
type               = "offchainreporting"
maxTaskDuration    = "30s"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xaA07d525B4006a2f927D79CA78a23A8ee680A32A"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout];
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "max task duration must be < observation timeout")
			},
		},
		{
			name: "individual max task duration > observation timeout should error",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = [
"/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "73e8966a78ca09bb912e9565cfb79fbe8a6048fab1f0cf49b18047c3895e0447"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "0xaA07d525B4006a2f927D79CA78a23A8ee680A32A"
observationTimeout = "10s"
observationSource = """
ds1          [type=bridge name=voter_turnout timeout="30s"];
"""
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "individual max task duration must be < observation timeout")
			},
		},
		{
			name: "sane defaults",
			toml: `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq"
p2pBootstrapPeers  = []
isBootstrapPeer    = true
`,
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.NoError(t, err)
				assert.Equal(t, os.ContractConfigConfirmations, uint16(3))
				assert.Equal(t, os.ObservationTimeout, models.Interval(10*time.Second))
				assert.Equal(t, os.BlockchainTimeout, models.Interval(20*time.Second))
				assert.Equal(t, os.ContractConfigTrackerSubscribeInterval, models.Interval(2*time.Minute))
				assert.Equal(t, os.ContractConfigTrackerPollInterval, models.Interval(1*time.Minute))
				assert.Len(t, os.P2PBootstrapPeers, 0)
			},
		},
		{
			name: "toml parse doesn't panic",
			toml: string(cltest.MustHexDecodeString("2222220d5c22223b22225c0d21222222")),
			assertion: func(t *testing.T, os offchainreporting.OracleSpec, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := services.ValidatedOracleSpecToml(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
