package v2

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

func deriveOwnerHex(t *testing.T, tenantID uint64, orgID string) string {
	t.Helper()
	derived, err := pkgworkflows.GenerateWorkflowOwnerAddress(strconv.FormatUint(tenantID, 10), orgID)
	require.NoError(t, err)
	return hex.EncodeToString(derived)
}

const criticalMismatchMsg = "centralized workflow owner does not match owner derived from its organization ID: possible data corruption or malicious workflow registry"

const testTenantID uint64 = 1

func Test_verifyCentralizedOwnerOrgMapping(t *testing.T) {
	t.Parallel()

	const orgID = "org-123"
	matchingOwner := deriveOwnerHex(t, testTenantID, orgID)

	t.Run("matching owner passes and does not log critical", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)

		err := verifyCentralizedOwnerOrgMapping(lggr, "grpc:my-source:v1", matchingOwner, orgID, testTenantID)

		require.NoError(t, err)
		assert.Empty(t, logs.FilterMessage(criticalMismatchMsg).All())
	})

	t.Run("matching owner is case-insensitive and 0x tolerant", func(t *testing.T) {
		t.Parallel()
		lggr, _ := logger.TestObserved(t, zapcore.DebugLevel)

		err := verifyCentralizedOwnerOrgMapping(lggr, "grpc:my-source:v1", "0x"+matchingOwner, orgID, testTenantID)
		require.NoError(t, err)
	})

	t.Run("mismatched owner logs critical and rejects", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)

		victimEOA := "1234567890123456789012345678901234567890"
		err := verifyCentralizedOwnerOrgMapping(lggr, "grpc:malicious-source:v1", victimEOA, orgID, testTenantID)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrCentralizedOwnerOrgMismatch)
		entries := logs.FilterMessage(criticalMismatchMsg).All()
		require.Len(t, entries, 1)
		assert.GreaterOrEqual(t, entries[0].Level, zapcore.ErrorLevel)
	})

	t.Run("empty organization ID fails closed without deriving", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)

		err := verifyCentralizedOwnerOrgMapping(lggr, "grpc:my-source:v1", matchingOwner, "", testTenantID)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrCentralizedOwnerOrgMismatch)
		// It should not emit the mismatch critical (no derivation happened)...
		assert.Empty(t, logs.FilterMessage(criticalMismatchMsg).All())
		// ...but should log the distinct missing-org critical.
		assert.NotEmpty(t, logs.FilterMessageSnippet("did not provide an organization ID").All())
	})

	t.Run("configured tenantID is used for derivation", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)

		const tenantID uint64 = 42
		ownerForTenant42 := deriveOwnerHex(t, tenantID, orgID)

		require.Error(t, verifyCentralizedOwnerOrgMapping(lggr, "grpc:my-source:v1", matchingOwner, orgID, tenantID))
		require.NoError(t, verifyCentralizedOwnerOrgMapping(lggr, "grpc:my-source:v1", ownerForTenant42, orgID, tenantID))
		assert.Len(t, logs.FilterMessage(criticalMismatchMsg).All(), 1)
	})
}

func Test_tenantIDFromSettings(t *testing.T) {
	t.Parallel()

	t.Run("zero tenant ID is rejected", func(t *testing.T) {
		t.Parallel()

		_, err := tenantIDFromSettings(t.Context(), nil)
		require.ErrorIs(t, err, ErrTenantIDNotConfigured)
	})
}

func Test_isCentralizedWorkflowSource(t *testing.T) {
	t.Parallel()

	assert.True(t, isCentralizedWorkflowSource("grpc:my-source:v1"))
	// The local file source is a dev/test mechanism and is not subject to owner/org
	// verification.
	assert.False(t, isCentralizedWorkflowSource("file:local:v1"))
	assert.False(t, isCentralizedWorkflowSource("contract:123:0xabc"))
	assert.False(t, isCentralizedWorkflowSource(""))
}

func Test_maybeVerifyCentralizedOwnerOrgMapping(t *testing.T) {
	t.Parallel()

	const orgID = "org-123"
	victimEOA := "1234567890123456789012345678901234567890"
	lggr, _ := logger.TestObserved(t, zapcore.DebugLevel)
	gateOpen := limits.NewGateLimiter(true)

	t.Run("skips on-chain contract sources", func(t *testing.T) {
		t.Parallel()

		err := maybeVerifyCentralizedOwnerOrgMapping(
			t.Context(),
			lggr,
			"contract:123:0xregistry",
			victimEOA,
			orgID,
			gateOpen,
			nil,
		)
		require.NoError(t, err)
	})

	t.Run("skips when gate is closed for centralized sources", func(t *testing.T) {
		t.Parallel()

		err := maybeVerifyCentralizedOwnerOrgMapping(
			t.Context(),
			lggr,
			"grpc:my-source:v1",
			victimEOA,
			orgID,
			limits.NewGateLimiter(false),
			nil,
		)
		require.NoError(t, err)
	})
}
