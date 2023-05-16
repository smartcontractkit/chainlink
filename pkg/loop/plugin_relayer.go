package loop

import (
	"context"
	"math/big"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

// PluginRelayerName is the name for [PluginRelayer]/[NewGRPCPluginRelayer].
const PluginRelayerName = "relayer"

type PluginRelayer interface {
	NewRelayer(ctx context.Context, config string, keystore Keystore) (Relayer, error)
}

func PluginRelayerHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_RELAYER_MAGIC_COOKIE",
		MagicCookieValue: "dae753d4542311b33cf041b930db0150647e806175c2818a0c88a9ab745e45aa",
	}
}

type Keystore interface {
	Accounts(ctx context.Context) (accounts []string, err error)
	// Sign returns data signed by account.
	// nil data can be used as a no-op to check for account existence.
	Sign(ctx context.Context, account string, data []byte) (signed []byte, err error)
}

type Relayer interface {
	types.Service

	NewConfigProvider(context.Context, types.RelayArgs) (types.ConfigProvider, error)
	NewMedianProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.MedianProvider, error)
	NewMercuryProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.MercuryProvider, error)

	ChainStatus(ctx context.Context, id string) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) (chains []types.ChainStatus, count int, err error)

	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)

	SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error
}

var _ plugin.GRPCPlugin = (*GRPCPluginRelayer)(nil)

// GRPCPluginRelayer implements [plugin.GRPCPlugin] for [PluginRelayer].
type GRPCPluginRelayer struct {
	plugin.NetRPCUnsupportedPlugin

	StopCh <-chan struct{}
	Logger logger.Logger

	PluginServer PluginRelayer

	pluginClient *pluginRelayerClient
}

func (p *GRPCPluginRelayer) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	pb.RegisterPluginRelayerServer(server, newPluginRelayerServer(p.StopCh, p.Logger, broker, p.PluginServer))
	return nil
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the PluginServer [PluginRelayer], updated with the new broker and conn.
func (p *GRPCPluginRelayer) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = newPluginRelayerClient(p.StopCh, p.Logger, broker, conn)
	} else {
		p.pluginClient.refresh(broker, conn)
	}
	return p.pluginClient, nil
}

func (p *GRPCPluginRelayer) ClientConfig() *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig:  PluginRelayerHandshakeConfig(),
		Plugins:          map[string]plugin.Plugin{PluginRelayerName: p},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}
}

var _ PluginRelayer = (*pluginRelayerClient)(nil)

type pluginRelayerClient struct {
	*pluginClient

	grpc pb.PluginRelayerClient
}

func newPluginRelayerClient(stopCh <-chan struct{}, lggr logger.Logger, broker *plugin.GRPCBroker, conn *grpc.ClientConn) *pluginRelayerClient {
	lggr = logger.Named(lggr, "PluginRelayerClient")
	pc := newPluginClient(stopCh, lggr, broker, conn)
	return &pluginRelayerClient{pluginClient: pc, grpc: pb.NewPluginRelayerClient(pc)}
}

func (p *pluginRelayerClient) NewRelayer(ctx context.Context, config string, keystore Keystore) (Relayer, error) {
	cc := p.newClientConn("Relayer", func(ctx context.Context) (id uint32, deps resources, err error) {
		var ksRes resource
		id, ksRes, err = p.serve("Keystore", func(s *grpc.Server) {
			pb.RegisterKeystoreServer(s, &keystoreServer{impl: keystore})
		})
		if err != nil {
			return
		}
		deps.Add(ksRes)

		reply, err := p.grpc.NewRelayer(ctx, &pb.NewRelayerRequest{
			Config:     config,
			KeystoreID: id,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.RelayerID, nil, nil
	})
	return newRelayerClient(p.brokerExt, cc), nil
}

type pluginRelayerServer struct {
	pb.UnimplementedPluginRelayerServer

	*brokerExt

	impl PluginRelayer
}

func newPluginRelayerServer(stopCh <-chan struct{}, lggr logger.Logger, broker *plugin.GRPCBroker, impl PluginRelayer) *pluginRelayerServer {
	lggr = logger.Named(lggr, "RelayerPluginServer")
	return &pluginRelayerServer{brokerExt: &brokerExt{stopCh, lggr, broker}, impl: impl}
}

func (p *pluginRelayerServer) NewRelayer(ctx context.Context, request *pb.NewRelayerRequest) (*pb.NewRelayerReply, error) {
	ksConn, err := p.broker.Dial(request.KeystoreID)
	if err != nil {
		return nil, ErrConnDial{Name: "Keystore", ID: request.KeystoreID, Err: err}
	}
	ksRes := resource{ksConn, "Keystore"}
	r, err := p.impl.NewRelayer(ctx, request.Config, newKeystoreClient(ksConn))
	if err != nil {
		p.closeAll(ksRes)
		return nil, err
	}
	err = r.Start(ctx)
	if err != nil {
		p.closeAll(ksRes)
		return nil, err
	}

	const name = "Relayer"
	rRes := resource{r, name}
	id, _, err := p.serve(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: r})
		pb.RegisterRelayerServer(s, newChainRelayerServer(r, p.brokerExt))
	}, rRes, ksRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewRelayerReply{RelayerID: id}, nil
}

var _ Keystore = (*keystoreClient)(nil)

type keystoreClient struct {
	grpc pb.KeystoreClient
}

func newKeystoreClient(cc *grpc.ClientConn) *keystoreClient {
	return &keystoreClient{pb.NewKeystoreClient(cc)}
}

func (k *keystoreClient) Accounts(ctx context.Context) (accounts []string, err error) {
	reply, err := k.grpc.Accounts(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return reply.Accounts, nil
}

func (k *keystoreClient) Sign(ctx context.Context, account string, data []byte) ([]byte, error) {
	reply, err := k.grpc.Sign(ctx, &pb.SignRequest{Account: account, Data: data})
	if err != nil {
		return nil, err
	}
	return reply.SignedData, nil
}

var _ pb.KeystoreServer = (*keystoreServer)(nil)

type keystoreServer struct {
	pb.UnimplementedKeystoreServer

	impl Keystore
}

func (k *keystoreServer) Accounts(ctx context.Context, _ *emptypb.Empty) (*pb.AccountsReply, error) {
	as, err := k.impl.Accounts(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.AccountsReply{Accounts: as}, nil
}

func (k *keystoreServer) Sign(ctx context.Context, request *pb.SignRequest) (*pb.SignReply, error) {
	signed, err := k.impl.Sign(ctx, request.Account, request.Data)
	if err != nil {
		return nil, err
	}
	return &pb.SignReply{SignedData: signed}, nil
}
