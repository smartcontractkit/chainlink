package vrfv2

import (
	"fmt"

	"github.com/rs/zerolog"
)

func logRandRequest(
	l zerolog.Logger,
	consumer string,
	coordinator string,
	subID uint64,
	minimumConfirmations uint16,
	callbackGasLimit uint32,
	numberOfWords uint32,
	randomnessRequestCountPerRequest uint16,
	randomnessRequestCountPerRequestDeviation uint16,
	keyhash [32]byte,
	keyNum int,
) {
	l.Info().
		Int("KeyNum", keyNum).
		Str("Consumer", consumer).
		Str("Coordinator", coordinator).
		Uint64("SubID", subID).
		Uint16("MinimumConfirmations", minimumConfirmations).
		Uint32("CallbackGasLimit", callbackGasLimit).
		Uint32("NumberOfWords", numberOfWords).
		Uint16("RandomnessRequestCountPerRequest", randomnessRequestCountPerRequest).
		Uint16("RandomnessRequestCountPerRequestDeviation", randomnessRequestCountPerRequestDeviation).
		Str("Keyhash", fmt.Sprintf("0x%x", keyhash)).
		Msg("Requesting randomness")
}
