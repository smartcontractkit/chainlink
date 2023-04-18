package loop

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ Relayer = (*relayerClient)(nil)

// relayerClient adapts a GRPC [pb.RelayerClient] to implement [Relayer].
type relayerClient struct {
	*lggrBroker

	relayer pb.RelayerClient
	service pb.ServiceClient
	closeFn func() error
}

func newRelayerClient(lb *lggrBroker, conn *grpc.ClientConn, deps ...*grpc.Server) *relayerClient {
	return &relayerClient{lb.named("ChainRelayerClient"), pb.NewRelayerClient(conn), pb.NewServiceClient(conn), func() error {
		for _, d := range deps {
			time.AfterFunc(time.Second, d.GracefulStop)
		}
		return conn.Close()
	}}
}

func (r *relayerClient) Start(ctx context.Context) error {
	_, err := r.service.Start(context.TODO(), &emptypb.Empty{})
	return err
}

func (r *relayerClient) Close() error {
	_, err := r.service.Close(context.TODO(), &emptypb.Empty{})
	err = multierr.Append(err, r.closeFn())
	return err
}

func (r *relayerClient) Ready() error {
	_, err := r.service.Ready(context.TODO(), &emptypb.Empty{})
	return err
}

func (r *relayerClient) Name() string { return r.lggr.Name() }

func (r *relayerClient) HealthReport() map[string]error {
	reply, err := r.service.HealthReport(context.TODO(), &emptypb.Empty{})
	if err != nil {
		return map[string]error{r.lggr.Name(): err}
	}
	hr := healthReport(reply.HealthReport)
	hr[r.lggr.Name()] = nil
	return hr
}

func (r *relayerClient) NewConfigProvider(ctx context.Context, rargs types.RelayArgs) (types.ConfigProvider, error) {
	reply, err := r.relayer.NewConfigProvider(ctx, &pb.NewConfigProviderRequest{
		RelayArgs: &pb.RelayArgs{
			ExternalJobID: rargs.ExternalJobID.Bytes(),
			JobID:         rargs.JobID,
			ContractID:    rargs.ContractID,
			New:           rargs.New,
			RelayConfig:   rargs.RelayConfig,
		},
	})
	if err != nil {
		return nil, err
	}
	providerConn, err := r.broker.Dial(reply.ConfigProviderID)
	if err != nil {
		return nil, ErrConnDial{Name: "ConfigProvider", ID: reply.ConfigProviderID, Err: err}
	}
	return newConfigProviderClient(r.named("ConfigProviderClient"), providerConn), nil
}

func (r *relayerClient) NewMedianProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	reply, err := r.relayer.NewMedianProvider(ctx, &pb.NewMedianProviderRequest{
		RelayArgs: &pb.RelayArgs{
			ExternalJobID: rargs.ExternalJobID.Bytes(),
			JobID:         rargs.JobID,
			ContractID:    rargs.ContractID,
			New:           rargs.New,
			RelayConfig:   rargs.RelayConfig,
		},
		PluginArgs: &pb.PluginArgs{
			TransmitterID: pargs.TransmitterID,
			PluginConfig:  pargs.PluginConfig,
		},
	})
	if err != nil {
		return nil, err
	}
	id := reply.MedianProviderID
	factoryConn, err := r.broker.Dial(id)
	if err != nil {
		return nil, ErrConnDial{Name: "MedianProvider", ID: id, Err: err}
	}
	return newMedianProviderClient(r.lggrBroker, factoryConn), nil
}

func (r *relayerClient) NewMercuryProvider(context.Context, types.RelayArgs, types.PluginArgs) (types.MercuryProvider, error) {
	return nil, errors.New("mercury is not supported")
}

func (r *relayerClient) ChainStatus(ctx context.Context, id string) (types.ChainStatus, error) {
	reply, err := r.relayer.ChainStatus(ctx, &pb.ChainStatusRequest{
		Id: id,
	})
	if err != nil {
		return types.ChainStatus{}, err
	}

	return types.ChainStatus{
		ID:      reply.Chain.Id,
		Enabled: reply.Chain.Enabled,
		Config:  reply.Chain.Config,
	}, nil
}

func (r *relayerClient) ChainStatuses(ctx context.Context, offset, limit int) (chains []types.ChainStatus, count int, err error) {
	var reply *pb.ChainStatusesReply
	reply, err = r.relayer.ChainStatuses(ctx, &pb.ChainStatusesRequest{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		return
	}
	count = int(reply.Count)
	for _, c := range reply.Chains {
		chains = append(chains, types.ChainStatus{
			ID:      c.Id,
			Enabled: c.Enabled,
			Config:  c.Config,
		})
	}

	return
}

func (r *relayerClient) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	reply, err := r.relayer.NodeStatuses(ctx, &pb.NodeStatusesRequest{
		Offset:   int32(offset),
		Limit:    int32(limit),
		ChainIDs: chainIDs,
	})
	if err != nil {
		return nil, -1, err
	}
	for _, n := range reply.Nodes {
		nodes = append(nodes, types.NodeStatus{
			ChainID: n.ChainID,
			Name:    n.Name,
			Config:  n.Config,
			State:   n.State,
		})
	}
	count = int(reply.Count)
	return
}

func (r *relayerClient) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	_, err := r.relayer.SendTx(ctx, &pb.SendTxRequest{
		ChainID:      chainID,
		From:         from,
		To:           to,
		Amount:       pb.NewBigIntFromInt(amount),
		BalanceCheck: balanceCheck,
	})
	return err
}

var _ pb.RelayerServer = (*relayerServer)(nil)

// relayerServer exposes [types.Relayer] as a GRPC [pb.RelayerServer].
type relayerServer struct {
	pb.UnimplementedRelayerServer

	*lggrBroker

	impl Relayer
}

func newChainRelayerServer(impl Relayer, lb *lggrBroker) *relayerServer {
	return &relayerServer{impl: impl, lggrBroker: lb.named("ChainRelayerServer")}
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
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &serviceServer{srv: cp, stop: func() {
		time.AfterFunc(time.Second, s.GracefulStop)
	}})
	pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: cp.OffchainConfigDigester()})
	pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: cp.ContractConfigTracker()})
	const name = "ConfigProvider"
	id, err := r.serve(s, name, resource{cp, name})
	if err != nil {
		return nil, err
	}
	return &pb.NewConfigProviderReply{ConfigProviderID: id}, nil
}

func (r *relayerServer) NewMedianProvider(ctx context.Context, request *pb.NewMedianProviderRequest) (*pb.NewMedianProviderReply, error) {
	exJobID, err := uuid.FromBytes(request.RelayArgs.ExternalJobID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid bytes for ExternalJobID: %w", err)
	}
	provider, err := r.impl.NewMedianProvider(ctx, types.RelayArgs{
		ExternalJobID: exJobID,
		JobID:         request.RelayArgs.JobID,
		ContractID:    request.RelayArgs.ContractID,
		New:           request.RelayArgs.New,
		RelayConfig:   request.RelayArgs.RelayConfig,
	}, types.PluginArgs{
		TransmitterID: request.PluginArgs.TransmitterID,
		PluginConfig:  request.PluginArgs.PluginConfig,
	})
	if err != nil {
		return nil, err
	}
	const name = "MedianProvider"
	providerRes := resource{name: name, Closer: provider}
	s := grpc.NewServer()
	pb.RegisterServiceServer(s, &serviceServer{srv: provider, stop: func() {
		time.AfterFunc(time.Second, s.GracefulStop)
		r.closeAll(providerRes)
	}})
	pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: provider.OffchainConfigDigester()})
	pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: provider.ContractConfigTracker()})
	pb.RegisterContractTransmitterServer(s, &contractTransmitterServer{impl: provider.ContractTransmitter()})
	pb.RegisterReportCodecServer(s, &reportCodecServer{impl: provider.ReportCodec()})
	pb.RegisterMedianContractServer(s, &medianContractServer{impl: provider.MedianContract()})
	pb.RegisterOnchainConfigCodecServer(s, &onchainConfigCodecServer{impl: provider.OnchainConfigCodec()})
	id, err := r.serve(s, name)
	if err != nil {
		r.closeAll(providerRes)
		return nil, err
	}
	return &pb.NewMedianProviderReply{MedianProviderID: id}, nil
}

func (r *relayerServer) NewMercuryProvider(ctx context.Context, request *pb.NewMercuryProviderRequest) (*pb.NewMercuryProviderReply, error) {
	return nil, errors.New("mercury is not supported")
}

func (r *relayerServer) ChainStatus(ctx context.Context, request *pb.ChainStatusRequest) (*pb.ChainStatusReply, error) {
	chain, err := r.impl.ChainStatus(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ChainStatusReply{Chain: &pb.ChainStatus{
		Id:      chain.ID,
		Enabled: chain.Enabled,
		Config:  chain.Config,
	}}, nil
}

func (r *relayerServer) ChainStatuses(ctx context.Context, request *pb.ChainStatusesRequest) (*pb.ChainStatusesReply, error) {
	chains, count, err := r.impl.ChainStatuses(ctx, int(request.Offset), int(request.Limit))
	if err != nil {
		return nil, err
	}
	reply := &pb.ChainStatusesReply{Count: int32(count)}
	for _, c := range chains {
		reply.Chains = append(reply.Chains, &pb.ChainStatus{
			Id:      c.ID,
			Enabled: c.Enabled,
			Config:  c.Config,
		})
	}
	return reply, nil
}

func (r *relayerServer) NodeStatuses(ctx context.Context, request *pb.NodeStatusesRequest) (*pb.NodeStatusesReply, error) {
	nodeConfigs, count, err := r.impl.NodeStatuses(ctx, int(request.Offset), int(request.Limit), request.ChainIDs...)
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
	return &pb.NodeStatusesReply{Nodes: nodes, Count: int32(count)}, nil
}
func (r *relayerServer) SendTx(ctx context.Context, request *pb.SendTxRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, r.impl.SendTx(ctx, request.ChainID, request.From, request.To, request.Amount.Int(), request.BalanceCheck)
}

func healthReport(s map[string]string) (hr map[string]error) {
	hr = make(map[string]error, len(s))
	for n, e := range s {
		hr[n] = errors.New(e)
	}
	return hr
}
