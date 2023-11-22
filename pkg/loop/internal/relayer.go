package internal

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ PluginRelayer = (*PluginRelayerClient)(nil)

type PluginRelayerClient struct {
	*pluginClient

	grpc pb.PluginRelayerClient
}

func NewPluginRelayerClient(broker Broker, brokerCfg BrokerConfig, conn *grpc.ClientConn) *PluginRelayerClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "PluginRelayerClient")
	pc := newPluginClient(broker, brokerCfg, conn)
	return &PluginRelayerClient{pluginClient: pc, grpc: pb.NewPluginRelayerClient(pc)}
}

func (p *PluginRelayerClient) NewRelayer(ctx context.Context, config string, keystore types.Keystore) (Relayer, error) {
	cc := p.newClientConn("Relayer", func(ctx context.Context) (id uint32, deps resources, err error) {
		var ksRes resource
		id, ksRes, err = p.serveNew("Keystore", func(s *grpc.Server) {
			pb.RegisterKeystoreServer(s, &keystoreServer{impl: keystore})
		})
		if err != nil {
			return 0, nil, fmt.Errorf("Failed to create relayer client: failed to serve keystore: %w", err)
		}
		deps.Add(ksRes)

		reply, err := p.grpc.NewRelayer(ctx, &pb.NewRelayerRequest{
			Config:     config,
			KeystoreID: id,
		})
		if err != nil {
			return 0, nil, fmt.Errorf("Failed to create relayer client: failed request: %w", err)
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

func RegisterPluginRelayerServer(server *grpc.Server, broker Broker, brokerCfg BrokerConfig, impl PluginRelayer) error {
	pb.RegisterPluginRelayerServer(server, newPluginRelayerServer(broker, brokerCfg, impl))
	return nil
}

func newPluginRelayerServer(broker Broker, brokerCfg BrokerConfig, impl PluginRelayer) *pluginRelayerServer {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "RelayerPluginServer")
	return &pluginRelayerServer{brokerExt: &brokerExt{broker, brokerCfg}, impl: impl}
}

func (p *pluginRelayerServer) NewRelayer(ctx context.Context, request *pb.NewRelayerRequest) (*pb.NewRelayerReply, error) {
	ksConn, err := p.dial(request.KeystoreID)
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
	id, _, err := p.serveNew(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: r})
		pb.RegisterRelayerServer(s, newChainRelayerServer(r, p.brokerExt))
	}, rRes, ksRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewRelayerReply{RelayerID: id}, nil
}

var _ types.Keystore = (*keystoreClient)(nil)

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

	impl types.Keystore
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

var _ Relayer = (*relayerClient)(nil)

// relayerClient adapts a GRPC [pb.RelayerClient] to implement [Relayer].
type relayerClient struct {
	*brokerExt
	*serviceClient

	relayer pb.RelayerClient
}

func newRelayerClient(b *brokerExt, conn grpc.ClientConnInterface) *relayerClient {
	b = b.withName("ChainRelayerClient")
	return &relayerClient{b, newServiceClient(b, conn), pb.NewRelayerClient(conn)}
}

func (r *relayerClient) NewConfigProvider(ctx context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	cc := r.newClientConn("ConfigProvider", func(ctx context.Context) (uint32, resources, error) {
		reply, err := r.relayer.NewConfigProvider(ctx, &pb.NewConfigProviderRequest{
			RelayArgs: &pb.RelayArgs{
				ExternalJobID: rargs.ExternalJobID[:],
				JobID:         rargs.JobID,
				ContractID:    rargs.ContractID,
				New:           rargs.New,
				RelayConfig:   rargs.RelayConfig,
			},
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.ConfigProviderID, nil, nil
	})
	return newConfigProviderClient(r.withName("ConfigProviderClient"), cc), nil
}

func (r *relayerClient) NewPluginProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.PluginProvider, error) {
	cc := r.newClientConn("PluginProvider", func(ctx context.Context) (uint32, resources, error) {
		reply, err := r.relayer.NewPluginProvider(ctx, &pb.NewPluginProviderRequest{
			RelayArgs: &pb.RelayArgs{
				ExternalJobID: rargs.ExternalJobID[:],
				JobID:         rargs.JobID,
				ContractID:    rargs.ContractID,
				New:           rargs.New,
				RelayConfig:   rargs.RelayConfig,
				ProviderType:  rargs.ProviderType,
			},
			PluginArgs: &pb.PluginArgs{
				TransmitterID: pargs.TransmitterID,
				PluginConfig:  pargs.PluginConfig,
			},
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.PluginProviderID, nil, nil
	})

	// TODO: Remove this when we have fully transitioned all relayers to running in LOOPPs.
	// This allows callers to type assert a PluginProvider into a product provider type (eg. MedianProvider)
	// for interoperability with legacy code.
	switch rargs.ProviderType {
	case string(types.Median):
		return newMedianProviderClient(r.brokerExt, cc), nil
	case string(types.GenericPlugin):
		return newPluginProviderClient(r.brokerExt, cc), nil
	default:
		return nil, fmt.Errorf("provider type not supported: %s", rargs.ProviderType)
	}
}

func (r *relayerClient) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	reply, err := r.relayer.GetChainStatus(ctx, &pb.GetChainStatusRequest{})
	if err != nil {
		return types.ChainStatus{}, err
	}

	return types.ChainStatus{
		ID:      reply.Chain.Id,
		Enabled: reply.Chain.Enabled,
		Config:  reply.Chain.Config,
	}, nil
}

func (r *relayerClient) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (nodes []types.NodeStatus, nextPageToken string, total int, err error) {
	reply, err := r.relayer.ListNodeStatuses(ctx, &pb.ListNodeStatusesRequest{
		PageSize:  pageSize,
		PageToken: pageToken,
	})
	if err != nil {
		return nil, "", -1, err
	}
	for _, n := range reply.Nodes {
		nodes = append(nodes, types.NodeStatus{
			ChainID: n.ChainID,
			Name:    n.Name,
			Config:  n.Config,
			State:   n.State,
		})
	}
	total = int(reply.Total)
	return
}

func (r *relayerClient) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	_, err := r.relayer.Transact(ctx, &pb.TransactionRequest{
		From:         from,
		To:           to,
		Amount:       pb.NewBigIntFromInt(amount),
		BalanceCheck: balanceCheck,
	})
	return err
}

var _ pb.RelayerServer = (*relayerServer)(nil)

// relayerServer exposes [Relayer] as a GRPC [pb.RelayerServer].
type relayerServer struct {
	pb.UnimplementedRelayerServer

	*brokerExt

	impl Relayer
}

func newChainRelayerServer(impl Relayer, b *brokerExt) *relayerServer {
	return &relayerServer{impl: impl, brokerExt: b.withName("ChainRelayerServer")}
}

func (r *relayerServer) NewConfigProvider(ctx context.Context, request *pb.NewConfigProviderRequest) (*pb.NewConfigProviderReply, error) {
	exJobID, err := uuid.FromBytes(request.RelayArgs.ExternalJobID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid bytes for ExternalJobID: %w", err)
	}
	cp, err := r.impl.NewConfigProvider(ctx, types.RelayArgs{
		ExternalJobID: exJobID,
		JobID:         request.RelayArgs.JobID,
		ContractID:    request.RelayArgs.ContractID,
		New:           request.RelayArgs.New,
		RelayConfig:   request.RelayArgs.RelayConfig,
	})
	if err != nil {
		return nil, err
	}
	err = cp.Start(ctx)
	if err != nil {
		return nil, err
	}

	const name = "ConfigProvider"
	id, _, err := r.serveNew(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: cp})
		pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: cp.OffchainConfigDigester()})
		pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: cp.ContractConfigTracker()})
	}, resource{cp, name})
	if err != nil {
		return nil, err
	}

	return &pb.NewConfigProviderReply{ConfigProviderID: id}, nil
}

func (r *relayerServer) NewPluginProvider(ctx context.Context, request *pb.NewPluginProviderRequest) (*pb.NewPluginProviderReply, error) {
	exJobID, err := uuid.FromBytes(request.RelayArgs.ExternalJobID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid bytes for ExternalJobID: %w", err)
	}
	relayArgs := types.RelayArgs{
		ExternalJobID: exJobID,
		JobID:         request.RelayArgs.JobID,
		ContractID:    request.RelayArgs.ContractID,
		New:           request.RelayArgs.New,
		RelayConfig:   request.RelayArgs.RelayConfig,
		ProviderType:  request.RelayArgs.ProviderType,
	}
	pluginArgs := types.PluginArgs{
		TransmitterID: request.PluginArgs.TransmitterID,
		PluginConfig:  request.PluginArgs.PluginConfig,
	}

	switch request.RelayArgs.ProviderType {
	case string(types.Median):
		id, err := r.newMedianProvider(ctx, relayArgs, pluginArgs)
		if err != nil {
			return nil, err
		}
		return &pb.NewPluginProviderReply{PluginProviderID: id}, nil
	case string(types.GenericPlugin):
		id, err := r.newPluginProvider(ctx, relayArgs, pluginArgs)
		if err != nil {
			return nil, err
		}
		return &pb.NewPluginProviderReply{PluginProviderID: id}, nil
	}

	return nil, fmt.Errorf("provider type not supported: %s", relayArgs.ProviderType)
}

func (r *relayerServer) newMedianProvider(ctx context.Context, relayArgs types.RelayArgs, pluginArgs types.PluginArgs) (uint32, error) {
	i, ok := r.impl.(MedianProvider)
	if !ok {
		return 0, errors.New("median not supported")
	}

	provider, err := i.NewMedianProvider(ctx, relayArgs, pluginArgs)
	if err != nil {
		return 0, err
	}
	err = provider.Start(ctx)
	if err != nil {
		return 0, err
	}
	const name = "MedianProvider"
	providerRes := resource{name: name, Closer: provider}

	id, _, err := r.serveNew(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: provider})
		pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: provider.OffchainConfigDigester()})
		pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: provider.ContractConfigTracker()})
		pb.RegisterContractTransmitterServer(s, &contractTransmitterServer{impl: provider.ContractTransmitter()})
		pb.RegisterReportCodecServer(s, &reportCodecServer{impl: provider.ReportCodec()})
		pb.RegisterMedianContractServer(s, &medianContractServer{impl: provider.MedianContract()})
		pb.RegisterOnchainConfigCodecServer(s, &onchainConfigCodecServer{impl: provider.OnchainConfigCodec()})
	}, providerRes)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *relayerServer) newPluginProvider(ctx context.Context, relayArgs types.RelayArgs, pluginArgs types.PluginArgs) (uint32, error) {
	provider, err := r.impl.NewPluginProvider(ctx, relayArgs, pluginArgs)
	if err != nil {
		return 0, err
	}
	err = provider.Start(ctx)
	if err != nil {
		return 0, err
	}
	const name = "PluginProvider"
	providerRes := resource{name: name, Closer: provider}

	id, _, err := r.serveNew(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: provider})
		pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: provider.OffchainConfigDigester()})
		pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: provider.ContractConfigTracker()})
		pb.RegisterContractTransmitterServer(s, &contractTransmitterServer{impl: provider.ContractTransmitter()})
	}, providerRes)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *relayerServer) GetChainStatus(ctx context.Context, request *pb.GetChainStatusRequest) (*pb.GetChainStatusReply, error) {
	chain, err := r.impl.GetChainStatus(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.GetChainStatusReply{Chain: &pb.ChainStatus{
		Id:      chain.ID,
		Enabled: chain.Enabled,
		Config:  chain.Config,
	}}, nil
}

func (r *relayerServer) ListNodeStatuses(ctx context.Context, request *pb.ListNodeStatusesRequest) (*pb.ListNodeStatusesReply, error) {
	nodeConfigs, nextPageToken, total, err := r.impl.ListNodeStatuses(ctx, request.PageSize, request.PageToken)
	if err != nil {
		return nil, err
	}
	var nodes []*pb.NodeStatus
	for _, n := range nodeConfigs {
		nodes = append(nodes, &pb.NodeStatus{
			ChainID: n.ChainID,
			Name:    n.Name,
			Config:  n.Config,
			State:   n.State,
		})
	}
	return &pb.ListNodeStatusesReply{Nodes: nodes, NextPageToken: nextPageToken, Total: int32(total)}, nil
}
func (r *relayerServer) Transact(ctx context.Context, request *pb.TransactionRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, r.impl.Transact(ctx, request.From, request.To, request.Amount.Int(), request.BalanceCheck)
}

func healthReport(s map[string]string) (hr map[string]error) {
	hr = make(map[string]error, len(s))
	for n, e := range s {
		var err error
		if e != "" {
			err = errors.New(e)
		}
		hr[n] = err
	}
	return hr
}

// RegisterStandAloneMedianProvider register the servers needed for a median plugin provider,
// this is a workaround to test the Node API on EVM until the EVM relayer is loopifyed
func RegisterStandAloneMedianProvider(s *grpc.Server, p types.MedianProvider) {
	pb.RegisterServiceServer(s, &serviceServer{srv: p})
	pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: p.OffchainConfigDigester()})
	pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: p.ContractConfigTracker()})
	pb.RegisterContractTransmitterServer(s, &contractTransmitterServer{impl: p.ContractTransmitter()})
	pb.RegisterReportCodecServer(s, &reportCodecServer{impl: p.ReportCodec()})
	pb.RegisterMedianContractServer(s, &medianContractServer{impl: p.MedianContract()})
	pb.RegisterOnchainConfigCodecServer(s, &onchainConfigCodecServer{impl: p.OnchainConfigCodec()})
}
