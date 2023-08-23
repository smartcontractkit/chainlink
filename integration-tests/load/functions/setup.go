package loadfunctions

import (
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
	"math/big"
)

type Contracts struct {
	LinkToken     contracts.LinkToken
	Coordinator   contracts.FunctionsCoordinator
	Router        contracts.FunctionsRouter
	ClientExample contracts.FunctionsClientExample
}

func SetupLocalLoadTestEnv(cfg *PerformanceConfig) (*test_env.CLClusterTestEnv, *Contracts, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		Build()
	if err != nil {
		return env, nil, err
	}
	lt, err := env.ContractLoader.LoadLINKToken(cfg.Common.LINKTokenAddr)
	if err != nil {
		return env, nil, err
	}
	coord, err := env.ContractLoader.LoadFunctionsCoordinator(cfg.Common.Coordinator)
	if err != nil {
		return env, nil, err
	}
	router, err := env.ContractLoader.LoadFunctionsRouter(cfg.Common.Router)
	if err != nil {
		return env, nil, err
	}
	clientExample, err := env.ContractLoader.LoadFunctionsClientExample(cfg.Common.ClientExample)
	if err != nil {
		return env, nil, err
	}
	if cfg.Common.SubscriptionID == 0 {
		subID, err := router.CreateSubscriptionWithConsumer(clientExample.Address())
		if err != nil {
			return env, nil, errors.Wrap(err, "failed to create a new subscription")
		}
		encodedSubId, err := chainlinkutils.ABIEncode(`[{"type":"uint64"}]`, subID)
		if err != nil {
			return env, nil, errors.Wrap(err, "failed to encode subscription ID for funding")
		}
		_, err = lt.TransferAndCall(router.Address(), big.NewInt(0).Mul(cfg.Common.Funding.SubFunds, big.NewInt(1e18)), encodedSubId)
		if err != nil {
			return env, nil, errors.Wrap(err, "failed to transferAndCall router, LINK funding")
		}
		cfg.Common.SubscriptionID = subID
	}
	return env, &Contracts{
		LinkToken:     lt,
		Coordinator:   coord,
		Router:        router,
		ClientExample: clientExample,
	}, nil
}
