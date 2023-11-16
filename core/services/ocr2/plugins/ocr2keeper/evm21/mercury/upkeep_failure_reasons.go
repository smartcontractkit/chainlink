package mercury

type MercuryUpkeepFailureReason uint8

const (
	// upkeep failure onchain reasons
	MercuryUpkeepFailureReasonNone                    MercuryUpkeepFailureReason = 0
	MercuryUpkeepFailureReasonTargetCheckReverted     MercuryUpkeepFailureReason = 3
	MercuryUpkeepFailureReasonUpkeepNotNeeded         MercuryUpkeepFailureReason = 4
	MercuryUpkeepFailureReasonMercuryCallbackReverted MercuryUpkeepFailureReason = 7
	// leaving a gap here for more onchain failure reasons in the future
	// upkeep failure offchain reasons
	MercuryUpkeepFailureReasonMercuryAccessNotAllowed MercuryUpkeepFailureReason = 32
	MercuryUpkeepFailureReasonInvalidRevertDataInput  MercuryUpkeepFailureReason = 34
)
