package web_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/core/services/web"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

func TestValidateWebJobSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, os job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
			type            = "web"
			schemaVersion   = 1
			observationSource   = """
				ds          [type=http method=GET url="https://chain.link/ETH-USD"];
				ds_parse    [type=jsonparse path="data,price"];
				ds -> ds_parse;
			"""
			`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.WebSpec)
				b, err := jsonapi.Marshal(s.WebSpec)
				require.NoError(t, err)
				var r job.WebSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
			},
		},
		{
			name: "invalid job name",
			toml: `
			type            = "webjob"
			schemaVersion   = 1
			observationSource   = """
			    ds          [type=http method=GET url="https://chain.link/ETH-USD"];
			    ds_parse    [type=jsonparse path="data,price"];
			    ds_multiply [type=multiply times=100];
			    ds -> ds_parse -> ds_multiply;
			"""
			`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.Equal(t, "unsupported type webjob", err.Error())
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := web.ValidateWebSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
