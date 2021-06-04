package feeds

import (
	"context"
)

//go:generate mockery --name Service --output ./mocks/ --case=underscore

type Service interface {
	CountManagerServices() (int64, error)
	GetManagerService(id int32) (*ManagerService, error)
	ListManagerServices() ([]ManagerService, error)
	RegisterManagerService(ms *ManagerService) (int32, error)
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

// RegisterManagerService creates a CSA key and registers a new ManagerService.
func (s *service) RegisterManagerService(ms *ManagerService) (int32, error) {
	return s.orm.CreateManagerService(context.Background(), ms)
}

// ListManagerServices lists all the manager services.
func (s *service) ListManagerServices() ([]ManagerService, error) {
	return s.orm.ListManagerServices(context.Background())
}

// GetManagerService gets a manager service by id.
func (s *service) GetManagerService(id int32) (*ManagerService, error) {
	return s.orm.GetManagerService(context.Background(), id)
}

// CountManagerServices gets the total number of manager services
func (s *service) CountManagerServices() (int64, error) {
	return s.orm.Count()
}
