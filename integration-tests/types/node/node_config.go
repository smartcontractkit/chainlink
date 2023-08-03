package node

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var NodeConfigTemplate = `RootDir = '/home/chainlink'

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

[[EVM.Nodes]]
WSURL = "{{ .EVM.WsUrl }}"
HTTPURL = "{{ .EVM.HttpUrl }}"
Name = "1337_primary_local_0"
SendOnly = false

[Feature]
LogPoller = true
FeedsManager = true
UICSAKeys = true

[OCR]
Enabled = true

[P2P]
[P2P.V1]
Enabled = true
ListenIP = '0.0.0.0'
ListenPort = 6690
`

type NodeConfigOpts struct {
	EVM struct {
		HttpUrl string
		WsUrl   string
	}
}

func ExecuteNodeConfigTemplate(opts NodeConfigOpts) (string, error) {
	t, err := template.New("node-config").Parse(NodeConfigTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, opts)

	return buf.String(), err
}
