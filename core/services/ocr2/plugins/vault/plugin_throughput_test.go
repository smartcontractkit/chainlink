package vault

// TestPlugin_ThroughputAnalysis measures the maximum number of requests we can
// process per OCR round for each request type before hitting the observation
// (512 KB) or state-transition (5 MB) size limits.
//
// It uses the production vault DON parameters: n=10, f=3, byzQuorumSize=7.
// All inputs are padded to their maximum allowed size so the results are
// conservative (worst-case) bounds.
//
// Run with -v to see the full size table:
//
//	go test -v -run TestPlugin_ThroughputAnalysis ./core/services/ocr2/plugins/vault/...

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
)

func TestPlugin_ThroughputAnalysis(t *testing.T) {
	// --- Vault DON parameters ---
	const (
		donN           = 10
		donF           = 3
		byzQuorumSize  = 2*donF + 1       // 7 — chosen observations for GetSecrets
		fPlusOne       = donF + 1         // 4 — chosen observations for all other types
		maxBatchSize   = 12               // VaultPluginBatchSizeLimit
		maxPendingBlob = 2 * maxBatchSize // max blob *handles* per Observation(); each handle may cover many requests via PendingQueueBlobItems
	)

	// --- Per-secret limits (from cresettings defaults and plugin config) ---
	const (
		maxSecretsPerRequest = 5          // VaultRequestBatchSizeLimit (GetSecrets cap)
		maxEncryptionKeysWC  = 2*donF + 1 // max encryption keys per secret request
		maxCiphertextBytes   = 2000       // MaxCiphertextLengthBytes
		// Shares are base64 strings.
		// Current base64 string limit = 600 bytes → raw binary = 600 * 3/4 = 450 bytes.
		maxShareBytes      = 600
		maxIdentifierBytes = 64 // key / namespace / owner each
	)

	// --- Size limits under test ---
	const (
		maxObservationBytes     = 512 * 1024      // 512 KB
		maxStateTransitionBytes = 5 * 1024 * 1024 // 5 MB
	)

	// --- BlobHandle wire size for this DON (n=10, f=3, byzQuorumSize=7) ---
	// Observation() packs local-queue items into up to maxPendingBlob blob payloads (legacy one item per blob or batched).
	// bytes in PendingQueueItems (not the full StoredPendingQueueItem).
	// Breakdown (see blobtypes.BlobHandleMarshalledBytesUpperBound):
	//   1 B  variant tag
	//   8 B  proto overhead
	//  32 B  ChunkDigestsRoot (mt.Digest = [32]byte)
	//   8 B  PayloadLength
	//   8 B  ExpirySeqNr
	//   8 B  Submitter
	//   7 × (8 proto + 64 ed25519 + 8 signer) = 7 × 80 = 560 B  availability sigs
	// total = 625 B per handle
	const blobHandleBytes = 625

	// --- Helper: max-size secret identifier ---
	makeID := func(i int) *vaultcommon.SecretIdentifier {
		return &vaultcommon.SecretIdentifier{
			Owner:     strings.Repeat("o", maxIdentifierBytes),
			Namespace: strings.Repeat("n", maxIdentifierBytes),
			Key:       fmt.Sprintf("%s%d", strings.Repeat("k", maxIdentifierBytes-2), i),
		}
	}

	// --- Helper: max-size encrypted value (hex of 2000 bytes → 4000 chars) ---
	maxEncValue := strings.Repeat("a", maxCiphertextBytes*2) // hex doubles the size

	// --- Helper: max-size share string ---
	maxShare := strings.Repeat("s", maxShareBytes)

	// --- Helper: max-size encryption key (hex of 32-byte curve25519 key) ---
	maxEncKey := func(i int) string {
		return fmt.Sprintf("%s%02d", strings.Repeat("e", maxIdentifierBytes-2), i)
	}

	// -------------------------------------------------------------------------
	// Observation builder: GetSecrets
	//
	// Production code (observeGetSecrets) sets only RequestType + Response —
	// the Request field is NOT included in GetSecrets observations.
	//
	// Response per secret: identifier + encrypted_value + numEncKeys×1 share.
	// -------------------------------------------------------------------------
	buildGetSecretsObsEntry := func(reqID string, numSecrets, numEncKeys int) *vaultcommon.Observation {
		resps := make([]*vaultcommon.SecretResponse, numSecrets)
		for i := range numSecrets {
			shares := make([]*vaultcommon.EncryptedShares, numEncKeys)
			for k := range numEncKeys {
				shares[k] = &vaultcommon.EncryptedShares{
					EncryptionKey: maxEncKey(k),
					Shares:        []string{maxShare},
				}
			}
			resps[i] = &vaultcommon.SecretResponse{
				Id: makeID(i),
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue:               maxEncValue,
						EncryptedDecryptionKeyShares: shares,
					},
				},
			}
		}
		return &vaultcommon.Observation{
			Id:          reqID,
			RequestType: vaultcommon.RequestType_GET_SECRETS,
			Response: &vaultcommon.Observation_GetSecretsResponse{
				GetSecretsResponse: &vaultcommon.GetSecretsResponse{Responses: resps},
			},
		}
	}

	// -------------------------------------------------------------------------
	// StateTransition outcome builder: GetSecrets
	//
	// stateTransitionGetSecrets aggregates shares from byzQuorumSize chosen
	// observations. Each chosen observation contributes 1 share per encryption
	// key, so the outcome has byzQuorumSize shares per (secret, enc-key) pair.
	// -------------------------------------------------------------------------
	buildGetSecretsOutcomeEntry := func(reqID string, numSecrets, numEncKeys int) *vaultcommon.Outcome {
		resps := make([]*vaultcommon.SecretResponse, numSecrets)
		for i := range numSecrets {
			aggShares := make([]*vaultcommon.EncryptedShares, numEncKeys)
			for k := range numEncKeys {
				shares := make([]string, byzQuorumSize)
				for j := range byzQuorumSize {
					shares[j] = fmt.Sprintf("%s%d", maxShare[:maxShareBytes-2], j)
				}
				aggShares[k] = &vaultcommon.EncryptedShares{
					EncryptionKey: maxEncKey(k),
					Shares:        shares,
				}
			}
			resps[i] = &vaultcommon.SecretResponse{
				Id: makeID(i),
				Result: &vaultcommon.SecretResponse_Data{
					Data: &vaultcommon.SecretData{
						EncryptedValue:               maxEncValue,
						EncryptedDecryptionKeyShares: aggShares,
					},
				},
			}
		}
		return &vaultcommon.Outcome{
			Id:          reqID,
			RequestType: vaultcommon.RequestType_GET_SECRETS,
			Response: &vaultcommon.Outcome_GetSecretsResponse{
				GetSecretsResponse: &vaultcommon.GetSecretsResponse{Responses: resps},
			},
		}
	}

	// -------------------------------------------------------------------------
	// Observation builder: CreateSecrets / UpdateSecrets
	//
	// Production code echoes the full request (including ciphertext) AND a
	// simple success/error response into the observation.
	// -------------------------------------------------------------------------
	buildCreateSecretsObsEntry := func(reqID string, numSecrets int) *vaultcommon.Observation {
		secrets := make([]*vaultcommon.EncryptedSecret, numSecrets)
		resps := make([]*vaultcommon.CreateSecretResponse, numSecrets)
		for i := range numSecrets {
			secrets[i] = &vaultcommon.EncryptedSecret{Id: makeID(i), EncryptedValue: maxEncValue}
			resps[i] = &vaultcommon.CreateSecretResponse{Id: makeID(i), Success: true}
		}
		return &vaultcommon.Observation{
			Id:          reqID,
			RequestType: vaultcommon.RequestType_CREATE_SECRETS,
			Request: &vaultcommon.Observation_CreateSecretsRequest{
				CreateSecretsRequest: &vaultcommon.CreateSecretsRequest{
					RequestId:        reqID,
					OrgId:            strings.Repeat("o", maxIdentifierBytes),
					WorkflowOwner:    strings.Repeat("w", maxIdentifierBytes),
					EncryptedSecrets: secrets,
				},
			},
			Response: &vaultcommon.Observation_CreateSecretsResponse{
				CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{Responses: resps},
			},
		}
	}

	// -------------------------------------------------------------------------
	// StateTransition outcome builder: CreateSecrets / UpdateSecrets
	//
	// stateTransitionCreateSecrets writes ciphertexts to the KV store and
	// returns only simple success/error responses — no ciphertext in the outcome.
	// -------------------------------------------------------------------------
	buildCreateSecretsOutcomeEntry := func(reqID string, numSecrets int) *vaultcommon.Outcome {
		resps := make([]*vaultcommon.CreateSecretResponse, numSecrets)
		for i := range numSecrets {
			resps[i] = &vaultcommon.CreateSecretResponse{Id: makeID(i), Success: true}
		}
		return &vaultcommon.Outcome{
			Id:          reqID,
			RequestType: vaultcommon.RequestType_CREATE_SECRETS,
			Response: &vaultcommon.Outcome_CreateSecretsResponse{
				CreateSecretsResponse: &vaultcommon.CreateSecretsResponse{Responses: resps},
			},
		}
	}

	// -------------------------------------------------------------------------
	// Shared helper: build an Observations proto for a batch of entries plus
	// a simulated pending queue. pendingItemBytes is the byte length of each
	// repeated PendingQueueItems entry — either marshalled BlobHandle (~625 B
	// with BroadcastBlob) or a full proto-marshalled StoredPendingQueueItem.
	// -------------------------------------------------------------------------
	buildObservations := func(entries []*vaultcommon.Observation, pendingItemBytes int) ([]byte, error) {
		pendingQueue := make([][]byte, maxPendingBlob)
		for i := range pendingQueue {
			pendingQueue[i] = make([]byte, pendingItemBytes)
		}
		obs := &vaultcommon.Observations{
			Observations:      entries,
			PendingQueueItems: pendingQueue,
			SortNonce:         make([]byte, 32),
		}
		return proto.Marshal(obs)
	}

	// -------------------------------------------------------------------------
	// Without BroadcastBlob: measure marshalled StoredPendingQueueItem size for
	// the two most expensive request types (for comparison / hypothetical inline
	// pending bytes).
	// -------------------------------------------------------------------------
	pendingItemSizeFor := func(req proto.Message) int {
		anyMsg, err := anypb.New(req)
		require.NoError(t, err)
		item := &vaultcommon.StoredPendingQueueItem{
			Id:   strings.Repeat("i", 36), // UUID-length request ID
			Item: anyMsg,
		}
		b, err := proto.Marshal(item)
		require.NoError(t, err)
		return len(b)
	}

	// Max-size GetSecretsRequest pending queue item
	getSecretsReqs := make([]*vaultcommon.SecretRequest, maxSecretsPerRequest)
	for i := range maxSecretsPerRequest {
		encKeys := make([]string, maxEncryptionKeysWC)
		for k := range maxEncryptionKeysWC {
			encKeys[k] = maxEncKey(k)
		}
		getSecretsReqs[i] = &vaultcommon.SecretRequest{
			Id:             makeID(i),
			EncryptionKeys: encKeys,
		}
	}
	getSecretsPendingItemSize := pendingItemSizeFor(&vaultcommon.GetSecretsRequest{
		Requests:      getSecretsReqs,
		OrgId:         strings.Repeat("o", maxIdentifierBytes),
		WorkflowOwner: strings.Repeat("w", maxIdentifierBytes),
	})

	// Max-size CreateSecretsRequest pending queue item
	createSecrets := make([]*vaultcommon.EncryptedSecret, maxSecretsPerRequest)
	for i := range maxSecretsPerRequest {
		createSecrets[i] = &vaultcommon.EncryptedSecret{Id: makeID(i), EncryptedValue: maxEncValue}
	}
	createSecretsPendingItemSize := pendingItemSizeFor(&vaultcommon.CreateSecretsRequest{
		RequestId:        strings.Repeat("r", 36),
		OrgId:            strings.Repeat("o", maxIdentifierBytes),
		WorkflowOwner:    strings.Repeat("w", maxIdentifierBytes),
		EncryptedSecrets: createSecrets,
	})

	t.Logf("Pending queue item sizes:")
	t.Logf("  BlobHandle (BroadcastBlob path):    %d bytes", blobHandleBytes)
	t.Logf("  GetSecretsRequest inline item:      %d bytes", getSecretsPendingItemSize)
	t.Logf("  CreateSecretsRequest inline item:   %d bytes", createSecretsPendingItemSize)

	buildOutcomes := func(entries []*vaultcommon.Outcome) ([]byte, error) {
		return proto.Marshal(&vaultcommon.Outcomes{Outcomes: entries})
	}

	mustMarshal := func(t *testing.T, fn func() ([]byte, error)) int {
		t.Helper()
		b, err := fn()
		require.NoError(t, err)
		return len(b)
	}

	type result struct {
		batchSize int
		obsBytes  int
		stBytes   int
	}

	printTable := func(t *testing.T, name string, results []result, obsLimit, stLimit int) {
		t.Helper()
		t.Logf("\n=== %s ===", name)
		t.Logf("  %-10s %-20s %-24s %-20s %-24s",
			"batch", "obs bytes", "obs (KB / limit KB)", "ST bytes", "ST (KB / limit KB)")
		for _, r := range results {
			obsKB := float64(r.obsBytes) / 1024
			stKB := float64(r.stBytes) / 1024
			obsLimitKB := float64(obsLimit) / 1024
			stLimitKB := float64(stLimit) / 1024
			obsOK := "✓"
			if r.obsBytes > obsLimit {
				obsOK = "✗ EXCEEDS"
			}
			stOK := "✓"
			if r.stBytes > stLimit {
				stOK = "✗ EXCEEDS"
			}
			t.Logf("  %-10d %-20d %.1f / %.1f KB %s    %-20d %.1f / %.1f KB %s",
				r.batchSize,
				r.obsBytes, obsKB, obsLimitKB, obsOK,
				r.stBytes, stKB, stLimitKB, stOK)
		}
	}

	// runBatchSweep sweeps batch sizes 1-10, builds obs+outcome protos, and
	// returns the results. pendingItemBytes controls each PendingQueueItems entry:
	// use blobHandleBytes for the BroadcastBlob path, or full StoredPendingQueueItem
	// size for the inline-payload comparison.
	runBatchSweep := func(
		t *testing.T,
		buildObs func(id string) *vaultcommon.Observation,
		buildOutcome func(id string) *vaultcommon.Outcome,
		pendingItemBytes int,
	) []result {
		t.Helper()
		var results []result
		for batchSize := 1; batchSize <= maxBatchSize; batchSize++ {
			obsEntries := make([]*vaultcommon.Observation, batchSize)
			stEntries := make([]*vaultcommon.Outcome, batchSize)
			for i := range batchSize {
				id := fmt.Sprintf("req-%d", i)
				obsEntries[i] = buildObs(id)
				stEntries[i] = buildOutcome(id)
			}
			results = append(results, result{
				batchSize: batchSize,
				obsBytes:  mustMarshal(t, func() ([]byte, error) { return buildObservations(obsEntries, pendingItemBytes) }),
				stBytes:   mustMarshal(t, func() ([]byte, error) { return buildOutcomes(stEntries) }),
			})
		}
		return results
	}

	summarize := func(t *testing.T, results []result) {
		t.Helper()
		maxObsBatch, maxSTBatch := 0, 0
		for _, r := range results {
			if r.obsBytes <= maxObservationBytes {
				maxObsBatch = r.batchSize
			}
			if r.stBytes <= maxStateTransitionBytes {
				maxSTBatch = r.batchSize
			}
		}
		effectiveMax := ourMin(maxObsBatch, maxSTBatch, maxBatchSize)
		t.Logf("  → obs limit: %d | ST limit: %d | effective max batch: %d", maxObsBatch, maxSTBatch, effectiveMax)
	}

	// =========================================================================
	// WITH BroadcastBlob — PendingQueueItems are compact marshalled BlobHandles
	// (matches production Observation() after broadcastBlobPayloads).
	// =========================================================================
	t.Run("WithBlobBroadcast", func(t *testing.T) {
		t.Run("GetSecrets_5secrets_7enckeys", func(t *testing.T) {
			results := runBatchSweep(t,
				func(id string) *vaultcommon.Observation {
					return buildGetSecretsObsEntry(id, maxSecretsPerRequest, maxEncryptionKeysWC)
				},
				func(id string) *vaultcommon.Outcome {
					return buildGetSecretsOutcomeEntry(id, maxSecretsPerRequest, maxEncryptionKeysWC)
				},
				blobHandleBytes,
			)
			printTable(t, "WITH BlobBroadcast — GetSecrets (5 secrets × 7 enc-keys)", results, maxObservationBytes, maxStateTransitionBytes)
			summarize(t, results)
		})

		t.Run("CreateSecrets_5secrets_2000ByteCiphertext", func(t *testing.T) {
			results := runBatchSweep(t,
				func(id string) *vaultcommon.Observation { return buildCreateSecretsObsEntry(id, maxSecretsPerRequest) },
				func(id string) *vaultcommon.Outcome { return buildCreateSecretsOutcomeEntry(id, maxSecretsPerRequest) },
				blobHandleBytes,
			)
			printTable(t, "WITH BlobBroadcast — CreateSecrets (5 secrets × 2000-byte ciphertext)", results, maxObservationBytes, maxStateTransitionBytes)
			summarize(t, results)
		})
	})

	// =========================================================================
	// WITHOUT BroadcastBlob — each PendingQueueItems entry is a full
	// StoredPendingQueueItem (worst-case size per active batch type).
	// =========================================================================
	// t.Run("WithoutBlobBroadcast", func(t *testing.T) {
	// 	t.Run("GetSecrets_5secrets_7enckeys", func(t *testing.T) {
	// 		results := runBatchSweep(t,
	// 			func(id string) *vaultcommon.Observation {
	// 				return buildGetSecretsObsEntry(id, maxSecretsPerRequest, maxEncryptionKeysWC)
	// 			},
	// 			func(id string) *vaultcommon.Outcome {
	// 				return buildGetSecretsOutcomeEntry(id, maxSecretsPerRequest, maxEncryptionKeysWC)
	// 			},
	// 			getSecretsPendingItemSize,
	// 		)
	// 		printTable(t, "WITHOUT BlobBroadcast — GetSecrets (5 secrets × 7 enc-keys)", results, maxObservationBytes, maxStateTransitionBytes)
	// 		summarize(t, results)
	// 	})

	// 	t.Run("CreateSecrets_5secrets_2000ByteCiphertext", func(t *testing.T) {
	// 		results := runBatchSweep(t,
	// 			func(id string) *vaultcommon.Observation { return buildCreateSecretsObsEntry(id, maxSecretsPerRequest) },
	// 			func(id string) *vaultcommon.Outcome { return buildCreateSecretsOutcomeEntry(id, maxSecretsPerRequest) },
	// 			createSecretsPendingItemSize,
	// 		)
	// 		printTable(t, "WITHOUT BlobBroadcast — CreateSecrets (5 secrets × 2000-byte ciphertext)", results, maxObservationBytes, maxStateTransitionBytes)
	// 		summarize(t, results)
	// 	})
	// })
}

func ourMin(vals ...int) int {
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
