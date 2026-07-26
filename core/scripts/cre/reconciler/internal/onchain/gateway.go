package onchain

import (
	"fmt"
	"slices"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	griddleinfra "github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	cre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func (d *Deployer) storeGatewayConnectors(desired *domain.DesiredState, cv *domain.ChartValues, state *domain.StateFile, topology *cre.Topology) {
	state.GatewayConnectors = nil
	if topology.GatewayConnectors == nil {
		return
	}

	gwNodes := cv.FindGatewayNodes()
	for i, config := range topology.GatewayConnectors.Configurations {
		if config == nil || config.GatewayConfiguration == nil {
			continue
		}
		gc := config.GatewayConfiguration
		wsURL := gc.WebSocketURL()
		donName := ""
		if i < len(gwNodes) {
			gwNode := gwNodes[i]
			namespace := cv.GetNodeNamespace(gwNode.Name)
			wsURL = fmt.Sprintf("ws://%s.%s.svc.cluster.local:%d", gwNode.Name, namespace, griddleinfra.GatewayWSPort)
			donName = gatewayDONNameForNode(desired, cv, gwNode.Name)
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

func gatewayDONNameForNode(desired *domain.DesiredState, cv *domain.ChartValues, nodeName string) string {
	for i := range desired.DONs {
		don := &desired.DONs[i]
		if !don.HasDONType("gateway") {
			continue
		}
		if slices.Contains(cv.NodeNamesForDONName(don.Name), nodeName) {
			return don.Name
		}
	}
	return ""
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
