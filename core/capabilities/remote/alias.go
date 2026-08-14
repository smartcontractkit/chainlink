// Package remote re-exports the DON-to-DON dispatcher implementation now living in
// capabilities/libs/x/don2don, so callers can keep importing this path unchanged.
package remote

import (
	don2don "github.com/smartcontractkit/capabilities/libs/x/don2don"
)

var ErrReceiverExists = don2don.ErrReceiverExists

type (
	CombinedClient      = don2don.CombinedClient
	DispatcherConfig    = don2don.DispatcherConfig
	DispatcherRateLimit = don2don.DispatcherRateLimit
	ParallelExecutor    = don2don.ParallelExecutor
	TriggerPublisher    = don2don.TriggerPublisher
	TriggerSubscriber   = don2don.TriggerSubscriber
)

var (
	NewCombinedClient    = don2don.NewCombinedClient
	NewDispatcher        = don2don.NewDispatcher
	NewParallelExecutor  = don2don.NewParallelExecutor
	NewTriggerPublisher  = don2don.NewTriggerPublisher
	NewTriggerSubscriber = don2don.NewTriggerSubscriber

	ValidateMessage   = don2don.ValidateMessage
	ToPeerID          = don2don.ToPeerID
	SanitizeLogString = don2don.SanitizeLogString
)
