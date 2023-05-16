package loop

import (
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
)

type pluginClient struct {
	atomicBroker
	atomicClient
	*brokerExt
}

func newPluginClient(stopCh <-chan struct{}, lggr logger.Logger, broker *plugin.GRPCBroker, conn *grpc.ClientConn) *pluginClient {
	var pc pluginClient
	pc.brokerExt = &brokerExt{stopCh: stopCh, lggr: lggr, broker: &pc.atomicBroker}
	pc.refresh(broker, conn)
	return &pc
}

func (p *pluginClient) refresh(broker *plugin.GRPCBroker, conn *grpc.ClientConn) {
	p.atomicBroker.store(broker)
	p.atomicClient.store(conn)
	p.lggr.Debugw("Refreshed pluginClient connection", "state", conn.GetState())
}
