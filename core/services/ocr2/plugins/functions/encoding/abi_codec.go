package encoding

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

type ReportCodec interface {
	EncodeReport(requests []*ProcessedRequest) ([]byte, error)
	DecodeReport(raw []byte) ([]*ProcessedRequest, error)
}

type reportCodecV1 struct {
	reportTypes abi.Arguments
}

func NewReportCodec(contractVersion uint32) (ReportCodec, error) {
	switch contractVersion {
	case 1:
		reportTypes, err := getReportTypesV1()
		if err != nil {
			return nil, err
		}
		return &reportCodecV1{reportTypes: reportTypes}, nil
	default:
		return nil, fmt.Errorf("unknown contract version: %d", contractVersion)
	}
}

func SliceToByte32(slice []byte) ([32]byte, error) {
	if len(slice) != 32 {
		return [32]byte{}, fmt.Errorf("input length is not 32 bytes: %d", len(slice))
	}
	var res [32]byte
	copy(res[:], slice[:32])
	return res, nil
}

func getReportTypesV1() (abi.Arguments, error) {
	bytes32ArrType, err := abi.NewType("bytes32[]", "", []abi.ArgumentMarshaling{})
	if err != nil {
		return nil, fmt.Errorf("unable to create an ABI type object for bytes32[]")
	}
	bytesArrType, err := abi.NewType("bytes[]", "", []abi.ArgumentMarshaling{})
	if err != nil {
		return nil, fmt.Errorf("unable to create an ABI type object for bytes[]")
	}
	return abi.Arguments([]abi.Argument{
		{Name: "ids", Type: bytes32ArrType},
		{Name: "results", Type: bytesArrType},
		{Name: "errors", Type: bytesArrType},
		{Name: "onchain_metadata", Type: bytesArrType},
		{Name: "processing_metadata", Type: bytesArrType},
	}), nil
}

func (c *reportCodecV1) EncodeReport(requests []*ProcessedRequest) ([]byte, error) {
	size := len(requests)
	if size == 0 {
		return []byte{}, nil
	}
	ids := make([][32]byte, size)
	results := make([][]byte, size)
	errors := make([][]byte, size)
	onchainMetadata := make([][]byte, size)
	processingMetadata := make([][]byte, size)
	for i := 0; i < size; i++ {
		var err error
		ids[i], err = SliceToByte32(requests[i].RequestID)
		if err != nil {
			return nil, err
		}
		results[i] = requests[i].Result
		errors[i] = requests[i].Error
		onchainMetadata[i] = requests[i].OnchainMetadata
		processingMetadata[i] = requests[i].CoordinatorContract
		// CallbackGasLimit is not ABI-encoded
	}
	return c.reportTypes.Pack(ids, results, errors, onchainMetadata, processingMetadata)
}

func (c *reportCodecV1) DecodeReport(raw []byte) ([]*ProcessedRequest, error) {
	reportElems := map[string]interface{}{}
	if err := c.reportTypes.UnpackIntoMap(reportElems, raw); err != nil {
		return nil, errors.WithMessage(err, "unable to unpack elements from raw report")
	}

	idsIface, idsOK := reportElems["ids"]
	resultsIface, resultsOK := reportElems["results"]
	errorsIface, errorsOK := reportElems["errors"]
	onchainMetaIface, onchainMetaOK := reportElems["onchain_metadata"]
	processingMetaIface, processingMetaOK := reportElems["processing_metadata"]
	if !idsOK || !resultsOK || !errorsOK || !onchainMetaOK || !processingMetaOK {
		return nil, fmt.Errorf("missing arrays in raw report, ids: %v, results: %v, errors: %v", idsOK, resultsOK, errorsOK)
	}

	ids, idsOK := idsIface.([][32]byte)
	results, resultsOK := resultsIface.([][]byte)
	errors, errorsOK := errorsIface.([][]byte)
	onchainMeta, onchainMetaOK := onchainMetaIface.([][]byte)
	processingMeta, processingMetaOK := processingMetaIface.([][]byte)
	if !idsOK || !resultsOK || !errorsOK || !onchainMetaOK || !processingMetaOK {
		return nil, fmt.Errorf("unable to cast part of raw report into array type, ids: %v, results: %v, errors: %v", idsOK, resultsOK, errorsOK)
	}

	size := len(ids)
	if len(results) != size || len(errors) != size || len(onchainMeta) != size || len(processingMeta) != size {
		return nil, fmt.Errorf("unequal sizes of arrays parsed from raw report, ids: %v, results: %v, errors: %v", len(ids), len(results), len(errors))
	}
	if size == 0 {
		return []*ProcessedRequest{}, nil
	}

	decoded := make([]*ProcessedRequest, size)
	for i := 0; i < size; i++ {
		decoded[i] = &ProcessedRequest{
			RequestID:           ids[i][:],
			Result:              results[i],
			Error:               errors[i],
			OnchainMetadata:     onchainMeta[i],
			CoordinatorContract: processingMeta[i],
			// CallbackGasLimit is not ABI-encoded
		}
	}
	return decoded, nil
}
