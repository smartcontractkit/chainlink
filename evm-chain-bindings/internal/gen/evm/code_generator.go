package evm

import (
	"bytes"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"log"
	"sync"
	"text/template"
)

type contractTemplateData struct {
	PackageName string
	CodeDetails CodeDetails
}

type chainComponentsConfigData struct {
	PackageName       string
	ChainReaderConfig map[string]interface{}
	ChainWriterConfig map[string]interface{}
}

type TemplateRegistry struct {
	contractTemplate          template.Template
	chainReaderConfigTemplate template.Template
}

func (c TemplateRegistry) ProcessContract(data contractTemplateData) ([]byte, error) {
	var buf bytes.Buffer
	err := c.contractTemplate.Execute(&buf, data)
	return buf.Bytes(), err
}

func (c TemplateRegistry) ProcessChainReaderConfig(data chainComponentsConfigData) ([]byte, error) {
	var buf bytes.Buffer
	err := c.chainReaderConfigTemplate.Execute(&buf, data)
	return buf.Bytes(), err
}

func getTemplate() (*TemplateRegistry, error) {
	once.Do(func() {
		contractTemplate, err := template.New("contractTemplate").Parse(ContractBindingTemplate)
		if err != nil {
			log.Fatalf("failed to parse contract template: %v", err)
		}
		chainReaderConfigTemplate, err := template.New("chainReaderConfigTemplate").Parse(ChainReaderConfigFactoryTemplate)
		if err != nil {
			log.Fatalf("failed to parse contract template: %v", err)
		}
		instance = &TemplateRegistry{
			contractTemplate:          *contractTemplate,
			chainReaderConfigTemplate: *chainReaderConfigTemplate,
		}

	})
	return instance, nil
}

var (
	instance *TemplateRegistry
	once     sync.Once
)

func GenerateContractBinding(packageName string, codeDetails CodeDetails) ([]byte, error) {
	data := contractTemplateData{
		PackageName: packageName,
		CodeDetails: codeDetails,
	}

	template := safeGetTemplate()

	return template.ProcessContract(data)
}

func safeGetTemplate() *TemplateRegistry {
	template, err := getTemplate()
	if err != nil {
		log.Fatalf("failed to parse contract template: %v", err)
	}
	return template
}

func GenerateChainReaderAndWriterFactory(packageName string, crConfig types.ChainReaderConfig, cwConfig types.ChainWriterConfig) ([]byte, error) {
	// Start processing from the root struct
	chainWriterConfigValueDefinition := GetValueDefinition(cwConfig)
	chainReaderConfigValueDefinition := GetValueDefinition(crConfig)

	bytes, err := safeGetTemplate().ProcessChainReaderConfig(chainComponentsConfigData{
		packageName,
		chainReaderConfigValueDefinition,
		chainWriterConfigValueDefinition,
	})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
