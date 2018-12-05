// +build sgx_enclave

package attestation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	result, err := Report()
	assert.NoError(t, err)
	assert.Equal(t, `{"report":{"key_id":[255,0,255,0,255,0,255,0,255,0,255,0,255,0,255,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"mac":[167,99,107,85,185,122,149,93,161,30,126,255,36,101,174,248]}}`, result)
}
