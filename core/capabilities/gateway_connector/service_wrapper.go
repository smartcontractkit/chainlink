package gateway_connector

import (
	"context"
	"errors"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/connector"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type serviceWrapper struct {
	services.StateMachine

	config    *WorkflowConnectorConfig
	keystore  keystore.Eth
	connector connector.GatewayConnector
	lggr      logger.Logger
}

// NOTE: this wrapper is needed to make sure that our services are started after Keystore.
func NewGatewayConnectorServiceWrapper(config *WorkflowConnectorConfig, keystore keystore.Eth, lggr logger.Logger) *serviceWrapper {
	return &serviceWrapper{
		config:   config,
		keystore: keystore,
		lggr:     lggr,
	}
}

func (e *serviceWrapper) Start(ctx context.Context) error {
	return e.StartOnce("GatewayConnectorServiceWrapper", func() error {
		// Extract default Eth key to use for Gateway auth.
		// TODO: handle multiple keys and allow for configuration.
		e.lggr.Infow("Starting GatewayConnectorServiceWrapper", "chainID", e.config.ChainIDForNodeKey)
		enabledKeys, err := e.keystore.EnabledKeysForChain(ctx, e.config.ChainIDForNodeKey)
		if err != nil {
			return err
		}
		if len(enabledKeys) == 0 {
			return errors.New("no available keys found")
		}
		signerKey := enabledKeys[0].ToEcdsaPrivKey()
		e.config.GatewayConnectorConfig.NodeAddress = enabledKeys[0].ID()

		handler, err := NewWorkflowConnectorHandler(e.config, signerKey, e.lggr)
		if err != nil {
			return err
		}
		e.connector, err = connector.NewGatewayConnector(e.config.GatewayConnectorConfig, handler, handler, clockwork.NewRealClock(), e.lggr)
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
