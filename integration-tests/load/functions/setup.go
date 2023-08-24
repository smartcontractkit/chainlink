package loadfunctions

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/docker/test_env"
	chainlinkutils "github.com/smartcontractkit/chainlink/v2/core/utils"
	"math/big"
)

type Contracts struct {
	LinkToken      contracts.LinkToken
	Coordinator    contracts.FunctionsCoordinator
	Router         contracts.FunctionsRouter
	LoadTestClient contracts.FunctionsLoadTestClient
}

func SetupLocalLoadTestEnv(cfg *PerformanceConfig) (*test_env.CLClusterTestEnv, *Contracts, error) {
	env, err := test_env.NewCLTestEnvBuilder().
		Build()
	if err != nil {
		return env, nil, err
	}
	env.EVMClient.GetNonceSetting()
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
	loadTestClient, err := env.ContractLoader.LoadFunctionsLoadTestClient(cfg.Common.LoadTestClient)
	if err != nil {
		return env, nil, err
	}
	if cfg.Common.SubscriptionID == 0 {
		log.Info().Msg("Creating new subscription")
		subID, err := router.CreateSubscriptionWithConsumer(loadTestClient.Address())
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
		LinkToken:      lt,
		Coordinator:    coord,
		Router:         router,
		LoadTestClient: loadTestClient,
	}, nil
}
