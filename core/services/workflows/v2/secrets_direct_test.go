package v2

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/workflowkey"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type scriptedVaultCall struct {
	req      *vault.GetSecretsRequest
	metadata capabilities.RequestMetadata
}

// scriptedVault returns errs[i] (or a successful response once errs is exhausted) for call i.
type scriptedVault struct {
	metadataCapturingVault
	errs  []error
	calls []scriptedVaultCall
}

func (s *scriptedVault) Execute(ctx context.Context, req capabilities.CapabilityRequest) (capabilities.CapabilityResponse, error) {
	vr := &vault.GetSecretsRequest{}
	if err := req.Payload.UnmarshalTo(vr); err != nil {
		return capabilities.CapabilityResponse{}, err
	}
	s.calls = append(s.calls, scriptedVaultCall{req: vr, metadata: req.Metadata})
	if i := len(s.calls) - 1; i < len(s.errs) {
		return capabilities.CapabilityResponse{}, s.errs[i]
	}
	p, err := anypb.New(&vault.GetSecretsResponse{})
	if err != nil {
		return capabilities.CapabilityResponse{}, err
	}
	return capabilities.CapabilityResponse{Payload: p}, nil
}

type staticKeyFetcher struct{}

func (staticKeyFetcher) GetEncryptionKeys(context.Context) ([]string, error) {
	return []string{"key"}, nil
}

func TestSecretsFetcher_GetSecretsDirectly(t *testing.T) {
	t.Parallel()
	activeFrom2000, err := settings.NewTOMLGetter([]byte(`
[global.PerWorkflow]
FeatureVaultGetSecretsDirectlyActivePeriod = '[2000-01-01 00:00:00 +0000 UTC,2100-01-01 00:00:00 +0000 UTC]'
`))
	require.NoError(t, err)
	enabledAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	disabledAt := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)

	newFetcher := func(t *testing.T, v *scriptedVault, getter settings.Getter, ts time.Time) RawSecretsFetcher {
		lggr := logger.TestLogger(t)
		reg := registry.NewRegistry(lggr)
		require.NoError(t, reg.Add(t.Context(), v))
		return NewSecretsFetcher(MetricsLabelerTest(t), reg, lggr,
			limits.WorkflowResourcePoolLimiter[int](5), limits.NewUpperBoundLimiter[int](5), getter,
			"", "1234567890abcdef1234567890abcdef12345678", "workflowName", "workflowID", "workflowExecID", ts,
			workflowkey.MustNewXXXTestingOnly(big.NewInt(1)), nil)
	}
	req := &sdkpb.GetSecretsRequest{CallbackId: 7, Requests: []*sdkpb.SecretRequest{{Id: "a"}}}
	skew := fmt.Errorf("error executing request: %w", vault.ErrSecretVersionSkew)

	t.Run("disabled by default", func(t *testing.T) {
		t.Parallel()
		v := &scriptedVault{}
		_, err := newFetcher(t, v, nil, enabledAt).GetRawSecrets(t.Context(), req, staticKeyFetcher{})
		require.NoError(t, err)
		require.Len(t, v.calls, 1)
		assert.False(t, v.calls[0].req.GetSecretsDirectly)
	})

	t.Run("gated on the execution timestamp", func(t *testing.T) {
		t.Parallel()
		v := &scriptedVault{}
		_, err := newFetcher(t, v, activeFrom2000, disabledAt).GetRawSecrets(t.Context(), req, staticKeyFetcher{})
		require.NoError(t, err)
		assert.False(t, v.calls[0].req.GetSecretsDirectly)
	})

	t.Run("enabled", func(t *testing.T) {
		t.Parallel()
		v := &scriptedVault{}
		_, err := newFetcher(t, v, activeFrom2000, enabledAt).GetRawSecrets(t.Context(), req, staticKeyFetcher{})
		require.NoError(t, err)
		require.Len(t, v.calls, 1)
		assert.True(t, v.calls[0].req.GetSecretsDirectly)
	})

	t.Run("version skew is retried once under a new reference ID", func(t *testing.T) {
		t.Parallel()
		v := &scriptedVault{errs: []error{skew}}
		_, err := newFetcher(t, v, activeFrom2000, enabledAt).GetRawSecrets(t.Context(), req, staticKeyFetcher{})
		require.NoError(t, err)
		require.Len(t, v.calls, 2)
		assert.Equal(t, "7", v.calls[0].metadata.ReferenceID)
		assert.Equal(t, "7-retry1", v.calls[1].metadata.ReferenceID)
		assert.Equal(t, v.calls[0].metadata.WorkflowExecutionID, v.calls[1].metadata.WorkflowExecutionID)
		assert.True(t, v.calls[1].req.GetSecretsDirectly)
	})

	t.Run("only one retry", func(t *testing.T) {
		t.Parallel()
		v := &scriptedVault{errs: []error{skew, skew}}
		_, err := newFetcher(t, v, activeFrom2000, enabledAt).GetRawSecrets(t.Context(), req, staticKeyFetcher{})
		require.ErrorIs(t, err, vault.ErrSecretVersionSkew)
		assert.Len(t, v.calls, 2)
	})

	t.Run("other errors are not retried", func(t *testing.T) {
		t.Parallel()
		v := &scriptedVault{errs: []error{errors.New("boom")}}
		_, err := newFetcher(t, v, activeFrom2000, enabledAt).GetRawSecrets(t.Context(), req, staticKeyFetcher{})
		require.ErrorContains(t, err, "boom")
		assert.Len(t, v.calls, 1)
	})
}
