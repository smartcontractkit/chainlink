package vault

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/vault/vaulttypes"
)

func writeUnobservablePendingQueueItem(t *testing.T, rdr *kv, id string) *vaultcommon.StoredPendingQueueItem {
	t.Helper()
	badAny, err := anypb.New(wrapperspb.String("not-a-vault-request"))
	require.NoError(t, err)
	item := &vaultcommon.StoredPendingQueueItem{Id: id, Item: badAny}
	require.NoError(t, newTestWriteStore(t, rdr).WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{item}))
	return item
}

// Regression: appendPendingQueueObservations walks the full pending queue and emits
// error contributions for unobservable items, and validatePendingQueueObservationsPrefix
// uses the same full walk — so honest observations are accepted rather than rejected
// and stalling the round.
func TestIncludeInvalid_UnobservableQueueItemObservationRoundTrip(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	writeUnobservablePendingQueueItem(t, kv, "bad-head")

	obs := observePendingQueueOnly(t, r, kv)
	require.Len(t, obs.Observations, 1)
	require.True(t, observationContributionIsErr(obs.Observations[0]))
	require.Equal(t, "bad-head", obs.Observations[0].Id)

	require.NoError(t, validatePendingQueueObservation(t, r, kv, obs))

	makeAO := func(observer commontypes.OracleID) types.AttributedObservation {
		obsb, merr := proto.Marshal(obs)
		require.NoError(t, merr)
		return types.AttributedObservation{Observer: observer, Observation: types.Observation(obsb)}
	}

	aos := []types.AttributedObservation{makeAO(0), makeAO(1), makeAO(2), makeAO(3)}
	reached, qerr := r.ObservationQuorum(t.Context(), 1, types.AttributedQuery{}, aos, kv, &blobber{})
	require.NoError(t, qerr)
	require.True(t, reached, "F+1 error contributions on head should satisfy observation quorum")

	out, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, kv, &blobber{})
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal([]byte(out), os))
	require.Len(t, os.Outcomes, 1)
	require.Equal(t, "bad-head", os.Outcomes[0].Id)
	require.Equal(t, vaultcommon.RequestType_UNKNOWN, os.Outcomes[0].RequestType)
}

func TestIncludeInvalid_UnobservableHeadWithValidTailAlignsByQueuePosition(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	badItem := writeUnobservablePendingQueueItem(t, kv, "aaa-bad-head")
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "my_secret"}
	delReq := &vaultcommon.DeleteSecretsRequest{RequestId: "zzz-valid-tail", Ids: []*vaultcommon.SecretIdentifier{id}}
	anyDel, err := anypb.New(delReq)
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, kv).WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{
		badItem,
		{Id: "zzz-valid-tail", Item: anyDel},
	}))

	obs := observePendingQueueOnly(t, r, kv)
	require.Len(t, obs.Observations, 2)
	require.True(t, observationContributionIsErr(obs.Observations[0]))
	require.Equal(t, "aaa-bad-head", obs.Observations[0].Id)
	require.True(t, observationContributionIsOk(obs.Observations[1]))
	require.Equal(t, "zzz-valid-tail", obs.Observations[1].Id)

	require.NoError(t, validatePendingQueueObservation(t, r, kv, obs))
}

func TestValidateObservation_AcceptsErrContributionForQueueItem(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret1"}
	delReq := &vaultcommon.DeleteSecretsRequest{RequestId: "request-1", Ids: []*vaultcommon.SecretIdentifier{id}}
	anyDel, err := anypb.New(delReq)
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, kv).WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{
		{Id: vaulttypes.KeyFor(id), Item: anyDel},
	}))

	obs := &vaultcommon.Observations{
		SortNonce: make([]byte, sortNonceLength),
		Observations: []*vaultcommon.Observation{
			observationToErrContribution(&vaultcommon.Observation{
				Id:          vaulttypes.KeyFor(id),
				RequestType: vaultcommon.RequestType_DELETE_SECRETS,
			}, "request is not valid"),
		},
	}
	obsb, err := proto.Marshal(obs)
	require.NoError(t, err)

	require.NoError(t, r.ValidateObservation(
		t.Context(),
		1,
		types.AttributedQuery{},
		types.AttributedObservation{Observer: 0, Observation: types.Observation(obsb)},
		kv,
		nil,
	))
}

func TestValidateObservation_IncludeInvalid_AcceptsNonMaximalPrefix(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
		withMaxObservationBytes(10*1024*1024),
	)

	rdr := &kv{m: make(map[string]response)}
	writeDeleteSecretsPendingQueueItems(t, rdr, "request-1", "request-2")
	fullObs := observePendingQueueOnly(t, r, rdr)
	require.Len(t, fullObs.Observations, 2)

	prefixObs := &vaultcommon.Observations{
		Observations:      fullObs.Observations[:1],
		PendingQueueItems: fullObs.PendingQueueItems,
		SortNonce:         fullObs.SortNonce,
	}
	require.NoError(t, validatePendingQueueObservation(t, r, rdr, prefixObs))
}

func TestObservationQuorum_IncludeInvalid_RequiresHeadContribution(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret1"}
	delReq := &vaultcommon.DeleteSecretsRequest{RequestId: "request-1", Ids: []*vaultcommon.SecretIdentifier{id}}
	anyDel, err := anypb.New(delReq)
	require.NoError(t, err)
	require.NoError(t, newTestWriteStore(t, kv).WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{
		{Id: vaulttypes.KeyFor(id), Item: anyDel},
	}))

	makeAO := func(observer commontypes.OracleID, errContribution bool) types.AttributedObservation {
		var item *vaultcommon.Observation
		if errContribution {
			item = observationToErrContribution(&vaultcommon.Observation{
				Id:          vaulttypes.KeyFor(id),
				RequestType: vaultcommon.RequestType_DELETE_SECRETS,
			}, "bad request")
		} else {
			item = &vaultcommon.Observation{
				Id:          vaulttypes.KeyFor(id),
				RequestType: vaultcommon.RequestType_DELETE_SECRETS,
				Response: &vaultcommon.Observation_DeleteSecretsResponse{
					DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
						Responses: []*vaultcommon.DeleteSecretResponse{{Id: id, Success: false}},
					},
				},
			}
		}
		obs := &vaultcommon.Observations{
			SortNonce:    make([]byte, sortNonceLength),
			Observations: []*vaultcommon.Observation{item},
		}
		obsb, merr := proto.Marshal(obs)
		require.NoError(t, merr)
		return types.AttributedObservation{Observer: observer, Observation: types.Observation(obsb)}
	}

	// N=4, F=1 -> need 3 observations for QuorumTwoFPlusOne; head needs f+1=2 Err or 2f+1=3 Ok.
	aos := []types.AttributedObservation{
		makeAO(0, true),
		makeAO(1, true),
		makeAO(2, false),
	}
	reached, qerr := r.ObservationQuorum(t.Context(), 1, types.AttributedQuery{}, aos, kv, nil)
	require.NoError(t, qerr)
	require.True(t, reached)

	aosInsufficient := []types.AttributedObservation{
		makeAO(0, true),
		makeAO(1, false),
		makeAO(2, false),
	}
	reached, qerr = r.ObservationQuorum(t.Context(), 1, types.AttributedQuery{}, aosInsufficient, kv, nil)
	require.NoError(t, qerr)
	require.False(t, reached)
}

func TestStateTransition_IncludeInvalid_ProcessesHeadBeforeStoppingOnTail(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	// Head sorts after tail alphabetically; include-invalid must still process head first.
	writeDeleteSecretsPendingQueueItems(t, kv, "zzz-head", "aaa-tail")

	secretID := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "my_secret"}
	makeDeleteObs := func(requestID string) *vaultcommon.Observation {
		return &vaultcommon.Observation{
			Id:          requestID,
			RequestType: vaultcommon.RequestType_DELETE_SECRETS,
			Request: &vaultcommon.Observation_DeleteSecretsRequest{
				DeleteSecretsRequest: &vaultcommon.DeleteSecretsRequest{
					RequestId: requestID,
					Ids:       []*vaultcommon.SecretIdentifier{secretID},
				},
			},
			Response: &vaultcommon.Observation_DeleteSecretsResponse{
				DeleteSecretsResponse: &vaultcommon.DeleteSecretsResponse{
					Responses: []*vaultcommon.DeleteSecretResponse{{Id: secretID, Success: false}},
				},
			},
		}
	}
	makeAO := func(observer commontypes.OracleID, includeTail bool) types.AttributedObservation {
		items := []*vaultcommon.Observation{makeDeleteObs("zzz-head")}
		if includeTail {
			items = append(items, makeDeleteObs("aaa-tail"))
		}
		obs := &vaultcommon.Observations{
			SortNonce:    make([]byte, sortNonceLength),
			Observations: items,
		}
		obsb, merr := proto.Marshal(obs)
		require.NoError(t, merr)
		return types.AttributedObservation{Observer: observer, Observation: types.Observation(obsb)}
	}

	// N=4, F=1: head has 2F+1 ok from all nodes; tail has only 2 ok (< 2F+1) and no F+1 err.
	aos := []types.AttributedObservation{
		makeAO(0, true),
		makeAO(1, true),
		makeAO(2, false),
		makeAO(3, false),
	}
	out, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, kv, &blobber{})
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal([]byte(out), os))
	require.Len(t, os.Outcomes, 1)
	require.Equal(t, "zzz-head", os.Outcomes[0].Id)
}

func TestStateTransition_IncludeInvalid_RejectsItemOnFPlusOneErr(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	id := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret1"}
	delReq := &vaultcommon.DeleteSecretsRequest{RequestId: "request-1", Ids: []*vaultcommon.SecretIdentifier{id}}
	anyDel, err := anypb.New(delReq)
	require.NoError(t, err)
	queueID := vaulttypes.KeyFor(id)
	require.NoError(t, newTestWriteStore(t, kv).WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{
		{Id: queueID, Item: anyDel},
	}))

	makeAO := func(observer commontypes.OracleID) types.AttributedObservation {
		item := observationToErrContribution(&vaultcommon.Observation{
			Id:          queueID,
			RequestType: vaultcommon.RequestType_DELETE_SECRETS,
		}, "invalid request")
		obs := &vaultcommon.Observations{
			SortNonce:    make([]byte, sortNonceLength),
			Observations: []*vaultcommon.Observation{item},
		}
		obsb, merr := proto.Marshal(obs)
		require.NoError(t, merr)
		return types.AttributedObservation{Observer: observer, Observation: types.Observation(obsb)}
	}

	aos := []types.AttributedObservation{makeAO(0), makeAO(1), makeAO(2), makeAO(3)}
	out, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, kv, &blobber{})
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal([]byte(out), os))
	require.Len(t, os.Outcomes, 1)
	require.Equal(t, queueID, os.Outcomes[0].Id)
	require.Contains(t, os.Outcomes[0].GetDeleteSecretsResponse().Responses[0].GetError(), "invalid request")
}

func TestBuildRejectedOutcome_FansOutPerItemErrors(t *testing.T) {
	t.Parallel()

	errMsg := "request is not valid"
	owner := "owner"

	id1 := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "secret1"}
	id2 := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "secret2"}

	t.Run("GetSecrets", func(t *testing.T) {
		t.Parallel()
		payload := &vaultcommon.GetSecretsRequest{
			Requests: []*vaultcommon.SecretRequest{
				{Id: id1},
				{Id: id2},
			},
		}
		outcome := buildRejectedOutcome("req-1", payload, vaultcommon.RequestType_GET_SECRETS, errMsg)
		resps := outcome.GetGetSecretsResponse().Responses
		require.Len(t, resps, 2)
		require.Equal(t, id1, resps[0].Id)
		require.Equal(t, id2, resps[1].Id)
		require.Equal(t, errMsg, resps[0].GetError())
		require.Equal(t, errMsg, resps[1].GetError())
	})

	t.Run("DeleteSecrets", func(t *testing.T) {
		t.Parallel()
		payload := &vaultcommon.DeleteSecretsRequest{
			RequestId: "req-1",
			Ids:       []*vaultcommon.SecretIdentifier{id1, id2},
		}
		outcome := buildRejectedOutcome("req-1", payload, vaultcommon.RequestType_DELETE_SECRETS, errMsg)
		resps := outcome.GetDeleteSecretsResponse().Responses
		require.Len(t, resps, 2)
		require.Equal(t, id1, resps[0].Id)
		require.Equal(t, id2, resps[1].Id)
		require.Equal(t, errMsg, resps[0].GetError())
		require.Equal(t, errMsg, resps[1].GetError())
	})

	t.Run("CreateSecrets", func(t *testing.T) {
		t.Parallel()
		payload := &vaultcommon.CreateSecretsRequest{
			RequestId: "req-1",
			EncryptedSecrets: []*vaultcommon.EncryptedSecret{
				{Id: id1, EncryptedValue: "a"},
				{Id: id2, EncryptedValue: "b"},
			},
		}
		outcome := buildRejectedOutcome("req-1", payload, vaultcommon.RequestType_CREATE_SECRETS, errMsg)
		resps := outcome.GetCreateSecretsResponse().Responses
		require.Len(t, resps, 2)
		require.Equal(t, id1, resps[0].Id)
		require.Equal(t, id2, resps[1].Id)
		require.Equal(t, errMsg, resps[0].GetError())
		require.Equal(t, errMsg, resps[1].GetError())
	})
}

func TestStateTransition_IncludeInvalid_RejectsMultiItemBatchOnFPlusOneErr(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	kv := &kv{m: make(map[string]response)}
	id1 := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "main", Key: "secret1"}
	id2 := &vaultcommon.SecretIdentifier{Owner: "owner", Namespace: "ns2", Key: "secret2"}
	delReq := &vaultcommon.DeleteSecretsRequest{
		RequestId: "request-1",
		Ids:       []*vaultcommon.SecretIdentifier{id1, id2},
	}
	anyDel, err := anypb.New(delReq)
	require.NoError(t, err)
	queueID := "request-1"
	require.NoError(t, newTestWriteStore(t, kv).WritePendingQueue(t.Context(), []*vaultcommon.StoredPendingQueueItem{
		{Id: queueID, Item: anyDel},
	}))

	makeAO := func(observer commontypes.OracleID) types.AttributedObservation {
		item := observationToErrContribution(&vaultcommon.Observation{
			Id:          queueID,
			RequestType: vaultcommon.RequestType_DELETE_SECRETS,
		}, "invalid request")
		obs := &vaultcommon.Observations{
			SortNonce:    make([]byte, sortNonceLength),
			Observations: []*vaultcommon.Observation{item},
		}
		obsb, merr := proto.Marshal(obs)
		require.NoError(t, merr)
		return types.AttributedObservation{Observer: observer, Observation: types.Observation(obsb)}
	}

	aos := []types.AttributedObservation{makeAO(0), makeAO(1), makeAO(2), makeAO(3)}
	out, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, kv, &blobber{})
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal([]byte(out), os))
	require.Len(t, os.Outcomes, 1)
	require.Equal(t, queueID, os.Outcomes[0].Id)

	resps := os.Outcomes[0].GetDeleteSecretsResponse().Responses
	require.Len(t, resps, 2)
	require.True(t, proto.Equal(id1, resps[0].Id))
	require.True(t, proto.Equal(id2, resps[1].Id))
	require.Contains(t, resps[0].GetError(), "invalid request")
	require.Contains(t, resps[1].GetError(), "invalid request")
}

// TestValidateContribution_GetSecretsSelfCheckedObsPassesValidateObservation is the
// core desync regression: an honest node stamps a GetSecrets item OK in the
// Observation self-check, and every peer's ValidateObservation must accept that
// same observation. Before unification, checkContribution skipped the response
// share checks that validateGetSecretsContribution enforced, so a valid
// self-checked OK could be rejected wholesale by peers and stall head-of-queue
// quorum. Both phases now call validateContribution, so the verdicts agree.
func TestValidateContribution_GetSecretsSelfCheckedObsPassesValidateObservation(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	rdr := &kv{m: make(map[string]response)}
	writeGetSecretsPendingQueueItems(t, rdr, pk, "get-secrets-1")

	obs := observePendingQueueOnly(t, r, rdr)
	require.Len(t, obs.Observations, 1)
	require.False(t, observationContributionIsErr(obs.Observations[0]), "self-check should stamp a valid GetSecrets item OK")
	require.True(t, observationContributionIsOk(obs.Observations[0]))

	require.NoError(t, validatePendingQueueObservation(t, r, rdr, obs),
		"ValidateObservation must accept the same observation the Observation self-check stamped OK")
}

// TestValidateContribution_GetSecretsOversizedShareRejected proves the shared
// validator now enforces GetSecrets response share-size limits in the
// Observation self-check path too. The old checkContribution skipped share
// checks, so an oversized share would be stamped OK locally and then rejected
// by peers. validateContribution must reject it so the self-check emits an error
// contribution instead.
func TestValidateContribution_GetSecretsOversizedShareRejected(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	owner := "0x0001020304050607080900010203040506070809"
	secretID := &vaultcommon.SecretIdentifier{Owner: owner, Namespace: "main", Key: "my_secret"}
	queueID := "get-secrets-oversized"
	req := &vaultcommon.GetSecretsRequest{
		Requests: []*vaultcommon.SecretRequest{
			{Id: secretID, EncryptionKeys: []string{"enc-key-1"}},
		},
	}
	anyReq, err := anypb.New(req)
	require.NoError(t, err)
	pendingItem := &vaultcommon.StoredPendingQueueItem{Id: queueID, Item: anyReq}

	o := &vaultcommon.Observation{
		Id:          queueID,
		RequestType: vaultcommon.RequestType_GET_SECRETS,
		Response: &vaultcommon.Observation_GetSecretsResponse{
			GetSecretsResponse: &vaultcommon.GetSecretsResponse{
				Responses: []*vaultcommon.SecretResponse{
					{
						Id: secretID,
						Result: &vaultcommon.SecretResponse_Data{
							Data: &vaultcommon.SecretData{
								EncryptedDecryptionKeyShares: []*vaultcommon.EncryptedShares{
									{EncryptionKey: "enc-key-1", BinaryShares: [][]byte{make([]byte, 10*1024*1024)}},
								},
							},
						},
					},
				},
			},
		},
	}

	err = r.validateContribution(t.Context(), pendingItem, o)
	require.Error(t, err, "validateContribution must reject oversized shares so the self-check stamps an error contribution")
	require.ErrorContains(t, err, "exceeds maximum size allowed")
}

// TestStateTransition_IncludeInvalid_GuardAcceptsValidGetSecrets proves the
// StateTransition guard (which re-runs validateContribution on each chosen
// observation) accepts observations that passed the Observation self-check and
// produces an outcome, so the three phases agree end-to-end.
func TestStateTransition_IncludeInvalid_GuardAcceptsValidGetSecrets(t *testing.T) {
	t.Parallel()
	_, pk, shares, err := tdh2easy.GenerateKeys(1, 3)
	require.NoError(t, err)
	r := newTestReportingPlugin(t,
		withKeys(pk, shares[0]),
		withOnchainCfg(4, 1),
	)

	rdr := &kv{m: make(map[string]response)}
	writeGetSecretsPendingQueueItems(t, rdr, pk, "get-secrets-1")

	obs := observePendingQueueOnly(t, r, rdr)
	require.Len(t, obs.Observations, 1)
	obsb, err := proto.Marshal(obs)
	require.NoError(t, err)

	makeAO := func(observer commontypes.OracleID) types.AttributedObservation {
		return types.AttributedObservation{Observer: observer, Observation: types.Observation(obsb)}
	}
	aos := []types.AttributedObservation{makeAO(0), makeAO(1), makeAO(2), makeAO(3)}

	out, err := r.StateTransition(t.Context(), 1, types.AttributedQuery{}, aos, rdr, &blobber{})
	require.NoError(t, err)

	os := &vaultcommon.Outcomes{}
	require.NoError(t, proto.Unmarshal([]byte(out), os))
	require.Len(t, os.Outcomes, 1)
	require.Equal(t, "get-secrets-1", os.Outcomes[0].Id)
	require.Equal(t, vaultcommon.RequestType_GET_SECRETS, os.Outcomes[0].RequestType)
}
