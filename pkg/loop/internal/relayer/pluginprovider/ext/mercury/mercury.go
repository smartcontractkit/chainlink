package mercury

import (
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	mercury_pb "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb/mercury"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2"

	mercury_v1_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v1"
	mercury_v2_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v2"
	mercury_v3_internal "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
	mercury_v1 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v1"
	mercury_v2 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v2"
	mercury_v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
)

var (
	_ types.MercuryProvider = (*ProviderClient)(nil)
	// in practice, inherited from pluginProviderClient.
	_ goplugin.GRPCClientConn = (*ProviderClient)(nil)
)

type ProviderClient struct {
	*ocr2.PluginProviderClient
	reportCodecV3      mercury_v3.ReportCodec
	reportCodecV2      mercury_v2.ReportCodec
	reportCodecV1      mercury_v1.ReportCodec
	onchainConfigCodec mercury.OnchainConfigCodec
	serverFetcher      mercury.ServerFetcher
	chainReader        types.ContractReader
	mercuryChainReader mercury.ChainReader
}

func NewProviderClient(b *net.BrokerExt, cc grpc.ClientConnInterface) *ProviderClient {
	m := &ProviderClient{PluginProviderClient: ocr2.NewPluginProviderClient(b.WithName("MercuryProviderClient"), cc)}

	m.reportCodecV1 = newReportCodecV1Client(mercury_v1_internal.NewReportCodecClient(cc))
	m.reportCodecV2 = newReportCodecV2Client(mercury_v2_internal.NewReportCodecClient(cc))
	m.reportCodecV3 = newReportCodecV3Client(mercury_v3_internal.NewReportCodecClient(cc))

	m.onchainConfigCodec = newOnchainConfigCodecClient(cc)
	m.serverFetcher = newServerFetcherClient(cc)
	m.mercuryChainReader = newChainReaderClient(cc)

	m.chainReader = chainreader.NewClient(b, cc)
	return m
}

func (m *ProviderClient) ReportCodecV3() mercury_v3.ReportCodec {
	return m.reportCodecV3
}

func (m *ProviderClient) ReportCodecV2() mercury_v2.ReportCodec {
	return m.reportCodecV2
}

func (m *ProviderClient) ReportCodecV1() mercury_v1.ReportCodec {
	return m.reportCodecV1
}

func (m *ProviderClient) OnchainConfigCodec() mercury.OnchainConfigCodec {
	return m.onchainConfigCodec
}

func (m *ProviderClient) ChainReader() types.ContractReader {
	return m.chainReader
}

func (m *ProviderClient) MercuryChainReader() mercury.ChainReader {
	return m.mercuryChainReader
}

func (m *ProviderClient) MercuryServerFetcher() mercury.ServerFetcher {
	return m.serverFetcher
}

func registerVersionAgnosticServices(s *grpc.Server, provider types.MercuryProvider) {
	mercury_pb.RegisterOnchainConfigCodecServer(s, newOnchainConfigCodecServer(provider.OnchainConfigCodec()))
	mercury_pb.RegisterServerFetcherServer(s, newServerFetcherServer(provider.MercuryServerFetcher()))
	mercury_pb.RegisterMercuryChainReaderServer(s, newChainReaderServer(provider.MercuryChainReader()))
}

// RegisterProviderServices registers the Mercury services with the given gRPC server.
// It registers all versions of the report codec service and the version-agnostic services.
// It should used by default, unless you need to register a specific version of the report codec service.
func RegisterProviderServices(s *grpc.Server, provider types.MercuryProvider) {
	registerVersionAgnosticServices(s, provider)

	mercury_pb.RegisterReportCodecV1Server(s, newReportCodecV1Server(s, provider.ReportCodecV1()))
	mercury_pb.RegisterReportCodecV2Server(s, newReportCodecV2Server(s, provider.ReportCodecV2()))
	mercury_pb.RegisterReportCodecV3Server(s, newReportCodecV3Server(s, provider.ReportCodecV3()))
}

// RegisterProviderServicesV1 registers the Mercury services with the given gRPC server.
// It registers only the v1 version of the report codec service and the version-agnostic services.
// It should be used when you only need to register the v1 version of the report codec service, and will cause
// gRPC to return an unimplemented error for the v2 and v3 versions.
func RegisterProviderServicesV1(s *grpc.Server, provider types.MercuryProvider) {
	registerVersionAgnosticServices(s, provider)

	mercury_pb.RegisterReportCodecV1Server(s, newReportCodecV1Server(s, provider.ReportCodecV1()))
	mercury_pb.RegisterReportCodecV2Server(s, mercury_pb.UnimplementedReportCodecV2Server{})
	mercury_pb.RegisterReportCodecV3Server(s, mercury_pb.UnimplementedReportCodecV3Server{})
}

// RegisterProviderServicesV2 registers the Mercury services with the given gRPC server.
// It registers only the v2 version of the report codec service and the version-agnostic services.
// It should be used when you only need to register the v2 version of the report codec service, and will cause
// gRPC to return an unimplemented error for the v1 and v3 versions.
func RegisterProviderServicesV2(s *grpc.Server, provider types.MercuryProvider) {
	registerVersionAgnosticServices(s, provider)

	mercury_pb.RegisterReportCodecV2Server(s, newReportCodecV2Server(s, provider.ReportCodecV2()))
	mercury_pb.RegisterReportCodecV1Server(s, mercury_pb.UnimplementedReportCodecV1Server{})
	mercury_pb.RegisterReportCodecV3Server(s, mercury_pb.UnimplementedReportCodecV3Server{})
}

// RegisterProviderServicesV3 registers the Mercury services with the given gRPC server.
// It registers only the v3 version of the report codec service and the version-agnostic services.
// It should be used when you only need to register the v3 version of the report codec service, and will cause
// gRPC to return an unimplemented error for the v1 and v2 versions.
func RegisterProviderServicesV3(s *grpc.Server, provider types.MercuryProvider) {
	registerVersionAgnosticServices(s, provider)

	mercury_pb.RegisterReportCodecV3Server(s, newReportCodecV3Server(s, provider.ReportCodecV3()))
	mercury_pb.RegisterReportCodecV1Server(s, mercury_pb.UnimplementedReportCodecV1Server{})
	mercury_pb.RegisterReportCodecV2Server(s, mercury_pb.UnimplementedReportCodecV2Server{})
}
