package ccvcommon

import (
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	executormonitoring "github.com/smartcontractkit/chainlink-ccv/executor/pkg/monitoring"
	verifiermonitoring "github.com/smartcontractkit/chainlink-ccv/verifier/pkg/monitoring"
)

func MetricViews() []sdkmetric.View {
	views := make([]sdkmetric.View, 0)
	views = append(views, executormonitoring.MetricViews()...)
	views = append(views, verifiermonitoring.MetricViews()...)
	return views
}
