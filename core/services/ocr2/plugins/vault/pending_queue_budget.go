package vault

import (
	"context"
	"slices"
	"strconv"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

const (
	deferredByBudgetPhaseProcessing = "processing"
	deferredByBudgetPhaseIngest     = "ingest"

	// kvBudgetReserveKeys and kvBudgetReserveBytes are subtracted from the
	// remaining budget before packing the pending queue, absorbing index-record
	// variance and any prediction drift.
	kvBudgetReserveKeys  = 4
	kvBudgetReserveBytes = 2048
)

// kvWriteCost is a write-set footprint in modified keys and cumulative bytes,
// using the same accounting as libocr's limitCheckWriteSet.
type kvWriteCost struct {
	keys  int
	bytes int
}

func (c kvWriteCost) add(o kvWriteCost) kvWriteCost {
	return kvWriteCost{keys: c.keys + o.keys, bytes: c.bytes + o.bytes}
}

// kvWriteBudgetTracker wraps the round's KeyValueStateReadWriter and counts
// the write set exactly, mirroring libocr's limitCheckWriteSet arithmetic: a
// modification of a key already in the write set only counts its final value;
// a new key costs one key plus key+value bytes; a delete of a new key costs
// one key plus key bytes.
type kvWriteBudgetTracker struct {
	ocr3_1types.KeyValueStateReadWriter

	modified map[string][]byte
	keys     int
	bytes    int
}

func newKVWriteBudgetTracker(w ocr3_1types.KeyValueStateReadWriter) *kvWriteBudgetTracker {
	return &kvWriteBudgetTracker{
		KeyValueStateReadWriter: w,
		modified:                map[string][]byte{},
	}
}

func (t *kvWriteBudgetTracker) modify(key, value []byte) {
	add, sub := 0, 0
	if prev, ok := t.modified[string(key)]; ok {
		add = len(value)
		sub = len(prev)
	} else {
		add = len(key) + len(value)
		t.keys++
	}
	t.bytes += add - sub
	t.modified[string(key)] = value
}

func (t *kvWriteBudgetTracker) Write(key, value []byte) error {
	if err := t.KeyValueStateReadWriter.Write(key, value); err != nil {
		return err
	}
	t.modify(key, value)
	return nil
}

func (t *kvWriteBudgetTracker) Delete(key []byte) error {
	if err := t.KeyValueStateReadWriter.Delete(key); err != nil {
		return err
	}
	t.modify(key, nil)
	return nil
}

func (t *kvWriteBudgetTracker) consumed() kvWriteCost {
	return kvWriteCost{keys: t.keys, bytes: t.bytes}
}

// projectedProcessedWriteCost returns the projected write-set cost of
// processing item during StateTransition, computed exactly from the item's
// payload and the store's current state. The store reads reflect this round's
// earlier writes, so intra-round sequencing — where earlier items' writes
// change what later items do — is accounted for. The projection is an upper
// bound on the realized cost: contributions carrying errors, undecodable
// ciphertext, duplicate or already-present identifiers, and owners over
// capacity all fail during processing before any write.
func projectedProcessedWriteCost(ctx context.Context, store ReadKVStore, item *vaultcommon.StoredPendingQueueItem, maxSecretsPerOwner int) kvWriteCost {
	if item == nil || item.Item == nil {
		return kvWriteCost{}
	}
	payload, err := item.Item.UnmarshalNew()
	if err != nil {
		return kvWriteCost{}
	}
	switch p := payload.(type) {
	case *vaultcommon.CreateSecretsRequest:
		return createSecretsProjectedWriteCost(ctx, store, p, maxSecretsPerOwner)
	case *vaultcommon.UpdateSecretsRequest:
		return updateSecretsProjectedWriteCost(p)
	case *vaultcommon.DeleteSecretsRequest:
		return deleteSecretsProjectedWriteCost(ctx, store, p)
	default:
		// Read-only request types (GetSecrets, ListSecretIdentifiers) and
		// unknown types never modify state.
		return kvWriteCost{}
	}
}

// createSecretsProjectedWriteCost projects the writes of a create request:
// one secret record per capacity-admitted new identifier plus one owner
// metadata rewrite per owner that gains an identifier.
func createSecretsProjectedWriteCost(ctx context.Context, store ReadKVStore, req *vaultcommon.CreateSecretsRequest, maxSecretsPerOwner int) kvWriteCost {
	// TODO: Remove secretsByOwner once we change EncryptedSecrets to inherit the owner from
	// the top level request instead of defining their own owner field
	secretsByOwner := map[string][]*vaultcommon.EncryptedSecret{}
	for _, s := range req.EncryptedSecrets {
		if s == nil || s.Id == nil {
			continue
		}
		secretsByOwner[s.Id.Owner] = append(secretsByOwner[s.Id.Owner], s)
	}

	cost := kvWriteCost{}
	for owner, secrets := range secretsByOwner {
		cost = cost.add(createOwnerProjectedWriteCost(ctx, store, owner, secrets, maxSecretsPerOwner))
	}
	return cost
}

// createOwnerProjectedWriteCost projects one owner's share of a create
// request. Identifiers already present in the owner's metadata and identifiers
// beyond maxSecretsPerOwner error out during processing before any write, so
// they cost nothing; admitted identifiers each write a secret record, and the
// metadata record's intermediate rewrites collapse in the write set to the
// final projected content.
func createOwnerProjectedWriteCost(ctx context.Context, store ReadKVStore, owner string, secrets []*vaultcommon.EncryptedSecret, maxSecretsPerOwner int) kvWriteCost {
	md, err := store.GetMetadata(ctx, owner)
	if err != nil {
		return kvWriteCost{}
	}
	current := md.GetSecretIdentifiers()

	present := map[string]struct{}{}
	for _, id := range current {
		present[vaulttypes.KeyFor(id)] = struct{}{}
	}

	// Identifiers are appended one per secret, deduplicated as processing's
	// request aggregation deduplicates them, each admitted only while the
	// owner stays within maxSecretsPerOwner.
	updated := slices.Clone(current)
	cost := kvWriteCost{}
	for _, s := range secrets {
		key := vaulttypes.KeyFor(s.Id)
		if _, ok := present[key]; ok {
			continue // key already exists; processing errors before any write
		}
		if len(updated)+1 > maxSecretsPerOwner {
			break // owner at capacity; later identifiers fail the same check
		}
		if slices.ContainsFunc(updated, func(id *vaultcommon.SecretIdentifier) bool {
			return vaulttypes.KeyFor(id) == key
		}) {
			continue // duplicate identifier collapses during processing
		}
		updated = append(updated, s.Id)
		cost = cost.add(secretRecordWriteCost(s))
	}
	if cost.keys == 0 {
		return kvWriteCost{}
	}
	return cost.add(kvWriteCost{
		keys:  1,
		bytes: len(metadataPrefix+owner) + proto.Size(&vaultcommon.StoredMetadata{SecretIdentifiers: updated}),
	})
}

// updateSecretsProjectedWriteCost projects the writes of an update request:
// one secret record rewrite per identifier. Updates never touch the owner
// metadata record — the identifier is unchanged — and identifiers whose secret
// does not exist error out during processing before any write, so each
// payload-derived record cost is an upper bound.
func updateSecretsProjectedWriteCost(req *vaultcommon.UpdateSecretsRequest) kvWriteCost {
	cost := kvWriteCost{}
	for _, s := range req.EncryptedSecrets {
		if s == nil || s.Id == nil {
			continue
		}
		cost = cost.add(secretRecordWriteCost(s))
	}
	return cost
}

// deleteSecretsProjectedWriteCost projects the writes of a delete request: one
// secret-record delete per identifier present in the owner's metadata plus one
// metadata rewrite per owner that loses an identifier.
func deleteSecretsProjectedWriteCost(ctx context.Context, store ReadKVStore, req *vaultcommon.DeleteSecretsRequest) kvWriteCost {
	idsByOwner := map[string][]*vaultcommon.SecretIdentifier{}
	for _, id := range req.Ids {
		if id == nil {
			continue
		}
		idsByOwner[id.Owner] = append(idsByOwner[id.Owner], id)
	}

	cost := kvWriteCost{}
	for owner, ids := range idsByOwner {
		cost = cost.add(deleteOwnerProjectedWriteCost(ctx, store, owner, ids))
	}
	return cost
}

// deleteOwnerProjectedWriteCost projects one owner's share of a delete
// request. Identifiers absent from the owner's metadata error out during
// processing before any write, duplicates are gone after the first removal,
// and the metadata record's intermediate rewrites collapse in the write set
// to the final projected content.
func deleteOwnerProjectedWriteCost(ctx context.Context, store ReadKVStore, owner string, ids []*vaultcommon.SecretIdentifier) kvWriteCost {
	md, err := store.GetMetadata(ctx, owner)
	if err != nil {
		return kvWriteCost{}
	}
	current := md.GetSecretIdentifiers()

	inMetadata := map[string]struct{}{}
	for _, id := range current {
		inMetadata[vaulttypes.KeyFor(id)] = struct{}{}
	}

	toDelete := map[string]struct{}{}
	cost := kvWriteCost{}
	for _, id := range ids {
		key := vaulttypes.KeyFor(id)
		if _, ok := inMetadata[key]; !ok {
			continue // absent from metadata; processing errors before any write
		}
		if _, ok := toDelete[key]; ok {
			continue // already removed by a duplicate earlier in the batch
		}
		toDelete[key] = struct{}{}
		cost.keys++                        // the secret record delete
		cost.bytes += len(keyPrefix + key) // libocr doesn't count the ciphertext bytes in the write cost
	}
	if len(toDelete) == 0 {
		return kvWriteCost{}
	}

	remaining := make([]*vaultcommon.SecretIdentifier, 0, len(current))
	for _, id := range current {
		if _, ok := toDelete[vaulttypes.KeyFor(id)]; !ok {
			remaining = append(remaining, id)
		}
	}
	return cost.add(kvWriteCost{
		keys:  1,
		bytes: len(metadataPrefix+owner) + proto.Size(&vaultcommon.StoredMetadata{SecretIdentifiers: remaining}),
	})
}

// secretRecordWriteCost returns the exact write-set cost of one secret record:
// its KV key plus the marshaled StoredSecret carrying the decoded ciphertext.
func secretRecordWriteCost(s *vaultcommon.EncryptedSecret) kvWriteCost {
	// EncryptedValue is hex-encoded on the wire; the stored record holds the
	// decoded ciphertext.
	ciphertext := len(s.EncryptedValue) / 2
	return kvWriteCost{
		keys:  1,
		bytes: len(keyPrefix+vaulttypes.KeyFor(s.Id)) + storedSecretWireSize(ciphertext),
	}
}

// storedSecretWireSize returns the exact marshaled size of a StoredSecret
// carrying n bytes of ciphertext. The single bytes field is omitted entirely
// when empty.
func storedSecretWireSize(n int) int {
	if n == 0 {
		return 0
	}
	return 1 + protowire.SizeVarint(uint64(n)) + n
}

// mandatoryPendingQueueRewriteFloor returns the exact minimum write-set cost
// of the pending queue rewrite in this round: WritePendingQueue always deletes
// all old positional keys and writes the index, even when it stores nothing.
// Phase 1 reserves this floor before admitting any processed write.
// libocr only counts the key and value bytes of the index record not the
// StoredPendingQueueItem records.
func mandatoryPendingQueueRewriteFloor(oldCount int) kvWriteCost {
	keyBytes := 0
	for j := range oldCount {
		keyBytes += pendingQueueItemKeyLen(j)
	}
	return kvWriteCost{
		keys:  oldCount + 1,
		bytes: keyBytes + len(pendingQueueIndex) + proto.Size(&vaultcommon.StoredPendingQueueIndex{}),
	}
}

// pendingQueueItemKeyLen returns the exact length of the positional pending
// queue item key at index i, matching kvstore.go's encoding.
func pendingQueueItemKeyLen(i int) int {
	return len(pendingQueueItemPrefix) + len(strconv.Itoa(i))
}

// packPendingQueueWithinKVBudget trims kept (already deterministically sorted)
// to the largest prefix whose pending-queue rewrite, performed on top of the
// writes already counted by consumed, fits within the remaining KV write-set
// budgets. Trimmed items are deferred; they remain in node-local queues and
// are re-broadcast next round until request TTL.
func packPendingQueueWithinKVBudget(oldCount int, kept []*vaultcommon.StoredPendingQueueItem, consumed kvWriteCost, keyLimit, byteLimit int) int {
	remainingKeys := keyLimit - consumed.keys - kvBudgetReserveKeys
	remainingBytes := byteLimit - consumed.bytes - kvBudgetReserveBytes
	if remainingKeys <= 0 || remainingBytes <= 0 {
		return 0
	}

	// Prefix sums of positional key lengths over the union of old and new
	// positional keys: delete-only keys beyond the packed prefix still cost
	// their key bytes.
	unionCount := max(oldCount, len(kept))
	positionalKeyBytes := make([]int, unionCount+1)
	for j := range unionCount {
		positionalKeyBytes[j+1] = positionalKeyBytes[j] + pendingQueueItemKeyLen(j)
	}

	packed := 0
	itemBytes := 0
	for packed < len(kept) {
		candidateItemBytes := itemBytes + proto.Size(kept[packed])
		positionalCount := max(oldCount, packed+1)
		candidateBytes := candidateItemBytes +
			positionalKeyBytes[positionalCount] +
			len(pendingQueueIndex) +
			proto.Size(&vaultcommon.StoredPendingQueueIndex{Length: int64(packed + 1)})
		// positionalCount modified item keys plus the index key.
		if positionalCount+1 > remainingKeys || candidateBytes > remainingBytes {
			break
		}
		itemBytes = candidateItemBytes
		packed++
	}
	return packed
}
