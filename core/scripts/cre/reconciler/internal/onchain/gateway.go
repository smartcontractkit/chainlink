package onchain

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	griddleinfra "github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

// storeGatewayConnectors persists each gateway connector's info, keyed by the
// gateway's own DON — matched by NodeUUID against topology.DonsMetadata
// (not by position against cv.FindGatewayNodes(), which is fragile once
// there's more than one gateway and nothing guarantees matching order).
func (d *Deployer) storeGatewayConnectors(cv *domain.ChartValues, state *domain.StateFile, topology *cre.Topology) {
	state.GatewayConnectors = nil
	if topology.GatewayConnectors == nil {
		return
	}

	for _, config := range topology.GatewayConnectors.Configurations {
		if config == nil || config.GatewayConfiguration == nil {
			continue
		}
		gc := config.GatewayConfiguration

		donName, nodeName := gatewayOwnDONAndNodeForUUID(topology, cv, gc.NodeUUID)
		wsURL := gc.WebSocketURL()
		if nodeName != "" {
			namespace := cv.GetNodeNamespace(nodeName)
			wsURL = fmt.Sprintf("ws://%s.%s.svc.cluster.local:%d", nodeName, namespace, griddleinfra.GatewayWSPort)
		}

		state.GatewayConnectors = append(state.GatewayConnectors, domain.GatewayConnectorState{
			NodeUUID:      gc.NodeUUID,
			AuthGatewayID: gc.AuthGatewayID,
			WebSocketURL:  wsURL,
			GatewayDonID:  gc.GatewayDonID,
			DONName:       donName,
		})
	}
}

// gatewayOwnDONAndNodeForUUID finds the gateway DON (and its single physical
// node) that owns a given gateway node's UUID, by scanning topology metadata
// rather than assuming any particular ordering.
func gatewayOwnDONAndNodeForUUID(topology *cre.Topology, cv *domain.ChartValues, nodeUUID string) (donName, nodeName string) {
	for _, donMeta := range topology.DonsMetadata.List() {
		gwNode, hasGateway := donMeta.Gateway()
		if !hasGateway || gwNode.UUID != nodeUUID {
			continue
		}
		if names := cv.NodeNamesForDONName(donMeta.Name); len(names) > 0 {
			nodeName = names[0]
		}
		return donMeta.Name, nodeName
	}
	return "", ""
}

// storeGatewayServiceConfigs copies the topology's gateway service configs (which
// Features' PreEnvStartup populated with handlers such as http-capabilities) into
// persisted state, so SyncJobs can restore them without re-running PreEnvStartup.
func (d *Deployer) storeGatewayServiceConfigs(state *domain.StateFile, topology *cre.Topology) {
	state.GatewayServiceConfigs = nil
	for _, svc := range topology.GatewayServiceConfigs {
		s := domain.GatewayServiceConfigState{
			ServiceName: svc.ServiceName,
			Handlers:    append([]string(nil), svc.Handlers...),
			DONs:        append([]string(nil), svc.DONs...),
		}
		if svc.Auth0 != nil {
			s.Auth0 = &domain.GatewayServiceAuth0State{
				IssuerURL: svc.Auth0.IssuerURL,
				Audience:  svc.Auth0.Audience,
				TenantID:  svc.Auth0.TenantID,
			}
		}
		state.GatewayServiceConfigs = append(state.GatewayServiceConfigs, s)
	}
}

// applyStoredGatewayServiceConfigs overwrites a freshly built topology's gateway
// service configs with the ones persisted during the on-chain phase (populated by
// Features' PreEnvStartup). No-op when nothing was persisted (e.g. non-gateway
// DONs), leaving the default web-api-capabilities handler from NewTopology.
func (d *Deployer) applyStoredGatewayServiceConfigs(state *domain.StateFile, topology *cre.Topology) {
	if len(state.GatewayServiceConfigs) == 0 {
		return
	}
	svcs := make([]cre.GatewayServiceConfig, 0, len(state.GatewayServiceConfigs))
	for _, s := range state.GatewayServiceConfigs {
		svc := cre.GatewayServiceConfig{
			ServiceName: s.ServiceName,
			Handlers:    append([]string(nil), s.Handlers...),
			DONs:        append([]string(nil), s.DONs...),
		}
		if s.Auth0 != nil {
			svc.Auth0 = &cre.GatewayServiceAuth0Config{
				IssuerURL: s.Auth0.IssuerURL,
				Audience:  s.Auth0.Audience,
				TenantID:  s.Auth0.TenantID,
			}
		}
		svcs = append(svcs, svc)
	}
	topology.GatewayServiceConfigs = svcs
}
