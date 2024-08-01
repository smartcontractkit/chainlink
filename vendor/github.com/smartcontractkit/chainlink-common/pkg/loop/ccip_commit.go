package loop

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/reportingplugin/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// CCIPCommitLOOPName is the name for [types.CCIPCommitFactoryGenerator]/[NewCommitLOOP].
const CCIPCommitLOOPName = "ccip_commit"

func PluginCCIPCommitHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		MagicCookieKey:   "CL_PLUGIN_CCIP_COMMIT_MAGIC_COOKIE",
		MagicCookieValue: "5a2d1527-6c0f-4c7e-8c96-00aa4bececd2",
	}
}

type CommitLoop struct {
	plugin.NetRPCUnsupportedPlugin

	BrokerConfig

	PluginServer types.CCIPCommitFactoryGenerator

	pluginClient *ccip.CommitLOOPClient
}

func (p *CommitLoop) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	return ccip.RegisterCommitLOOPServer(server, broker, p.BrokerConfig, p.PluginServer)
}

// GRPCClient implements [plugin.GRPCPlugin] and returns the pluginClient [types.CCIPCommitFactoryGenerator], updated with the new broker and conn.
func (p *CommitLoop) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	if p.pluginClient == nil {
		p.pluginClient = ccip.NewCommitLOOPClient(broker, p.BrokerConfig, conn)
	} else {
		p.pluginClient.Refresh(broker, conn)
	}

	return types.CCIPCommitFactoryGenerator(p.pluginClient), nil
}

func (p *CommitLoop) ClientConfig() *plugin.ClientConfig {
	clientConfig := &plugin.ClientConfig{
		HandshakeConfig: PluginCCIPCommitHandshakeConfig(),
		Plugins:         map[string]plugin.Plugin{CCIPCommitLOOPName: p},
	}
	return ManagedGRPCClientConfig(clientConfig, p.BrokerConfig)
}

var _ ocrtypes.ReportingPluginFactory = (*CommitFactoryService)(nil)

// CommitFactoryService is a [types.Service] that maintains an internal [types.CCIPCommitFactoryGenerator].
type CommitFactoryService struct {
	goplugin.PluginService[*CommitLoop, types.ReportingPluginFactory]
}

// NewCommitService returns a new [*CommitFactoryService].
// cmd must return a new exec.Cmd each time it is called.
func NewCommitService(lggr logger.Logger, grpcOpts GRPCOpts, cmd func() *exec.Cmd, provider types.CCIPCommitProvider) *CommitFactoryService {
	newService := func(ctx context.Context, instance any) (types.ReportingPluginFactory, error) {
		plug, ok := instance.(types.CCIPCommitFactoryGenerator)
		if !ok {
			return nil, fmt.Errorf("expected CCIPCommitFactoryGenerator but got %T", instance)
		}
		return plug.NewCommitFactory(ctx, provider)
	}
	stopCh := make(chan struct{})
	lggr = logger.Named(lggr, "CCIPCommitService")
	var cfs CommitFactoryService
	broker := BrokerConfig{StopCh: stopCh, Logger: lggr, GRPCOpts: grpcOpts}
	cfs.Init(CCIPCommitLOOPName, &CommitLoop{BrokerConfig: broker}, newService, lggr, cmd, stopCh)
	return &cfs
}

func (m *CommitFactoryService) NewReportingPlugin(config ocrtypes.ReportingPluginConfig) (ocrtypes.ReportingPlugin, ocrtypes.ReportingPluginInfo, error) {
	if err := m.Wait(); err != nil {
		return nil, ocrtypes.ReportingPluginInfo{}, err
	}
	return m.Service.NewReportingPlugin(config)
}
