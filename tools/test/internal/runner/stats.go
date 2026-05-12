package runner

import (
	"math"

	"charm.land/lipgloss/v2"

	"github.com/smartcontractkit/chainlink/v2/tools/test/internal/termstyle"
)

func ciStyleForGap(gap float64) lipgloss.Style {
	switch {
	case gap <= 0.10:
		return termstyle.OK
	case gap <= 0.30:
		return termstyle.Accent
	default:
		return termstyle.Bad
	}
}

// WilsonScoreInterval calculates our confidence interval for how (non) flaky a test is.
// We use a 95% confidence interval, so z = 1.96.
// https://en.wikipedia.org/wiki/Binomial_proportion_confidence_interval
func WilsonScoreInterval(k, n int, z float64) (lower float64, upper float64) {
	if n == 0 {
		return 0, 0
	}
	if z == 0 {
		z = 1.96
	}
	p := float64(k) / float64(n)
	nf := float64(n)
	zSq := z * z
	denom := 1.0 + zSq/nf
	center := (p + zSq/(2.0*nf)) / denom
	margin := (z / denom) * math.Sqrt(p*(1.0-p)/nf+zSq/(4.0*nf*nf))
	lower = math.Max(0.0, center-margin)
	upper = math.Min(1.0, center+margin)
	return lower, upper
}
