package mercury

type MercuryUpkeepState uint8

// NOTE: This enum should be kept in sync with evmregistry/v21/encoding/interface.go
// TODO (AUTO-7928) Remove this duplication
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
