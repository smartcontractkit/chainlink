package webhook_test

import (
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
)

func TestValidateWebJobSpec(t *testing.T) {
	var tt = []struct {
		name      string
		toml      string
		assertion func(t *testing.T, spec job.Job, err error)
	}{
		{
			name: "valid spec",
			toml: `
			type            = "webhook"
			schemaVersion   = 1
            jobID           = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
			observationSource   = """
				ds          [type=http method=GET url="https://chain.link/ETH-USD"];
				ds_parse    [type=jsonparse path="data,price"];
				ds -> ds_parse;
			"""
			`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.NoError(t, err)
				require.NotNil(t, s.WebhookSpec)
				b, err := jsonapi.Marshal(s.WebhookSpec)
				require.NoError(t, err)
				var r job.WebhookSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
				require.Equal(t, "0eec7e1dd0d2476ca1a872dfb6633f46", r.OnChainJobSpecID.String())
			},
		},
		{
			name: "invalid job name",
			toml: `
			type            = "webhookjob"
			schemaVersion   = 1
            jobID           = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
			observationSource   = """
			    ds          [type=http method=GET url="https://chain.link/ETH-USD"];
			    ds_parse    [type=jsonparse path="data,price"];
			    ds_multiply [type=multiply times=100];
			    ds -> ds_parse -> ds_multiply;
			"""
			`,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				assert.Equal(t, "unsupported type webhookjob", err.Error())
			},
		},
		{
			name: "missing jobID",
			toml: `
            type            = "webhookjob"
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
				require.Equal(t, webhook.ErrMissingJobID, errors.Cause(err))
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s, err := webhook.ValidateWebhookSpec(tc.toml)
			tc.assertion(t, s, err)
		})
	}
}
