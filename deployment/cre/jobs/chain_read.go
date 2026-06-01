package jobs

import (
	"errors"
	"fmt"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
)

type ChainReadCapabilityJobSpecInput struct {
	OCRContractQualifier string   `json:"ocrContractQualifier" yaml:"ocrContractQualifier"`
	OCRChainSelector     uint64   `json:"ocrChainSelector" yaml:"ocrChainSelector"`
	BootstrapperOCR3Urls []string `json:"bootstrapperOCR3Urls" yaml:"bootstrapperOCR3Urls"`
}

func VerifyChainReadCapabilityPreconditions(e cldf.Environment, input ChainReadCapabilityJobSpecInput) error {
	if input.OCRChainSelector == 0 {
		return errors.New("ocr chain selector is required")
	}
	if len(input.BootstrapperOCR3Urls) == 0 {
		return errors.New("at least one bootstrapper OCR3 URL is required")
	}
	for i, u := range input.BootstrapperOCR3Urls {
		if u == "" {
			return fmt.Errorf("bootstrapper OCR3 URL at index %d is empty", i)
		}
	}
	if input.OCRContractQualifier == "" {
		return errors.New("ocr contract qualifier is required")
	}

	ocrAddrRefKey := pkg.GetOCR3CapabilityAddressRefKey(input.OCRChainSelector, input.OCRContractQualifier)
	if _, err := e.DataStore.Addresses().Get(ocrAddrRefKey); err != nil {
		return fmt.Errorf("failed to get OCR contract address for ref key %s: %w", ocrAddrRefKey, err)
	}

	return nil
}

type ChainReadCapabilityConfig struct {
	ObservationPollerWorkersCount uint          `json:"observationPollerWorkersCount,omitempty" yaml:"observationPollerWorkersCount,omitempty"`
	ObservationPollPeriod         time.Duration `json:"observationPollPeriod,omitempty" yaml:"observationPollPeriod,omitempty"`
	ChainHeightPollPeriod         time.Duration `json:"chainHeightPollPeriod,omitempty" yaml:"chainHeightPollPeriod,omitempty"`
	UnknownRequestsTTL            time.Duration `json:"unknownRequestsTTL,omitempty" yaml:"unknownRequestsTTL,omitempty"`
}
