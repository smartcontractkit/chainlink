package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
)

// JobSpecNotice is sent to the External Initiator when JobSpecs are created.
type JobSpecNotice struct {
	JobID  *models.ID  `json:"jobId"`
	Type   string      `json:"type"`
	Params models.JSON `json:"params,omitempty"`
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

func newNotifyHTTPRequest(jsn JobSpecNotice, ei models.ExternalInitiator) (*http.Request, error) {
	buf, err := json.Marshal(jsn)
	if err != nil {
		return nil, errors.Wrap(err, "new Job Spec notification")
	}
	req, err := http.NewRequest(http.MethodPost, ei.URL.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	setHeaders(req, ei)
	return req, nil
}

// NewExternalInitiatorManager returns the concrete externalInitiatorManager
func NewExternalInitiatorManager() *externalInitiatorManager {
	return &externalInitiatorManager{}
}

type externalInitiatorManager struct{}

// Notify sends a POST notification to the External Initiator
// responsible for initiating the Job Spec.
func (externalInitiatorManager) Notify(
	js models.JobSpec,
	store *store.Store,
) error {
	initrs := js.InitiatorsFor(models.InitiatorExternal)
	if len(initrs) > 1 {
		return errors.New("must have one or less External Initiators")
	}
	if len(initrs) == 0 {
		return nil
	}
	initr := initrs[0]

	ei, err := store.FindExternalInitiatorByName(initr.Name)
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

	req, err := newNotifyHTTPRequest(*notice, ei)
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

func (externalInitiatorManager) DeleteJob(db *gorm.DB, jobID *models.ID) error {
	ei, err := findExternalInitiatorForJob(db, jobID)
	if err != nil {
		return errors.Wrapf(err, "error looking up external initiator for job with id %s", jobID)
	}
	if ei == nil {
		return nil
	}
	if ei.URL == nil {
		return nil
	}
	req, err := newDeleteJobFromExternalInitiatorHTTPRequest(*ei, jobID)
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

func newDeleteJobFromExternalInitiatorHTTPRequest(ei models.ExternalInitiator, id *models.ID) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", ei.URL.String(), id), bytes.NewBuffer(nil))
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

func findExternalInitiatorForJob(db *gorm.DB, id *models.ID) (exi *models.ExternalInitiator, err error) {
	exi = new(models.ExternalInitiator)
	err = db.
		Joins("INNER JOIN initiators ON initiators.name = external_initiators.name").
		Where("initiators.job_spec_id = ?", id).
		First(exi).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return exi, err
}

type NullExternalInitiatorManager struct{}

func (NullExternalInitiatorManager) Notify(models.JobSpec, *store.Store) error {
	return nil
}

func (NullExternalInitiatorManager) DeleteJob(db *gorm.DB, jobID *models.ID) error {
	return nil
}
