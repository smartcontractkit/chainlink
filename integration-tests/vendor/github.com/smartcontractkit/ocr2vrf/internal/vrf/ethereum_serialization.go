package vrf

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/gethwrappers/vrfbeacon"
	"github.com/smartcontractkit/ocr2vrf/internal/util"
	vrf_types "github.com/smartcontractkit/ocr2vrf/types"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
)

type EthereumReportSerializer struct {
	G kyber.Group
}

func (e *EthereumReportSerializer) DeserializeReport(
	rawReport []byte,
) (vrf_types.AbstractReport, error) {
	arguments := vrfABI().Methods["exposeType"].Inputs
	parsedReport, err := arguments.Unpack(rawReport)
	if err != nil {
		return vrf_types.AbstractReport{}, err
	}
	type reportType = vrfbeacon.VRFBeaconReportReport
	report := reportType{}
	abi.ConvertType(parsedReport[0], &report)
	abstractReport, err := e.ConvertToAbstractReport(report)
	if err != nil {
		return vrf_types.AbstractReport{}, err
	}
	return abstractReport, err
}

func (e *EthereumReportSerializer) MaxReportLength() uint {

	panic("implement me")
}

func (e *EthereumReportSerializer) ReportLength(abstractReport vrf_types.AbstractReport) uint {

	panic("implement me")
}

var _ vrf_types.ReportSerializer = &EthereumReportSerializer{}

func (e *EthereumReportSerializer) SerializeReport(
	report vrf_types.AbstractReport,
) ([]byte, error) {
	beaconReport, err := e.ConvertToBeaconReport(report)
	if err != nil {
		return nil, err
	}
	arguments := vrfABI().Methods["exposeType"].Inputs
	serializedReport, err := arguments.Pack(beaconReport)
	if err != nil {
		return nil, err
	}
	return serializedReport, err
}

func (e *EthereumReportSerializer) ConvertToBeaconReport(
	report vrf_types.AbstractReport,
) (vrfbeacon.VRFBeaconReportReport, error) {
	vrfOutputs := make(
		[]vrfbeacon.VRFBeaconTypesVRFOutput, 0, len(report.Outputs),
	)
	emptyReport := vrfbeacon.VRFBeaconReportReport{}
	for _, output := range report.Outputs {
		p := e.G.Point()
		if err := p.UnmarshalBinary(output.VRFProof[:]); err != nil {
			return emptyReport, errors.Wrap(err, "while unmarshalling vrf proof")
		}
		x, y := big.NewInt(0), big.NewInt(0)
		if !p.Equal(e.G.Point().Null()) {
			x, y = affineCoordinates(p)
		}
		vrfProof := vrfbeacon.ECCArithmeticG1Point{P: [2]*big.Int{x, y}}
		type callbackType = vrfbeacon.VRFBeaconTypesCostedCallback
		callbacks := make([]callbackType, 0, len(output.Callbacks))
		for _, callback := range output.Callbacks {
			beaconCostedCallback := callbackType{
				Callback: vrfbeacon.VRFBeaconTypesCallback{
					RequestID:      big.NewInt(0).SetUint64(callback.RequestID),
					NumWords:       callback.NumWords,
					Requester:      callback.Requester,
					Arguments:      callback.Arguments,
					GasAllowance:   callback.GasAllowance,
					SubID:          callback.SubscriptionID,
					GasPrice:       callback.GasPrice,
					WeiPerUnitLink: callback.WeiPerUnitLink,
				},
				Price: callback.Price,
			}
			callbacks = append(callbacks, beaconCostedCallback)
		}
		vrfOutput := vrfbeacon.VRFBeaconTypesVRFOutput{
			BlockHeight:       output.BlockHeight,
			ConfirmationDelay: big.NewInt(int64(output.ConfirmationDelay)),
			VrfOutput:         vrfProof,
			Callbacks:         callbacks,
		}
		vrfOutputs = append(vrfOutputs, vrfOutput)
	}
	onchainReport := vrfbeacon.VRFBeaconReportReport{
		Outputs:            vrfOutputs,
		JuelsPerFeeCoin:    report.JuelsPerFeeCoin,
		ReasonableGasPrice: report.ReasonableGasPrice,
		RecentBlockHeight:  report.RecentBlockHeight,
		RecentBlockHash:    report.RecentBlockHash,
	}
	return onchainReport, nil
}

func (e *EthereumReportSerializer) ConvertToAbstractReport(
	report vrfbeacon.VRFBeaconReportReport,
) (vrf_types.AbstractReport, error) {
	abstractOutputs := make([]vrf_types.AbstractVRFOutput, 0, len(report.Outputs))
	for _, out := range report.Outputs {
		xCoordinate := mod.NewInt(out.VrfOutput.P[0], bn256.P)
		yCoordinate := mod.NewInt(out.VrfOutput.P[1], bn256.P)
		vrfG1Point, err := altbn_128.CoordinatesToG1(xCoordinate, yCoordinate)
		if err != nil {
			return vrf_types.AbstractReport{},
				util.WrapErrorf(
					err,
					"could not parse VRF proof (0x%x, 0x%x)",
					xCoordinate,
					yCoordinate,
				)
		}
		vrfG1PointBinary, err := vrfG1Point.MarshalBinary()
		var vrfProof [32]byte
		copy(vrfProof[:], vrfG1PointBinary[:])

		if err != nil {
			errMsg := "while unmarshalling vrf proof"
			return vrf_types.AbstractReport{}, util.WrapError(err, errMsg)
		}
		var abstractCallbacks []vrf_types.AbstractCostedCallbackRequest
		for _, c := range out.Callbacks {
			abstractCallback := vrf_types.AbstractCostedCallbackRequest{
				out.BlockHeight,
				uint32(out.ConfirmationDelay.Uint64()),
				c.Callback.SubID,
				c.Price,
				c.Callback.RequestID.Uint64(),
				c.Callback.NumWords,
				c.Callback.Requester,
				c.Callback.Arguments,
				c.Callback.GasAllowance,
				c.Callback.GasPrice,
				c.Callback.WeiPerUnitLink,
			}
			abstractCallbacks = append(abstractCallbacks, abstractCallback)
		}
		abstractOutput := vrf_types.AbstractVRFOutput{
			out.BlockHeight,
			uint32(out.ConfirmationDelay.Uint64()),
			vrfProof,
			abstractCallbacks,
		}
		abstractOutputs = append(abstractOutputs, abstractOutput)
	}
	abstractReport := vrf_types.AbstractReport{
		abstractOutputs,
		report.JuelsPerFeeCoin,
		report.ReasonableGasPrice,
		report.RecentBlockHeight,
		report.RecentBlockHash,
	}
	return abstractReport, nil
}

func vrfABI() *abi.ABI {
	rv, err := abi.JSON(
		strings.NewReader(vrfbeacon.VRFBeaconReportMetaData.ABI),
	)
	if err != nil {
		panic(err)
	}
	return &rv
}
