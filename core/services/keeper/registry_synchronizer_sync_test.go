package keeper

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// GetUpkeepFailure implements the upkeepGetter interface with an induced error and nil
// config response.
type GetUpkeepFailure struct{}

var errGetUpkeep = errors.New("chain connection error example")

func (g *GetUpkeepFailure) GetUpkeep(opts *bind.CallOpts, id *big.Int) (*UpkeepConfig, error) {
	return nil, fmt.Errorf("%w [%s]: getConfig v1.%d", ErrContractCallFailure, errGetUpkeep, RegistryVersion_1_2)
}

func TestSyncUpkeepWithCallback_UpkeepNotFound(t *testing.T) {
	log, logObserver := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	synchronizer := &RegistrySynchronizer{
		logger: log.(logger.SugaredLogger),
	}

	addr := ethkey.EIP55Address(testutils.NewAddress().Hex())
	registry := Registry{
		ContractAddress: addr,
	}

	o, ok := new(big.Int).SetString("5032485723458348569331745", 10)
	if !ok {
		t.FailNow()
	}

	id := utils.NewBig(o)
	count := 0
	doneFunc := func() {
		count++
	}

	getter := &GetUpkeepFailure{}
	synchronizer.syncUpkeepWithCallback(getter, registry, id, doneFunc)

	// logs should have the upkeep identifier included in the error context properly formatted
	require.Equal(t, 1, logObserver.Len())

	keys := map[string]bool{}
	for _, entry := range logObserver.All() {
		for _, field := range entry.Context {
			keys[field.Key] = true

			switch field.Key {
			case "error":
				require.Equal(t, "failed to get upkeep config: failure in calling contract [chain connection error example]: getConfig v1.2", field.String)
			case "upkeepID":
				require.Equal(t, fmt.Sprintf("UPx%064s", "429ab990419450db80821"), field.String)
			case "registryContract":
				require.Equal(t, addr.Hex(), field.String)
			}
		}
	}

	require.Equal(t, map[string]bool{"upkeepID": true, "error": true, "registryContract": true}, keys)
	require.Equal(t, 1, count, "callback function should run")
}
