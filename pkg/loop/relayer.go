package loop

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/smartcontractkit/chainlink-relay/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

var _ Relayer = (*relayerClient)(nil)

// relayerClient adapts a GRPC [pb.RelayerClient] to implement [Relayer].
type relayerClient struct {
	*brokerExt
	*serviceClient

	relayer pb.RelayerClient
}

func newRelayerClient(b *brokerExt, conn grpc.ClientConnInterface) *relayerClient {
	b = b.named("ChainRelayerClient")
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
	return newConfigProviderClient(r.named("ConfigProviderClient"), cc), nil
}

func (r *relayerClient) NewMedianProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MedianProvider, error) {
	cc := r.newClientConn("MedianProvider", func(ctx context.Context) (uint32, resources, error) {
		reply, err := r.relayer.NewMedianProvider(ctx, &pb.NewMedianProviderRequest{
			RelayArgs: &pb.RelayArgs{
				ExternalJobID: rargs.ExternalJobID[:],
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
			return 0, nil, err
		}
		return reply.MedianProviderID, nil, nil
	})
	return newMedianProviderClient(r.brokerExt, cc), nil
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

	*brokerExt

	impl Relayer
}

func newChainRelayerServer(impl Relayer, b *brokerExt) *relayerServer {
	return &relayerServer{impl: impl, brokerExt: b.named("ChainRelayerServer")}
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
	id, _, err := r.serve(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: cp})
		pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: cp.OffchainConfigDigester()})
		pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: cp.ContractConfigTracker()})
	}, resource{cp, name})
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
	err = provider.Start(ctx)
	if err != nil {
		return nil, err
	}
	const name = "MedianProvider"
	providerRes := resource{name: name, Closer: provider}

	id, _, err := r.serve(name, func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &serviceServer{srv: provider})
		pb.RegisterOffchainConfigDigesterServer(s, &offchainConfigDigesterServer{impl: provider.OffchainConfigDigester()})
		pb.RegisterContractConfigTrackerServer(s, &contractConfigTrackerServer{impl: provider.ContractConfigTracker()})
		pb.RegisterContractTransmitterServer(s, &contractTransmitterServer{impl: provider.ContractTransmitter()})
		pb.RegisterReportCodecServer(s, &reportCodecServer{impl: provider.ReportCodec()})
		pb.RegisterMedianContractServer(s, &medianContractServer{impl: provider.MedianContract()})
		pb.RegisterOnchainConfigCodecServer(s, &onchainConfigCodecServer{impl: provider.OnchainConfigCodec()})
	}, providerRes)
	if err != nil {
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
		var err error
		if e != "" {
			err = errors.New(e)
		}
		hr[n] = err
	}
	return hr
}
