package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"go.uber.org/multierr"

	uuid "github.com/satori/go.uuid"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name ExternalInitiatorManager --output ./mocks/ --case=underscore

// ExternalInitiatorManager manages HTTP requests to remote external initiators
type ExternalInitiatorManager interface {
	Notify(webhookSpecID int32) error
	DeleteJob(webhookSpecID int32) error
	FindExternalInitiatorByName(name string) (models.ExternalInitiator, error)
}

//go:generate mockery --name HTTPClient --output ./mocks/ --case=underscore
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type externalInitiatorManager struct {
	db         *gorm.DB
	httpclient HTTPClient
}

var _ ExternalInitiatorManager = (*externalInitiatorManager)(nil)

// NewExternalInitiatorManager returns the concrete externalInitiatorManager
func NewExternalInitiatorManager(db *gorm.DB, httpclient HTTPClient) *externalInitiatorManager {
	return &externalInitiatorManager{db, httpclient}
}

// Notify sends a POST notification to the External Initiator
// responsible for initiating the Job Spec.
func (m externalInitiatorManager) Notify(webhookSpecID int32) error {
	eiWebhookSpecs, jobID, err := m.Load(webhookSpecID)
	if err != nil {
		return err
	}
	for _, eiWebhookSpec := range eiWebhookSpecs {
		ei := eiWebhookSpec.ExternalInitiator
		if ei.URL == nil {
			continue
		}
		notice := JobSpecNotice{
			JobID:  jobID,
			Type:   ei.Name,
			Params: eiWebhookSpec.Spec,
		}
		buf, err := json.Marshal(notice)
		if err != nil {
			return errors.Wrap(err, "new Job Spec notification")
		}
		req, err := newNotifyHTTPRequest(buf, ei)
		if err != nil {
			return errors.Wrap(err, "creating notify HTTP request")
		}
		resp, err := m.httpclient.Do(req)
		if err != nil {
			return errors.Wrap(err, "could not notify '%s' (%s)")
		}
		if err := resp.Body.Close(); err != nil {
			return err
		}
		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			return fmt.Errorf(" notify '%s' (%s) received bad response '%s'", ei.Name, ei.URL, resp.Status)
		}
	}
	return nil
}

func (m externalInitiatorManager) Load(webhookSpecID int32) (eiWebhookSpecs []job.ExternalInitiatorWebhookSpec, jobID uuid.UUID, err error) {
	row := m.db.Raw("SELECT external_job_id FROM jobs WHERE webhook_spec_id = ?", webhookSpecID).Row()
	err = multierr.Combine(
		errors.Wrapf(row.Scan(&jobID), "failed to load job ID from job for webhook spec with ID %d", webhookSpecID),
		errors.Wrapf(m.db.Where("webhook_spec_id = ?", webhookSpecID).Preload("ExternalInitiator").Find(&eiWebhookSpecs).Error, "failed to load external_initiator_webhook_specs for webhook_spec_id %d", webhookSpecID),
	)
	return
}

func (m externalInitiatorManager) DeleteJob(webhookSpecID int32) error {
	eiWebhookSpecs, jobID, err := m.Load(webhookSpecID)
	if err != nil {
		return err
	}
	for _, eiWebhookSpec := range eiWebhookSpecs {
		ei := eiWebhookSpec.ExternalInitiator
		if ei.URL == nil {
			continue
		}

		req, err := newDeleteJobFromExternalInitiatorHTTPRequest(ei, jobID)
		if err != nil {
			return errors.Wrap(err, "creating delete HTTP request")
		}
		resp, err := m.httpclient.Do(req)
		if err != nil {
			return errors.Wrapf(err, "could not delete job from remote external initiator at %s", req.URL)
		}
		if err := resp.Body.Close(); err != nil {
			return err
		}
		if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
			return fmt.Errorf(" notify '%s' (%s) received bad response '%s'", ei.Name, ei.URL, resp.Status)
		}
	}
	return nil
}

func (m externalInitiatorManager) FindExternalInitiatorByName(name string) (models.ExternalInitiator, error) {
	var exi models.ExternalInitiator
	return exi, m.db.First(&exi, "lower(name) = lower(?)", name).Error
}

// JobSpecNotice is sent to the External Initiator when JobSpecs are created.
type JobSpecNotice struct {
	JobID  uuid.UUID   `json:"jobId"`
	Type   string      `json:"type"`
	Params models.JSON `json:"params,omitempty"`
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

func (NullExternalInitiatorManager) Notify(int32) error    { return nil }
func (NullExternalInitiatorManager) DeleteJob(int32) error { return nil }
func (NullExternalInitiatorManager) FindExternalInitiatorByName(name string) (models.ExternalInitiator, error) {
	return models.ExternalInitiator{}, nil
}
