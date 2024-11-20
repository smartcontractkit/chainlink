package actions

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/lib/client"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
)

// SetMockServerWithUSDCAttestation responds with a mock attestation for any msgHash
// The path is set with regex to match any path that starts with /v1/attestations
func SetMockServerWithUSDCAttestation(
	killGrave *ctftestenv.Killgrave,
	mockserver *ctfClient.MockserverClient,
) error {
	path := "/v1/attestations"
	response := struct {
		Status      string `json:"status"`
		Attestation string `json:"attestation"`
		Error       string `json:"error"`
	}{
		Status:      "complete",
		Attestation: "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b",
	}
	if killGrave == nil && mockserver == nil {
		return fmt.Errorf("both killgrave and mockserver are nil")
	}
	log.Info().Str("path", path).Msg("setting attestation-api response for any msgHash")
	if killGrave != nil {
		err := killGrave.SetAnyValueResponse(fmt.Sprintf("%s/{_hash:.*}", path), []string{http.MethodGet}, response)
		if err != nil {
			return fmt.Errorf("failed to set killgrave server value: %w", err)
		}
	}
	if mockserver != nil {
		err := mockserver.SetAnyValueResponse(fmt.Sprintf("%s/.*", path), response)
		if err != nil {
			return fmt.Errorf("failed to set mockserver value: %w URL = %s", err, fmt.Sprintf("%s/%s/.*", mockserver.LocalURL(), path))
		}
	}
	return nil
}
