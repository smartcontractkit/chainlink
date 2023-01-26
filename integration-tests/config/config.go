package config

var (
	BaseOCRP2PV1Config = `[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`

	BaseOCR2VRFTomlConfig = `[Feature]
LogPoller = true

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
AnnounceAddresses = ["0.0.0.0:6690"]
ListenAddresses = ["0.0.0.0:6690"]`

	DefaultOCR2VRFNetworkDetailTomlConfig = `FinalityDepth = 5
[EVM.GasEstimator]
LimitDefault = 1400000
PriceMax = 100000000000
FeeCapDefault = 100000000000`
)
