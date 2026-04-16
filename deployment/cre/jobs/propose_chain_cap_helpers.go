package jobs

import (
	"errors"
	"fmt"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

type commonCapFields struct {
	Environment          string
	Domain               string
	Zone                 string
	DONName              string
	ChainSelector        uint64
	OCRChainSelector     uint64
	BootstrapperOCR3Urls []string
	OCRContractQualifier string
	DeltaStage           time.Duration
}

func validateCommonFields(f commonCapFields) error {
	if f.Environment == "" {
		return errors.New("environment is required")
	}
	if f.Domain == "" {
		return errors.New("domain is required")
	}
	if f.Zone == "" {
		return errors.New("zone is required")
	}
	if f.DONName == "" {
		return errors.New("donName is required")
	}
	if f.ChainSelector == 0 {
		return errors.New("chain selector is required")
	}
	if f.OCRChainSelector == 0 {
		return errors.New("ocr chain selector is required")
	}
	if len(f.BootstrapperOCR3Urls) == 0 {
		return errors.New("at least one bootstrapper OCR3 URL is required")
	}
	for i, u := range f.BootstrapperOCR3Urls {
		if u == "" {
			return fmt.Errorf("bootstrapper OCR3 URL at index %d is empty", i)
		}
	}
	if f.OCRContractQualifier == "" {
		return errors.New("ocr contract qualifier is required")
	}
	if f.DeltaStage <= 0 {
		return fmt.Errorf("deltaStage (%s) must be greater than 0", f.DeltaStage)
	}
	return nil
}

type resolvedAddresses struct {
	ForwarderAddress string
}

func resolveContractAddresses(
	e cldf.Environment,
	ocrChainSelector uint64,
	ocrQualifier string,
	fwdChainSelector uint64,
	fwdQualifier string,
) (resolvedAddresses, error) {
	ocrAddrRefKey := pkg.GetOCR3CapabilityAddressRefKey(ocrChainSelector, ocrQualifier)
	if _, err := e.DataStore.Addresses().Get(ocrAddrRefKey); err != nil {
		return resolvedAddresses{}, fmt.Errorf("failed to get OCR contract address for ref key %s: %w", ocrAddrRefKey, err)
	}

	fwdAddrRefKey := pkg.GetKeystoneForwarderCapabilityAddressRefKey(fwdChainSelector, fwdQualifier)
	fwdAddress, err := e.DataStore.Addresses().Get(fwdAddrRefKey)
	if err != nil {
		return resolvedAddresses{}, fmt.Errorf("failed to get CRE forwarder address for ref key %q: %w", fwdAddrRefKey, err)
	}

	return resolvedAddresses{ForwarderAddress: fwdAddress.Address}, nil
}

func validateOverrideNetwork(got, expected, nodeID string) error {
	if got != "" && got != expected {
		return fmt.Errorf("network in override config must be %q if set; got %q for node %s", expected, got, nodeID)
	}
	return nil
}

func validateOverrideForwarder(got, expected, nodeID string) error {
	if got != "" && got != expected {
		return fmt.Errorf(
			"CRE forwarder address in override config (%s) does not match address from data store (%s) for node %s; "+
				"this field is auto-populated and can be omitted",
			got, expected, nodeID,
		)
	}
	return nil
}

func proposeAndReport(
	e cldf.Environment,
	job pkg.StandardCapabilityJob,
	nodeIDToConfig map[string]string,
	domain, donName, zone string,
) (cldf.ChangesetOutput, error) {
	report, err := operations.ExecuteSequence(
		e.OperationsBundle,
		operations2.ProposeStandardCapabilityJob,
		operations2.ProposeStandardCapabilityJobDeps{Env: e},
		operations2.ProposeStandardCapabilityJobInput{
			Job:            job,
			NodeIDToConfig: nodeIDToConfig,
			Domain:         domain,
			DONName:        donName,

			DONFilters: []offchain.TargetDONFilter{
				{Key: "zone", Value: zone},
			},
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to propose standard capability job: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
