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

// DefaultPremiumPercentage is the premium percentage used to reserve the
// worst-case fee when the job spec does not set premiumPercentage. It covers
// every premium configured on current deployments (highest observed is 84) with
// headroom. Set premiumPercentage in the job spec to match a chain's configured
// premium exactly; the coordinator caps premiums at 155.
const DefaultPremiumPercentage uint32 = 100

// withPremium applies the coordinator's premium percentage to a gas cost.
func withPremium(cost *big.Int, premiumPercentage uint32) *big.Int {
	return cost.Mul(cost, big.NewInt(int64(100+premiumPercentage))).Quo(cost, big.NewInt(100))
}

// EstimateFeeJuels estimates the amount of link needed to fulfill a request
// given the callback gas limit, the gas overhead, the gas price, the premium
// percentage and the wei per unit link. An error is returned if the wei per
// unit link provided is zero.
func EstimateFeeJuels(callbackGasLimit, gasOverhead uint32, maxGasPriceWei, weiPerUnitLink *big.Int, premiumPercentage uint32) (*big.Int, error) {
	if weiPerUnitLink.Cmp(big.NewInt(0)) == 0 {
		return nil, errors.New("wei per unit link is zero")
	}
	maxGasUsed := big.NewInt(int64(callbackGasLimit + gasOverhead))
	costWei := maxGasUsed.Mul(maxGasUsed, maxGasPriceWei)
	costWei = withPremium(costWei, premiumPercentage)
	// Multiply by 1e18 first so that we don't lose a ton of digits due to truncation when we divide
	// by weiPerUnitLink
	numerator := costWei.Mul(costWei, big.NewInt(1e18))
	costJuels := numerator.Quo(numerator, weiPerUnitLink)
	return costJuels, nil
}

// EstimateFeeWei estimates the amount of wei needed to fulfill a request
func EstimateFeeWei(callbackGasLimit, gasOverhead uint32, maxGasPriceWei *big.Int, premiumPercentage uint32) (*big.Int, error) {
	maxGasUsed := big.NewInt(int64(callbackGasLimit + gasOverhead))
	costWei := maxGasUsed.Mul(maxGasUsed, maxGasPriceWei)
	return withPremium(costWei, premiumPercentage), nil
}

// gasOverhead returns the gas reserved for coordinator verification and
// accounting: the job spec override if set, otherwise the measured default.
func (lsn *listenerV2) gasOverhead() uint32 {
	if lsn.job.VRFSpec.GasOverhead != nil {
		return *lsn.job.VRFSpec.GasOverhead
	}
	return GasProofVerification
}

// premiumPercentage returns the premium percentage used to reserve worst-case
// fees: the job spec override if set, otherwise the default.
func (lsn *listenerV2) premiumPercentage() uint32 {
	if lsn.job.VRFSpec.PremiumPercentage != nil {
		return *lsn.job.VRFSpec.PremiumPercentage
	}
	return DefaultPremiumPercentage
}

// reservedFee returns the amount to reserve from a subscription for a request:
// the larger of the simulated charge and the worst-case charge (full callback
// gas limit at the gas lane maximum price, including the premium).
func reservedFee(p *vrfPipelineResult) *big.Int {
	switch {
	case p.maxFee == nil:
		return p.fundsNeeded
	case p.fundsNeeded == nil:
		return p.maxFee
	case p.maxFee.Cmp(p.fundsNeeded) >= 0:
		return p.maxFee
	default:
		return p.fundsNeeded
	}
}

// WorstCaseGasFloorMargin covers the fulfillment event emission and counter
// writes not included in the gas overhead.
const WorstCaseGasFloorMargin uint32 = 10_000

// WorstCaseGasFloor returns the gas limit needed to fulfill a request whose
// callback uses its entire callbackGasLimit: the callback allotment plus the
// EIP-150 retained gas required by _callWithExactGas, and the coordinator's
// verification and accounting work (gasOverhead).
func WorstCaseGasFloor(callbackGasLimit, gasOverhead uint32) uint64 {
	return uint64(callbackGasLimit)*64/63 +
		uint64(gasOverhead) +
		uint64(WorstCaseGasFloorMargin)
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
