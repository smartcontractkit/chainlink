package vault

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"google.golang.org/protobuf/proto"

	vaultcommon "github.com/smartcontractkit/chainlink-common/pkg/capabilities/actions/vault"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func (r *ReportingPlugin) observeListSecretIdentifiers(ctx context.Context, reader ReadKVStore, req proto.Message, o *vaultcommon.Observation) {
	tp := req.(*vaultcommon.ListSecretIdentifiersRequest)
	o.RequestType = vaultcommon.RequestType_LIST_SECRET_IDENTIFIERS
	o.Request = &vaultcommon.Observation_ListSecretIdentifiersRequest{
		ListSecretIdentifiersRequest: tp,
	}
	l := r.lggr.With("requestId", tp.RequestId, "requestType", "ListSecretIdentifiers", "owner", tp.Owner)

	resp, err := r.processListSecretIdentifiersRequest(ctx, l, reader, tp)
	if err != nil {
		l.Debugw("failed to process list secret identifiers request", "error", err)
		o.Response = &vaultcommon.Observation_ListSecretIdentifiersResponse{
			ListSecretIdentifiersResponse: &vaultcommon.ListSecretIdentifiersResponse{
				Error:   err.Error(),
				Success: false,
			},
		}
		return
	}

	l.Debugw("observed list secret identifiers request")
	o.Response = &vaultcommon.Observation_ListSecretIdentifiersResponse{
		ListSecretIdentifiersResponse: resp,
	}
}

func (r *ReportingPlugin) processListSecretIdentifiersRequest(ctx context.Context, l logger.Logger, reader ReadKVStore, req *vaultcommon.ListSecretIdentifiersRequest) (*vaultcommon.ListSecretIdentifiersResponse, error) {
	if req.Owner == "" {
		return nil, errors.New("invalid request: owner cannot be empty")
	}

	md, err := reader.GetMetadata(ctx, req.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for owner: %w", err)
	}

	if md == nil {
		// No metadata, so the list is empty.
		// The user hasn't added any items to the vault DON yet.
		l.Debugw("successfully read metadata for owner: no metadata found, returning empty list")
		return &vaultcommon.ListSecretIdentifiersResponse{Identifiers: []*vaultcommon.SecretIdentifier{}, Success: true}, nil
	}

	sort.Slice(md.SecretIdentifiers, func(i, j int) bool {
		if md.SecretIdentifiers[i].Namespace == md.SecretIdentifiers[j].Namespace {
			return md.SecretIdentifiers[i].Key < md.SecretIdentifiers[j].Key
		}
		return md.SecretIdentifiers[i].Namespace < md.SecretIdentifiers[j].Namespace
	})

	if req.Namespace == "" {
		return &vaultcommon.ListSecretIdentifiersResponse{Identifiers: md.SecretIdentifiers, Success: true}, nil
	}

	si := []*vaultcommon.SecretIdentifier{}
	for _, id := range md.SecretIdentifiers {
		if id.Namespace == req.Namespace {
			si = append(si, id)
		}
	}

	return &vaultcommon.ListSecretIdentifiersResponse{
		Identifiers: si,
		Success:     true,
	}, nil
}

func (r *ReportingPlugin) stateTransitionListSecretIdentifiers(ctx context.Context, _ WriteKVStore, chosen []*vaultcommon.Observation, o *vaultcommon.Outcome) {
	// All of the logic for the ListSecretIdentifiers request is in the
	// observation phase. This returns the observations in sorted order,
	// so we can just take the first aggregated response and use it as the outcome.
	first := chosen[0]
	if !r.optimizationsEnabled(ctx) {
		o.Request = &vaultcommon.Outcome_ListSecretIdentifiersRequest{
			ListSecretIdentifiersRequest: first.GetListSecretIdentifiersRequest(),
		}
	}
	o.Response = &vaultcommon.Outcome_ListSecretIdentifiersResponse{
		ListSecretIdentifiersResponse: first.GetListSecretIdentifiersResponse(),
	}
}
