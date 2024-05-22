package median

import (
	"context"

	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/reportingplugin/ocr2"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	medianprovider "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ext/median"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
)

var _ core.PluginMedian = (*PluginMedianClient)(nil)

type PluginMedianClient struct {
	*goplugin.PluginClient
	*goplugin.ServiceClient

	median pb.PluginMedianClient
}

func NewPluginMedianClient(broker net.Broker, brokerCfg net.BrokerConfig, conn *grpc.ClientConn) *PluginMedianClient {
	brokerCfg.Logger = logger.Named(brokerCfg.Logger, "PluginMedianClient")
	pc := goplugin.NewPluginClient(broker, brokerCfg, conn)
	return &PluginMedianClient{PluginClient: pc, median: pb.NewPluginMedianClient(pc), ServiceClient: goplugin.NewServiceClient(pc.BrokerExt, pc)}
}

func (m *PluginMedianClient) NewMedianFactory(ctx context.Context, provider types.MedianProvider, dataSource, juelsPerFeeCoin, gasPriceSubunits median.DataSource, errorLog core.ErrorLog) (types.ReportingPluginFactory, error) {
	cc := m.NewClientConn("MedianPluginFactory", func(ctx context.Context) (id uint32, deps net.Resources, err error) {
		dataSourceID, dsRes, err := m.ServeNew("DataSource", func(s *grpc.Server) {
			pb.RegisterDataSourceServer(s, newDataSourceServer(dataSource))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(dsRes)

		juelsPerFeeCoinDataSourceID, juelsPerFeeCoinDataSourceRes, err := m.ServeNew("JuelsPerFeeCoinDataSource", func(s *grpc.Server) {
			pb.RegisterDataSourceServer(s, newDataSourceServer(juelsPerFeeCoin))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(juelsPerFeeCoinDataSourceRes)

		gasPriceSubunitsDataSourceID, gasPriceSubunitsDataSourceRes, err := m.ServeNew("GasPriceSubunitsDataSource", func(s *grpc.Server) {
			pb.RegisterDataSourceServer(s, newDataSourceServer(gasPriceSubunits))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(gasPriceSubunitsDataSourceRes)

		var (
			providerID  uint32
			providerRes net.Resource
		)
		if grpcProvider, ok := provider.(goplugin.GRPCClientConn); ok {
			providerID, providerRes, err = m.Serve("MedianProvider", proxy.NewProxy(grpcProvider.ClientConn()))
		} else {
			providerID, providerRes, err = m.ServeNew("MedianProvider", func(s *grpc.Server) {
				medianprovider.RegisterProviderServices(s, provider)
			})
		}
		if err != nil {
			return 0, nil, err
		}
		deps.Add(providerRes)

		errorLogID, errorLogRes, err := m.ServeNew("ErrorLog", func(s *grpc.Server) {
			pb.RegisterErrorLogServer(s, errorlog.NewServer(errorLog))
		})
		if err != nil {
			return 0, nil, err
		}
		deps.Add(errorLogRes)

		reply, err := m.median.NewMedianFactory(ctx, &pb.NewMedianFactoryRequest{
			MedianProviderID:             providerID,
			DataSourceID:                 dataSourceID,
			JuelsPerFeeCoinDataSourceID:  juelsPerFeeCoinDataSourceID,
			GasPriceSubunitsDataSourceID: gasPriceSubunitsDataSourceID,
			ErrorLogID:                   errorLogID,
		})
		if err != nil {
			return 0, nil, err
		}
		return reply.ReportingPluginFactoryID, nil, nil
	})
	return ocr2.NewReportingPluginFactoryClient(m.PluginClient.BrokerExt, cc), nil
}

var _ pb.PluginMedianServer = (*pluginMedianServer)(nil)

type pluginMedianServer struct {
	pb.UnimplementedPluginMedianServer

	*net.BrokerExt
	impl core.PluginMedian
}

func RegisterPluginMedianServer(server *grpc.Server, broker net.Broker, brokerCfg net.BrokerConfig, impl core.PluginMedian) error {
	pb.RegisterPluginMedianServer(server, newPluginMedianServer(&net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}, impl))
	return nil
}

func newPluginMedianServer(b *net.BrokerExt, mp core.PluginMedian) *pluginMedianServer {
	return &pluginMedianServer{BrokerExt: b.WithName("PluginMedian"), impl: mp}
}

func (m *pluginMedianServer) NewMedianFactory(ctx context.Context, request *pb.NewMedianFactoryRequest) (*pb.NewMedianFactoryReply, error) {
	dsConn, err := m.Dial(request.DataSourceID)
	if err != nil {
		return nil, net.ErrConnDial{Name: "DataSource", ID: request.DataSourceID, Err: err}
	}
	dsRes := net.Resource{Closer: dsConn, Name: "DataSource"}
	dataSource := newDataSourceClient(dsConn)

	juelsConn, err := m.Dial(request.JuelsPerFeeCoinDataSourceID)
	if err != nil {
		m.CloseAll(dsRes)
		return nil, net.ErrConnDial{Name: "JuelsPerFeeCoinDataSource", ID: request.JuelsPerFeeCoinDataSourceID, Err: err}
	}
	juelsRes := net.Resource{Closer: juelsConn, Name: "JuelsPerFeeCoinDataSource"}
	juelsPerFeeCoin := newDataSourceClient(juelsConn)

	gasPriceSubunitsConn, err := m.Dial(request.GasPriceSubunitsDataSourceID)
	if err != nil {
		m.CloseAll(dsRes, juelsRes)
		return nil, net.ErrConnDial{Name: "GasPriceSubunitsDataSource", ID: request.GasPriceSubunitsDataSourceID, Err: err}
	}
	gasPriceSubunitsRes := net.Resource{Closer: gasPriceSubunitsConn, Name: "GasPriceSubunitsDataSource"}
	gasPriceSubunits := newDataSourceClient(gasPriceSubunitsConn)

	providerConn, err := m.Dial(request.MedianProviderID)
	if err != nil {
		m.CloseAll(dsRes, juelsRes, gasPriceSubunitsRes)
		return nil, net.ErrConnDial{Name: "MedianProvider", ID: request.MedianProviderID, Err: err}
	}
	providerRes := net.Resource{Closer: providerConn, Name: "MedianProvider"}
	provider := medianprovider.NewProviderClient(m.BrokerExt, providerConn)

	errorLogConn, err := m.Dial(request.ErrorLogID)
	if err != nil {
		m.CloseAll(dsRes, juelsRes, gasPriceSubunitsRes, providerRes)
		return nil, net.ErrConnDial{Name: "ErrorLog", ID: request.ErrorLogID, Err: err}
	}
	errorLogRes := net.Resource{Closer: errorLogConn, Name: "ErrorLog"}
	errorLog := errorlog.NewClient(errorLogConn)

	factory, err := m.impl.NewMedianFactory(ctx, provider, dataSource, juelsPerFeeCoin, gasPriceSubunits, errorLog)
	if err != nil {
		m.CloseAll(dsRes, juelsRes, gasPriceSubunitsRes, providerRes, errorLogRes)
		return nil, err
	}

	id, _, err := m.ServeNew("ReportingPluginProvider", func(s *grpc.Server) {
		pb.RegisterServiceServer(s, &goplugin.ServiceServer{Srv: factory})
		pb.RegisterReportingPluginFactoryServer(s, ocr2.NewReportingPluginFactoryServer(factory, m.BrokerExt))
	}, dsRes, juelsRes, gasPriceSubunitsRes, providerRes, errorLogRes)
	if err != nil {
		return nil, err
	}

	return &pb.NewMedianFactoryReply{ReportingPluginFactoryID: id}, nil
}
