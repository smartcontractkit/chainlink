package job_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/registry"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/globalconfig"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/localcapmgr"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/cresettings"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

// These tests use a real, committing database (heavyweight), unlike pgtest which runs each test
// in a single transaction that is never committed: commit vs rollback is exactly what they test.

var errInjected = errors.New("injected failure after the job change")

type capRegNode struct {
	spawner   job.Spawner
	projector *cresettings.CapRegistryProjector
	gc        *globalconfig.GlobalConfig
}

type capRegEnv struct {
	db  *sqlx.DB
	orm job.ORM
	cfg interface {
		Database() config.Database
	}
}

func newCapRegEnv(t *testing.T) capRegEnv {
	t.Helper()
	cfg, db := heavyweight.FullTestDBV2(t, nil)
	lggr := logger.TestLogger(t)
	keyStore := cltest.NewKeyStore(t, db)
	orm := NewTestORM(t, db, pipeline.NewORM(db, lggr, cfg.JobPipeline().MaxSuccessfulRuns()), bridges.NewORM(db), keyStore)
	return capRegEnv{db: db, orm: orm, cfg: cfg}
}

// boot starts a node: a fresh GlobalConfig, projector and spawner over the shared database. The
// projector is not started unless startProjector is set, so tests control exactly when
// committed state is observed (Refresh).
func (e capRegEnv) boot(t *testing.T, startProjector bool) capRegNode {
	t.Helper()
	lggr := logger.TestLogger(t)
	n := capRegNode{gc: globalconfig.New()}
	n.projector = cresettings.NewCapRegistryProjector(lggr, e.db, n.gc)
	if startProjector {
		require.NoError(t, n.projector.Start(t.Context()))
		t.Cleanup(func() { assert.NoError(t, n.projector.Close()) })
	}
	d := cresettings.NewDelegate(lggr, &loop.AtomicSettings{}, &loop.AtomicSettings{}, n.projector)
	n.spawner = job.NewSpawner(e.orm, e.cfg.Database(), noopChecker{}, map[job.Type]job.Delegate{job.CRESettings: d}, lggr, nil)
	require.NoError(t, n.spawner.Start(t.Context()))
	t.Cleanup(func() { assert.NoError(t, n.spawner.Close()) })
	return n
}

func capRegPayload(version uint64, interval string) string {
	return fmt.Sprintf(`{"version":%d,"dons":{"7":{"capabilityConfigs":{"cron@1.0.0":{"specConfig":{"fields":{"interval":{"stringValue":%q}}}}}}}}`, version, interval)
}

func appliedVersion(gc *globalconfig.GlobalConfig) uint64 {
	_, v := gc.Load()
	return v
}

// replaceInTx deletes old and creates next in one transaction, the way the feeds manager
// approves a new spec version. fail makes the transaction roll back after both changes.
func replaceInTx(ctx context.Context, db *sqlx.DB, n capRegNode, oldID int32, next *job.Job, during func(), fail bool) error {
	return sqlutil.TransactDataSource(ctx, db, nil, func(tx sqlutil.DataSource) error {
		if oldID != 0 {
			if err := n.spawner.DeleteJob(ctx, tx, oldID); err != nil {
				return err
			}
		}
		if during != nil {
			during()
		}
		if next != nil {
			if err := n.spawner.CreateJob(ctx, tx, next); err != nil {
				return err
			}
		}
		if during != nil {
			during()
		}
		if fail {
			return errInjected
		}
		return nil
	})
}

func TestCapRegistry_TransactionRollbackLeavesRuntimeUnchanged(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	env := newCapRegEnv(t)
	n := env.boot(t, false)

	v1 := capRegistrySpec(t, capRegPayload(1, "10"))
	require.NoError(t, n.spawner.CreateJob(ctx, nil, &v1))
	require.NoError(t, n.projector.Refresh(ctx))
	raw, v := n.gc.Load()
	require.Equal(t, uint64(1), v)
	require.JSONEq(t, capRegPayload(1, "10"), raw)

	assertStillV1 := func() {
		t.Helper()
		// Unchanged right after the rollback (nothing was applied inside the transaction)...
		assert.Equal(t, uint64(1), appliedVersion(n.gc), "runtime config changed by a rolled-back transaction")
		// ...and after observing committed state.
		require.NoError(t, n.projector.Refresh(ctx))
		raw, v := n.gc.Load()
		assert.Equal(t, uint64(1), v)
		assert.JSONEq(t, capRegPayload(1, "10"), raw)
		committed, err := job.LoadCapabilitiesRegistrySpec(ctx, env.db)
		require.NoError(t, err)
		require.NotNil(t, committed)
		assert.Equal(t, v1.CRESettingsSpec.Hash, committed.Hash)
	}

	// The steps below are sequential: each starts from the state the previous one left.

	{ // failed replacement
		v2 := capRegistrySpec(t, capRegPayload(2, "20"))
		during := func() { assert.Equal(t, uint64(1), appliedVersion(n.gc), "no uncommitted state may be applied") }
		require.ErrorIs(t, replaceInTx(ctx, env.db, n, v1.ID, &v2, during, true), errInjected)
		assertStillV1()
	}

	{ // failed delete
		during := func() { assert.Equal(t, uint64(1), appliedVersion(n.gc), "no uncommitted state may be applied") }
		require.ErrorIs(t, replaceInTx(ctx, env.db, n, v1.ID, nil, during, true), errInjected)
		assertStillV1()
	}

	{ // rolled-back high-water mark
		// The failed replacement above inserted version 2 with payload A before rolling back. If
		// the high-water mark had survived the rollback, a different version-2 payload would now
		// be stale.
		other := capRegistrySpec(t, capRegPayload(2, "25"))
		require.NoError(t, replaceInTx(ctx, env.db, n, v1.ID, &other, nil, false))
		require.NoError(t, n.projector.Refresh(ctx))
		raw, v := n.gc.Load()
		assert.Equal(t, uint64(2), v)
		assert.JSONEq(t, capRegPayload(2, "25"), raw)

		// Committed delete, then a failed create: stays cleared.
		require.NoError(t, n.spawner.DeleteJob(ctx, nil, other.ID))
		require.NoError(t, n.projector.Refresh(ctx))
		require.Equal(t, uint64(0), appliedVersion(n.gc))
		v3 := capRegistrySpec(t, capRegPayload(3, "30"))
		require.ErrorIs(t, replaceInTx(ctx, env.db, n, 0, &v3, nil, true), errInjected)
		require.NoError(t, n.projector.Refresh(ctx))
		reg, v := n.gc.LoadParsed()
		assert.Nil(t, reg)
		assert.Equal(t, uint64(0), v)
	}
}

// TestCapRegistry_DeleteStaleSubmitRestart covers durable version enforcement: once a version
// has been committed, an older payload is rejected even after the job carrying it is deleted
// and the node restarts.
func TestCapRegistry_DeleteStaleSubmitRestart(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	env := newCapRegEnv(t)
	first := env.boot(t, false)

	v5 := capRegistrySpec(t, capRegPayload(5, "50"))
	require.NoError(t, first.spawner.CreateJob(ctx, nil, &v5))
	require.NoError(t, first.spawner.DeleteJob(ctx, nil, v5.ID))
	require.NoError(t, first.projector.Refresh(ctx))
	require.Equal(t, uint64(0), appliedVersion(first.gc))

	stale := capRegistrySpec(t, capRegPayload(4, "40"))
	require.ErrorIs(t, first.spawner.CreateJob(ctx, nil, &stale), job.ErrCRESettingsCapRegistryStale)
	sameVersionOtherPayload := capRegistrySpec(t, capRegPayload(5, "55"))
	require.ErrorIs(t, first.spawner.CreateJob(ctx, nil, &sameVersionOtherPayload), job.ErrCRESettingsCapRegistryStale)
	require.NoError(t, first.projector.Refresh(ctx))
	require.Equal(t, uint64(0), appliedVersion(first.gc))

	// Restart: nothing stale was persisted, so nothing is applied on boot.
	second := env.boot(t, true)
	committed, err := job.LoadCapabilitiesRegistrySpec(ctx, env.db)
	require.NoError(t, err)
	assert.Nil(t, committed)
	require.Never(t, func() bool { return appliedVersion(second.gc) != 0 }, 300*time.Millisecond, 20*time.Millisecond)

	// The high-water mark survived the restart.
	staleAgain := capRegistrySpec(t, capRegPayload(4, "40"))
	require.ErrorIs(t, second.spawner.CreateJob(ctx, nil, &staleAgain), job.ErrCRESettingsCapRegistryStale)

	// Re-creating exactly the payload that set the mark is allowed, and is applied.
	again := capRegistrySpec(t, capRegPayload(5, "50"))
	require.NoError(t, second.spawner.CreateJob(ctx, nil, &again))
	require.Eventually(t, func() bool { return appliedVersion(second.gc) == 5 }, 5*time.Second, 20*time.Millisecond)
}

// TestCapRegistry_ConcurrentReconcileDuringReplacement replaces the capabilities_registry job
// repeatedly (delete + create in one transaction, held open for a while between and after the
// two changes) while the projector polls quickly and the LocalCapabilityManager reconciles
// continuously. No observer may ever see the payload cleared, and no capability may ever be
// launched with the legacy (on-chain) value in between two offchain versions.
func TestCapRegistry_ConcurrentReconcileDuringReplacement(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	env := newCapRegEnv(t)
	n := env.boot(t, true)

	current := capRegistrySpec(t, capRegPayload(1, "off-1"))
	require.NoError(t, n.spawner.CreateJob(ctx, nil, &current))
	require.Eventually(t, func() bool { return appliedVersion(n.gc) == 1 }, 5*time.Second, 10*time.Millisecond)

	rec := &capRegLaunches{}
	mgr, err := localcapmgr.NewLocalCapabilityManager(logger.TestLogger(t), capRegLocalCfg{}, rec.newServices, n.gc, true)
	require.NoError(t, err)
	require.NoError(t, mgr.Start(ctx))
	t.Cleanup(func() { assert.NoError(t, mgr.Close()) })
	dons := []registry.DON{{
		ID:                       7,
		CapabilityConfigurations: map[string]registry.CapabilityConfiguration{"cron@1.0.0": onchainCronConfig(t, "onchain")},
	}}
	require.NoError(t, mgr.Reconcile(ctx, dons))

	var stop atomic.Bool
	var clearedSeen atomic.Int64
	var wg sync.WaitGroup
	wg.Go(func() { // registry-driven reconciles
		for !stop.Load() {
			assert.NoError(t, mgr.Reconcile(ctx, dons))
		}
	})
	wg.Go(func() { // direct observer of the runtime registry
		for !stop.Load() {
			if reg, _ := n.gc.LoadParsed(); reg == nil {
				clearedSeen.Add(1)
			}
		}
	})

	const replacements = 5
	for i := 2; i <= replacements+1; i++ {
		next := capRegistrySpec(t, capRegPayload(uint64(i), fmt.Sprintf("off-%d", i)))
		hold := func() { time.Sleep(100 * time.Millisecond) } // let polls/reconciles run mid-transaction
		require.NoError(t, replaceInTx(ctx, env.db, n, current.ID, &next, hold, false))
		current = next
	}
	require.Eventually(t, func() bool { return appliedVersion(n.gc) == replacements+1 }, 5*time.Second, 10*time.Millisecond)
	require.NoError(t, mgr.Reconcile(ctx, dons))
	stop.Store(true)
	wg.Wait()

	assert.Zero(t, clearedSeen.Load(), "the runtime registry was observed cleared during a replacement")
	intervals := rec.intervals(t)
	require.NotEmpty(t, intervals)
	for _, iv := range intervals {
		assert.NotEqual(t, "onchain", iv, "a capability was launched with the legacy fallback during a replacement")
	}
	assert.Equal(t, fmt.Sprintf("off-%d", replacements+1), intervals[len(intervals)-1])
}

func onchainCronConfig(t *testing.T, interval string) registry.CapabilityConfiguration {
	t.Helper()
	vm, err := values.NewMap(map[string]any{"interval": interval})
	require.NoError(t, err)
	b, err := proto.Marshal(&capabilitiespb.CapabilityConfig{SpecConfig: values.ProtoMap(vm)})
	require.NoError(t, err)
	return registry.CapabilityConfiguration{Config: b}
}

type capRegLaunches struct {
	mu      sync.Mutex
	configs []string
}

func (r *capRegLaunches) newServices(_ context.Context, _ string, _ uint32, _ string, configJSON string, _ *ocrtypes.ContractConfig) ([]job.ServiceCtx, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configs = append(r.configs, configJSON)
	return nil, nil
}

func (r *capRegLaunches) intervals(t *testing.T) []string {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]string, 0, len(r.configs))
	for _, c := range r.configs {
		var m struct {
			Interval string `json:"interval"`
		}
		require.NoError(t, json.Unmarshal([]byte(c), &m))
		out = append(out, m.Interval)
	}
	return out
}

type capRegLocalCfg struct{}

func (capRegLocalCfg) RegistryBasedLaunchAllowlist() []string                 { return []string{"cron@1.0.0"} }
func (capRegLocalCfg) Capabilities() map[string]config.CapabilityNodeConfig   { return nil }
func (capRegLocalCfg) IsAllowlisted(capabilityID string) bool                 { return capabilityID == "cron@1.0.0" }
func (capRegLocalCfg) GetCapabilityConfig(string) config.CapabilityNodeConfig { return nil }
func (capRegLocalCfg) UseOffchainRegistry() bool                              { return true }

// TestCapRegistry_DBDiscriminatorConstraints inserts rows directly, bypassing validation, to
// show the schema alone guarantees that every row routed to the offchain capabilities registry
// has config_type = 'capabilities_registry' (the column covered by the single-job unique index)
// and that no other row can carry an offchain payload.
func TestCapRegistry_DBDiscriminatorConstraints(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	env := newCapRegEnv(t)

	insert := func(settings, configType, offchain string) error {
		_, err := env.db.ExecContext(ctx, `INSERT INTO cre_settings_specs (settings, hash, config_type, offchain_config, created_at, updated_at)
			VALUES ($1, 'h', $2, $3, NOW(), NOW())`, settings, configType, offchain)
		return err
	}
	const payload = `{"version":1}`

	for _, tc := range []struct {
		name                         string
		settings, configType, offchn string
	}{
		{"payload without config_type", "", "", payload},
		{"payload with settings config_type", `Foo = "bar"`, "settings", payload},
		{"payload with shard_assignment config_type", "", "shard_assignment", payload},
		{"case variant", "", "Capabilities_Registry", payload},
		{"padded variant", "", " capabilities_registry", payload},
		{"capabilities_registry without payload", "", "capabilities_registry", ""},
		{"capabilities_registry with settings", `config_type = "capabilities_registry"`, "capabilities_registry", payload},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.ErrorContains(t, insert(tc.settings, tc.configType, tc.offchn), "violates check constraint")
		})
	}

	// Legacy and settings-family rows are unaffected.
	require.NoError(t, insert(`Foo = "bar"`, "", ""))
	require.NoError(t, insert(`config_type = "shard_assignment"`, "", ""))
	require.NoError(t, insert(`Foo = "bar"`, "settings", ""))

	// Exactly one capabilities_registry row.
	require.NoError(t, insert("", "capabilities_registry", payload))
	require.ErrorContains(t, insert("", "capabilities_registry", `{"version":2}`), "idx_cre_settings_specs_single_capabilities_registry")
}
