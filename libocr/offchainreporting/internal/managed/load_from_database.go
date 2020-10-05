package managed

import (
	"context"

	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
)

func loadConfigFromDatabase(ctx context.Context, database types.Database, logger types.Logger) *types.ContractConfig {
	cc, err := database.ReadConfig(ctx)
	if err != nil {
		logger.Error("loadConfigFromDatabase: Error during Database.ReadConfig", types.LogFields{
			"error": err,
		})
		return nil
	}

	if cc == nil {
		logger.Info("loadConfigFromDatabase: Database.ReadConfig returned nil, no configuration to restore", nil)
		return nil
	}

	return cc
}
