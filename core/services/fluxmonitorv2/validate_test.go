package fluxmonitorv2

import (
	"regexp"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	var tt = []struct {
		name       string
		toml       string
		setGlobals func(t *testing.T, c *orm.Config)
		assertion  func(t *testing.T, os job.SpecDB, err error)
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
			assertion: func(t *testing.T, s job.SpecDB, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.FluxMonitorSpec)
				b, err := jsonapi.Marshal(s.FluxMonitorSpec)
				require.NoError(t, err)
				var r job.FluxMonitorSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
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
ds1 -> ds1_parse -> answer1;
"""
`,
			assertion: func(t *testing.T, s job.SpecDB, err error) {
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
			assertion: func(t *testing.T, s job.SpecDB, err error) {
				require.Error(t, err)
				assert.EqualError(t, err, "pollTimer.period must be equal or greater than 500ms, got 400ms")
			},
			setGlobals: func(t *testing.T, c *orm.Config) {
				c.Set("DEFAULT_HTTP_TIMEOUT", "2s")
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := orm.NewConfig()
			if tc.setGlobals != nil {
				tc.setGlobals(t, c)
			}
			s, err := ValidatedFluxMonitorSpec(c, tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
