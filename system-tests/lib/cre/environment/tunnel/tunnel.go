package tunnel

import "context"

type EndpointRef struct {
	ComponentID  string
	EndpointName string
	Scheme       string
	Host         string
	Port         int
	OriginalURL  string
}

type TunnelBinding struct {
	EndpointRef
	LocalPort int
	LocalURL  string
}

type Manager interface {
	Start(ctx context.Context, refs []EndpointRef) ([]TunnelBinding, error)
	Stop(ctx context.Context) error
	IsStarted() bool
}

type Provider interface {
	Open(ctx context.Context, ref EndpointRef) (TunnelBinding, error)
	Close(ctx context.Context, binding TunnelBinding) error
	Name() string
}
