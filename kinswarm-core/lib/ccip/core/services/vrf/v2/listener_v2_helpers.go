package v2

import (
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
)

func uniqueReqs(reqs []pendingRequest) int {
	s := map[string]struct{}{}
	for _, r := range reqs {
		s[r.req.RequestID().String()] = struct{}{}
	}
	return len(s)
}

// GasProofVerification is an upper limit on the gas used for verifying the VRF proof on-chain.
// It can be used to estimate the amount of LINK or native needed to fulfill a request.
const GasProofVerification uint32 = 200_000

// EstimateFeeJuels estimates the amount of link needed to fulfill a request
// given the callback gas limit, the gas price, and the wei per unit link.
// An error is returned if the wei per unit link provided is zero.
func EstimateFeeJuels(callbackGasLimit uint32, maxGasPriceWei, weiPerUnitLink *big.Int) (*big.Int, error) {
	if weiPerUnitLink.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("wei per unit link is zero")
	}
	maxGasUsed := big.NewInt(int64(callbackGasLimit + GasProofVerification))
	costWei := maxGasUsed.Mul(maxGasUsed, maxGasPriceWei)
	// Multiply by 1e18 first so that we don't lose a ton of digits due to truncation when we divide
	// by weiPerUnitLink
	numerator := costWei.Mul(costWei, big.NewInt(1e18))
	costJuels := numerator.Quo(numerator, weiPerUnitLink)
	return costJuels, nil
}

// EstimateFeeWei estimates the amount of wei needed to fulfill a request
func EstimateFeeWei(callbackGasLimit uint32, maxGasPriceWei *big.Int) (*big.Int, error) {
	maxGasUsed := big.NewInt(int64(callbackGasLimit + GasProofVerification))
	costWei := maxGasUsed.Mul(maxGasUsed, maxGasPriceWei)
	return costWei, nil
}

// observeRequestSimDuration records the time between the given requests simulations or
// the time until it's first simulation, whichever is applicable.
// Cases:
// 1. Never simulated: in this case, we want to observe the time until simulated
// on the utcTimestamp field of the pending request.
// 2. Simulated before: in this case, lastTry will be set to a non-zero time value,
// in which case we'd want to use that as a relative point from when we last tried
// the request.
func observeRequestSimDuration(jobName string, extJobID uuid.UUID, vrfVersion vrfcommon.Version, pendingReqs []pendingRequest) {
	now := time.Now().UTC()
	for _, request := range pendingReqs {
		// First time around lastTry will be zero because the request has not been
		// simulated yet. It will be updated every time the request is simulated (in the event
		// the request is simulated multiple times, due to it being underfunded).
		if request.lastTry.IsZero() {
			vrfcommon.MetricTimeUntilInitialSim.
				WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).
				Observe(float64(now.Sub(request.utcTimestamp)))
		} else {
			vrfcommon.MetricTimeBetweenSims.
				WithLabelValues(jobName, extJobID.String(), string(vrfVersion)).
				Observe(float64(now.Sub(request.lastTry)))
		}
	}
}

func ptr[T any](t T) *T { return &t }

func isProofVerificationError(errMsg string) bool {
	// See VRF.sol for all these messages
	// NOTE: it's unclear which of these errors are impossible and which
	// may actually happen, so including them all to be safe.
	errMessages := []string{
		"invalid x-ordinate",
		"invalid y-ordinate",
		"zero scalar",
		"invZ must be inverse of z",
		"bad witness",
		"points in sum must be distinct",
		"First mul check failed",
		"Second mul check failed",
		"public key is not on curve",
		"gamma is not on curve",
		"cGammaWitness is not on curve",
		"sHashWitness is not on curve",
		"addr(c*pk+s*g)!=_uWitness",
		"invalid proof",
	}
	for _, msg := range errMessages {
		if strings.Contains(errMsg, msg) {
			return true
		}
	}
	return false
}
