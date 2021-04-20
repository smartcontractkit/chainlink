package web

import (
	"github.com/smartcontractkit/chainlink/core/gracefulpanic"
	log "github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Web struct {
	chDone       chan struct{}
	chStop       chan struct{}
	jobID        int32
	logger       *log.Logger
	runner       pipeline.Runner
	pipelineSpec pipeline.Spec
}

// NewFromJobSpec() - instantiates a Web singleton to execute the given job pipelineSpec
func NewFromJobSpec(
	jobSpec job.Job,
	runner pipeline.Runner,
) (*Web, error) {
	pipelineSpec := jobSpec.PipelineSpec

	logger := log.CreateLogger(
		log.Default.With(
			"jobID", jobSpec.ID,
		),
	)

	chDone := make(chan struct{})
	chStop := make(chan struct{})

	return &Web{chDone, chStop, jobSpec.ID, logger, runner, *pipelineSpec}, nil
}

// Start implements the job.Service interface.
func (web *Web) Start() error {
	web.logger.Debug("Web: Start")
	go gracefulpanic.WrapRecover(func() {
		defer close(web.chDone)
		web.runWebPipeline()
	})

	return nil
}

// Close stops this instance and cleans up pipeline resource.
func (web *Web) Close() error {
	web.logger.Debug("Web: Closing")
	close(web.chStop)
	<-web.chDone

	return nil
}

func (web *Web) runWebPipeline() {
	ctx, cancel := utils.ContextFromChan(web.chStop)
	defer cancel()
	_, _, err := web.runner.ExecuteAndInsertFinishedRun(ctx, web.pipelineSpec, pipeline.JSONSerializable{}, *web.logger, true)
	if err != nil {
		web.logger.Errorf("Error executing new runWebRequest for jobSpec ID %v", web.jobID)
	}
}
