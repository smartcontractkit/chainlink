package fluxmonitorv2

import (
	"regexp"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/job"
	coreorm "github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	var tt = []struct {
		name       string
		toml       string
		setGlobals func(t *testing.T, c *coreorm.Config)
		assertion  func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
type              = "fluxmonitor"
schemaVersion       = 1
name                = "example flux monitor spec"
contractAddress   = "0x3cCad4715152693fE3BC4460591e3D3Fbd071b42"
precision = 2
threshold = 0.5
absoluteThreshold = 0.0 

idleTimerPeriod = "1s"
idleTimerDisabled = false

pollTimerPeriod = "1m"
pollTimerDisabled = false

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
				assert.Equal(t, int32(2), spec.Precision)
				assert.Equal(t, float32(0.5), spec.Threshold)
				assert.Equal(t, float32(0), spec.AbsoluteThreshold)
				assert.Equal(t, 1*time.Second, spec.IdleTimerPeriod)
				assert.Equal(t, false, spec.IdleTimerDisabled)
				assert.Equal(t, 1*time.Minute, spec.PollTimerPeriod)
				assert.Equal(t, false, spec.PollTimerDisabled)
				assert.Equal(t, assets.NewLink(1000000000000000000), spec.MinPayment)
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
precision = 2
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
precision = 2
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
				assert.EqualError(t, err, "pollTimer.period must be equal or greater than 500ms, got 400ms")
			},
			setGlobals: func(t *testing.T, c *coreorm.Config) {
				c.Set("DEFAULT_HTTP_TIMEOUT", "2s")
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := coreorm.NewConfig()
			if tc.setGlobals != nil {
				tc.setGlobals(t, c)
			}
			s, err := ValidatedFluxMonitorSpec(c, tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
