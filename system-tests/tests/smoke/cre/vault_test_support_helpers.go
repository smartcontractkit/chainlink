package cre

import (
	"encoding/json"
	"fmt"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	ctflinkingservice "github.com/smartcontractkit/chainlink-testing-framework/framework/components/linkingservice"
	ctfvaultjwtissuer "github.com/smartcontractkit/chainlink-testing-framework/framework/components/vaultjwtissuer"
)

const (
	defaultVaultJWTIssuerKeyID = ctfvaultjwtissuer.DefaultJWTIssuerKeyID
	defaultVaultJWTAudience    = ctfvaultjwtissuer.DefaultAudience
)

type vaultTestJWTIssuer = ctfvaultjwtissuer.Client
type vaultJWTTokenClaims = ctfvaultjwtissuer.TokenClaims
type vaultTestLinkingService = ctflinkingservice.AdminClient

func newVaultDockerizedTestLinkingService(out *ctflinkingservice.Output) *vaultTestLinkingService {
	return ctflinkingservice.NewAdminClientFromOutput(out)
}

func newLocalVaultTestJWTIssuer() (*vaultTestJWTIssuer, error) {
	return ctfvaultjwtissuer.NewClient("http://127.0.0.1:18123", "http://127.0.0.1:18123")
}

func newVaultDockerizedTestJWTIssuer(localURL, dockerURL string) (*vaultTestJWTIssuer, error) {
	return ctfvaultjwtissuer.NewClient(localURL, dockerURL)
}

func outboundVaultRequestDigest(req jsonrpc.Request[json.RawMessage]) (string, error) {
	outboundReq := outboundRequestWithoutAuth(req)
	digest, err := outboundReq.Digest()
	if err != nil {
		return "", fmt.Errorf("failed to compute outbound request digest: %w", err)
	}
	return digest, nil
}
