package automation

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

type Configurator struct {
	Automation *Automation `toml:"automation"`
}

type Automation struct {
}

type EVMConfigData struct {
	LogPollInterval     string
	ChainID             string
	WsUrl               string
	HttpUrl             string
	LinkContractAddress *string
	FinalityDepth       *int
	FinalityTag         *bool
	SafeTagSupported    *bool
	HeadTracker         *HeadTrackerData
	GasEstimator        *GasEstimatorData
}

type HeadTrackerData struct {
	HistoryDepth int
}

type GasEstimatorData struct {
	Mode         string
	LimitDefault int64
}

func NewAutomationConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.Automation = cfg.Automation
	return nil
}

func (m *Configurator) Store(path string) error {
	if err := products.Store(".", m); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateCLNodesBlockchainConfig(ctx context.Context, bc *blockchain.Input) (string, error) {
	L.Info().Msg("Applying default CL nodes configuration")
	// configure node set and generate CL nodes configs
	config := `
	[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[Log]
Level = 'debug'
JSONConsole = true

[Log.File]
MaxSize = '0b'

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
HTTPWriteTimeout = '3m'
SessionTimeout = '999h0m0s'

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']
AnnounceAddresses = ['0.0.0.0:6690']`

	_ = `

LogPollInterval="500ms"
BackupLogPollerBlockDelay = 0
FinalityDepth = 10
FinalityTagEnabled = false`

	netConfigTemplate := `
[[EVM]]
AutoCreateKey = true
MinContractPayment = 0
BlockBackfillDepth = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.0000001 link'

ChainID = '{{.ChainID}}'
LogPollInterval = '{{.LogPollInterval}}'

{{- if .BackupLogPollerBlockDelay}}
BackupLogPollerBlockDelay = {{.BackupLogPollerBlockDelay}}
{{- end}}
{{- if .LinkContractAddress}}
LinkContractAddress = '{{.LinkContractAddress}}'
{{- end}}

{{- if .FinalityDepth}}
FinalityDepth = {{.FinalityDepth}}
{{- end}}
{{- if .FinalityTag}}
FinalityTag = {{.FinalityTag}}
{{- end}}
{{- if .SafeTagSupported}}
SafeTagSupported = {{.SafeTagSupported}}
{{- end}}

{{- if .HeadTracker}}
[HeadTracker]
HistoryDepth = {{.HeadTracker.HistoryDepth}}
{{- end}}
{{- if .GasEstimator}}
[GasEstimator]
Mode = '{{.GasEstimator.Mode}}'
LimitDefault = {{.GasEstimator.LimitDefault}}
{{- end}}

[[EVM.Nodes]]
Name = 'default'
WsUrl = '{{.WsUrl}}'
HttpUrl = '{{.HttpUrl}}'
`

	tmpl, err := template.New("evmConfig").Parse(netConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := EVMConfigData{
		LinkContractAddress: "", // TODO: set from bc.Out or other source
		ChainID:             bc.Out.ChainID,
		FinalityDepth:       nil, // TODO: set to &value if needed
		FinalityTag:         nil, // TODO: set to &value if needed
		WsUrl:               bc.Out.Nodes[0].InternalWSUrl,
		HttpUrl:             bc.Out.Nodes[0].InternalHTTPUrl,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	L.Info().Msg("Nodes network configuration is finished")
	return buf.String(), nil
}
