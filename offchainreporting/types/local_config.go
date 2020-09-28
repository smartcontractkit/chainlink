package types

import "time"

// LocalConfig contains oracle-specific configuration details which are not
// mandated by the on-chain configuration specification via OffchainAggregator.SetConfig
type LocalConfig struct {
	// Timeout for making observations.
	// (This is necessary because an oracle's operations are serialized, so
	// blocking forever on an observation would break the oracle.)
	DataSourceTimeout time.Duration

	// Timeout for blockchain interactions (mediated through
	// ContractConfigTracker and ContractTransmitter).
	// (This is necessary because an oracle's operations are serialized, so
	// blocking forever on a chain interaction would break the oracle.)
	BlockchainTimeout time.Duration

	// Polling interval at which ContractConfigTracker is queried for
	// updated on-chain configurations. Recommended values are between
	// fifteen seconds and two minutes.
	ContractConfigTrackerPollInterval time.Duration

	// Interval at which we try to establish a subscription on ContractConfigTracker
	// if one doesn't exist. Recommended values are between two and five minutes.
	ContractConfigTrackerSubscribeInterval time.Duration

	// Number of block confirmations to wait for before enacting an on-chain
	// configuration change. This value doesn't need to be very high (in
	// particular, it does not need to protect against malicious re-orgs).
	// Since configuration changes create some overhead, and mini-reorgs
	// are fairly common, recommended values are between two and ten.
	//
	// Malicious re-orgs are not any more of concern here than they are in
	// blockchain applications in general: Since nodes check the contract for the
	// latest config every ContractConfigTrackerPollInterval.Seconds(), they will
	// come to a common view of the current config within any interval longer than
	// that, as long as the latest setConfig transaction in the longest chain is
	// stable. They will thus be able to continue reporting after the poll
	// interval, unless an adversary is able to repeatedly re-org the transaction
	// out during every poll interval, which would amount to the capability to
	// censor any transaction.
	ContractConfigConfirmations uint16
}
