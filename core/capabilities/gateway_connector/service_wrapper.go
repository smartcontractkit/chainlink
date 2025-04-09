package gatewayconnector

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
)

type Keystore interface {
	keys.AddressChecker
	keys.MessageSigner
}

type ServiceWrapper struct {
	services.StateMachine
	stopCh services.StopChan

	config    config.GatewayConnector
	keystore  Keystore
	connector connector.GatewayConnector
	lggr      logger.Logger
	clock     clockwork.Clock
}

func translateConfigs(f config.GatewayConnector) connector.ConnectorConfig {
	r := connector.ConnectorConfig{}
	r.NodeAddress = f.NodeAddress()
	r.DonId = f.DonID()

	if len(f.Gateways()) != 0 {
		r.Gateways = make([]connector.ConnectorGatewayConfig, len(f.Gateways()))
		for index, element := range f.Gateways() {
			r.Gateways[index] = connector.ConnectorGatewayConfig{Id: element.ID(), URL: element.URL()}
		}
	}

	r.WsClientConfig = network.WebSocketClientConfig{HandshakeTimeoutMillis: f.WSHandshakeTimeoutMillis()}
	r.AuthMinChallengeLen = f.AuthMinChallengeLen()
	r.AuthTimestampToleranceSec = f.AuthTimestampToleranceSec()
	return r
}

// NOTE: this wrapper is needed to make sure that our services are started after Keystore.
func NewGatewayConnectorServiceWrapper(config config.GatewayConnector, keystore Keystore, clock clockwork.Clock, lggr logger.Logger) *ServiceWrapper {
	return &ServiceWrapper{
		stopCh:   make(services.StopChan),
		config:   config,
		keystore: keystore,
		clock:    clock,
		lggr:     lggr.Named("GatewayConnectorServiceWrapper"),
	}
}

func (e *ServiceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		conf := e.config
		nodeAddress := conf.NodeAddress()
		configuredNodeAddress := common.HexToAddress(nodeAddress)
		err := e.keystore.CheckEnabled(ctx, configuredNodeAddress)
		if err != nil {
			return err
		}

		translated := translateConfigs(conf)
		e.connector, err = connector.NewGatewayConnector(&translated, e, e.clock, e.lggr)
		if err != nil {
			return err
		}
		return e.connector.Start(ctx)
	})
}

func (e *ServiceWrapper) Sign(data ...[]byte) ([]byte, error) {
	ctx, cancel := e.stopCh.NewCtx()
	defer cancel()
	account := common.HexToAddress(e.config.NodeAddress())
	return e.keystore.SignMessage(ctx, account, gwcommon.Flatten(data...))
}

func (e *ServiceWrapper) Close() error {
	return e.StopOnce("GatewayConnectorServiceWrapper", func() (err error) {
		close(e.stopCh)
		return e.connector.Close()
	})
}

func (e *ServiceWrapper) HealthReport() map[string]error {
	return map[string]error{e.Name(): e.Healthy()}
}

func (e *ServiceWrapper) Name() string {
	return e.lggr.Name()
}

func (e *ServiceWrapper) GetGatewayConnector() connector.GatewayConnector {
	return e.connector
}
