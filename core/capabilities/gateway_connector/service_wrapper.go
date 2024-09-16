package gatewayconnector

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	gwcommon "github.com/smartcontractkit/chainlink/v2/core/services/gateway/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

type ServiceWrapper struct {
	services.StateMachine

	config    config.GatewayConnector
	keystore  keystore.Eth
	connector connector.GatewayConnector
	signerKey *ecdsa.PrivateKey
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
func NewGatewayConnectorServiceWrapper(config config.GatewayConnector, keystore keystore.Eth, clock clockwork.Clock, lggr logger.Logger) *ServiceWrapper {
	return &ServiceWrapper{
		config:   config,
		keystore: keystore,
		clock:    clock,
		lggr:     lggr,
	}
}

func (e *ServiceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		conf := e.config
		nodeAddress := conf.NodeAddress()
		chainID, _ := new(big.Int).SetString(conf.ChainIDForNodeKey(), 0)
		enabledKeys, err := e.keystore.EnabledKeysForChain(ctx, chainID)
		if err != nil {
			return err
		}
		if len(enabledKeys) == 0 {
			return errors.New("no available keys found")
		}
		configuredNodeAddress := common.HexToAddress(nodeAddress)
		idx := slices.IndexFunc(enabledKeys, func(key ethkey.KeyV2) bool { return key.Address == configuredNodeAddress })

		if idx == -1 {
			return errors.New("key for configured node address not found")
		}
		e.signerKey = enabledKeys[idx].ToEcdsaPrivKey()
		if enabledKeys[idx].ID() != nodeAddress {
			return errors.New("node address mismatch")
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
	return gwcommon.SignData(e.signerKey, data...)
}

func (e *ServiceWrapper) Close() error {
	return e.StopOnce("GatewayConnectorServiceWrapper", func() (err error) {
		return e.connector.Close()
	})
}

func (e *ServiceWrapper) HealthReport() map[string]error {
	return map[string]error{e.Name(): e.Healthy()}
}

func (e *ServiceWrapper) Name() string {
	return "GatewayConnectorServiceWrapper"
}

func (e *ServiceWrapper) GetGatewayConnector() connector.GatewayConnector {
	return e.connector
}
