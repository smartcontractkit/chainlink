package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smartcontractkit/chainlink/core/logger"
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(ExternalInitiatorAccessKeyHeader, ei.OutgoingToken)
	req.Header.Set(ExternalInitiatorSecretHeader, ei.OutgoingSecret)
	return req, nil
}

// NotifyExternalInitiator sends a POST notification to the External Initiator
// responsible for initiating the Job Spec.
func NotifyExternalInitiator(
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
