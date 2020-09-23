package types

import "time"

type LocalConfig struct {
	DataSourceTimeout time.Duration

	BlockchainTimeout time.Duration

	ContractConfigTrackerPollInterval time.Duration

	ContractConfigTrackerSubscribeInterval time.Duration

	ContractConfigConfirmations uint16
}
