package webhook

import (
	"crypto/subtle"
	"database/sql"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"gopkg.in/guregu/null.v4"
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

func (ea *eiAuthorizer) CanRun(jobUUID uuid.UUID) (bool, error) {
	var eiName null.String
	row := ea.db.Raw(`
SELECT external_initiator_name FROM webhook_specs
JOIN jobs ON webhook_specs.id = jobs.webhook_spec_id
AND jobs.external_job_id = ?`, jobUUID).Row()

	err := row.Scan(&eiName)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if eiName.Valid {
		return subtle.ConstantTimeCompare([]byte(ea.ei.Name), []byte(eiName.String)) == 1, nil
	}
	// nil external_initiator_name means this webhook spec can be run by any EI
	return true, nil
}

func (*alwaysAuthorizer) CanRun(uuid.UUID) (bool, error) {
	return true, nil
}

func (*neverAuthorizer) CanRun(uuid.UUID) (bool, error) {
	return false, nil
}
