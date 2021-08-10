package webhook

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
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

func NewAuthorizer(db *gorm.DB, user *models.User, ei *models.ExternalInitiator) Authorizer {
	if user != nil {
		return &alwaysAuthorizer{}
	} else if ei != nil {
		return NewEIAuthorizer(db, *ei)
	}
	return &neverAuthorizer{}
}

type eiAuthorizer struct {
	db *gorm.DB
	ei models.ExternalInitiator
}

func NewEIAuthorizer(db *gorm.DB, ei models.ExternalInitiator) *eiAuthorizer {
	return &eiAuthorizer{db, ei}
}

func (ea *eiAuthorizer) CanRun(ctx context.Context, config AuthorizerConfig, jobUUID uuid.UUID) (can bool, err error) {
	if !config.FeatureExternalInitiators() {
		return false, nil
	}
	row := ea.db.WithContext(ctx).Raw(`
SELECT EXISTS (
	SELECT 1 FROM external_initiator_webhook_specs
	JOIN jobs ON external_initiator_webhook_specs.webhook_spec_id = jobs.webhook_spec_id
	AND jobs.external_job_id = ?
	AND external_initiator_webhook_specs.external_initiator_id = ?
)`, jobUUID, ea.ei.ID).Row()

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
