package ocr2key

import (
	"bytes"
	cryptorand "crypto/rand"
	"io"
	"math/big"

	"github.com/NethermindEth/juno/pkg/crypto/pedersen"
	starksig "github.com/NethermindEth/juno/pkg/crypto/signature"
	"github.com/NethermindEth/juno/pkg/crypto/weierstrass"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/starkkey"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var _ ocrtypes.OnchainKeyring = &starknetKeyring{}

var starkCurve = weierstrass.Stark()

type starknetKeyring struct {
	privateKey starksig.PrivateKey
}

func newStarknetKeyring(material io.Reader) (*starknetKeyring, error) {
	privKey, err := starksig.GenerateKey(starkCurve, material)
	if err != nil {
		return nil, err
	}
	return &starknetKeyring{privateKey: *privKey}, err
}

func (sk *starknetKeyring) PublicKey() ocrtypes.OnchainPublicKey {
	return starkkey.PubToStarkKey(sk.privateKey.PublicKey)
}

func (sk *starknetKeyring) reportToSigData(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	var dataArray []*big.Int
	rawReportContext := evmutil.RawReportContext(reportCtx)
	dataArray = append(dataArray, new(big.Int).SetBytes(rawReportContext[0][:]))
	dataArray = append(dataArray, new(big.Int).SetBytes(rawReportContext[1][:]))
	dataArray = append(dataArray, new(big.Int).SetBytes(rawReportContext[2][:]))

	// TODO: report hashing needs to be finalized and verified before merging

	// assert minimum length (timestamp + observers + observation number + observations + juelsPerFeeCoin)
	// 4 bytes (uint32) + 32 bytes + 1 byte + obs number * 16 bytes + 16 bytes
	timestampLen := 4
	observersLen := 32
	obsCountLen := 1
	obsLen := 16
	juelsLen := 16
	min := timestampLen + observersLen + obsCountLen + 0*obsLen + juelsLen
	if len(report) < min {
		return []byte{}, errors.Errorf("report did not meet minimum length, got %d want %d", len(report), min)
	}
	obsCount := int(report[timestampLen+observersLen])
	expectedLen := min + obsLen*obsCount
	if len(report) != expectedLen {
		return []byte{}, errors.Errorf("incorrect report length, got %d want %d", len(report), expectedLen)
	}
	dataArray = append(dataArray, new(big.Int).SetBytes(report[:timestampLen]))
	dataArray = append(dataArray, new(big.Int).SetBytes(report[timestampLen:timestampLen+observersLen]))
	dataArray = append(dataArray, new(big.Int).SetBytes(report[timestampLen+observersLen:timestampLen+observersLen+obsCountLen]))
	// hash observations
	for i := 0; i < obsCount; i++ {
		obsData := report[timestampLen+observersLen+obsCountLen+obsLen*i : timestampLen+observersLen+obsCountLen+obsLen*(i+1)]
		dataArray = append(dataArray, new(big.Int).SetBytes(obsData))
	}
	dataArray = append(dataArray, new(big.Int).SetBytes(report[expectedLen-juelsLen:]))

	hash := pedersen.ArrayDigest(dataArray...)
	return hash.Bytes(), nil
}

func (sk *starknetKeyring) Sign(reportCtx ocrtypes.ReportContext, report ocrtypes.Report) ([]byte, error) {
	hash, err := sk.reportToSigData(reportCtx, report)
	if err != nil {
		return []byte{}, err
	}

	r, s, err := starksig.Sign(cryptorand.Reader, &sk.privateKey, hash)
	if err != nil {
		return []byte{}, err
	}

	// construct signature using SEC encoding (instead of ASN.1 DER)
	// simpler to decode later on if needed
	// https://bitcoin.stackexchange.com/questions/92680/what-are-the-der-signature-and-sec-format
	buff := bytes.NewBuffer([]byte{0x04})
	if _, err := buff.Write(math.PaddedBigBytes(r, 32)); err != nil {
		return []byte{}, err
	}
	if _, err := buff.Write(math.PaddedBigBytes(s, 32)); err != nil {
		return []byte{}, err
	}

	out := buff.Bytes()
	if len(out) != sk.MaxSignatureLength() {
		return []byte{}, errors.Errorf("unexpected signature size, got %d want %d", len(out), sk.MaxSignatureLength())
	}
	return out, nil
}

func (sk *starknetKeyring) Verify(publicKey ocrtypes.OnchainPublicKey, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signature []byte) bool {
	var keys [2]starksig.PublicKey

	// convert OnchainPublicKey (starkkey) into ecdsa public keys (prepend 2 or 3 to indicate +/- Y coord)
	prefix := []byte{2, 3}
	for i := 0; i < len(prefix); i++ {
		keys[i] = starksig.PublicKey{Curve: starkCurve}

		// prepend sign byte
		compressedKey := append([]byte{prefix[i]}, publicKey...)
		keys[i].X, keys[i].Y = weierstrass.UnmarshalCompressed(starkCurve, compressedKey)

		// handle invalid publicKey
		if keys[i].X == nil || keys[i].Y == nil {
			return false
		}
	}

	// check valid signature length
	if len(signature) != sk.MaxSignatureLength() {
		return false
	}

	hash, err := sk.reportToSigData(reportCtx, report)
	if err != nil {
		return false
	}

	r := new(big.Int).SetBytes(signature[1:33])
	s := new(big.Int).SetBytes(signature[33:65])

	return starksig.Verify(&keys[0], hash, r, s) || starksig.Verify(&keys[1], hash, r, s)
}

func (sk *starknetKeyring) MaxSignatureLength() int {
	return 65
}

func (sk *starknetKeyring) marshal() ([]byte, error) {
	// https://github.com/ethereum/go-ethereum/blob/07508ac0e9695df347b9dd00d418c25151fbb213/crypto/crypto.go#L159
	return math.PaddedBigBytes(sk.privateKey.D, sk.privateKeyLen()), nil
}

func (sk *starknetKeyring) privateKeyLen() int {
	// https://github.com/NethermindEth/juno/blob/3e71279632d82689e5af03e26693ca5c58a2376e/pkg/crypto/weierstrass/weierstrass.go#L377
	N := starkCurve.Params().N
	bitSize := N.BitLen()
	return (bitSize + 7) / 8 // 32
}

func (sk *starknetKeyring) unmarshal(in []byte) error {
	// enforce byte length
	if len(in) != sk.privateKeyLen() {
		return errors.Errorf("unexpected seed size, got %d want %d", len(in), sk.privateKeyLen())
	}

	sk.privateKey.D = new(big.Int).SetBytes(in)
	sk.privateKey.PublicKey.Curve = starkCurve
	sk.privateKey.PublicKey.X, sk.privateKey.PublicKey.Y = starkCurve.ScalarBaseMult(in)
	return nil
}
