package vault

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"testing"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const (
	fuzzN           = 4
	fuzzF           = 1
	fuzzOwner       = "0x1111111111111111111111111111111111111111"
	fuzzRequest     = "request-1"
	fuzzByzantineID = 3
)

// unavailableBlobs lets a Byzantine observation reference blob handles that no oracle ever broadcast.
type unavailableBlobs struct{}

func (unavailableBlobs) FetchBlob(context.Context, ocr3_1types.BlobHandle) ([]byte, error) {
	return nil, errors.New("blob unavailable")
}

func fuzzPendingItem(t *testing.T, pk *tdh2easy.PublicKey, kind uint8, owner, namespace, key string, secret []byte, corrupt bool) proto.Message {
	t.Helper()
	id := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: namespace, Key: key}

	encrypt := func() string {
		var label [32]byte
		copy(label[12:], common.HexToAddress(owner).Bytes())
		ct, err := tdh2easy.EncryptWithLabel(pk, secret, label)
		require.NoError(t, err)
		raw, err := ct.Marshal()
		require.NoError(t, err)
		if corrupt && len(raw) > 0 {
			raw[len(raw)/2] ^= 0xff
		}
		return hex.EncodeToString(raw)
	}

	switch kind % 5 {
	case 0:
		return &vaultcommon.GetSecretsRequest{Requests: []*vaultcommon.SecretRequest{{Id: id}}}
	case 1:
		return &vaultcommon.CreateSecretsRequest{EncryptedSecrets: []*vaultcommon.EncryptedSecret{{Id: id, EncryptedValue: encrypt()}}}
	case 2:
		return &vaultcommon.UpdateSecretsRequest{EncryptedSecrets: []*vaultcommon.EncryptedSecret{{Id: id, EncryptedValue: encrypt()}}}
	case 3:
		return &vaultcommon.DeleteSecretsRequest{Ids: []*vaultcommon.SecretIdentifier{id}}
	default:
		return &vaultcommon.ListSecretIdentifiersRequest{Owner: owner, Namespace: namespace}
	}
}

func FuzzReportingPlugin_StateTransitionAcrossOracles(f *testing.F) {
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(f, err)

	f.Add(uint64(1), uint8(0), true, fuzzOwner, "main", "secret", []byte("value"), false, false, false, []byte{})
	f.Add(uint64(1), uint8(1), false, fuzzOwner, "main", "secret", []byte("value"), false, false, false, []byte{})
	f.Add(uint64(2), uint8(2), true, fuzzOwner, "main", "secret", []byte("value"), true, false, false, []byte{0x0a, 0x00})
	f.Add(uint64(3), uint8(3), true, fuzzOwner, "main", "secret", []byte{}, false, false, false, []byte{0xff})
	f.Add(uint64(4), uint8(4), true, "owner", "", "", []byte{}, false, false, false, []byte{})
	// A Byzantine oracle sends a well-formed observation computed from the opposite KV state, so the oracles disagree;
	// with a second oracle doing the same, two different observations each have F+1 supporters.
	for kind := range uint8(5) {
		for _, preexisting := range []bool{true, false} {
			f.Add(uint64(6), kind, preexisting, fuzzOwner, "main", "secret", []byte("value"), false, true, false, []byte{})
			f.Add(uint64(6), kind, preexisting, fuzzOwner, "main", "secret", []byte("value"), false, true, true, []byte{})
		}
	}
	danglingBlob, err := proto.Marshal(&vaultcommon.Observations{
		SortNonce:         bytes.Repeat([]byte{'a'}, sortNonceLength),
		PendingQueueItems: [][]byte{{0x00}},
	})
	require.NoError(f, err)
	f.Add(uint64(5), uint8(0), true, fuzzOwner, "main", "secret", []byte("value"), false, false, false, danglingBlob)

	f.Fuzz(func(t *testing.T, seqNr uint64, kind uint8, preexisting bool, owner, namespace, key string, secret []byte, corrupt, byzantineFlipsState, honestOracleFlipsState bool, byzantine []byte) {
		if !utf8.ValidString(owner) || !utf8.ValidString(namespace) || !utf8.ValidString(key) {
			t.Skip("proto3 strings must be valid UTF-8")
		}
		if len(secret) == 0 {
			secret = []byte("x")
		}
		item := fuzzPendingItem(t, pk, kind, owner, namespace, key, secret, corrupt)
		anyItem, err := anypb.New(item)
		require.NoError(t, err)

		id := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: namespace, Key: key}
		kvWith := func(secretExists bool) *kv {
			rdr := &kv{m: make(map[string]response)}
			w := newTestWriteStore(t, rdr)
			require.NoError(t, w.WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{{Id: fuzzRequest, Item: anyItem}}))
			if secretExists {
				require.NoError(t, w.WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{EncryptedSecret: []byte("ciphertext")}))
				require.NoError(t, w.WriteMetadata(t.Context(), owner, &vaultcommon.StoredMetadata{SecretIdentifiers: []*vaultcommon.SecretIdentifier{id}}))
			}
			return rdr
		}
		freshKV := func() *kv { return kvWith(preexisting) }
		newPlugin := func(i int) *ReportingPlugin {
			return newTestReportingPlugin(t, withKeys(pk, shares[i]), withOnchainCfg(fuzzN, fuzzF), withLggr(logger.Nop()), withBatchSize(10))
		}

		aos := make([]types.AttributedObservation, 0, fuzzN)
		for i := range fuzzN - 1 {
			state := freshKV()
			if honestOracleFlipsState && i == fuzzN-2 {
				state = kvWith(!preexisting)
			}
			data, err := newPlugin(i).Observation(t.Context(), seqNr, types.AttributedQuery{}, state, &blobber{})
			require.NoError(t, err, "an honest oracle must be able to observe a well-formed pending request")
			aos = append(aos, types.AttributedObservation{Observer: commontypes.OracleID(i), Observation: data})
		}
		if byzantineFlipsState {
			data, err := newPlugin(fuzzByzantineID-1).Observation(t.Context(), seqNr, types.AttributedQuery{}, kvWith(!preexisting), &blobber{})
			require.NoError(t, err)
			byzantine = data
		}
		aos = append(aos, types.AttributedObservation{Observer: fuzzByzantineID, Observation: byzantine})

		validator := newPlugin(0)
		var validated []types.AttributedObservation
		for i, ao := range aos {
			err := validator.ValidateObservation(t.Context(), seqNr, types.AttributedQuery{}, ao, freshKV(), unavailableBlobs{})
			if i < fuzzN-1 && !honestOracleFlipsState {
				require.NoError(t, err, "honest observation must validate")
			}
			if err == nil {
				validated = append(validated, ao)
			}
		}

		type result struct {
			precursor ocr3_1types.ReportsPlusPrecursor
			errMsg    string
			kv        map[string]response
		}
		run := func(oracle int) result {
			rdr := freshKV()
			p, err := newPlugin(oracle).StateTransition(t.Context(), seqNr, types.AttributedQuery{}, validated, rdr, unavailableBlobs{})
			res := result{precursor: p, kv: rdr.m}
			if err != nil {
				res.errMsg = err.Error()
			}
			return res
		}

		first := run(0)
		// Several runs per oracle so a dependence on map iteration order has a chance to show.
		for range 3 {
			for oracle := range fuzzN - 1 {
				next := run(oracle)
				require.Equal(t, first.errMsg, next.errMsg)
				require.Equal(t, first.precursor, next.precursor, "StateTransition must be deterministic")
				require.Equal(t, first.kv, next.kv, "KV writes must be deterministic")
			}
		}
		if first.errMsg != "" {
			return
		}

		want, wantErr := newPlugin(0).Reports(t.Context(), seqNr, first.precursor)
		for oracle := range fuzzN - 1 {
			got, err := newPlugin(oracle).Reports(t.Context(), seqNr, first.precursor)
			require.Equal(t, wantErr == nil, err == nil)
			require.Equal(t, want, got, "Reports must not depend on the oracle's private key share")
		}
	})
}
