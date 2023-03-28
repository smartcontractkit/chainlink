package feeds

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	pb "github.com/smartcontractkit/chainlink/v2/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr"
	ocr2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/crypto"
)

//go:generate mockery --quiet --name Service --output ./mocks/ --case=underscore
//go:generate mockery --quiet --dir ./proto --name FeedsManagerClient --output ./mocks/ --case=underscore

var (
	ErrOCR2Disabled         = errors.New("ocr2 is disabled")
	ErrOCRDisabled          = errors.New("ocr is disabled")
	ErrSingleFeedsManager   = errors.New("only a single feeds manager is supported")
	ErrJobAlreadyExists     = errors.New("a job for this contract address already exists - please use the 'force' option to replace it")
	ErrFeedsManagerDisabled = errors.New("feeds manager is disabled")

	promJobProposalRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feeds_job_proposal_requests",
		Help: "Metric to track job proposal requests",
	})

	promJobProposalCounts = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "feeds_job_proposal_count",
		Help: "Number of job proposals for the node partitioned by status.",
	}, []string{
		// Job Proposal status
		"status",
	})
)

// Service represents a behavior of the feeds service
type Service interface {
	Start(ctx context.Context) error
	Close() error

	CountManagers() (int64, error)
	GetManager(id int64) (*FeedsManager, error)
	ListManagers() ([]FeedsManager, error)
	ListManagersByIDs(ids []int64) ([]FeedsManager, error)
	RegisterManager(ctx context.Context, params RegisterManagerParams) (int64, error)
	UpdateManager(ctx context.Context, mgr FeedsManager) error

	CreateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error)
	DeleteChainConfig(ctx context.Context, id int64) (int64, error)
	GetChainConfig(id int64) (*ChainConfig, error)
	ListChainConfigsByManagerIDs(mgrIDs []int64) ([]ChainConfig, error)
	UpdateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error)

	DeleteJob(ctx context.Context, args *DeleteJobArgs) (int64, error)
	IsJobManaged(ctx context.Context, jobID int64) (bool, error)
	ProposeJob(ctx context.Context, args *ProposeJobArgs) (int64, error)
	RevokeJob(ctx context.Context, args *RevokeJobArgs) (int64, error)
	SyncNodeInfo(ctx context.Context, id int64) error

	CountJobProposalsByStatus() (*JobProposalCounts, error)
	GetJobProposal(id int64) (*JobProposal, error)
	ListJobProposals() ([]JobProposal, error)
	ListJobProposalsByManagersIDs(ids []int64) ([]JobProposal, error)

	ApproveSpec(ctx context.Context, id int64, force bool) error
	CancelSpec(ctx context.Context, id int64) error
	GetSpec(id int64) (*JobProposalSpec, error)
	ListSpecsByJobProposalIDs(ids []int64) ([]JobProposalSpec, error)
	RejectSpec(ctx context.Context, id int64) error
	UpdateSpecDefinition(ctx context.Context, id int64, spec string) error

	Unsafe_SetConnectionsManager(ConnectionsManager)
}

type service struct {
	utils.StartStopOnce

	orm          ORM
	jobORM       job.ORM
	q            pg.Q
	csaKeyStore  keystore.CSA
	p2pKeyStore  keystore.P2P
	ocr1KeyStore keystore.OCR
	ocr2KeyStore keystore.OCR2
	jobSpawner   job.Spawner
	cfg          Config
	connMgr      ConnectionsManager
	chainSet     evm.ChainSet
	lggr         logger.Logger
	version      string
}

// NewService constructs a new feeds service
func NewService(
	orm ORM,
	jobORM job.ORM,
	db *sqlx.DB,
	jobSpawner job.Spawner,
	keyStore keystore.Master,
	cfg Config,
	chainSet evm.ChainSet,
	lggr logger.Logger,
	version string,
) *service {
	lggr = lggr.Named("Feeds")
	svc := &service{
		orm:          orm,
		jobORM:       jobORM,
		q:            pg.NewQ(db, lggr, cfg),
		jobSpawner:   jobSpawner,
		p2pKeyStore:  keyStore.P2P(),
		csaKeyStore:  keyStore.CSA(),
		ocr1KeyStore: keyStore.OCR(),
		ocr2KeyStore: keyStore.OCR2(),
		cfg:          cfg,
		connMgr:      newConnectionsManager(lggr),
		chainSet:     chainSet,
		lggr:         lggr,
		version:      version,
	}

	return svc
}

type RegisterManagerParams struct {
	Name         string
	URI          string
	PublicKey    crypto.PublicKey
	ChainConfigs []ChainConfig
}

// RegisterManager registers a new ManagerService and attempts to establish a
// connection.
//
// Only a single feeds manager is currently supported.
func (s *service) RegisterManager(ctx context.Context, params RegisterManagerParams) (int64, error) {
	count, err := s.CountManagers()
	if err != nil {
		return 0, err
	}
	if count >= 1 {
		return 0, ErrSingleFeedsManager
	}

	mgr := FeedsManager{
		Name:      params.Name,
		URI:       params.URI,
		PublicKey: params.PublicKey,
	}

	var id int64
	q := s.q.WithOpts(pg.WithParentCtx(context.Background()))
	err = q.Transaction(func(tx pg.Queryer) error {
		var txerr error

		id, txerr = s.orm.CreateManager(&mgr, pg.WithQueryer(tx))
		if err != nil {
			return txerr
		}

		if _, txerr = s.orm.CreateBatchChainConfig(params.ChainConfigs, pg.WithQueryer(tx)); txerr != nil {
			return txerr
		}

		return nil
	})

	privkey, err := s.getCSAPrivateKey()
	if err != nil {
		return 0, err
	}

	// Establish a connection
	mgr.ID = id
	s.connectFeedManager(ctx, mgr, privkey)

	return id, nil
}

// SyncNodeInfo syncs the node's information with FMS
func (s *service) SyncNodeInfo(ctx context.Context, id int64) error {
	// Get the FMS RPC client
	fmsClient, err := s.connMgr.GetClient(id)
	if err != nil {
		return errors.Wrap(err, "could not fetch client")
	}

	cfgs, err := s.orm.ListChainConfigsByManagerIDs([]int64{id})
	if err != nil {
		return errors.Wrap(err, "could not fetch chain configs")
	}

	cfgMsgs := make([]*pb.ChainConfig, 0, len(cfgs))
	for _, cfg := range cfgs {
		cfgMsg, msgErr := s.newChainConfigMsg(cfg)
		if msgErr != nil {
			s.lggr.Errorf("SyncNodeInfo: %v", msgErr)

			continue
		}

		cfgMsgs = append(cfgMsgs, cfgMsg)
	}

	if _, err = fmsClient.UpdateNode(ctx, &pb.UpdateNodeRequest{
		Version:      s.version,
		ChainConfigs: cfgMsgs,
	}); err != nil {
		return err
	}

	return nil
}

// UpdateManager updates the feed manager details, takes down the
// connection and reestablishes a new connection with the updated public key.
func (s *service) UpdateManager(ctx context.Context, mgr FeedsManager) error {
	q := s.q.WithOpts(pg.WithParentCtx(ctx))
	err := q.Transaction(func(tx pg.Queryer) error {
		txerr := s.orm.UpdateManager(mgr, pg.WithQueryer(tx))
		if txerr != nil {
			return errors.Wrap(txerr, "could not update manager")
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err := s.restartConnection(ctx, mgr); err != nil {
		s.lggr.Errorf("could not restart FMS connection: %w", err)
	}

	return nil
}

// ListManagerServices lists all the manager services.
func (s *service) ListManagers() ([]FeedsManager, error) {
	managers, err := s.orm.ListManagers()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a list of managers")
	}

	for i := range managers {
		managers[i].IsConnectionActive = s.connMgr.IsConnected(managers[i].ID)
	}

	return managers, nil
}

// GetManager gets a manager service by id.
func (s *service) GetManager(id int64) (*FeedsManager, error) {
	manager, err := s.orm.GetManager(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get manager by ID")
	}

	manager.IsConnectionActive = s.connMgr.IsConnected(manager.ID)
	return manager, nil
}

// ListManagersByIDs get managers services by ids.
func (s *service) ListManagersByIDs(ids []int64) ([]FeedsManager, error) {
	managers, err := s.orm.ListManagersByIDs(ids)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list managers by IDs")
	}

	for _, manager := range managers {
		manager.IsConnectionActive = s.connMgr.IsConnected(manager.ID)
	}

	return managers, nil
}

// CountManagers gets the total number of manager services
func (s *service) CountManagers() (int64, error) {
	return s.orm.CountManagers()
}

// CreateChainConfig creates a chain config.
func (s *service) CreateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error) {
	var err error
	if cfg.AdminAddress != "" {
		_, err = common.NewMixedcaseAddressFromString(cfg.AdminAddress)
		if err != nil {
			return 0, fmt.Errorf("invalid admin address: %v", cfg.AdminAddress)
		}
	}

	id, err := s.orm.CreateChainConfig(cfg)
	if err != nil {
		return 0, errors.Wrap(err, "CreateChainConfig failed")
	}

	mgr, err := s.orm.GetManager(cfg.FeedsManagerID)
	if err != nil {
		return 0, errors.Wrap(err, "CreateChainConfig: failed to fetch manager")
	}

	if err := s.SyncNodeInfo(ctx, mgr.ID); err != nil {
		s.lggr.Infof("FMS: Unable to sync node info: %w", err)
	}

	return id, nil
}

// DeleteChainConfig deletes the chain config by id.
func (s *service) DeleteChainConfig(ctx context.Context, id int64) (int64, error) {
	cfg, err := s.orm.GetChainConfig(id)
	if err != nil {
		return 0, errors.Wrap(err, "DeleteChainConfig failed: could not get chain config")
	}

	_, err = s.orm.DeleteChainConfig(id)
	if err != nil {
		return 0, errors.Wrap(err, "DeleteChainConfig failed")
	}

	mgr, err := s.orm.GetManager(cfg.FeedsManagerID)
	if err != nil {
		return 0, errors.Wrap(err, "DeleteChainConfig: failed to fetch manager")
	}

	if err := s.SyncNodeInfo(ctx, mgr.ID); err != nil {
		s.lggr.Infof("FMS: Unable to sync node info: %w", err)
	}

	return id, nil
}

func (s *service) GetChainConfig(id int64) (*ChainConfig, error) {
	cfg, err := s.orm.GetChainConfig(id)
	if err != nil {
		return nil, errors.Wrap(err, "GetChainConfig failed")
	}

	return cfg, nil
}

func (s *service) ListChainConfigsByManagerIDs(mgrIDs []int64) ([]ChainConfig, error) {
	cfgs, err := s.orm.ListChainConfigsByManagerIDs(mgrIDs)

	return cfgs, errors.Wrap(err, "ListChainConfigsByManagerIDs failed")
}

func (s *service) UpdateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error) {
	var err error
	if cfg.AdminAddress != "" {
		_, err = common.NewMixedcaseAddressFromString(cfg.AdminAddress)
		if err != nil {
			return 0, fmt.Errorf("invalid admin address: %v", cfg.AdminAddress)
		}
	}

	id, err := s.orm.UpdateChainConfig(cfg)
	if err != nil {
		return 0, errors.Wrap(err, "UpdateChainConfig failed")
	}

	ccfg, err := s.orm.GetChainConfig(cfg.ID)
	if err != nil {
		return 0, errors.Wrap(err, "UpdateChainConfig failed: could not get chain config")
	}

	if err := s.SyncNodeInfo(ctx, ccfg.FeedsManagerID); err != nil {
		s.lggr.Infof("FMS: Unable to sync node info: %w", err)
	}

	return id, nil
}

// Lists all JobProposals
//
// When we support multiple feed managers, we will need to change this to filter
// by feeds manager
func (s *service) ListJobProposals() ([]JobProposal, error) {
	return s.orm.ListJobProposals()
}

// ListJobProposalsByManagersIDs gets job proposals by feeds managers IDs
func (s *service) ListJobProposalsByManagersIDs(ids []int64) ([]JobProposal, error) {
	return s.orm.ListJobProposalsByManagersIDs(ids)
}

// DeleteJobArgs are the arguments to provide to the DeleteJob method.
type DeleteJobArgs struct {
	FeedsManagerID int64
	RemoteUUID     uuid.UUID
}

// DeleteJob deletes a job proposal if it exist. The feeds manager id check
// ensures that only the intended feed manager can make this request.
func (s *service) DeleteJob(ctx context.Context, args *DeleteJobArgs) (int64, error) {
	proposal, err := s.orm.GetJobProposalByRemoteUUID(args.RemoteUUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, errors.Wrap(err, "GetJobProposalByRemoteUUID failed to check existence of job proposal")
		}

		return 0, errors.Wrap(err, "GetJobProposalByRemoteUUID did not find any proposals to delete")
	}

	// Ensure that if the job proposal exists, that it belongs to the feeds
	// manager which previously proposed a job using the remote UUID.
	if args.FeedsManagerID != proposal.FeedsManagerID {
		return 0, errors.New("cannot delete a job proposal belonging to another feeds manager")
	}

	pctx := pg.WithParentCtx(ctx)
	if err = s.orm.DeleteProposal(proposal.ID, pctx); err != nil {
		s.lggr.Errorw("Failed to delete the proposal", "error", err)

		return 0, errors.Wrap(err, "DeleteProposal failed")
	}

	return proposal.ID, nil
}

// RevokeJobArgs are the arguments to provide the RevokeJob method
type RevokeJobArgs struct {
	FeedsManagerID int64
	RemoteUUID     uuid.UUID
}

// RevokeJob revokes a pending job proposal if it exist. The feeds manager
// id check ensures that only the intended feed manager can make this request.
func (s *service) RevokeJob(ctx context.Context, args *RevokeJobArgs) (int64, error) {
	proposal, err := s.orm.GetJobProposalByRemoteUUID(args.RemoteUUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, errors.Wrap(err, "GetJobProposalByRemoteUUID failed to check existence of job proposal")
		}

		return 0, errors.Wrap(err, "GetJobProposalByRemoteUUID did not find any proposals to revoke")
	}

	// Ensure that if the job proposal exists, that it belongs to the feeds
	// manager which previously proposed a job using the remote UUID.
	if args.FeedsManagerID != proposal.FeedsManagerID {
		return 0, errors.New("cannot revoke a job proposal belonging to another feeds manager")
	}

	// get the latest spec for the proposal
	latest, err := s.orm.GetLatestSpec(proposal.ID)
	if err != nil {
		return 0, errors.Wrap(err, "GetLatestSpec failed to get latest spec")
	}

	if canRevoke := s.isRevokable(proposal.Status, latest.Status); !canRevoke {
		return 0, errors.New("only pending job proposals can be revoked")
	}

	pctx := pg.WithParentCtx(ctx)
	if err = s.orm.RevokeSpec(latest.ID, pctx); err != nil {
		s.lggr.Errorw("Failed to revoke the proposal", "error", err)

		return 0, errors.Wrap(err, "RevokeSpec failed")
	}

	return proposal.ID, nil
}

// ProposeJobArgs are the arguments to provide to the ProposeJob method.
type ProposeJobArgs struct {
	FeedsManagerID int64
	RemoteUUID     uuid.UUID
	Multiaddrs     pq.StringArray
	Version        int32
	Spec           string
}

// ProposeJob creates a job proposal if it does not exist. If it already exists
// and a new version is provided, a new spec is created.
//
// The feeds manager id check exists for support of multiple feeds managers in
// the future so that in the (very slim) off chance that the same uuid is
// generated by another feeds manager or they maliciously send an existing uuid
// belonging to another feeds manager, we do not update it.
func (s *service) ProposeJob(ctx context.Context, args *ProposeJobArgs) (int64, error) {
	// Validate the args
	if err := s.validateProposeJobArgs(*args); err != nil {
		return 0, err
	}

	existing, err := s.orm.GetJobProposalByRemoteUUID(args.RemoteUUID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return 0, errors.Wrap(err, "failed to check existence of job proposal")
		}
	}

	// Validation for existing job proposals
	if err == nil {
		// Ensure that if the job proposal exists, that it belongs to the feeds
		// manager which previously proposed a job using the remote UUID.
		if args.FeedsManagerID != existing.FeedsManagerID {
			return 0, errors.New("cannot update a job proposal belonging to another feeds manager")
		}

		// Check the version being proposed has not been previously proposed.
		var exists bool
		exists, err = s.orm.ExistsSpecByJobProposalIDAndVersion(existing.ID, args.Version)
		if err != nil {
			return 0, errors.Wrap(err, "failed to check existence of spec")
		}

		if exists {
			return 0, errors.New("proposed job spec version already exists")
		}
	}

	var id int64
	q := s.q.WithOpts(pg.WithParentCtx(ctx))
	err = q.Transaction(func(tx pg.Queryer) error {
		var txerr error

		// Parse the Job Spec TOML to extract the name
		name := extractName(args.Spec)

		// Upsert job proposal
		id, txerr = s.orm.UpsertJobProposal(&JobProposal{
			Name:           name,
			RemoteUUID:     args.RemoteUUID,
			Status:         JobProposalStatusPending,
			FeedsManagerID: args.FeedsManagerID,
			Multiaddrs:     args.Multiaddrs,
		}, pg.WithQueryer(tx))
		if txerr != nil {
			return errors.Wrap(txerr, "failed to upsert job proposal")
		}

		// Create the spec version
		_, txerr = s.orm.CreateSpec(JobProposalSpec{
			Definition:    args.Spec,
			Status:        SpecStatusPending,
			Version:       args.Version,
			JobProposalID: id,
		}, pg.WithQueryer(tx))
		if txerr != nil {
			return errors.Wrap(txerr, "failed to create spec")
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	// Track the given job proposal request
	promJobProposalRequest.Inc()

	if err = s.observeJobProposalCounts(); err != nil {
		return 0, err
	}

	return id, nil
}

// GetJobProposal gets a job proposal by id.
func (s *service) GetJobProposal(id int64) (*JobProposal, error) {
	return s.orm.GetJobProposal(id)
}

// CountJobProposalsByStatus returns the count of job proposals with a given status.
func (s *service) CountJobProposalsByStatus() (*JobProposalCounts, error) {
	return s.orm.CountJobProposalsByStatus()
}

// RejectSpec rejects a spec.
func (s *service) RejectSpec(ctx context.Context, id int64) error {
	pctx := pg.WithParentCtx(ctx)

	spec, err := s.orm.GetSpec(id, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal spec")
	}

	// Validate
	if spec.Status != SpecStatusPending {
		return errors.New("must be a pending job proposal spec")
	}

	proposal, err := s.orm.GetJobProposal(spec.JobProposalID, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal")
	}

	fmsClient, err := s.connMgr.GetClient(proposal.FeedsManagerID)
	if err != nil {
		return errors.Wrap(err, "fms rpc client is not connected")
	}

	q := s.q.WithOpts(pctx)
	err = q.Transaction(func(tx pg.Queryer) error {
		if err = s.orm.RejectSpec(id, pg.WithQueryer(tx)); err != nil {
			return err
		}

		if _, err = fmsClient.RejectedJob(ctx, &pb.RejectedJobRequest{
			Uuid:    proposal.RemoteUUID.String(),
			Version: int64(spec.Version),
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "could not reject job proposal")
	}

	err = s.observeJobProposalCounts()

	return err
}

// IsJobManaged determines is a job is managed by the Feeds Manager.
func (s *service) IsJobManaged(ctx context.Context, jobID int64) (bool, error) {
	return s.orm.IsJobManaged(jobID, pg.WithParentCtx(ctx))
}

// ApproveSpec approves a spec for a job proposal and creates a job with the
// spec.
func (s *service) ApproveSpec(ctx context.Context, id int64, force bool) error {
	pctx := pg.WithParentCtx(ctx)

	spec, err := s.orm.GetSpec(id, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal spec")
	}

	proposal, err := s.orm.GetJobProposal(spec.JobProposalID, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal")
	}

	if err = s.isApprovable(proposal.Status, proposal.ID, spec.Status, spec.ID); err != nil {
		return err
	}

	logger := s.lggr.With(
		"job_proposal_id", proposal.ID,
		"job_proposal_spec_id", id,
	)

	fmsClient, err := s.connMgr.GetClient(proposal.FeedsManagerID)
	if err != nil {
		logger.Errorw("Failed to get FMS Client", "error", err)

		return errors.Wrap(err, "fms rpc client")
	}

	j, err := s.generateJob(spec.Definition)
	if err != nil {
		return errors.Wrap(err, "could not generate job from spec")
	}

	// Check that the bridges exist
	if err = s.jobORM.AssertBridgesExist(j.Pipeline); err != nil {
		logger.Errorw("Failed to approve job spec due to bridge check", "err", err.Error())

		return errors.Wrap(err, "failed to approve job spec due to bridge check")
	}

	address, evmChainID, err := s.getAddressAndEVMChainIDFromJob(j)
	if err != nil {
		return err
	}

	q := s.q.WithOpts(pctx)
	err = q.Transaction(func(tx pg.Queryer) error {
		var (
			txerr error

			pgOpts = pg.WithQueryer(tx)
		)

		// Remove the existing job, continuing if no job is found
		existingJobID, txerr := s.jobORM.FindJobIDByAddress(address, evmChainID, pgOpts)
		if txerr != nil {
			// Return an error if the repository errors. If there is a not found
			// error we want to continue with approving the job.
			if !errors.Is(txerr, sql.ErrNoRows) {
				return errors.Wrap(txerr, "FindJobIDByAddress failed")
			}
		}

		// Remove the existing job since a job was found
		if txerr == nil {
			// Do not proceed to remove the running job unless the force flag is true
			if !force {
				return ErrJobAlreadyExists
			}

			// Check if the job is managed by FMS
			approvedSpec, serr := s.orm.GetApprovedSpec(proposal.ID, pgOpts)
			if serr != nil {
				if !errors.Is(serr, sql.ErrNoRows) {
					logger.Errorw("Failed to get approved spec", "error", serr)

					// Return an error for any other errors fetching the
					// approved spec
					return errors.Wrap(serr, "GetApprovedSpec failed")
				}
			}

			// If a spec is found, cancel the existing job spec
			if serr == nil {
				if cerr := s.orm.CancelSpec(approvedSpec.ID, pgOpts); cerr != nil {
					logger.Errorw("Failed to delete the cancel the spec", "error", cerr)

					return cerr
				}
			}

			// Delete the job
			if serr = s.jobSpawner.DeleteJob(existingJobID, pgOpts); serr != nil {
				logger.Errorw("Failed to delete the job", "error", serr)

				return errors.Wrap(serr, "DeleteJob failed")
			}
		}

		// Create the job
		if txerr = s.jobSpawner.CreateJob(j, pgOpts); txerr != nil {
			logger.Errorw("Failed to create job", "error", txerr)

			return txerr
		}

		// Approve the job proposal spec
		if txerr = s.orm.ApproveSpec(id, j.ExternalJobID, pgOpts); txerr != nil {
			logger.Errorw("Failed to approve spec", "error", txerr)

			return txerr
		}

		// Send to FMS Client
		if _, txerr = fmsClient.ApprovedJob(ctx, &pb.ApprovedJobRequest{
			Uuid:    proposal.RemoteUUID.String(),
			Version: int64(spec.Version),
		}); txerr != nil {
			logger.Errorw("Failed to approve job to FMS", "error", txerr)

			return txerr
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "could not approve job proposal")
	}

	if err = s.observeJobProposalCounts(); err != nil {
		logger.Errorw("Failed to push metrics for job approval", err)
	}

	return nil
}

// CancelSpec cancels a spec for a job proposal.
func (s *service) CancelSpec(ctx context.Context, id int64) error {
	pctx := pg.WithParentCtx(ctx)

	spec, err := s.orm.GetSpec(id, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal spec")
	}

	if spec.Status != SpecStatusApproved {
		return errors.New("must be an approved job proposal spec")
	}

	jp, err := s.orm.GetJobProposal(spec.JobProposalID, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "orm: job proposal")
	}

	fmsClient, err := s.connMgr.GetClient(jp.FeedsManagerID)
	if err != nil {
		return errors.Wrap(err, "fms rpc client")
	}

	q := s.q.WithOpts(pctx)
	err = q.Transaction(func(tx pg.Queryer) error {
		if err = s.orm.CancelSpec(id, pg.WithQueryer(tx)); err != nil {
			return err
		}

		// Delete the job
		var j job.Job
		j, err = s.jobORM.FindJobByExternalJobID(jp.ExternalJobID.UUID, pg.WithQueryer(tx))
		if err != nil {
			return errors.Wrap(err, "FindJobByExternalJobID failed")
		}

		if err = s.jobSpawner.DeleteJob(j.ID, pg.WithQueryer(tx)); err != nil {
			return errors.Wrap(err, "DeleteJob failed")
		}

		// Send to FMS Client
		if _, err = fmsClient.CancelledJob(ctx, &pb.CancelledJobRequest{
			Uuid:    jp.RemoteUUID.String(),
			Version: int64(spec.Version),
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = s.observeJobProposalCounts()

	return err
}

// ListSpecsByJobProposalIDs gets the specs which belong to the job proposal ids.
func (s *service) ListSpecsByJobProposalIDs(ids []int64) ([]JobProposalSpec, error) {
	return s.orm.ListSpecsByJobProposalIDs(ids)
}

// GetSpec gets the spec details by id.
func (s *service) GetSpec(id int64) (*JobProposalSpec, error) {
	return s.orm.GetSpec(id)
}

// UpdateSpecDefinition updates the spec's TOML definition.
func (s *service) UpdateSpecDefinition(ctx context.Context, id int64, defn string) error {
	pctx := pg.WithParentCtx(ctx)

	spec, err := s.orm.GetSpec(id, pctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Wrap(err, "job proposal spec does not exist")
		}

		return errors.Wrap(err, "database error")
	}

	if !spec.CanEditDefinition() {
		return errors.New("must be a pending or cancelled spec")
	}

	// Update the spec definition
	if err = s.orm.UpdateSpecDefinition(id, defn, pctx); err != nil {
		return errors.Wrap(err, "could not update job proposal")
	}

	return nil
}

// Start starts the service.
func (s *service) Start(ctx context.Context) error {
	return s.StartOnce("FeedsService", func() error {
		privkey, err := s.getCSAPrivateKey()
		if err != nil {
			return err
		}

		// We only support a single feeds manager right now
		mgrs, err := s.ListManagers()
		if err != nil {
			return err
		}
		if len(mgrs) < 1 {
			s.lggr.Info("no feeds managers registered")

			return nil
		}

		mgr := mgrs[0]
		s.connectFeedManager(ctx, mgr, privkey)

		err = s.observeJobProposalCounts()

		return err
	})
}

// Close shuts down the service
func (s *service) Close() error {
	return s.StopOnce("FeedsService", func() error {
		// This blocks until it finishes
		s.connMgr.Close()

		return nil
	})
}

// connectFeedManager connects to a feeds manager
func (s *service) connectFeedManager(ctx context.Context, mgr FeedsManager, privkey []byte) {
	s.connMgr.Connect(ConnectOpts{
		FeedsManagerID: mgr.ID,
		URI:            mgr.URI,
		Privkey:        privkey,
		Pubkey:         mgr.PublicKey,
		Handlers: &RPCHandlers{
			feedsManagerID: mgr.ID,
			svc:            s,
		},
		OnConnect: func(pb.FeedsManagerClient) {
			// Sync the node's information with FMS once connected
			err := s.SyncNodeInfo(ctx, mgr.ID)
			if err != nil {
				s.lggr.Infof("Error syncing node info: %v", err)
			}
		},
	})
}

// getCSAPrivateKey gets the server's CSA private key
func (s *service) getCSAPrivateKey() (privkey []byte, err error) {
	// Fetch the server's public key
	keys, err := s.csaKeyStore.GetAll()
	if err != nil {
		return privkey, err
	}
	if len(keys) < 1 {
		return privkey, errors.New("CSA key does not exist")
	}
	return keys[0].Raw(), nil
}

// observeJobProposalCounts is a helper method that queries the repository for the count of
// job proposals by status and then updates prometheus gauges.
func (s *service) observeJobProposalCounts() error {
	counts, err := s.CountJobProposalsByStatus()
	if err != nil {
		return errors.Wrap(err, "failed to fetch counts of job proposals")
	}

	// Transform counts into prometheus metrics.
	metrics := counts.toMetrics()

	// Set the prometheus gauge metrics.
	for _, status := range []JobProposalStatus{JobProposalStatusPending, JobProposalStatusApproved,
		JobProposalStatusCancelled, JobProposalStatusRejected} {

		status := status

		promJobProposalCounts.With(prometheus.Labels{"status": string(status)}).Set(metrics[status])
	}

	return nil
}

// Unsafe_SetConnectionsManager sets the ConnectionsManager on the service.
//
// We need to be able to inject a mock for the client to facilitate integration
// tests.
//
// ONLY TO BE USED FOR TESTING.
func (s *service) Unsafe_SetConnectionsManager(connMgr ConnectionsManager) {
	s.connMgr = connMgr
}

// getAddressAndChainIDFromJob extracts the address and evmChainID from a job
func (s *service) getAddressAndEVMChainIDFromJob(j *job.Job) (ethkey.EIP55Address, *utils.Big, error) {
	var address ethkey.EIP55Address
	var evmChainID *utils.Big

	switch j.Type {
	case job.OffchainReporting:
		address = j.OCROracleSpec.ContractAddress
		evmChainID = j.OCROracleSpec.EVMChainID
	case job.OffchainReporting2:
		eipAddress, addrErr := ethkey.NewEIP55Address(j.OCR2OracleSpec.ContractID)
		if addrErr != nil {
			return eipAddress, nil, errors.Wrap(addrErr, "failed to create EIP55Address from OCR2 job spec")
		}

		evmChain, chainErr := job.EVMChainForJob(j, s.chainSet)
		if chainErr != nil {
			return eipAddress, nil, errors.Wrap(chainErr, "failed to get evmChainID from OCR2 job spec")
		}

		evmChainID = utils.NewBig(evmChain.ID())
		address = eipAddress
	case job.Bootstrap:
		eipAddress, addrErr := ethkey.NewEIP55Address(j.BootstrapSpec.ContractID)
		if addrErr != nil {
			return eipAddress, nil, errors.Wrap(addrErr, "failed to create EIP55Address from Bootstrap job spec")
		}

		evmChain, chainErr := job.EVMChainForBootstrapJob(j, s.chainSet)
		if chainErr != nil {
			return eipAddress, nil, errors.Wrap(chainErr, "failed to get evmChainID from Bootstrap job spec")
		}

		evmChainID = utils.NewBig(evmChain.ID())
		address = eipAddress
	case job.FluxMonitor:
		address = j.FluxMonitorSpec.ContractAddress
		evmChainID = j.FluxMonitorSpec.EVMChainID
	default:
		return address, nil, errors.Errorf("unsupported job type when approving job proposal specs: %s", j.Type)
	}

	return address, evmChainID, nil
}

// generateJob validates and generates a job from a spec.
func (s *service) generateJob(spec string) (*job.Job, error) {
	jobType, err := job.ValidateSpec(spec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse job spec TOML")
	}

	var js job.Job
	switch jobType {
	case job.OffchainReporting:
		if !s.cfg.Dev() && !s.cfg.FeatureOffchainReporting() {
			return nil, ErrOCRDisabled
		}
		js, err = ocr.ValidatedOracleSpecToml(s.chainSet, spec)
	case job.OffchainReporting2:
		if !s.cfg.Dev() && !s.cfg.FeatureOffchainReporting2() {
			return nil, ErrOCR2Disabled
		}
		js, err = ocr2.ValidatedOracleSpecToml(s.cfg, spec)
	case job.Bootstrap:
		if !s.cfg.Dev() && !s.cfg.FeatureOffchainReporting2() {
			return nil, ErrOCR2Disabled
		}
		js, err = ocrbootstrap.ValidatedBootstrapSpecToml(spec)
	case job.FluxMonitor:
		js, err = fluxmonitorv2.ValidatedFluxMonitorSpec(s.cfg, spec)
	default:
		return nil, errors.Errorf("unknown job type: %s", jobType)

	}
	if err != nil {
		return nil, err
	}

	return &js, nil
}

// newChainConfigMsg generates a chain config protobuf message.
func (s *service) newChainConfigMsg(cfg ChainConfig) (*pb.ChainConfig, error) {
	// Only supports EVM Chains
	if cfg.ChainType != "EVM" {
		return nil, errors.New("unsupported chain type")
	}

	ocr1Cfg, err := s.newOCR1ConfigMsg(cfg.OCR1Config)
	if err != nil {
		return nil, err
	}

	ocr2Cfg, err := s.newOCR2ConfigMsg(cfg.OCR2Config)
	if err != nil {
		return nil, err
	}

	return &pb.ChainConfig{
		Chain: &pb.Chain{
			Id:   cfg.ChainID,
			Type: pb.ChainType_CHAIN_TYPE_EVM,
		},
		AccountAddress:    cfg.AccountAddress,
		AdminAddress:      cfg.AdminAddress,
		FluxMonitorConfig: s.newFluxMonitorConfigMsg(cfg.FluxMonitorConfig),
		Ocr1Config:        ocr1Cfg,
		Ocr2Config:        ocr2Cfg,
	}, nil
}

// newFMConfigMsg generates a FMConfig protobuf message. Flux Monitor does not
// have any configuration but this is here for consistency.
func (*service) newFluxMonitorConfigMsg(cfg FluxMonitorConfig) *pb.FluxMonitorConfig {
	return &pb.FluxMonitorConfig{Enabled: cfg.Enabled}
}

// newOCR1ConfigMsg generates a OCR1Config protobuf message.
func (s *service) newOCR1ConfigMsg(cfg OCR1Config) (*pb.OCR1Config, error) {
	if !cfg.Enabled {
		return &pb.OCR1Config{Enabled: false}, nil
	}

	msg := &pb.OCR1Config{
		Enabled:     true,
		IsBootstrap: cfg.IsBootstrap,
		Multiaddr:   cfg.Multiaddr.ValueOrZero(),
	}

	// Fetch the P2P key bundle
	if cfg.P2PPeerID.Valid {
		peerID, err := p2pkey.MakePeerID(cfg.P2PPeerID.String)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid peer id: %s", cfg.P2PPeerID.String)
		}
		p2pKey, err := s.p2pKeyStore.Get(peerID)
		if err != nil {
			return nil, errors.Wrapf(err, "p2p key not found: %s", cfg.P2PPeerID.String)
		}

		msg.P2PKeyBundle = &pb.OCR1Config_P2PKeyBundle{
			PeerId:    p2pKey.PeerID().String(),
			PublicKey: p2pKey.PublicKeyHex(),
		}
	}

	if cfg.KeyBundleID.Valid {
		ocrKey, err := s.ocr1KeyStore.Get(cfg.KeyBundleID.String)
		if err != nil {
			return nil, errors.Wrapf(err, "ocr key not found: %s", cfg.KeyBundleID.String)
		}

		msg.OcrKeyBundle = &pb.OCR1Config_OCRKeyBundle{
			BundleId:              ocrKey.GetID(),
			ConfigPublicKey:       ocrkey.ConfigPublicKey(ocrKey.PublicKeyConfig()).String(),
			OffchainPublicKey:     ocrKey.OffChainSigning.PublicKey().String(),
			OnchainSigningAddress: ocrKey.OnChainSigning.Address().String(),
		}
	}

	return msg, nil
}

// newOCR2ConfigMsg generates a OCR2Config protobuf message.
func (s *service) newOCR2ConfigMsg(cfg OCR2Config) (*pb.OCR2Config, error) {
	if !cfg.Enabled {
		return &pb.OCR2Config{Enabled: false}, nil
	}

	msg := &pb.OCR2Config{
		Enabled:     true,
		IsBootstrap: cfg.IsBootstrap,
		Multiaddr:   cfg.Multiaddr.ValueOrZero(),
		Plugins: &pb.OCR2Config_Plugins{
			Commit:  cfg.Plugins.Commit,
			Execute: cfg.Plugins.Execute,
			Median:  cfg.Plugins.Median,
			Mercury: cfg.Plugins.Mercury,
		},
	}

	// Fetch the P2P key bundle
	if cfg.P2PPeerID.Valid {
		peerID, err := p2pkey.MakePeerID(cfg.P2PPeerID.String)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid peer id: %s", cfg.P2PPeerID.String)
		}
		p2pKey, err := s.p2pKeyStore.Get(peerID)
		if err != nil {
			return nil, errors.Wrapf(err, "p2p key not found: %s", cfg.P2PPeerID.String)
		}

		msg.P2PKeyBundle = &pb.OCR2Config_P2PKeyBundle{
			PeerId:    p2pKey.PeerID().String(),
			PublicKey: p2pKey.PublicKeyHex(),
		}
	}

	// Fetch the OCR Key Bundle
	if cfg.KeyBundleID.Valid {
		ocrKey, err := s.ocr2KeyStore.Get(cfg.KeyBundleID.String)
		if err != nil {
			return nil, errors.Wrapf(err, "ocr key not found: %s", cfg.KeyBundleID.String)
		}

		ocrConfigPublicKey := ocrKey.ConfigEncryptionPublicKey()
		ocrOffChainPublicKey := ocrKey.OffchainPublicKey()

		msg.OcrKeyBundle = &pb.OCR2Config_OCRKeyBundle{
			BundleId:              ocrKey.ID(),
			ConfigPublicKey:       hex.EncodeToString(ocrConfigPublicKey[:]),
			OffchainPublicKey:     hex.EncodeToString(ocrOffChainPublicKey[:]),
			OnchainSigningAddress: ocrKey.OnChainPublicKey(),
		}
	}

	return msg, nil
}

func (s *service) validateProposeJobArgs(args ProposeJobArgs) error {
	// Validate the job spec
	j, err := s.generateJob(args.Spec)
	if err != nil {
		return errors.Wrap(err, "failed to generate a job based on spec")
	}

	// Validate bootstrap multiaddrs which are only allowed for OCR jobs
	if len(args.Multiaddrs) > 0 && j.Type != job.OffchainReporting && j.Type != job.OffchainReporting2 {
		return errors.New("only OCR job type supports multiaddr")
	}

	return nil
}

func (s *service) restartConnection(ctx context.Context, mgr FeedsManager) error {
	s.lggr.Infof("Restarting connection")

	if err := s.connMgr.Disconnect(mgr.ID); err != nil {
		s.lggr.Info("Feeds Manager not connected, attempting to connect")
	}

	// Establish a new connection
	privkey, err := s.getCSAPrivateKey()
	if err != nil {
		return err
	}

	s.connectFeedManager(ctx, mgr, privkey)

	return nil
}

// extractName extracts the name from the TOML returning an null string if
// there is an error.
func extractName(defn string) null.String {
	spec := struct {
		Name null.String
	}{}

	if err := toml.Unmarshal([]byte(defn), &spec); err != nil {
		return null.StringFromPtr(nil)
	}

	return spec.Name
}

// isApprovable returns nil if a spec can be approved based on the current
// proposal and spec status, and if it can't be approved, the reason as an
// error.
func (s *service) isApprovable(propStatus JobProposalStatus, proposalID int64, specStatus SpecStatus, specID int64) error {
	if propStatus == JobProposalStatusDeleted {
		return errors.New("cannot approve spec for a deleted job proposal")
	}

	if propStatus == JobProposalStatusRevoked {
		return errors.New("cannot approve spec for a revoked job proposal")
	}

	switch specStatus {
	case SpecStatusApproved:
		return errors.New("cannot approve an approved spec")
	case SpecStatusRejected:
		return errors.New("cannot approve a rejected spec")
	case SpecStatusRevoked:
		return errors.New("cannot approve a revoked spec")
	case SpecStatusCancelled:
		// Allowed to approve a cancelled job if it is the latest job
		latest, serr := s.orm.GetLatestSpec(proposalID)
		if serr != nil {
			return errors.Wrap(serr, "failed to get latest spec")
		}

		if latest.ID != specID {
			return errors.New("cannot approve a cancelled spec")
		}

		return nil
	case SpecStatusPending:
		return nil
	default:
		return errors.New("invalid job spec status")
	}
}

func (s *service) isRevokable(propStatus JobProposalStatus, specStatus SpecStatus) bool {
	return propStatus == JobProposalStatusPending && specStatus == SpecStatusPending
}

var _ Service = &NullService{}

// NullService defines an implementation of the Feeds Service that is used
// when the Feeds Service is disabled.
type NullService struct{}

//revive:disable
func (ns NullService) Start(ctx context.Context) error { return nil }
func (ns NullService) Close() error                    { return nil }
func (ns NullService) ApproveSpec(ctx context.Context, id int64, force bool) error {
	return ErrFeedsManagerDisabled
}
func (ns NullService) CountManagers() (int64, error) { return 0, nil }
func (ns NullService) CountJobProposalsByStatus() (*JobProposalCounts, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) CancelSpec(ctx context.Context, id int64) error {
	return ErrFeedsManagerDisabled
}
func (ns NullService) GetJobProposal(id int64) (*JobProposal, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) ListSpecsByJobProposalIDs(ids []int64) ([]JobProposalSpec, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) GetManager(id int64) (*FeedsManager, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) ListManagersByIDs(ids []int64) ([]FeedsManager, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) GetSpec(id int64) (*JobProposalSpec, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) ListManagers() ([]FeedsManager, error) { return nil, nil }
func (ns NullService) CreateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) GetChainConfig(id int64) (*ChainConfig, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) DeleteChainConfig(ctx context.Context, id int64) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) ListChainConfigsByManagerIDs(mgrIDs []int64) ([]ChainConfig, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) UpdateChainConfig(ctx context.Context, cfg ChainConfig) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) ListJobProposals() ([]JobProposal, error) { return nil, nil }
func (ns NullService) ListJobProposalsByManagersIDs(ids []int64) ([]JobProposal, error) {
	return nil, ErrFeedsManagerDisabled
}
func (ns NullService) ProposeJob(ctx context.Context, args *ProposeJobArgs) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) DeleteJob(ctx context.Context, args *DeleteJobArgs) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) RevokeJob(ctx context.Context, args *RevokeJobArgs) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) RegisterManager(ctx context.Context, params RegisterManagerParams) (int64, error) {
	return 0, ErrFeedsManagerDisabled
}
func (ns NullService) RejectSpec(ctx context.Context, id int64) error {
	return ErrFeedsManagerDisabled
}
func (ns NullService) SyncNodeInfo(ctx context.Context, id int64) error { return nil }
func (ns NullService) UpdateManager(ctx context.Context, mgr FeedsManager) error {
	return ErrFeedsManagerDisabled
}
func (ns NullService) IsJobManaged(ctx context.Context, jobID int64) (bool, error) {
	return false, nil
}
func (ns NullService) UpdateSpecDefinition(ctx context.Context, id int64, spec string) error {
	return ErrFeedsManagerDisabled
}
func (ns NullService) Unsafe_SetConnectionsManager(_ ConnectionsManager) {}

//revive:enable
