package vault

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
)

func TestPendingQueueItemBytes_ListObservationAtLeastOutcome(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1), withBatchSize(10))

	req := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId:     "rid",
		Owner:         "owner",
		Namespace:     "ns",
		OrgId:         "org",
		WorkflowOwner: "wf",
	}
	anyReq, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "rid", Item: anyReq}

	obsB, err := pendingQueueItemObservationBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	outB, err := pendingQueueItemOutcomeBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)

	require.Positive(t, obsB)
	require.Positive(t, outB)
	require.GreaterOrEqual(t, obsB, outB, "list observation includes request; outcome is response only")
}

func TestPendingQueueItemBytes_GetSecretsOutcomeLargerThanObservation(t *testing.T) {
	ctx := t.Context()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1), withBatchSize(10), withKeys(pk, shares[0]))

	req := &vaultcommon.GetSecretsRequest{
		OrgId: "org", WorkflowOwner: "wf",
		Requests: []*vaultcommon.SecretRequest{
			{
				Id:             &vaultcommon.SecretIdentifier{Key: "k", Owner: "owner", Namespace: "ns"},
				EncryptionKeys: []string{"abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"},
			},
		},
	}
	anyReq, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "req-1", Item: anyReq}

	obsB, err := pendingQueueItemObservationBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	outB, err := pendingQueueItemOutcomeBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)

	require.Greater(t, outB, obsB, "outcome aggregates 2F+1 shares per key")
}

func TestPendingQueueItemBytes_GetSecretsMonotonicInEncryptionKeys(t *testing.T) {
	ctx := t.Context()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1), withKeys(pk, shares[0]))
	id := &vaultcommon.SecretIdentifier{Key: "k", Owner: "owner", Namespace: "ns"}
	mk := func(keys ...string) *vaultcommon.StoredPendingQueueItem {
		req := &vaultcommon.GetSecretsRequest{
			OrgId: "o", WorkflowOwner: "w",
			Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: keys}},
		}
		a, ierr := anypb.New(req)
		require.NoError(t, ierr)
		return &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	}
	// Empty EncryptionKeys is budgeted as one recipient (keys==0 -> keys=1 in sizing), so the
	// 0-key and 1-key cases can tie; monotonicity applies once there is at least one explicit key.
	b1, err := pendingQueueItemObservationBytes(ctx, mk("aa"), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	b2, err := pendingQueueItemObservationBytes(ctx, mk("aa", "bb"), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	b3, err := pendingQueueItemObservationBytes(ctx, mk("aa", "bb", "cc"), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	require.Less(t, b1, b2)
	require.Less(t, b2, b3)
}

func TestPendingQueueItemBytes_GetSecretsMonotonicInF(t *testing.T) {
	ctx := t.Context()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	cfg := makeReportingPluginConfig(t, 10, pk, shares[0], 100, 2000, 64, 64, 64, 10, 0)
	applyTestByteBudgets(t, cfg, ocr3types.ReportingPluginConfig{N: 7, F: 2})

	id := &vaultcommon.SecretIdentifier{Key: "k", Owner: "owner", Namespace: "ns"}
	req := &vaultcommon.GetSecretsRequest{
		OrgId: "o", WorkflowOwner: "w",
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{"k1"}}},
	}
	a, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}

	outLowF, err := pendingQueueItemOutcomeBytes(ctx, item, cfg, 1)
	require.NoError(t, err)
	outHighF, err := pendingQueueItemOutcomeBytes(ctx, item, cfg, 2)
	require.NoError(t, err)
	require.Less(t, outLowF, outHighF)
}

func TestPendingQueueItemBytes_CreateSecretsMonotonicBatch(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	mk := func(n int) *vaultcommon.StoredPendingQueueItem {
		secrets := make([]*vaultcommon.EncryptedSecret, n)
		for i := range n {
			secrets[i] = &vaultcommon.EncryptedSecret{
				Id:             &vaultcommon.SecretIdentifier{Key: fmt.Sprintf("k%d", i), Owner: "o", Namespace: "n"},
				EncryptedValue: "00",
			}
		}
		req := &vaultcommon.CreateSecretsRequest{
			RequestId: "r", OrgId: "org", WorkflowOwner: "wf", EncryptedSecrets: secrets,
		}
		a, err := anypb.New(req)
		require.NoError(t, err)
		return &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	}
	b1, err := pendingQueueItemObservationBytes(ctx, mk(1), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	b3, err := pendingQueueItemObservationBytes(ctx, mk(3), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	require.Less(t, b1, b3)
}

func TestPendingQueueItemBytes_DeleteSecretsMonotonicIds(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	mk := func(n int) *vaultcommon.StoredPendingQueueItem {
		ids := make([]*vaultcommon.SecretIdentifier, n)
		for i := range n {
			ids[i] = &vaultcommon.SecretIdentifier{Key: fmt.Sprintf("k%d", i), Owner: "o", Namespace: "n"}
		}
		req := &vaultcommon.DeleteSecretsRequest{RequestId: "r", OrgId: "org", WorkflowOwner: "wf", Ids: ids}
		a, err := anypb.New(req)
		require.NoError(t, err)
		return &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	}
	b1, err := pendingQueueItemObservationBytes(ctx, mk(1), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	b4, err := pendingQueueItemObservationBytes(ctx, mk(4), r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	require.Less(t, b1, b4)
}

func TestPendingQueueItemBytes_ListPerOwnerScopedLimiterUsesRequestOwner(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	msl, err := limits.MakeUpperBoundLimiter(limits.Factory{}, cresettings.Default.PerOwner.VaultSecretsLimit)
	require.NoError(t, err)
	r.cfg.MaxSecretsPerOwner = msl
	req := &vaultcommon.ListSecretIdentifiersRequest{
		RequestId: "r", Owner: "someowner", Namespace: "n", OrgId: "org", WorkflowOwner: "wf",
	}
	a, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	_, err = pendingQueueItemObservationBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
	_, err = pendingQueueItemOutcomeBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.NoError(t, err)
}

func TestPendingQueueItemBytes_ListScalesWithMaxSecretsPerOwner(t *testing.T) {
	ctx := t.Context()
	rSmall := newTestReportingPlugin(t, withOnchainCfg(4, 1), withMaxSecretsPerOwner(3))
	rLarge := newTestReportingPlugin(t, withOnchainCfg(4, 1), withMaxSecretsPerOwner(20))
	req := &vaultcommon.ListSecretIdentifiersRequest{RequestId: "r", Owner: "o", Namespace: "n", OrgId: "org", WorkflowOwner: "wf"}
	a, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	small, err := pendingQueueItemObservationBytes(ctx, item, rSmall.cfg, rSmall.onchainCfg.F)
	require.NoError(t, err)
	large, err := pendingQueueItemObservationBytes(ctx, item, rLarge.cfg, rLarge.onchainCfg.F)
	require.NoError(t, err)
	require.Less(t, small, large)
}

func TestPendingQueueItemBytes_UnknownPayloadErrors(t *testing.T) {
	ctx := t.Context()
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1))
	a, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	_, err = pendingQueueItemObservationBytes(ctx, item, r.cfg, r.onchainCfg.F)
	require.Error(t, err)
}

func TestPendingQueueItemBytes_GetSecretsObservationIgnoresF(t *testing.T) {
	ctx := t.Context()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t, withOnchainCfg(4, 1), withKeys(pk, shares[0]))
	id := &vaultcommon.SecretIdentifier{Key: "k", Owner: "owner", Namespace: "ns"}
	req := &vaultcommon.GetSecretsRequest{
		OrgId: "o", WorkflowOwner: "w",
		Requests: []*vaultcommon.SecretRequest{{Id: id, EncryptionKeys: []string{"k1"}}},
	}
	a, err := anypb.New(req)
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: "id", Item: a}
	b1, err := pendingQueueItemObservationBytes(ctx, item, r.cfg, 1)
	require.NoError(t, err)
	b9, err := pendingQueueItemObservationBytes(ctx, item, r.cfg, 9)
	require.NoError(t, err)
	require.Equal(t, b1, b9)
}
