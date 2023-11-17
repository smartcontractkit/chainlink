package internal

import (
	"context"
	"math"

	"google.golang.org/grpc"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

var _ types.ConfigProvider = (*configProviderClient)(nil)

type configProviderClient struct {
	*serviceClient
	offchainDigester libocr.OffchainConfigDigester
	contractTracker  libocr.ContractConfigTracker
}

func newConfigProviderClient(b *brokerExt, cc grpc.ClientConnInterface) *configProviderClient {
	c := &configProviderClient{serviceClient: newServiceClient(b, cc)}
	c.offchainDigester = &offchainConfigDigesterClient{b, pb.NewOffchainConfigDigesterClient(cc)}
	c.contractTracker = &contractConfigTrackerClient{pb.NewContractConfigTrackerClient(cc)}
	return c
}

func (c *configProviderClient) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return c.offchainDigester
}

func (c *configProviderClient) ContractConfigTracker() libocr.ContractConfigTracker {
	return c.contractTracker
}

var _ libocr.OffchainConfigDigester = (*offchainConfigDigesterClient)(nil)

type offchainConfigDigesterClient struct {
	*brokerExt
	grpc pb.OffchainConfigDigesterClient
}

func (o *offchainConfigDigesterClient) ConfigDigest(config libocr.ContractConfig) (digest libocr.ConfigDigest, err error) {
	ctx, cancel := o.stopCtx()
	defer cancel()

	var reply *pb.ConfigDigestReply
	reply, err = o.grpc.ConfigDigest(ctx, &pb.ConfigDigestRequest{
		ContractConfig: pbContractConfig(config),
	})
	if err != nil {
		return
	}
	if l := len(reply.ConfigDigest); l != 32 {
		err = ErrConfigDigestLen(l)
		return
	}
	copy(digest[:], reply.ConfigDigest)
	return
}

func (o *offchainConfigDigesterClient) ConfigDigestPrefix() (libocr.ConfigDigestPrefix, error) {
	ctx, cancel := o.stopCtx()
	defer cancel()

	reply, err := o.grpc.ConfigDigestPrefix(ctx, &pb.ConfigDigestPrefixRequest{})
	if err != nil {
		return 0, err
	}
	return libocr.ConfigDigestPrefix(reply.ConfigDigestPrefix), nil
}

var _ pb.OffchainConfigDigesterServer = (*offchainConfigDigesterServer)(nil)

type offchainConfigDigesterServer struct {
	pb.UnimplementedOffchainConfigDigesterServer
	impl libocr.OffchainConfigDigester
}

func (o *offchainConfigDigesterServer) ConfigDigest(ctx context.Context, request *pb.ConfigDigestRequest) (*pb.ConfigDigestReply, error) {
	if request.ContractConfig.F > math.MaxUint8 {
		return nil, ErrUint8Bounds{Name: "F", U: request.ContractConfig.F}
	}
	cc := libocr.ContractConfig{
		ConfigCount:           request.ContractConfig.ConfigCount,
		F:                     uint8(request.ContractConfig.F),
		OnchainConfig:         request.ContractConfig.OnchainConfig,
		OffchainConfigVersion: request.ContractConfig.OffchainConfigVersion,
		OffchainConfig:        request.ContractConfig.OffchainConfig,
	}
	copy(cc.ConfigDigest[:], request.ContractConfig.ConfigDigest)
	for _, s := range request.ContractConfig.Signers {
		cc.Signers = append(cc.Signers, s)
	}
	for _, t := range request.ContractConfig.Transmitters {
		cc.Transmitters = append(cc.Transmitters, libocr.Account(t))
	}
	cd, err := o.impl.ConfigDigest(cc)
	if err != nil {
		return nil, err
	}
	return &pb.ConfigDigestReply{ConfigDigest: cd[:]}, nil
}

func (o *offchainConfigDigesterServer) ConfigDigestPrefix(ctx context.Context, request *pb.ConfigDigestPrefixRequest) (*pb.ConfigDigestPrefixReply, error) {
	p, err := o.impl.ConfigDigestPrefix()
	if err != nil {
		return nil, err
	}
	return &pb.ConfigDigestPrefixReply{ConfigDigestPrefix: uint32(p)}, nil
}

var _ libocr.ContractConfigTracker = (*contractConfigTrackerClient)(nil)

type contractConfigTrackerClient struct {
	grpc pb.ContractConfigTrackerClient
}

func (c *contractConfigTrackerClient) Notify() <-chan struct{} { return nil }

func (c *contractConfigTrackerClient) LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest libocr.ConfigDigest, err error) {
	var reply *pb.LatestConfigDetailsReply
	reply, err = c.grpc.LatestConfigDetails(ctx, &pb.LatestConfigDetailsRequest{})
	if err != nil {
		return
	}
	changedInBlock = reply.ChangedInBlock
	if l := len(reply.ConfigDigest); l != 32 {
		err = ErrConfigDigestLen(l)
		return
	}
	copy(configDigest[:], reply.ConfigDigest)
	return
}

func (c *contractConfigTrackerClient) LatestConfig(ctx context.Context, changedInBlock uint64) (cfg libocr.ContractConfig, err error) {
	var reply *pb.LatestConfigReply
	reply, err = c.grpc.LatestConfig(ctx, &pb.LatestConfigRequest{
		ChangedInBlock: changedInBlock,
	})
	if err != nil {
		return
	}
	if l := len(reply.ContractConfig.ConfigDigest); l != 32 {
		err = ErrConfigDigestLen(l)
		return
	}
	copy(cfg.ConfigDigest[:], reply.ContractConfig.ConfigDigest)
	cfg.ConfigCount = reply.ContractConfig.ConfigCount
	for _, s := range reply.ContractConfig.Signers {
		cfg.Signers = append(cfg.Signers, s)
	}
	for _, t := range reply.ContractConfig.Transmitters {
		cfg.Transmitters = append(cfg.Transmitters, libocr.Account(t))
	}
	if reply.ContractConfig.F > math.MaxUint8 {
		err = ErrUint8Bounds{Name: "F", U: reply.ContractConfig.F}
		return
	}
	cfg.F = uint8(reply.ContractConfig.F)
	cfg.OnchainConfig = reply.ContractConfig.OnchainConfig
	cfg.OffchainConfigVersion = reply.ContractConfig.OffchainConfigVersion
	cfg.OffchainConfig = reply.ContractConfig.OffchainConfig

	return
}

func (c *contractConfigTrackerClient) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	var reply *pb.LatestBlockHeightReply
	reply, err = c.grpc.LatestBlockHeight(ctx, &pb.LatestBlockHeightRequest{})
	if err != nil {
		return
	}
	blockHeight = reply.BlockHeight
	return
}

var _ pb.ContractConfigTrackerServer = (*contractConfigTrackerServer)(nil)

type contractConfigTrackerServer struct {
	pb.UnimplementedContractConfigTrackerServer
	impl libocr.ContractConfigTracker
}

func (c *contractConfigTrackerServer) LatestConfigDetails(ctx context.Context, request *pb.LatestConfigDetailsRequest) (*pb.LatestConfigDetailsReply, error) {
	changedInBlock, configDigest, err := c.impl.LatestConfigDetails(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.LatestConfigDetailsReply{ChangedInBlock: changedInBlock, ConfigDigest: configDigest[:]}, nil
}

func (c *contractConfigTrackerServer) LatestConfig(ctx context.Context, request *pb.LatestConfigRequest) (*pb.LatestConfigReply, error) {
	cc, err := c.impl.LatestConfig(ctx, request.ChangedInBlock)
	if err != nil {
		return nil, err
	}
	return &pb.LatestConfigReply{ContractConfig: pbContractConfig(cc)}, nil
}

func (c *contractConfigTrackerServer) LatestBlockHeight(ctx context.Context, request *pb.LatestBlockHeightRequest) (*pb.LatestBlockHeightReply, error) {
	blockHeight, err := c.impl.LatestBlockHeight(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.LatestBlockHeightReply{BlockHeight: blockHeight}, nil
}

func pbContractConfig(cc libocr.ContractConfig) *pb.ContractConfig {
	r := &pb.ContractConfig{
		ConfigDigest:          cc.ConfigDigest[:],
		ConfigCount:           cc.ConfigCount,
		F:                     uint32(cc.F),
		OnchainConfig:         cc.OnchainConfig,
		OffchainConfigVersion: cc.OffchainConfigVersion,
		OffchainConfig:        cc.OffchainConfig,
	}
	for _, s := range cc.Signers {
		r.Signers = append(r.Signers, s)
	}
	for _, t := range cc.Transmitters {
		r.Transmitters = append(r.Transmitters, string(t))
	}
	return r
}
