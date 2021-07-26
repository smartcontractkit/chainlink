package webhook

import (
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gorm.io/gorm"
)

type (
	eiAuthorizer struct {
		db *gorm.DB
		ei models.ExternalInitiator
	}

	alwaysAuthorizer struct{}
	neverAuthorizer  struct{}

	Authorizer interface {
		CanRun(jobUUID uuid.UUID) (bool, error)
	}
)

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

func NewEIAuthorizer(db *gorm.DB, ei models.ExternalInitiator) *eiAuthorizer {
	return &eiAuthorizer{db, ei}
}

func (ea *eiAuthorizer) CanRun(jobUUID uuid.UUID) (can bool, err error) {
	row := ea.db.Raw(`
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

func (*alwaysAuthorizer) CanRun(uuid.UUID) (bool, error) {
	return true, nil
}

func (*neverAuthorizer) CanRun(uuid.UUID) (bool, error) {
	return false, nil
}
