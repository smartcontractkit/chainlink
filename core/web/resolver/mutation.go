package resolver

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/auth"
	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/cron"
	"github.com/smartcontractkit/chainlink/core/services/directrequest"
	"github.com/smartcontractkit/chainlink/core/services/feeds"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keeper"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
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

type createFeedsManagerInput struct {
	Name                   string
	URI                    string
	PublicKey              string
	JobTypes               []JobType
	IsBootstrapPeer        bool
	BootstrapPeerMultiaddr *string
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
		Name:                   bridges.TaskType(args.Input.Name),
		URL:                    webURL,
		Confirmations:          uint32(args.Input.Confirmations),
		MinimumContractPayment: minContractPayment,
	}

	taskType, err := bridges.NewTaskType(string(args.ID))
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
	ID graphql.ID
}) (*DeleteNodePayloadResolver, error) {
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
			return NewDeleteNodePayloadResolver(nil, err), nil
		}

		return nil, err
	}

	err = r.App.EVMORM().DeleteNode(int64(id))
	if err != nil {
		if errors.Is(err, evm.ErrNoRowsAffected) {
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

	taskType, err := bridges.NewTaskType(string(args.ID))
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

	id, err := stringutils.ToInt64(string(args.ID))
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

	id, err := stringutils.ToInt64(string(action.jpID))
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

	chain, err := r.App.GetChainSet().Add(id.ToInt(), *chainCfg)
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

	chain, err := r.App.GetChainSet().Configure(id.ToInt(), args.Input.Enabled, *chainCfg)
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

	err = r.App.GetChainSet().Remove(id.ToInt())
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
		return NewCreateJobPayload(nil, map[string]string{
			"TOML spec": errors.Wrap(err, "failed to parse TOML").Error(),
		}), nil
	}

	var jb job.Job
	config := r.App.GetConfig()
	switch jbt {
	case job.OffchainReporting:
		jb, err = offchainreporting.ValidatedOracleSpecToml(r.App.GetChainSet(), args.Input.TOML)
		if !config.Dev() && !config.FeatureOffchainReporting() {
			return nil, errors.New("The Offchain Reporting feature is disabled by configuration")
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
	default:
		return NewCreateJobPayload(nil, map[string]string{
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

	return NewCreateJobPayload(&jb, nil), nil
}
