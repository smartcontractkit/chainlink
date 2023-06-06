package internal

import (
	"google.golang.org/grpc"
)

type pluginClient struct {
	atomicBroker
	atomicClient
	*brokerExt
}

func newPluginClient(broker Broker, brokerCfg BrokerConfig, conn *grpc.ClientConn) *pluginClient {
	var pc pluginClient
	pc.brokerExt = &brokerExt{&pc.atomicBroker, brokerCfg}
	pc.Refresh(broker, conn)
	return &pc
}

func (p *pluginClient) Refresh(broker Broker, conn *grpc.ClientConn) {
	p.atomicBroker.store(broker)
	p.atomicClient.store(conn)
	p.Logger.Debugw("Refreshed pluginClient connection", "state", conn.GetState())
}
