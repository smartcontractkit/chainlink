package products

import (
	"bytes"
	"text/template"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

type clNodeConfigParams struct {
	ChainID                       string
	InternalBlockchainNodeWSURL   string
	InternalBlockchainNodeHTTPURL string
}

func DefaultLegacyCLNodeConfig(in []*blockchain.Input) (string, error) {
	L.Info().Msg("Applying default CL nodes configuration")
	const configTemplate = `[[EVM]]
LogPollInterval = '1s'
BlockBackfillDepth = 100
ChainID = '{{.ChainID}}'
MinIncomingConfirmations = 1
MinContractPayment = '0.0 link'

[[EVM.Nodes]]
Name = 'default'
WsUrl = '{{.InternalBlockchainNodeWSURL}}'
HttpUrl = '{{.InternalBlockchainNodeHTTPURL}}'

[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[OCR2]
Enabled = true
SimulateTransactions = false
DefaultTransactionQueueDepth = 1

[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']

[Log]
JSONConsole = true
Level = 'debug'

[Pyroscope]
ServerAddress = 'http://pyroscope:4040'
Environment = 'local'

[WebServer]
SessionTimeout = '999h0m0s'
HTTPWriteTimeout = '3m'
SecureCookies = false
HTTPPort = 6688

[WebServer.TLS]
HTTPSPort = 0

[WebServer.RateLimit]
Authenticated = 5000
Unauthenticated = 5000

[JobPipeline]

[JobPipeline.HTTPRequest]
DefaultTimeout = '1m'

[Log.File]
MaxSize = '0b'`

	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return "", err
	}

	node := in[0].Out.Nodes[0]
	chainID := in[0].Out.ChainID

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, clNodeConfigParams{
		ChainID:                       chainID,
		InternalBlockchainNodeWSURL:   node.InternalWSUrl,
		InternalBlockchainNodeHTTPURL: node.InternalHTTPUrl,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
