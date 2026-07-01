package standardcapabilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseTenantID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config string
		want   uint64
		wantOK bool
	}{
		{name: "empty config", config: "", want: 0, wantOK: false},
		{name: "no auth0", config: `{"schedule":"1s"}`, want: 0, wantOK: false},
		{name: "auth0 without tenantID", config: `{"auth0":{}}`, want: 0, wantOK: false},
		{name: "auth0 zero tenantID", config: `{"auth0":{"tenantID":0}}`, want: 0, wantOK: false},
		{name: "auth0 tenantID set", config: `{"auth0":{"tenantID":1}}`, want: 1, wantOK: true},
		{name: "auth0 tenantID with other fields", config: `{"foo":"bar","auth0":{"tenantID":42}}`, want: 42, wantOK: true},
		{name: "invalid json", config: `{not json`, want: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := parseTenantID(tt.config)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}
