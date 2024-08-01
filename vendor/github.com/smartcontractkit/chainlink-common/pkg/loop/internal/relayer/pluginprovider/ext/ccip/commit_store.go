package ccip

import (
	"context"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// CommitStoreGRPCClient implements [cciptypes.CommitStoreReader] by wrapping a
// [ccippb.CommitStoreReaderGRPCClient] grpc client.
// It is used by a ReportingPlugin to call the CommitStoreReader service, which
// is hosted by the relayer
type CommitStoreGRPCClient struct {
	client ccippb.CommitStoreReaderClient

	//  brokerExt is use to allocate and serve the gas estimator server.
	//  must be the same as that used by the server
	//  TODO BCF-3061: note unsure if this really has to change for the proxy case or not
	//  marking it so that it is considered when we implement the proxy.
	//  the reason it may not need to change is that the gas estimator server is
	//  a static resource of the offramp reader server. It is not created directly by the client.
	b    *net.BrokerExt
	conn grpc.ClientConnInterface
}

func NewCommitStoreReaderGRPCClient(brokerExt *net.BrokerExt, cc grpc.ClientConnInterface) *CommitStoreGRPCClient {
	return &CommitStoreGRPCClient{client: ccippb.NewCommitStoreReaderClient(cc), b: brokerExt, conn: cc}
}

// CommitStoreGRPCServer implements [ccippb.CommitStoreReaderServer] by wrapping a
// [cciptypes.CommitStoreReader] implementation.
// This server is hosted by the relayer and is called ReportingPlugin via
// the [CommitStoreGRPCClient]
type CommitStoreGRPCServer struct {
	ccippb.UnimplementedCommitStoreReaderServer

	impl ccip.CommitStoreReader

	//  brokerExt is use to allocate and serve the gas estimator server.
	//  must be the same as that used by the server
	//  TODO BCF-3061. see the comment in OffRampReaderGRPCClient for more details
	b                    *net.BrokerExt
	gasEstimatorServerID uint32 // allocated by the broker on creation of the off ramp server

	deps []io.Closer
}

func NewCommitStoreReaderGRPCServer(impl ccip.CommitStoreReader, brokerExt *net.BrokerExt) (*CommitStoreGRPCServer, error) {
	estimator, err := impl.GasPriceEstimator(context.Background())
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	estimatorHandler := NewCommitGasEstimatorGRPCServer(estimator)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	estimatorID, spawnedServer, err := brokerExt.ServeNew("CommitStoreReader.OffRampGasEstimator", func(s *grpc.Server) {
		ccippb.RegisterGasPriceEstimatorCommitServer(s, estimatorHandler)
	})
	if err != nil {
		return nil, err
	}

	var deps []io.Closer
	deps = append(deps, impl, spawnedServer)

	return &CommitStoreGRPCServer{impl: impl, deps: deps, b: brokerExt, gasEstimatorServerID: estimatorID}, nil
}

// ensure the types are satisfied
var _ ccippb.CommitStoreReaderServer = (*CommitStoreGRPCServer)(nil)
var _ ccip.CommitStoreReader = (*CommitStoreGRPCClient)(nil)
var _ goplugin.GRPCClientConn = (*CommitStoreGRPCClient)(nil)

// ClientConn implements goplugin.GRPCClientConn.
func (c *CommitStoreGRPCClient) ClientConn() grpc.ClientConnInterface {
	return c.conn
}

// ChangeConfig implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (ccip.Address, error) {
	resp, err := c.client.ChangeConfig(ctx, &ccippb.CommitStoreChangeConfigRequest{
		OnchainConfig:  onchainConfig,
		OffchainConfig: offchainConfig,
	})
	if err != nil {
		return ccip.Address(""), err
	}
	return ccip.Address(resp.Address), nil
}

func (c *CommitStoreGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), c.client)
}

// DecodeCommitReport implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) DecodeCommitReport(ctx context.Context, report []byte) (ccip.CommitStoreReport, error) {
	resp, err := c.client.DecodeCommitReport(ctx, &ccippb.DecodeCommitReportRequest{EncodedReport: report})
	if err != nil {
		return ccip.CommitStoreReport{}, err
	}
	return commitStoreReport(resp.Report)
}

// EncodeCommitReport implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) EncodeCommitReport(ctx context.Context, report ccip.CommitStoreReport) ([]byte, error) {
	pb, err := commitStoreReportPB(report)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.EncodeCommitReport(ctx, &ccippb.EncodeCommitReportRequest{Report: pb})
	if err != nil {
		return nil, err
	}
	return resp.EncodedReport, nil
}

// GasPriceEstimator implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) GasPriceEstimator(ctx context.Context) (ccip.GasPriceEstimatorCommit, error) {
	resp, err := c.client.GetCommitGasPriceEstimator(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: this works because the broker is shared and the id refers to a resource served by the broker
	gasEstimatorConn, err := c.b.Dial(resp.GasPriceEstimatorId)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup gas estimator service for off ramp reader at %d: %w", resp.GasPriceEstimatorId, err)
	}
	// need to wrap grpc offRamp into the desired interface
	gasEstimator := NewCommitGasEstimatorGRPCClient(gasEstimatorConn)
	// need to hydrate the gas price estimator from the server id
	return gasEstimator, nil
}

// GetAcceptedCommitReportsGteTimestamp implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confirmations int) ([]ccip.CommitStoreReportWithTxMeta, error) {
	resp, err := c.client.GetAcceptedCommitReportsGteTimestamp(ctx, &ccippb.GetAcceptedCommitReportsGteTimestampRequest{
		Timestamp:     timestamppb.New(ts),
		Confirmations: uint64(confirmations),
	})
	if err != nil {
		return nil, err
	}
	return commitStoreReportWithTxMetaSlice(resp.Reports)
}

// GetCommitReportMatchingSeqNum implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confirmations int) ([]ccip.CommitStoreReportWithTxMeta, error) {
	resp, err := c.client.GetCommitReportMatchingSequenceNumber(ctx, &ccippb.GetCommitReportMatchingSequenceNumberRequest{
		SequenceNumber: seqNum,
		Confirmations:  uint64(confirmations),
	})
	if err != nil {
		return nil, err
	}
	return commitStoreReportWithTxMetaSlice(resp.Reports)
}

// GetCommitStoreStaticConfig implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) GetCommitStoreStaticConfig(ctx context.Context) (ccip.CommitStoreStaticConfig, error) {
	resp, err := c.client.GetCommitStoreStaticConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return ccip.CommitStoreStaticConfig{}, err
	}
	return ccip.CommitStoreStaticConfig{
		ChainSelector:       resp.StaticConfig.ChainSelector,
		SourceChainSelector: resp.StaticConfig.SourceChainSelector,
		OnRamp:              ccip.Address(resp.StaticConfig.OnRamp),
		ArmProxy:            ccip.Address(resp.StaticConfig.ArmProxy),
	}, nil
}

// GetExpectedNextSequenceNumber implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) GetExpectedNextSequenceNumber(ctx context.Context) (uint64, error) {
	resp, err := c.client.GetExpectedNextSequenceNumber(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}
	return resp.SequenceNumber, nil
}

// GetLatestPriceEpochAndRound implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) GetLatestPriceEpochAndRound(ctx context.Context) (uint64, error) {
	resp, err := c.client.GetLatestPriceEpochAndRound(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}
	return resp.EpochAndRound, nil
}

// IsBlessed implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	resp, err := c.client.IsBlessed(ctx, &ccippb.IsBlessedRequest{Root: root[:]})
	if err != nil {
		return false, err
	}
	return resp.IsBlessed, nil
}

// IsDestChainHealthy implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) IsDestChainHealthy(ctx context.Context) (bool, error) {
	resp, err := c.client.IsDestChainHealthy(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.IsHealthy, nil
}

// IsDown implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) IsDown(ctx context.Context) (bool, error) {
	resp, err := c.client.IsDown(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}
	return resp.IsDown, nil
}

// OffchainConfig implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) OffchainConfig(ctx context.Context) (ccip.CommitOffchainConfig, error) {
	resp, err := c.client.GetOffchainConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return ccip.CommitOffchainConfig{}, err
	}
	return ccip.CommitOffchainConfig{
		GasPriceDeviationPPB:   resp.OffchainConfig.GasPriceDeviationPpb,
		GasPriceHeartBeat:      resp.OffchainConfig.GasPriceHeartbeat.AsDuration(),
		TokenPriceDeviationPPB: resp.OffchainConfig.TokenPriceDeviationPpb,
		TokenPriceHeartBeat:    resp.OffchainConfig.TokenPriceHeartbeat.AsDuration(),
		InflightCacheExpiry:    resp.OffchainConfig.InflightCacheExpiry.AsDuration(),
	}, nil
}

// VerifyExecutionReport implements ccip.CommitStoreReader.
func (c *CommitStoreGRPCClient) VerifyExecutionReport(ctx context.Context, report ccip.ExecReport) (bool, error) {
	resp, err := c.client.VerifyExecutionReport(ctx, &ccippb.VerifyExecutionReportRequest{Report: executionReportPB(report)})
	if err != nil {
		return false, err
	}
	return resp.IsValid, nil
}

// Server implementation

// ChangeConfig implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) ChangeConfig(ctx context.Context, req *ccippb.CommitStoreChangeConfigRequest) (*ccippb.CommitStoreChangeConfigResponse, error) {
	addr, err := c.impl.ChangeConfig(ctx, req.OnchainConfig, req.OffchainConfig)
	if err != nil {
		return nil, err
	}
	return &ccippb.CommitStoreChangeConfigResponse{Address: string(addr)}, nil
}

func (c *CommitStoreGRPCServer) Close(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(c.deps).Close()
}

// DecodeCommitReport implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) DecodeCommitReport(ctx context.Context, req *ccippb.DecodeCommitReportRequest) (*ccippb.DecodeCommitReportResponse, error) {
	r, err := c.impl.DecodeCommitReport(ctx, req.EncodedReport)
	if err != nil {
		return nil, err
	}
	pb, err := commitStoreReportPB(r)
	if err != nil {
		return nil, err
	}
	return &ccippb.DecodeCommitReportResponse{Report: pb}, nil
}

// EncodeCommitReport implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) EncodeCommitReport(ctx context.Context, req *ccippb.EncodeCommitReportRequest) (*ccippb.EncodeCommitReportResponse, error) {
	r, err := commitStoreReport(req.Report)
	if err != nil {
		return nil, err
	}
	encoded, err := c.impl.EncodeCommitReport(ctx, r)
	if err != nil {
		return nil, err
	}
	return &ccippb.EncodeCommitReportResponse{EncodedReport: encoded}, nil
}

// GetAcceptedCommitReportsGteTimestamp implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, req *ccippb.GetAcceptedCommitReportsGteTimestampRequest) (*ccippb.GetAcceptedCommitReportsGteTimestampResponse, error) {
	reports, err := c.impl.GetAcceptedCommitReportsGteTimestamp(ctx, req.Timestamp.AsTime(), int(req.Confirmations))
	if err != nil {
		return nil, err
	}
	pbReports, err := commitStoreReportWithTxMetaPBSlice(reports)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetAcceptedCommitReportsGteTimestampResponse{Reports: pbReports}, nil
}

// GetCommitGasPriceEstimator implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetCommitGasPriceEstimator(ctx context.Context, req *emptypb.Empty) (*ccippb.GetCommitGasPriceEstimatorResponse, error) {
	return &ccippb.GetCommitGasPriceEstimatorResponse{GasPriceEstimatorId: c.gasEstimatorServerID}, nil
}

// GetCommitStoreStaticConfig implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetCommitStoreStaticConfig(ctx context.Context, req *emptypb.Empty) (*ccippb.GetCommitStoreStaticConfigResponse, error) {
	config, err := c.impl.GetCommitStoreStaticConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetCommitStoreStaticConfigResponse{StaticConfig: &ccippb.CommitStoreStaticConfig{
		ChainSelector:       config.ChainSelector,
		SourceChainSelector: config.SourceChainSelector,
		OnRamp:              string(config.OnRamp),
		ArmProxy:            string(config.ArmProxy),
	}}, nil
}

// GetExpectedNextSequenceNumber implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetExpectedNextSequenceNumber(ctx context.Context, req *emptypb.Empty) (*ccippb.GetExpectedNextSequenceNumberResponse, error) {
	seqNum, err := c.impl.GetExpectedNextSequenceNumber(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetExpectedNextSequenceNumberResponse{SequenceNumber: seqNum}, nil
}

// GetLatestPriceEpochAndRound implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetLatestPriceEpochAndRound(ctx context.Context, req *emptypb.Empty) (*ccippb.GetLatestPriceEpochAndRoundResponse, error) {
	epoch, err := c.impl.GetLatestPriceEpochAndRound(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetLatestPriceEpochAndRoundResponse{EpochAndRound: epoch}, nil
}

// GetOffchainConfig implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetOffchainConfig(ctx context.Context, req *emptypb.Empty) (*ccippb.GetOffchainConfigResponse, error) {
	config, err := c.impl.OffchainConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetOffchainConfigResponse{
		OffchainConfig: &ccippb.CommitOffchainConfig{
			GasPriceDeviationPpb:   config.GasPriceDeviationPPB,
			GasPriceHeartbeat:      durationpb.New(config.GasPriceHeartBeat),
			TokenPriceDeviationPpb: config.TokenPriceDeviationPPB,
			TokenPriceHeartbeat:    durationpb.New(config.TokenPriceHeartBeat),
			InflightCacheExpiry:    durationpb.New(config.InflightCacheExpiry),
		},
	}, nil
}

// GeteCommitReportMatchingSequenceNumber implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) GetCommitReportMatchingSequenceNumber(ctx context.Context, req *ccippb.GetCommitReportMatchingSequenceNumberRequest) (*ccippb.GetCommitReportMatchingSequenceNumberResponse, error) {
	reports, err := c.impl.GetCommitReportMatchingSeqNum(ctx, req.SequenceNumber, int(req.Confirmations))
	if err != nil {
		return nil, err
	}
	pbReports, err := commitStoreReportWithTxMetaPBSlice(reports)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetCommitReportMatchingSequenceNumberResponse{Reports: pbReports}, nil
}

// IsBlessed implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) IsBlessed(ctx context.Context, req *ccippb.IsBlessedRequest) (*ccippb.IsBlessedResponse, error) {
	r, err := merkleRoot(req.Root)
	if err != nil {
		return nil, err
	}
	blessed, err := c.impl.IsBlessed(ctx, r)
	if err != nil {
		return nil, err
	}
	return &ccippb.IsBlessedResponse{IsBlessed: blessed}, nil
}

// IsDestChainHealthy implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) IsDestChainHealthy(ctx context.Context, req *emptypb.Empty) (*ccippb.IsDestChainHealthyResponse, error) {
	healthy, err := c.impl.IsDestChainHealthy(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.IsDestChainHealthyResponse{IsHealthy: healthy}, nil
}

// IsDown implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) IsDown(ctx context.Context, req *emptypb.Empty) (*ccippb.IsDownResponse, error) {
	down, err := c.impl.IsDown(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.IsDownResponse{IsDown: down}, nil
}

// VerifyExecutionReport implements ccippb.CommitStoreReaderServer.
func (c *CommitStoreGRPCServer) VerifyExecutionReport(ctx context.Context, req *ccippb.VerifyExecutionReportRequest) (*ccippb.VerifyExecutionReportResponse, error) {
	r, err := execReport(req.Report)
	if err != nil {
		return nil, err
	}
	valid, err := c.impl.VerifyExecutionReport(ctx, r)
	if err != nil {
		return nil, err
	}
	return &ccippb.VerifyExecutionReportResponse{IsValid: valid}, nil
}

// AddDep adds a closer to the list of dependencies that will be closed when the server is closed.
func (c *CommitStoreGRPCServer) AddDep(dep io.Closer) *CommitStoreGRPCServer {
	c.deps = append(c.deps, dep)
	return c
}

func commitStoreReport(pb *ccippb.CommitStoreReport) (ccip.CommitStoreReport, error) {
	root, err := merkleRoot(pb.MerkleRoot)
	if err != nil {
		return ccip.CommitStoreReport{}, fmt.Errorf("cannot convert merkle root: %w", err)
	}
	out := ccip.CommitStoreReport{
		TokenPrices: tokenPrices(pb.TokenPrices),
		GasPrices:   gasPrices(pb.GasPrices),
		Interval:    commitStoreInterval(pb.Interval),
		MerkleRoot:  root,
	}
	return out, nil
}

func tokenPrices(pb []*ccippb.TokenPrice) []ccip.TokenPrice {
	out := make([]ccip.TokenPrice, len(pb))
	for i, p := range pb {
		out[i] = tokenPrice(p)
	}
	return out
}

func gasPrices(pb []*ccippb.GasPrice) []ccip.GasPrice {
	out := make([]ccip.GasPrice, len(pb))
	for i, p := range pb {
		out[i] = gasPrice(p)
	}
	return out
}

func gasPrice(pb *ccippb.GasPrice) ccip.GasPrice {
	return ccip.GasPrice{
		DestChainSelector: pb.DestChainSelector,
		Value:             pb.Value.Int(),
	}
}

func commitStoreInterval(pb *ccippb.CommitStoreInterval) ccip.CommitStoreInterval {
	return ccip.CommitStoreInterval{
		Min: pb.Min,
		Max: pb.Max,
	}
}

func merkleRoot(pb []byte) ([32]byte, error) {
	if len(pb) != 32 {
		return [32]byte{}, fmt.Errorf("expected 32 bytes, got %d", len(pb))
	}
	var out [32]byte
	copy(out[:], pb)
	return out, nil
}

func commitStoreReportPB(r ccip.CommitStoreReport) (*ccippb.CommitStoreReport, error) {
	if len(r.MerkleRoot) != 32 {
		return nil, fmt.Errorf("invalid merkle root: expected 32 bytes, got %d", len(r.MerkleRoot))
	}
	pb := &ccippb.CommitStoreReport{
		TokenPrices: tokenPricesPB(r.TokenPrices),
		GasPrices:   gasPricesPB(r.GasPrices),
		Interval:    commitStoreIntervalPB(r.Interval),
		MerkleRoot:  r.MerkleRoot[:],
	}
	return pb, nil
}

func tokenPricesPB(r []ccip.TokenPrice) []*ccippb.TokenPrice {
	out := make([]*ccippb.TokenPrice, len(r))
	for i, p := range r {
		out[i] = tokenPricePB(p)
	}
	return out
}

func gasPricesPB(r []ccip.GasPrice) []*ccippb.GasPrice {
	out := make([]*ccippb.GasPrice, len(r))
	for i, p := range r {
		out[i] = gasPricePB(p)
	}
	return out
}

func commitStoreIntervalPB(r ccip.CommitStoreInterval) *ccippb.CommitStoreInterval {
	return &ccippb.CommitStoreInterval{
		Min: r.Min,
		Max: r.Max,
	}
}

func commitStoreReportWithTxMetaSlice(pb []*ccippb.CommitStoreReportWithTxMeta) ([]ccip.CommitStoreReportWithTxMeta, error) {
	out := make([]ccip.CommitStoreReportWithTxMeta, len(pb))
	var err error
	for i, p := range pb {
		out[i], err = commitStoreReportWithTxMeta(p)
		if err != nil {
			return nil, fmt.Errorf("cannot convert commit store report with tx meta: %w", err)
		}
	}
	return out, nil
}

func commitStoreReportWithTxMeta(pb *ccippb.CommitStoreReportWithTxMeta) (ccip.CommitStoreReportWithTxMeta, error) {
	r, err := commitStoreReport(pb.Report)
	if err != nil {
		return ccip.CommitStoreReportWithTxMeta{}, fmt.Errorf("cannot convert commit store report: %w", err)
	}
	return ccip.CommitStoreReportWithTxMeta{
		TxMeta:            txMeta(pb.TxMeta),
		CommitStoreReport: r,
	}, nil
}

func commitStoreReportWithTxMetaPBSlice(r []ccip.CommitStoreReportWithTxMeta) ([]*ccippb.CommitStoreReportWithTxMeta, error) {
	out := make([]*ccippb.CommitStoreReportWithTxMeta, len(r))
	var err error
	for i, p := range r {
		out[i], err = commitStoreReportWithTxMetaPB(p)
		if err != nil {
			return nil, fmt.Errorf("cannot convert commit store report %v at %d with tx meta: %w", p, i, err)
		}
	}
	return out, nil
}

func commitStoreReportWithTxMetaPB(r ccip.CommitStoreReportWithTxMeta) (*ccippb.CommitStoreReportWithTxMeta, error) {
	report, err := commitStoreReportPB(r.CommitStoreReport)
	if err != nil {
		return nil, fmt.Errorf("cannot convert commit store report: %w", err)
	}
	return &ccippb.CommitStoreReportWithTxMeta{
		TxMeta: txMetaToPB(r.TxMeta),
		Report: report,
	}, nil
}
