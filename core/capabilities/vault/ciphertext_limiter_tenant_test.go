package vault_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	pkgconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vault "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

// Regression tests for the pre-auth owner-scoped ciphertext limiter issue: checking
// MaxCiphertextLengthLimiter with an owner-scoped context registers a per-owner
// tenant (with a persistent background updater goroutine) in the limiter, so it
// must never be consulted before the request is authorized.

type recordedCiphertextCheck struct {
	owner  string
	amount pkgconfig.Size
}

// recordingCiphertextLimiter records every Check call's owner tenant and amount.
// errFor optionally returns an error for a given amount.
type recordingCiphertextLimiter struct {
	mu     sync.Mutex
	checks []recordedCiphertextCheck
	errFor func(amount pkgconfig.Size) error
}

var _ limits.BoundLimiter[pkgconfig.Size] = (*recordingCiphertextLimiter)(nil)

func (r *recordingCiphertextLimiter) Check(ctx context.Context, amount pkgconfig.Size) error {
	cre := contexts.CREValue(ctx)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks = append(r.checks, recordedCiphertextCheck{owner: cre.Owner, amount: amount})
	if r.errFor != nil {
		return r.errFor(amount)
	}
	return nil
}

func (r *recordingCiphertextLimiter) Limit(context.Context) (pkgconfig.Size, error) {
	return 0, nil
}

func (r *recordingCiphertextLimiter) Close() error { return nil }

func (r *recordingCiphertextLimiter) recorded() []recordedCiphertextCheck {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]recordedCiphertextCheck(nil), r.checks...)
}

func mustNewRecordingValidator(t *testing.T) (*vault.RequestValidator, *recordingCiphertextLimiter) {
	t.Helper()
	recorder := &recordingCiphertextLimiter{}
	validator := vault.NewRequestValidator(
		limits.NewUpperBoundLimiter(10),
		recorder,
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
		limits.NewUpperBoundLimiter[pkgconfig.Size](64*pkgconfig.Byte),
	)
	return validator, recorder
}

func mustWriteMethodParams(t *testing.T, method, requestID string, secrets []*vaultcommon.EncryptedSecret) *json.RawMessage {
	t.Helper()

	var payload any
	switch method {
	case vaulttypes.MethodSecretsCreate:
		payload = vaultcommon.CreateSecretsRequest{RequestId: requestID, EncryptedSecrets: secrets}
	case vaulttypes.MethodSecretsUpdate:
		payload = vaultcommon.UpdateSecretsRequest{RequestId: requestID, EncryptedSecrets: secrets}
	default:
		t.Fatalf("unsupported write method %s", method)
	}

	raw, err := json.Marshal(payload)
	require.NoError(t, err)
	return (*json.RawMessage)(&raw)
}

func mustWriteRequest(t *testing.T, method string, secrets []*vaultcommon.EncryptedSecret) jsonrpc.Request[json.RawMessage] {
	t.Helper()
	params := mustWriteMethodParams(t, method, "req-1", secrets)
	return jsonrpc.Request[json.RawMessage]{
		Version: jsonrpc.JsonRpcVersion,
		ID:      "req-1",
		Method:  method,
		Params:  params,
	}
}

func TestGatewayVaultRequestProcessor_ProcessRequest_UnauthorizedWriteNeverTouchesCiphertextLimiter(t *testing.T) {
	t.Parallel()

	for _, method := range []string{vaulttypes.MethodSecretsCreate, vaulttypes.MethodSecretsUpdate} {
		for _, stripOwnerPrefix := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/stripOwnerPrefix=%t", method, stripOwnerPrefix), func(t *testing.T) {
				t.Parallel()

				validator, recorder := mustNewRecordingValidator(t)

				// One-byte hex value under a fresh owner passes structure validation
				// (publicKey is nil so label validation is skipped), reaching authorization.
				secrets := []*vaultcommon.EncryptedSecret{
					{Id: &vaultcommon.SecretIdentifier{Owner: "0xnewowner", Key: "k"}, EncryptedValue: "00"},
				}
				req := mustWriteRequest(t, method, secrets)

				authorizer := vaultcapmocks.NewAuthorizer(t)
				authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(nil, errors.New("not authorized"))

				processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, stripOwnerPrefix)
				_, err := processor.ProcessRequest(t.Context(), &req, nil)
				require.Error(t, err)
				require.ErrorContains(t, err, "request not authorized")
				require.Empty(t, recorder.recorded(), "owner-scoped ciphertext limiter must not be consulted before authorization")
			})
		}
	}
}

func TestGatewayVaultRequestProcessor_ProcessRequest_AuthorizedWriteChecksCiphertextLimiterWithAuthorizedOwner(t *testing.T) {
	t.Parallel()

	for _, method := range []string{vaulttypes.MethodSecretsCreate, vaulttypes.MethodSecretsUpdate} {
		for _, stripOwnerPrefix := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/stripOwnerPrefix=%t", method, stripOwnerPrefix), func(t *testing.T) {
				t.Parallel()

				validator, recorder := mustNewRecordingValidator(t)
				owner := "0xauthorized"

				secrets := []*vaultcommon.EncryptedSecret{
					{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k1"}, EncryptedValue: "00"},
					{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k2"}, EncryptedValue: "abab"},
				}
				req := mustWriteRequest(t, method, secrets)

				authorizer := vaultcapmocks.NewAuthorizer(t)
				authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(vault.NewAuthResult("", owner, "digest", 0), nil)

				processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, stripOwnerPrefix)
				authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
				require.NoError(t, err)
				require.Equal(t, owner, authorized.AuthResult.AuthorizedOwner())

				expectedTenantOwner := contexts.CRE{Owner: owner}.Normalized().Owner
				checks := recorder.recorded()
				require.Len(t, checks, 2)
				for _, check := range checks {
					require.Equal(t, expectedTenantOwner, check.owner)
				}
				require.Equal(t, pkgconfig.Byte, checks[0].amount)
				require.Equal(t, 2*pkgconfig.Byte, checks[1].amount)
			})
		}
	}
}

func TestGatewayVaultRequestProcessor_ProcessRequest_AuthorizedWriteRejectsOversizedCiphertext(t *testing.T) {
	t.Parallel()

	for _, method := range []string{vaulttypes.MethodSecretsCreate, vaulttypes.MethodSecretsUpdate} {
		t.Run(method, func(t *testing.T) {
			t.Parallel()

			validator, recorder := mustNewRecordingValidator(t)
			recorder.errFor = func(amount pkgconfig.Size) error {
				return limits.ErrorBoundLimited[pkgconfig.Size]{Limit: pkgconfig.Byte, Amount: amount}
			}
			owner := "0xauthorized"

			secrets := []*vaultcommon.EncryptedSecret{
				{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k"}, EncryptedValue: "abab"},
			}
			req := mustWriteRequest(t, method, secrets)

			authorizer := vaultcapmocks.NewAuthorizer(t)
			authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(vault.NewAuthResult("", owner, "digest", 0), nil)

			processor := mustNewGatewayVaultRequestProcessor(t, validator, authorizer, false)
			_, err := processor.ProcessRequest(t.Context(), &req, nil)
			require.Error(t, err)
			require.True(t, vault.IsInvalidVaultParamsError(err))
			require.ErrorContains(t, err, "ciphertext size exceeds maximum allowed size")
			require.Len(t, recorder.recorded(), 1)
		})
	}
}

func TestRequestValidator_ValidateEncryptedSecretsStructure_SkipsCiphertextSizeCheck(t *testing.T) {
	t.Parallel()

	validator, recorder := mustNewRecordingValidator(t)

	oversized := strings.Repeat("00", 4096)
	err := validator.ValidateEncryptedSecretsStructure(t.Context(), nil, "req-1", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xabc", Key: "k"}, EncryptedValue: oversized},
	}, true)
	require.NoError(t, err)
	require.Empty(t, recorder.recorded())
}

func TestRequestValidator_ValidateEncryptedSecretsStructure_StillRejectsInvalidStructure(t *testing.T) {
	t.Parallel()

	validator, recorder := mustNewRecordingValidator(t)
	validSecret := &vaultcommon.EncryptedSecret{
		Id: &vaultcommon.SecretIdentifier{Owner: "0xabc", Key: "k"}, EncryptedValue: "00",
	}

	err := validator.ValidateEncryptedSecretsStructure(t.Context(), nil, "", []*vaultcommon.EncryptedSecret{validSecret}, true)
	require.ErrorContains(t, err, "request ID must not be empty")

	err = validator.ValidateEncryptedSecretsStructure(t.Context(), nil, "req-1", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "", Key: "k"}, EncryptedValue: "00"},
	}, true)
	require.ErrorContains(t, err, "owner cannot be empty")

	require.Empty(t, recorder.recorded())
}

func TestRequestValidator_ValidateCiphertextSizes_UsesAuthorizedOwnerPerItem(t *testing.T) {
	t.Parallel()

	validator, recorder := mustNewRecordingValidator(t)
	owner := "0xauthorized"

	err := validator.ValidateCiphertextSizes(t.Context(), owner, []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k1"}, EncryptedValue: "00"},
		{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k2"}, EncryptedValue: "abab"},
	})
	require.NoError(t, err)

	expectedTenantOwner := contexts.CRE{Owner: owner}.Normalized().Owner
	checks := recorder.recorded()
	require.Len(t, checks, 2)
	for _, check := range checks {
		require.Equal(t, expectedTenantOwner, check.owner)
	}
	require.Equal(t, pkgconfig.Byte, checks[0].amount)
	require.Equal(t, 2*pkgconfig.Byte, checks[1].amount)
}

func TestRequestValidator_ValidateCiphertextSizes_RejectsOversizedItemWithIndex(t *testing.T) {
	t.Parallel()

	validator, recorder := mustNewRecordingValidator(t)
	recorder.errFor = func(amount pkgconfig.Size) error {
		if amount > pkgconfig.Byte {
			return limits.ErrorBoundLimited[pkgconfig.Size]{Limit: pkgconfig.Byte, Amount: amount}
		}
		return nil
	}

	err := validator.ValidateCiphertextSizes(t.Context(), "0xauthorized", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xauthorized", Key: "k1"}, EncryptedValue: "00"},
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xauthorized", Key: "k2"}, EncryptedValue: "abab"},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "secret encrypted value at index 1 is invalid")
	require.ErrorContains(t, err, "ciphertext size exceeds maximum allowed size")
}

func TestRequestValidator_ValidateCiphertextSizes_RejectsNonHexValue(t *testing.T) {
	t.Parallel()

	validator, _ := mustNewRecordingValidator(t)

	err := validator.ValidateCiphertextSizes(t.Context(), "0xauthorized", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xauthorized", Key: "k"}, EncryptedValue: "zz"},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to decode encrypted value")
}

func TestRequestValidator_ValidateCiphertextSizes_RejectsNilItem(t *testing.T) {
	t.Parallel()

	validator, _ := mustNewRecordingValidator(t)

	err := validator.ValidateCiphertextSizes(t.Context(), "0xauthorized", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xauthorized", Key: "k1"}, EncryptedValue: "00"},
		nil,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "encrypted secret must not be nil at index 1")
}
