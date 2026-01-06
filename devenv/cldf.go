package devenv

/*
This code is an example if product uses CLD, CLDF integrations
*/

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfevmprovider "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/provider"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	csav1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/csa"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
)

const (
	AnvilKey0 = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
)

const LinkToken cldf.ContractType = "LinkToken"

type JobDistributor struct {
	nodev1.NodeServiceClient
	jobv1.JobServiceClient
	csav1.CSAServiceClient
	WSRPC string
}

type JDConfig struct {
	GRPC  string
	WSRPC string
}

// LoadCLDFEnvironment loads CLDF environment with a memory data store and JD client.
func LoadCLDFEnvironment(in *Cfg) (cldf.Environment, error) {
	ctx := context.Background()

	getCtx := func() context.Context {
		return ctx
	}

	// This only generates a brand new datastore and does not load any existing data.
	// We will need to figure out how data will be persisted and loaded in the future.
	ds := datastore.NewMemoryDataStore().Seal()

	lggr, err := logger.NewWith(func(config *zap.Config) {
		config.Development = true
		config.Encoding = "console"
	})
	if err != nil {
		return cldf.Environment{}, fmt.Errorf("failed to create logger: %w", err)
	}

	blockchains, err := loadCLDFChains(in.Blockchains)
	if err != nil {
		return cldf.Environment{}, fmt.Errorf("failed to load CLDF chains: %w", err)
	}

	jd, err := NewJDClient(ctx, JDConfig{
		GRPC:  in.JD.Out.ExternalGRPCUrl,
		WSRPC: in.JD.Out.ExternalWSRPCUrl,
	})
	if err != nil {
		return cldf.Environment{},
			fmt.Errorf("failed to load offchain client: %w", err)
	}

	opBundle := operations.NewBundle(
		getCtx,
		lggr,
		operations.NewMemoryReporter(),
		operations.WithOperationRegistry(operations.NewOperationRegistry()),
	)

	return cldf.Environment{
		Name:              "local",
		Logger:            lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         ds,
		Offchain:          jd,
		GetContext:        getCtx,
		OperationsBundle:  opBundle,
		BlockChains:       cldfchain.NewBlockChainsFromSlice(blockchains),
	}, nil
}

func loadCLDFChains(bcis []*blockchain.Input) ([]cldfchain.BlockChain, error) {
	blockchains := make([]cldfchain.BlockChain, 0)
	for _, bci := range bcis {
		switch bci.Type {
		case "anvil":
			bc, err := loadEVMChain(bci)
			if err != nil {
				return blockchains, fmt.Errorf("failed to load EVM chain %s: %w", bci.ChainID, err)
			}

			blockchains = append(blockchains, bc)
		default:
			return blockchains, fmt.Errorf("unsupported chain type %s", bci.Type)
		}
	}

	return blockchains, nil
}

func loadEVMChain(bci *blockchain.Input) (cldfchain.BlockChain, error) {
	if bci.Out == nil {
		return nil, fmt.Errorf("output configuration for %s blockchain %s is not set", bci.Type, bci.ChainID)
	}

	chainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(bci.ChainID, chainsel.FamilyEVM)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain details for %s: %w", bci.ChainID, err)
	}

	chain, err := cldfevmprovider.NewRPCChainProvider(
		chainDetails.ChainSelector,
		cldfevmprovider.RPCChainProviderConfig{
			DeployerTransactorGen: cldfevmprovider.TransactorFromRaw(
				AnvilKey0,
			),
			RPCs: []cldf.RPC{
				{
					Name:               "default",
					WSURL:              bci.Out.Nodes[0].ExternalWSUrl,
					HTTPURL:            bci.Out.Nodes[0].ExternalHTTPUrl,
					PreferredURLScheme: cldf.URLSchemePreferenceHTTP,
				},
			},
			ConfirmFunctor: cldfevmprovider.ConfirmFuncGeth(1 * time.Minute),
		},
	).Initialize(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize EVM chain %s: %w", bci.ChainID, err)
	}

	return chain, nil
}

// NewJDClient creates a new JobDistributor client.
func NewJDClient(ctx context.Context, cfg JDConfig) (cldf.OffchainClient, error) {
	conn, err := NewJDConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}
	jd := &JobDistributor{
		WSRPC:             cfg.WSRPC,
		NodeServiceClient: nodev1.NewNodeServiceClient(conn),
		JobServiceClient:  jobv1.NewJobServiceClient(conn),
		CSAServiceClient:  csav1.NewCSAServiceClient(conn),
	}

	return jd, err
}

func (jd JobDistributor) GetCSAPublicKey(ctx context.Context) (string, error) {
	keypairs, err := jd.ListKeypairs(ctx, &csav1.ListKeypairsRequest{})
	if err != nil {
		return "", err
	}
	if keypairs == nil || len(keypairs.Keypairs) == 0 {
		return "", errors.New("no keypairs found")
	}
	csakey := keypairs.Keypairs[0].PublicKey
	return csakey, nil
}

// ProposeJob proposes jobs through the jobService and accepts the proposed job on selected node based on ProposeJobRequest.NodeId.
func (jd JobDistributor) ProposeJob(ctx context.Context, in *jobv1.ProposeJobRequest, opts ...grpc.CallOption) (*jobv1.ProposeJobResponse, error) {
	res, err := jd.JobServiceClient.ProposeJob(ctx, in, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to propose job. err: %w", err)
	}
	if res.Proposal == nil {
		return nil, errors.New("failed to propose job. err: proposal is nil")
	}

	return res, nil
}

// NewJDConnection creates new gRPC connection with JobDistributor.
func NewJDConnection(cfg JDConfig) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	interceptors := []grpc.UnaryClientInterceptor{}

	if len(interceptors) > 0 {
		opts = append(opts, grpc.WithChainUnaryInterceptor(interceptors...))
	}

	conn, err := grpc.NewClient(cfg.GRPC, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect Job Distributor service. Err: %w", err)
	}

	return conn, nil
}
