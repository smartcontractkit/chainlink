package types

import "context"

type ForwarderManager[ADDR any] interface {
	Name() string
	Start(ctx context.Context) error
	// Name change b/c no distinction b/w EOA and contracts for some chains
	GetForwarderFor(addr ADDR) (forwarder ADDR, err error)
	GetForwardedPayload(dest ADDR, origPayload []byte) ([]byte, error)
	Close() error
	HealthReport() map[string]error
}
