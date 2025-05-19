package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/static"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// ExternalInitiatorManager manages HTTP requests to remote external initiators
type ExternalInitiatorManager interface {
	Notify(ctx context.Context, webhookSpecID int32) error
	DeleteJob(ctx context.Context, webhookSpecID int32) error
	FindExternalInitiatorByName(ctx context.Context, name string) (bridges.ExternalInitiator, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type externalInitiatorManager struct {
	ds         sqlutil.DataSource
	httpclient HTTPClient
}

var _ ExternalInitiatorManager = (*externalInitiatorManager)(nil)

// NewExternalInitiatorManager returns the concrete externalInitiatorManager
func NewExternalInitiatorManager(ds sqlutil.DataSource, httpclient HTTPClient) *externalInitiatorManager {
	return &externalInitiatorManager{
		ds:         ds,
		httpclient: httpclient,
	}
}

// Notify sends a POST notification to the External Initiator
// responsible for initiating the Job Spec.
func (m *externalInitiatorManager) Notify(ctx context.Context, webhookSpecID int32) error {
	eiWebhookSpecs, jobID, err := m.Load(ctx, webhookSpecID)
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
		req, err := newNotifyHTTPRequest(ctx, buf, ei)
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

func (m *externalInitiatorManager) Load(ctx context.Context, webhookSpecID int32) (eiWebhookSpecs []job.ExternalInitiatorWebhookSpec, jobID uuid.UUID, err error) {
	err = sqlutil.Transact(ctx, func(ds sqlutil.DataSource) *externalInitiatorManager {
		return NewExternalInitiatorManager(ds, m.httpclient)
	}, m.ds, nil, func(tx *externalInitiatorManager) error {
		if err = tx.ds.GetContext(ctx, &jobID, "SELECT external_job_id FROM jobs WHERE webhook_spec_id = $1", webhookSpecID); err != nil {
			if err = errors.Wrapf(err, "failed to load job ID from job for webhook spec with ID %d", webhookSpecID); err != nil {
				return err
			}
		}
		if err = tx.ds.SelectContext(ctx, &eiWebhookSpecs, "SELECT * FROM external_initiator_webhook_specs WHERE external_initiator_webhook_specs.webhook_spec_id = $1", webhookSpecID); err != nil {
			if err = errors.Wrapf(err, "failed to load external_initiator_webhook_specs for webhook_spec_id %d", webhookSpecID); err != nil {
				return err
			}
		}
		if err = tx.eagerLoadExternalInitiator(ctx, eiWebhookSpecs); err != nil {
			if err = errors.Wrapf(err, "failed to preload ExternalInitiator for webhook_spec_id %d", webhookSpecID); err != nil {
				return err
			}
		}
		return nil
	})

	return
}

func (m *externalInitiatorManager) eagerLoadExternalInitiator(ctx context.Context, txs []job.ExternalInitiatorWebhookSpec) error {
	var ids []int64
	for _, tx := range txs {
		ids = append(ids, tx.ExternalInitiatorID)
	}
	if len(ids) == 0 {
		return nil
	}
	var externalInitiators []bridges.ExternalInitiator
	if err := m.ds.SelectContext(ctx, &externalInitiators, `SELECT * FROM external_initiators WHERE external_initiators.id = ANY($1);`, pq.Array(ids)); err != nil {
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

func (m *externalInitiatorManager) DeleteJob(ctx context.Context, webhookSpecID int32) error {
	eiWebhookSpecs, jobID, err := m.Load(ctx, webhookSpecID)
	if err != nil {
		return err
	}
	for _, eiWebhookSpec := range eiWebhookSpecs {
		ei := eiWebhookSpec.ExternalInitiator
		if ei.URL == nil {
			continue
		}

		req, err := newDeleteJobFromExternalInitiatorHTTPRequest(ctx, ei, jobID)
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

func (m *externalInitiatorManager) FindExternalInitiatorByName(ctx context.Context, name string) (bridges.ExternalInitiator, error) {
	var exi bridges.ExternalInitiator
	err := m.ds.GetContext(ctx, &exi, "SELECT * FROM external_initiators WHERE lower(external_initiators.name) = lower($1)", name)
	return exi, err
}

// JobSpecNotice is sent to the External Initiator when JobSpecs are created.
type JobSpecNotice struct {
	JobID  uuid.UUID   `json:"jobId"`
	Type   string      `json:"type"`
	Params models.JSON `json:"params,omitempty"`
}

func newNotifyHTTPRequest(ctx context.Context, buf []byte, ei bridges.ExternalInitiator) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ei.URL.String(), bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}
	setHeaders(req, ei)
	return req, nil
}

func newDeleteJobFromExternalInitiatorHTTPRequest(ctx context.Context, ei bridges.ExternalInitiator, jobID uuid.UUID) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", ei.URL.String(), jobID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
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

func (NullExternalInitiatorManager) Notify(context.Context, int32) error    { return nil }
func (NullExternalInitiatorManager) DeleteJob(context.Context, int32) error { return nil }
func (NullExternalInitiatorManager) FindExternalInitiatorByName(ctx context.Context, name string) (bridges.ExternalInitiator, error) {
	return bridges.ExternalInitiator{}, nil
}
