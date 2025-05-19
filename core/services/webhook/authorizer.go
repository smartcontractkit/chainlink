package webhook

import (
	"context"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
)

type AuthorizerConfig interface {
	ExternalInitiatorsEnabled() bool
}

type Authorizer interface {
	CanRun(ctx context.Context, config AuthorizerConfig, jobUUID uuid.UUID) (bool, error)
}

var (
	_ Authorizer = &eiAuthorizer{}
	_ Authorizer = &alwaysAuthorizer{}
	_ Authorizer = &neverAuthorizer{}
)

func NewAuthorizer(ds sqlutil.DataSource, user *sessions.User, ei *bridges.ExternalInitiator) Authorizer {
	if user != nil {
		return &alwaysAuthorizer{}
	} else if ei != nil {
		return NewEIAuthorizer(ds, *ei)
	}
	return &neverAuthorizer{}
}

type eiAuthorizer struct {
	ds sqlutil.DataSource
	ei bridges.ExternalInitiator
}

func NewEIAuthorizer(ds sqlutil.DataSource, ei bridges.ExternalInitiator) *eiAuthorizer {
	return &eiAuthorizer{ds, ei}
}

func (ea *eiAuthorizer) CanRun(ctx context.Context, config AuthorizerConfig, jobUUID uuid.UUID) (can bool, err error) {
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

type alwaysAuthorizer struct{}

func (*alwaysAuthorizer) CanRun(context.Context, AuthorizerConfig, uuid.UUID) (bool, error) {
	return true, nil
}

type neverAuthorizer struct{}

func (*neverAuthorizer) CanRun(context.Context, AuthorizerConfig, uuid.UUID) (bool, error) {
	return false, nil
}
