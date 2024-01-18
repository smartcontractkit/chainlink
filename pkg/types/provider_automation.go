package types

import "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

// AutomationProvider provides components needed for the automation OCR2 plugin.
type AutomationProvider interface {
	PluginProvider
	Registry() automation.Registry
	Encoder() automation.Encoder
	TransmitEventProvider() automation.EventProvider
	BlockSubscriber() automation.BlockSubscriber
	PayloadBuilder() automation.PayloadBuilder
	UpkeepStateStore() automation.UpkeepStateStore
	LogEventProvider() automation.LogEventProvider
	LogRecoverer() automation.LogRecoverer
	UpkeepProvider() automation.ConditionalUpkeepProvider
}
