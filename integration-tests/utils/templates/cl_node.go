package templates

import (
	"github.com/google/uuid"
)

type NodeSecretsTemplate struct {
	PgDbName string
	PgHost   string
	PgPort   string
}

func (c NodeSecretsTemplate) String() (string, error) {
	tpl := `
[Database]
URL = 'postgresql://postgres:test@{{ .PgHost }}:{{ .PgPort }}/{{ .PgDbName }}?sslmode=disable' # Required
AllowSimplePasswords = true

[Password]
Keystore = '................' # Required

[Mercury.Credentials.cred1]
# URL = 'http://host.docker.internal:3000/reports'
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
`
	return MarshallTemplate(c, uuid.NewString(), tpl)
}

type NodeConfigOpts struct {
	EVM struct {
		HttpUrl string
		WsUrl   string
	}
	VRFv2Opts *NodeVRFv2ConfigOpts
}

type NodeVRFv2ConfigOpts struct {
	Key          string
	PriceMaxGwei int
}

func (c NodeConfigOpts) String() (string, error) {
	tpl := `RootDir = '/home/chainlink'

[Database]
MaxIdleConns = 20
MaxOpenConns = 40
MigrateOnStartup = true

[Log]
Level = 'debug'
JSONConsole = true

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
SessionTimeout = '999h0m0s'

[WebServer.TLS]
HTTPSPort = 0

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 100

[[EVM]]
ChainID = "1337"
AutoCreateKey = true
finalityDepth = 1
MinContractPayment = '0'

{{if .VRFv2Opts}}
[[EVM.KeySpecific]]
Key = '{{ .VRFv2Opts.Key }}'

[EVM.KeySpecific.GasEstimator]
PriceMax = '{{ .VRFv2Opts.PriceMaxGwei }} gwei'

[EVM.GasEstimator]
LimitDefault = 3500000
[EVM.Transactions]
MaxQueued = 10000
{{end}}

[[EVM.Nodes]]
WSURL = "{{ .EVM.WsUrl }}"
HTTPURL = "{{ .EVM.HttpUrl }}"
Name = "1337_primary_local_0"
SendOnly = false

[Feature]
LogPoller = true
FeedsManager = true
UICSAKeys = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690
`
	return MarshallTemplate(c, uuid.NewString(), tpl)
}
