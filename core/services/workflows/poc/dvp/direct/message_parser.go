package direct

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/yudaprama/iso20022/model"
	"github.com/yudaprama/iso20022/pacs"
	"github.com/yudaprama/iso20022/sese"
)

// For simplicity, I'm only parsing the newest version of the messages, this could need to be changed?

type ISO20022Document interface {
	// VerifyRequiredAttestations TODO determine what type of signature would sign the Merkel Tree
	VerifyRequiredAttestations(key *ecdsa.PublicKey) error
	private()
}

func ParseMessage(rawXml []byte) (ISO20022Document, error) {
	bd := &basicSwiftDocument{}
	if err := xml.Unmarshal(rawXml, bd); err != nil {
		return nil, err
	}

	parser, ok := docIdToParser[bd.XMLName.Space]
	if !ok {
		return nil, fmt.Errorf("unsupported message type: %s", bd.XMLName.Space)
	}

	return parser(rawXml)
}

type basicSwiftDocument struct {
	XMLName xml.Name
}

// I'm going to parse the latest of each example because it's annoying to do all for now
// We can use xsdgen to generate the files, the repos I've found that do it don't have all the message we need.
// The repo I'm pulling in seems manual.  Regardless, we should generate getters
// for fields so we can make an interface that abstracts the version.
// Then the types would embed one of those instead of the specific version.

type SecuritiesSettlementTransactionInstruction struct {
	sese.Document02300107
}

func (s *SecuritiesSettlementTransactionInstruction) VerifyRequiredAttestations() error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (s *SecuritiesSettlementTransactionInstruction) private() {}

type SecuritiesSettlementTransactionStatusAdvice struct {
	sese.Document02400208
}

func (s SecuritiesSettlementTransactionStatusAdvice) VerifyRequiredAttestations() error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (s SecuritiesSettlementTransactionStatusAdvice) private() {}

type FIToFIPaymentStatusReport struct {
	pacs.Document00200108
}

func (F FIToFIPaymentStatusReport) VerifyRequiredAttestations() error {
	//TODO implement me
	panic("implement me")
	return nil
}

func (F FIToFIPaymentStatusReport) private() {}

func verifyAttestations(key *ecdsa.PublicKey, supDatas []*model.SupplementaryData1, fields map[string]any) error {
	var merkelTree *model.SupplementaryData1
	for _, supData := range supDatas {
		if supData.PlaceAndName != nil && *supData.PlaceAndName == "MT" {
			merkelTree = supData
		}
	}

	merkelBytes, err := parseMerkel(merkelTree)
	if err != nil {
		return fmt.Errorf("error getting Merkel Tree bytes: %w", err)
	}


	// TODO verify signature against bytes of Merkel Tree
}

func parseMerkel(data *model.SupplementaryData1) (*merkletree.Tree, error) {
	if data == nil {
		return nil, errors.New("data not found")
	}

	hData := data.Envelope
	if hData == nil {
		return nil, fmt.Errorf("missing signature on Merkel Tree")
	}

	parts := strings.Split(string(*hData), "\n")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Merkel Tree signature format")
	}

	mt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("error decoding Merkel Tree: %w", err)
	}

	sig, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding signature: %w", err)
	}

	pkRaw, err := x509.ParsePKIXPublicKey(sig)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %w", err)
	}

	pk, ok := pkRaw.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}

	if !ecdsa.VerifyASN1(pk, mt, sig) {
		return nil, errors.New("signature on Merkel Tree is invalid")
	}

	merkleTree.
}

// TODO this would need to be verified, I copied one example from ChatGPT and just ran with it
var docIdToParser = map[string]func([]byte) (ISO20022Document, error){
	"urn:iso:std:iso:20022:tech:xsd:sese.023.001.07": func(rawXml []byte) (ISO20022Document, error) {
		msg := &SecuritiesSettlementTransactionInstruction{}
		err := xml.Unmarshal(rawXml, &msg.Document02300107)
		return msg, err
	},
	"urn:iso:std:iso:20022:tech:xsd:sese.024.002.08": func(rawXml []byte) (ISO20022Document, error) {
		msg := &SecuritiesSettlementTransactionStatusAdvice{}
		err := xml.Unmarshal(rawXml, &msg.Document02400208)
		return msg, err
	},
	"urn:iso:std:iso:20022:tech:xsd:pacs.002.001.08": func(bytes []byte) (ISO20022Document, error) {
		msg := &FIToFIPaymentStatusReport{}
		err := xml.Unmarshal(bytes, &msg.Document00200108)
		return msg, err
	},
}
