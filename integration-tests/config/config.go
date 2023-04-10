package config

var (
	BaseOCRP2PV1Config = `[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690`

	BaseOCR2Config = `[Feature]
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
LimitDefault = 3_500_000
PriceMax = 100000000000
FeeCapDefault = 100000000000`

	BaseMercuryTomlConfig = `[Feature]
LogPoller = true

[Log]
Level = 'debug'
JSONConsole = true

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false

[WebServer.TLS]
HTTPSPort = 0

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 100

[JobPipeline]
MaxSuccessfulRuns = 0

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']`

	TelemetryIngressConfig = `[TelemetryIngress]
UniConn = false 
Logging = true 
ServerPubKey = '8fa807463ad73f9ee855cfd60ba406dcf98a2855b3dd8af613107b0f6890a707'
URL = 'oti:1337' 
`
)
