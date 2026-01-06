package ocr2

/*
This code was extracted from chainlink/v2 and chainlink/v2/deployment
It should be moved either to CTF or to a separate module in chainlink/v2/deployment to expose product interface
without dependency hell
*/

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml/v2"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

// TaskJobSpec represents an OCR2 job that is given to other nodes, meant to communicate with the bootstrap node,
// and provide their answers.
type TaskJobSpec struct {
	OCR2OracleSpec    OracleSpec
	Name              string `toml:"name"`
	JobType           string `toml:"type"`
	MaxTaskDuration   string `toml:"maxTaskDuration"`
	ObservationSource string `toml:"observationSource"`
	ForwardingAllowed bool   `toml:"forwardingAllowed"`
}

// OracleSpec defines the job spec for OCR2 jobs.
// Relay config is chain specific config for a relay (chain adapter).
type OracleSpec struct {
	UpdatedAt                         time.Time            `toml:"-"`
	CreatedAt                         time.Time            `toml:"-"`
	OnchainSigningStrategy            JSONConfig           `toml:"onchainSigningStrategy"`
	FeedID                            *common.Hash         `toml:"feedID"`
	PluginConfig                      JSONConfig           `toml:"pluginConfig"`
	RelayConfig                       JSONConfig           `toml:"relayConfig"`
	PluginType                        types.OCR2PluginType `toml:"pluginType"`
	ChainID                           string               `toml:"chainID"`
	ContractID                        string               `toml:"contractID"`
	Relay                             string               `toml:"relay"`
	P2PV2Bootstrappers                pq.StringArray       `toml:"p2pv2Bootstrappers"`
	OCRKeyBundleID                    null.String          `toml:"ocrKeyBundleID"`
	TransmitterID                     null.String          `toml:"transmitterID"`
	MonitoringEndpoint                null.String          `toml:"monitoringEndpoint"`
	ContractConfigTrackerPollInterval Interval             `toml:"contractConfigTrackerPollInterval"`
	BlockchainTimeout                 Interval             `toml:"blockchainTimeout"`
	ID                                int32                `toml:"-"`
	ContractConfigConfirmations       uint16               `toml:"contractConfigConfirmations"`
	CaptureEATelemetry                bool                 `toml:"captureEATelemetry"`
	CaptureAutomationCustomTelemetry  bool                 `toml:"captureAutomationCustomTelemetry"`
	AllowNoBootstrappers              bool                 `toml:"allowNoBootstrappers"`
}

// JSONConfig is a map for config properties which are encoded as JSON in the database by implementing
// sql.Scanner and driver.Valuer.
type JSONConfig map[string]any

// Interval represents a time.Duration stored as a Postgres interval type.
type Interval time.Duration

type ocr2Config interface {
	DefaultTransactionQueueDepth() uint32
	SimulateTransactions() bool
}

// Type returns the type of the job.
func (o *TaskJobSpec) Type() string { return o.JobType }

// String representation of the job.
func (o *TaskJobSpec) String() (string, error) {
	var feedID string
	if o.OCR2OracleSpec.FeedID != nil {
		feedID = o.OCR2OracleSpec.FeedID.Hex()
	}
	relayConfig, err := toml.Marshal(struct {
		RelayConfig JSONConfig `toml:"relayConfig"`
	}{RelayConfig: o.OCR2OracleSpec.RelayConfig})
	if err != nil {
		return "", fmt.Errorf("failed to marshal relay config: %w", err)
	}
	specWrap := struct {
		PluginConfig             map[string]any
		RelayConfig              string
		OCRKeyBundleID           string
		ObservationSource        string
		ContractID               string
		FeedID                   string
		Relay                    string
		PluginType               string
		Name                     string
		MaxTaskDuration          string
		JobType                  string
		TransmitterID            string
		MonitoringEndpoint       string
		P2PV2Bootstrappers       []string
		BlockchainTimeout        time.Duration
		TrackerSubscribeInterval time.Duration
		TrackerPollInterval      time.Duration
		ContractConfirmations    uint16
		ForwardingAllowed        bool
	}{
		Name:                  o.Name,
		JobType:               o.JobType,
		ForwardingAllowed:     o.ForwardingAllowed,
		MaxTaskDuration:       o.MaxTaskDuration,
		ContractID:            o.OCR2OracleSpec.ContractID,
		FeedID:                feedID,
		Relay:                 o.OCR2OracleSpec.Relay,
		PluginType:            string(o.OCR2OracleSpec.PluginType),
		RelayConfig:           string(relayConfig),
		PluginConfig:          o.OCR2OracleSpec.PluginConfig,
		P2PV2Bootstrappers:    o.OCR2OracleSpec.P2PV2Bootstrappers,
		OCRKeyBundleID:        o.OCR2OracleSpec.OCRKeyBundleID.String,
		MonitoringEndpoint:    o.OCR2OracleSpec.MonitoringEndpoint.String,
		TransmitterID:         o.OCR2OracleSpec.TransmitterID.String,
		BlockchainTimeout:     o.OCR2OracleSpec.BlockchainTimeout.Duration(),
		ContractConfirmations: o.OCR2OracleSpec.ContractConfigConfirmations,
		TrackerPollInterval:   o.OCR2OracleSpec.ContractConfigTrackerPollInterval.Duration(),
		ObservationSource:     o.ObservationSource,
	}
	ocr2TemplateString := `
type                                   = "{{ .JobType }}"
name                                   = "{{.Name}}"
forwardingAllowed                      = {{.ForwardingAllowed}}
{{- if .MaxTaskDuration}}
maxTaskDuration                        = "{{ .MaxTaskDuration }}" {{end}}
{{- if .PluginType}}
pluginType                             = "{{ .PluginType }}" {{end}}
relay                                  = "{{.Relay}}"
schemaVersion                          = 1
contractID                             = "{{.ContractID}}"
{{- if .FeedID}}
feedID                                 = "{{.FeedID}}"
{{end}}
{{- if eq .JobType "offchainreporting2" }}
ocrKeyBundleID                         = "{{.OCRKeyBundleID}}" {{end}}
{{- if eq .JobType "offchainreporting2" }}
transmitterID                          = "{{.TransmitterID}}" {{end}}
{{- if .BlockchainTimeout}}
blockchainTimeout                      = "{{.BlockchainTimeout}}"
{{end}}
{{- if .ContractConfirmations}}
contractConfigConfirmations            = {{.ContractConfirmations}}
{{end}}
{{- if .TrackerPollInterval}}
contractConfigTrackerPollInterval      = "{{.TrackerPollInterval}}"
{{end}}
{{- if .TrackerSubscribeInterval}}
contractConfigTrackerSubscribeInterval = "{{.TrackerSubscribeInterval}}"
{{end}}
{{- if .P2PV2Bootstrappers}}
p2pv2Bootstrappers                     = [{{range .P2PV2Bootstrappers}}"{{.}}",{{end}}]{{end}}
{{- if .MonitoringEndpoint}}
monitoringEndpoint                     = "{{.MonitoringEndpoint}}" {{end}}
{{- if .ObservationSource}}
observationSource                      = """
{{.ObservationSource}}
"""{{end}}
{{if eq .JobType "offchainreporting2" }}
[pluginConfig]{{range $key, $value := .PluginConfig}}
{{$key}} = {{$value}}{{end}}
{{end}}
{{.RelayConfig}}
`
	return MarshallTemplate(specWrap, "OCR2 Job", ocr2TemplateString)
}

// MarshallTemplate Helper to marshall templates.
func MarshallTemplate(jobSpec any, name, templateString string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, jobSpec)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

// Bytes returns the raw bytes.
func (r JSONConfig) Bytes() []byte {
	b, _ := json.Marshal(r)
	return b
}

// Value returns this instance serialized for database storage.
func (r JSONConfig) Value() (driver.Value, error) {
	return json.Marshal(r)
}

// Scan reads the database value and returns an instance.
func (r *JSONConfig) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.Errorf("expected bytes got %T", b)
	}
	return json.Unmarshal(b, &r)
}

func (r JSONConfig) MercuryCredentialName() (string, error) {
	url, ok := r["mercuryCredentialName"]
	if !ok {
		return "", nil
	}
	name, ok := url.(string)
	if !ok {
		return "", fmt.Errorf("expected string mercuryCredentialName but got: %T", url)
	}
	return name, nil
}

func (r JSONConfig) ApplyDefaultsOCR2(cfg ocr2Config) {
	_, ok := r["defaultTransactionQueueDepth"]
	if !ok {
		r["defaultTransactionQueueDepth"] = cfg.DefaultTransactionQueueDepth()
	}
	_, ok = r["simulateTransactions"]
	if !ok {
		r["simulateTransactions"] = cfg.SimulateTransactions()
	}
}

// NewInterval creates Interval for specified duration.
func NewInterval(d time.Duration) *Interval {
	i := new(Interval)
	*i = Interval(d)
	return i
}

func (i Interval) Duration() time.Duration {
	return time.Duration(i)
}

// MarshalText implements the text.Marshaler interface.
func (i Interval) MarshalText() ([]byte, error) {
	return []byte(time.Duration(i).String()), nil
}

// UnmarshalText implements the text.Unmarshaler interface.
func (i *Interval) UnmarshalText(input []byte) error {
	v, err := time.ParseDuration(string(input))
	if err != nil {
		return err
	}
	*i = Interval(v)
	return nil
}

func (i *Interval) Scan(v any) error {
	if v == nil {
		*i = Interval(time.Duration(0))
		return nil
	}
	asInt64, is := v.(int64)
	if !is {
		return errors.Errorf("models.Interval#Scan() wanted int64, got %T", v)
	}
	*i = Interval(time.Duration(asInt64) * time.Nanosecond)
	return nil
}

func (i Interval) Value() (driver.Value, error) {
	return time.Duration(i).Nanoseconds(), nil
}

func (i Interval) IsZero() bool {
	return time.Duration(i) == time.Duration(0)
}
