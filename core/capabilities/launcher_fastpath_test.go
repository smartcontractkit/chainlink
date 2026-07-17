package capabilities

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func TestAggregatorFactoryForExecutableClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		capID       string
		method      string
		remoteDONF  int
		wantFactory bool
	}{
		{
			name:        "vault getSecrets returns factory",
			capID:       "vault@1.0.0",
			method:      vault.MethodGetSecrets,
			remoteDONF:  1,
			wantFactory: true,
		},
		{
			name:        "non-vault capability returns nil",
			capID:       "streams@1.0.0",
			method:      vault.MethodGetSecrets,
			remoteDONF:  1,
			wantFactory: false,
		},
		{
			name:        "vault non-getSecrets method returns nil",
			capID:       "vault@1.0.0",
			method:      "vault.secrets.create",
			remoteDONF:  1,
			wantFactory: false,
		},
		{
			name:        "vault id prefix but different capability returns nil",
			capID:       "vault2@1.0.0",
			method:      vault.MethodGetSecrets,
			remoteDONF:  1,
			wantFactory: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			factory := aggregatorFactoryForExecutableClient(tc.capID, tc.method, tc.remoteDONF)
			if tc.wantFactory {
				assert.NotNil(t, factory)
			} else {
				assert.Nil(t, factory)
			}
		})
	}
}
