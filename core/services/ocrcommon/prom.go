package ocrcommon

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/cast"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	PromBridgeJsonParseValues = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bridge_json_parse_values",
		Help: "Values returned by json_parse for bridge task",
	},
		[]string{"job_id", "job_name", "bridge_name", "task_id"})

	PromOcrMedianValues = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ocr_median_values",
		Help: "Median value returned by ocr job",
	},
		[]string{"job_id", "job_name"})
)

// promSetBridgeParseMetrics will parse pipeline.TaskRunResults for bridge tasks, get the pipeline.TaskTypeJSONParse task and update prometheus metrics with it
func promSetBridgeParseMetrics(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults) {
	if ds.jb.Type.String() != pipeline.OffchainReportingJobType && ds.jb.Type.String() != pipeline.OffchainReporting2JobType {
		return
	}

	for _, trr := range *trrs {
		if trr.Task.Type() == pipeline.TaskTypeBridge {
			nextTask := trrs.GetNextTaskOf(trr)

			if nextTask != nil && nextTask.Task.Type() == pipeline.TaskTypeJSONParse {
				fetchedValue := cast.ToFloat64(nextTask.Result.Value)

				PromBridgeJsonParseValues.WithLabelValues(fmt.Sprintf("%d", ds.jb.ID), ds.jb.Name.String, trr.Task.(*pipeline.BridgeTask).Name, trr.Task.DotID()).Set(fetchedValue)
			}
		}
	}
}

// promSetFinalResultMetrics will check if job is pipeline.OffchainReportingJobType or pipeline.OffchainReporting2JobType then send the pipeline.FinalResult to prometheus
func promSetFinalResultMetrics(ds *inMemoryDataSource, finalResult *pipeline.FinalResult) {
	if ds.jb.Type.String() != pipeline.OffchainReportingJobType && ds.jb.Type.String() != pipeline.OffchainReporting2JobType {
		return
	}

	singularResult, err := finalResult.SingularResult()
	if err != nil {
		return
	}

	finalResultDecimal, err := utils.ToDecimal(singularResult.Value)
	if err != nil {
		return
	}
	finalResultFloat, _ := finalResultDecimal.Float64()
	PromOcrMedianValues.WithLabelValues(fmt.Sprintf("%d", ds.jb.ID), ds.jb.Name.String).Set(finalResultFloat)
}
