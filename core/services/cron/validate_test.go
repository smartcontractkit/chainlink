package cron_test

import (
	"strings"
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

func TestValidatedCronJobSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
type            = "cron"
schemaVersion   = 1
schedule        = "CRON_TZ=UTC 0 0 1 1 * *"
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
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
			name: "no timezone",
			toml: `
type            = "cron"
schemaVersion   = 1
schedule        = "0 0 1 1 * *"
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), "cron schedule must specify a time zone using CRON_TZ"))
			},
		},
		{
			name: "invalid cron schedule",
			toml: `
type            = "cron"
schemaVersion   = 1
schedule        = "CRON_TZ=UTC x x"
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), "invalid cron schedule"))
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := cron.ValidatedCronSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
