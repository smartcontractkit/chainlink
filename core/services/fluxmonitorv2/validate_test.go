package fluxmonitorv2

import (
	"regexp"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/utils/tomlutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testcfg struct{}

func (testcfg) DefaultHTTPTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(2 * time.Second)
}

func TestValidate(t *testing.T) {
	t.Parallel()
	var tt = []struct {
		name      string
		toml      string
		config    ValidationConfig
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
type              = "fluxmonitor"
schemaVersion       = 1
name                = "example flux monitor spec"
contractAddress   = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
threshold = 0.5
absoluteThreshold = 0.0

idleTimerDisabled = true

pollTimerPeriod = "1m"
pollTimerDisabled = false

drumbeatEnabled = true
drumbeatSchedule = "@every 1m"
drumbeatRandomDelay = "10s"

minPayment = 1000000000000000000

observationSource = """
// data source 1
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds1_parse [type=jsonparse path="latest"];

// data source 2
ds2 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds2_parse [type=jsonparse path="latest"];

ds1 -> ds1_parse -> answer1;
ds2 -> ds2_parse -> answer1;

answer1 [type=median index=0];
"""
`,
			assertion: func(t *testing.T, j job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, j.FluxMonitorSpec)
				spec := j.FluxMonitorSpec
				assert.Equal(t, "example flux monitor spec", j.Name.String)
				assert.Equal(t, "fluxmonitor", j.Type.String())
				assert.Equal(t, uint32(1), j.SchemaVersion)
				assert.Equal(t, "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42", j.FluxMonitorSpec.ContractAddress.String())
				assert.Equal(t, tomlutils.Float32(0.5), spec.Threshold)
				assert.Equal(t, tomlutils.Float32(0), spec.AbsoluteThreshold)
				assert.Equal(t, true, spec.IdleTimerDisabled)
				assert.Equal(t, 1*time.Minute, spec.PollTimerPeriod)
				assert.Equal(t, false, spec.PollTimerDisabled)
				assert.Equal(t, true, spec.DrumbeatEnabled)
				assert.Equal(t, "@every 1m", spec.DrumbeatSchedule)
				assert.Equal(t, 10*time.Second, spec.DrumbeatRandomDelay)
				assert.Equal(t, false, spec.PollTimerDisabled)
				assert.Equal(t, assets.NewLinkFromJuels(1000000000000000000), spec.MinPayment)
				assert.NotZero(t, j.Pipeline)
			},
		},
		{
			name: "invalid contract addr",
			toml: `
type              = "fluxmonitor"
schemaVersion       = 1
name                = "example flux monitor spec"
contractAddress   = "0x3CCad4715152693fE3BC4460591e3D3Fbd071b42"
threshold = 0.5
absoluteThreshold = 0.0

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "1m"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}"];
ds1_parse [type=jsonparse path="latest"];
ds1 -> ds1_parse;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Nil(t, s.FluxMonitorSpec)
				require.Error(t, err)
				assert.Regexp(t, regexp.MustCompile("^.*is not a valid EIP55 formatted address$"), err.Error())
			},
		},
		{
			name: "invalid poll interval",
			toml: `
type              = "fluxmonitor"
schemaVersion       = 1
name                = "example flux monitor spec"
contractAddress   = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
maxTaskDuration = "1s"
threshold = 0.5
absoluteThreshold = 0.0

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "400ms"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}" timeout="500ms"];
ds1_parse [type=jsonparse path="latest"];
ds1 -> ds1_parse;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "PollTimerPeriod (400ms) must be equal or greater than the smallest value of MaxTaskDuration param, JobPipeline.HTTPRequest.DefaultTimeout config var, or MinTimeout of all tasks (500ms)")
			},
		},
		{
			name: "drumbeat and idle both active",
			toml: `
type              = "fluxmonitor"
schemaVersion       = 1
name                = "example flux monitor spec"
contractAddress   = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
maxTaskDuration = "1s"
threshold = 0.5
absoluteThreshold = 0.0

idleTimerDisabled = false
idleTimerPeriod = "1s"

drumbeatEnabled = true
drumbeatSchedule = "@every 1m"

pollTimerPeriod = "800ms"
pollTimerDisabled = false

observationSource = """
ds1 [type=http method=GET url="https://pricesource1.com" requestData="{\\"coin\\": \\"ETH\\", \\"market\\": \\"USD\\"}" timeout="500ms"];
ds1_parse [type=jsonparse path="latest"];
ds1 -> ds1_parse;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "When the drumbeat ticker is enabled, the idle timer must be disabled. Please set IdleTimerDisabled to true")
			},
		},
		{
			name: "integer thresholds",
			toml: `
type = "fluxmonitor"
schemaVersion = 1
name = "ADA / USD version 3 contract 0x3e4a23dB81D1F1268983f0CE78F1a9dC329A5b36 1624906849640"
contractAddress = "0x3e4a23dB81D1F1268983f0CE78F1a9dC329A5b36"
precision = 8
threshold = 2
idleTimerPeriod = "1m0s"
idleTimerDisabled = false
pollTimerPeriod = "1m0s"
pollTimerDisabled = false
maxTaskDuration = "0s"
observationSource = """
  // Node definitions.
  feed0 [method=POST name="bridge-coinmarketcap" requestData="{\\"data\\":{\\"from\\":\\"ADA\\",\\"to\\":\\"USD\\"}}" type=bridge];
  jsonparse0 [ path="data,result" type=jsonparse ];  
  feed0 -> jsonparse0;
  jsonparse0 -> median;
  feed1 [method=POST name="bridge-kaiko" requestData="{\\"data\\":{\\"from\\":\\"ADA\\",\\"to\\":\\"USD\\"}}" type=bridge];
  feed1 -> jsonparse1;
  jsonparse1 -> median;
  jsonparse1 [path="data,result" type=jsonparse];
  feed2 [method=POST name="bridge-nomics" requestData="{\\"data\\":{\\"from\\":\\"ADA\\",\\"to\\":\\"USD\\"}}" type=bridge];
  jsonparse2 [path="data,result" type=jsonparse];
  feed2 -> jsonparse2;
  jsonparse2 -> median;
  // Edge definitions.
  median [type=median];
  multiply0 [times=100000000 type=multiply];
  median -> multiply0;
"""
externalJobID = "cfa3fa6b-2850-446b-b973-8f4c3b29d519"
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidatedFluxMonitorSpec(testcfg{}, tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
