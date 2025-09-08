package vault

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestTransmitter(t *testing.T) {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.DebugLevel)
	store := requests.NewStore[*vaulttypes.Request]()
	handler := requests.NewHandler[*vaulttypes.Request, *vaulttypes.Response](lggr, store, clockwork.NewFakeClock(), time.Hour)
	servicetest.Run(t, handler)
	transmitter := NewTransmitter(lggr, types.Account("0x1"), handler)

	id1 := &vault.SecretIdentifier{
		Owner:     "owner",
		Namespace: "main",
		Key:       "secret2",
	}
	req1 := &vault.GetSecretsRequest{
		Requests: []*vault.SecretRequest{
			{
				Id: id1,
			},
		},
	}

	ch := make(chan *vaulttypes.Response, 1)
	err := store.Add(&vaulttypes.Request{
		Payload:      req1,
		ResponseChan: ch,
		IDVal:        vaulttypes.KeyFor(id1),
	})
	require.NoError(t, err)

	value := "encrypted-value"
	resp1 := &vault.GetSecretsResponse{
		Responses: []*vault.SecretResponse{
			{
				Id:     id1,
				Result: &vault.SecretResponse_Data{Data: &vault.SecretData{EncryptedValue: value}},
			},
		},
	}
	expectedOutcome1 := &vault.Outcome{
		Id:          vaulttypes.KeyFor(id1),
		RequestType: vault.RequestType_GET_SECRETS,
		Request: &vault.Outcome_GetSecretsRequest{
			GetSecretsRequest: req1,
		},
		Response: &vault.Outcome_GetSecretsResponse{
			GetSecretsResponse: resp1,
		},
	}

	eopb, err := proto.Marshal(&vault.Outcomes{Outcomes: []*vault.Outcome{expectedOutcome1}})
	require.NoError(t, err)

	r := &ReportingPlugin{
		lggr: lggr,
		onchainCfg: ocr3types.ReportingPluginConfig{
			N: 4,
			F: 1,
		},
		store: store,
		cfg: &ReportingPluginConfig{
			BatchSize:                         10,
			PublicKey:                         nil,
			PrivateKeyShare:                   nil,
			MaxSecretsPerOwner:                1,
			MaxCiphertextLengthBytes:          1024,
			MaxIdentifierOwnerLengthBytes:     100,
			MaxIdentifierNamespaceLengthBytes: 100,
			MaxIdentifierKeyLengthBytes:       100,
		},
	}

	seqNr := uint64(1)

	rs, err := r.Reports(t.Context(), seqNr, eopb)
	require.NoError(t, err)

	assert.Len(t, rs, 1)
	report := rs[0]

	err = transmitter.Transmit(
		t.Context(),
		types.ConfigDigest([32]byte{0: 1}),
		1,
		report.ReportWithInfo,
		[]types.AttributedOnchainSignature{
			types.AttributedOnchainSignature{Signature: []byte{0: 2}},
			types.AttributedOnchainSignature{Signature: []byte{0: 3}},
		},
	)
	require.NoError(t, err)

	resp := <-ch
	assert.Equal(t, report.ReportWithInfo.Report, types.Report(resp.Payload))
	assert.Equal(t, "REPORT_FORMAT_PROTOBUF", resp.Format)
	assert.Equal(t, vaulttypes.KeyFor(id1), resp.ID)
}
