package ccip

import (
	"context"
	"fmt"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/reportingplugin/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	ccipprovider "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// CommitLOOPClient is a client is run on the core node to connect to the commit LOOP server.
type CommitLOOPClient struct {
	// hashicorp plugin client
	*goplugin.PluginClient
	// client to base service
	*goplugin.ServiceClient

	// creates new commit factory instances
	generator ccippb.CommitFactoryGeneratorClient
}

func NewCommitLOOPClient(broker net.Broker, brokerCfg net.BrokerConfig, conn *grpc.ClientConn) *CommitLOOPClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "CommitLOOPClient")
	pc := goplugin.NewPluginClient(broker, brokerCfg, conn)
	return &CommitLOOPClient{
		PluginClient:  pc,
		ServiceClient: goplugin.NewServiceClient(pc.BrokerExt, pc),
		generator:     ccippb.NewCommitFactoryGeneratorClient(pc),
	}
}

// NewCommitFactory creates a new reporting plugin factory client.
// In practice this client is called by the core node.
// The reporting plugin factory client is a client to the LOOP server, which
// is run as an external process via hashicorp plugin. If the given provider is a GRPCClientConn, then the provider is proxied to the
// to the relayer, which is its own process via hashicorp plugin. If the provider is not a GRPCClientConn, then the provider is a local
// to the core node. The core must wrap the provider in a grpc server and serve it locally.
func (c *CommitLOOPClient) NewCommitFactory(ctx context.Context, provider types.CCIPCommitProvider) (types.ReportingPluginFactory, error) {
	newCommitClientFn := func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		// TODO are there any local resources that need to be passed to the executor and started as a server?

		// the proxyable resources are the Provider,  which may or may not be local to the client process. (legacy vs loopp)
		var (
			providerID       uint32
			providerResource net.Resource
		)
		if grpcProvider, ok := provider.(goplugin.GRPCClientConn); ok {
			// TODO: BCF-3061 ccip provider can create new services. the proxying needs to be augmented
			// to intercept and route to the created services. also, need to prevent leaks.
			providerID, providerResource, err = c.Serve("CommitProvider", proxy.NewProxy(grpcProvider.ClientConn()))
		} else {
			// loop client runs in the core node. if the provider is not a grpc client conn, then we are in legacy mode
			// and need to serve all the required services locally.
			providerID, providerResource, err = c.ServeNew("CommitProvider", func(s *grpc.Server) {
				ccipprovider.RegisterCommitProviderServices(s, provider, c.BrokerExt)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerResource)

		resp, err := c.generator.NewCommitFactory(ctx, &ccippb.NewCommitFactoryRequest{
			ProviderServiceId: providerID,
		})
		if err != nil {
			return 0, nil, err
		}
		return resp.CommitFactoryServiceId, deps, nil
	}
	cc := c.NewClientConn("CommitFactory", newCommitClientFn)
	return ocr2.NewReportingPluginFactoryClient(c.BrokerExt, cc), nil
}

// CommitLOOPServer is a server that runs the commit LOOP.
type CommitLOOPServer struct {
	ccippb.UnimplementedCommitFactoryGeneratorServer

	*net.BrokerExt
	impl types.CCIPCommitFactoryGenerator
}

func RegisterCommitLOOPServer(s *grpc.Server, b net.Broker, cfg net.BrokerConfig, impl types.CCIPCommitFactoryGenerator) error {
	ext := &net.BrokerExt{Broker: b, BrokerConfig: cfg}
	ccippb.RegisterCommitFactoryGeneratorServer(s, newCommitLOOPServer(impl, ext))
	return nil
}

func newCommitLOOPServer(impl types.CCIPCommitFactoryGenerator, b *net.BrokerExt) *CommitLOOPServer {
	return &CommitLOOPServer{impl: impl, BrokerExt: b.WithName("CommitLOOPServer")}
}

func (r *CommitLOOPServer) NewCommitFactory(ctx context.Context, request *ccippb.NewCommitFactoryRequest) (*ccippb.NewCommitFactoryResponse, error) {
	var err error
	var deps net.Resources
	defer func() {
		if err != nil {
			r.CloseAll(deps...)
		}
	}()

	// lookup the provider service
	providerConn, err := r.Dial(request.ProviderServiceId)
	if err != nil {
		return nil, net.ErrConnDial{Name: "CommitProvider", ID: request.ProviderServiceId, Err: err}
	}
	deps.Add(net.Resource{Closer: providerConn, Name: "CommitProvider"})
	provider := ccipprovider.NewCommitProviderClient(r.BrokerExt, providerConn)

	factory, err := r.impl.NewCommitFactory(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to create new commit factory: %w", err)
	}

	id, _, err := r.ServeNew("CommitFactory", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		pb.RegisterReportingPluginFactoryServer(s, ocr2.NewReportingPluginFactoryServer(factory, r.BrokerExt))
	}, deps...)
	if err != nil {
		return nil, fmt.Errorf("failed to serve new commit factory: %w", err)
	}
	return &ccippb.NewCommitFactoryResponse{CommitFactoryServiceId: id}, nil
}
