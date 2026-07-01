package v2

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pkgworkflows "github.com/smartcontractkit/chainlink-common/pkg/workflows"
)

func Test_isCentralizedWorkflowSource(t *testing.T) {
	t.Parallel()

	assert.True(t, isCentralizedWorkflowSource("grpc:my-source:v1"))
	assert.True(t, isCentralizedWorkflowSource("file:my-source:v1"))
	assert.False(t, isCentralizedWorkflowSource("contract:1:0xabc"))
	assert.False(t, isCentralizedWorkflowSource(""))
}

// deriveOwnerHex returns the hex-encoded (no 0x prefix) workflow owner address that
// the centralized owner/org verification expects for the given org ID.
func deriveOwnerHex(t *testing.T, orgID string) string {
	t.Helper()
	derived, err := pkgworkflows.GenerateWorkflowOwnerAddress(strconv.FormatUint(defaultTenantID, 10), orgID)
	require.NoError(t, err)
	return hex.EncodeToString(derived)
}

const criticalMismatchMsg = "centralized workflow owner does not match owner derived from its organization ID: possible data corruption or malicious workflow registry"

func Test_verifyCentralizedOwnerOrgMapping(t *testing.T) {
	t.Parallel()

	const orgID = "org-123"
	matchingOwner := deriveOwnerHex(t, orgID)

	t.Run("matching owner passes and does not log critical", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
		h := &eventHandler{lggr: lggr}

		err := h.verifyCentralizedOwnerOrgMapping("grpc:my-source:v1", matchingOwner, orgID)

		require.NoError(t, err)
		assert.Empty(t, logs.FilterMessage(criticalMismatchMsg).All())
	})

	t.Run("matching owner is case-insensitive and 0x tolerant", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
		h := &eventHandler{lggr: lggr}

		err := h.verifyCentralizedOwnerOrgMapping("grpc:my-source:v1", "0x"+matchingOwner, orgID)

		require.NoError(t, err)
		assert.Empty(t, logs.FilterMessage(criticalMismatchMsg).All())
	})

	t.Run("mismatched owner logs critical and rejects (spoofed on-chain owner)", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
		h := &eventHandler{lggr: lggr}

		// An attacker-claimed on-chain EOA that is not derived from orgID.
		victimEOA := "1234567890123456789012345678901234567890"
		err := h.verifyCentralizedOwnerOrgMapping("grpc:malicious-source:v1", victimEOA, orgID)

		require.Error(t, err)
		require.ErrorIs(t, err, ErrCentralizedOwnerOrgMismatch)

		entries := logs.FilterMessage(criticalMismatchMsg).All()
		require.Len(t, entries, 1)
		// Critical is a remapping of DPanic, but falls back to Error ('[crit]' prefix)
		// on loggers that don't support it; accept either.
		assert.GreaterOrEqual(t, entries[0].Level, zapcore.ErrorLevel)
	})

	t.Run("empty orgID skips verification (no error, no critical)", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
		h := &eventHandler{lggr: lggr}

		err := h.verifyCentralizedOwnerOrgMapping("grpc:my-source:v1", matchingOwner, "")

		require.NoError(t, err)
		assert.Empty(t, logs.FilterMessage(criticalMismatchMsg).All())
	})

	t.Run("configured tenantID is used for derivation", func(t *testing.T) {
		t.Parallel()
		lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)

		const tenantID uint64 = 42
		h := &eventHandler{lggr: lggr}
		h.SetTenantID(tenantID)

		derived, err := pkgworkflows.GenerateWorkflowOwnerAddress(strconv.FormatUint(tenantID, 10), orgID)
		require.NoError(t, err)
		ownerForTenant42 := hex.EncodeToString(derived)

		// The default-tenant owner must be rejected under tenant 42...
		require.Error(t, h.verifyCentralizedOwnerOrgMapping("grpc:my-source:v1", matchingOwner, orgID))
		// ...while the tenant-42-derived owner passes.
		require.NoError(t, h.verifyCentralizedOwnerOrgMapping("grpc:my-source:v1", ownerForTenant42, orgID))

		assert.Len(t, logs.FilterMessage(criticalMismatchMsg).All(), 1)
	})
}
