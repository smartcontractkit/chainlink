package types

import "time"

const EnableDangerousDevelopmentMode = "enable dangerous development mode"

// LocalConfig contains oracle-specific configuration details which are not
// mandated by the on-chain configuration specification via OffchainAggregator.SetConfig
type LocalConfig struct {
	// Timeout for blockchain queries (mediated through
	// ContractConfigTracker and ContractTransmitter).
	// (This is necessary because an oracle's operations are serialized, so
	// blocking forever on a chain interaction would break the oracle.)
	BlockchainTimeout time.Duration

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
	//
	// Note that 1 confirmation implies that the transaction/event has been mined in one block.
	// 0 confirmations would imply that the event would be recognised before it has even been mined, which is not currently supported.
	// e.g.
	// Current block height: 42
	// Changed in block height: 43
	// Contract config confirmations: 1
	// STILL PENDING
	//
	// Current block height: 43
	// Changed in block height: 43
	// Contract config confirmations: 1
	// CONFIRMED
	ContractConfigConfirmations uint16

	// SkipContractConfigConfirmations allows to disable the confirmations check entirely
	// This can be useful in some cases e.g. L2 which has instant finality and
	// where local block numbers do not match the on-chain value returned from
	// block.number
	SkipContractConfigConfirmations bool

	// Polling interval at which ContractConfigTracker is queried for
	// updated on-chain configurations. Recommended values are between
	// fifteen seconds and two minutes.
	ContractConfigTrackerPollInterval time.Duration

	// Interval at which we try to establish a subscription on ContractConfigTracker
	// if one doesn't exist. Recommended values are between two and five minutes.
	ContractConfigTrackerSubscribeInterval time.Duration

	// Timeout for ContractTransmitter.Transmit calls.
	ContractTransmitterTransmitTimeout time.Duration

	// Timeout for database interactions.
	// (This is necessary because an oracle's operations are serialized, so
	// blocking forever on an observation would break the oracle.)
	DatabaseTimeout time.Duration

	// Timeout for making observations using the DataSource.Observe method.
	// (This is necessary because an oracle's operations are serialized, so
	// blocking forever on an observation would break the oracle.)
	DataSourceTimeout time.Duration

	// After DataSourceTimeout expires, we additionally wait for this grace
	// period for DataSource.Observe to return a result, before forcibly moving
	// on.
	DataSourceGracePeriod time.Duration

	// DANGER, this turns off all kinds of sanity checks. May be useful for testing.
	// Set this to EnableDangerousDevelopmentMode to turn on dev mode.
	DevelopmentMode string
}
