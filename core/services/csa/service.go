package csa

import (
	"context"
	"errors"

	storeorm "github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore

type Service interface {
	CreateCSAKey() (*CSAKey, error)
	ListCSAKeys() ([]CSAKey, error)
	CountCSAKeys() (int64, error)
}

type service struct {
	cfg          *storeorm.Config
	orm          ORM
	scryptParams utils.ScryptParams
}

func NewService(cfg *storeorm.Config, orm ORM, scryptParams utils.ScryptParams) Service {
	return &service{
		cfg:          cfg,
		orm:          orm,
		scryptParams: scryptParams,
	}
}

// CreateCSAKey creates a new CSA key
func (s *service) CreateCSAKey() (*CSAKey, error) {
	// Ensure you can only have one CSA at a time. This is a temporary
	// restriction until we are able to handle multiple CSA keys in the
	// communication channel
	count, err := s.orm.CountCSAKeys()
	if err != nil {
		return nil, err
	}

	if count >= 1 {
		return nil, errors.New("can only have 1 CSA key")
	}

	key, err := NewCSAKey(s.cfg.GetKeystorePassword(), s.scryptParams)
	if err != nil {
		return nil, err
	}

	id, err := s.orm.CreateCSAKey(context.Background(), key)
	if err != nil {
		return nil, err
	}

	key, err = s.orm.GetCSAKey(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// ListCSAKeys lists all CSA keys.
func (s *service) ListCSAKeys() ([]CSAKey, error) {
	return s.orm.ListCSAKeys(context.Background())
}

// CountCSAKeys counts the total number of CSA keys.
func (s *service) CountCSAKeys() (int64, error) {
	return s.orm.CountCSAKeys()
}
