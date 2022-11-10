package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --quiet --name ExternalInitiatorManager --output ./mocks/ --case=underscore

// ExternalInitiatorManager manages HTTP requests to remote external initiators
type ExternalInitiatorManager interface {
	Notify(webhookSpecID int32) error
	DeleteJob(webhookSpecID int32) error
	FindExternalInitiatorByName(name string) (bridges.ExternalInitiator, error)
}

//go:generate mockery --quiet --name HTTPClient --output ./mocks/ --case=underscore
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type externalInitiatorManager struct {
	q          pg.Q
	httpclient HTTPClient
}

var _ ExternalInitiatorManager = (*externalInitiatorManager)(nil)

// NewExternalInitiatorManager returns the concrete externalInitiatorManager
func NewExternalInitiatorManager(db *sqlx.DB, httpclient HTTPClient, lggr logger.Logger, cfg pg.QConfig) *externalInitiatorManager {
	namedLogger := lggr.Named("ExternalInitiatorManager")
	return &externalInitiatorManager{
		q:          pg.NewQ(db, namedLogger, cfg),
		httpclient: httpclient,
	}
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
			return fmt.Errorf(" notify '%s' (%s) received bad response '%d: %s'", ei.Name, ei.URL, resp.StatusCode, resp.Status)
		}
	}
	return nil
}

func (m externalInitiatorManager) Load(webhookSpecID int32) (eiWebhookSpecs []job.ExternalInitiatorWebhookSpec, jobID uuid.UUID, err error) {
	err = m.q.Transaction(func(tx pg.Queryer) error {
		if err = tx.Get(&jobID, "SELECT external_job_id FROM jobs WHERE webhook_spec_id = $1", webhookSpecID); err != nil {
			if err = errors.Wrapf(err, "failed to load job ID from job for webhook spec with ID %d", webhookSpecID); err != nil {
				return err
			}
		}
		if err = tx.Select(&eiWebhookSpecs, "SELECT * FROM external_initiator_webhook_specs WHERE external_initiator_webhook_specs.webhook_spec_id = $1", webhookSpecID); err != nil {
			if err = errors.Wrapf(err, "failed to load external_initiator_webhook_specs for webhook_spec_id %d", webhookSpecID); err != nil {
				return err
			}
		}
		if err = m.eagerLoadExternalInitiator(tx, eiWebhookSpecs); err != nil {
			if err = errors.Wrapf(err, "failed to preload ExternalInitiator for webhook_spec_id %d", webhookSpecID); err != nil {
				return err
			}
		}
		return nil
	})

	return
}

func (m externalInitiatorManager) eagerLoadExternalInitiator(q pg.Queryer, txs []job.ExternalInitiatorWebhookSpec) error {
	var ids []int64
	for _, tx := range txs {
		ids = append(ids, tx.ExternalInitiatorID)
	}
	if len(ids) == 0 {
		return nil
	}
	var externalInitiators []bridges.ExternalInitiator
	if err := sqlx.Select(q, &externalInitiators, `SELECT * FROM external_initiators WHERE external_initiators.id = ANY($1);`, pq.Array(ids)); err != nil {
		return err
	}

	eiMap := make(map[int64]bridges.ExternalInitiator)
	for _, externalInitiator := range externalInitiators {
		eiMap[externalInitiator.ID] = externalInitiator
	}

	for i := range txs {
		txs[i].ExternalInitiator = eiMap[txs[i].ExternalInitiatorID]
	}
	return nil
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
			return fmt.Errorf(" delete '%s' (%s) received bad response '%d: %s'", ei.Name, ei.URL, resp.StatusCode, resp.Status)
		}
	}
	return nil
}

func (m externalInitiatorManager) FindExternalInitiatorByName(name string) (bridges.ExternalInitiator, error) {
	var exi bridges.ExternalInitiator
	err := m.q.Get(&exi, "SELECT * FROM external_initiators WHERE lower(external_initiators.name) = lower($1)", name)
	return exi, err
}

// JobSpecNotice is sent to the External Initiator when JobSpecs are created.
type JobSpecNotice struct {
	JobID  uuid.UUID   `json:"jobId"`
	Type   string      `json:"type"`
	Params models.JSON `json:"params,omitempty"`
}

func newNotifyHTTPRequest(buf []byte, ei bridges.ExternalInitiator) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, ei.URL.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	setHeaders(req, ei)
	return req, nil
}

func newDeleteJobFromExternalInitiatorHTTPRequest(ei bridges.ExternalInitiator, jobID uuid.UUID) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", ei.URL.String(), jobID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(req, ei)
	return req, nil
}

func setHeaders(req *http.Request, ei bridges.ExternalInitiator) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(static.ExternalInitiatorAccessKeyHeader, ei.OutgoingToken)
	req.Header.Set(static.ExternalInitiatorSecretHeader, ei.OutgoingSecret)
}

type NullExternalInitiatorManager struct{}

var _ ExternalInitiatorManager = (*NullExternalInitiatorManager)(nil)

func (NullExternalInitiatorManager) Notify(int32) error    { return nil }
func (NullExternalInitiatorManager) DeleteJob(int32) error { return nil }
func (NullExternalInitiatorManager) FindExternalInitiatorByName(name string) (bridges.ExternalInitiator, error) {
	return bridges.ExternalInitiator{}, nil
}
