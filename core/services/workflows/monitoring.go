package workflows

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PromExecutionTimeMS = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "capability_execution_time_ms",
			Help: "Metric representing the execution time in milliseconds",
		},
		[]string{"keystone_type", "id"}, //i.e. trigger, cron
	)
	PromTaskRunSuccessCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "capability_runs_count",
			Help: "Metric representing the number of runs completed successfully",
		},
		[]string{"keystone_type", "id"}, //i.e. consensus, ocr
	)
	PromTaskRunFaultCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "capability_runs_fault_count",
			Help: "Metric representing the number of runs with an application fault",
		},
		[]string{"keystone_type", "id"}, //i.e. target, write
	)
	PromTaskRunInvalidCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "capability_runs_invalid_count",
			Help: "Metric representing the number of runs with an application fault",
		},
		[]string{"keystone_type", "id"},
	)
	PromTaskRunUnauthorizedCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "capability_runs_unauthorized_count",
			Help: "Metric representing the number of runs with an application fault",
		},
		[]string{"keystone_type", "id"},
	)
	PromTaskRunNoResourceCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "capability_runs_no_resource_count",
			Help: "Metric representing the number of runs with an application fault",
		},
		[]string{"keystone_type", "id"},
	)
)
