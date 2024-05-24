package starkkey

import (
	"bytes"
	"io"
	"math/big"

	"github.com/pkg/errors"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

var _ types.OnchainKeyring = &OCR2Key{}

type OCR2Key struct {
	Key
}

func NewOCR2Key(material io.Reader) (*OCR2Key, error) {
	k, err := GenerateKey(material)

	return &OCR2Key{k}, err
}

func (sk *OCR2Key) PublicKey() types.OnchainPublicKey {
	ans := new(felt.Felt).SetBytes(sk.pub.X.Bytes()).Bytes()
	return ans[:]
}

func ReportToSigData(reportCtx types.ReportContext, report types.Report) (*big.Int, error) {
	var dataArray []*big.Int

	rawReportContext := rawReportContext(reportCtx)
	dataArray = append(dataArray, new(big.Int).SetBytes(rawReportContext[0][:]))
	dataArray = append(dataArray, new(big.Int).SetBytes(rawReportContext[1][:]))
	dataArray = append(dataArray, new(big.Int).SetBytes(rawReportContext[2][:]))

	// split report into separate felts for hashing
	splitReport, err := splitReport(report)
	if err != nil {
		return &big.Int{}, err
	}
	for i := 0; i < len(splitReport); i++ {
		dataArray = append(dataArray, new(big.Int).SetBytes(splitReport[i]))
	}

	hash, err := curve.Curve.ComputeHashOnElements(dataArray)
	if err != nil {
		return &big.Int{}, err
	}
	return hash, nil
}

func (sk *OCR2Key) Sign(reportCtx types.ReportContext, report types.Report) ([]byte, error) {
	hash, err := ReportToSigData(reportCtx, report)
	if err != nil {
		return []byte{}, err
	}
	r, s, err := curve.Curve.Sign(hash, sk.priv)
	if err != nil {
		return []byte{}, err
	}

	// enforce s <= N/2 to prevent signature malleability
	if s.Cmp(new(big.Int).Rsh(curve.Curve.N, 1)) > 0 {
		s.Sub(curve.Curve.N, s)
	}

	// encoding: public key (32 bytes) + r (32 bytes) + s (32 bytes)
	buff := bytes.NewBuffer([]byte(sk.PublicKey()))
	if _, err := buff.Write(padBytes(r.Bytes(), byteLen)); err != nil {
		return []byte{}, err
	}
	if _, err := buff.Write(padBytes(s.Bytes(), byteLen)); err != nil {
		return []byte{}, err
	}

	out := buff.Bytes()
	if len(out) != sk.MaxSignatureLength() {
		return []byte{}, errors.Errorf("unexpected signature size, got %d want %d", len(out), sk.MaxSignatureLength())
	}
	return out, nil
}

func (sk *OCR2Key) Sign3(digest types.ConfigDigest, seqNr uint64, r types.Report) (signature []byte, err error) {
	return nil, errors.New("not implemented")
}

func (sk *OCR2Key) Verify(publicKey types.OnchainPublicKey, reportCtx types.ReportContext, report types.Report, signature []byte) bool {
	// check valid signature length
	if len(signature) != sk.MaxSignatureLength() {
		return false
	}

	// convert OnchainPublicKey (starkkey) into ecdsa public keys (prepend 2 or 3 to indicate +/- Y coord)
	var keys [2]PublicKey
	keys[0].X = new(big.Int).SetBytes(publicKey)
	keys[0].Y = curve.Curve.GetYCoordinate(keys[0].X)

	// When there is no point with the provided x-coordinate, the GetYCoordinate function returns the nil value.
	if keys[0].Y == nil {
		return false
	}

	keys[1].X = keys[0].X
	keys[1].Y = new(big.Int).Mul(keys[0].Y, big.NewInt(-1))

	hash, err := ReportToSigData(reportCtx, report)
	if err != nil {
		return false
	}

	r := new(big.Int).SetBytes(signature[32:64])
	s := new(big.Int).SetBytes(signature[64:])

	// Only allow canonical signatures to avoid signature malleability. Verify s <= N/2
	if s.Cmp(new(big.Int).Rsh(curve.Curve.N, 1)) == 1 {
		return false
	}

	return curve.Curve.Verify(hash, r, s, keys[0].X, keys[0].Y) || curve.Curve.Verify(hash, r, s, keys[1].X, keys[1].Y)
}

func (sk *OCR2Key) Verify3(publicKey types.OnchainPublicKey, cd types.ConfigDigest, seqNr uint64, r types.Report, signature []byte) bool {
	return false
}

func (sk *OCR2Key) MaxSignatureLength() int {
	return 32 + 32 + 32 // publickey + r + s
}

func (sk *OCR2Key) Marshal() ([]byte, error) {
	return padBytes(sk.priv.Bytes(), sk.privateKeyLen()), nil
}

func (sk *OCR2Key) privateKeyLen() int {
	// https://github.com/NethermindEth/juno/blob/3e71279632d82689e5af03e26693ca5c58a2376e/pkg/crypto/weierstrass/weierstrass.go#L377
	return 32
}

func (sk *OCR2Key) Unmarshal(in []byte) error {
	// enforce byte length
	if len(in) != sk.privateKeyLen() {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), sk.privateKeyLen())
	}

	sk.Key = Raw(in).Key()
	return nil
}

func splitReport(report types.Report) ([][]byte, error) {
	chunkSize := 32
	if len(report)%chunkSize != 0 {
		return [][]byte{}, errors.New("invalid report length")
	}

	// order is guaranteed by buildReport:
	//   observation_timestamp
	//   observers
	//   observations_len
	//   observations
	//   juels_per_fee_coin
	//   gas_price
	slices := [][]byte{}
	for i := 0; i < len(report)/chunkSize; i++ {
		idx := i * chunkSize
		slices = append(slices, report[idx:(idx+chunkSize)])
	}

	return slices, nil
}

// NOTE: this should sit in the ocr2 package but that causes import cycles
func rawReportContext(repctx types.ReportContext) [3][32]byte {
	rawReportContext := evmutil.RawReportContext(repctx)
	// NOTE: Ensure extra_hash is 31 bytes with first byte blanked out
	// libocr generates a 32 byte extraHash but we need to fit it into a felt
	rawReportContext[2][0] = 0
	return rawReportContext
}
