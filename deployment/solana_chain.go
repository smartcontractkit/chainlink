package deployment

const (
	RouterProgramName               = "ccip_router"
	OffRampProgramName              = "ccip_offramp"
	FeeQuoterProgramName            = "fee_quoter"
	BurnMintTokenPoolProgramName    = "burnmint_token_pool"
	LockReleaseTokenPoolProgramName = "lockrelease_token_pool"
	AccessControllerProgramName     = "access_controller"
	TimelockProgramName             = "timelock"
	McmProgramName                  = "mcm"
	RMNRemoteProgramName            = "rmn_remote"
	ReceiverProgramName             = "test_ccip_receiver"
)

// https://docs.google.com/document/d/1Fk76lOeyS2z2X6MokaNX_QTMFAn5wvSZvNXJluuNV1E/edit?tab=t.0#heading=h.uij286zaarkz
// https://docs.google.com/document/d/1nCNuam0ljOHiOW0DUeiZf4ntHf_1Bw94Zi7ThPGoKR4/edit?tab=t.0#heading=h.hju45z55bnqd
var SolanaProgramBytes = map[string]int{
	RouterProgramName:               5 * 1024 * 1024,
	OffRampProgramName:              1.5 * 1024 * 1024, // router should be redeployed but it does support upgrades if required (big fixes etc.)
	FeeQuoterProgramName:            5 * 1024 * 1024,
	BurnMintTokenPoolProgramName:    3 * 1024 * 1024,
	LockReleaseTokenPoolProgramName: 3 * 1024 * 1024,
	AccessControllerProgramName:     1 * 1024 * 1024,
	TimelockProgramName:             1 * 1024 * 1024,
	McmProgramName:                  1 * 1024 * 1024,
	RMNRemoteProgramName:            3 * 1024 * 1024,
}
