package cre

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	ctflinkingservice "github.com/smartcontractkit/chainlink-testing-framework/framework/components/linkingservice"
	ctfvaultjwtissuer "github.com/smartcontractkit/chainlink-testing-framework/framework/components/vaultjwtissuer"
)

const (
	defaultVaultJWTIssuerKeyID = ctfvaultjwtissuer.DefaultJWTIssuerKeyID
	defaultVaultJWTAudience    = "https://vault.test.chain.link"
)

type vaultTestLinkingService struct {
	adminURL string
}

func newVaultDockerizedTestLinkingService(out *ctflinkingservice.Output) *vaultTestLinkingService {
	if out == nil {
		return nil
	}

	return &vaultTestLinkingService{
		adminURL: out.LocalAdminURL,
	}
}

func (s *vaultTestLinkingService) SetOwnerOrg(owner, orgID string) {
	if s == nil || s.adminURL == "" {
		panic("linking service admin URL is not configured")
	}

	payload, err := json.Marshal(map[string]string{
		"workflowOwner": owner,
		"orgID":         orgID,
	})
	if err != nil {
		panic(fmt.Errorf("failed to marshal linking service admin request: %w", err))
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		strings.TrimSuffix(s.adminURL, "/")+"/admin/link",
		bytes.NewReader(payload),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create linking service admin request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("failed to call linking service admin endpoint: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		panic(fmt.Errorf("unexpected linking service admin status: %d", resp.StatusCode))
	}
}

type vaultJWTTokenClaims struct {
	OrgID         string
	WorkflowOwner string
	RequestDigest string
	Issuer        string
	Audience      string
	Subject       string
	JWTID         string
	KeyID         string
	IssuedAt      time.Time
	ExpiresAt     time.Time
	ExtraClaims   map[string]any
}

type vaultTestJWTIssuer struct {
	localURL     string
	dockerURL    string
	privateKey   *rsa.PrivateKey
	defaultKeyID string
}

func newLocalVaultTestJWTIssuer() (*vaultTestJWTIssuer, error) {
	privateKey, err := parseDefaultVaultJWTSigningKey()
	if err != nil {
		return nil, err
	}

	return &vaultTestJWTIssuer{
		localURL:     normalizeVaultIssuerURL("http://127.0.0.1:18123"),
		dockerURL:    normalizeVaultIssuerURL("http://127.0.0.1:18123"),
		privateKey:   privateKey,
		defaultKeyID: defaultVaultJWTIssuerKeyID,
	}, nil
}

func newVaultDockerizedTestJWTIssuer(localURL, dockerURL string) (*vaultTestJWTIssuer, error) {
	privateKey, err := parseDefaultVaultJWTSigningKey()
	if err != nil {
		return nil, err
	}

	return &vaultTestJWTIssuer{
		localURL:     normalizeVaultIssuerURL(localURL),
		dockerURL:    normalizeVaultIssuerURL(dockerURL),
		privateKey:   privateKey,
		defaultKeyID: defaultVaultJWTIssuerKeyID,
	}, nil
}

func (i *vaultTestJWTIssuer) LocalIssuerURL() string {
	return i.localURL
}

func (i *vaultTestJWTIssuer) DockerIssuerURL() string {
	return i.dockerURL
}

func (i *vaultTestJWTIssuer) MintToken(claims vaultJWTTokenClaims) (string, error) {
	if i == nil || i.privateKey == nil {
		return "", errors.New("JWT issuer signing key is not configured")
	}
	if claims.KeyID == "" {
		claims.KeyID = i.defaultKeyID
	}
	if claims.Issuer == "" {
		claims.Issuer = i.LocalIssuerURL()
	}
	if claims.Audience == "" {
		claims.Audience = defaultVaultJWTAudience
	}

	return signVaultTestJWT(i.privateKey, claims)
}

func signVaultTestJWT(privateKey *rsa.PrivateKey, claims vaultJWTTokenClaims) (string, error) {
	if privateKey == nil {
		return "", errors.New("private key is required")
	}
	if claims.KeyID == "" {
		return "", errors.New("kid is required")
	}
	if claims.Issuer == "" {
		return "", errors.New("issuer is required")
	}
	if claims.OrgID == "" {
		return "", errors.New("org_id is required")
	}
	if claims.RequestDigest == "" {
		return "", errors.New("request_digest is required")
	}

	now := time.Now().UTC()
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = now
	}
	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = claims.IssuedAt.Add(5 * time.Minute)
	}
	if claims.Subject == "" {
		claims.Subject = claims.OrgID
	}
	if claims.Audience == "" {
		claims.Audience = defaultVaultJWTAudience
	}

	authorizationDetails := []map[string]string{
		{
			"type":  "request_digest",
			"value": claims.RequestDigest,
		},
	}
	if claims.WorkflowOwner != "" {
		authorizationDetails = append(authorizationDetails, map[string]string{
			"type":  "workflow_owner",
			"value": claims.WorkflowOwner,
		})
	}

	tokenClaims := jwt.MapClaims{
		"iss":                   claims.Issuer,
		"aud":                   claims.Audience,
		"sub":                   claims.Subject,
		"iat":                   jwt.NewNumericDate(claims.IssuedAt),
		"exp":                   jwt.NewNumericDate(claims.ExpiresAt),
		"org_id":                claims.OrgID,
		"authorization_details": authorizationDetails,
	}
	if claims.JWTID != "" {
		tokenClaims["jti"] = claims.JWTID
	}
	for key, value := range claims.ExtraClaims {
		tokenClaims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)
	token.Header["kid"] = claims.KeyID

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

func parseDefaultVaultJWTSigningKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(defaultVaultJWTPrivateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode default JWT signing key PEM")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default JWT signing key: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("default JWT signing key is not RSA")
	}

	return rsaKey, nil
}

func normalizeVaultIssuerURL(raw string) string {
	if raw == "" {
		return raw
	}
	if strings.HasSuffix(raw, "/") {
		return raw
	}
	return raw + "/"
}

func outboundVaultRequestDigest(req jsonrpc.Request[json.RawMessage]) (string, error) {
	outboundReq := outboundRequestWithoutAuth(req)
	return outboundReq.Digest()
}

const defaultVaultJWTPrivateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDYhEVPZ8YdC3Va
DGZ2hWPt+VYptOt0heTulBOwBW0ESavpfvokLYGFu+bLkGhIw365nCFw0eulLZYN
tD4nzq7F5Swtb2iIaDK19PBVNcukU/CY6j44KC1eomyaOvPXKWKwcc7qxjy9bIyA
TyOmOlxNxcNRSjL2SOApFkzb8M/RymHlMT/RY5ubytvjcbQgn2gy19U7HuNLYW1P
gviAAMY635u0A+HAxXx83lQSz9gy08/uBarmKAd2OadCA8cNiTSYyfUS6m1pycA7
j8ZHY75xL4hm+p2PJd9V1x3Z4S1TpZDIj+YAG/v4ZHB1vLTLoPIgwLEqwGRRWijl
sbdUZRd9AgMBAAECggEAGCiWFTiWheof43bLvgC/OC/gedHajctc0nQKSFMqqVZR
DMIixgOf1pyzMVaBFFFf4/T0VELQAMO34PqSDt4EaUdbaQxrxQCfW+cjI9bXTJQj
HeTRIXH2Mf98j67xQzo2bUqdlFufLmGcwbpS13rejrz4wKq/SfSyslLvK4FQpu8x
5J9ntn2wdgeUQCm62FyuNPxFMBldcovnwf9bbojTjMAatWfyF++W8OAcRqZCab1H
1WNPyhBqG5vDVMtgBdTkwZHqI01B+ozMnBLuEhsLVzvQWE79ZouWtU76GIeFlr0n
bC/3uWq9LBo1kEbLIPucxYA14ytWfpQwUvy1k11s4QKBgQD4dz2fVYSVb6hn0Pon
EQtunruNB7F2JlobY2s3C7aBKs+l48J16whKFcqHUA6NpuSvyUhFTqIpxM0LXdar
6nWu4Yw0kbqACJOHXuG71VhfkUgRJMOZoC/V0RKudoTwWDzFgNXvYF3bqtpmQDW7
2dUrSJ+jMOU7eCzXOdHDTFGhbQKBgQDfFQT/NACHapIn5w6c1Dha6fy7t1Z6A2zw
bUUzAh5C1kZ8yeDrkVfr5Ys+Y7Am/tfFteXO2XRSGH5yqq9YHVr0RihavqX72FGT
YY2rmyht+JjnZ3y+vOG5LXePR9tilvGei3jH0lTRPdwKpa6feHKry9MBx5xmqKqQ
xKRmyXaUUQKBgQCcOp3MqgEL1YGWhZhFKDp/+98B9mxnVgYiYojvu7Wt0jVuoZ+M
dZRowPrvyi7ccqwou+9tZNwiV1R2aTKqNmp44+k8xMT37GyXGdnmOWev77HY1b0H
w+lQEH4mpO9CELlllnTuZzGdBfj9gjJHQ9j9tlRqUDxTAGVxjzGOE1bgoQKBgQCu
DxmCAlIzVqzJY5hcN53tGcrvsKJRu2CBy9CFdy6jWctPzLipNROT5Nubh27HTmqP
QlkX50XCVIg88f60UttH44HTJBQgh+1GgIRolDycaa7sRyvnKzs4IEi8TAXaTAok
eZB44Rz60jhhOlsg5HscnoF6TwQyeYH0SOo5pRHXsQKBgQCY/pua7PceD5ZQ4lae
Pi5E9LzPjoeFegVgAP7bRUeC21nzLZlKYOcRCV2WkGLsz60bZm+7VEyFZmrrFoTE
58G0eCLCUq3Dj+NPfIvXNWwSuUAdDspWOBSCyENP+y+jLzIa2OtCj+KJe6Oe28pf
CcSeCJqr6aLeDRPcuD7yUat1OA==
-----END PRIVATE KEY-----`
