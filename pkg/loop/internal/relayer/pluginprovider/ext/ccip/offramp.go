package ccip

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	ccippb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
)

// OffRampReaderGRPCClient implement [cciptypes.OffRampReader] by wrapping a grpc client connection
// this client will be used by the CCIP loop service to communicate with the offramp reader
type OffRampReaderGRPCClient struct {
	client ccippb.OffRampReaderClient

	//  brokerExt is use to allocate and serve the gas estimator server.
	//  must be the same as that used by the server
	//  TODO BCF-3061: note unsure if this really has to change for the proxy case or not
	//  marking it so that it is considered when we implement the proxy.
	//  the reason it may not need to change is that the gas estimator server is
	//  a static resource of the offramp reader server. It is not created directly by the client.
	b *net.BrokerExt

	conn grpc.ClientConnInterface
}

// NewOffRampReaderGRPCClient creates a new OffRampReaderGRPCClient. It is used by the reporting plugin to call the offramp reader service.
// The client is created by wrapping a grpc client connection. It requires a brokerExt to allocate and serve the gas estimator server.
// *must* be the same broker used by the server BCF-3061
func NewOffRampReaderGRPCClient(brokerExt *net.BrokerExt, cc grpc.ClientConnInterface) *OffRampReaderGRPCClient {
	return &OffRampReaderGRPCClient{client: ccippb.NewOffRampReaderClient(cc), b: brokerExt, conn: cc}
}

// OffRampReaderGRPCServer implements [ccippb.OffRampReaderServer] by wrapping a [cciptypes.OffRampReader] implementation.
type OffRampReaderGRPCServer struct {
	ccippb.UnimplementedOffRampReaderServer

	impl cciptypes.OffRampReader

	//  brokerExt is use to allocate and serve the gas estimator server.
	//  must be the same as that used by the server
	//  TODO BCF-3061. see the comment in OffRampReaderGRPCClient for more details
	b                    *net.BrokerExt
	gasEstimatorServerID uint32 // allocated by the broker on creation of the off ramp server

	// must support multiple close handlers because the offramp reader server needs to serve the gas estimator server,
	// which needs to be close when the offramp reader server is closed, as well as the offramp reader server itself
	deps []io.Closer
}

// NewOffRampReaderGRPCServer creates a new OffRampReaderGRPCServer. It is used by the relayer to serve the offramp reader service.
// The server is created by wrapping a [cciptypes.OffRampReader] implementation. It requires a brokerExt to allocate and serve the gas estimator server.
// *must* be the same broker used by the client. BCF-3061
func NewOffRampReaderGRPCServer(impl cciptypes.OffRampReader, brokerExt *net.BrokerExt) (*OffRampReaderGRPCServer, error) {
	// offramp reader server needs to serve the gas estimator server
	estimator, err := impl.GasPriceEstimator(context.Background())
	if err != nil {
		return nil, err
	}
	// wrap the reader in a grpc server and serve it
	estimatorHandler := NewExecGasEstimatorGRPCServer(estimator)
	// the id is handle to the broker, we will need it on the other side to dial the resource
	estimatorID, spawnedServer, err := brokerExt.ServeNew("OffRampReader.OffRampGasEstimator", func(s *grpc.Server) {
		ccippb.RegisterGasPriceEstimatorExecServer(s, estimatorHandler)
	})
	if err != nil {
		return nil, err
	}

	var toClose []io.Closer
	toClose = append(toClose, impl, spawnedServer)
	return &OffRampReaderGRPCServer{
		impl:                 impl,
		gasEstimatorServerID: estimatorID,
		b:                    brokerExt,
		deps:                 toClose}, nil
}

// ensure the types are satisfied
var _ cciptypes.OffRampReader = (*OffRampReaderGRPCClient)(nil)
var _ ccippb.OffRampReaderServer = (*OffRampReaderGRPCServer)(nil)

// Address i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) Address(ctx context.Context) (cciptypes.Address, error) {
	resp, err := o.client.Address(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.Address), nil
}

// ChangeConfig implements [github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) ChangeConfig(ctx context.Context, onchainConfig []byte, offchainConfig []byte) (cciptypes.Address, cciptypes.Address, error) {
	resp, err := o.client.ChangeConfig(ctx, &ccippb.ChangeConfigRequest{
		OnchainConfig:  onchainConfig,
		OffchainConfig: offchainConfig,
	})
	if err != nil {
		return cciptypes.Address(""), cciptypes.Address(""), err
	}

	return cciptypes.Address(resp.OnchainConfigAddress), cciptypes.Address(resp.OffchainConfigAddress), nil
}

func (o *OffRampReaderGRPCClient) ClientConn() grpc.ClientConnInterface {
	return o.conn
}

func (o *OffRampReaderGRPCClient) Close() error {
	return shutdownGRPCServer(context.Background(), o.client)
}

// CurrentRateLimiterState i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) CurrentRateLimiterState(ctx context.Context) (cciptypes.TokenBucketRateLimit, error) {
	resp, err := o.client.CurrentRateLimiterState(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.TokenBucketRateLimit{}, err
	}
	return tokenBucketRateLimit(resp.RateLimiter), nil
}

// DecodeExecutionReport [github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) DecodeExecutionReport(ctx context.Context, report []byte) (cciptypes.ExecReport, error) {
	resp, err := o.client.DecodeExecutionReport(ctx, &ccippb.DecodeExecutionReportRequest{
		Report: report,
	})
	if err != nil {
		return cciptypes.ExecReport{}, err
	}

	return execReport(resp.Report)
}

// EncodeExecutionReport [github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) EncodeExecutionReport(ctx context.Context, report cciptypes.ExecReport) ([]byte, error) {
	reportPB := executionReportPB(report)

	resp, err := o.client.EncodeExecutionReport(ctx, &ccippb.EncodeExecutionReportRequest{
		Report: reportPB,
	})
	if err != nil {
		return nil, err
	}
	return resp.Report, nil
}

// GasPriceEstimator [github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GasPriceEstimator(ctx context.Context) (cciptypes.GasPriceEstimatorExec, error) {
	resp, err := o.client.GasPriceEstimator(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	// TODO BCF-3061: this works because the broker is shared and the id refers to a resource served by the broker
	gasEstimatorConn, err := o.b.Dial(uint32(resp.EstimatorServiceId))
	if err != nil {
		return nil, fmt.Errorf("failed to lookup gas estimator service for off ramp reader at %d: %w", resp.EstimatorServiceId, err)
	}
	// need to wrap grpc offRamp into the desired interface
	gasEstimator := NewExecGasEstimatorGRPCClient(gasEstimatorConn)
	// need to hydrate the gas price estimator from the server id
	return gasEstimator, nil
}

// GetExecutionState i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GetExecutionState(ctx context.Context, sequenceNumber uint64) (uint8, error) {
	resp, err := o.client.GetExecutionState(ctx, &ccippb.GetExecutionStateRequest{
		SeqNum: sequenceNumber,
	})
	if err != nil {
		return 0, err
	}
	return uint8(resp.ExecutionState), nil
}

// GetExecutionStateChangesBetweenSeqNums i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64, confirmations int) ([]cciptypes.ExecutionStateChangedWithTxMeta, error) {
	resp, err := o.client.GetExecutionStateChanges(ctx, &ccippb.GetExecutionStateChangesRequest{
		MinSeqNum:     seqNumMin,
		MaxSeqNum:     seqNumMax,
		Confirmations: int64(confirmations),
	})
	if err != nil {
		return nil, err
	}
	return executionStateChangedWithTxMetaSlice(resp.ExecutionStateChanges), nil
}

// GetSendersNonce i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) ListSenderNonces(ctx context.Context, senders []cciptypes.Address) (map[cciptypes.Address]uint64, error) {
	stringSenders := make([]string, len(senders))
	for i, s := range senders {
		stringSenders[i] = string(s)
	}

	resp, err := o.client.ListSenderNonces(ctx, &ccippb.ListSenderNoncesRequest{
		Senders: stringSenders,
	})
	if err != nil {
		return nil, err
	}
	return senderToNonceMapping(resp.GetNonceMapping()), nil
}

// GetSourceToDestTokensMapping i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GetSourceToDestTokensMapping(ctx context.Context) (map[cciptypes.Address]cciptypes.Address, error) {
	resp, err := o.client.GetSourceToDestTokensMapping(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return sourceToDestTokensMapping(resp.TokenMappings), nil
}

// GetStaticConfig i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GetStaticConfig(ctx context.Context) (cciptypes.OffRampStaticConfig, error) {
	resp, err := o.client.GetStaticConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.OffRampStaticConfig{}, err
	}
	return cciptypes.OffRampStaticConfig{
		CommitStore:         cciptypes.Address(resp.Config.CommitStore),
		ChainSelector:       resp.Config.ChainSelector,
		SourceChainSelector: resp.Config.SourceChainSelector,
		OnRamp:              cciptypes.Address(resp.Config.OnRamp),
		PrevOffRamp:         cciptypes.Address(resp.Config.PrevOffRamp),
		ArmProxy:            cciptypes.Address(resp.Config.ArmProxy),
	}, nil
}

// GetTokens i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GetTokens(ctx context.Context) (cciptypes.OffRampTokens, error) {
	resp, err := o.client.GetTokens(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.OffRampTokens{}, err
	}
	return offRampTokens(resp.Tokens), nil
}

// GetRouter i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) GetRouter(ctx context.Context) (cciptypes.Address, error) {
	resp, err := o.client.GetRouter(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.Address(""), err
	}
	return cciptypes.Address(resp.Router), nil
}

// OffchainConfig i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) OffchainConfig(ctx context.Context) (cciptypes.ExecOffchainConfig, error) {
	resp, err := o.client.OffchainConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.ExecOffchainConfig{}, err
	}
	return offChainConfig(resp.Config)
}

// OnchainConfig i[github.com/smartcontractkit/chainlink-common/pkg/types/ccip.OffRampReader]
func (o *OffRampReaderGRPCClient) OnchainConfig(ctx context.Context) (cciptypes.ExecOnchainConfig, error) {
	resp, err := o.client.OnchainConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return cciptypes.ExecOnchainConfig{}, err
	}
	return cciptypes.ExecOnchainConfig{
		PermissionLessExecutionThresholdSeconds: resp.Config.PermissionlessExecThresholdSeconds.AsDuration(),
		Router:                                  cciptypes.Address(resp.Config.Router),
	}, nil
}

// Server implementation of OffRampReader

// Address implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) Address(ctx context.Context, req *emptypb.Empty) (*ccippb.OffRampAddressResponse, error) {
	addr, err := o.impl.Address(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.OffRampAddressResponse{Address: string(addr)}, nil
}

// ChangeConfig implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) ChangeConfig(ctx context.Context, req *ccippb.ChangeConfigRequest) (*ccippb.ChangeConfigResponse, error) {
	onchainAddr, offchainAddr, err := o.impl.ChangeConfig(ctx, req.OnchainConfig, req.OffchainConfig)
	if err != nil {
		return nil, err
	}
	return &ccippb.ChangeConfigResponse{
		OnchainConfigAddress:  string(onchainAddr),
		OffchainConfigAddress: string(offchainAddr),
	}, nil
}

// Close implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) Close(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, services.MultiCloser(o.deps).Close()
}

// CurrentRateLimiterState implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) CurrentRateLimiterState(ctx context.Context, req *emptypb.Empty) (*ccippb.CurrentRateLimiterStateResponse, error) {
	state, err := o.impl.CurrentRateLimiterState(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.CurrentRateLimiterStateResponse{RateLimiter: tokenBucketRateLimitPB(state)}, nil
}

// DecodeExecutionReport implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) DecodeExecutionReport(ctx context.Context, req *ccippb.DecodeExecutionReportRequest) (*ccippb.DecodeExecutionReportResponse, error) {
	report, err := o.impl.DecodeExecutionReport(ctx, req.Report)
	if err != nil {
		return nil, err
	}
	return &ccippb.DecodeExecutionReportResponse{Report: executionReportPB(report)}, nil
}

// EncodeExecutionReport implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) EncodeExecutionReport(ctx context.Context, req *ccippb.EncodeExecutionReportRequest) (*ccippb.EncodeExecutionReportResponse, error) {
	report, err := execReport(req.Report)
	if err != nil {
		return nil, err
	}

	encoded, err := o.impl.EncodeExecutionReport(ctx, report)
	if err != nil {
		return nil, err
	}
	return &ccippb.EncodeExecutionReportResponse{Report: encoded}, nil
}

// GasPriceEstimator implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GasPriceEstimator(ctx context.Context, req *emptypb.Empty) (*ccippb.GasPriceEstimatorResponse, error) {
	return &ccippb.GasPriceEstimatorResponse{EstimatorServiceId: int32(o.gasEstimatorServerID)}, nil
}

// GetExecutionState implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GetExecutionState(ctx context.Context, req *ccippb.GetExecutionStateRequest) (*ccippb.GetExecutionStateResponse, error) {
	state, err := o.impl.GetExecutionState(ctx, req.SeqNum)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetExecutionStateResponse{ExecutionState: uint32(state)}, nil
}

// GetExecutionStateChanges implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GetExecutionStateChanges(ctx context.Context, req *ccippb.GetExecutionStateChangesRequest) (*ccippb.GetExecutionStateChangesResponse, error) {
	changes, err := o.impl.GetExecutionStateChangesBetweenSeqNums(ctx, req.MinSeqNum, req.MaxSeqNum, int(req.Confirmations))
	if err != nil {
		return nil, err
	}
	return &ccippb.GetExecutionStateChangesResponse{ExecutionStateChanges: executionStateChangedWithTxMetaSliceToPB(changes)}, nil
}

// GetSendersNonce implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) ListSenderNonces(ctx context.Context, req *ccippb.ListSenderNoncesRequest) (*ccippb.ListSenderNoncesResponse, error) {
	senders := make([]cciptypes.Address, len(req.Senders))
	for i, s := range req.Senders {
		senders[i] = cciptypes.Address(s)
	}

	resp, err := o.impl.ListSenderNonces(ctx, senders)
	if err != nil {
		return nil, err
	}
	return &ccippb.ListSenderNoncesResponse{NonceMapping: senderToNonceMappingToPB(resp)}, nil
}

// GetSourceToDestTokensMapping implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GetSourceToDestTokensMapping(ctx context.Context, req *emptypb.Empty) (*ccippb.GetSourceToDestTokensMappingResponse, error) {
	mapping, err := o.impl.GetSourceToDestTokensMapping(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetSourceToDestTokensMappingResponse{TokenMappings: sourceDestTokenMappingToPB(mapping)}, nil
}

// GetStaticConfig implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GetStaticConfig(ctx context.Context, req *emptypb.Empty) (*ccippb.GetStaticConfigResponse, error) {
	config, err := o.impl.GetStaticConfig(ctx)
	if err != nil {
		return nil, err
	}

	pbConfig := ccippb.OffRampStaticConfig{
		CommitStore:         string(config.CommitStore),
		ChainSelector:       config.ChainSelector,
		SourceChainSelector: config.SourceChainSelector,
		OnRamp:              string(config.OnRamp),
		PrevOffRamp:         string(config.PrevOffRamp),
		ArmProxy:            string(config.ArmProxy),
	}
	return &ccippb.GetStaticConfigResponse{Config: &pbConfig}, nil
}

// GetTokens implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GetTokens(ctx context.Context, req *emptypb.Empty) (*ccippb.GetTokensResponse, error) {
	tokens, err := o.impl.GetTokens(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetTokensResponse{Tokens: offRampTokensToPB(tokens)}, nil
}

// GetRouter implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) GetRouter(ctx context.Context, req *emptypb.Empty) (*ccippb.GetRouterResponse, error) {
	router, err := o.impl.GetRouter(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.GetRouterResponse{Router: string(router)}, nil
}

// OffchainConfig implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) OffchainConfig(ctx context.Context, req *emptypb.Empty) (*ccippb.OffchainConfigResponse, error) {
	config, err := o.impl.OffchainConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &ccippb.OffchainConfigResponse{Config: offChainConfigToPB(config)}, nil
}

// OnchainConfig implements ccippb.OffRampReaderServer.
func (o *OffRampReaderGRPCServer) OnchainConfig(ctx context.Context, req *emptypb.Empty) (*ccippb.OnchainConfigResponse, error) {
	config, err := o.impl.OnchainConfig(ctx)
	if err != nil {
		return nil, err
	}
	pbConfig := ccippb.ExecOnchainConfig{
		PermissionlessExecThresholdSeconds: durationpb.New(config.PermissionLessExecutionThresholdSeconds),
		Router:                             string(config.Router),
	}
	return &ccippb.OnchainConfigResponse{Config: &pbConfig}, nil
}

// AddDep adds a closer to the list of dependencies that will be closed when the server is closed.
func (o *OffRampReaderGRPCServer) AddDep(dep io.Closer) *OffRampReaderGRPCServer {
	o.deps = append(o.deps, dep)
	return o
}

// Conversion functions and helpers

func tokenBucketRateLimit(pb *ccippb.TokenPoolRateLimit) cciptypes.TokenBucketRateLimit {
	return cciptypes.TokenBucketRateLimit{
		Tokens:      pb.Tokens.Int(),
		LastUpdated: pb.LastUpdated,
		IsEnabled:   pb.IsEnabled,
		Capacity:    pb.Capacity.Int(),
		Rate:        pb.Rate.Int(),
	}
}

func tokenBucketRateLimitPB(state cciptypes.TokenBucketRateLimit) *ccippb.TokenPoolRateLimit {
	return &ccippb.TokenPoolRateLimit{
		Tokens:      pb.NewBigIntFromInt(state.Tokens),
		LastUpdated: state.LastUpdated,
		IsEnabled:   state.IsEnabled,
		Capacity:    pb.NewBigIntFromInt(state.Capacity),
		Rate:        pb.NewBigIntFromInt(state.Rate),
	}
}

func execReport(pb *ccippb.ExecutionReport) (cciptypes.ExecReport, error) {
	proofs, err := byte32Slice(pb.Proofs)
	if err != nil {
		return cciptypes.ExecReport{}, fmt.Errorf("execReport: invalid proofs: %w", err)
	}
	msgs, err := evm2EVMMessageSlice(pb.EvmToEvmMessages)
	if err != nil {
		return cciptypes.ExecReport{}, fmt.Errorf("execReport: invalid messages: %w", err)
	}

	return cciptypes.ExecReport{
		Messages:          msgs,
		OffchainTokenData: offchainTokenData(pb.OffchainTokenData),
		Proofs:            proofs,
		ProofFlagBits:     pb.ProofFlagBits.Int(),
	}, nil
}

func evm2EVMMessageSlice(in []*ccippb.EVM2EVMMessage) ([]cciptypes.EVM2EVMMessage, error) {
	out := make([]cciptypes.EVM2EVMMessage, len(in))
	for i, m := range in {
		decodedMsg, err := evm2EVMMessage(m)
		if err != nil {
			return nil, err
		}
		out[i] = decodedMsg
	}
	return out, nil
}

func offchainTokenData(in []*ccippb.TokenData) [][][]byte {
	out := make([][][]byte, len(in))
	for i, b := range in {
		out[i] = b.Data
	}
	return out
}

func byte32Slice(in [][]byte) ([][32]byte, error) {
	out := make([][32]byte, len(in))
	for i, b := range in {
		if len(b) != 32 {
			return nil, fmt.Errorf("byte32Slice: invalid length %d", len(b))
		}
		copy(out[i][:], b)
	}
	return out, nil
}

func executionReportPB(report cciptypes.ExecReport) *ccippb.ExecutionReport {
	return &ccippb.ExecutionReport{
		EvmToEvmMessages:  evm2EVMMessageSliceToPB(report.Messages),
		OffchainTokenData: offchainTokenDataToPB(report.OffchainTokenData),
		Proofs:            byte32SliceToPB(report.Proofs),
		ProofFlagBits:     pb.NewBigIntFromInt(report.ProofFlagBits),
	}
}

func evm2EVMMessageSliceToPB(in []cciptypes.EVM2EVMMessage) []*ccippb.EVM2EVMMessage {
	out := make([]*ccippb.EVM2EVMMessage, len(in))
	for i, m := range in {
		out[i] = evm2EVMMessageToPB(m)
	}
	return out
}

func offchainTokenDataToPB(in [][][]byte) []*ccippb.TokenData {
	out := make([]*ccippb.TokenData, len(in))
	for i, b := range in {
		out[i] = &ccippb.TokenData{Data: b}
	}
	return out
}

func byte32SliceToPB(in [][32]byte) [][]byte {
	out := make([][]byte, len(in))
	for i, b := range in {
		out[i] = make([]byte, 32)
		copy(out[i][:], b[:])
	}
	return out
}

func evm2EVMMessageToPB(m cciptypes.EVM2EVMMessage) *ccippb.EVM2EVMMessage {
	return &ccippb.EVM2EVMMessage{
		SequenceNumber:      m.SequenceNumber,
		GasLimit:            pb.NewBigIntFromInt(m.GasLimit),
		Nonce:               m.Nonce,
		MessageId:           m.MessageID[:],
		SourceChainSelector: m.SourceChainSelector,
		Sender:              string(m.Sender),
		Receiver:            string(m.Receiver),
		Strict:              m.Strict,
		FeeToken:            string(m.FeeToken),
		FeeTokenAmount:      pb.NewBigIntFromInt(m.FeeTokenAmount),
		Data:                m.Data,
		TokenAmounts:        tokenAmountSliceToPB(m.TokenAmounts),
		SourceTokenData:     m.SourceTokenData,
	}
}

func tokenAmountSliceToPB(tokenAmounts []cciptypes.TokenAmount) []*ccippb.TokenAmount {
	res := make([]*ccippb.TokenAmount, len(tokenAmounts))
	for i, t := range tokenAmounts {
		res[i] = &ccippb.TokenAmount{
			Token:  string(t.Token),
			Amount: pb.NewBigIntFromInt(t.Amount),
		}
	}
	return res
}

func executionStateChangedWithTxMetaSlice(in []*ccippb.ExecutionStateChangeWithTxMeta) []cciptypes.ExecutionStateChangedWithTxMeta {
	out := make([]cciptypes.ExecutionStateChangedWithTxMeta, len(in))
	for i, m := range in {
		out[i] = executionStateChangedWithTxMeta(m)
	}
	return out
}

func executionStateChangedWithTxMetaSliceToPB(in []cciptypes.ExecutionStateChangedWithTxMeta) []*ccippb.ExecutionStateChangeWithTxMeta {
	out := make([]*ccippb.ExecutionStateChangeWithTxMeta, len(in))
	for i, m := range in {
		out[i] = executionStateChangedWithTxMetaToPB(m)
	}
	return out
}

func executionStateChangedWithTxMeta(in *ccippb.ExecutionStateChangeWithTxMeta) cciptypes.ExecutionStateChangedWithTxMeta {
	return cciptypes.ExecutionStateChangedWithTxMeta{
		TxMeta: txMeta(in.TxMeta),
		ExecutionStateChanged: cciptypes.ExecutionStateChanged{
			SequenceNumber: in.ExecutionStateChange.SeqNum,
			Finalized:      in.ExecutionStateChange.Finalized,
		},
	}
}

func executionStateChangedWithTxMetaToPB(in cciptypes.ExecutionStateChangedWithTxMeta) *ccippb.ExecutionStateChangeWithTxMeta {
	return &ccippb.ExecutionStateChangeWithTxMeta{
		TxMeta: txMetaToPB(in.TxMeta),
		ExecutionStateChange: &ccippb.ExecutionStateChange{
			SeqNum:    in.ExecutionStateChanged.SequenceNumber,
			Finalized: in.ExecutionStateChanged.Finalized,
		},
	}
}
func txMetaToPB(in cciptypes.TxMeta) *ccippb.TxMeta {
	return &ccippb.TxMeta{
		BlockTimestampUnixMilli: in.BlockTimestampUnixMilli,
		BlockNumber:             in.BlockNumber,
		TxHash:                  in.TxHash,
		LogIndex:                in.LogIndex,
	}
}

func offRampTokens(in *ccippb.OffRampTokens) cciptypes.OffRampTokens {
	source := make([]cciptypes.Address, len(in.SourceTokens))
	for i, t := range in.SourceTokens {
		source[i] = cciptypes.Address(t)
	}
	dest := make([]cciptypes.Address, len(in.DestinationTokens))
	for i, t := range in.DestinationTokens {
		dest[i] = cciptypes.Address(t)
	}
	destPool := make(map[cciptypes.Address]cciptypes.Address)
	for k, v := range in.DestinationPool {
		destPool[cciptypes.Address(k)] = cciptypes.Address(v)
	}

	return cciptypes.OffRampTokens{
		SourceTokens:      source,
		DestinationTokens: dest,
		DestinationPool:   destPool,
	}
}

func offRampTokensToPB(in cciptypes.OffRampTokens) *ccippb.OffRampTokens {
	source := make([]string, len(in.SourceTokens))
	for i, t := range in.SourceTokens {
		source[i] = string(t)
	}
	dest := make([]string, len(in.DestinationTokens))
	for i, t := range in.DestinationTokens {
		dest[i] = string(t)
	}
	destPool := make(map[string]string)
	for k, v := range in.DestinationPool {
		destPool[string(k)] = string(v)
	}

	return &ccippb.OffRampTokens{
		SourceTokens:      source,
		DestinationTokens: dest,
		DestinationPool:   destPool,
	}
}

func offChainConfig(in *ccippb.ExecOffchainConfig) (cciptypes.ExecOffchainConfig, error) {
	cachedExpiry, err := config.NewDuration(in.InflightCacheExpiry.AsDuration())
	if err != nil {
		return cciptypes.ExecOffchainConfig{}, fmt.Errorf("offChainConfig: invalid InflightCacheExpiry: %w", err)
	}
	rootSnoozeTime, err := config.NewDuration(in.RootSnoozeTime.AsDuration())
	if err != nil {
		return cciptypes.ExecOffchainConfig{}, fmt.Errorf("offChainConfig: invalid RootSnoozeTime: %w", err)
	}

	return cciptypes.ExecOffchainConfig{
		DestOptimisticConfirmations: in.DestOptimisticConfirmations,
		BatchGasLimit:               in.BatchGasLimit,
		RelativeBoostPerWaitHour:    in.RelativeBoostPerWaitHour,
		InflightCacheExpiry:         cachedExpiry,
		RootSnoozeTime:              rootSnoozeTime,
	}, nil
}

func offChainConfigToPB(in cciptypes.ExecOffchainConfig) *ccippb.ExecOffchainConfig {
	return &ccippb.ExecOffchainConfig{
		DestOptimisticConfirmations: in.DestOptimisticConfirmations,
		BatchGasLimit:               in.BatchGasLimit,
		RelativeBoostPerWaitHour:    in.RelativeBoostPerWaitHour,
		InflightCacheExpiry:         durationpb.New(in.InflightCacheExpiry.Duration()),
		RootSnoozeTime:              durationpb.New(in.RootSnoozeTime.Duration()),
	}
}

func sourceToDestTokensMapping(in map[string]string) map[cciptypes.Address]cciptypes.Address {
	out := make(map[cciptypes.Address]cciptypes.Address)
	for k, v := range in {
		out[cciptypes.Address(k)] = cciptypes.Address(v)
	}
	return out
}

func sourceDestTokenMappingToPB(in map[cciptypes.Address]cciptypes.Address) map[string]string {
	out := make(map[string]string)
	for k, v := range in {
		out[string(k)] = string(v)
	}
	return out
}

func senderToNonceMapping(in map[string]uint64) map[cciptypes.Address]uint64 {
	out := make(map[cciptypes.Address]uint64, len(in))
	for k, v := range in {
		out[cciptypes.Address(k)] = v
	}
	return out
}

func senderToNonceMappingToPB(in map[cciptypes.Address]uint64) map[string]uint64 {
	out := make(map[string]uint64, len(in))
	for k, v := range in {
		out[string(k)] = v
	}
	return out
}
