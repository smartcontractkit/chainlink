package helpers

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
)

// RequirePrometheusQueryEventually polls an instant PromQL query against the
// local observability stack's Prometheus (framework.LocalPrometheusBaseURL,
// deployed by `env start --with-dashboards`) until accept returns true for the
// sum of the query's result values. An empty result sums to 0, so a query for
// a metric that has not been emitted yet keeps polling instead of failing -
// pair it with an accept that requires a positive value.
func RequirePrometheusQueryEventually(t *testing.T, query string, timeout, interval time.Duration, accept func(sum float64) bool, msg string, msgArgs ...any) {
	t.Helper()

	pc := framework.NewPrometheusQueryClient(framework.LocalPrometheusBaseURL)

	require.Eventuallyf(t, func() bool {
		resp, err := pc.Query(query, time.Now())
		if err != nil {
			framework.L.Debug().Err(err).Str("query", query).Msg("Prometheus query not ready yet")
			return false
		}

		sum := 0.0
		for _, result := range resp.Data.Result {
			if len(result.Value) < 2 {
				continue
			}
			value, parseErr := strconv.ParseFloat(fmt.Sprint(result.Value[1]), 64)
			if parseErr != nil {
				continue
			}
			sum += value
		}
		return accept(sum)
	}, timeout, interval, msg, msgArgs...)
}
