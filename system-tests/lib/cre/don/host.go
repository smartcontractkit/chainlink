package don

import (
	"fmt"

	coregateway "github.com/smartcontractkit/chainlink/v2/core/services/gateway"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

// Unfortunately, we need to construct some of these URLs before any environment is created, because they are used
// in CL node configs. This introduces a coupling between Helm charts used by CRIB and Docker container names used by CTFv2.
func InternalHost(nodeIndex int, nodeType cre.CapabilityFlag, donName string, infraInput infra.Input) string {
	if infraInput.Type == infra.CRIB {
		if nodeType == cre.BootstrapNode {
			return fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		return fmt.Sprintf("%s-%d", donName, nodeIndex)
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func InternalGatewayHost(nodeIndex int, nodeType cre.CapabilityFlag, donName string, infraInput infra.Input) string {
	if infraInput.Type == infra.CRIB {
		host := fmt.Sprintf("%s-%d", donName, nodeIndex)
		if nodeType == cre.BootstrapNode {
			host = fmt.Sprintf("%s-bt-%d", donName, nodeIndex)
		}
		host += "-gtwnode"

		return host
	}

	return fmt.Sprintf("%s-node%d", donName, nodeIndex)
}

func ExternalGatewayHost(nodeIndex int, nodeType cre.CapabilityFlag, donName string, infraInput infra.Input) string {
	if infraInput.Type == infra.CRIB {
		return infraInput.CRIB.Namespace + "-gateway.main.stage.cldev.sh"
	}

	return "localhost"
}

func ExternalGatewayPort(infraInput infra.Input, handlerType coregateway.HandlerType) int {
	if infraInput.Type == infra.CRIB {
		return 80
	}

	switch handlerType {
	case coregateway.WebAPICapabilitiesType:
		return 5002
	case coregateway.VaultHandlerType:
		return 15002
	default:
		return 5002
	}
}
