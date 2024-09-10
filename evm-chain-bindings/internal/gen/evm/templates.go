package evm

const ContractBindingTemplate = `// Code generated evm-bindings; DO NOT EDIT.

package {{.PackageName}}

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"math/big"
)

//CodeDetails methods inputs and outputs structs
type {{.CodeDetails.ContractStructName}} struct {
	ContractReader types.ContractReader
	ChainWriter types.ChainWriter
}

{{ range $key, $value := .CodeDetails.Structs }}
type {{ $key }} struct {
	{{- range $value.Fields }}
	{{ .Name }} {{ .Type.Print }}
	{{- end }}
}
{{ end }}

{{range .CodeDetails.Functions}}
{{- if .RequiresTransaction }}
func (b {{$.CodeDetails.ContractStructName}}) {{.Name}}(ctx context.Context, {{ if .Input }}input {{.Input.Name}},{{ end }} txId string, toAddress string, meta *types.TxMeta{{ if .IsPayable }}, value *big.Int{{ end }}) error {
	return b.ChainWriter.SubmitTransaction(ctx, "{{$.CodeDetails.ContractName}}", "{{.Name}}", {{ if .Input }}input{{ else }}nil{{ end}}, txId, toAddress, meta, {{ if .IsPayable }}value{{ end }}big.NewInt(0))
}
{{- else }}
func (b {{$.CodeDetails.ContractStructName}}) {{.Name}}(ctx context.Context, {{ if .Input }}input {{.Input.Name}},{{ end }} confidence primitives.ConfidenceLevel) {{ if .Output }}({{ .Output.Print }}, error){{ else }}error {{ end }}{
	{{ if .Output.IsStruct }}output := {{.Output.Print }}{}{{ else }}var output {{ .Output.Print }}{{ end }}
	err := b.ContractReader.GetLatestValue(ctx, "{{$.CodeDetails.ContractName}}", "{{.Name}}", confidence, {{ if .Input }}input{{ else }}nil{{ end }}, &output)
	return output, err	
}
{{- end }}
{{end}}
`

const ChainReaderConfigFactoryTemplate = `package {{.PackageName}}

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func NewChainReaderConfig() types.ChainReaderConfig {
    chainReaderConfig := {{- template "fieldInit" .ChainReaderConfig }}
	return chainReaderConfig
}

func NewChainWriterConfig(maxGasPrice assets.Wei, defaultGasPrice uint64) types.ChainWriterConfig {
    chainWriterConfig := {{- template "fieldInit" .ChainWriterConfig }}
	chainWriterConfig.MaxGasPrice = &maxGasPrice
	for _, contract := range chainWriterConfig.Contracts{
		for _, chainWriterDefinition := range contract.Configs {
			chainWriterDefinition.GasLimit = defaultGasPrice
		}
	}
	return chainWriterConfig
}

{{- define "fieldInit" -}}
{{ if .IsPointer }}&{{ end -}}
{{- if .IsStruct -}}
{{ .GoType }}{
	{{ range $key, $value := .Fields }}{{ $key }}: {{ template "fieldInit" $value }},
	{{ end -}} 
	}
{{- else if .IsMap }}map[{{ .KeyType }}]{{ .ValueType }}{
	{{ range $key, $value := .MapValues }}"{{ $key }}": {{ template "fieldInit" $value }},{{- end }}
	}
{{- else if .IsArray }}{{ .GoType }}{
	{{ range $value := .ArrayValues }}{{ template "fieldInit" $value }},{{ end }}
}
{{- else }}{{ .Value }}{{- end }}
{{- end -}}
`

//const ChainReaderConfigFactoryTemplate = `// Code generated evm-bindings; DO NOT EDIT.
//package {{.PackageName}}
//
//import (
//	"github.com/goccy/go-json"
//	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
//)
//
//func NewChainReaderConfig() (types.ChainReaderConfig, error) {
//	var jsonData = []byte({{ .ChainReaderConfigJson }})
//	chainReaderConfig := types.ChainReaderConfig{}
//	err := json.Unmarshal(jsonData, chainReaderConfig)
//	return chainReaderConfig, err
//}
//`
