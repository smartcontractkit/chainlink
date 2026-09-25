package vault

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaultutils"
)

// limitCheckingKV mirrors libocr's limitCheckWriteSet on top of the kv test
// store: it errors a write/delete that would push the write set past the
// limits, exactly as the real transaction would. StateTransition-level budget
// tests run against it — a round that defers correctly never trips it.
type limitCheckingKV struct {
	inner     *kv
	keyLimit  int
	byteLimit int

	modified map[string][]byte
	keys     int
	bytes    int
}

func newLimitCheckingKV(inner *kv, keyLimit, byteLimit int) *limitCheckingKV {
	return &limitCheckingKV{
		inner:     inner,
		keyLimit:  keyLimit,
		byteLimit: byteLimit,
		modified:  map[string][]byte{},
	}
}

func (l *limitCheckingKV) Read(key []byte) ([]byte, error) {
	return l.inner.Read(key)
}

func (l *limitCheckingKV) modify(key, value []byte) error {
	add, sub := 0, 0
	if prev, ok := l.modified[string(key)]; ok {
		add = len(value)
		sub = len(prev)
	} else {
		add = len(key) + len(value)
		if l.keys+1 > l.keyLimit {
			return fmt.Errorf("keys %d exceed limit %d", l.keys+1, l.keyLimit)
		}
		l.keys++
	}
	if l.bytes-sub+add > l.byteLimit {
		return fmt.Errorf("keys + values length %d exceeds limit %d", l.bytes-sub+add, l.byteLimit)
	}
	l.bytes = l.bytes - sub + add
	l.modified[string(key)] = value
	return nil
}

func (l *limitCheckingKV) Write(key, value []byte) error {
	if err := l.modify(key, value); err != nil {
		return err
	}
	return l.inner.Write(key, value)
}

func (l *limitCheckingKV) Delete(key []byte) error {
	if err := l.modify(key, nil); err != nil {
		return err
	}
	return l.inner.Delete(key)
}

var _ ocr3_1types.KeyValueStateReadWriter = (*limitCheckingKV)(nil)

type failingReadWriter struct {
	err error
}

func (f *failingReadWriter) Read([]byte) ([]byte, error) { return nil, f.err }

func (f *failingReadWriter) Write([]byte, []byte) error { return f.err }

func (f *failingReadWriter) Delete([]byte) error { return f.err }

func TestKVWriteBudgetTracker_MirrorsLimitCheckWriteSetArithmetic(t *testing.T) {
	t.Parallel()
	base := &kv{m: map[string]response{}}
	tracker := newKVWriteBudgetTracker(base)

	v10, v20, v7, v5 := make([]byte, 10), make([]byte, 20), make([]byte, 7), make([]byte, 5)

	// Hand-computed expectations mirroring limitCheckWriteSet.modify: a new
	// key costs 1 key plus key+value bytes; rewriting a key already in the
	// write set counts only its final value; a delete of a new key costs the
	// key bytes.
	steps := []struct {
		op        string
		key       string
		value     []byte
		wantKeys  int
		wantBytes int
	}{
		{"write", "a", v10, 1, 1 + 10},
		{"write", "a", v20, 1, 1 + 20},
		{"delete", "a", nil, 1, 1},
		{"write", "bb", v5, 2, 1 + 2 + 5},
		{"delete", "ccc", nil, 3, 1 + 2 + 5 + 3},
		{"write", "a", v7, 3, 1 + 2 + 5 + 3 + 7},
	}
	for i, s := range steps {
		switch s.op {
		case "write":
			require.NoError(t, tracker.Write([]byte(s.key), s.value), "step %d", i)
		case "delete":
			require.NoError(t, tracker.Delete([]byte(s.key)), "step %d", i)
		}
		got := tracker.consumed()
		assert.Equal(t, s.wantKeys, got.keys, "keys after step %d", i)
		assert.Equal(t, s.wantBytes, got.bytes, "bytes after step %d", i)
	}

	// Writes pass through to the underlying store.
	assert.Equal(t, v7, base.m["a"].data)
	assert.Nil(t, base.m["ccc"].data)

	// Underlying errors are propagated and not counted.
	tracker = newKVWriteBudgetTracker(&failingReadWriter{err: errors.New("boom")})
	require.Error(t, tracker.Write([]byte("x"), v10))
	assert.Equal(t, kvWriteCost{}, tracker.consumed())
}

func TestProjectedProcessedWriteCost(t *testing.T) {
	t.Parallel()
	const maxSecretsPerOwner = 5
	const decodedCiphertext = 100

	newItem := func(payload proto.Message) *vaultcommon.StoredPendingQueueItem {
		anyPayload, err := anypb.New(payload)
		require.NoError(t, err)
		return &vaultcommon.StoredPendingQueueItem{Id: "item", Item: anyPayload}
	}
	id := func(owner, key string) *vaultcommon.SecretIdentifier {
		return &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: key}
	}
	encSecret := func(owner, key string) *vaultcommon.EncryptedSecret {
		return &vaultcommon.EncryptedSecret{
			Id:             id(owner, key),
			EncryptedValue: strings.Repeat("ab", decodedCiphertext),
		}
	}
	// mdStore builds a read store over owner metadata records.
	mdStore := func(t *testing.T, md map[string]*vaultcommon.StoredMetadata) ReadKVStore {
		t.Helper()
		underlying := &kv{m: map[string]response{}}
		for owner, m := range md {
			underlying.m[metadataPrefix+owner] = response{data: protoMarshal(t, m)}
		}
		return newTestReadStore(t, underlying)
	}
	secretRecordCost := func(identifier *vaultcommon.SecretIdentifier) kvWriteCost {
		return kvWriteCost{
			keys:  1,
			bytes: len(keyPrefix+vaulttypes.KeyFor(identifier)) + storedSecretWireSize(decodedCiphertext),
		}
	}
	metadataCost := func(owner string, identifiers ...*vaultcommon.SecretIdentifier) kvWriteCost {
		return kvWriteCost{
			keys:  1,
			bytes: len(metadataPrefix+owner) + proto.Size(&vaultcommon.StoredMetadata{SecretIdentifiers: identifiers}),
		}
	}

	t.Run("nil and unprocessable items cost nothing", func(t *testing.T) {
		t.Parallel()
		store := mdStore(t, nil)
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, nil, maxSecretsPerOwner))
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, &vaultcommon.StoredPendingQueueItem{Id: "x"}, maxSecretsPerOwner))
		bogus := &vaultcommon.StoredPendingQueueItem{Id: "x", Item: &anypb.Any{TypeUrl: "bogus"}}
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, bogus, maxSecretsPerOwner))
	})

	t.Run("read-only request types cost nothing", func(t *testing.T) {
		t.Parallel()
		store := mdStore(t, nil)
		get := newItem(&vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id("o", "k")}}})
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, get, maxSecretsPerOwner))
		list := newItem(&vaultcommon.ListSecretIdentifiersRequest{RequestId: "r", Owner: "o"})
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, list, maxSecretsPerOwner))
	})

	t.Run("create projects admitted adds plus the metadata rewrite", func(t *testing.T) {
		t.Parallel()
		store := mdStore(t, nil)
		create := newItem(&vaultcommon.CreateSecretsRequest{
			RequestId: "r",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				encSecret("o1", "k1"), encSecret("o1", "k2"), encSecret("o2", "k1"),
			},
		})
		want := secretRecordCost(id("o1", "k1")).
			add(secretRecordCost(id("o1", "k2"))).
			add(secretRecordCost(id("o2", "k1"))).
			add(metadataCost("o1", id("o1", "k1"), id("o1", "k2"))).
			add(metadataCost("o2", id("o2", "k1")))
		assert.Equal(t, want, projectedProcessedWriteCost(t.Context(), store, create, maxSecretsPerOwner))
	})

	t.Run("create skips already-present identifiers and capacity", func(t *testing.T) {
		t.Parallel()
		existing := id("o1", "k0")
		store := mdStore(t, map[string]*vaultcommon.StoredMetadata{
			"o1": {SecretIdentifiers: []*vaultcommon.SecretIdentifier{existing}},
		})
		create := newItem(&vaultcommon.CreateSecretsRequest{
			RequestId: "r",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				// k0 already exists; with maxSecretsPerOwner=2 only k1 is
				// admitted before k2 fails the capacity check.
				encSecret("o1", "k0"), encSecret("o1", "k1"), encSecret("o1", "k2"),
			},
		})
		want := secretRecordCost(id("o1", "k1")).
			add(metadataCost("o1", existing, id("o1", "k1")))
		assert.Equal(t, want, projectedProcessedWriteCost(t.Context(), store, create, 2))
	})

	t.Run("create at capacity writes nothing", func(t *testing.T) {
		t.Parallel()
		store := mdStore(t, map[string]*vaultcommon.StoredMetadata{
			"o1": {SecretIdentifiers: []*vaultcommon.SecretIdentifier{id("o1", "k0")}},
		})
		create := newItem(&vaultcommon.CreateSecretsRequest{
			RequestId:        "r",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{encSecret("o1", "k1")},
		})
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, create, 1))
	})

	t.Run("update projects secret records only", func(t *testing.T) {
		t.Parallel()
		store := mdStore(t, nil)
		update := newItem(&vaultcommon.UpdateSecretsRequest{
			RequestId:        "r",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{encSecret("o1", "k1"), encSecret("o2", "k1")},
		})
		want := secretRecordCost(id("o1", "k1")).add(secretRecordCost(id("o2", "k1")))
		assert.Equal(t, want, projectedProcessedWriteCost(t.Context(), store, update, maxSecretsPerOwner))
	})

	t.Run("delete projects key deletes plus the metadata rewrite", func(t *testing.T) {
		t.Parallel()
		k1, k2, k3 := id("o1", "k1"), id("o1", "k2"), id("o2", "k3")
		store := mdStore(t, map[string]*vaultcommon.StoredMetadata{
			"o1": {SecretIdentifiers: []*vaultcommon.SecretIdentifier{k1, k2}},
			"o2": {SecretIdentifiers: []*vaultcommon.SecretIdentifier{k3}},
		})
		del := newItem(&vaultcommon.DeleteSecretsRequest{
			RequestId: "r",
			// k1 is duplicated and one id is unknown: both cost nothing.
			Ids: []*vaultcommon.SecretIdentifier{k3, k1, k1, id("ghost", "k")},
		})
		want := metadataCost("o1", k2).
			add(metadataCost("o2")).
			add(kvWriteCost{keys: 1, bytes: len(keyPrefix + vaulttypes.KeyFor(k1))}).
			add(kvWriteCost{keys: 1, bytes: len(keyPrefix + vaulttypes.KeyFor(k3))})
		assert.Equal(t, want, projectedProcessedWriteCost(t.Context(), store, del, maxSecretsPerOwner))
	})

	t.Run("delete of absent ids costs nothing", func(t *testing.T) {
		t.Parallel()
		store := mdStore(t, map[string]*vaultcommon.StoredMetadata{
			"o1": {SecretIdentifiers: []*vaultcommon.SecretIdentifier{id("o1", "k0")}},
		})
		del := newItem(&vaultcommon.DeleteSecretsRequest{
			RequestId: "r",
			Ids:       []*vaultcommon.SecretIdentifier{id("ghost", "k")},
		})
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, del, maxSecretsPerOwner))
	})

	t.Run("metadata read errors cost nothing", func(t *testing.T) {
		t.Parallel()
		store := newTestReadStore(t, &failingReadWriter{err: errors.New("boom")})
		create := newItem(&vaultcommon.CreateSecretsRequest{
			RequestId:        "r",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{encSecret("o1", "k1")},
		})
		del := newItem(&vaultcommon.DeleteSecretsRequest{
			RequestId: "r",
			Ids:       []*vaultcommon.SecretIdentifier{id("o1", "k1")},
		})
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, create, maxSecretsPerOwner))
		assert.Equal(t, kvWriteCost{}, projectedProcessedWriteCost(t.Context(), store, del, maxSecretsPerOwner))
	})
}

// TestDeleteSecretsProjectedWriteCost_MatchesRealizedDeletes pins the delete
// projection's exactness: replaying DeleteSecret's exact per-id write sequence
// (owner metadata rewrite plus secret key delete) over the same starting state
// realizes exactly the projected write set — identifiers absent from the
// metadata, and duplicates, error out during processing before any write, so
// both the projection and the replay skip them identically.
func TestDeleteSecretsProjectedWriteCost_MatchesRealizedDeletes(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(2))

	component := func() string { return strings.Repeat("c", 1+rng.Intn(64)) }

	for iteration := range 200 {
		maxSecretsPerOwner := 1 + rng.Intn(10)

		// Seed owners with at most maxSecretsPerOwner identifiers each,
		// mirroring the create-path invariant DeleteSecret relies on.
		seeded := map[string][]*vaultcommon.SecretIdentifier{}
		for o := range 1 + rng.Intn(4) {
			owner := fmt.Sprintf("owner-%d-%s", o, component())
			for range rng.Intn(maxSecretsPerOwner + 1) {
				seeded[owner] = append(seeded[owner], &vaultcommon.SecretIdentifier{
					Owner:     owner,
					Namespace: component(),
					Key:       component(),
				})
			}
		}

		// The batch mixes realizable deletes, duplicates (the second delete
		// errors out before any write), and sometimes an unknown id (same).
		var ids []*vaultcommon.SecretIdentifier
		for _, list := range seeded {
			for _, sid := range list {
				if rng.Intn(2) == 0 {
					ids = append(ids, sid)
					if rng.Intn(4) == 0 {
						ids = append(ids, sid)
					}
				}
			}
		}
		if rng.Intn(2) == 0 {
			ids = append(ids, &vaultcommon.SecretIdentifier{Owner: "ghost", Namespace: "main", Key: "unknown"})
		}
		rng.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
		if len(ids) > vaulttypes.MaxBatchSize {
			ids = ids[:vaulttypes.MaxBatchSize]
		}

		// Project the batch from a read store over the seeded metadata.
		underlying := &kv{m: map[string]response{}}
		for owner, list := range seeded {
			underlying.m[metadataPrefix+owner] = response{data: protoMarshal(t, &vaultcommon.StoredMetadata{SecretIdentifiers: list})}
		}
		req := &vaultcommon.DeleteSecretsRequest{RequestId: "r", Ids: ids}
		projected := deleteSecretsProjectedWriteCost(t.Context(), newTestReadStore(t, underlying), req)

		// Replay the exact write sequence DeleteSecret performs per id:
		// rewrite the owner metadata with the id removed, then delete the
		// secret key.
		tracker := newKVWriteBudgetTracker(&kv{m: map[string]response{}})
		state := map[string][]*vaultcommon.SecretIdentifier{}
		for owner, list := range seeded {
			state[owner] = append([]*vaultcommon.SecretIdentifier(nil), list...)
		}
		for _, sid := range ids {
			if sid == nil {
				continue
			}
			var next []*vaultcommon.SecretIdentifier
			found := false
			for _, x := range state[sid.Owner] {
				if vaulttypes.KeyFor(x) == vaulttypes.KeyFor(sid) {
					found = true
					continue
				}
				next = append(next, x)
			}
			if !found {
				continue
			}
			state[sid.Owner] = next
			b, err := proto.Marshal(&vaultcommon.StoredMetadata{SecretIdentifiers: next})
			require.NoError(t, err)
			require.NoError(t, tracker.Write([]byte(metadataPrefix+sid.Owner), b))
			require.NoError(t, tracker.Delete([]byte(keyPrefix+vaulttypes.KeyFor(sid))))
		}

		assert.Equal(t, projected, tracker.consumed(),
			"iteration %d: the delete projection must match DeleteSecret's realized write set", iteration)
	}
}

func TestPackPendingQueueWithinKVBudget(t *testing.T) {
	t.Parallel()
	const defaultKeys = 300
	const defaultBytes = 1468006

	item := func(id string, pad int) *vaultcommon.StoredPendingQueueItem {
		anyPayload, err := anypb.New(&vaultcommon.ListSecretIdentifiersRequest{RequestId: "r", Owner: "o"})
		require.NoError(t, err)
		return &vaultcommon.StoredPendingQueueItem{Id: id + string(make([]byte, pad)), Item: anyPayload}
	}

	t.Run("depth tracks the configured key budget", func(t *testing.T) {
		t.Parallel()
		kept := make([]*vaultcommon.StoredPendingQueueItem, 500)
		for i := range kept {
			kept[i] = item("id", 0)
		}
		// remaining = 400 - 0 - 4 = 396 keys; positionalCount+1 <= 396 caps
		// the packed prefix at 395. There is no separate depth cap: the packer
		// tracks the configured limits, so the committed depth stays within
		// keyLimit minus the reserve and a delete-only purge always fits.
		packed := packPendingQueueWithinKVBudget(0, kept, kvWriteCost{}, 400, defaultBytes)
		assert.Equal(t, 395, packed)
	})

	t.Run("key budget trims the packed prefix", func(t *testing.T) {
		t.Parallel()
		kept := make([]*vaultcommon.StoredPendingQueueItem, 400)
		for i := range kept {
			kept[i] = item("id", 0)
		}
		// remaining = 300 - 0 - 4 = 296 keys; positionalCount+1 <= 296 caps
		// the packed prefix at 295.
		packed := packPendingQueueWithinKVBudget(0, kept, kvWriteCost{}, defaultKeys, defaultBytes)
		assert.Equal(t, 295, packed)
	})

	t.Run("positional overlap lets new items reuse old queue keys", func(t *testing.T) {
		t.Parallel()
		kept := make([]*vaultcommon.StoredPendingQueueItem, 100)
		for i := range kept {
			kept[i] = item("id", 0)
		}
		// A 250-deep old queue already owns positional keys 0..249; packing
		// 100 new items rewrites that prefix instead of adding keys, so the
		// union stays 250 and everything fits.
		packed := packPendingQueueWithinKVBudget(250, kept, kvWriteCost{}, defaultKeys, defaultBytes)
		assert.Equal(t, 100, packed)
	})

	t.Run("byte budget trims large items", func(t *testing.T) {
		t.Parallel()
		kept := make([]*vaultcommon.StoredPendingQueueItem, 100)
		for i := range kept {
			kept[i] = item("id", 4096)
		}
		// Each item's marshaled size is dominated by the 4KB id padding; a
		// 30KB byte budget packs only a handful.
		packed := packPendingQueueWithinKVBudget(0, kept, kvWriteCost{}, defaultKeys, 30*1024)
		assert.Less(t, packed, 10)
		assert.Positive(t, packed)
	})

	t.Run("consumed budget leaves nothing to pack", func(t *testing.T) {
		t.Parallel()
		kept := []*vaultcommon.StoredPendingQueueItem{item("id", 0)}
		assert.Equal(t, 0, packPendingQueueWithinKVBudget(0, kept, kvWriteCost{keys: defaultKeys - kvBudgetReserveKeys, bytes: 0}, defaultKeys, defaultBytes))
		assert.Equal(t, 0, packPendingQueueWithinKVBudget(0, kept, kvWriteCost{keys: 0, bytes: defaultBytes - kvBudgetReserveBytes}, defaultKeys, defaultBytes))
	})

	t.Run("empty candidates pack nothing", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 0, packPendingQueueWithinKVBudget(10, nil, kvWriteCost{}, defaultKeys, defaultBytes))
	})

	t.Run("deterministic", func(t *testing.T) {
		t.Parallel()
		kept := make([]*vaultcommon.StoredPendingQueueItem, 50)
		for i := range kept {
			kept[i] = item("id", i*10)
		}
		first := packPendingQueueWithinKVBudget(7, kept, kvWriteCost{keys: 12, bytes: 5000}, defaultKeys, defaultBytes)
		second := packPendingQueueWithinKVBudget(7, kept, kvWriteCost{keys: 12, bytes: 5000}, defaultKeys, defaultBytes)
		assert.Equal(t, first, second)
	})
}

// TestPackPendingQueueWithinKVBudget_NeverExceedsLimits is the packing
// property: for random shapes, the packed rewrite performed on top of the
// already-consumed write set never exceeds the limits, and the packer is
// maximal — one more item would not fit within the reserved remainder.
func TestPackPendingQueueWithinKVBudget_NeverExceedsLimits(t *testing.T) {
	t.Parallel()
	rng := rand.New(rand.NewSource(1))

	item := func(id string, pad int) *vaultcommon.StoredPendingQueueItem {
		anyPayload, err := anypb.New(&vaultcommon.ListSecretIdentifiersRequest{RequestId: "r", Owner: "o"})
		require.NoError(t, err)
		return &vaultcommon.StoredPendingQueueItem{Id: id + string(make([]byte, pad)), Item: anyPayload}
	}

	// simulate performs the exact write sequence WritePendingQueue would.
	simulate := func(oldCount int, kept []*vaultcommon.StoredPendingQueueItem) kvWriteCost {
		tracker := newKVWriteBudgetTracker(&kv{m: map[string]response{}})
		for j := range oldCount {
			require.NoError(t, tracker.Delete([]byte(pendingQueueItemPrefix+strconv.Itoa(j))))
		}
		for j, it := range kept {
			require.NoError(t, tracker.Write([]byte(pendingQueueItemPrefix+strconv.Itoa(j)), protoMarshal(t, it)))
		}
		require.NoError(t, tracker.Write([]byte(pendingQueueIndex), protoMarshal(t, &vaultcommon.StoredPendingQueueIndex{Length: int64(len(kept))})))
		return tracker.consumed()
	}

	for iteration := range 200 {
		keyLimit := 1 + rng.Intn(400)
		byteLimit := 1 + rng.Intn(200000)

		// Only states the gates can produce are tested: the committed queue
		// depth was itself packed by a prior round (its full rewrite —
		// positional keys plus index — fits both budgets minus reserve), and
		// the consumed budget obeys the Phase 1 gate's floor reservation
		// (consumed+floor fits the limits). The packer's guarantee holds for
		// exactly these reachable states.
		maxOldCount := max(0, keyLimit-kvBudgetReserveKeys-1)
		for maxOldCount > 0 && mandatoryPendingQueueRewriteFloor(maxOldCount).bytes > byteLimit-kvBudgetReserveBytes {
			maxOldCount--
		}
		oldCount := rng.Intn(maxOldCount + 1)
		floor := mandatoryPendingQueueRewriteFloor(oldCount)
		maxConsumedKeys := max(0, keyLimit-floor.keys)
		maxConsumedBytes := max(0, byteLimit-floor.bytes)

		keptCount := rng.Intn(300)
		kept := make([]*vaultcommon.StoredPendingQueueItem, keptCount)
		for i := range kept {
			kept[i] = item(fmt.Sprintf("id-%d-", i), rng.Intn(1024))
		}
		consumed := kvWriteCost{keys: rng.Intn(maxConsumedKeys + 1), bytes: rng.Intn(maxConsumedBytes + 1)}

		packed := packPendingQueueWithinKVBudget(oldCount, kept, consumed, keyLimit, byteLimit)

		// Guarantee: the packed rewrite plus the consumed writes fits.
		realized := simulate(oldCount, kept[:packed]).add(consumed)
		assert.LessOrEqual(t, realized.keys, keyLimit, "iteration %d keys", iteration)
		assert.LessOrEqual(t, realized.bytes, byteLimit, "iteration %d bytes", iteration)

		// Maximality: one more item would not fit within the reserved
		// remainder, unless the candidate count already bound us.
		if packed < keptCount {
			realizedPlusOne := simulate(oldCount, kept[:packed+1]).add(consumed)
			assert.True(t,
				realizedPlusOne.keys > keyLimit-kvBudgetReserveKeys ||
					realizedPlusOne.bytes > byteLimit-kvBudgetReserveBytes,
				"iteration %d: packed %d of %d but one more fit (realized+1 = %+v, limits %d/%d, consumed %+v)",
				iteration, packed, keptCount, realizedPlusOne, keyLimit, byteLimit, consumed)
		}
	}
}

// budgetCreateFixture bundles everything a round needs for one create request:
// the committed pending-queue item plus the observation contribution that
// lets StateTransition process it.
type budgetCreateFixture struct {
	id   *vaultcommon.SecretIdentifier
	req  *vaultcommon.CreateSecretsRequest
	resp *vaultcommon.CreateSecretsResponse
	item *vaultcommon.StoredPendingQueueItem
}

// newBudgetCreateFixtures builds n create requests, each with a distinct
// address owner so every processed request realizes its full write cost (one
// secret record plus one owner metadata record) and every ciphertext carries
// the owner label the contribution guard verifies.
func newBudgetCreateFixtures(t *testing.T, pk *tdh2easy.PublicKey, prefix string, n int) []budgetCreateFixture {
	t.Helper()
	fixtures := make([]budgetCreateFixture, n)
	for i := range fixtures {
		ownerAddr := common.Address{byte(i + 1)}
		enc, err := vaultutils.EncryptSecretWithWorkflowOwner("secret value", pk, ownerAddr)
		require.NoError(t, err)
		id := &vaultcommon.SecretIdentifier{
			Owner:     ownerAddr.Hex(),
			Namespace: "main",
			Key:       fmt.Sprintf("%s_secret_%d", prefix, i),
		}
		req := &vaultcommon.CreateSecretsRequest{
			RequestId: fmt.Sprintf("%s_req_%d", prefix, i),
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: id, EncryptedValue: enc},
			},
		}
		resp := &vaultcommon.CreateSecretsResponse{
			Responses: []*vaultcommon.CreateSecretResponse{{Id: id, Success: false, Error: ""}},
		}
		anyReq, err := anypb.New(req)
		require.NoError(t, err)
		fixtures[i] = budgetCreateFixture{
			id:   id,
			req:  req,
			resp: resp,
			item: &vaultcommon.StoredPendingQueueItem{Id: vaulttypes.KeyFor(id), Item: anyReq},
		}
	}
	return fixtures
}

// attributedObservations replicates one observations blob across three
// observers, matching the convergence pattern of the existing StateTransition
// tests.
func attributedObservations(t *testing.T, observations ...observation) []types.AttributedObservation {
	t.Helper()
	obsb := marshalObservations(t, observations...)
	return []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(obsb)},
		{Observer: 1, Observation: types.Observation(obsb)},
		{Observer: 2, Observation: types.Observation(obsb)},
	}
}

func createObservations(fixtures []budgetCreateFixture) []observation {
	observations := make([]observation, 0, len(fixtures))
	for _, f := range fixtures {
		observations = append(observations, observation{id: f.id, req: f.req, resp: f.resp})
	}
	return observations
}

func budgetRound(t *testing.T, r *ReportingPlugin, underlying *kv, keyLimit, byteLimit int, seqNr uint64, aos []types.AttributedObservation, blobFetcher ocr3_1types.BlobFetcher) *vaultcommon.Outcomes {
	t.Helper()
	// A fresh enforcing writer per round mirrors the real per-round transaction.
	rw := newLimitCheckingKV(underlying, keyLimit, byteLimit)
	reportPrecursor, err := r.StateTransition(t.Context(), seqNr, types.AttributedQuery{}, aos, rw, blobFetcher)
	require.NoError(t, err, "round %d must not fail: the budget gates defer instead of erroring", seqNr)
	outcomes := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal(reportPrecursor, outcomes))
	return outcomes
}

// TestStateTransition_KVBudgetGate_DefersProcessingAndDrainsOverRounds is the
// legacy deep-queue regression: a committed queue deeper than the budget can
// process in one round must not fail the round. Deferred items leave the
// committed queue, re-enter via re-broadcast the next round, and the backlog
// drains over rounds without ever tripping the enforcing writer.
func TestStateTransition_KVBudgetGate_DefersProcessingAndDrainsOverRounds(t *testing.T) {
	t.Parallel()
	// keyLimit 15 with a 6-deep queue: the mandatory rewrite floor (7 keys)
	// plus four full creates (2 keys each) lands exactly on the limit, so the
	// fifth item is deferred and the round still succeeds.
	const keyLimit = 15
	const byteLimit = 100000

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withKVWriteBudgetLimits(keyLimit, byteLimit))

	underlying := &kv{m: map[string]response{}}
	fixtures := newBudgetCreateFixtures(t, pk, "drain", 6)
	items := make([]*vaultcommon.StoredPendingQueueItem, len(fixtures))
	for i, f := range fixtures {
		items[i] = f.item
	}
	require.NoError(t, newTestWriteStore(t, underlying).WritePendingQueue(t.Context(), items))

	rs := newTestReadStore(t, underlying)

	// The blobber backs round 2's re-broadcast of the two deferred items.
	bf := &blobber{blobs: [][]byte{protoMarshal(t, fixtures[4].item), protoMarshal(t, fixtures[5].item)}}
	r.unmarshalBlob = bf.unmarshalBlob

	// Round 1: six processable creates against a budget that fits four.
	outcomes := budgetRound(t, r, underlying, keyLimit, byteLimit, 1, attributedObservations(t, createObservations(fixtures)...), bf)
	assert.Len(t, outcomes.Outcomes, 4, "budget admits exactly four creates")
	pq, err := rs.GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Empty(t, pq, "deferred items leave the committed queue and retry via re-broadcast")

	// Round 2: the two deferred items re-enter the queue via blob handles.
	handles := &vaultcommon.Observations{
		PendingQueueItems: [][]byte{{0}, {1}},
		SortNonce:         testSortNonce(),
	}
	handleObs, err := proto.Marshal(handles)
	require.NoError(t, err)
	aos := []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(handleObs)},
		{Observer: 1, Observation: types.Observation(handleObs)},
		{Observer: 2, Observation: types.Observation(handleObs)},
	}
	outcomes = budgetRound(t, r, underlying, keyLimit, byteLimit, 2, aos, bf)
	assert.Empty(t, outcomes.Outcomes, "round 2 only re-ingests; nothing is processable yet")
	pq, err = rs.GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Len(t, pq, 2, "both deferred items re-enter the queue")

	// Round 3: the re-ingested items process within the now-shallow floor.
	outcomes = budgetRound(t, r, underlying, keyLimit, byteLimit, 3, attributedObservations(t, createObservations(fixtures[4:6])...), bf)
	assert.Len(t, outcomes.Outcomes, 2, "re-ingested items process once observable")
	pq, err = rs.GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Empty(t, pq)

	// The backlog fully drained: all six secrets exist.
	for _, f := range fixtures {
		secret, err := rs.GetSecret(t.Context(), f.id)
		require.NoError(t, err)
		assert.NotNil(t, secret, "secret for %s", vaulttypes.KeyFor(f.id))
	}
}

// TestStateTransition_KVBudgetGate_PacksIngestAgainstConsumedBudget is the
// stacked worst-case regression: the same round both processes a full old
// queue and ingests new items. The packer must trim the ingest against the
// budget already consumed by processing so the enforcing writer never trips.
func TestStateTransition_KVBudgetGate_PacksIngestAgainstConsumedBudget(t *testing.T) {
	t.Parallel()
	// keyLimit 18: Phase 1 floor (5) plus four creates (8 keys) consumes 8;
	// the packer's remaining 18-8-4=6 keys cap the ingest prefix at 5 of 8.
	const keyLimit = 18
	const byteLimit = 100000

	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withKVWriteBudgetLimits(keyLimit, byteLimit))

	underlying := &kv{m: map[string]response{}}
	oldFixtures := newBudgetCreateFixtures(t, pk, "old", 4)
	oldItems := make([]*vaultcommon.StoredPendingQueueItem, len(oldFixtures))
	for i, f := range oldFixtures {
		oldItems[i] = f.item
	}
	require.NoError(t, newTestWriteStore(t, underlying).WritePendingQueue(t.Context(), oldItems))

	newFixtures := newBudgetCreateFixtures(t, pk, "new", 8)
	newItems := make([]*vaultcommon.StoredPendingQueueItem, len(newFixtures))
	blobs := make([][]byte, len(newFixtures))
	for i, f := range newFixtures {
		newItems[i] = f.item
		blobs[i] = protoMarshal(t, f.item)
	}

	// Observations carry both the processing contributions for the old queue
	// and blob handles ingesting the new items.
	obs := &vaultcommon.Observations{
		Observations:      make([]*vaultcommon.Observation, 0, len(oldFixtures)),
		PendingQueueItems: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}},
		SortNonce:         testSortNonce(),
	}
	for _, f := range oldFixtures {
		obs.Observations = append(obs.Observations, &vaultcommon.Observation{
			Id:          vaulttypes.KeyFor(f.id),
			RequestType: vaultcommon.RequestType_CREATE_SECRETS,
			Request:     &vaultcommon.Observation_CreateSecretsRequest{CreateSecretsRequest: f.req},
			Response:    &vaultcommon.Observation_CreateSecretsResponse{CreateSecretsResponse: f.resp},
		})
	}
	obsb, err := proto.Marshal(obs)
	require.NoError(t, err)
	aos := []types.AttributedObservation{
		{Observer: 0, Observation: types.Observation(obsb)},
		{Observer: 1, Observation: types.Observation(obsb)},
		{Observer: 2, Observation: types.Observation(obsb)},
	}

	bf := &blobber{blobs: blobs}
	r.unmarshalBlob = bf.unmarshalBlob

	outcomes := budgetRound(t, r, underlying, keyLimit, byteLimit, 1, aos, bf)
	assert.Len(t, outcomes.Outcomes, 4, "the old queue fully processes within the floor reserve")

	pq, err := newTestReadStore(t, underlying).GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Len(t, pq, 5, "ingest is packed into the remaining budget; the tail defers to the next round")
}

// TestStateTransition_KVBudgetGate_UnderRealizedItemKeepsTrackerExact proves
// the one-item lookahead never under-reserves and the tracker stays exact: a
// create whose secret already exists passes the gate (its payload bounds two
// keys) but realizes zero writes during processing, and the round still
// succeeds with an exact consumed total.
func TestStateTransition_KVBudgetGate_UnderRealizedItemKeepsTrackerExact(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withKeys(pk, shares[0]), withKVWriteBudgetLimits(50, 100000))

	underlying := &kv{m: map[string]response{}}
	fixtures := newBudgetCreateFixtures(t, pk, "under", 2)
	fresh, existing := fixtures[0], fixtures[1]

	// The second item's secret already exists (both the record and its
	// metadata entry), so processing it writes nothing and returns a user
	// error — even though its observation, formed before the state was
	// checked, carries no error.
	seed := newTestWriteStore(t, underlying)
	require.NoError(t, seed.WriteSecret(t.Context(), existing.id, &vaultcommon.StoredSecret{EncryptedSecret: []byte("existing")}))
	require.NoError(t, seed.WriteMetadata(t.Context(), existing.id.Owner, &vaultcommon.StoredMetadata{
		SecretIdentifiers: []*vaultcommon.SecretIdentifier{existing.id},
	}))
	require.NoError(t, seed.WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{fresh.item, existing.item}))

	outcomes := budgetRound(t, r, underlying, 50, 100000, 1, attributedObservations(
		t,
		observation{id: fresh.id, req: fresh.req, resp: fresh.resp},
		observation{id: existing.id, req: existing.req, resp: existing.resp},
	), nil)
	require.Len(t, outcomes.Outcomes, 2)
	assert.True(t, outcomes.Outcomes[0].GetCreateSecretsResponse().Responses[0].Success)
	assert.False(t, outcomes.Outcomes[1].GetCreateSecretsResponse().Responses[0].Success)
	assert.Contains(t, outcomes.Outcomes[1].GetCreateSecretsResponse().Responses[0].GetError(), "key already exists")

	secret, err := newTestReadStore(t, underlying).GetSecret(t.Context(), fresh.id)
	require.NoError(t, err)
	assert.NotNil(t, secret)
}

// budgetDeleteFixture bundles a round's worth of delete-batch state: a
// committed queue item for a full-batch delete of one owner's secrets, the
// observations that make it processable, and the item's worst-case bound.
// The plugin is built per subtest so its budget limits match the derived
// byte limit exactly.
type budgetDeleteFixture struct {
	pk        *tdh2easy.PublicKey
	share     *tdh2easy.PrivateShare
	kv        *kv
	ids       []*vaultcommon.SecretIdentifier
	req       *vaultcommon.DeleteSecretsRequest
	resp      *vaultcommon.DeleteSecretsResponse
	item      *vaultcommon.StoredPendingQueueItem
	projected kvWriteCost
}

func newBudgetDeleteFixture(t *testing.T, maxSecretsPerOwner int) budgetDeleteFixture {
	t.Helper()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)

	underlying := &kv{m: map[string]response{}}
	seed := newTestWriteStore(t, underlying)

	ids := make([]*vaultcommon.SecretIdentifier, vaulttypes.MaxBatchSize)
	for i := range ids {
		ids[i] = &vaultcommon.SecretIdentifier{
			Owner:     "delowner",
			Namespace: "main",
			Key:       fmt.Sprintf("del_key_%d", i),
		}
		require.NoError(t, seed.WriteSecret(t.Context(), ids[i], &vaultcommon.StoredSecret{EncryptedSecret: []byte("ciphertext")}))
	}

	req := &vaultcommon.DeleteSecretsRequest{RequestId: "del_req_1", Ids: ids}
	resp := &vaultcommon.DeleteSecretsResponse{}
	for _, id := range ids {
		resp.Responses = append(resp.Responses, &vaultcommon.DeleteSecretResponse{Id: id})
	}
	anyReq, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: vaulttypes.KeyFor(ids[0]), Item: anyReq}
	require.NoError(t, seed.WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{item}))

	wc := projectedProcessedWriteCost(t.Context(), newTestReadStore(t, underlying), item, maxSecretsPerOwner)
	require.Positive(t, wc.bytes)
	return budgetDeleteFixture{
		pk:        pk,
		share:     shares[0],
		kv:        underlying,
		ids:       ids,
		req:       req,
		resp:      resp,
		item:      item,
		projected: wc,
	}
}

// deleteBudgetRound runs one StateTransition round for the delete fixture at
// the given byte limit, with the plugin's budget gate and the enforcing
// writer both configured from that limit.
func deleteBudgetRound(t *testing.T, f budgetDeleteFixture, maxSecretsPerOwner, keyLimit, byteLimit int) *vaultcommon.Outcomes {
	t.Helper()
	r := newTestReportingPlugin(t,
		withKeys(f.pk, f.share),
		withMaxSecretsPerOwner(maxSecretsPerOwner),
		withKVWriteBudgetLimits(keyLimit, byteLimit),
	)
	return budgetRound(t, r, f.kv, keyLimit, byteLimit, 1,
		attributedObservations(t, observation{id: f.ids[0], req: f.req, resp: f.resp}), nil)
}

// TestStateTransition_KVBudgetGate_DeleteBatchBoundsRealizedWrites is the
// delete-batch regression: a full-batch delete processes at exactly the
// gate's projection (rewrite floor plus the exact projected write set)
// without tripping the enforcing writer, and one byte less defers the batch.
func TestStateTransition_KVBudgetGate_DeleteBatchBoundsRealizedWrites(t *testing.T) {
	t.Parallel()
	const maxSecretsPerOwner = 10
	const keyLimit = 50

	t.Run("processes at the projected bound", func(t *testing.T) {
		t.Parallel()
		f := newBudgetDeleteFixture(t, maxSecretsPerOwner)
		byteLimit := mandatoryPendingQueueRewriteFloor(1).add(f.projected).bytes

		outcomes := deleteBudgetRound(t, f, maxSecretsPerOwner, keyLimit, byteLimit)
		require.Len(t, outcomes.Outcomes, 1)
		delResp := outcomes.Outcomes[0].GetDeleteSecretsResponse()
		require.Len(t, delResp.Responses, len(f.ids))
		for _, dr := range delResp.Responses {
			assert.True(t, dr.Success, dr.GetError())
		}

		rs := newTestReadStore(t, f.kv)
		for _, id := range f.ids {
			secret, err := rs.GetSecret(t.Context(), id)
			require.NoError(t, err)
			assert.Nil(t, secret, "secret for %s", vaulttypes.KeyFor(id))
		}
		pq, err := rs.GetPendingQueue(t.Context())
		require.NoError(t, err)
		assert.Empty(t, pq)
	})

	t.Run("defers one byte below the projection", func(t *testing.T) {
		t.Parallel()
		f := newBudgetDeleteFixture(t, maxSecretsPerOwner)
		byteLimit := mandatoryPendingQueueRewriteFloor(1).add(f.projected).bytes - 1

		outcomes := deleteBudgetRound(t, f, maxSecretsPerOwner, keyLimit, byteLimit)
		assert.Empty(t, outcomes.Outcomes, "the delete batch defers to the next round")

		rs := newTestReadStore(t, f.kv)
		for _, id := range f.ids {
			secret, err := rs.GetSecret(t.Context(), id)
			require.NoError(t, err)
			assert.NotNil(t, secret, "secret for %s", vaulttypes.KeyFor(id))
		}
		pq, err := rs.GetPendingQueue(t.Context())
		require.NoError(t, err)
		assert.Empty(t, pq, "deferred item leaves the committed queue and retries via re-broadcast")
	})
}

// TestPurgeStalledPendingQueue_FitsAtMaxPackedDepth pins the depth invariant:
// the packer bounds the committed queue depth by the configured key budget
// minus the reserve, so a queue at that depth can always be purged with a
// delete-only write set inside the same key budget.
func TestPurgeStalledPendingQueue_FitsAtMaxPackedDepth(t *testing.T) {
	t.Parallel()
	const keyLimit = 300
	const byteLimit = 1468006

	kept := make([]*vaultcommon.StoredPendingQueueItem, 500)
	for i := range kept {
		anyPayload, err := anypb.New(&vaultcommon.ListSecretIdentifiersRequest{RequestId: "r", Owner: "o"})
		require.NoError(t, err)
		kept[i] = &vaultcommon.StoredPendingQueueItem{Id: fmt.Sprintf("purge-item-%d", i), Item: anyPayload}
	}
	// The packer's max depth under the default key budget: remaining =
	// 300 - 0 - 4 = 296 keys, so the packed prefix caps at 295.
	depth := packPendingQueueWithinKVBudget(0, kept, kvWriteCost{}, keyLimit, byteLimit)
	require.Positive(t, depth)
	require.Less(t, depth, keyLimit, "the packed depth must stay within the key budget")

	underlying := &kv{m: map[string]response{}}
	require.NoError(t, newTestWriteStore(t, underlying).WritePendingQueue(t.Context(), kept[:depth]))

	r := newTestReportingPlugin(t)
	rw := newLimitCheckingKV(underlying, keyLimit, byteLimit)
	store := NewWriteStore(rw, newTestMetrics(t))

	_, err := r.purgeStalledPendingQueue(t.Context(), logger.TestSugared(t), store, 2)
	require.NoError(t, err, "delete-only purge at the packer's max depth must fit the key budget")

	pq, err := newTestReadStore(t, underlying).GetPendingQueue(t.Context())
	require.NoError(t, err)
	assert.Empty(t, pq)
}
