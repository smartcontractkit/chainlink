package vault_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	vault "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault"
	vaultcapmocks "github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// Regression tests for the pre-auth owner-scoped ciphertext limiter issue: checking
// MaxCiphertextLengthLimiter with an owner-scoped context registers a per-owner
// tenant (with a persistent background updater goroutine) in the limiter, so it
// must never be consulted before the request is authorized.

type recordedCiphertextCheck struct {
	org    string
	owner  string
	amount commonconfig.Size
}

// recordingCiphertextLimiter records every Check call's owner tenant and amount.
// errFor optionally returns an error for a given amount.
type recordingCiphertextLimiter struct {
	mu     sync.Mutex
	checks []recordedCiphertextCheck
	errFor func(amount commonconfig.Size) error
}

var _ limits.BoundLimiter[commonconfig.Size] = (*recordingCiphertextLimiter)(nil)

func (r *recordingCiphertextLimiter) Check(ctx context.Context, amount commonconfig.Size) error {
	cre := contexts.CREValue(ctx)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks = append(r.checks, recordedCiphertextCheck{org: cre.Org, owner: cre.Owner, amount: amount})
	if r.errFor != nil {
		return r.errFor(amount)
	}
	return nil
}

func (r *recordingCiphertextLimiter) Limit(context.Context) (commonconfig.Size, error) {
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
		limits.NewUpperBoundLimiter[commonconfig.Size](64*commonconfig.Byte),
		limits.NewUpperBoundLimiter[commonconfig.Size](64*commonconfig.Byte),
		limits.NewUpperBoundLimiter[commonconfig.Size](64*commonconfig.Byte),
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
				require.Equal(t, commonconfig.Byte, checks[0].amount)
				require.Equal(t, 2*commonconfig.Byte, checks[1].amount)
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
			recorder.errFor = func(amount commonconfig.Size) error {
				return limits.ErrorBoundLimited[commonconfig.Size]{Limit: commonconfig.Byte, Amount: amount}
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

	err := validator.ValidateCiphertextSizes(t.Context(), "", owner, []*vaultcommon.EncryptedSecret{
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
	require.Equal(t, commonconfig.Byte, checks[0].amount)
	require.Equal(t, 2*commonconfig.Byte, checks[1].amount)
}

func TestRequestValidator_ValidateCiphertextSizes_RejectsOversizedItemWithIndex(t *testing.T) {
	t.Parallel()

	validator, recorder := mustNewRecordingValidator(t)
	recorder.errFor = func(amount commonconfig.Size) error {
		if amount > commonconfig.Byte {
			return limits.ErrorBoundLimited[commonconfig.Size]{Limit: commonconfig.Byte, Amount: amount}
		}
		return nil
	}

	err := validator.ValidateCiphertextSizes(t.Context(), "", "0xauthorized", []*vaultcommon.EncryptedSecret{
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

	err := validator.ValidateCiphertextSizes(t.Context(), "", "0xauthorized", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xauthorized", Key: "k"}, EncryptedValue: "zz"},
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to decode encrypted value")
}

func TestRequestValidator_ValidateCiphertextSizes_RejectsNilItem(t *testing.T) {
	t.Parallel()

	validator, _ := mustNewRecordingValidator(t)

	err := validator.ValidateCiphertextSizes(t.Context(), "", "0xauthorized", []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xauthorized", Key: "k1"}, EncryptedValue: "00"},
		nil,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "encrypted secret must not be nil at index 1")
}

// stubOrgResolver maps owners to orgs from a static table, or fails with err if set.
type stubOrgResolver struct {
	services.Service
	orgByOwner map[string]string
	err        error
	calls      atomic.Int32
}

func (s *stubOrgResolver) Get(_ context.Context, owner string) (string, error) {
	s.calls.Add(1)
	if s.err != nil {
		return "", s.err
	}
	return s.orgByOwner[owner], nil
}

func mustNewGatewayVaultRequestProcessorWithOrgResolver(t *testing.T, validator *vault.RequestValidator, authorizer vault.Authorizer, orgResolver orgresolver.OrgResolver) *vault.GatewayVaultRequestProcessor {
	t.Helper()
	processor, err := vault.NewGatewayVaultRequestProcessor(validator, authorizer, orgResolver, false, logger.TestLogger(t))
	require.NoError(t, err)
	return processor
}

func TestGatewayVaultRequestProcessor_ProcessRequest_CiphertextLimiterOrg(t *testing.T) {
	t.Parallel()

	const owner = "0xauthorized"
	tests := []struct {
		name        string
		jwtOrgID    string
		orgResolver *stubOrgResolver
		expectedOrg string
	}{
		{name: "JWT org claim wins over resolver", jwtOrgID: "org-jwt", orgResolver: &stubOrgResolver{orgByOwner: map[string]string{owner: "org-resolved"}}, expectedOrg: "org-jwt"},
		{name: "resolver used without JWT org claim", orgResolver: &stubOrgResolver{orgByOwner: map[string]string{owner: "org-resolved"}}, expectedOrg: "org-resolved"},
		{name: "resolver error continues without org", orgResolver: &stubOrgResolver{err: errors.New("linking service down")}, expectedOrg: ""},
		{name: "nil resolver continues without org", expectedOrg: ""},
	}
	for _, tc := range tests {
		for _, method := range []string{vaulttypes.MethodSecretsCreate, vaulttypes.MethodSecretsUpdate} {
			t.Run(tc.name+"/"+method, func(t *testing.T) {
				t.Parallel()

				validator, recorder := mustNewRecordingValidator(t)
				req := mustWriteRequest(t, method, []*vaultcommon.EncryptedSecret{
					{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k"}, EncryptedValue: "00"},
				})

				authorizer := vaultcapmocks.NewAuthorizer(t)
				authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(vault.NewAuthResult(tc.jwtOrgID, owner, "digest", 0), nil)

				var resolver orgresolver.OrgResolver
				if tc.orgResolver != nil {
					resolver = tc.orgResolver
				}
				processor := mustNewGatewayVaultRequestProcessorWithOrgResolver(t, validator, authorizer, resolver)
				authorized, err := processor.ProcessRequest(t.Context(), &req, nil)
				require.NoError(t, err)
				require.Equal(t, tc.expectedOrg, authorized.OrgID)

				checks := recorder.recorded()
				require.Len(t, checks, 1)
				require.Equal(t, tc.expectedOrg, checks[0].org)
			})
		}
	}
}

func TestGatewayVaultRequestProcessor_ProcessRequest_UnauthorizedWriteNeverResolvesOrg(t *testing.T) {
	t.Parallel()

	validator, _ := mustNewRecordingValidator(t)
	req := mustWriteRequest(t, vaulttypes.MethodSecretsCreate, []*vaultcommon.EncryptedSecret{
		{Id: &vaultcommon.SecretIdentifier{Owner: "0xnewowner", Key: "k"}, EncryptedValue: "00"},
	})

	authorizer := vaultcapmocks.NewAuthorizer(t)
	authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(nil, errors.New("not authorized"))

	orgResolver := &stubOrgResolver{orgByOwner: map[string]string{"0xnewowner": "org-1"}}
	processor := mustNewGatewayVaultRequestProcessorWithOrgResolver(t, validator, authorizer, orgResolver)
	_, err := processor.ProcessRequest(t.Context(), &req, nil)
	require.ErrorContains(t, err, "request not authorized")
	require.Zero(t, orgResolver.calls.Load(), "org must not be resolved before authorization")
}

// TestGatewayVaultRequestProcessor_ProcessRequest_CiphertextLimitOrgOverride proves end-to-end
// that an org-level settings override of the owner-scoped VaultCiphertextSizeLimit is applied
// once the org of the authorized owner is resolved.
func TestGatewayVaultRequestProcessor_ProcessRequest_CiphertextLimitOrgOverride(t *testing.T) {
	t.Parallel()

	const (
		raisedOwner = "0x1111111111111111111111111111111111aaaa"
		normalOwner = "0x2222222222222222222222222222222222bbbb"
	)
	getter, err := settings.NewJSONGetter([]byte(`{
		"global": {"PerOwner": {"VaultCiphertextSizeLimit": "2kb"}},
		"org": {"org-raised": {"PerOwner": {"VaultCiphertextSizeLimit": "5kb"}}}
	}`))
	require.NoError(t, err)
	validator, err := vault.NewRequestValidatorFromLimitsFactory(limits.Factory{Settings: getter})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, validator.Close()) })

	orgResolver := &stubOrgResolver{orgByOwner: map[string]string{raisedOwner: "org-raised", normalOwner: "org-normal"}}
	threeKB := strings.Repeat("00", 3*1024)

	for owner, wantErr := range map[string]bool{raisedOwner: false, normalOwner: true} {
		req := mustWriteRequest(t, vaulttypes.MethodSecretsCreate, []*vaultcommon.EncryptedSecret{
			{Id: &vaultcommon.SecretIdentifier{Owner: owner, Key: "k"}, EncryptedValue: threeKB},
		})
		authorizer := vaultcapmocks.NewAuthorizer(t)
		authorizer.EXPECT().AuthorizeRequest(t.Context(), mock.Anything).Return(vault.NewAuthResult("", owner, "digest", 0), nil)

		processor := mustNewGatewayVaultRequestProcessorWithOrgResolver(t, validator, authorizer, orgResolver)
		_, err := processor.ProcessRequest(t.Context(), &req, nil)
		if wantErr {
			require.ErrorContains(t, err, "ciphertext size exceeds maximum allowed size", owner)
		} else {
			require.NoError(t, err, owner)
		}
	}
}
