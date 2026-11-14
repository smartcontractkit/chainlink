package vault

import (
	"encoding/hex"
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

	f.Add(uint64(1), uint8(0), true, fuzzOwner, "main", "secret", []byte("value"), false, []byte{})
	f.Add(uint64(1), uint8(1), false, fuzzOwner, "main", "secret", []byte("value"), false, []byte{})
	f.Add(uint64(2), uint8(2), true, fuzzOwner, "main", "secret", []byte("value"), true, []byte{0x0a, 0x00})
	f.Add(uint64(3), uint8(3), true, fuzzOwner, "main", "secret", []byte{}, false, []byte{0xff})
	f.Add(uint64(4), uint8(4), true, "owner", "", "", []byte{}, false, []byte{})

	f.Fuzz(func(t *testing.T, seqNr uint64, kind uint8, preexisting bool, owner, namespace, key string, secret []byte, corrupt bool, byzantine []byte) {
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
		freshKV := func() *kv {
			rdr := &kv{m: make(map[string]response)}
			w := newTestWriteStore(t, rdr)
			require.NoError(t, w.WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{{Id: fuzzRequest, Item: anyItem}}))
			if preexisting {
				require.NoError(t, w.WriteSecret(t.Context(), id, &vaultcommon.StoredSecret{EncryptedSecret: []byte("ciphertext")}))
				require.NoError(t, w.WriteMetadata(t.Context(), owner, &vaultcommon.StoredMetadata{SecretIdentifiers: []*vaultcommon.SecretIdentifier{id}}))
			}
			return rdr
		}
		newPlugin := func(i int) *ReportingPlugin {
			return newTestReportingPlugin(t, withKeys(pk, shares[i]), withOnchainCfg(fuzzN, fuzzF), withLggr(logger.Nop()), withBatchSize(10))
		}

		var aos []types.AttributedObservation
		for i := range fuzzN - 1 {
			data, err := newPlugin(i).Observation(t.Context(), seqNr, types.AttributedQuery{}, freshKV(), &blobber{})
			if err != nil {
				return
			}
			aos = append(aos, types.AttributedObservation{Observer: commontypes.OracleID(i), Observation: data})
		}
		aos = append(aos, types.AttributedObservation{Observer: fuzzByzantineID, Observation: byzantine})

		validator := newPlugin(0)
		var validated []types.AttributedObservation
		for i, ao := range aos {
			err := validator.ValidateObservation(t.Context(), seqNr, types.AttributedQuery{}, ao, freshKV(), &blobber{})
			if i < fuzzN-1 {
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
			p, err := newPlugin(oracle).StateTransition(t.Context(), seqNr, types.AttributedQuery{}, validated, rdr, &blobber{})
			res := result{precursor: p, kv: rdr.m}
			if err != nil {
				res.errMsg = err.Error()
			}
			return res
		}

		first := run(0)
		for oracle := range fuzzN - 1 {
			next := run(oracle)
			require.Equal(t, first.errMsg, next.errMsg)
			require.Equal(t, first.precursor, next.precursor, "StateTransition must be deterministic")
			require.Equal(t, first.kv, next.kv, "KV writes must be deterministic")
		}
		if first.errMsg != "" {
			return
		}

		r1, err1 := newPlugin(0).Reports(t.Context(), seqNr, first.precursor)
		r2, err2 := newPlugin(0).Reports(t.Context(), seqNr, first.precursor)
		require.Equal(t, err1 == nil, err2 == nil)
		require.Equal(t, r1, r2, "Reports must be deterministic")
	})
}
