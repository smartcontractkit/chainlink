package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	iservicesdk "github.com/irisnet/service-sdk-go"
	iservice "github.com/irisnet/service-sdk-go/service"
	"github.com/irisnet/service-sdk-go/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type IritaServiceTracker struct {
	Service      *Service
	store        *store.Store
	runManager   RunManager
	startedMutex sync.RWMutex
	started      bool
}

func NewIritaServiceTracker(
	client iservicesdk.ServiceClient,
	store *store.Store,
	runManager RunManager,
) *IritaServiceTracker {
	return &IritaServiceTracker{
		Service: &Service{
			Store:            store,
			Client:           client,
			RunManager:       runManager,
			jobSubscriptions: map[string]types.Subscription{},
			jobsMutex:        &sync.RWMutex{},
		},
		store:      store,
		runManager: runManager,
	}
}

func (ist *IritaServiceTracker) Start() error {
	ist.startedMutex.Lock()
	defer ist.startedMutex.Unlock()
	if ist.started {
		return errors.New("Irita service tracker already started")
	}

	if err := ist.Service.Start(); err != nil {
		return err
	}

	ist.started = true

	return ist.store.Jobs(
		func(j *models.JobSpec) bool {
			ist.addJob(j)
			return true
		},
		models.InitiatorIritaLog,
	)
}

func (ist *IritaServiceTracker) Stop() {
	ist.startedMutex.Lock()
	defer ist.startedMutex.Unlock()
	if ist.started {
		ist.Service.Stop()
		ist.started = false
	}
}

func (ist *IritaServiceTracker) addJob(job *models.JobSpec) {
	ist.Service.AddJob(*job)
}

func (ist *IritaServiceTracker) AddJob(job models.JobSpec) {
	ist.startedMutex.RLock()
	defer ist.startedMutex.RUnlock()
	if !ist.started {
		return
	}
	ist.addJob(&job)
}

func (ist *IritaServiceTracker) RemoveJob(ID *models.ID) error {
	return ist.Service.RemoveJob(ID)
}

type Service struct {
	Client           iservicesdk.ServiceClient
	Store            *store.Store
	RunManager       RunManager
	jobSubscriptions map[string]types.Subscription
	jobsMutex        *sync.RWMutex
	done             chan struct{}
}

func (s *Service) Start() error {
	s.done = make(chan struct{})
	return nil
}

func (s *Service) Stop() {
	close(s.done)
}

func (s *Service) AddJob(job models.JobSpec) {
	for _, initiator := range job.InitiatorsFor(models.InitiatorIritaLog) {
		go s.Subscribe(initiator, job)
	}
}

func (s *Service) RemoveJob(ID *models.ID) error {
	s.jobsMutex.Lock()
	sub, ok := s.jobSubscriptions[ID.String()]
	delete(s.jobSubscriptions, ID.String())
	numberJobSubscriptions.Set(float64(len(s.jobSubscriptions)))
	s.jobsMutex.Unlock()

	if !ok {
		return fmt.Errorf("JobSubscriber#RemoveJob: job %s not found", ID)
	}
	_ = s.Client.Unsubscribe(sub)
	return nil
}

func (s *Service) addSubscription(jobID *models.ID, sub types.Subscription) {
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()
	s.jobSubscriptions[jobID.String()] = sub
}

func (s *Service) Subscribe(initiator models.Initiator, job models.JobSpec) {
	builder := types.NewEventQueryBuilder().AddCondition(
		types.NewCond(
			"new_batch_request_provider",
			"provider",
		).EQ(
			types.EventValue(initiator.IritaServiceProvider),
		),
	).AddCondition(
		types.NewCond(
			"new_batch_request",
			"service_name",
		).EQ(
			types.EventValue(initiator.IritaServiceName),
		),
	)

	ch := make(chan iservice.QueryServiceRequestResponse)

	providerAcc, _ := types.AccAddressFromBech32(initiator.IritaServiceProvider)

	sub, err := s.Client.SubscribeNewBlock(
		builder,
		func(block types.EventDataNewBlock) {
			for _, request := range s.GetServiceResquest(
				block.ResultEndBlock.Events,
				initiator.IritaServiceName,
				providerAcc,
			) {
				ch <- request
			}
		},
	)
	if err != nil {
		panic(err)
	}

	s.addSubscription(job.ID, sub)

	for {
		select {
		case request := <-ch:
			now := time.Now()
			if !job.Started(now) || job.Ended(now) {
				return
			}

			jobRun, err := s.RunManager.Create(
				job.ID,
				&initiator,
				nil,
				&models.RunRequest{},
			)
			if err != nil {
				logger.Error(err.Error())
				return
			}

			// Add to memory
			store.AddToMemory(jobRun.GetID(), &store.ServiceRequset{
				RequestResponse: request,
				Provider:        initiator.IritaServiceProvider,
			})

		case <-s.done:
			_ = s.Client.Unsubscribe(sub)
			println("done")
		}
	}
}

func (s *Service) GetServiceResquest(
	events types.StringEvents,
	serviceName string,
	provider types.AccAddress,
) []iservice.QueryServiceRequestResponse {
	providerAddr, err := events.GetValue("new_batch_request_provider", "provder")
	if err != nil {
		return nil
	}

	if providerAddr != provider.String() {
		return nil
	}

	svcName, err := events.GetValue("new_batch_request_provider", "service_name")
	if err != nil {
		return nil
	}

	if svcName != serviceName {
		return nil
	}

	var ids []string
	reqIDsStr, err := events.GetValue("new_batch_request_provider", "requests")
	if err != nil {
		return nil
	}

	var idsTemp []string
	_ = json.Unmarshal([]byte(reqIDsStr), &idsTemp)
	ids = append(ids, idsTemp...)

	var requests []iservice.QueryServiceRequestResponse

	for _, reqID := range ids {
		request, _ := s.Client.QueryServiceRequest(reqID)
		requests = append(requests, request)
	}

	return requests
}
