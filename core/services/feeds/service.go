package feeds

import (
	"context"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore

type Service interface {
	CountManagers() (int64, error)
	GetManager(id int32) (*FeedsManager, error)
	ListManagers() ([]FeedsManager, error)
	RegisterManager(ms *FeedsManager) (int32, error)
}

type service struct {
	orm ORM
}

// NewService constructs a new feeds service
func NewService(orm ORM) Service {
	return &service{
		orm: orm,
	}
}

// RegisterManager registers a new ManagerService.
func (s *service) RegisterManager(ms *FeedsManager) (int32, error) {
	return s.orm.CreateManager(context.Background(), ms)
}

// ListManagerServices lists all the manager services.
func (s *service) ListManagers() ([]FeedsManager, error) {
	return s.orm.ListManagers(context.Background())
}

// GetManager gets a manager service by id.
func (s *service) GetManager(id int32) (*FeedsManager, error) {
	return s.orm.GetManager(context.Background(), id)
}

// CountManagerServices gets the total number of manager services
func (s *service) CountManagers() (int64, error) {
	return s.orm.CountManagers()
}
