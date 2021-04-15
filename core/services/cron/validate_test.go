package cron

import (
	"regexp"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

func TestValidateCronJobSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
type            = "cronjob"
schemaVersion   = 1
name            = "example cron spec"
cronSchedule 	= "0 0 0 1 1 *"
toAddress       = "0xc1912fEE45d61C87Cc5EA59DaE31190FFFFf232d"
oraclePayment 	= 1
observationSource   = """
    ds          [type=http method=GET url="https://chain.link/ETH-USD"];
    ds_parse    [type=jsonparse path="data,price"];
    ds_multiply [type=multiply times=100];
    ds_uint256  [type=ethuint256]
    ds -> ds_parse -> ds_multiply -> ds_uint256;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.CronSpec)
				b, err := jsonapi.Marshal(s.CronSpec)
				require.NoError(t, err)
				var r job.CronSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
			},
		},
		{
			name: "invalid cron schedule",
			toml: `
type            = "cronjob"
schemaVersion   = 1
name            = "invalid cron spec"
cronSchedule	= "x x"
toAddress       = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
oraclePayment 	= 1
observationSource   = """
    ds          [type=http method=GET url="https://chain.link/ETH-USD"];
    ds_parse    [type=jsonparse path="data,price"];
    ds_multiply [type=multiply times=100];
    ds_uint256  [type=ethuint256]
    ds -> ds_parse -> ds_multiply -> ds_uint256;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.Regexp(t, regexp.MustCompile("^.*error parsing cron schedule: Expected 5 to 6 fields, found 2: x x$"), err.Error())
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := ValidateCronSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
