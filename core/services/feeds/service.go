package feeds

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	pb "github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/core/services/fluxmonitorv2"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore
//go:generate mockery --dir ./proto --name FeedsManagerClient --output ./mocks/ --case=underscore

var (
	ErrOCRDisabled        = errors.New("ocr is disabled")
	ErrSingleFeedsManager = errors.New("only a single feeds manager is supported")
	ErrBootstrapXorJobs   = errors.New("feeds manager cannot be bootstrap while having assigned job types")

	promJobProposalRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feeds_job_proposal_requests",
		Help: "Metric to track job proposal requests",
	})
)

// Service represents a behavior of the feeds service
type Service interface {
	Start() error
	Close() error

	CountManagers() (int64, error)
	GetManager(id int64) (*FeedsManager, error)
	ListManagersByIDs(ids []int64) ([]FeedsManager, error)
	ListManagers() ([]FeedsManager, error)
	RegisterManager(ms *FeedsManager) (int64, error)
	UpdateManager(ctx context.Context, mgr FeedsManager) error

	ProposeJob(ctx context.Context, args *ProposeJobArgs) (int64, error)
	SyncNodeInfo(id int64) error
	IsJobManaged(ctx context.Context, jobID int64) (bool, error)

	GetJobProposal(id int64) (*JobProposal, error)
	ListJobProposalsByManagersIDs(ids []int64) ([]JobProposal, error)
	ListJobProposals() ([]JobProposal, error)

	ApproveSpec(ctx context.Context, id int64) error
	CancelSpec(ctx context.Context, id int64) error
	GetSpec(id int64) (*JobProposalSpec, error)
	ListSpecsByJobProposalIDs(ids []int64) ([]JobProposalSpec, error)
	RejectSpec(ctx context.Context, id int64) error
	UpdateSpecDefinition(ctx context.Context, id int64, spec string) error

	Unsafe_SetConnectionsManager(ConnectionsManager)
}

type service struct {
	utils.StartStopOnce

	orm         ORM
	jobORM      job.ORM
	q           pg.Q
	csaKeyStore keystore.CSA
	ethKeyStore keystore.Eth
	p2pKeyStore keystore.P2P
	jobSpawner  job.Spawner
	cfg         Config
	connMgr     ConnectionsManager
	chainSet    evm.ChainSet
	lggr        logger.Logger
	version     string
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
		orm:         orm,
		jobORM:      jobORM,
		q:           pg.NewQ(db, lggr, cfg),
		jobSpawner:  jobSpawner,
		p2pKeyStore: keyStore.P2P(),
		csaKeyStore: keyStore.CSA(),
		ethKeyStore: keyStore.Eth(),
		cfg:         cfg,
		connMgr:     newConnectionsManager(lggr),
		chainSet:    chainSet,
		lggr:        lggr,
		version:     version,
	}

	return svc
}

// RegisterManager registers a new ManagerService and attempts to establish a
// connection.
//
// Only a single feeds manager is currently supported.
func (s *service) RegisterManager(mgr *FeedsManager) (int64, error) {
	count, err := s.CountManagers()
	if err != nil {
		return 0, err
	}
	if count >= 1 {
		return 0, ErrSingleFeedsManager
	}

	if mgr.IsOCRBootstrapPeer && len(mgr.JobTypes) > 0 {
		return 0, ErrBootstrapXorJobs
	}

	id, err := s.orm.CreateManager(mgr)
	if err != nil {
		return 0, err
	}

	privkey, err := s.getCSAPrivateKey()
	if err != nil {
		return 0, err
	}

	// Establish a connection
	mgr.ID = id
	s.connectFeedManager(*mgr, privkey)

	return id, nil
}

// SyncNodeInfo syncs the node's information with FMS
func (s *service) SyncNodeInfo(id int64) error {
	mgr, err := s.GetManager(id)
	if err != nil {
		return err
	}

	jobtypes := []pb.JobType{}
	for _, jt := range mgr.JobTypes {
		switch jt {
		case JobTypeFluxMonitor:
			jobtypes = append(jobtypes, pb.JobType_JOB_TYPE_FLUX_MONITOR)
		case JobTypeOffchainReporting:
			jobtypes = append(jobtypes, pb.JobType_JOB_TYPE_OCR)
		default:
			// NOOP
		}
	}

	// Assemble EVM keys
	evmKeys, err := s.ethKeyStore.SendingKeys()
	if err != nil {
		return err
	}

	evmKeyStates, err := s.ethKeyStore.GetStatesForKeys(evmKeys)
	if err != nil {
		return err
	}

	accounts := make([]*pb.Account, 0, len(evmKeyStates))
	for _, k := range evmKeyStates {
		accounts = append(accounts, &pb.Account{
			ChainType: pb.ChainType_CHAIN_TYPE_EVM,
			ChainId:   k.EVMChainID.String(),
			Address:   k.Address.String(),
		})
	}

	// Make the remote call to FMS
	fmsClient, err := s.connMgr.GetClient(id)
	if err != nil {
		return errors.Wrap(err, "could not fetch client")
	}

	chainIDs := []int64{}
	for _, c := range s.chainSet.Chains() {
		chainIDs = append(chainIDs, c.ID().Int64())
	}

	_, err = fmsClient.UpdateNode(context.Background(), &pb.UpdateNodeRequest{
		JobTypes:           jobtypes,
		ChainIds:           chainIDs,
		IsBootstrapPeer:    mgr.IsOCRBootstrapPeer,
		BootstrapMultiaddr: mgr.OCRBootstrapPeerMultiaddr.ValueOrZero(),
		Version:            s.version,
		Accounts:           accounts,
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdateManager updates the feed manager details, takes down the
// connection and reestablishes a new connection with the updated public key.
func (s *service) UpdateManager(ctx context.Context, mgr FeedsManager) error {
	err := s.orm.UpdateManager(mgr, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "could not update manager")
	}

	s.lggr.Infof("Restarting connection")

	if err = s.connMgr.Disconnect(mgr.ID); err != nil {
		s.lggr.Info("Feeds Manager not connected, attempting to connect")
	}

	// Establish a new connection
	privkey, err := s.getCSAPrivateKey()
	if err != nil {
		return err
	}

	s.connectFeedManager(mgr, privkey)

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

// CountManagerServices gets the total number of manager services
func (s *service) CountManagers() (int64, error) {
	return s.orm.CountManagers()
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
	// Track the given job proposal request
	promJobProposalRequest.Inc()

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

	// TODO - Use parent context
	var id int64
	q := s.q.WithOpts(pg.WithParentCtx(context.Background()))
	err = q.Transaction(func(tx pg.Queryer) error {
		var txerr error
		// Upsert job proposal
		id, txerr = s.orm.UpsertJobProposal(&JobProposal{
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

	return id, nil
}

// GetJobProposal gets a job proposal by id.
func (s *service) GetJobProposal(id int64) (*JobProposal, error) {
	return s.orm.GetJobProposal(id)
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

	return nil
}

// IsJobManaged determines is a job is managed by the Feeds Manager.
func (s *service) IsJobManaged(ctx context.Context, jobID int64) (bool, error) {
	return s.orm.IsJobManaged(jobID, pg.WithParentCtx(ctx))
}

// ApproveSpec approves a spec for a job proposal and creates a job with the
// spec.
func (s *service) ApproveSpec(ctx context.Context, id int64) error {
	pctx := pg.WithParentCtx(ctx)

	spec, err := s.orm.GetSpec(id, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal spec")
	}

	if spec.Status != SpecStatusPending && spec.Status != SpecStatusCancelled {
		return errors.New("must be a pending or cancelled job proposal")
	}

	proposal, err := s.orm.GetJobProposal(spec.JobProposalID, pctx)
	if err != nil {
		return errors.Wrap(err, "orm: job proposal")
	}

	fmsClient, err := s.connMgr.GetClient(proposal.FeedsManagerID)
	if err != nil {
		return errors.Wrap(err, "fms rpc client")
	}

	j, err := s.generateJob(spec.Definition)
	if err != nil {
		return errors.Wrap(err, "could not generate job from spec")
	}

	q := s.q.WithOpts(pctx)
	err = q.Transaction(func(tx pg.Queryer) error {
		// Create the job
		if err = s.jobSpawner.CreateJob(j, pg.WithQueryer(tx)); err != nil {
			return err
		}

		// Approve the job proposal spec
		if err = s.orm.ApproveSpec(id, j.ExternalJobID, pg.WithQueryer(tx)); err != nil {
			return err
		}

		// Send to FMS Client
		if _, err = fmsClient.ApprovedJob(ctx, &pb.ApprovedJobRequest{
			Uuid:    proposal.RemoteUUID.String(),
			Version: int64(spec.Version),
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "could not approve job proposal")
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
func (s *service) Start() error {
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
			return errors.New("no feeds managers registered")
		}

		mgr := mgrs[0]
		s.connectFeedManager(mgr, privkey)

		return nil
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
func (s *service) connectFeedManager(mgr FeedsManager, privkey []byte) {
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
			err := s.SyncNodeInfo(mgr.ID)
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

// Unsafe_SetConnectionsManager sets the ConnectionsManager on the service.
//
// We need to be able to inject a mock for the client to facilitate integration
// tests.
//
// ONLY TO BE USED FOR TESTING.
func (s *service) Unsafe_SetConnectionsManager(connMgr ConnectionsManager) {
	s.connMgr = connMgr
}

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
		js, err = offchainreporting.ValidatedOracleSpecToml(s.chainSet, spec)
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

func (s *service) validateProposeJobArgs(args ProposeJobArgs) error {
	// Validate the job spec
	j, err := s.generateJob(args.Spec)
	if err != nil {
		return errors.Wrap(err, "failed to generate a job based on spec")
	}

	// Validate bootstrap multiaddrs which are only allowed for OCR jobs
	if len(args.Multiaddrs) > 0 && j.Type != job.OffchainReporting {
		return errors.New("only OCR job type supports multiaddr")
	}

	return nil
}

var _ Service = &NullService{}

// NullService defines an implementation of the Feeds Service that is used
// when the Feeds Service is disabled.
type NullService struct{}

//revive:disable
func (ns NullService) Start() error { return nil }
func (ns NullService) Close() error { return nil }
func (ns NullService) ApproveSpec(ctx context.Context, id int64) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) ApproveJobProposal(ctx context.Context, id int64) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) CountManagers() (int64, error) { return 0, nil }
func (ns NullService) CancelSpec(ctx context.Context, id int64) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) CreateJobProposal(jp *JobProposal) (int64, error) {
	return 0, errors.New("feeds manager is disabled")
}
func (ns NullService) GetJobProposal(id int64) (*JobProposal, error) {
	return nil, errors.New("feeds manager is disabled")
}
func (ns NullService) ListSpecsByJobProposalIDs(ids []int64) ([]JobProposalSpec, error) {
	return nil, errors.New("feeds manager is disabled")
}
func (ns NullService) GetManager(id int64) (*FeedsManager, error) {
	return nil, errors.New("feeds manager is disabled")
}
func (ns NullService) ListManagersByIDs(ids []int64) ([]FeedsManager, error) {
	return nil, errors.New("feeds manager is disabled")
}
func (ns NullService) GetSpec(id int64) (*JobProposalSpec, error) {
	return nil, errors.New("feeds manager is disabled")
}
func (ns NullService) ListManagers() ([]FeedsManager, error)    { return nil, nil }
func (ns NullService) ListJobProposals() ([]JobProposal, error) { return nil, nil }
func (ns NullService) ListJobProposalsByManagersIDs(ids []int64) ([]JobProposal, error) {
	return nil, errors.New("feeds manager is disabled")
}
func (ns NullService) ProposeJob(ctx context.Context, args *ProposeJobArgs) (int64, error) {
	return 0, errors.New("feeds manager is disabled")
}
func (ns NullService) RegisterManager(ms *FeedsManager) (int64, error) {
	return 0, errors.New("feeds manager is disabled")
}
func (ns NullService) RejectSpec(ctx context.Context, id int64) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) SyncNodeInfo(id int64) error { return nil }
func (ns NullService) UpdateJobProposalSpec(ctx context.Context, id int64, spec string) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) UpdateManager(ctx context.Context, mgr FeedsManager) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) IsJobManaged(ctx context.Context, jobID int64) (bool, error) { return false, nil }
func (ns NullService) UpdateSpecDefinition(ctx context.Context, id int64, spec string) error {
	return errors.New("feeds manager is disabled")
}
func (ns NullService) Unsafe_SetConnectionsManager(_ ConnectionsManager) {}

//revive:enable
