package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/services/webhook"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/core/utils/stringutils"
	webauth "github.com/smartcontractkit/chainlink/core/web/auth"
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

// CreateBridge creates a new bridge.
func (r *Resolver) CreateBridge(ctx context.Context, args struct{ Input createBridgeInput }) (*CreateBridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	var webURL models.WebURL
	if len(args.Input.URL) != 0 {
		rURL, err := url.ParseRequestURI(args.Input.URL)
		if err != nil {
			return nil, err
		}
		webURL = models.WebURL(*rURL)
	}
	minContractPayment := &assets.Link{}
	if err := minContractPayment.UnmarshalText([]byte(args.Input.MinimumContractPayment)); err != nil {
		return nil, err
	}

	btr := &bridges.BridgeTypeRequest{
		Name:                   bridges.BridgeName(args.Input.Name),
		URL:                    webURL,
		Confirmations:          uint32(args.Input.Confirmations),
		MinimumContractPayment: minContractPayment,
	}

	bta, bt, err := bridges.NewBridgeType(btr)
	if err != nil {
		return nil, err
	}
	orm := r.App.BridgeORM()
	if err = ValidateBridgeType(btr); err != nil {
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

func (r *Resolver) DeleteCSAKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteCSAKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().CSA().Delete(string(args.ID))
	if err != nil {
		if errors.As(err, &keystore.KeyNotFoundError{}) {
			return NewDeleteCSAKeyPayload(csakey.KeyV2{}, err), nil
		}

		return nil, err
	}

	return NewDeleteCSAKeyPayload(key, nil), nil
}

type createFeedsManagerChainConfigInput struct {
	FeedsManagerID     string
	ChainID            string
	ChainType          string
	AccountAddr        string
	AdminAddr          string
	FluxMonitorEnabled bool
	OCR1Enabled        bool
	OCR1IsBootstrap    *bool
	OCR1Multiaddr      *string
	OCR1P2PPeerID      *string
	OCR1KeyBundleID    *string
	OCR2Enabled        bool
	OCR2IsBootstrap    *bool
	OCR2Multiaddr      *string
	OCR2P2PPeerID      *string
	OCR2KeyBundleID    *string
}

func (r *Resolver) CreateFeedsManagerChainConfig(ctx context.Context, args struct {
	Input *createFeedsManagerChainConfigInput
}) (*CreateFeedsManagerChainConfigPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	fsvc := r.App.GetFeedsService()

	fmID, err := stringutils.ToInt64(args.Input.FeedsManagerID)
	if err != nil {
		return nil, err
	}

	ctype, err := feeds.NewChainType(args.Input.ChainType)
	if err != nil {
		return nil, err
	}

	params := feeds.ChainConfig{
		FeedsManagerID: fmID,
		ChainID:        args.Input.ChainID,
		ChainType:      ctype,
		AccountAddress: args.Input.AccountAddr,
		AdminAddress:   args.Input.AdminAddr,
		FluxMonitorConfig: feeds.FluxMonitorConfig{
			Enabled: args.Input.FluxMonitorEnabled,
		},
	}

	if args.Input.OCR1Enabled {
		params.OCR1Config = feeds.OCR1Config{
			Enabled:     args.Input.OCR1Enabled,
			IsBootstrap: *args.Input.OCR1IsBootstrap,
			Multiaddr:   null.StringFromPtr(args.Input.OCR1Multiaddr),
			P2PPeerID:   null.StringFromPtr(args.Input.OCR1P2PPeerID),
			KeyBundleID: null.StringFromPtr(args.Input.OCR1KeyBundleID),
		}
	}

	if args.Input.OCR2Enabled {
		params.OCR2Config = feeds.OCR2Config{
			Enabled:     args.Input.OCR2Enabled,
			IsBootstrap: *args.Input.OCR2IsBootstrap,
			Multiaddr:   null.StringFromPtr(args.Input.OCR2Multiaddr),
			P2PPeerID:   null.StringFromPtr(args.Input.OCR2P2PPeerID),
			KeyBundleID: null.StringFromPtr(args.Input.OCR2KeyBundleID),
		}
	}

	id, err := fsvc.CreateChainConfig(params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCreateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	ccfg, err := fsvc.GetChainConfig(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCreateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	return NewCreateFeedsManagerChainConfigPayload(ccfg, nil, nil), nil
}

func (r *Resolver) DeleteFeedsManagerChainConfig(ctx context.Context, args struct {
	ID string
}) (*DeleteFeedsManagerChainConfigPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(args.ID)
	if err != nil {
		return nil, err
	}

	fsvc := r.App.GetFeedsService()

	ccfg, err := fsvc.GetChainConfig(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteFeedsManagerChainConfigPayload(nil, err), nil
		}

		return nil, err
	}

	if _, err := fsvc.DeleteChainConfig(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteFeedsManagerChainConfigPayload(nil, err), nil
		}

		return nil, err
	}

	return NewDeleteFeedsManagerChainConfigPayload(ccfg, nil), nil
}

type updateFeedsManagerChainConfigInput struct {
	AccountAddr        string
	AdminAddr          string
	FluxMonitorEnabled bool
	OCR1Enabled        bool
	OCR1IsBootstrap    *bool
	OCR1Multiaddr      *string
	OCR1P2PPeerID      *string
	OCR1KeyBundleID    *string
	OCR2Enabled        bool
	OCR2IsBootstrap    *bool
	OCR2Multiaddr      *string
	OCR2P2PPeerID      *string
	OCR2KeyBundleID    *string
}

func (r *Resolver) UpdateFeedsManagerChainConfig(ctx context.Context, args struct {
	ID    string
	Input *updateFeedsManagerChainConfigInput
}) (*UpdateFeedsManagerChainConfigPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	fsvc := r.App.GetFeedsService()

	id, err := stringutils.ToInt64(args.ID)
	if err != nil {
		return nil, err
	}

	params := feeds.ChainConfig{
		ID:             id,
		AccountAddress: args.Input.AccountAddr,
		AdminAddress:   args.Input.AdminAddr,
		FluxMonitorConfig: feeds.FluxMonitorConfig{
			Enabled: args.Input.FluxMonitorEnabled,
		},
	}

	if args.Input.OCR1Enabled {
		params.OCR1Config = feeds.OCR1Config{
			Enabled:     args.Input.OCR1Enabled,
			IsBootstrap: *args.Input.OCR1IsBootstrap,
			Multiaddr:   null.StringFromPtr(args.Input.OCR1Multiaddr),
			P2PPeerID:   null.StringFromPtr(args.Input.OCR1P2PPeerID),
			KeyBundleID: null.StringFromPtr(args.Input.OCR1KeyBundleID),
		}
	}

	if args.Input.OCR2Enabled {
		params.OCR2Config = feeds.OCR2Config{
			Enabled:     args.Input.OCR2Enabled,
			IsBootstrap: *args.Input.OCR2IsBootstrap,
			Multiaddr:   null.StringFromPtr(args.Input.OCR2Multiaddr),
			P2PPeerID:   null.StringFromPtr(args.Input.OCR2P2PPeerID),
			KeyBundleID: null.StringFromPtr(args.Input.OCR2KeyBundleID),
		}
	}

	id, err = fsvc.UpdateChainConfig(params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	ccfg, err := fsvc.GetChainConfig(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	return NewUpdateFeedsManagerChainConfigPayload(ccfg, nil, nil), nil
}

type createFeedsManagerInput struct {
	Name      string
	URI       string
	PublicKey string
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

	params := feeds.RegisterManagerParams{
		Name:      args.Input.Name,
		URI:       args.Input.URI,
		PublicKey: *publicKey,
	}

	feedsService := r.App.GetFeedsService()

	id, err := feedsService.RegisterManager(params)
	if err != nil {
		if errors.Is(err, feeds.ErrSingleFeedsManager) {
			return NewCreateFeedsManagerPayload(nil, err, nil), nil
		}
		return nil, err
	}

	mgr, err := feedsService.GetManager(id)
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
	ID    graphql.ID
	Input updateBridgeInput
}) (*UpdateBridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	var webURL models.WebURL
	if len(args.Input.URL) != 0 {
		rURL, err := url.ParseRequestURI(args.Input.URL)
		if err != nil {
			return nil, err
		}
		webURL = models.WebURL(*rURL)
	}
	minContractPayment := &assets.Link{}
	if err := minContractPayment.UnmarshalText([]byte(args.Input.MinimumContractPayment)); err != nil {
		return nil, err
	}

	btr := &bridges.BridgeTypeRequest{
		Name:                   bridges.BridgeName(args.Input.Name),
		URL:                    webURL,
		Confirmations:          uint32(args.Input.Confirmations),
		MinimumContractPayment: minContractPayment,
	}

	taskType, err := bridges.ParseBridgeName(string(args.ID))
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
	if err := ValidateBridgeType(btr); err != nil {
		return nil, err
	}

	if err := orm.UpdateBridgeType(&bridge, btr); err != nil {
		return nil, err
	}

	return NewUpdateBridgePayload(&bridge, nil), nil
}

type updateFeedsManagerInput struct {
	Name      string
	URI       string
	PublicKey string
}

func (r *Resolver) UpdateFeedsManager(ctx context.Context, args struct {
	ID    graphql.ID
	Input *updateFeedsManagerInput
}) (*UpdateFeedsManagerPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	publicKey, err := crypto.PublicKeyFromHex(args.Input.PublicKey)
	if err != nil {
		return NewUpdateFeedsManagerPayload(nil, nil, map[string]string{
			"input/publicKey": "invalid hex value",
		}), nil
	}

	mgr := &feeds.FeedsManager{
		ID:        id,
		URI:       args.Input.URI,
		Name:      args.Input.Name,
		PublicKey: *publicKey,
	}

	feedsService := r.App.GetFeedsService()

	if err = feedsService.UpdateManager(ctx, *mgr); err != nil {
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

	return NewCreateOCRKeyBundlePayload(&key), nil
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

	node, err := r.App.EVMORM().CreateNode(types.Node{
		Name:       args.Input.Name,
		EVMChainID: args.Input.EVMChainID,
		WSURL:      args.Input.WSURL,
		HTTPURL:    args.Input.HTTPURL,
		SendOnly:   args.Input.SendOnly,
	})
	if err != nil {
		return nil, err
	}

	return NewCreateNodePayloadResolver(&node), nil
}

func (r *Resolver) DeleteNode(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteNodePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt32(string(args.ID))
	if err != nil {
		return nil, err
	}

	node, err := r.App.GetChains().EVM.GetNode(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteNodePayloadResolver(nil, err), nil
		}

		return nil, err
	}

	err = r.App.EVMORM().DeleteNode(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Sending the SQL error as the expected error to happen
			// though the prior check should take this into consideration
			// so this should never happen anyway
			return NewDeleteNodePayloadResolver(nil, sql.ErrNoRows), nil
		}

		return nil, err
	}

	return NewDeleteNodePayloadResolver(&node, nil), nil
}

func (r *Resolver) DeleteBridge(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteBridgePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	taskType, err := bridges.ParseBridgeName(string(args.ID))
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

	jobsUsingBridge, err := r.App.JobORM().FindJobIDsWithBridge(string(args.ID))
	if err != nil {
		return nil, err
	}
	if len(jobsUsingBridge) > 0 {
		return NewDeleteBridgePayload(nil, fmt.Errorf("bridge has jobs associated with it")), nil
	}

	if err = orm.DeleteBridgeType(&bt); err != nil {
		return nil, err
	}

	return NewDeleteBridgePayload(&bt, nil), nil
}

func (r *Resolver) CreateP2PKey(ctx context.Context) (*CreateP2PKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().P2P().Create()
	if err != nil {
		return nil, err
	}

	return NewCreateP2PKeyPayload(key), nil
}

func (r *Resolver) DeleteP2PKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteP2PKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	keyID, err := p2pkey.MakePeerID(string(args.ID))
	if err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().P2P().Delete(keyID)
	if err != nil {
		if errors.As(err, &keystore.KeyNotFoundError{}) {
			return NewDeleteP2PKeyPayload(p2pkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewDeleteP2PKeyPayload(key, nil), nil
}

func (r *Resolver) CreateVRFKey(ctx context.Context) (*CreateVRFKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().VRF().Create()
	if err != nil {
		return nil, err
	}

	return NewCreateVRFKeyPayloadResolver(key), nil
}

func (r *Resolver) DeleteVRFKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteVRFKeyPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().VRF().Delete(string(args.ID))
	if err != nil {
		if errors.Is(errors.Cause(err), keystore.ErrMissingVRFKey) {
			return NewDeleteVRFKeyPayloadResolver(vrfkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewDeleteVRFKeyPayloadResolver(key, nil), nil
}

// ApproveJobProposalSpec approves the job proposal spec.
func (r *Resolver) ApproveJobProposalSpec(ctx context.Context, args struct {
	ID    graphql.ID
	Force *bool
}) (*ApproveJobProposalSpecPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	forceApprove := false
	if args.Force != nil {
		forceApprove = *args.Force
	}

	feedsSvc := r.App.GetFeedsService()
	if err = feedsSvc.ApproveSpec(ctx, id, forceApprove); err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, feeds.ErrJobAlreadyExists) {
			return NewApproveJobProposalSpecPayload(nil, err), nil
		}
		return nil, err
	}

	spec, err := feedsSvc.GetSpec(id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return NewApproveJobProposalSpecPayload(spec, err), nil
}

// CancelJobProposalSpec cancels the job proposal spec.
func (r *Resolver) CancelJobProposalSpec(ctx context.Context, args struct {
	ID graphql.ID
}) (*CancelJobProposalSpecPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	feedsSvc := r.App.GetFeedsService()
	if err = feedsSvc.CancelSpec(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCancelJobProposalSpecPayload(nil, err), nil
		}

		return nil, err
	}

	spec, err := feedsSvc.GetSpec(id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return NewCancelJobProposalSpecPayload(spec, err), nil
}

// RejectJobProposalSpec rejects the job proposal spec.
func (r *Resolver) RejectJobProposalSpec(ctx context.Context, args struct {
	ID graphql.ID
}) (*RejectJobProposalSpecPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	feedsSvc := r.App.GetFeedsService()
	if err = feedsSvc.RejectSpec(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewRejectJobProposalSpecPayload(nil, err), nil
		}

		return nil, err
	}

	spec, err := feedsSvc.GetSpec(id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return NewRejectJobProposalSpecPayload(spec, err), nil
}

// UpdateJobProposalSpecDefinition updates the spec definition.
func (r *Resolver) UpdateJobProposalSpecDefinition(ctx context.Context, args struct {
	ID    graphql.ID
	Input *struct{ Definition string }
}) (*UpdateJobProposalSpecDefinitionPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	feedsSvc := r.App.GetFeedsService()

	err = feedsSvc.UpdateSpecDefinition(ctx, id, args.Input.Definition)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateJobProposalSpecDefinitionPayload(nil, err), nil
		}

		return nil, err
	}

	spec, err := feedsSvc.GetSpec(id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return NewUpdateJobProposalSpecDefinitionPayload(spec, err), nil
}

func (r *Resolver) SetServicesLogLevels(ctx context.Context, args struct {
	Input struct{ Config LogLevelConfig }
}) (*SetServicesLogLevelsPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	if args.Input.Config.HeadTracker != nil {
		inputErrs, err := r.setServiceLogLevel(ctx, logger.HeadTracker, *args.Input.Config.HeadTracker)
		if inputErrs != nil {
			return NewSetServicesLogLevelsPayload(nil, inputErrs), nil
		}
		if err != nil {
			return nil, err
		}
	}

	if args.Input.Config.FluxMonitor != nil {
		inputErrs, err := r.setServiceLogLevel(ctx, logger.FluxMonitor, *args.Input.Config.FluxMonitor)
		if inputErrs != nil {
			return NewSetServicesLogLevelsPayload(nil, inputErrs), nil
		}
		if err != nil {
			return nil, err
		}
	}

	if args.Input.Config.Keeper != nil {
		inputErrs, err := r.setServiceLogLevel(ctx, logger.Keeper, *args.Input.Config.Keeper)
		if inputErrs != nil {
			return NewSetServicesLogLevelsPayload(nil, inputErrs), nil
		}
		if err != nil {
			return nil, err
		}
	}

	return NewSetServicesLogLevelsPayload(&args.Input.Config, nil), nil
}

func (r *Resolver) setServiceLogLevel(ctx context.Context, svcName string, logLvl LogLevel) (map[string]string, error) {
	var lvl zapcore.Level
	svcLvl := FromLogLevel(logLvl)

	err := lvl.UnmarshalText([]byte(svcLvl))
	if err != nil {
		return map[string]string{
			svcName + "/" + svcLvl: "invalid log level",
		}, nil
	}

	if err = r.App.SetServiceLogLevel(ctx, svcName, lvl); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) UpdateUserPassword(ctx context.Context, args struct {
	Input UpdatePasswordInput
}) (*UpdatePasswordPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	session, ok := webauth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return nil, errors.New("couldn't retrieve user session")
	}

	dbUser, err := r.App.SessionORM().FindUser()
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(args.Input.OldPassword, dbUser.HashedPassword) {
		return NewUpdatePasswordPayload(nil, map[string]string{
			"oldPassword": "old password does not match",
		}), nil
	}

	if err = r.App.SessionORM().ClearNonCurrentSessions(session.SessionID); err != nil {
		return nil, clearSessionsError{}
	}

	err = r.App.SessionORM().SetPassword(&dbUser, args.Input.NewPassword)
	if err != nil {
		return nil, failedPasswordUpdateError{}
	}

	return NewUpdatePasswordPayload(session.User, nil), nil
}

func (r *Resolver) SetSQLLogging(ctx context.Context, args struct {
	Input struct{ Enabled bool }
}) (*SetSQLLoggingPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	r.App.GetConfig().SetLogSQL(args.Input.Enabled)

	return NewSetSQLLoggingPayload(args.Input.Enabled), nil
}

func (r *Resolver) CreateAPIToken(ctx context.Context, args struct {
	Input struct{ Password string }
}) (*CreateAPITokenPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	dbUser, err := r.App.SessionORM().FindUser()
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(args.Input.Password, dbUser.HashedPassword) {
		return NewCreateAPITokenPayload(nil, map[string]string{
			"password": "incorrect password",
		}), nil
	}

	newToken, err := r.App.SessionORM().CreateAndSetAuthToken(&dbUser)
	if err != nil {
		return nil, err
	}

	return NewCreateAPITokenPayload(newToken, nil), nil
}

func (r *Resolver) DeleteAPIToken(ctx context.Context, args struct {
	Input struct{ Password string }
}) (*DeleteAPITokenPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	dbUser, err := r.App.SessionORM().FindUser()
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(args.Input.Password, dbUser.HashedPassword) {
		return NewDeleteAPITokenPayload(nil, map[string]string{
			"password": "incorrect password",
		}), nil
	}

	err = r.App.SessionORM().DeleteAuthToken(&dbUser)
	if err != nil {
		return nil, err
	}

	return NewDeleteAPITokenPayload(&auth.Token{
		AccessKey: dbUser.TokenKey.String,
	}, nil), nil
}

func (r *Resolver) CreateChain(ctx context.Context, args struct {
	Input struct {
		ID                 graphql.ID
		Config             ChainConfigInput
		KeySpecificConfigs []*KeySpecificChainConfigInput
	}
}) (*CreateChainPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	var id utils.Big
	err := id.UnmarshalText([]byte(args.Input.ID))
	if err != nil {
		return nil, err
	}

	chainCfg, inputErrs := ToChainConfig(args.Input.Config)
	if len(inputErrs) > 0 {
		return NewCreateChainPayload(nil, inputErrs), nil
	}

	if args.Input.KeySpecificConfigs != nil {
		sCfgs := make(map[string]types.ChainCfg)

		for _, cfg := range args.Input.KeySpecificConfigs {
			if cfg != nil {
				sCfg, inputErrs := ToChainConfig(cfg.Config)
				if len(inputErrs) > 0 {
					return NewCreateChainPayload(nil, inputErrs), nil
				}

				sCfgs[cfg.Address] = *sCfg
			}
		}

		chainCfg.KeySpecific = sCfgs
	}

	chain, err := r.App.GetChains().EVM.Add(ctx, id, chainCfg)
	if err != nil {
		return nil, err
	}

	return NewCreateChainPayload(&chain, nil), nil
}

func (r *Resolver) UpdateChain(ctx context.Context, args struct {
	ID    graphql.ID
	Input struct {
		Enabled            bool
		Config             ChainConfigInput
		KeySpecificConfigs []*KeySpecificChainConfigInput
	}
}) (*UpdateChainPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	var id utils.Big
	err := id.UnmarshalText([]byte(args.ID))
	if err != nil {
		return nil, err
	}

	chainCfg, inputErrs := ToChainConfig(args.Input.Config)
	if len(inputErrs) > 0 {
		return NewUpdateChainPayload(nil, inputErrs, nil), nil
	}

	if args.Input.KeySpecificConfigs != nil {
		sCfgs := make(map[string]types.ChainCfg)

		for _, cfg := range args.Input.KeySpecificConfigs {
			if cfg != nil {
				sCfg, inputErrs := ToChainConfig(cfg.Config)
				if len(inputErrs) > 0 {
					return NewUpdateChainPayload(nil, inputErrs, nil), nil
				}

				sCfgs[cfg.Address] = *sCfg
			}
		}

		chainCfg.KeySpecific = sCfgs
	}

	chain, err := r.App.GetChains().EVM.Configure(ctx, id, args.Input.Enabled, chainCfg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateChainPayload(nil, nil, err), nil
		}

		return nil, err
	}

	return NewUpdateChainPayload(&chain, nil, nil), nil
}

func (r *Resolver) DeleteChain(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteChainPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	var id utils.Big
	err := id.UnmarshalText([]byte(args.ID))
	if err != nil {
		return nil, err
	}

	chain, err := r.App.EVMORM().Chain(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteChainPayload(nil, err), nil
		}

		return nil, err
	}

	err = r.App.GetChains().EVM.Remove(id)
	if err != nil {
		return nil, err
	}

	return NewDeleteChainPayload(&chain, nil), nil
}

func (r *Resolver) CreateJob(ctx context.Context, args struct {
	Input struct {
		TOML string
	}
}) (*CreateJobPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	jbt, err := job.ValidateSpec(args.Input.TOML)
	if err != nil {
		return NewCreateJobPayload(r.App, nil, map[string]string{
			"TOML spec": errors.Wrap(err, "failed to parse TOML").Error(),
		}), nil
	}

	var jb job.Job
	config := r.App.GetConfig()
	switch jbt {
	case job.OffchainReporting:
		jb, err = ocr.ValidatedOracleSpecToml(r.App.GetChains().EVM, args.Input.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting() {
			return nil, errors.New("The Offchain Reporting feature is disabled by configuration")
		}
	case job.OffchainReporting2:
		jb, err = validate.ValidatedOracleSpecToml(r.App.GetConfig(), args.Input.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting2() {
			return nil, errors.New("The Offchain Reporting 2 feature is disabled by configuration")
		}
	case job.DirectRequest:
		jb, err = directrequest.ValidatedDirectRequestSpec(args.Input.TOML)
	case job.FluxMonitor:
		jb, err = fluxmonitorv2.ValidatedFluxMonitorSpec(config, args.Input.TOML)
	case job.Keeper:
		jb, err = keeper.ValidatedKeeperSpec(args.Input.TOML)
	case job.Cron:
		jb, err = cron.ValidatedCronSpec(args.Input.TOML)
	case job.VRF:
		jb, err = vrf.ValidatedVRFSpec(args.Input.TOML)
	case job.Webhook:
		jb, err = webhook.ValidatedWebhookSpec(args.Input.TOML, r.App.GetExternalInitiatorManager())
	case job.BlockhashStore:
		jb, err = blockhashstore.ValidatedSpec(args.Input.TOML)
	case job.Bootstrap:
		jb, err = ocrbootstrap.ValidatedBootstrapSpecToml(args.Input.TOML)
	default:
		return NewCreateJobPayload(r.App, nil, map[string]string{
			"Job Type": fmt.Sprintf("unknown job type: %s", jbt),
		}), nil
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = r.App.AddJobV2(ctx, &jb)
	if err != nil {
		return nil, err
	}

	return NewCreateJobPayload(r.App, &jb, nil), nil
}

func (r *Resolver) DeleteJob(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteJobPayloadResolver, error) {
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
			return NewDeleteJobPayload(r.App, nil, err), nil
		}

		return nil, err
	}

	err = r.App.DeleteJob(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteJobPayload(r.App, nil, err), nil
		}

		return nil, err
	}

	return NewDeleteJobPayload(r.App, &j, nil), nil
}

func (r *Resolver) DismissJobError(ctx context.Context, args struct {
	ID graphql.ID
}) (*DismissJobErrorPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	specErr, err := r.App.JobORM().FindSpecError(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDismissJobErrorPayload(nil, err), nil
		}

		return nil, err
	}

	err = r.App.JobORM().DismissError(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDismissJobErrorPayload(nil, err), nil
		}

		return nil, err
	}

	return NewDismissJobErrorPayload(&specErr, nil), nil
}

func (r *Resolver) RunJob(ctx context.Context, args struct {
	ID graphql.ID
}) (*RunJobPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	jobID, err := stringutils.ToInt32(string(args.ID))
	if err != nil {
		return nil, err
	}

	jobRunID, err := r.App.RunJobV2(ctx, jobID, nil)
	if err != nil {
		if errors.Is(err, webhook.ErrJobNotExists) {
			return NewRunJobPayload(nil, r.App, err), nil
		}

		return nil, err
	}

	plnRun, err := r.App.PipelineORM().FindRun(jobRunID)
	if err != nil {
		return nil, err
	}

	return NewRunJobPayload(&plnRun, r.App, nil), nil
}

func (r *Resolver) SetGlobalLogLevel(ctx context.Context, args struct {
	Level LogLevel
}) (*SetGlobalLogLevelPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	var lvl zapcore.Level
	logLvl := FromLogLevel(args.Level)

	err := lvl.UnmarshalText([]byte(logLvl))
	if err != nil {
		return NewSetGlobalLogLevelPayload("", map[string]string{
			"level": "invalid log level",
		}), nil
	}

	if err := r.App.SetLogLevel(lvl); err != nil {
		return nil, err
	}

	return NewSetGlobalLogLevelPayload(args.Level, nil), nil
}

// CreateOCR2KeyBundle resolves a create OCR2 Key bundle mutation
func (r *Resolver) CreateOCR2KeyBundle(ctx context.Context, args struct {
	ChainType OCR2ChainType
}) (*CreateOCR2KeyBundlePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	ct := FromOCR2ChainType(args.ChainType)
	key, err := r.App.GetKeyStore().OCR2().Create(chaintype.ChainType(ct))
	if err != nil {
		// Not covering the	`chaintype.ErrInvalidChainType` since the GQL model would prevent a non-accepted chain-type from being received
		return nil, err
	}

	return NewCreateOCR2KeyBundlePayload(&key), nil
}

// DeleteOCR2KeyBundle resolves a create OCR2 Key bundle mutation
func (r *Resolver) DeleteOCR2KeyBundle(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteOCR2KeyBundlePayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id := string(args.ID)
	key, err := r.App.GetKeyStore().OCR2().Get(id)
	if err != nil {
		return NewDeleteOCR2KeyBundlePayloadResolver(nil, err), nil
	}

	err = r.App.GetKeyStore().OCR2().Delete(id)
	if err != nil {
		return nil, err
	}

	return NewDeleteOCR2KeyBundlePayloadResolver(&key, nil), nil
}
