package loop

import (
	"context"
	"math/big"
	"time"

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
		ProtocolVersion:  0,
		MagicCookieKey:   "magic-key-relayer-todo",
		MagicCookieValue: "magic-value-relayer-todo",
	}
}

func PluginRelayerClientConfig(lggr logger.Logger) *plugin.ClientConfig {
	return &plugin.ClientConfig{
		HandshakeConfig: PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			PluginRelayerName: NewGRPCPluginRelayer(nil, lggr),
		},
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	}
}

type Keystore interface {
	Keys(ctx context.Context) (accounts []string, err error)
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

	//TODO return info? s/SendTokens? remove https://smartcontract-it.atlassian.net/browse/BCF-2111
	SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error
}

var _ plugin.GRPCPlugin = (*grpcPluginRelayer)(nil)

// grpcPluginRelayer implements [plugin.GRPCPlugin] for [PluginRelayer].
type grpcPluginRelayer struct {
	plugin.NetRPCUnsupportedPlugin

	impl PluginRelayer
	lggr logger.Logger
}

// NewGRPCPluginRelayer returns a new [plugin.Plugin] (which really only implements [plugin.GRPCPlugin]).
// [PluginRelayer] is required for servers.
func NewGRPCPluginRelayer(cp PluginRelayer, lggr logger.Logger) plugin.Plugin {
	return &grpcPluginRelayer{impl: cp, lggr: lggr}
}

func (p *grpcPluginRelayer) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	pb.RegisterPluginRelayerServer(server, newRelayerPluginServer(p.lggr, broker, p.impl))
	return nil
}

// GRPCClient implements [plugin.GRPCPlugin] and returns a [PluginRelayer].
func (p *grpcPluginRelayer) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return newRelayerPluginClient(p.lggr, broker, conn), nil
}

var _ PluginRelayer = (*relayerPluginClient)(nil)

type relayerPluginClient struct {
	*lggrBroker

	grpc pb.PluginRelayerClient
}

func newRelayerPluginClient(lggr logger.Logger, broker *plugin.GRPCBroker, conn *grpc.ClientConn) *relayerPluginClient {
	lggr = logger.Named(lggr, "RelayerPluginClient")
	return &relayerPluginClient{lggrBroker: &lggrBroker{lggr, broker}, grpc: pb.NewPluginRelayerClient(conn)}
}

func (p *relayerPluginClient) NewRelayer(ctx context.Context, config string, keystore Keystore) (Relayer, error) {
	ksSrv := grpc.NewServer()
	pb.RegisterKeystoreServer(ksSrv, &keystoreServer{impl: keystore})
	id, err := p.serve(ksSrv, "Keystore")
	if err != nil {
		return nil, err
	}
	reply, err := p.grpc.NewRelayer(ctx, &pb.NewRelayerRequest{
		Config:     config,
		KeystoreID: id,
	})
	if err != nil {
		ksSrv.Stop()
		return nil, err
	}
	conn, err := p.broker.Dial(reply.RelayerID)
	if err != nil {
		ksSrv.Stop()
		return nil, ErrConnDial{Name: "Relayer", ID: reply.RelayerID, Err: err}
	}
	return newRelayerClient(p.lggrBroker, conn, ksSrv), nil
}

type relayerPluginServer struct {
	pb.UnimplementedPluginRelayerServer

	*lggrBroker

	impl PluginRelayer
}

func newRelayerPluginServer(lggr logger.Logger, broker *plugin.GRPCBroker, impl PluginRelayer) *relayerPluginServer {
	lggr = logger.Named(lggr, "RelayerPluginServer")
	return &relayerPluginServer{lggrBroker: &lggrBroker{lggr, broker}, impl: impl}
}

func (p *relayerPluginServer) NewRelayer(ctx context.Context, request *pb.NewRelayerRequest) (*pb.NewRelayerReply, error) {
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
	relayerSrv := grpc.NewServer()
	pb.RegisterServiceServer(relayerSrv, &serviceServer{srv: r, stop: func() {
		time.AfterFunc(time.Second, relayerSrv.GracefulStop)
	}})
	pb.RegisterRelayerServer(relayerSrv, newChainRelayerServer(r, p.lggrBroker))
	const name = "Relayer"
	id, err := p.serve(relayerSrv, name, resource{r, name}, ksRes)
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

func (k *keystoreClient) Keys(ctx context.Context) (accounts []string, err error) {
	reply, err := k.grpc.Keys(ctx, &emptypb.Empty{})
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

func (k *keystoreServer) Keys(ctx context.Context, _ *emptypb.Empty) (*pb.KeysReply, error) {
	as, err := k.impl.Keys(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.KeysReply{Accounts: as}, nil
}

func (k *keystoreServer) Sign(ctx context.Context, request *pb.SignRequest) (*pb.SignReply, error) {
	signed, err := k.impl.Sign(ctx, request.Account, request.Data)
	if err != nil {
		return nil, err
	}
	return &pb.SignReply{SignedData: signed}, nil
}
