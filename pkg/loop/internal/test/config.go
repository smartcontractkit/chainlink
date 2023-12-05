package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type staticConfigProvider struct{}

// TODO validate start/Close calls?
func (s staticConfigProvider) Start(ctx context.Context) error { return nil }

func (s staticConfigProvider) Close() error { return nil }

func (s staticConfigProvider) Ready() error { panic("unimplemented") }

func (s staticConfigProvider) Name() string { panic("unimplemented") }

func (s staticConfigProvider) HealthReport() map[string]error { panic("unimplemented") }

func (s staticConfigProvider) OffchainConfigDigester() libocr.OffchainConfigDigester {
	return staticOffchainConfigDigester{}
}

func (s staticConfigProvider) ContractConfigTracker() libocr.ContractConfigTracker {
	return staticContractConfigTracker{}
}

type staticOffchainConfigDigester struct{}

func (s staticOffchainConfigDigester) ConfigDigest(config libocr.ContractConfig) (libocr.ConfigDigest, error) {
	if !assert.ObjectsAreEqual(contractConfig, config) {
		return libocr.ConfigDigest{}, fmt.Errorf("expected contract config %v but got %v", configDigest, config)
	}
	return configDigest, nil
}

func (s staticOffchainConfigDigester) ConfigDigestPrefix() (libocr.ConfigDigestPrefix, error) {
	return configDigestPrefix, nil
}

type staticContractConfigTracker struct{}

func (s staticContractConfigTracker) Notify() <-chan struct{} { return nil }

func (s staticContractConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, libocr.ConfigDigest, error) {
	return changedInBlock, configDigest, nil
}

func (s staticContractConfigTracker) LatestConfig(ctx context.Context, cib uint64) (libocr.ContractConfig, error) {
	if changedInBlock != cib {
		return libocr.ContractConfig{}, fmt.Errorf("expected changed in block %d but got %d", changedInBlock, cib)
	}
	return contractConfig, nil
}

func (s staticContractConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	return blockHeight, nil
}

type staticContractTransmitter struct{}

func (s staticContractTransmitter) Transmit(ctx context.Context, rc libocr.ReportContext, r libocr.Report, ss []libocr.AttributedOnchainSignature) error {
	if !assert.ObjectsAreEqual(reportContext, rc) {
		return fmt.Errorf("expected report context %v but got %v", reportContext, report)
	}
	if !bytes.Equal(report, r) {
		return fmt.Errorf("expected report %x but got %x", report, r)
	}
	if !assert.ObjectsAreEqual(sigs, ss) {
		return fmt.Errorf("expected signatures %v but got %v", sigs, ss)
	}
	return nil
}

func (s staticContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (libocr.ConfigDigest, uint32, error) {
	return configDigest, epoch, nil
}

func (s staticContractTransmitter) FromAccount() (libocr.Account, error) {
	return account, nil
}

type staticChainReader struct{}

func (c staticChainReader) GetMaxEncodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return 0, errors.New("not used for these test")
}

func (c staticChainReader) GetMaxDecodingSize(ctx context.Context, n int, itemType string) (int, error) {
	return 0, errors.New("not used for these test")
}

func (c staticChainReader) Encode(ctx context.Context, item any, itemType string) ([]byte, error) {
	return nil, errors.New("not used for these test")
}

func (c staticChainReader) Decode(ctx context.Context, raw []byte, into any, itemType string) error {
	return errors.New("not used for these test")
}

func (c staticChainReader) GetLatestValue(ctx context.Context, bc types.BoundContract, method string, params, returnVal any) error {
	if !assert.ObjectsAreEqual(bc, boundContract) {
		return fmt.Errorf("expected report context %v but got %v", boundContract, bc)
	}
	if method != medianContractGenericMethod {
		return fmt.Errorf("expected generic contract method %v but got %v", medianContractGenericMethod, method)
	}
	gotParams, ok := params.(*map[string]any)
	if !ok {
		return fmt.Errorf("Invalid parameter type received in GetLatestValue. Expected %T but received %T", gotParams, params)
	}
	if (*gotParams)["param1"] != getLatestValueParams["param1"] || (*gotParams)["param2"] != getLatestValueParams["param2"] {
		return fmt.Errorf("Wrong params value received in GetLatestValue. Expected %v but received %v", getLatestValueParams, *gotParams)
	}

	ret, ok := returnVal.(*map[string]any)
	if !ok {
		return fmt.Errorf("Wrong type passed for retVal param to GetLatestValue impl (expected %T instead of %T", ret, returnVal)
	}

	(*ret)["ret1"] = latestValue["ret1"]
	(*ret)["ret2"] = latestValue["ret2"]

	return nil
}
