package service

import "github.com/smartcontractkit/chainlink/core/services/health"

type (
	Service interface {
		Start() error
		Close() error
		health.Checkable
	}
)
