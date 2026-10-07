package localcapmgr

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

// configCapture records the configJSON each capability is launched with, so the test can assert
// the effective (offchain-wins-with-fallback) config without a real plugin.
type configCapture struct {
	mu   sync.Mutex
	last map[string]string // capID -> configJSON
}

func (c *configCapture) builder(_ context.Context, capID string, _ uint32, _ string, configJSON string, _ *ocrtypes.ContractConfig) ([]job.ServiceCtx, error) {
	c.mu.Lock()
	c.last[capID] = configJSON
	c.mu.Unlock()
	return []job.ServiceCtx{&mockService{}}, nil
}

func (c *configCapture) get(capID string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last[capID]
}

// validateCapRegistrySpec runs a capabilities_registry cresettings TOML (config_type and the
// payload in the settings TOML) through the real ingestion validator, returning the validated job
// (with the computed Hash).
func validateCapRegistrySpec(t *testing.T, raw string) job.Job {
	t.Helper()
	settings := fmt.Sprintf("config_type = \"capabilities_registry\"\noffchain_config = '''%s'''\n", raw)
	toml := fmt.Sprintf("type = \"cresettings\"\nschemaVersion = 1\nexternalJobID = %q\nsettings = \"\"\"\n%s\"\"\"\n", uuid.NewString(), strings.ReplaceAll(settings, `\`, `\\`))
	jb, err := cresettings.ValidatedCRESettingsSpec(toml)
	require.NoError(t, err)
	require.NotNil(t, jb.CRESettingsSpec)
	return jb
}

// applyAsDelegate mirrors what the cresettings delegate does when a capabilities_registry job
// starts: it stores the validated payload into the runtime GlobalConfig. This test exercises
// validation + GlobalConfig + the launcher.
func applyAsDelegate(t *testing.T, gc *globalconfig.GlobalConfig, raw string) {
	t.Helper()
	jb := validateCapRegistrySpec(t, raw)
	require.NoError(t, gc.Store(globalconfig.Update{Raw: raw, Hash: jb.CRESettingsSpec.Hash}))
}

// TestOffchainRegistry_EndToEnd walks the full node-side path with the cutover ON:
// validate a capabilities_registry spec -> apply it to the runtime GlobalConfig (as the
// delegate would) -> the LocalCapabilityManager launches the capability with offchain-wins
// config -> a newer payload re-reconciles reactively via Subscribe, without a new Reconcile call.
func TestOffchainRegistry_EndToEnd(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	lggr := testLogger(t)

	gc := globalconfig.New()

	// --- Ingestion: validate the v1 payload and apply it as the delegate would. ---
	rawV1, err := marshalOffchainRegistry(offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		1: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})},
	}))
	require.NoError(t, err)
	applyAsDelegate(t, gc, rawV1)

	_, version := gc.Load()
	require.Equal(t, uint64(1), version)

	// --- Launcher wired to the same GlobalConfig, cutover ON. ---
	capture := &configCapture{last: map[string]string{}}
	localCfg := &testLocalCapabilities{
		allowlisted: map[string]bool{"cron@1.0.0": true},
		useOffchain: true,
	}
	m, err := NewLocalCapabilityManager(lggr, localCfg, capture.builder, gc, true)
	require.NoError(t, err)
	require.NoError(t, m.Start(ctx))
	t.Cleanup(func() { require.NoError(t, m.Close()) })

	// On-chain DON provides interval=20 and region=us. Offchain overrides interval->30;
	// region has no offchain value, so it falls back to the on-chain value.
	dons := []registry.DON{{
		ID:   1,
		Name: testDONName(1),
		CapabilityConfigurations: map[string]registry.CapabilityConfiguration{
			"cron@1.0.0": onchainSpecConfig(t, map[string]any{"interval": "20", "region": "us"}),
		},
	}}
	require.NoError(t, m.Reconcile(ctx, dons))

	// Offchain wins on interval; region falls back to on-chain. (json.Marshal sorts map keys.)
	assert.JSONEq(t, `{"interval":"30","region":"us"}`, capture.get("cron@1.0.0"))

	// Cross-validation telemetry sees the capability as matched, and flags config_mismatch
	// because turning the gate on changed the launch config.
	reg, v := gc.LoadParsed()
	check := m.(*localCapabilityManager).computeOffchainCrossCheck(reg, v, dons)
	assert.Equal(t, int64(1), check.matchedCaps)
	assert.Equal(t, int64(1), check.divergences[divergenceConfigMismatch])

	// --- Reactive path: a newer payload re-reconciles via Subscribe, with no new Reconcile. ---
	rawV2, err := marshalOffchainRegistry(offchainReg(2, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		1: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "99"})},
	}))
	require.NoError(t, err)
	applyAsDelegate(t, gc, rawV2)

	require.Eventually(t, func() bool {
		return capture.get("cron@1.0.0") == `{"interval":"99","region":"us"}`
	}, 2*time.Second, 10*time.Millisecond, "launcher should re-reconcile to the v2 offchain config via Subscribe")
}

// TestOffchainRegistry_EndToEnd_GateOff confirms the backwards-compatible default: with the
// cutover off, the same offchain payload is ignored for config and the capability launches with
// the on-chain/TOML config unchanged.
func TestOffchainRegistry_EndToEnd_GateOff(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	lggr := testLogger(t)

	gc := globalconfig.New()
	rawV1, err := marshalOffchainRegistry(offchainReg(1, map[uint32]map[string]*capabilitiespb.CapabilityConfig{
		1: {"cron@1.0.0": specConfigCap(t, map[string]any{"interval": "30"})},
	}))
	require.NoError(t, err)
	applyAsDelegate(t, gc, rawV1)

	capture := &configCapture{last: map[string]string{}}
	localCfg := &testLocalCapabilities{
		allowlisted: map[string]bool{"cron@1.0.0": true},
		useOffchain: false, // gate OFF (default)
	}
	m, err := NewLocalCapabilityManager(lggr, localCfg, capture.builder, gc, false)
	require.NoError(t, err)
	require.NoError(t, m.Start(ctx))
	t.Cleanup(func() { require.NoError(t, m.Close()) })

	dons := []registry.DON{{
		ID:   1,
		Name: testDONName(1),
		CapabilityConfigurations: map[string]registry.CapabilityConfiguration{
			"cron@1.0.0": onchainSpecConfig(t, map[string]any{"interval": "20", "region": "us"}),
		},
	}}
	require.NoError(t, m.Reconcile(ctx, dons))

	// Offchain ignored for config: on-chain values only.
	assert.JSONEq(t, `{"interval":"20","region":"us"}`, capture.get("cron@1.0.0"))
}
