package ccip

import (
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink-ccip/execute/tokendata/lbtc"
	ctftestenv "github.com/smartcontractkit/chainlink-testing-framework/lib/docker/test_env"
	"github.com/smartcontractkit/chainlink-testing-framework/parrot"
)

const (
	usdcAttestationWildcardPath = "/v1/attestations/*"
	lbtcAttestationPath         = "/bridge/v1/deposits/getByHash"
)

// SetMockServerWithUSDCAttestation registers Circle CCTP-style JSON on the Parrot mock server.
// The CCIP USDC client issues GET {AttestationAPI}/v1/attestations/{messageHash}.
func SetMockServerWithUSDCAttestation(p *ctftestenv.Parrot, isUSDCAttestationMissing bool) error {
	if p == nil || p.Client == nil {
		return fmt.Errorf("parrot mock adapter is not initialized")
	}
	// Re-registering the same route fails on Parrot; allow repeated setup.
	_ = p.Client.DeleteRoute(&parrot.Route{Method: http.MethodGet, Path: usdcAttestationWildcardPath})

	body := `{
			"status": "complete",
			"attestation": "0x9049623e91719ef2aa63c55f357be2529b0e7122ae552c18aff8db58b4633c4d3920ff03d3a6d1ddf11f06bf64d7fd60d45447ac81f527ba628877dc5ca759651b08ffae25a6d3b1411749765244f0a1c131cbfe04430d687a2e12fd9d2e6dc08e118ad95d94ad832332cf3c4f7a4f3da0baa803b7be024b02db81951c0f0714de1b"
		}`
	if isUSDCAttestationMissing {
		body = `{
				"status": "pending",
				"error": "internal error"
			}`
	}
	return p.Client.RegisterRoute(&parrot.Route{
		Method:             http.MethodGet,
		Path:               usdcAttestationWildcardPath,
		RawResponseBody:    body,
		ResponseStatusCode: http.StatusOK,
	})
}

// SetMockServerWithLBTCAttestation registers Lombard LBTC attestation API responses on the Parrot mock server.
// The CCIP LBTC client issues POST {AttestationAPI}/bridge/v1/deposits/getByHash with JSON body {"messageHash":[...]}.
func SetMockServerWithLBTCAttestation(p *ctftestenv.Parrot, isAttestationMissing bool) error {
	if p == nil || p.Client == nil {
		return fmt.Errorf("parrot mock adapter is not initialized")
	}
	_ = p.Client.DeleteRoute(&parrot.Route{Method: http.MethodPost, Path: lbtcAttestationPath})

	var response lbtc.AttestationResponse
	if isAttestationMissing {
		response = lbtc.AttestationResponse{
			Code:    3,
			Message: "invalid hash",
		}
	} else {
		response = lbtc.AttestationResponse{
			Attestations: []lbtc.Attestation{
				{
					MessageHash: "0xdee9d5a70c34ab6ad3d3be55cc81b8f3dbd7aaf4070d7f1046b239e4995df489",
					Status:      "NOTARIZATION_STATUS_SESSION_APPROVED",
					Data:        "0x0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000000e45c70a5050000000000000000000000000000000000000000000000000000000000000061000000000000000000000000ca571682d1478ab3f7fcbcbade6e4954de3a96760000000000000000000000000000000000000000000000000000000000014a34000000000000000000000000ca571682d1478ab3f7fcbcbade6e4954de3a96760000000000000000000000004b431813bcf797bf9bf93890656618ac80a1d5d20000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000024000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000000000000000000000000000000000000000000040fd53ff0dd6da6873e12afe8ac0b4e2c1c92ac5edf940ba53cf2a1ae2f70dbf4a7bbd6b5949b2bb511d1cbfd3e90ebb12dd6bf20074a3c5b67732f63571363d6b000000000000000000000000000000000000000000000000000000000000004094aa83e1524340ed3365b6ef061cb337c593ace76ca9565b984a8695f7292edf2aa55673ed153fe3282c18bfab6383fcdc23f96fefb0246264d6f12769cf34b0000000000000000000000000000000000000000000000000000000000000004052a309783debf3682b377c309e105fb288d0acf7aae352ea02b306cd11506aee7f418fb1a13284c9262243d69120d5064f1c442f652c4f03b4ff0071f7e5923a00000000000000000000000000000000000000000000000000000000000000406dd9501ab5af88098f2443634c5196c5ceddfab27bb109d7cd8d464dfe0c86bf36d5dad799a9c755fb30ff00aaee4eabeb8cbc2380e3903f260d24833aa26a51",
				},
			},
		}
	}
	return p.Client.RegisterRoute(&parrot.Route{
		Method:             http.MethodPost,
		Path:               lbtcAttestationPath,
		ResponseBody:       response,
		ResponseStatusCode: http.StatusOK,
	})
}
