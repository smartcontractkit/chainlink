package feeds

import (
	"context"
	"crypto/ed25519"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink/core/logger"
	pb "github.com/smartcontractkit/chainlink/core/services/feeds/proto"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/wsrpc"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore

type Service interface {
	Start() error
	Close() error

	CountManagers() (int64, error)
	CreateJobProposal(jp *JobProposal) (int64, error)
	GetManager(id int64) (*FeedsManager, error)
	ListManagers() ([]FeedsManager, error)
	RegisterManager(ms *FeedsManager) (int64, error)
}

type service struct {
	utils.StartStopOnce

	mu     sync.Mutex
	chDone chan struct{}
	wgDone sync.WaitGroup

	orm       ORM
	ks        keystore.CSAKeystoreInterface
	fmsClient pb.FeedsManagerClient
}

// NewService constructs a new feeds service
func NewService(orm ORM, ks keystore.CSAKeystoreInterface) Service {
	svc := &service{
		chDone: make(chan struct{}),
		orm:    orm,
		ks:     ks,
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
		return 0, errors.New("only a single feeds manager is supported")
	}

	id, err := s.orm.CreateManager(context.Background(), mgr)
	if err != nil {
		return 0, err
	}

	privkey, err := s.getCSAPrivateKey()
	if err != nil {
		return 0, err
	}

	// Establish a connection
	s.connect(mgr.URI, privkey, mgr.PublicKey, id)

	return id, nil
}

// ListManagerServices lists all the manager services.
func (s *service) ListManagers() ([]FeedsManager, error) {
	return s.orm.ListManagers(context.Background())
}

// GetManager gets a manager service by id.
func (s *service) GetManager(id int64) (*FeedsManager, error) {
	return s.orm.GetManager(context.Background(), id)
}

// CountManagerServices gets the total number of manager services
func (s *service) CountManagers() (int64, error) {
	return s.orm.CountManagers()
}

func (s *service) CreateJobProposal(jp *JobProposal) (int64, error) {
	return s.orm.CreateJobProposal(context.Background(), jp)
}

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

		s.connect(mgr.URI, privkey, mgr.PublicKey, mgr.ID)

		return nil
	})
}

func (s *service) Close() error {
	return s.StopOnce("FeedsService", func() error {
		close(s.chDone)
		s.wgDone.Wait()
		return nil
	})
}

// Connect attempts to establish a connection to the Feeds Manager.
//
// In the future we will connect to multiple Feeds Managers
func (s *service) connect(uri string, privkey []byte, pubkey []byte, feedsManagerID int64) {
	s.wgDone.Add(1)

	go func() {
		defer s.wgDone.Done()

		conn, err := wsrpc.Dial(uri,
			wsrpc.WithTransportCreds(privkey, ed25519.PublicKey(pubkey)),
		)
		if err != nil {
			logger.Infof("Error connecting to Feeds Manager server: %v", err)
			return
		}
		defer conn.Close()

		logger.Infow("[Feeds Manager] Connected to Feeds Manager", "feedsManagerID", feedsManagerID)

		// Initialize a new wsrpc client to make RPC calls
		s.mu.Lock()
		s.fmsClient = pb.NewFeedsManagerClient(conn)
		s.mu.Unlock()

		// Initialize RPC call handlers on the client connection
		pb.RegisterNodeServiceServer(conn, &RPCHandlers{
			feedsManagerID: feedsManagerID,
			svc:            s,
		})

		// Wait for close
		<-s.chDone
	}()
}

// getCSAPrivateKey gets the server's CSA private key
func (s *service) getCSAPrivateKey() (privkey []byte, err error) {
	// Fetch the server's public key
	keys, err := s.ks.ListCSAKeys()
	if err != nil {
		return privkey, err
	}
	if len(keys) < 1 {
		return privkey, errors.New("CSA key does not exist")
	}

	privkey, err = s.ks.Unsafe_GetUnlockedPrivateKey(keys[0].PublicKey)
	if err != nil {
		return []byte{}, err
	}

	return privkey, nil
}
