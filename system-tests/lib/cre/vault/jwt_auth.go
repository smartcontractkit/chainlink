package vault

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	vaultcap "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
)

const (
	DefaultJWTIssuerKeyID   = "vault-jwt-test-key"
	DefaultJWTAudience      = "https://vault.test.chain.link"
	DefaultJWTLifetime      = 5 * time.Minute
	defaultJWTPrivateKeyPEM = `-----BEGIN PRIVATE KEY-----
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
)

var (
	ErrMissingOrgID         = errors.New("org_id is required")
	ErrMissingRequestDigest = errors.New("request_digest is required")
	ErrMissingIssuer        = errors.New("issuer is required")
	ErrMissingKeyID         = errors.New("kid is required")
	ErrMissingPrivateKey    = errors.New("private key is required")
)

type jwtWebKey struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwtWebKeySet struct {
	Keys []jwtWebKey `json:"keys"`
}

// JWTTokenClaims describes the claims shape expected by Vault's JWT authorizer.
type JWTTokenClaims struct {
	OrgID         string
	WorkflowOwner string
	RequestDigest string
	// TenantID is emitted as urn:chainlink:tenant_id; when zero, Vault JWT tests default it to 1.
	TenantID uint64
	// Scopes are OAuth scopes (e.g. create:secrets) required for Vault JWT authorization.
	Scopes      []string
	Issuer      string
	Audience    string
	Subject     string
	JWTID       string
	KeyID       string
	IssuedAt    time.Time
	ExpiresAt   time.Time
	ExtraClaims map[string]any
}

// TestJWTIssuer is a minimal fake Auth0-style issuer for local CRE and system tests.
// It serves a JWKS endpoint and can mint RS256 JWTs matching Vault's expected claims.
type TestJWTIssuer struct {
	server       *http.Server
	listener     net.Listener
	signers      map[string]*rsa.PrivateKey
	defaultKeyID string
	mu           sync.RWMutex
}

// NewTestJWTIssuer creates a fake issuer with one generated RSA key and starts serving JWKS immediately.
func NewTestJWTIssuer() (*TestJWTIssuer, error) {
	return NewTestJWTIssuerOnAddr("0.0.0.0:0")
}

// NewTestJWTIssuerOnAddr creates a fake issuer bound to the provided TCP address.
func NewTestJWTIssuerOnAddr(listenAddr string) (*TestJWTIssuer, error) {
	privateKey, err := parseDefaultJWTSigningKey()
	if err != nil {
		return nil, err
	}

	return NewTestJWTIssuerWithKeysOnAddr(map[string]*rsa.PrivateKey{
		DefaultJWTIssuerKeyID: privateKey,
	}, DefaultJWTIssuerKeyID, listenAddr)
}

func parseDefaultJWTSigningKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(defaultJWTPrivateKeyPEM))
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

// NewTestJWTIssuerWithKeys creates a fake issuer backed by the provided key set.
func NewTestJWTIssuerWithKeys(signers map[string]*rsa.PrivateKey, defaultKeyID string) (*TestJWTIssuer, error) {
	return NewTestJWTIssuerWithKeysOnAddr(signers, defaultKeyID, "0.0.0.0:0")
}

// NewTestJWTIssuerWithKeysOnAddr creates a fake issuer backed by the provided key set and listen address.
func NewTestJWTIssuerWithKeysOnAddr(signers map[string]*rsa.PrivateKey, defaultKeyID, listenAddr string) (*TestJWTIssuer, error) {
	if len(signers) == 0 {
		return nil, errors.New("at least one signer is required")
	}
	if _, ok := signers[defaultKeyID]; !ok {
		return nil, fmt.Errorf("default signer %q is not present in key set", defaultKeyID)
	}
	if listenAddr == "" {
		listenAddr = "0.0.0.0:0"
	}

	// #nosec G102 -- test-only JWKS server must be reachable from Dockerized nodes during local CRE runs.
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate JWKS listener: %w", err)
	}

	issuer := &TestJWTIssuer{
		listener:     listener,
		signers:      signers,
		defaultKeyID: defaultKeyID,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", issuer.handleJWKS)
	mux.HandleFunc("/.well-known/openid-configuration", issuer.handleOpenIDConfiguration)

	issuer.server = &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		_ = issuer.server.Serve(listener)
	}()

	return issuer, nil
}

// Close shuts down the fake issuer.
func (i *TestJWTIssuer) Close() error {
	if i == nil || i.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return i.server.Shutdown(ctx)
}

// LocalIssuerURL returns an issuer URL that the local test process can use.
func (i *TestJWTIssuer) LocalIssuerURL() string {
	return i.baseURL("http://127.0.0.1")
}

// DockerIssuerURL returns an issuer URL that Dockerized Chainlink nodes can use.
func (i *TestJWTIssuer) DockerIssuerURL() string {
	return i.baseURL(framework.HostDockerInternal())
}

// MintToken signs a Vault-compatible JWT with one of the issuer's registered keys.
func (i *TestJWTIssuer) MintToken(claims JWTTokenClaims) (string, error) {
	if claims.KeyID == "" {
		claims.KeyID = i.defaultKeyID
	}
	if claims.Issuer == "" {
		claims.Issuer = i.LocalIssuerURL()
	}
	if claims.Audience == "" {
		claims.Audience = DefaultJWTAudience
	}

	i.mu.RLock()
	privateKey := i.signers[claims.KeyID]
	i.mu.RUnlock()
	if privateKey == nil {
		return "", fmt.Errorf("signer %q not registered in issuer", claims.KeyID)
	}

	return SignTestJWT(privateKey, claims)
}

func (i *TestJWTIssuer) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	keySet := jwtWebKeySet{
		Keys: make([]jwtWebKey, 0, len(i.signers)),
	}
	for keyID, privateKey := range i.signers {
		keySet.Keys = append(keySet.Keys, rsaPublicKeyToJWK(keyID, &privateKey.PublicKey))
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(keySet)
}

func (i *TestJWTIssuer) handleOpenIDConfiguration(w http.ResponseWriter, r *http.Request) {
	issuerURL := requestBaseURL(r)
	response := map[string]string{
		"issuer":   issuerURL,
		"jwks_uri": strings.TrimSuffix(issuerURL, "/") + "/.well-known/jwks.json",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (i *TestJWTIssuer) baseURL(host string) string {
	if i == nil || i.listener == nil {
		return ""
	}

	tcpAddr, ok := i.listener.Addr().(*net.TCPAddr)
	if !ok {
		return ""
	}

	return withPort(host, tcpAddr.Port)
}

// GenerateJWTSigningKey creates an RSA signing key suitable for RS256 JWTs.
func GenerateJWTSigningKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA signing key: %w", err)
	}
	return privateKey, nil
}

// SignTestJWT signs a Vault-compatible RS256 JWT.
func SignTestJWT(privateKey *rsa.PrivateKey, claims JWTTokenClaims) (string, error) {
	if privateKey == nil {
		return "", ErrMissingPrivateKey
	}
	if claims.KeyID == "" {
		return "", ErrMissingKeyID
	}
	if claims.Issuer == "" {
		return "", ErrMissingIssuer
	}
	if claims.OrgID == "" {
		return "", ErrMissingOrgID
	}
	if claims.RequestDigest == "" {
		return "", ErrMissingRequestDigest
	}

	now := time.Now().UTC()
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = now
	}
	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = claims.IssuedAt.Add(DefaultJWTLifetime)
	}
	if claims.Subject == "" {
		claims.Subject = claims.OrgID
	}
	if claims.Audience == "" {
		claims.Audience = DefaultJWTAudience
	}
	tenantID := claims.TenantID
	if tenantID == 0 {
		tenantID = 1
	}

	tokenClaims := jwt.MapClaims{
		"iss":                           claims.Issuer,
		"aud":                           claims.Audience,
		"sub":                           claims.Subject,
		"iat":                           jwt.NewNumericDate(claims.IssuedAt),
		"exp":                           jwt.NewNumericDate(claims.ExpiresAt),
		"org_id":                        claims.OrgID,
		vaultcap.ClaimChainlinkTenantID: strconv.FormatUint(tenantID, 10),
		vaultcap.ClaimVaultSecretManagementEnabled: "true",
		"authorization_details": []map[string]string{
			{
				"type":  "request_digest",
				"value": claims.RequestDigest,
			},
		},
	}

	if len(claims.Scopes) > 0 {
		tokenClaims["scope"] = strings.Join(claims.Scopes, " ")
	}

	if claims.WorkflowOwner != "" {
		tokenClaims["authorization_details"] = []map[string]string{
			{
				"type":  "request_digest",
				"value": claims.RequestDigest,
			},
			{
				"type":  "workflow_owner",
				"value": claims.WorkflowOwner,
			},
		}
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

// ComputeRequestDigest mirrors the digest computation used by Vault's authorizer.
func ComputeRequestDigest(req jsonrpc.Request[json.RawMessage]) (string, error) {
	return req.Digest()
}

// ComputeRawRequestDigest computes a Vault request digest from the marshalled JSON-RPC request body.
func ComputeRawRequestDigest(requestBody []byte) (string, error) {
	var req jsonrpc.Request[json.RawMessage]
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return "", fmt.Errorf("failed to decode JSON-RPC request: %w", err)
	}
	return ComputeRequestDigest(req)
}

func rsaPublicKeyToJWK(keyID string, publicKey *rsa.PublicKey) jwtWebKey {
	return jwtWebKey{
		Kid: keyID,
		Alg: "RS256",
		Kty: "RSA",
		Use: "sig",
		N:   base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes()),
	}
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}

	return withPort(scheme+"://"+r.Host, portFromRequest(r))
}

func portFromRequest(r *http.Request) int {
	if r.URL.Port() != "" {
		port, _ := strconv.Atoi(r.URL.Port())
		return port
	}
	if _, rawPort, err := net.SplitHostPort(r.Host); err == nil {
		port, _ := strconv.Atoi(rawPort)
		return port
	}
	return 0
}

func withPort(rawBase string, port int) string {
	base, err := url.Parse(rawBase)
	if err != nil {
		return rawBase
	}

	host := base.Hostname()
	if host == "" {
		host = strings.Trim(rawBase, "/")
	}

	if port > 0 {
		base.Host = net.JoinHostPort(host, strconv.Itoa(port))
	} else {
		base.Host = host
	}
	base.Path = "/"
	base.RawPath = ""
	base.RawQuery = ""
	base.Fragment = ""

	return base.String()
}
