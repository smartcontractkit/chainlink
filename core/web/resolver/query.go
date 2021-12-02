package resolver

import (
	"context"
	"database/sql"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
)

// Bridge retrieves a bridges by name.
func (r *Resolver) Bridge(ctx context.Context, args struct{ ID graphql.ID }) (*BridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	name, err := bridges.NewTaskType(string(args.ID))
	if err != nil {
		return nil, err
	}

	bridge, err := r.App.BridgeORM().FindBridge(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewBridgePayload(bridge, err), nil
		}

		return nil, err
	}

	return NewBridgePayload(bridge, nil), nil
}

// Bridges retrieves a paginated list of bridges.
func (r *Resolver) Bridges(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*BridgesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	brdgs, count, err := r.App.BridgeORM().BridgeTypes(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewBridgesPayload(brdgs, int32(count)), nil
}

// Chain retrieves a chain by id.
func (r *Resolver) Chain(ctx context.Context, args struct{ ID graphql.ID }) (*ChainPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id := utils.Big{}
	err := id.UnmarshalText([]byte(args.ID))
	if err != nil {
		return nil, err
	}

	chain, err := r.App.EVMORM().Chain(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewChainPayload(chain, err), nil
		}

		return nil, err
	}

	return NewChainPayload(chain, nil), nil
}

// Chains retrieves a paginated list of chains.
func (r *Resolver) Chains(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*ChainsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	page, count, err := r.App.EVMORM().Chains(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewChainsPayload(page, int32(count)), nil
}

// FeedsManager retrieves a feeds manager by id.
func (r *Resolver) FeedsManager(ctx context.Context, args struct{ ID graphql.ID }) (*FeedsManagerPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	mgr, err := r.App.GetFeedsService().GetManager(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewFeedsManagerPayload(nil, err), nil
		}

		return nil, err
	}

	return NewFeedsManagerPayload(mgr, nil), nil
}

func (r *Resolver) FeedsManagers(ctx context.Context) (*FeedsManagersPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	mgrs, err := r.App.GetFeedsService().ListManagers()
	if err != nil {
		return nil, err
	}

	return NewFeedsManagersPayload(mgrs), nil
}

// Job retrieves a job by id.
func (r *Resolver) Job(ctx context.Context, args struct{ ID graphql.ID }) (*JobPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt32(string(args.ID))
	if err != nil {
		return nil, err
	}

	j, err := r.App.JobORM().FindJobTx(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewJobPayload(r.App, nil, err), nil

		}

		return nil, err
	}

	return NewJobPayload(r.App, &j, nil), nil
}

// Jobs fetches a paginated list of jobs
func (r *Resolver) Jobs(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*JobsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	jobs, count, err := r.App.JobORM().FindJobs(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewJobsPayload(r.App, jobs, int32(count)), nil
}

func (r *Resolver) OCRKeyBundles(ctx context.Context) (*OCRKeyBundlesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	ocrKeyBundles, err := r.App.GetKeyStore().OCR().GetAll()
	if err != nil {
		return nil, err
	}

	return NewOCRKeyBundlesPayloadResolver(ocrKeyBundles), nil
}

func (r *Resolver) CSAKeys(ctx context.Context) (*CSAKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().CSA().GetAll()
	if err != nil {
		return nil, err
	}

	return NewCSAKeysResolver(keys), nil
}

// Features retrieves each featured enabled by boolean mapping
func (r *Resolver) Features(ctx context.Context) (*FeaturesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	return NewFeaturesPayloadResolver(r.App.GetConfig()), nil
}

// Node retrieves a node by ID
func (r *Resolver) Node(ctx context.Context, args struct{ ID graphql.ID }) (*NodePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt32(string(args.ID))
	if err != nil {
		return nil, err
	}

	node, err := r.App.EVMORM().Node(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewNodePayloadResolver(nil, err), nil
		}
		return nil, err
	}

	return NewNodePayloadResolver(&node, nil), nil
}

func (r *Resolver) P2PKeys(ctx context.Context) (*P2PKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	p2pKeys, err := r.App.GetKeyStore().P2P().GetAll()
	if err != nil {
		return nil, err
	}

	return NewP2PKeysPayload(p2pKeys), nil
}

// VRFKeys fetches all VRF keys.
func (r *Resolver) VRFKeys(ctx context.Context) (*VRFKeysPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keys, err := r.App.GetKeyStore().VRF().GetAll()
	if err != nil {
		return nil, err
	}

	return NewVRFKeysPayloadResolver(keys), nil
}

// VRFKey fetches the VRF key with the given ID.
func (r *Resolver) VRFKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*VRFKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().VRF().Get(string(args.ID))
	if err != nil {
		if errors.Cause(err) == keystore.ErrMissingVRFKey {
			return NewVRFKeyPayloadResolver(vrfkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewVRFKeyPayloadResolver(key, nil), err
}

// JobProposal retrieves a job proposal by ID
func (r *Resolver) JobProposal(ctx context.Context, args struct {
	ID graphql.ID
}) (*JobProposalPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	jp, err := r.App.GetFeedsService().GetJobProposal(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewJobProposalPayload(nil, err), nil
		}

		return nil, err
	}

	return NewJobProposalPayload(jp, err), nil
}

// Nodes retrieves a paginated list of nodes.
func (r *Resolver) Nodes(ctx context.Context, args struct {
	Offset *int32
	Limit  *int32
}) (*NodesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	nodes, count, err := r.App.EVMORM().Nodes(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewNodesPayload(nodes, int32(count)), nil
}
