package trigger_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
)

func TestBaseTrigger_CRE_MissingOrgID_DoesNotPersistOrResend(t *testing.T) {
	lggr := logger.Test(t)
	ctx := t.Context()

	getter, err := settings.NewJSONGetter([]byte(`{
		"global": {
			"PerOrg": {"BaseTriggerRetransmitEnabled": "true"},
			"BaseTriggerRetryInterval": "20ms",
			"BaseTriggerMaxRetries": "100"
		}
	}`))
	require.NoError(t, err)

	store := capabilities.NewMemEventStore()
	b, err := capabilities.NewBaseTriggerCapabilityWithCRESettings(ctx, store,
		func() *wrapperspb.BytesValue { return &wrapperspb.BytesValue{} },
		lggr, "testCap", getter)
	require.NoError(t, err)

	ch := make(chan capabilities.TriggerAndId[*wrapperspb.BytesValue], 10)
	b.RegisterTrigger("trig", ch)
	require.NoError(t, b.Start(ctx))
	t.Cleanup(func() { b.Stop() })

	msg := &wrapperspb.BytesValue{Value: []byte("payload")}
	anyMsg, err := anypb.New(msg)
	require.NoError(t, err)
	te := capabilities.TriggerEvent{ID: "e1", Payload: anyMsg}

	// No contexts.WithCRE: org ID is empty while CRE settings are enabled — retransmit must stay off.
	require.NoError(t, b.DeliverEvent(ctx, te, "trig"))

	recs, err := store.List(ctx)
	require.NoError(t, err)
	require.Empty(t, recs, "missing org must not insert a pending row when using CRE settings")

	select {
	case <-ch:
	default:
		t.Fatal("expected immediate inbox delivery even when retransmit path is disabled")
	}

	select {
	case extra := <-ch:
		t.Fatalf("did not expect retransmit/resend without a persisted pending event, got %+v", extra)
	case <-time.After(200 * time.Millisecond):
	}
}
