package por

import (
	"context"
)

type ContractReader interface {
	// GetLatestTransmittedReportDetails retrieves the latest (even unfinalized) transmission details from the contract on-chain.
	// (Some of the returned values, such as latestTimestamp, are not used in the current implementation but are included for future extensibility.)
	GetLatestTransmittedReportDetails(
		ctx context.Context,
		chain ChainSelector,
	) (
		details TransmittedReportDetails,
		err error,
	)
}
