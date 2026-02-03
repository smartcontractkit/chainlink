package ccvcommon

import (
	"slices"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	executormonitoring "github.com/smartcontractkit/chainlink-ccv/executor/pkg/monitoring"
	verifiermonitoring "github.com/smartcontractkit/chainlink-ccv/verifier/pkg/monitoring"
)

func MetricViews() []sdkmetric.View {
	return slices.Concat(executormonitoring.MetricViews(), verifiermonitoring.MetricViews())
}
