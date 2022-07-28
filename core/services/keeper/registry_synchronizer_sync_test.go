package keeper

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// GetUpkeepFailure implements the UpkeepGetter interface with an induced error and nil
// config response.
type GetUpkeepFailure struct{}

var GetUpkeepError = errors.New("upkeep config not found")

func (g *GetUpkeepFailure) GetUpkeep(opts *bind.CallOpts, id *big.Int) (*UpkeepConfig, error) {
	return nil, GetUpkeepError
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

	id := utils.NewBigI(3)
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
				require.Equal(t, GetUpkeepError.Error(), field.String)
			case "upkeepID":
				require.Equal(t, "3", field.String)
			case "registryContract":
				require.Equal(t, addr.Hex(), field.String)
			}
		}
	}

	require.Equal(t, map[string]bool{"upkeepID": true, "error": true, "registryContract": true}, keys)
}
