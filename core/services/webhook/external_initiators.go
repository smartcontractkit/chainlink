package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/job"

	uuid "github.com/satori/go.uuid"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name ExternalInitiatorManager --output ./mocks/ --case=underscore

type (
	// ExternalInitiatorManager manages HTTP requests to remote external initiators
	ExternalInitiatorManager interface {
		Notify(models.JobSpec) error
		NotifyV2(jobID uuid.UUID, initrName string, initrSpec *models.JSON) error
		DeleteJob(jobID models.JobID) error
		DeleteJobV2(job job.Job) error
		FindExternalInitiatorByName(name string) (models.ExternalInitiator, error)
	}

	externalInitiatorManager struct {
		db *gorm.DB
	}
)

var _ ExternalInitiatorManager = (*externalInitiatorManager)(nil)

// NewExternalInitiatorManager returns the concrete externalInitiatorManager
func NewExternalInitiatorManager(db *gorm.DB) *externalInitiatorManager {
	return &externalInitiatorManager{db: db}
}

// Notify sends a POST notification to the External Initiator
// responsible for initiating the Job Spec.
func (m externalInitiatorManager) Notify(js models.JobSpec) error {
	initrs := js.InitiatorsFor(models.InitiatorExternal)
	if len(initrs) > 1 {
		return errors.New("must have one or less External Initiators")
	}
	if len(initrs) == 0 {
		return nil
	}
	initr := initrs[0]

	ei, err := m.FindExternalInitiatorByName(initr.Name)
	if err != nil {
		return errors.Wrap(err, "external initiator")
	}
	if ei.URL == nil {
		return nil
	}
	notice, err := NewJobSpecNotice(initr, js)
	if err != nil {
		return errors.Wrap(err, "new Job Spec notification")
	}
	buf, err := json.Marshal(notice)
	if err != nil {
		return errors.Wrap(err, "new Job Spec notification")
	}
	req, err := newNotifyHTTPRequest(buf, ei)
	if err != nil {
		return errors.Wrap(err, "creating notify HTTP request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not notify '%s' (%s)")
	}
	defer logger.ErrorIfCalling(resp.Body.Close)
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf(" notify '%s' (%s) received bad response '%s'", ei.Name, ei.URL, resp.Status)
	}
	return nil
}

// NotifyV2 sends a POST notification to the External Initiator
// responsible for initiating the Job Spec.
func (m externalInitiatorManager) NotifyV2(
	jobID uuid.UUID,
	initrName string,
	initrSpec *models.JSON,
) error {
	ei, err := m.FindExternalInitiatorByName(initrName)
	if err != nil {
		return errors.Wrap(err, "external initiator")
	} else if ei.URL == nil {
		return nil
	} else if initrSpec == nil {
		return errors.New("body must be defined")
	}

	notice := JobSpecNoticeV2{
		JobID:  jobID,
		Type:   initrName,
		Params: *initrSpec,
	}
	buf, err := json.Marshal(notice)
	if err != nil {
		return errors.Wrap(err, "new Job Spec notification")
	}
	req, err := newNotifyHTTPRequest(buf, ei)
	if err != nil {
		return errors.Wrap(err, "creating notify HTTP request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "could not notify '%s' (%s)")
	}
	defer logger.ErrorIfCalling(resp.Body.Close)
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf(" notify '%s' (%s) received bad response '%s'", ei.Name, ei.URL, resp.Status)
	}
	return nil
}

func (m externalInitiatorManager) DeleteJob(jobID models.JobID) error {
	ei, err := m.findExternalInitiatorForJob(jobID)
	if err != nil {
		return errors.Wrapf(err, "error looking up external initiator for job with id %s", jobID)
	}
	if ei == nil {
		return nil
	}
	if ei.URL == nil {
		return nil
	}
	req, err := newDeleteJobFromExternalInitiatorHTTPRequest(*ei, uuid.UUID(jobID))
	if err != nil {
		return errors.Wrap(err, "creating delete HTTP request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "could not delete job from remote external initiator at %s", req.URL)
	}
	defer logger.ErrorIfCalling(resp.Body.Close)
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf(" notify '%s' (%s) received bad response '%s'", ei.Name, ei.URL, resp.Status)
	}
	return nil
}

func (m externalInitiatorManager) DeleteJobV2(jb job.Job) error {
	if jb.ExternalInitiator == nil {
		return errors.Errorf("no external initiator found for job %d", jb.ID)
	} else if jb.ExternalInitiator.URL == nil {
		return nil
	}

	req, err := newDeleteJobFromExternalInitiatorHTTPRequest(*jb.ExternalInitiator, jb.ExternalJobID)
	if err != nil {
		return errors.Wrap(err, "creating delete HTTP request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "could not delete job from remote external initiator at %s", req.URL)
	}
	defer logger.ErrorIfCalling(resp.Body.Close)
	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return fmt.Errorf(" notify '%s' (%s) received bad response '%s'", jb.ExternalInitiator.Name, jb.ExternalInitiator.URL, resp.Status)
	}
	return nil
}

func (m externalInitiatorManager) FindExternalInitiatorByName(name string) (models.ExternalInitiator, error) {
	var exi models.ExternalInitiator
	return exi, m.db.First(&exi, "lower(name) = lower(?)", name).Error
}

func (m externalInitiatorManager) findExternalInitiatorForJob(id models.JobID) (exi *models.ExternalInitiator, err error) {
	exi = new(models.ExternalInitiator)
	err = m.db.
		Joins("INNER JOIN initiators ON initiators.name = external_initiators.name").
		Where("initiators.job_spec_id = ?", id).
		First(exi).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return exi, err
}

// JobSpecNotice is sent to the External Initiator when JobSpecs are created.
type JobSpecNoticeV2 struct {
	JobID  uuid.UUID   `json:"jobId"`
	Type   string      `json:"type"`
	Params models.JSON `json:"params,omitempty"`
}

// JobSpecNotice is sent to the External Initiator when JobSpecs are created.
type JobSpecNotice struct {
	JobID  models.JobID `json:"jobId"`
	Type   string       `json:"type"`
	Params models.JSON  `json:"params,omitempty"`
}

// NewJobSpecNotice returns a new JobSpec.
func NewJobSpecNotice(initiator models.Initiator, js models.JobSpec) (*JobSpecNotice, error) {
	if initiator.Body == nil {
		return nil, errors.New("body must be defined")
	}
	return &JobSpecNotice{
		JobID:  js.ID,
		Type:   initiator.Type,
		Params: *initiator.Body,
	}, nil
}

func newNotifyHTTPRequest(buf []byte, ei models.ExternalInitiator) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, ei.URL.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	setHeaders(req, ei)
	return req, nil
}

func newDeleteJobFromExternalInitiatorHTTPRequest(ei models.ExternalInitiator, jobID uuid.UUID) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", ei.URL.String(), jobID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(req, ei)
	return req, nil
}

func setHeaders(req *http.Request, ei models.ExternalInitiator) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(static.ExternalInitiatorAccessKeyHeader, ei.OutgoingToken)
	req.Header.Set(static.ExternalInitiatorSecretHeader, ei.OutgoingSecret)
}

type NullExternalInitiatorManager struct{}

var _ ExternalInitiatorManager = (*NullExternalInitiatorManager)(nil)

func (NullExternalInitiatorManager) Notify(models.JobSpec) error {
	return nil
}

func (NullExternalInitiatorManager) NotifyV2(jobUUID uuid.UUID, initrName string, initrSpec *models.JSON) error {
	return nil
}

func (NullExternalInitiatorManager) DeleteJob(jobID models.JobID) error {
	return nil
}

func (NullExternalInitiatorManager) DeleteJobV2(jb job.Job) error {
	return nil
}

func (NullExternalInitiatorManager) FindExternalInitiatorByName(name string) (models.ExternalInitiator, error) {
	return models.ExternalInitiator{}, nil
}
