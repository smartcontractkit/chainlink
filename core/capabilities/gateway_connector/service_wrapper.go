package gateway_connector

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	gwConfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/network"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

type serviceWrapper struct {
	services.StateMachine

	config    *config.GatewayConnector
	keystore  keystore.Eth
	connector connector.GatewayConnector
	lggr      logger.Logger
}

func translateConfigs(f config.GatewayConnector) connector.ConnectorConfig {
	r := connector.ConnectorConfig{}
	if f.NodeAddress != nil {
		r.NodeAddress = f.NodeAddress()
	}

	if f.DonID != nil {
		r.DonId = f.DonID()
	}

	if f.Gateways != nil {
		r.Gateways = make([]connector.ConnectorGatewayConfig, len(r.Gateways))
		for index, element := range f.Gateways() {
			r.Gateways[index] = connector.ConnectorGatewayConfig{Id: element.ID(), URL: element.URL()}
		}
	}

	if !reflect.ValueOf(f.WSHandshakeTimeoutMillis).IsZero() {
		r.WsClientConfig = network.WebSocketClientConfig{HandshakeTimeoutMillis: f.WSHandshakeTimeoutMillis()}
	}

	if f.AuthMinChallengeLen != nil {
		r.AuthMinChallengeLen = f.AuthMinChallengeLen()
	}

	if f.AuthTimestampToleranceSec != nil {
		r.AuthTimestampToleranceSec = f.AuthTimestampToleranceSec()
	}
	return r
}

// NOTE: this wrapper is needed to make sure that our services are started after Keystore.
func NewGatewayConnectorServiceWrapper(config *gwConfig.GatewayConnector, keystore keystore.Eth, lggr logger.Logger) *serviceWrapper {
	return &serviceWrapper{
		config:   config,
		keystore: keystore,
		lggr:     lggr,
	}
}

func (e *serviceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		conf := *e.config
		e.lggr.Infow("Starting GatewayConnectorServiceWrapper", "chainID", conf.ChainIDForNodeKey())
		chainId, _ := new(big.Int).SetString(conf.ChainIDForNodeKey(), 0)
		enabledKeys, err := e.keystore.EnabledKeysForChain(ctx, chainId)
		if err != nil {
			return err
		}
		if len(enabledKeys) == 0 {
			return errors.New("no available keys found")
		}
		configuredNodeAddress := common.HexToAddress(conf.NodeAddress())
		idx := slices.IndexFunc(enabledKeys, func(key ethkey.KeyV2) bool { return key.Address == configuredNodeAddress })
		if idx == -1 {
			return errors.New("key for configured node address not found")
		}
		signerKey := enabledKeys[idx].ToEcdsaPrivKey()
		if enabledKeys[idx].ID() != conf.NodeAddress() {
			return errors.New("node address mismatch")
		}
		// cannot assign to conf.NodeAddress (neither addressable nor a map index expression
		// I think I'm still confused about mixing interfaces and structures, and how to use a setter on an interface.
		conf.NodeAddress = enabledKeys[idx].ID()

		handler, err := NewWorkflowConnectorHandler(e.config, signerKey, e.lggr)
		if err != nil {
			return err
		}
		translated := translateConfigs(conf)
		e.connector, err = connector.NewGatewayConnector(&translated, handler, handler, clockwork.NewRealClock(), e.lggr)
		if err != nil {
			return err
		}
		handler.SetConnector(e.connector)

		return e.connector.Start(ctx)
	})
}

func (e *serviceWrapper) Close() error {
	return e.StopOnce("WorkflowConnectorHandler", func() (err error) {
		return e.connector.Close()
	})
}

func (e *serviceWrapper) Ready() error {
	return nil
}

func (e *serviceWrapper) HealthReport() map[string]error {
	return nil
}

func (e *serviceWrapper) Name() string {
	return "GatewayConnectorServiceWrapper"
}
