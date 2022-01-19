package keystore

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"

	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
)

type encryptedKeyRing struct {
	UpdatedAt     time.Time
	EncryptedKeys []byte
}

func (ekr encryptedKeyRing) Decrypt(password string) (keyRing, error) {
	if len(ekr.EncryptedKeys) == 0 {
		return newKeyRing(), nil
	}
	var cryptoJSON gethkeystore.CryptoJSON
	err := json.Unmarshal(ekr.EncryptedKeys, &cryptoJSON)
	if err != nil {
		return keyRing{}, err
	}
	marshalledRawKeyRingJson, err := gethkeystore.DecryptDataV3(cryptoJSON, adulteratedPassword(password))
	if err != nil {
		return keyRing{}, err
	}
	var rawKeys rawKeyRing
	err = json.Unmarshal(marshalledRawKeyRingJson, &rawKeys)
	if err != nil {
		return keyRing{}, err
	}
	ring, err := rawKeys.keys()
	if err != nil {
		return keyRing{}, err
	}
	return ring, nil
}

type keyStates struct {
	Eth map[string]*ethkey.State
}

func newKeyStates() keyStates {
	return keyStates{
		Eth: make(map[string]*ethkey.State),
	}
}

func (ks keyStates) validate(kr keyRing) (err error) {
	for id := range kr.Eth {
		_, exists := ks.Eth[id]
		if !exists {
			err = multierr.Combine(err, errors.Errorf("key %s is missing state", id))
		}
	}

	return err
}

type keyRing struct {
	CSA    map[string]csakey.KeyV2
	Eth    map[string]ethkey.KeyV2
	OCR    map[string]ocrkey.KeyV2
	OCR2   map[string]ocr2key.KeyBundle
	P2P    map[string]p2pkey.KeyV2
	Solana map[string]solkey.Key
	Terra  map[string]terrakey.Key
	VRF    map[string]vrfkey.KeyV2
}

func newKeyRing() keyRing {
	return keyRing{
		CSA:    make(map[string]csakey.KeyV2),
		Eth:    make(map[string]ethkey.KeyV2),
		OCR:    make(map[string]ocrkey.KeyV2),
		OCR2:   make(map[string]ocr2key.KeyBundle),
		P2P:    make(map[string]p2pkey.KeyV2),
		Solana: make(map[string]solkey.Key),
		Terra:  make(map[string]terrakey.Key),
		VRF:    make(map[string]vrfkey.KeyV2),
	}
}

func (kr *keyRing) Encrypt(password string, scryptParams utils.ScryptParams) (ekr encryptedKeyRing, err error) {
	marshalledRawKeyRingJson, err := json.Marshal(kr.raw())
	if err != nil {
		return ekr, err
	}
	cryptoJSON, err := gethkeystore.EncryptDataV3(
		marshalledRawKeyRingJson,
		[]byte(adulteratedPassword(password)),
		scryptParams.N,
		scryptParams.P,
	)
	if err != nil {
		return ekr, errors.Wrapf(err, "could not encrypt key ring")
	}
	encryptedKeys, err := json.Marshal(&cryptoJSON)
	if err != nil {
		return ekr, errors.Wrapf(err, "could not encode cryptoJSON")
	}
	return encryptedKeyRing{
		EncryptedKeys: encryptedKeys,
	}, nil
}

func (kr *keyRing) raw() (rawKeys rawKeyRing) {
	for _, csaKey := range kr.CSA {
		rawKeys.CSA = append(rawKeys.CSA, csaKey.Raw())
	}
	for _, ethKey := range kr.Eth {
		rawKeys.Eth = append(rawKeys.Eth, ethKey.Raw())
	}
	for _, ocrKey := range kr.OCR {
		rawKeys.OCR = append(rawKeys.OCR, ocrKey.Raw())
	}
	for _, ocr2key := range kr.OCR2 {
		rawKeys.OCR2 = append(rawKeys.OCR2, ocr2key.Raw())
	}
	for _, p2pKey := range kr.P2P {
		rawKeys.P2P = append(rawKeys.P2P, p2pKey.Raw())
	}
	for _, solkey := range kr.Solana {
		rawKeys.Solana = append(rawKeys.Solana, solkey.Raw())
	}
	for _, terrakey := range kr.Terra {
		rawKeys.Terra = append(rawKeys.Terra, terrakey.Raw())
	}
	for _, vrfKey := range kr.VRF {
		rawKeys.VRF = append(rawKeys.VRF, vrfKey.Raw())
	}
	return rawKeys
}

func (kr *keyRing) logPubKeys(lggr logger.Logger) {
	lggr = lggr.Named("KeyRing")
	var csaIDs []string
	for _, CSAKey := range kr.CSA {
		csaIDs = append(csaIDs, CSAKey.ID())
	}
	var ethIDs []string
	for _, ETHKey := range kr.Eth {
		ethIDs = append(ethIDs, ETHKey.ID())
	}
	var ocrIDs []string
	for _, OCRKey := range kr.OCR {
		ocrIDs = append(ocrIDs, OCRKey.ID())
	}
	var ocr2IDs []string
	for _, OCR2Key := range kr.OCR2 {
		ocr2IDs = append(ocr2IDs, OCR2Key.ID())
	}
	var p2pIDs []string
	for _, P2PKey := range kr.P2P {
		p2pIDs = append(p2pIDs, P2PKey.ID())
	}
	var solanaIDs []string
	for _, solanaKey := range kr.Solana {
		solanaIDs = append(solanaIDs, solanaKey.ID())
	}
	var terraIDs []string
	for _, terraKey := range kr.Terra {
		terraIDs = append(terraIDs, terraKey.ID())
	}
	var vrfIDs []string
	for _, VRFKey := range kr.VRF {
		vrfIDs = append(vrfIDs, VRFKey.ID())
	}
	if len(csaIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d CSA keys", len(csaIDs)), "keys", csaIDs)
	}
	if len(ethIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d ETH keys", len(ethIDs)), "keys", ethIDs)
	}
	if len(ocrIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d OCR keys", len(ocrIDs)), "keys", ocrIDs)
	}
	if len(ocr2IDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d OCR2 keys", len(ocr2IDs)), "keys", ocr2IDs)
	}
	if len(p2pIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d P2P keys", len(p2pIDs)), "keys", p2pIDs)
	}
	if len(solanaIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d Solana keys", len(solanaIDs)), "keys", solanaIDs)
	}
	if len(terraIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d Terra keys", len(terraIDs)), "keys", terraIDs)
	}
	if len(vrfIDs) > 0 {
		lggr.Infow(fmt.Sprintf("Unlocked %d VRF keys", len(vrfIDs)), "keys", vrfIDs)
	}
}

// rawKeyRing is an intermediate struct for encrypting / decrypting keyRing
// it holds only the essential key information to avoid adding unnecessary data
// (like public keys) to the database
type rawKeyRing struct {
	Eth    []ethkey.Raw
	CSA    []csakey.Raw
	OCR    []ocrkey.Raw
	OCR2   []ocr2key.Raw
	P2P    []p2pkey.Raw
	Solana []solkey.Raw
	Terra  []terrakey.Raw
	VRF    []vrfkey.Raw
}

func (rawKeys rawKeyRing) keys() (keyRing, error) {
	keyRing := newKeyRing()
	for _, rawCSAKey := range rawKeys.CSA {
		csaKey := rawCSAKey.Key()
		keyRing.CSA[csaKey.ID()] = csaKey
	}
	for _, rawETHKey := range rawKeys.Eth {
		ethKey := rawETHKey.Key()
		keyRing.Eth[ethKey.ID()] = ethKey
	}
	for _, rawOCRKey := range rawKeys.OCR {
		ocrKey := rawOCRKey.Key()
		keyRing.OCR[ocrKey.ID()] = ocrKey
	}
	for _, rawOCR2Key := range rawKeys.OCR2 {
		ocr2Key := rawOCR2Key.Key()
		keyRing.OCR2[ocr2Key.ID()] = ocr2Key
	}
	for _, rawP2PKey := range rawKeys.P2P {
		p2pKey := rawP2PKey.Key()
		keyRing.P2P[p2pKey.ID()] = p2pKey
	}
	for _, rawSolKey := range rawKeys.Solana {
		solKey := rawSolKey.Key()
		keyRing.Solana[solKey.ID()] = solKey
	}
	for _, rawTerraKey := range rawKeys.Terra {
		terraKey := rawTerraKey.Key()
		keyRing.Terra[terraKey.ID()] = terraKey
	}
	for _, rawVRFKey := range rawKeys.VRF {
		vrfKey := rawVRFKey.Key()
		keyRing.VRF[vrfKey.ID()] = vrfKey
	}
	return keyRing, nil
}

// adulteration prevents the password from getting used in the wrong place
func adulteratedPassword(password string) string {
	return "master-password-" + password
}
