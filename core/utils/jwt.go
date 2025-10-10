package utils

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

const (
	maxJWTExpiryDuration = 5 * time.Minute // Maximum allowed expiry duration
)

// Option is a function type that allows configuring CreateRequestJWT.
type Option func(*jwtOptions)

type jwtOptions struct {
	expiryDuration *time.Duration
	issuer         *string  // New field for optional issuer
	audience       []string // New field for optional audience
	subject        *string  // New field for optional subject
}

// VerifyOption is a function type that allows configuring VerifyRequestJWT.
type VerifyOption func(*verifyOptions)

type verifyOptions struct {
	maxExpiryDuration *time.Duration
}

func WithExpiry(d time.Duration) Option {
	return func(opts *jwtOptions) {
		opts.expiryDuration = &d
	}
}

func WithIssuer(issuer string) Option {
	return func(opts *jwtOptions) {
		opts.issuer = &issuer
	}
}

func WithAudience(audience []string) Option {
	return func(opts *jwtOptions) {
		opts.audience = audience
	}
}

func WithSubject(subject string) Option {
	return func(opts *jwtOptions) {
		opts.subject = &subject
	}
}

func WithMaxExpiryDuration(d time.Duration) VerifyOption {
	return func(opts *verifyOptions) {
		opts.maxExpiryDuration = &d
	}
}

type SigningMethodEth struct{}

var EthereumSigningMethod = &SigningMethodEth{}

func init() {
	// registering a custom implementation of the ETH signing method here
	jwt.RegisterSigningMethod(EthereumSigningMethod.Alg(), func() jwt.SigningMethod {
		return EthereumSigningMethod
	})
}

func (m *SigningMethodEth) Alg() string {
	return "ETH"
}

// Sign signs the given signing string using the given key
// key is expected to be an *ecdsa.PrivateKey
// returns the signature as a 65-byte array (r, s, v) with v being 0 or 1
func (m *SigningMethodEth) Sign(signingString string, key any) ([]byte, error) {
	var ecdsaKey *ecdsa.PrivateKey
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		ecdsaKey = k
	default:
		return nil, jwt.ErrInvalidKeyType
	}
	signature, err := GenerateEthSignature(ecdsaKey, []byte(signingString))
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Verify verifies the given signature for the given signing string using the given public key
// key is expected to be a gethcommon.Address
func (m *SigningMethodEth) Verify(signingString string, signature []byte, key any) error {
	var ethAddr gethcommon.Address
	switch k := key.(type) {
	case gethcommon.Address:
		ethAddr = k
	default:
		return jwt.ErrInvalidKeyType
	}
	recoveredAddr, err := GetSignersEthAddress([]byte(signingString), signature)
	if err != nil {
		return err
	}
	if !bytes.Equal(recoveredAddr.Bytes(), ethAddr.Bytes()) {
		return jwt.ErrSignatureInvalid
	}
	return nil
}

type JWTClaims struct {
	Digest string `json:"digest"`
	jwt.RegisteredClaims
}

// CreateRequestJWT creates an unsigned JWT for a JSON-RPC request
// JWT has 3 parts: header, payload, and signature as shown below
// header:
//
//	{
//		alg: "ES256K",
//		typ: "JWT"
//	}
//
// payload:
//
//	{
//		digest: "<request-digest>",      // 32 byte hex string with "0x" prefix
//		jti: "<unique-id>",              // JWT ID (UUID) for replay protection (REQUIRED)
//		iss: "ethereum-address",         // Ethereum address of the issuer
//		exp: <timestamp>,                // expiration time (Unix timestamp)
//		iat: <timestamp>                 // issued at time (Unix timestamp)
//	}
//
// sample payload:
//
//	{
//	  "digest": "0x4a1f2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a",
//	  "jti": "550e8400-e29b-41d4-a716-446655440000",
//	  "iss": "0xc40B35B10c5003C182e300a9a34F6ff559eB746d",
//	  "exp": 1717596700,
//	  "iat": 1717596400
//	}
//
// signature: ETH signature of the header and payload using the private key
func CreateRequestJWT[T any](req jsonrpc.Request[T], opts ...Option) (*jwt.Token, error) {
	// Apply options
	options := &jwtOptions{}
	for _, opt := range opts {
		opt(options)
	}

	expiryDuration := maxJWTExpiryDuration
	if options.expiryDuration != nil {
		expiryDuration = *options.expiryDuration
	}

	digest, err := req.Digest()
	if err != nil {
		return nil, err
	}

	var issuer string
	if options.issuer != nil {
		issuer = *options.issuer
	}

	var subject string
	if options.subject != nil {
		subject = *options.subject
	}

	var audience []string
	if options.audience != nil {
		audience = options.audience
	}

	now := time.Now()
	jti := uuid.New().String()

	claims := JWTClaims{
		Digest: "0x" + digest,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    issuer,
			Subject:   subject,
			Audience:  jwt.ClaimStrings(audience),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	return jwt.NewWithClaims(&SigningMethodEth{}, claims), nil
}

func splitToken(tokenString string) (string, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", "", errors.New("invalid JWT format: expected 3 parts")
	}
	signedString := parts[0] + "." + parts[1]
	signature := parts[2]
	return signedString, signature, nil
}

// VerifyRequestJWT verifies a signed JWT for a JSON-RPC request
// It recovers and returns the public key used to sign the JWT, checks the issuer, validates the digest,
// and performs all validations done by jwt.ParseWithClaims() including expiration checks.
func VerifyRequestJWT[T any](tokenString string, req jsonrpc.Request[T], opts ...VerifyOption) (*JWTClaims, gethcommon.Address, error) {
	options := &verifyOptions{}
	for _, opt := range opts {
		opt(options)
	}

	maxExpiryDuration := maxJWTExpiryDuration
	if options.maxExpiryDuration != nil {
		maxExpiryDuration = *options.maxExpiryDuration
	}
	signedString, signature, err := splitToken(tokenString)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}
	decodedSignature, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}
	pubKey, err := GetSignersEthAddress([]byte(signedString), decodedSignature)
	if err != nil {
		return nil, gethcommon.Address{}, err
	}
	verifiedToken, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != EthereumSigningMethod.Alg() {
			return nil, jwt.ErrSignatureInvalid
		}
		if _, ok := token.Method.(*SigningMethodEth); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, gethcommon.Address{}, err
	}
	verifiedClaims, ok := verifiedToken.Claims.(*JWTClaims)
	if !ok {
		return nil, gethcommon.Address{}, jwt.ErrTokenInvalidClaims
	}
	if !verifiedToken.Valid {
		return nil, gethcommon.Address{}, jwt.ErrTokenInvalidClaims
	}
	reqDigest, err := req.Digest()
	if err != nil {
		return nil, gethcommon.Address{}, err
	}
	if verifiedClaims.Digest != "0x"+reqDigest {
		return nil, gethcommon.Address{}, errors.New("JWT digest does not match request digest")
	}
	if verifiedClaims.ID == "" {
		return nil, gethcommon.Address{}, errors.New("JWT ID (jti) is required but missing")
	}
	if verifiedClaims.ExpiresAt == nil {
		return nil, gethcommon.Address{}, errors.New("expiredAt (exp) is required but missing")
	}
	if verifiedClaims.IssuedAt == nil {
		return nil, gethcommon.Address{}, errors.New("issuedAt (iat) is required but missing")
	}
	duration := verifiedClaims.ExpiresAt.Sub(verifiedClaims.IssuedAt.Time)
	if duration > maxExpiryDuration {
		return nil, gethcommon.Address{}, fmt.Errorf("expiry duration exceeds maximum allowed %.0f minutes", maxExpiryDuration.Minutes())
	}

	return verifiedClaims, pubKey, nil
}
