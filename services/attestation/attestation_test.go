// +build sgx_enclave

package attestation

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AttestationReport struct {
	Report struct {
		Body struct {
			ReportData []byte `json:"report_data"`
			MrEnclave  []byte `json:"mr_enclave"`
		} `json:"body"`
		KeyID []byte `json:"key_id"`
		Mac   []byte `json:"mac"`
	} `json:"report"`
}

func TestReport(t *testing.T) {
	result, err := Report()
	assert.NoError(t, err)

	var report AttestationReport
	err = json.Unmarshal([]byte(result), &report)
	assert.NoError(t, err)

	// Report now contains a nonce so we can only assert on its structure
	assert.Len(t, report.Report.Body.ReportData, 64)
	assert.Len(t, report.Report.Body.MrEnclave, 32)
	assert.Len(t, report.Report.KeyID, 32)
	assert.Len(t, report.Report.Mac, 16)
}
