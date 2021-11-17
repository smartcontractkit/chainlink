package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	"github.com/graph-gophers/graphql-go"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/chainlink/core/utils/crypto"
	"github.com/smartcontractkit/chainlink/core/web/auth"
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

	return NewCreateP2PKeyPayloadResolver(key), nil
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
			return NewDeleteP2PKeyPayloadResolver(p2pkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewDeleteP2PKeyPayloadResolver(key, nil), nil
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
		if errors.Cause(err) == keystore.ErrMissingVRFKey {
			return NewDeleteVRFKeyPayloadResolver(vrfkey.KeyV2{}, err), nil
		}
		return nil, err
	}

	return NewDeleteVRFKeyPayloadResolver(key, nil), nil
}

func (r *Resolver) ApproveJobProposal(ctx context.Context, args struct {
	ID graphql.ID
}) (*ApproveJobProposalPayloadResolver, error) {
	jp, err := r.executeJobProposalAction(ctx, jobProposalAction{
		args.ID, approve,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewApproveJobProposalPayload(nil, err), nil
		}

		return nil, err
	}

	return NewApproveJobProposalPayload(jp, nil), nil
}

func (r *Resolver) CancelJobProposal(ctx context.Context, args struct {
	ID graphql.ID
}) (*CancelJobProposalPayloadResolver, error) {
	jp, err := r.executeJobProposalAction(ctx, jobProposalAction{
		args.ID, cancel,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewCancelJobProposalPayload(nil, err), nil
		}

		return nil, err
	}

	return NewCancelJobProposalPayload(jp, nil), nil
}

func (r *Resolver) RejectJobProposal(ctx context.Context, args struct {
	ID graphql.ID
}) (*RejectJobProposalPayloadResolver, error) {
	jp, err := r.executeJobProposalAction(ctx, jobProposalAction{
		args.ID, reject,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewRejectJobProposalPayload(nil, err), nil
		}

		return nil, err
	}

	return NewRejectJobProposalPayload(jp, nil), nil
}

func (r *Resolver) UpdateJobProposalSpec(ctx context.Context, args struct {
	ID    graphql.ID
	Input *struct{ Spec string }
}) (*UpdateJobProposalSpecPayloadResolver, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(string(args.ID), 10, 64)
	if err != nil {
		return nil, err
	}

	feedsSvc := r.App.GetFeedsService()

	err = feedsSvc.UpdateJobProposalSpec(ctx, id, args.Input.Spec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateJobProposalSpecPayload(nil, err), nil
		}

		return nil, err
	}

	jp, err := r.App.GetFeedsService().GetJobProposal(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewUpdateJobProposalSpecPayload(nil, err), nil
		}

		return nil, err
	}

	return NewUpdateJobProposalSpecPayload(jp, nil), nil
}

type jobProposalAction struct {
	jpID graphql.ID
	name JobProposalAction
}

func (r *Resolver) executeJobProposalAction(ctx context.Context, action jobProposalAction) (*feeds.JobProposal, error) {
	if err := authenticateUser(ctx); err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(string(action.jpID), 10, 64)
	if err != nil {
		return nil, err
	}

	feedsSvc := r.App.GetFeedsService()

	switch action.name {
	case approve:
		err = feedsSvc.ApproveJobProposal(ctx, id)
	case cancel:
		err = feedsSvc.CancelJobProposal(ctx, id)
	case reject:
		err = feedsSvc.RejectJobProposal(ctx, id)
	default:
		return nil, errors.New("invalid job proposal action")
	}

	if err != nil {
		return nil, err
	}

	jp, err := r.App.GetFeedsService().GetJobProposal(id)
	if err != nil {
		return nil, err
	}

	return jp, nil
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

	session, ok := auth.GetGQLAuthenticatedSession(ctx)
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

	if err := r.App.GetConfig().SetLogSQLStatements(args.Input.Enabled); err != nil {
		return nil, err
	}

	return NewSetSQLLoggingPayload(args.Input.Enabled), nil
}
