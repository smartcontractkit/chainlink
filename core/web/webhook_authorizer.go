package web

import (
	"context"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
)

// WebhookRunAuthorizerConfig is satisfied by JobPipeline config.
type WebhookRunAuthorizerConfig interface {
	ExternalInitiatorsEnabled() bool
}

// WebhookRunAuthorizer authorizes legacy webhook runs via external job UUID.
type WebhookRunAuthorizer interface {
	CanRun(ctx context.Context, config WebhookRunAuthorizerConfig, jobUUID uuid.UUID) (bool, error)
}

var (
	_ WebhookRunAuthorizer = &webhookEIAuthorizer{}
	_ WebhookRunAuthorizer = &webhookAlwaysAuthorizer{}
	_ WebhookRunAuthorizer = &webhookNeverAuthorizer{}
)

func NewWebhookRunAuthorizer(ds sqlutil.DataSource, user *sessions.User, ei *bridges.ExternalInitiator) WebhookRunAuthorizer {
	if user != nil {
		return &webhookAlwaysAuthorizer{}
	} else if ei != nil {
		return newWebhookEIAuthorizer(ds, *ei)
	}
	return &webhookNeverAuthorizer{}
}

type webhookEIAuthorizer struct {
	ds sqlutil.DataSource
	ei bridges.ExternalInitiator
}

func newWebhookEIAuthorizer(ds sqlutil.DataSource, ei bridges.ExternalInitiator) *webhookEIAuthorizer {
	return &webhookEIAuthorizer{ds: ds, ei: ei}
}

func (ea *webhookEIAuthorizer) CanRun(ctx context.Context, config WebhookRunAuthorizerConfig, jobUUID uuid.UUID) (can bool, err error) {
	if !config.ExternalInitiatorsEnabled() {
		return false, nil
	}
	row := ea.ds.QueryRowxContext(ctx, `
SELECT EXISTS (
	SELECT 1 FROM external_initiator_webhook_specs
	JOIN jobs ON external_initiator_webhook_specs.webhook_spec_id = jobs.webhook_spec_id
	AND jobs.external_job_id = $1
	AND external_initiator_webhook_specs.external_initiator_id = $2
)`, jobUUID, ea.ei.ID)

	err = row.Scan(&can)
	if err != nil {
		return false, err
	}
	return can, nil
}

type webhookAlwaysAuthorizer struct{}

func (*webhookAlwaysAuthorizer) CanRun(context.Context, WebhookRunAuthorizerConfig, uuid.UUID) (bool, error) {
	return true, nil
}

type webhookNeverAuthorizer struct{}

func (*webhookNeverAuthorizer) CanRun(context.Context, WebhookRunAuthorizerConfig, uuid.UUID) (bool, error) {
	return false, nil
}
