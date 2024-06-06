package resolver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	"github.com/smartcontractkit/chainlink/v2/core/auth"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/cron"
	"github.com/smartcontractkit/chainlink/v2/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/standardcapabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/webhook"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/v2/core/utils/stringutils"
	webauth "github.com/smartcontractkit/chainlink/v2/core/web/auth"
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
	if err := authenticateUserCanEdit(ctx); err != nil {
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
	if err = ValidateBridgeTypeUniqueness(ctx, btr, orm); err != nil {
		return nil, err
	}
	if err := orm.CreateBridgeType(ctx, bt); err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.BridgeCreated, map[string]interface{}{
		"bridgeName":                   bta.Name,
		"bridgeConfirmations":          bta.Confirmations,
		"bridgeMinimumContractPayment": bta.MinimumContractPayment,
		"bridgeURL":                    bta.URL,
	})

	return NewCreateBridgePayload(*bt, bta.IncomingToken), nil
}

func (r *Resolver) CreateCSAKey(ctx context.Context) (*CreateCSAKeyPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().CSA().Create(ctx)
	if err != nil {
		if errors.Is(err, keystore.ErrCSAKeyExists) {
			return NewCreateCSAKeyPayload(nil, err), nil
		}

		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.CSAKeyCreated, map[string]interface{}{
		"CSAPublicKey": key.PublicKey,
		"CSVersion":    key.Version,
	})

	return NewCreateCSAKeyPayload(&key, nil), nil
}

func (r *Resolver) DeleteCSAKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteCSAKeyPayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().CSA().Delete(ctx, string(args.ID))
	if err != nil {
		if errors.As(err, &keystore.KeyNotFoundError{}) {
			return NewDeleteCSAKeyPayload(csakey.KeyV2{}, err), nil
		}

		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.CSAKeyDeleted, map[string]interface{}{"id": args.ID})

	return NewDeleteCSAKeyPayload(key, nil), nil
}

type createFeedsManagerChainConfigInput struct {
	FeedsManagerID       string
	ChainID              string
	ChainType            string
	AccountAddr          string
	AccountAddrPubKey    *string
	AdminAddr            string
	FluxMonitorEnabled   bool
	OCR1Enabled          bool
	OCR1IsBootstrap      *bool
	OCR1Multiaddr        *string
	OCR1P2PPeerID        *string
	OCR1KeyBundleID      *string
	OCR2Enabled          bool
	OCR2IsBootstrap      *bool
	OCR2Multiaddr        *string
	OCR2ForwarderAddress *string
	OCR2P2PPeerID        *string
	OCR2KeyBundleID      *string
	OCR2Plugins          string
}

func (r *Resolver) CreateFeedsManagerChainConfig(ctx context.Context, args struct {
	Input *createFeedsManagerChainConfigInput
}) (*CreateFeedsManagerChainConfigPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	if args.Input.AccountAddrPubKey != nil {
		params.AccountAddressPublicKey = null.StringFromPtr(args.Input.AccountAddrPubKey)
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
		var plugins feeds.Plugins
		if err = plugins.Scan(args.Input.OCR2Plugins); err != nil {
			return nil, err
		}

		params.OCR2Config = feeds.OCR2ConfigModel{
			Enabled:          args.Input.OCR2Enabled,
			IsBootstrap:      *args.Input.OCR2IsBootstrap,
			Multiaddr:        null.StringFromPtr(args.Input.OCR2Multiaddr),
			ForwarderAddress: null.StringFromPtr(args.Input.OCR2ForwarderAddress),
			P2PPeerID:        null.StringFromPtr(args.Input.OCR2P2PPeerID),
			KeyBundleID:      null.StringFromPtr(args.Input.OCR2KeyBundleID),
			Plugins:          plugins,
		}
	}

	id, err := fsvc.CreateChainConfig(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCreateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	ccfg, err := fsvc.GetChainConfig(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCreateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	fmj, _ := json.Marshal(ccfg)
	r.App.GetAuditLogger().Audit(audit.FeedsManChainConfigCreated, map[string]interface{}{"feedsManager": fmj})

	return NewCreateFeedsManagerChainConfigPayload(ccfg, nil, nil), nil
}

func (r *Resolver) DeleteFeedsManagerChainConfig(ctx context.Context, args struct {
	ID string
}) (*DeleteFeedsManagerChainConfigPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(args.ID)
	if err != nil {
		return nil, err
	}

	fsvc := r.App.GetFeedsService()

	ccfg, err := fsvc.GetChainConfig(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteFeedsManagerChainConfigPayload(nil, err), nil
		}

		return nil, err
	}

	if _, err := fsvc.DeleteChainConfig(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteFeedsManagerChainConfigPayload(nil, err), nil
		}

		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.FeedsManChainConfigDeleted, map[string]interface{}{"id": args.ID})

	return NewDeleteFeedsManagerChainConfigPayload(ccfg, nil), nil
}

type updateFeedsManagerChainConfigInput struct {
	AccountAddr          string
	AccountAddrPubKey    *string
	AdminAddr            string
	FluxMonitorEnabled   bool
	OCR1Enabled          bool
	OCR1IsBootstrap      *bool
	OCR1Multiaddr        *string
	OCR1P2PPeerID        *string
	OCR1KeyBundleID      *string
	OCR2Enabled          bool
	OCR2IsBootstrap      *bool
	OCR2Multiaddr        *string
	OCR2ForwarderAddress *string
	OCR2P2PPeerID        *string
	OCR2KeyBundleID      *string
	OCR2Plugins          string
}

func (r *Resolver) UpdateFeedsManagerChainConfig(ctx context.Context, args struct {
	ID    string
	Input *updateFeedsManagerChainConfigInput
}) (*UpdateFeedsManagerChainConfigPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	if args.Input.AccountAddrPubKey != nil {
		params.AccountAddressPublicKey = null.StringFromPtr(args.Input.AccountAddrPubKey)
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
		var plugins feeds.Plugins
		if err = plugins.Scan(args.Input.OCR2Plugins); err != nil {
			return nil, err
		}

		params.OCR2Config = feeds.OCR2ConfigModel{
			Enabled:          args.Input.OCR2Enabled,
			IsBootstrap:      *args.Input.OCR2IsBootstrap,
			Multiaddr:        null.StringFromPtr(args.Input.OCR2Multiaddr),
			ForwarderAddress: null.StringFromPtr(args.Input.OCR2ForwarderAddress),
			P2PPeerID:        null.StringFromPtr(args.Input.OCR2P2PPeerID),
			KeyBundleID:      null.StringFromPtr(args.Input.OCR2KeyBundleID),
			Plugins:          plugins,
		}
	}

	id, err = fsvc.UpdateChainConfig(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	ccfg, err := fsvc.GetChainConfig(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateFeedsManagerChainConfigPayload(nil, err, nil), nil
		}

		return nil, err
	}

	fmj, _ := json.Marshal(ccfg)
	r.App.GetAuditLogger().Audit(audit.FeedsManChainConfigUpdated, map[string]interface{}{"feedsManager": fmj})

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
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	id, err := feedsService.RegisterManager(ctx, params)
	if err != nil {
		if errors.Is(err, feeds.ErrSingleFeedsManager) {
			return NewCreateFeedsManagerPayload(nil, err, nil), nil
		}
		return nil, err
	}

	mgr, err := feedsService.GetManager(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCreateFeedsManagerPayload(nil, err, nil), nil
		}

		return nil, err
	}

	mgrj, _ := json.Marshal(mgr)
	r.App.GetAuditLogger().Audit(audit.FeedsManCreated, map[string]interface{}{"mgrj": mgrj})

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
	if err := authenticateUserCanEdit(ctx); err != nil {
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
	bridge, err := orm.FindBridge(ctx, taskType)
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

	if err := orm.UpdateBridgeType(ctx, &bridge, btr); err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.BridgeUpdated, map[string]interface{}{
		"bridgeName":                   bridge.Name,
		"bridgeConfirmations":          bridge.Confirmations,
		"bridgeMinimumContractPayment": bridge.MinimumContractPayment,
		"bridgeURL":                    bridge.URL,
	})

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
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	mgr, err = feedsService.GetManager(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateFeedsManagerPayload(nil, err, nil), nil
		}

		return nil, err
	}

	mgrj, _ := json.Marshal(mgr)
	r.App.GetAuditLogger().Audit(audit.FeedsManUpdated, map[string]interface{}{"mgrj": mgrj})

	return NewUpdateFeedsManagerPayload(mgr, nil, nil), nil
}

func (r *Resolver) CreateOCRKeyBundle(ctx context.Context) (*CreateOCRKeyBundlePayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().OCR().Create(ctx)
	if err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.OCRKeyBundleCreated, map[string]interface{}{
		"ocrKeyBundleID":                      key.ID(),
		"ocrKeyBundlePublicKeyAddressOnChain": key.PublicKeyAddressOnChain(),
	})

	return NewCreateOCRKeyBundlePayload(&key), nil
}

func (r *Resolver) DeleteOCRKeyBundle(ctx context.Context, args struct {
	ID string
}) (*DeleteOCRKeyBundlePayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
		return nil, err
	}

	deletedKey, err := r.App.GetKeyStore().OCR().Delete(ctx, args.ID)
	if err != nil {
		if errors.As(err, &keystore.KeyNotFoundError{}) {
			return NewDeleteOCRKeyBundlePayloadResolver(ocrkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.OCRKeyBundleDeleted, map[string]interface{}{"id": args.ID})
	return NewDeleteOCRKeyBundlePayloadResolver(deletedKey, nil), nil
}

func (r *Resolver) DeleteBridge(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteBridgePayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	taskType, err := bridges.ParseBridgeName(string(args.ID))
	if err != nil {
		return NewDeleteBridgePayload(nil, err), nil
	}

	orm := r.App.BridgeORM()
	bt, err := orm.FindBridge(ctx, taskType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewDeleteBridgePayload(nil, err), nil
		}

		return nil, err
	}

	jobsUsingBridge, err := r.App.JobORM().FindJobIDsWithBridge(ctx, string(args.ID))
	if err != nil {
		return nil, err
	}
	if len(jobsUsingBridge) > 0 {
		return NewDeleteBridgePayload(nil, fmt.Errorf("bridge has jobs associated with it")), nil
	}

	if err = orm.DeleteBridgeType(ctx, &bt); err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.BridgeDeleted, map[string]interface{}{"name": bt.Name})
	return NewDeleteBridgePayload(&bt, nil), nil
}

func (r *Resolver) CreateP2PKey(ctx context.Context) (*CreateP2PKeyPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().P2P().Create(ctx)
	if err != nil {
		return nil, err
	}

	const keyType = "Ed25519"
	r.App.GetAuditLogger().Audit(audit.KeyCreated, map[string]interface{}{
		"type":         "p2p",
		"id":           key.ID(),
		"p2pPublicKey": key.PublicKeyHex(),
		"p2pPeerID":    key.PeerID(),
		"p2pType":      keyType,
	})

	return NewCreateP2PKeyPayload(key), nil
}

func (r *Resolver) DeleteP2PKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteP2PKeyPayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
		return nil, err
	}

	keyID, err := p2pkey.MakePeerID(string(args.ID))
	if err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().P2P().Delete(ctx, keyID)
	if err != nil {
		if errors.As(err, &keystore.KeyNotFoundError{}) {
			return NewDeleteP2PKeyPayload(p2pkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.KeyDeleted, map[string]interface{}{
		"type": "p2p",
		"id":   args.ID,
	})

	return NewDeleteP2PKeyPayload(key, nil), nil
}

func (r *Resolver) CreateVRFKey(ctx context.Context) (*CreateVRFKeyPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().VRF().Create(ctx)
	if err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.KeyCreated, map[string]interface{}{
		"type":                "vrf",
		"id":                  key.ID(),
		"vrfPublicKey":        key.PublicKey,
		"vrfPublicKeyAddress": key.PublicKey.Address(),
	})

	return NewCreateVRFKeyPayloadResolver(key), nil
}

func (r *Resolver) DeleteVRFKey(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteVRFKeyPayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
		return nil, err
	}

	key, err := r.App.GetKeyStore().VRF().Delete(ctx, string(args.ID))
	if err != nil {
		if errors.Is(errors.Cause(err), keystore.ErrMissingVRFKey) {
			return NewDeleteVRFKeyPayloadResolver(vrfkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.KeyDeleted, map[string]interface{}{
		"type": "vrf",
		"id":   args.ID,
	})

	return NewDeleteVRFKeyPayloadResolver(key, nil), nil
}

// ApproveJobProposalSpec approves the job proposal spec.
func (r *Resolver) ApproveJobProposalSpec(ctx context.Context, args struct {
	ID    graphql.ID
	Force *bool
}) (*ApproveJobProposalSpecPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	spec, err := feedsSvc.GetSpec(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	specj, _ := json.Marshal(spec)
	r.App.GetAuditLogger().Audit(audit.JobProposalSpecApproved, map[string]interface{}{"spec": specj})

	return NewApproveJobProposalSpecPayload(spec, err), nil
}

// CancelJobProposalSpec cancels the job proposal spec.
func (r *Resolver) CancelJobProposalSpec(ctx context.Context, args struct {
	ID graphql.ID
}) (*CancelJobProposalSpecPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	spec, err := feedsSvc.GetSpec(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	specj, _ := json.Marshal(spec)
	r.App.GetAuditLogger().Audit(audit.JobProposalSpecCanceled, map[string]interface{}{"spec": specj})

	return NewCancelJobProposalSpecPayload(spec, err), nil
}

// RejectJobProposalSpec rejects the job proposal spec.
func (r *Resolver) RejectJobProposalSpec(ctx context.Context, args struct {
	ID graphql.ID
}) (*RejectJobProposalSpecPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	spec, err := feedsSvc.GetSpec(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	specj, _ := json.Marshal(spec)
	r.App.GetAuditLogger().Audit(audit.JobProposalSpecRejected, map[string]interface{}{"spec": specj})

	return NewRejectJobProposalSpecPayload(spec, err), nil
}

// UpdateJobProposalSpecDefinition updates the spec definition.
func (r *Resolver) UpdateJobProposalSpecDefinition(ctx context.Context, args struct {
	ID    graphql.ID
	Input *struct{ Definition string }
}) (*UpdateJobProposalSpecDefinitionPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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

	spec, err := feedsSvc.GetSpec(ctx, id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	specj, _ := json.Marshal(spec)
	r.App.GetAuditLogger().Audit(audit.JobProposalSpecUpdated, map[string]interface{}{"spec": specj})

	return NewUpdateJobProposalSpecDefinitionPayload(spec, err), nil
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

	dbUser, err := r.App.AuthenticationProvider().FindUser(ctx, session.User.Email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(args.Input.OldPassword, dbUser.HashedPassword) {
		r.App.GetAuditLogger().Audit(audit.PasswordResetAttemptFailedMismatch, map[string]interface{}{"user": dbUser.Email})

		return NewUpdatePasswordPayload(nil, map[string]string{
			"oldPassword": "old password does not match",
		}), nil
	}

	if err = r.App.AuthenticationProvider().ClearNonCurrentSessions(ctx, session.SessionID); err != nil {
		return nil, clearSessionsError{}
	}

	err = r.App.AuthenticationProvider().SetPassword(ctx, &dbUser, args.Input.NewPassword)
	if err != nil {
		return nil, failedPasswordUpdateError{}
	}

	r.App.GetAuditLogger().Audit(audit.PasswordResetSuccess, map[string]interface{}{"user": dbUser.Email})
	return NewUpdatePasswordPayload(session.User, nil), nil
}

func (r *Resolver) SetSQLLogging(ctx context.Context, args struct {
	Input struct{ Enabled bool }
}) (*SetSQLLoggingPayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
		return nil, err
	}

	r.App.GetConfig().SetLogSQL(args.Input.Enabled)

	if args.Input.Enabled {
		r.App.GetAuditLogger().Audit(audit.ConfigSqlLoggingEnabled, map[string]interface{}{})
	} else {
		r.App.GetAuditLogger().Audit(audit.ConfigSqlLoggingDisabled, map[string]interface{}{})
	}

	return NewSetSQLLoggingPayload(args.Input.Enabled), nil
}

func (r *Resolver) CreateAPIToken(ctx context.Context, args struct {
	Input struct{ Password string }
}) (*CreateAPITokenPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	session, ok := webauth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return nil, errors.New("Failed to obtain current user from context")
	}
	dbUser, err := r.App.AuthenticationProvider().FindUser(ctx, session.User.Email)
	if err != nil {
		return nil, err
	}

	err = r.App.AuthenticationProvider().TestPassword(ctx, dbUser.Email, args.Input.Password)
	if err != nil {
		r.App.GetAuditLogger().Audit(audit.APITokenCreateAttemptPasswordMismatch, map[string]interface{}{"user": dbUser.Email})

		return NewCreateAPITokenPayload(nil, map[string]string{
			"password": "incorrect password",
		}), nil
	}

	newToken, err := r.App.AuthenticationProvider().CreateAndSetAuthToken(ctx, &dbUser)
	if err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.APITokenCreated, map[string]interface{}{"user": dbUser.Email})
	return NewCreateAPITokenPayload(newToken, nil), nil
}

func (r *Resolver) DeleteAPIToken(ctx context.Context, args struct {
	Input struct{ Password string }
}) (*DeleteAPITokenPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	session, ok := webauth.GetGQLAuthenticatedSession(ctx)
	if !ok {
		return nil, errors.New("Failed to obtain current user from context")
	}
	dbUser, err := r.App.AuthenticationProvider().FindUser(ctx, session.User.Email)
	if err != nil {
		return nil, err
	}

	err = r.App.AuthenticationProvider().TestPassword(ctx, dbUser.Email, args.Input.Password)
	if err != nil {
		r.App.GetAuditLogger().Audit(audit.APITokenDeleteAttemptPasswordMismatch, map[string]interface{}{"user": dbUser.Email})

		return NewDeleteAPITokenPayload(nil, map[string]string{
			"password": "incorrect password",
		}), nil
	}

	err = r.App.AuthenticationProvider().DeleteAuthToken(ctx, &dbUser)
	if err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.APITokenDeleted, map[string]interface{}{"user": dbUser.Email})

	return NewDeleteAPITokenPayload(&auth.Token{
		AccessKey: dbUser.TokenKey.String,
	}, nil), nil
}

func (r *Resolver) CreateJob(ctx context.Context, args struct {
	Input struct {
		TOML string
	}
}) (*CreateJobPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
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
		jb, err = ocr.ValidatedOracleSpecToml(config, r.App.GetRelayers().LegacyEVMChains(), args.Input.TOML)
		if !config.OCR().Enabled() {
			return nil, errors.New("The Offchain Reporting feature is disabled by configuration")
		}
	case job.OffchainReporting2:
		jb, err = validate.ValidatedOracleSpecToml(ctx, r.App.GetConfig().OCR2(), r.App.GetConfig().Insecure(), args.Input.TOML, r.App.GetLoopRegistrarConfig())
		if !config.OCR2().Enabled() {
			return nil, errors.New("The Offchain Reporting 2 feature is disabled by configuration")
		}
	case job.DirectRequest:
		jb, err = directrequest.ValidatedDirectRequestSpec(args.Input.TOML)
	case job.FluxMonitor:
		jb, err = fluxmonitorv2.ValidatedFluxMonitorSpec(config.JobPipeline(), args.Input.TOML)
	case job.Keeper:
		jb, err = keeper.ValidatedKeeperSpec(args.Input.TOML)
	case job.Cron:
		jb, err = cron.ValidatedCronSpec(args.Input.TOML)
	case job.VRF:
		jb, err = vrfcommon.ValidatedVRFSpec(args.Input.TOML)
	case job.Webhook:
		jb, err = webhook.ValidatedWebhookSpec(ctx, args.Input.TOML, r.App.GetExternalInitiatorManager())
	case job.BlockhashStore:
		jb, err = blockhashstore.ValidatedSpec(args.Input.TOML)
	case job.BlockHeaderFeeder:
		jb, err = blockheaderfeeder.ValidatedSpec(args.Input.TOML)
	case job.Bootstrap:
		jb, err = ocrbootstrap.ValidatedBootstrapSpecToml(args.Input.TOML)
	case job.Gateway:
		jb, err = gateway.ValidatedGatewaySpec(args.Input.TOML)
	case job.Workflow:
		jb, err = workflows.ValidatedWorkflowSpec(args.Input.TOML)
	case job.StandardCapabilities:
		jb, err = standardcapabilities.ValidatedStandardCapabilitiesSpec(args.Input.TOML)
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

	jbj, _ := json.Marshal(jb)
	r.App.GetAuditLogger().Audit(audit.JobCreated, map[string]interface{}{"job": string(jbj)})

	return NewCreateJobPayload(r.App, &jb, nil), nil
}

func (r *Resolver) DeleteJob(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteJobPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt32(string(args.ID))
	if err != nil {
		return nil, err
	}

	j, err := r.App.JobORM().FindJobWithoutSpecErrors(ctx, id)
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

	r.App.GetAuditLogger().Audit(audit.JobDeleted, map[string]interface{}{"id": args.ID})
	return NewDeleteJobPayload(r.App, &j, nil), nil
}

func (r *Resolver) DismissJobError(ctx context.Context, args struct {
	ID graphql.ID
}) (*DismissJobErrorPayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	id, err := stringutils.ToInt64(string(args.ID))
	if err != nil {
		return nil, err
	}

	specErr, err := r.App.JobORM().FindSpecError(ctx, id)
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

	r.App.GetAuditLogger().Audit(audit.JobErrorDismissed, map[string]interface{}{"id": args.ID})
	return NewDismissJobErrorPayload(&specErr, nil), nil
}

func (r *Resolver) RunJob(ctx context.Context, args struct {
	ID graphql.ID
}) (*RunJobPayloadResolver, error) {
	if err := authenticateUserCanRun(ctx); err != nil {
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

	plnRun, err := r.App.PipelineORM().FindRun(ctx, jobRunID)
	if err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.JobRunSet, map[string]interface{}{"jobID": args.ID, "jobRunID": jobRunID, "planRunID": plnRun})
	return NewRunJobPayload(&plnRun, r.App, nil), nil
}

func (r *Resolver) SetGlobalLogLevel(ctx context.Context, args struct {
	Level LogLevel
}) (*SetGlobalLogLevelPayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
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

	r.App.GetAuditLogger().Audit(audit.GlobalLogLevelSet, map[string]interface{}{"logLevel": args.Level})
	return NewSetGlobalLogLevelPayload(args.Level, nil), nil
}

// CreateOCR2KeyBundle resolves a create OCR2 Key bundle mutation
func (r *Resolver) CreateOCR2KeyBundle(ctx context.Context, args struct {
	ChainType OCR2ChainType
}) (*CreateOCR2KeyBundlePayloadResolver, error) {
	if err := authenticateUserCanEdit(ctx); err != nil {
		return nil, err
	}

	ct := FromOCR2ChainType(args.ChainType)

	key, err := r.App.GetKeyStore().OCR2().Create(ctx, chaintype.ChainType(ct))
	if err != nil {
		// Not covering the	`chaintype.ErrInvalidChainType` since the GQL model would prevent a non-accepted chain-type from being received
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.OCR2KeyBundleCreated, map[string]interface{}{
		"ocrKeyID":                        key.ID(),
		"ocrKeyChainType":                 key.ChainType(),
		"ocrKeyConfigEncryptionPublicKey": key.ConfigEncryptionPublicKey(),
		"ocrKeyOffchainPublicKey":         key.OffchainPublicKey(),
		"ocrKeyMaxSignatureLength":        key.MaxSignatureLength(),
		"ocrKeyPublicKey":                 key.PublicKey(),
	})

	return NewCreateOCR2KeyBundlePayload(&key), nil
}

// DeleteOCR2KeyBundle resolves a create OCR2 Key bundle mutation
func (r *Resolver) DeleteOCR2KeyBundle(ctx context.Context, args struct {
	ID graphql.ID
}) (*DeleteOCR2KeyBundlePayloadResolver, error) {
	if err := authenticateUserIsAdmin(ctx); err != nil {
		return nil, err
	}

	id := string(args.ID)
	key, err := r.App.GetKeyStore().OCR2().Get(id)
	if err != nil {
		return NewDeleteOCR2KeyBundlePayloadResolver(nil, err), nil
	}

	err = r.App.GetKeyStore().OCR2().Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	r.App.GetAuditLogger().Audit(audit.OCR2KeyBundleDeleted, map[string]interface{}{"id": id})
	return NewDeleteOCR2KeyBundlePayloadResolver(&key, nil), nil
}
