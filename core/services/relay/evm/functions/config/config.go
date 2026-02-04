package config

type RelayConfig struct {
	MinIncomingConfirmations        uint32 `json:"minIncomingConfirmations"`
	MinRequestConfirmations         uint32 `json:"minRequestConfirmations"`
	MinResponseConfirmations        uint32 `json:"minResponseConfirmations"`
	LogPollerCacheDurationSec       uint32 `json:"logPollerCacheDurationSec"` // Duration to cache previously detected request or response logs such that they can be filtered when calling logpoller_wrapper.LatestEvents()
	PastBlocksToPoll                uint32 `json:"pastBlocksToPoll"`
	DONID                           string `json:"donID"`
	ContractVersion                 uint32 `json:"contractVersion"`
	ContractUpdateCheckFrequencySec uint32 `json:"contractUpdateCheckFrequencySec"`
}
