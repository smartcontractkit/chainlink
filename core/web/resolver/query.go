package resolver

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Bridge retrieves a bridges by name.
func (r *Resolver) Bridge(ctx context.Context, args struct{ Name string }) (*BridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	name, err := bridges.NewTaskType(args.Name)
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
	Offset *int
	Limit  *int
}) (*BridgesPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	bridges, count, err := r.App.BridgeORM().BridgeTypes(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewBridgesPayload(bridges, int32(count)), nil
}

// Chain retrieves a chain by id.
func (r *Resolver) Chain(ctx context.Context, args struct{ ID graphql.ID }) (*ChainResolver, error) {
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
		return nil, err
	}

	return NewChain(chain), nil
}

// Chains retrieves a paginated list of chains.
func (r *Resolver) Chains(ctx context.Context, args struct {
	Offset *int
	Limit  *int
}) ([]*ChainResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	page, _, err := r.App.EVMORM().Chains(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewChains(page), nil
}

// FeedsManager retrieves a feeds manager by id.
func (r *Resolver) FeedsManager(ctx context.Context, args struct{ ID graphql.ID }) (*FeedsManagerPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, err
	}

	mgr, err := r.App.GetFeedsService().GetManager(int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewFeedsManagerPayload(nil), nil
		}

		return nil, err
	}

	return NewFeedsManagerPayload(mgr), nil
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

	id, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, err
	}

	j, err := r.App.JobORM().FindJobTx(int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewJobPayload(nil, err), nil

		}

		return nil, err
	}

	return NewJobPayload(&j, nil), nil
}

// Jobs fetches a paginated list of jobs
func (r *Resolver) Jobs(ctx context.Context, args struct {
	Offset *int
	Limit  *int
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

	return NewJobsPayload(jobs, int32(count)), nil
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
