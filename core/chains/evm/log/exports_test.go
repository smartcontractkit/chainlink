package log

import (
	"math/big"

	"github.com/smartcontractkit/chainlink/core/logger"
)

func NewRegistrations(logger logger.Logger, evmChainID big.Int) *registrations {
	return newRegistrations(logger, evmChainID)
}

func RegistrationsAddSubscriber(r *registrations, sub *subscriber) bool {
	return r.addSubscriber(sub)
}
func RegistrationsRmSubscriber(r *registrations, sub *subscriber) bool {
	return r.removeSubscriber(sub)
}

func NewSubscriber(l Listener, opts ListenerOpts) *subscriber {
	return &subscriber{l, opts}
}
