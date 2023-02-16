package ocrcommon

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var promBridgeTaskValues = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "bridge_task_values",
	Help: "-",
},
	[]string{"contract", "job_id", "job_name", "bridge_name", "task_id"})

// inMemoryDataSource is an abstraction over the process of initiating a pipeline run
// and returning the result. Additionally, it converts the result to an
// ocrtypes.Observation (*big.Int), as expected by the offchain reporting library.
type inMemoryDataSource struct {
	pipelineRunner pipeline.Runner
	jb             job.Job
	spec           pipeline.Spec
	lggr           logger.Logger

	current bridges.BridgeMetaData
	mu      sync.RWMutex
}

// dataSource uses inMemoryDataSource and implements capturing the result to be stored in the DB
type dataSource struct {
	inMemoryDataSource
	runResults chan<- pipeline.Run
}

// dataSourceV2 implements dataSource with the proper Observe return type
type dataSourceV2 struct {
	dataSource
}

func (ds *dataSourceV2) Observe(ctx context.Context) (*big.Int, error) {
	return ds.dataSource.Observe(ctx)
}

func NewDataSourceV1(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, runResults chan<- pipeline.Run) ocrtypes.DataSource {
	return &dataSource{
		inMemoryDataSource: inMemoryDataSource{
			pipelineRunner: pr,
			jb:             jb,
			spec:           spec,
			lggr:           lggr,
		},
		runResults: runResults,
	}
}

func NewDataSourceV2(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger, runResults chan<- pipeline.Run) median.DataSource {
	return &dataSourceV2{
		dataSource: dataSource{
			inMemoryDataSource: inMemoryDataSource{
				pipelineRunner: pr,
				jb:             jb,
				spec:           spec,
				lggr:           lggr,
			},
			runResults: runResults,
		},
	}
}

func NewInMemoryDataSource(pr pipeline.Runner, jb job.Job, spec pipeline.Spec, lggr logger.Logger) median.DataSource {
	return &inMemoryDataSource{
		pipelineRunner: pr,
		jb:             jb,
		spec:           spec,
		lggr:           lggr,
	}
}

var _ ocrtypes.DataSource = (*dataSource)(nil)

func (ds *inMemoryDataSource) updateAnswer(a *big.Int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.current = bridges.BridgeMetaData{
		LatestAnswer: a,
		UpdatedAt:    big.NewInt(time.Now().Unix()),
	}
}

func (ds *inMemoryDataSource) currentAnswer() (*big.Int, *big.Int) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.current.LatestAnswer, ds.current.UpdatedAt
}

// The context passed in here has a timeout of (ObservationTimeout + ObservationGracePeriod).
// Upon context cancellation, its expected that we return any usable values within ObservationGracePeriod.
func (ds *inMemoryDataSource) executeRun(ctx context.Context) (pipeline.Run, pipeline.FinalResult, error) {
	md, err := bridges.MarshalBridgeMetaData(ds.currentAnswer())
	if err != nil {
		ds.lggr.Warnw("unable to attach metadata for run", "err", err)
	}

	vars := pipeline.NewVarsFrom(map[string]interface{}{
		"jb": map[string]interface{}{
			"databaseID":    ds.jb.ID,
			"externalJobID": ds.jb.ExternalJobID,
			"name":          ds.jb.Name.ValueOrZero(),
		},
		"jobRun": map[string]interface{}{
			"meta": md,
		},
	})

	run, trrs, err := ds.pipelineRunner.ExecuteRun(ctx, ds.spec, vars, ds.lggr)
	if err != nil {
		return pipeline.Run{}, pipeline.FinalResult{}, errors.Wrapf(err, "error executing run for spec ID %v", ds.spec.ID)
	}
	finalResult := trrs.FinalResult(ds.lggr)

	extractBridgeMetrics(ds, &trrs, &finalResult)

	return run, finalResult, err
}

func extractBridgeMetrics(ds *inMemoryDataSource, trrs *pipeline.TaskRunResults, finalResult *pipeline.FinalResult) {
	for _, trr := range *trrs {
		if trr.Task.Type() == pipeline.TaskTypeBridge {

			if nextTask := trrs.GetNextTaskOf(trr); nextTask != nil && nextTask.Task.Type() == pipeline.TaskTypeJSONParse {
				fetchedValue, err := getFloat(nextTask.Result.Value)
				if err == nil {
					contractAddress := ""
					if ds.jb.Type == "offchainreporting2" {
						contractAddress = ds.jb.OCR2OracleSpec.ContractID
					}

					promBridgeTaskValues.WithLabelValues(contractAddress, fmt.Sprintf("%d", ds.jb.ID), ds.jb.Name.String, trr.Task.(*pipeline.BridgeTask).Name, trr.Task.DotID()).Set(fetchedValue)
				}
			}
		}
	}
}

var floatType = reflect.TypeOf(float64(0))

func getFloat(unk interface{}) (float64, error) {
	if unk == nil {
		return 0, fmt.Errorf("cannot convert nil to float64")
	}
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

// parse uses the finalResult into a big.Int and stores it in the bridge metadata
func (ds *inMemoryDataSource) parse(finalResult pipeline.FinalResult) (*big.Int, error) {
	result, err := finalResult.SingularResult()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting singular result for job ID %v", ds.spec.JobID)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	asDecimal, err := utils.ToDecimal(result.Value)
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert observation to decimal")
	}
	ds.updateAnswer(asDecimal.BigInt())
	return asDecimal.BigInt(), nil
}

// Observe without saving to DB
func (ds *inMemoryDataSource) Observe(ctx context.Context) (*big.Int, error) {
	_, finalResult, err := ds.executeRun(ctx)
	if err != nil {
		return nil, err
	}
	return ds.parse(finalResult)
}

// Observe with saving to DB
func (ds *dataSource) Observe(ctx context.Context) (ocrtypes.Observation, error) {
	run, finalResult, err := ds.inMemoryDataSource.executeRun(ctx)
	if err != nil {
		return nil, err
	}

	// Do the database write in a non-blocking fashion
	// so we can return the observation results immediately.
	// This is helpful in the case of a blocking API call, where
	// we reach the passed in context deadline and we want to
	// immediately return any result we have and do not want to have
	// a db write block that.
	select {
	case ds.runResults <- run:
	default:
		// If we're unable to enqueue a write, still return the value we have but warn.
		ds.lggr.Warnf("unable to enqueue run save for job ID %d, buffer full", ds.inMemoryDataSource.spec.JobID)
		return ds.inMemoryDataSource.parse(finalResult)
	}

	return ds.inMemoryDataSource.parse(finalResult)
}
