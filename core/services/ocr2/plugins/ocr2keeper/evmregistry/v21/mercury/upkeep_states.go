package mercury

type MercuryUpkeepState uint8

const (
	NoPipelineError        MercuryUpkeepState = 0
	RpcFlakyFailure        MercuryUpkeepState = 3
	MercuryFlakyFailure    MercuryUpkeepState = 4
	PackUnpackDecodeFailed MercuryUpkeepState = 5
	MercuryUnmarshalError  MercuryUpkeepState = 6
	InvalidMercuryRequest  MercuryUpkeepState = 7
	InvalidMercuryResponse MercuryUpkeepState = 8 // this will only happen if Mercury server sends bad responses
	UpkeepNotAuthorized    MercuryUpkeepState = 9
)
