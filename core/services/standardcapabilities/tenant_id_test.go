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
		{name: "no tenantID", config: `{"schedule":"1s"}`, want: 0, wantOK: false},
		{name: "zero tenantID", config: `{"tenantID":0}`, want: 0, wantOK: false},
		{name: "tenantID set", config: `{"tenantID":1}`, want: 1, wantOK: true},
		{name: "tenantID with other fields", config: `{"foo":"bar","tenantID":42}`, want: 42, wantOK: true},
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
