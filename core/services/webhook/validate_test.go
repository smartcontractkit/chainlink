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
	"github.com/smartcontractkit/chainlink/core/store/models"
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
            externalJobID           = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
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
				require.Equal(t, "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46", s.ExternalJobID.String())
			},
		},
		{
			name: "invalid job name",
			toml: `
			type            = "webhookjob"
			schemaVersion   = 1
            externalJobID           = "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46"
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
            externalJobID           = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
            externalInitiatorName = "foo"
            externalInitiatorSpec = "foo"
            observationSource   = """
                ds          [type=http method=GET url="https://chain.link/ETH-USD"];
                ds_parse    [type=jsonparse path="data,price"];
                ds -> ds_parse;
            """
            `,
			findError: errors.New("foo"),
			assertion: func(t *testing.T, s job.Job, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "EI does exist",
			toml: `
            type            = "webhook"
            schemaVersion   = 1
            externalJobID = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
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
				require.Equal(t, "0eec7e1d-d0d2-476c-a1a8-72dfb6633f46", s.ExternalJobID.String())
				var r job.WebhookSpec
				err = jsonapi.Unmarshal(b, &r)
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			eim := new(webhookmocks.ExternalInitiatorManager)
			eim.On("FindExternalInitiatorByName", mock.Anything).Return(models.ExternalInitiator{}, tc.findError)
			s, err := webhook.ValidatedWebhookSpec(tc.toml, eim)
			tc.assertion(t, s, err)
		})
	}
}
