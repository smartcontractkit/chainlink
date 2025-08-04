package utils

import (
	"encoding/json"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
)

func TestSigningMethodEth_Sign(t *testing.T) {
	t.Run("valid ECDSA key", func(t *testing.T) {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		sm := &SigningMethodEth{}
		signingString := "test.signing.string"

		signature, err := sm.Sign(signingString, privateKey)
		require.NoError(t, err)
		require.Len(t, signature, 65)
	})

	t.Run("invalid key type", func(t *testing.T) {
		sm := &SigningMethodEth{}
		signingString := "test.signing.string"

		signature, err := sm.Sign(signingString, "invalid-key")
		require.Nil(t, signature)
		require.Equal(t, jwt.ErrInvalidKeyType, err)
	})
}

func TestSigningMethodEth_Verify(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	t.Run("valid signature and address", func(t *testing.T) {
		sm := &SigningMethodEth{}
		signingString := "test.signing.string"

		signature, err := sm.Sign(signingString, privateKey)
		require.NoError(t, err)

		err = sm.Verify(signingString, signature, address)
		require.NoError(t, err)
	})

	t.Run("wrong address", func(t *testing.T) {
		sm := &SigningMethodEth{}
		signingString := "test.signing.string"

		signature, err := sm.Sign(signingString, privateKey)
		require.NoError(t, err)

		wrongAddress := gethcommon.HexToAddress("0x3E1Cc9C5aBEa1B87f9C918A7160A37a182f27e8b")
		err = sm.Verify(signingString, signature, wrongAddress)
		require.Equal(t, jwt.ErrSignatureInvalid, err)
	})

	t.Run("invalid key type", func(t *testing.T) {
		sm := &SigningMethodEth{}
		signingString := "test.signing.string"
		signature := make([]byte, 65)

		err := sm.Verify(signingString, signature, "invalid-key")
		require.Equal(t, jwt.ErrInvalidKeyType, err)
	})

	t.Run("invalid signature", func(t *testing.T) {
		sm := &SigningMethodEth{}
		signingString := "test.signing.string"
		invalidSignature := make([]byte, 64) // Wrong length

		err := sm.Verify(signingString, invalidSignature, address)
		require.Error(t, err)
	})
}

func TestWithExpiry(t *testing.T) {
	duration := 30 * time.Minute
	opts := &jwtOptions{}

	WithExpiry(duration)(opts)

	require.NotNil(t, opts.expiryDuration)
	require.Equal(t, duration, *opts.expiryDuration)
}

func TestWithIssuer(t *testing.T) {
	issuer := "test-issuer"
	opts := &jwtOptions{}

	WithIssuer(issuer)(opts)

	require.NotNil(t, opts.issuer)
	require.Equal(t, issuer, *opts.issuer)
}

func TestWithAudience(t *testing.T) {
	audience := []string{"aud1", "aud2"}
	opts := &jwtOptions{}

	WithAudience(audience)(opts)

	require.Equal(t, audience, opts.audience)
}

func TestWithSubject(t *testing.T) {
	subject := "test-subject"
	opts := &jwtOptions{}

	WithSubject(subject)(opts)

	require.NotNil(t, opts.subject)
	require.Equal(t, subject, *opts.subject)
}

func testRequest(t *testing.T) jsonrpc.Request[json.RawMessage] {
	params := json.RawMessage(`{"num":3}`)
	return jsonrpc.Request[json.RawMessage]{
		Version: "2.0",
		ID:      "test-id",
		Method:  "test-method",
		Params:  &params,
	}
}

func TestCreateRequestJWT(t *testing.T) {
	req := testRequest(t)

	t.Run("default options", func(t *testing.T) {
		token, err := CreateRequestJWT(req)
		require.NoError(t, err)
		require.NotNil(t, token)

		claims, ok := token.Claims.(JWTClaims)
		require.True(t, ok)
		digest, err := req.Digest()
		require.NoError(t, err)
		require.Equal(t, "0x"+digest, claims.Digest)
		require.Empty(t, claims.Issuer)
		require.Empty(t, claims.Subject)
		require.Empty(t, claims.Audience)
		require.NotNil(t, claims.ExpiresAt)
		require.NotNil(t, claims.IssuedAt)
	})

	t.Run("with custom options", func(t *testing.T) {
		customExpiry := 2 * time.Hour
		issuer := "test-issuer"
		audience := []string{"aud1", "aud2"}
		subject := "test-subject"

		token, err := CreateRequestJWT(req,
			WithExpiry(customExpiry),
			WithIssuer(issuer),
			WithAudience(audience),
			WithSubject(subject),
		)
		require.NoError(t, err)

		claims, ok := token.Claims.(JWTClaims)
		require.True(t, ok)
		require.Equal(t, issuer, claims.Issuer)
		require.Equal(t, subject, claims.Subject)
		require.Equal(t, jwt.ClaimStrings(audience), claims.Audience)

		// Verify expiry duration
		expectedExpiry := claims.IssuedAt.Add(customExpiry)
		require.Equal(t, expectedExpiry.Unix(), claims.ExpiresAt.Unix())
	})
}

func TestSplitToken(t *testing.T) {
	t.Run("valid token format", func(t *testing.T) {
		tokenString := "header.payload.signature"

		signedString, signature, err := splitToken(tokenString)
		require.NoError(t, err)
		require.Equal(t, "header.payload", signedString)
		require.Equal(t, "signature", signature)
	})

	t.Run("invalid token format - too few parts", func(t *testing.T) {
		tokenString := "header.payload"

		_, _, err := splitToken(tokenString)
		require.EqualError(t, err, "invalid JWT format: expected 3 parts")
	})

	t.Run("invalid token format - too many parts", func(t *testing.T) {
		tokenString := "header.payload.signature.extra"

		_, _, err := splitToken(tokenString)
		require.EqualError(t, err, "invalid JWT format: expected 3 parts")
	})
}

func TestVerifyRequestJWT_Integration(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	req := testRequest(t)

	t.Run("valid JWT verification", func(t *testing.T) {
		token, err := CreateRequestJWT(req)
		require.NoError(t, err)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		claims, recoveredAddr, err := VerifyRequestJWT(tokenString, req)
		require.NoError(t, err)

		expectedAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
		require.Equal(t, expectedAddr, recoveredAddr)
		digest, err := req.Digest()
		require.NoError(t, err)
		require.Equal(t, "0x"+digest, claims.Digest)
	})

	t.Run("invalid token format", func(t *testing.T) {
		invalidToken := "invalid.token"

		_, _, err := VerifyRequestJWT(invalidToken, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid JWT format")
	})

	t.Run("digest mismatch", func(t *testing.T) {
		now := time.Now()
		claims := JWTClaims{
			Digest: "0x123", // different digest
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(defaultJWTExpiryDuration)),
				IssuedAt:  jwt.NewNumericDate(now),
			},
		}

		token := jwt.NewWithClaims(&SigningMethodEth{}, claims)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		_, _, err = VerifyRequestJWT(tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "JWT digest does not match request digest")
	})

	t.Run("expired token", func(t *testing.T) {
		// Create token with past expiry
		token, err := CreateRequestJWT(req, WithExpiry(-time.Hour))
		require.NoError(t, err)

		tokenString, err := token.SignedString(privateKey)
		require.NoError(t, err)

		_, _, err = VerifyRequestJWT(tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "token is expired")
	})

	t.Run("wrong signing method", func(t *testing.T) {
		digest, err := req.Digest()
		require.NoError(t, err)
		claims := JWTClaims{
			Digest: "0x" + digest,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("secret"))
		require.NoError(t, err)

		_, _, err = VerifyRequestJWT(tokenString, req)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid signature: signature length must be 65 bytes")
	})
}

func TestSigningMethodRegistration(t *testing.T) {
	method := jwt.GetSigningMethod("ETH")
	require.NotNil(t, method)

	ethMethod, ok := method.(*SigningMethodEth)
	require.True(t, ok)
	require.Equal(t, "ETH", ethMethod.Alg())
}
