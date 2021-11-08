package resolver

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
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
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

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
	if err := orm.CreateBridgeType(bt); err != nil {
		return nil, err
	}

	return NewCreateBridgePayload(*bt, bta.IncomingToken), nil
}

type createFeedsManagerInput struct {
	Name                   string
	URI                    string
	PublicKey              string
	JobTypes               []JobType
	IsBootstrapPeer        bool
	BootstrapPeerMultiaddr *string
}

func (r *Resolver) CreateCSAKey(ctx context.Context) (*CreateCSAKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().CSA().Create()
	if err != nil {
		if errors.Is(err, keystore.ErrCSAKeyExists) {
			return NewCreateCSAKeyPayload(nil, err), nil
		}

		return nil, err
	}

	return NewCreateCSAKeyPayload(&key, nil), nil
}

func (r *Resolver) CreateFeedsManager(ctx context.Context, args struct {
	Input *createFeedsManagerInput
}) (*CreateFeedsManagerPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	publicKey, err := crypto.PublicKeyFromHex(args.Input.PublicKey)
	if err != nil {
		return NewCreateFeedsManagerPayload(nil, nil, map[string]string{
			"input/publicKey": "invalid hex value",
		}), nil
	}

	// convert enum job types
	jobTypes := pq.StringArray{}
	for _, jt := range args.Input.JobTypes {
		jobTypes = append(jobTypes, FromJobTypeInput(jt))
	}

	mgr := &feeds.FeedsManager{
		Name:                      args.Input.Name,
		URI:                       args.Input.URI,
		PublicKey:                 *publicKey,
		JobTypes:                  jobTypes,
		IsOCRBootstrapPeer:        args.Input.IsBootstrapPeer,
		OCRBootstrapPeerMultiaddr: null.StringFromPtr(args.Input.BootstrapPeerMultiaddr),
	}

	feedsService := r.App.GetFeedsService()

	id, err := feedsService.RegisterManager(mgr)
	if err != nil {
		if errors.Is(err, feeds.ErrSingleFeedsManager) {
			return NewCreateFeedsManagerPayload(nil, err, nil), nil
		}

		return nil, err
	}

	mgr, err = feedsService.GetManager(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCreateFeedsManagerPayload(nil, err, nil), nil
		}

		return nil, err
	}

	return NewCreateFeedsManagerPayload(mgr, nil, nil), nil
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
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

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
		return NewUpdateBridgePayload(nil, err), nil
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

	return NewUpdateBridgePayload(&bridge, nil), nil
}

type updateFeedsManagerInput struct {
	Name                   string
	URI                    string
	PublicKey              string
	JobTypes               []JobType
	IsBootstrapPeer        bool
	BootstrapPeerMultiaddr *string
}

func (r *Resolver) UpdateFeedsManager(ctx context.Context, args struct {
	ID    graphql.ID
	Input *updateFeedsManagerInput
}) (*UpdateFeedsManagerPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(string(args.ID), 10, 32)
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.PublicKeyFromHex(args.Input.PublicKey)
	if err != nil {
		return NewUpdateFeedsManagerPayload(nil, nil, map[string]string{
			"input/publicKey": "invalid hex value",
		}), nil
	}

	// convert enum job types
	jobTypes := pq.StringArray{}
	for _, jt := range args.Input.JobTypes {
		jobTypes = append(jobTypes, FromJobTypeInput(jt))
	}

	mgr := &feeds.FeedsManager{
		ID:                        id,
		URI:                       args.Input.URI,
		Name:                      args.Input.Name,
		PublicKey:                 *publicKey,
		JobTypes:                  jobTypes,
		IsOCRBootstrapPeer:        args.Input.IsBootstrapPeer,
		OCRBootstrapPeerMultiaddr: null.StringFromPtr(args.Input.BootstrapPeerMultiaddr),
	}

	feedsService := r.App.GetFeedsService()

	err = feedsService.UpdateFeedsManager(ctx, *mgr)
	if err != nil {
		return nil, err
	}

	mgr, err = feedsService.GetManager(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateFeedsManagerPayload(nil, err, nil), nil
		}

		return nil, err
	}

	return NewUpdateFeedsManagerPayload(mgr, nil, nil), nil
}

func (r *Resolver) CreateOCRKeyBundle(ctx context.Context) (*CreateOCRKeyBundlePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().OCR().Create()
	if err != nil {
		return nil, err
	}

	return NewCreateOCRKeyBundlePayloadResolver(key), nil
}

func (r *Resolver) DeleteOCRKeyBundle(ctx context.Context, args struct {
	ID string
}) (*DeleteOCRKeyBundlePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	deletedKey, err := r.App.GetKeyStore().OCR().Delete(args.ID)
	if err != nil {
		if errors.As(err, &keystore.KeyNotFoundError{}) {
			return NewDeleteOCRKeyBundlePayloadResolver(ocrkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewDeleteOCRKeyBundlePayloadResolver(deletedKey, nil), nil
}

func (r *Resolver) CreateNode(ctx context.Context, args struct {
	Input *types.NewNode
}) (*CreateNodePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	node, err := r.App.EVMORM().CreateNode(*args.Input)
	if err != nil {
		return nil, err
	}

	return NewCreateNodePayloadResolver(&node), nil
}

func (r *Resolver) DeleteNode(ctx context.Context, args struct {
	ID int32
}) (*DeleteNodePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	node, err := r.App.EVMORM().Node(args.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteNodePayloadResolver(nil, err), nil
		}

		return nil, err
	}

	err = r.App.EVMORM().DeleteNode(int64(args.ID))
	if err != nil {
		if errors.Is(err, evm.ErrNoRowsAffected) {
			return NewDeleteNodePayloadResolver(nil, err), nil
		}

		return nil, err
	}

	return NewDeleteNodePayloadResolver(&node, nil), nil
}

func (r *Resolver) DeleteBridge(ctx context.Context, args struct {
	Name string
}) (*DeleteBridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	taskType, err := bridges.NewTaskType(args.Name)
	if err != nil {
		return NewDeleteBridgePayload(nil, err), nil
	}

	orm := r.App.BridgeORM()
	bt, err := orm.FindBridge(taskType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteBridgePayload(nil, err), nil
		}

		return nil, err
	}

	jobsUsingBridge, err := r.App.JobORM().FindJobIDsWithBridge(args.Name)
	if err != nil {
		if len(jobsUsingBridge) > 0 {
			return NewDeleteBridgePayload(nil, err), nil
		}

		return nil, err
	}

	if err = orm.DeleteBridgeType(&bt); err != nil {
		return nil, err
	}

	return NewDeleteBridgePayload(&bt, nil), nil
}
