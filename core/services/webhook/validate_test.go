package webhook_test

import (
	"testing"

	"github.com/manyminds/api2go/jsonapi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	webhookmocks "github.com/smartcontractkit/chainlink/core/services/webhook/mocks"
)

func TestValidatedWebJobSpec(t *testing.T) {
	t.Parallel()
	var tt = []struct {
		name      string
		toml      string
		findError error
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
			findError: nil,
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
			findError: nil,
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
			findError: nil,
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
				require.Equal(t, webhook.ErrMissingJobID, errors.Cause(err))
			},
		},
		{
			name: "EI does not exist",
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
			findError: errors.New("foo"),
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
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			eim := new(webhookmocks.ExternalInitiatorManager)
			eim.On("FindExternalInitiatorByName", mock.Anything).Return(tc.findError)
			s, err := webhook.ValidatedWebhookSpec(tc.toml, eim)
			tc.assertion(t, s, err)
		})
	}
}
