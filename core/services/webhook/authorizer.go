package webhook

import (
	"context"
	"database/sql"

	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/sessions"
)

type AuthorizerConfig interface {
	FeatureExternalInitiators() bool
}

type Authorizer interface {
	CanRun(ctx context.Context, config AuthorizerConfig, jobUUID uuid.UUID) (bool, error)
}

var (
	_ Authorizer = &eiAuthorizer{}
	_ Authorizer = &alwaysAuthorizer{}
	_ Authorizer = &neverAuthorizer{}
)

func NewAuthorizer(db *sql.DB, user *sessions.User, ei *bridges.ExternalInitiator) Authorizer {
	if user != nil {
		return &alwaysAuthorizer{}
	} else if ei != nil {
		return NewEIAuthorizer(db, *ei)
	}
	return &neverAuthorizer{}
}

type eiAuthorizer struct {
	db *sql.DB
	ei bridges.ExternalInitiator
}

func NewEIAuthorizer(db *sql.DB, ei bridges.ExternalInitiator) *eiAuthorizer {
	return &eiAuthorizer{db, ei}
}

func (ea *eiAuthorizer) CanRun(ctx context.Context, config AuthorizerConfig, jobUUID uuid.UUID) (can bool, err error) {
	if !config.FeatureExternalInitiators() {
		return false, nil
	}
	row := ea.db.QueryRowContext(ctx, `
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
