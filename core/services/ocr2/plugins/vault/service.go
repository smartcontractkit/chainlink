package vault

import "context"

type Service struct{}

func (s *Service) Start(ctx context.Context) error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func NewService() *Service {
	return &Service{}
}
