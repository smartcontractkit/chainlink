package sets

import (
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	docker_blockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/docker"
	k8s_blockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/kubernetes"
)

func NewDeployerSet(commonLogger logger.Logger, testLogger zerolog.Logger, provider *infra.Provider) map[blockchains.ChainFamily]blockchains.Deployer {
	if provider.IsDocker() {
		return docker_blockchains.NewDeployerSet(commonLogger)
	}

	return k8s_blockchains.NewDeployerSet(commonLogger, testLogger, provider.CRIB.Namespace, environment.CribConfigsDir)
}
