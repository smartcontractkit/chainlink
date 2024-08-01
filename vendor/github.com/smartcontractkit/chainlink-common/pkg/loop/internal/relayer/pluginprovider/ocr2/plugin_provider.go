package ocr2

import (
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/goplugin"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/net"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type PluginProviderClient struct {
	*ConfigProviderClient
	contractTransmitter libocr.ContractTransmitter
	chainReader         types.ContractReader
	codec               types.Codec
}

var _ types.PluginProvider = (*PluginProviderClient)(nil)

// in practice, inherited from configProviderClient.
var _ goplugin.GRPCClientConn = (*PluginProviderClient)(nil)

func NewPluginProviderClient(b *net.BrokerExt, cc grpc.ClientConnInterface) *PluginProviderClient {
	p := &PluginProviderClient{ConfigProviderClient: NewConfigProviderClient(b.WithName("PluginProviderClient"), cc)}
	p.contractTransmitter = &contractTransmitterClient{b, pb.NewContractTransmitterClient(cc)}
	p.chainReader = chainreader.NewClient(b, cc)
	p.codec = chainreader.NewCodecClient(b, cc)
	return p
}

func (p *PluginProviderClient) ContractTransmitter() libocr.ContractTransmitter {
	return p.contractTransmitter
}

func (p *PluginProviderClient) ChainReader() types.ContractReader {
	return p.chainReader
}

func (p *PluginProviderClient) Codec() types.Codec {
	return p.codec
}

type PluginProviderServer struct{}

func (p PluginProviderServer) ConnToProvider(conn grpc.ClientConnInterface, broker net.Broker, brokerCfg net.BrokerConfig) types.PluginProvider {
	be := &net.BrokerExt{Broker: broker, BrokerConfig: brokerCfg}
	return NewPluginProviderClient(be, conn)
}
