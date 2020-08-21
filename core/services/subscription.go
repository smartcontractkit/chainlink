package services

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// JobSubscription is the interface for listening for requests from the Ethereum node for a particular job.
// The only thing you can do with it is Unsubscribe.
type JobSubscription struct {
	// Job needs to be public for callers of StartJobSubscription()
	Job            models.JobSpec
	logBroadcaster logBroadcasterRegistration
	runManager     runManagerCreate
	listeners      []*listener
	logger         jsLogger
}

func StartJobSubscription(job models.JobSpec, logBroadcaster logBroadcasterRegistration, runManager runManagerCreate, logger jsLogger) (JobSubscription, error) {
	js := JobSubscription{
		Job:            job,
		logBroadcaster: logBroadcaster,
		runManager:     runManager,
		logger:         logger,
	}
	js.listeners = js.makeAndRegisterListeners()
	if len(js.listeners) == 0 {
		return js, fmt.Errorf("unable to subscribe to any logs, check earlier errors in this message, and the initiator types, jobID=%s", job.ID)
	}
	return js, nil
}

func (js JobSubscription) makeAndRegisterListeners() []*listener {
	initiators := js.Job.InitiatorsFor(models.LogBasedChainlinkJobInitiators...)
	listeners := []*listener{}
	for _, initiator := range initiators {
		l := &listener{
			jobID:      js.Job.ID,
			initiator:  &initiator,
			runManager: js.runManager,
			logger:     js.logger,
		}
		if !js.logBroadcaster.Register(initiator.InitiatorParams.Address, l) {
			js.logger.Errorw("unable to register handler because log broadcaster has not started", "jobID", js.Job.ID, "initiator", initiator.ID)
		} else {
			listeners = append(listeners, l)
			js.logger.Debugw("handler registered to log broadcaster", "jobID", js.Job.ID, "initiator", initiator.ID)
		}
	}
	return listeners
}

// Unsubscribe stops listening for logs for all the initiators of a job.
func (js *JobSubscription) Unsubscribe() {
	for _, l := range js.listeners {
		js.logBroadcaster.Unregister(l.initiator.InitiatorParams.Address, l)
		js.logger.Debugw("handler unregistered from log broadcaster", "jobID", js.Job.ID, "initiator", l.initiator.ID)
	}
}

// listener subscribes to logs for one initiator in a job. If the log is relevant, it will be forwarde to the RunManager.
// listener implements eth.LogListener
type listener struct {
	jobID      *models.ID
	initiator  *models.Initiator // We need the whole initator for RunManager
	runManager runManagerCreate
	logger     jsLogger
}

var _ eth.LogListener = &listener{}

// HandleLog gets called by LogBroadcaster whenever a new log from the Ethereum matches this listener's job address.
// The log is ignored if it doesn't match this listeners initiator or if it's not valid.
// Otherwise a JobRun is created for it and passed to the RunManager for execution.
func (l *listener) HandleLog(lb eth.LogBroadcast, err error) {
	if err != nil {
		// TODO: why would we receive an error from LogBroadcaseter?
		l.logger.Errorw("received error from LogBroadcaster", "jobID", l.jobID, "initiator", l.initiator.ID, "error", err)
		return
	}
	if !l.isMatchingInitiator(lb) {
		return
	}
	lr := l.toLogRequest(lb)
	if !lr.Validate() {
		l.logger.Infow("log failed validation", lr.ForLogger()...)
		return
	}
	if err := lr.ValidateRequester(); err != nil {
		l.logger.Errorw("log failed requester validation", append(lr.ForLogger(), "error", err)...)
		if _, rmErr := l.runManager.CreateErrored(l.jobID, *l.initiator, err); rmErr != nil {
			l.logger.Errorw("failed to create JobRun in the errored state", append(lr.ForLogger(), "error", rmErr))
		}
		return
	}
	rr, err := lr.RunRequest()
	if err != nil {
		l.logger.Errorw("failed to unmarshal log into a RunRequest", append(lr.ForLogger(), "error", err)...)
		if _, rmErr := l.runManager.CreateErrored(l.jobID, *l.initiator, err); rmErr != nil {
			l.logger.Errorw("failed to create JobRun in the errored state", append(lr.ForLogger(), "error", rmErr))
		}
		return
	}
	_, err = l.runManager.Create(l.jobID, l.initiator, lr.BlockNumber(), &rr)
	if err != nil {
		l.logger.Errorw("failed persist JobRun from a RunRequest", append(lr.ForLogger(), "error", err)...)
	}
	return
}

func (l *listener) JobID() *models.ID {
	return l.jobID
}

// Noops

func (l *listener) OnConnect() {
	l.logger.Debugw("connected to LogBroadcaster", "jobID", l.jobID, "initiator", l.initiator.ID)
}

func (l *listener) OnDisconnect() {
	l.logger.Debugw("disconnected from LogBroadcaster", "jobID", l.jobID, "initiator", l.initiator.ID)
}

// Helpers

func (l *listener) isMatchingInitiator(lb eth.LogBroadcast) bool {
	// TODO make this more efficient!
	logTopics := lb.Log().RawLog().Topics
	initiatorTopics := l.initiator.InitiatorParams.Topics
	for _, topics := range initiatorTopics {
		if reflect.DeepEqual(topics, logTopics) {
			return true
		}
	}
	return false
}

func (l *listener) toLogRequest(lb eth.LogBroadcast) models.LogRequest {
	ile := models.InitiatorLogEvent{
		Initiator: *l.initiator,
		Log:       lb.Log().RawLog(),
	}
	return ile.LogRequest()
}

// runManagerCreate is a subset of the RunManager interface that is needed in the subscription.
type runManagerCreate interface {
	Create(jobSpecID *models.ID, initiator *models.Initiator, creationHeight *big.Int, runRequest *models.RunRequest) (*models.JobRun, error)
	CreateErrored(jobSpecID *models.ID, initiator models.Initiator, err error) (*models.JobRun, error)
}

// logBroadcasterRegistration is a subset of eth.LogBroadcaster with only Register and Unregister.
type logBroadcasterRegistration interface {
	Register(address common.Address, listener eth.LogListener) (connected bool)
	Unregister(address common.Address, listener eth.LogListener)
}

// jsLogger is a subset of the global logger
type jsLogger interface {
	Errorw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
}
