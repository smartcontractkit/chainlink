package resolver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/jackc/pgconn"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Resolver struct {
	App chainlink.Application
}

type createBridgeInput struct {
	Name                   string
	URL                    string
	Confirmations          int32
	MinimumContractPayment string
}

// Bridge retrieves a bridges by name.
func (r *Resolver) CreateBridge(ctx context.Context, args struct{ Input createBridgeInput }) (*CreateBridgePayloadResolver, error) {
	var webURL models.WebURL
	if len(args.Input.URL) != 0 {
		url, err := url.ParseRequestURI(args.Input.URL)
		if err != nil {
			return nil, err
		}
		webURL = models.WebURL(*url)
	}
	minContractPayment := &assets.Link{}
	if err := minContractPayment.UnmarshalText([]byte(args.Input.MinimumContractPayment)); err != nil {
		return nil, err
	}

	btr := &bridges.BridgeTypeRequest{
		Name:                   bridges.TaskType(args.Input.Name),
		URL:                    webURL,
		Confirmations:          uint32(args.Input.Confirmations),
		MinimumContractPayment: minContractPayment,
	}

	bta, bt, err := bridges.NewBridgeType(btr)
	if err != nil {
		return nil, err
	}
	orm := r.App.BridgeORM()
	if err = ValidateBridgeType(btr, orm); err != nil {
		return nil, err
	}
	if err = ValidateBridgeTypeUniqueness(btr, orm); err != nil {
		return nil, err
	}
	if e := orm.CreateBridgeType(bt); e != nil {
		return nil, err
	}
	if err != nil {
		switch e := err.(type) {
		case *pgconn.PgError:
			var apiErr error
			if e.ConstraintName == "external_initiators_name_key" {
				apiErr = fmt.Errorf("bridge Type %v conflict", bt.Name)
			} else {
				apiErr = err
			}
			return nil, apiErr
		default:
			return nil, err
		}
	}

	return NewCreateBridgePayload(*bt, bta.IncomingToken), nil
}

// Bridge retrieves a bridges by name.
func (r *Resolver) Bridge(ctx context.Context, args struct{ Name string }) (*BridgeResolver, error) {
	name, err := bridges.NewTaskType(args.Name)
	if err != nil {
		return nil, err
	}

	bridge, err := r.App.BridgeORM().FindBridge(name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("bridge not found")
	}
	if err != nil {
		return nil, err
	}

	return NewBridge(bridge), nil
}

// Bridges retrieves a paginated list of bridges.
func (r *Resolver) Bridges(ctx context.Context, args struct {
	Offset *int
	Limit  *int
}) ([]*BridgeResolver, error) {
	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	bridges, _, err := r.App.BridgeORM().BridgeTypes(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewBridges(bridges), nil
}

type updateBridgeInput struct {
	Name                   string
	URL                    string
	Confirmations          int32
	MinimumContractPayment string
}

func (r *Resolver) UpdateBridge(ctx context.Context, args struct {
	Name  string
	Input updateBridgeInput
}) (*UpdateBridgePayloadResolver, error) {
	var webURL models.WebURL
	if len(args.Input.URL) != 0 {
		url, err := url.ParseRequestURI(args.Input.URL)
		if err != nil {
			return nil, err
		}
		webURL = models.WebURL(*url)
	}
	minContractPayment := &assets.Link{}
	if err := minContractPayment.UnmarshalText([]byte(args.Input.MinimumContractPayment)); err != nil {
		return nil, err
	}

	btr := &bridges.BridgeTypeRequest{
		Name:                   bridges.TaskType(args.Input.Name),
		URL:                    webURL,
		Confirmations:          uint32(args.Input.Confirmations),
		MinimumContractPayment: minContractPayment,
	}

	taskType, err := bridges.NewTaskType(args.Name)
	if err != nil {
		return nil, err
	}

	// Find the bridge
	orm := r.App.BridgeORM()
	bridge, err := orm.FindBridge(taskType)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("bridge not found")
	}

	if err != nil {
		return nil, err
	}

	// Update the bridge
	if err := ValidateBridgeType(btr, orm); err != nil {
		return nil, err
	}

	if err := orm.UpdateBridgeType(&bridge, btr); err != nil {
		return nil, err
	}

	return NewUpdateBridgePayload(bridge), nil
}

// Chain retrieves a chain by id.
func (r *Resolver) Chain(ctx context.Context, args struct{ ID graphql.ID }) (*ChainResolver, error) {
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
	offset := pageOffset(args.Offset)
	limit := pageLimit(args.Limit)

	page, _, err := r.App.EVMORM().Chains(offset, limit)
	if err != nil {
		return nil, err
	}

	return NewChains(page), nil
}

// FeedsManager retrieves a feeds manager by id.
func (r *Resolver) FeedsManager(ctx context.Context, args struct{ ID graphql.ID }) (*FeedsManagerResolver, error) {
	id, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, err
	}

	mgr, err := r.App.GetFeedsService().GetManager(int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("feeds manager not found")
		}

		return nil, err
	}

	return NewFeedsManager(*mgr), nil
}

func (r *Resolver) FeedsManagers() ([]*FeedsManagerResolver, error) {
	mgrs, err := r.App.GetFeedsService().ListManagers()
	if err != nil {
		return nil, err
	}

	return NewFeedsManagers(mgrs), nil
}
