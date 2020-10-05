package protocol

import (
	"bytes"
	"log"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink/libocr/gethwrappers/testoffchainaggregator"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/internal/signature"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type ContractReport struct {
	Ctx    signature.ReportingContext
	Values OracleValues
	Sig    []byte
}

type OracleValue struct {
	ID    types.OracleID
	Value observation.Observation
}

type OracleValues []OracleValue

type ContractReportWithSignatures struct {
	ContractReport
	Signatures [][]byte
}

func (o OracleValue) Equal(o2 OracleValue) bool {
	return (o.ID == o2.ID) && o.Value.Equal(o2.Value)
}

func (os OracleValues) Median() (observation.Observation, error) {
	if len(os) == 0 {
		return observation.Observation{}, errors.Errorf(
			"can't take median of empty list")
	}
	return os[len(os)/2].Value, nil
}

func (c *ContractReport) wireMessage() []byte {
	observations, err := c.observations()
	if err != nil {
		log.Println("could not serialize observations for transmission: ", err.Error())
		return nil
	}
	tag := c.Ctx.DomainSeparationTag()
	return append(tag[:], observations...)
}

func (c ContractReport) Equal(c2 ContractReport) bool {
	if (len(c.Values) != len(c2.Values)) || !c.Ctx.Equal(c2.Ctx) {
		return false
	}
	for i, v := range c.Values {
		ov := c2.Values[i]
		if !v.Equal(ov) {
			return false
		}
	}
	return true
}

func (c *ContractReport) observations() ([]byte, error) {
	serializedObservers, err := c.observers()
	if err != nil {
		return nil, err
	}
	var bigObservations []byte
	for _, v := range c.Values {
		sv := v.Value.Marshal()
		bigObservations = append(bigObservations, sv...)
	}
	return append(serializedObservers[:], bigObservations...), nil
}

func (c *ContractReport) observers() (rv common.Hash, err error) {
	if len(c.Values) > 32 {
		return rv, errors.Errorf("too many values! can only handle 32, got %d",
			len(c.Values))
	}
	for i, value := range c.Values {
		id := int(value.ID)
		if id < 0 || i > 31 {
			return [32]byte{}, errors.Errorf(
				"Oracle index %d for %#+v is out of range", id, value)
		}
		rv[i] = byte(uint8(id))
	}
	return rv, nil
}

func (c ContractReportWithSignatures) Equals(c2 ContractReportWithSignatures) bool {
	if (!c.ContractReport.Equal(c2.ContractReport)) ||
		(len(c.Signatures) != len(c2.Signatures)) {
		return false
	}
	for i := range c.Signatures {
		if !bytes.Equal(c.Signatures[i], c2.Signatures[i]) {
			return false
		}
	}
	return true
}

func (c *ContractReportWithSignatures) collateSignatures() (rs, ss [][32]byte, vs [32]byte) {
	for i, sig := range c.Signatures {
		rs = append(rs, common.BytesToHash(sig[:32]))
		ss = append(ss, common.BytesToHash(sig[32:64]))
		vs[i] = sig[64]
	}
	return rs, ss, vs
}

func (c *ContractReportWithSignatures) TransmissionArgs() (report []byte, rs,
	ss [][32]byte, vs [32]byte, err error) {
	report, err = c.ContractReport.OnChainReport()
	if err != nil {
		return nil, nil, nil, [32]byte{}, errors.Wrapf(err,
			"while constructing report for on-chain transmission")
	}
	rs, ss, vs = c.collateSignatures()
	return report, rs, ss, vs, nil
}

func getReportTypes() abi.Arguments {
	uABI, err := abi.JSON(strings.NewReader(testoffchainaggregator.TestOffchainAggregatorABI))
	if err != nil {
		panic(err)
	}
	return uABI.Methods["testDecodeReport"].Outputs
}

func (c *ContractReport) onChainObservations() (rv []*big.Int) {
	for _, v := range c.Values {
		rv = append(rv, v.Value.GoEthereumValue())
	}
	return rv
}

var reportTypes = getReportTypes()

func (c *ContractReport) OnChainReport() ([]byte, error) {
	observers, err := c.observers()
	if err != nil {
		return nil, errors.Wrapf(err, "while collating observers for onChainReport")
	}
	return reportTypes.Pack(c.Ctx.DomainSeparationTag(), observers, c.onChainObservations())
}

func (c *ContractReport) Sign(signer func([]byte) ([]byte, error)) error {
	report, err := c.OnChainReport()
	if err != nil {
		return err
	}
	c.Sig, err = signer(report)
	if err != nil {
		return errors.Wrapf(err, "while signing on-chain report")
	}
	return nil
}

func (c *ContractReport) verify(a types.OnChainSigningAddress) (err error) {
	report, err := c.OnChainReport()
	if err != nil {
		return err
	}
	var dummyID types.OracleID
	address := map[types.OnChainSigningAddress]types.OracleID{a: dummyID}
	_, err = signature.VerifyOnChain(report, c.Sig, address)
	return err
}

func (c *ContractReportWithSignatures) VerifySignatures(
	as signature.EthAddresses,
) error {
	report, err := c.OnChainReport()
	if err != nil {
		return errors.Wrapf(err,
			"while serializing report to check signatures on it")
	}
	seen := make(map[types.OracleID]bool)
	for _, sig := range c.Signatures {
		if oid, err := signature.VerifyOnChain(report, sig, as); err != nil {
			return errors.Wrapf(err,
				"while checking a signature on a report, 0x%x", sig)
		} else {
			if seen[oid] {
				return errors.Errorf("oracle #%d signed more than once", oid)
			}
			seen[oid] = true
		}
	}
	return nil
}
