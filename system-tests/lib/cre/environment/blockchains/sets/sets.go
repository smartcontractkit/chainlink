package sets

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	docker_blockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/docker"
	k8s_blockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/kubernetes"
)

func NewDeployerSet(testLogger zerolog.Logger, provider *infra.Provider) map[blockchain.ChainFamily]blockchains.Deployer {
	if provider.IsDocker() {
		return docker_blockchains.NewDeployerSet()
	}

	return k8s_blockchains.NewDeployerSet(testLogger, provider.CRIB.Namespace, environment.CribConfigsDir)
}
