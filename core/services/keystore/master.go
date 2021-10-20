package keystore

import (
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

func New(db *gorm.DB, scryptParams utils.ScryptParams) *Master {
	return &Master{
		eth: newEthKeyStore(db, scryptParams),
		csa: newCSAKeyStore(db, scryptParams),
		ocr: newOCRKeyStore(db, scryptParams),
		vrf: newVRFKeyStore(db, scryptParams),
	}
}

type Master struct {
	eth *Eth
	csa *CSA
	ocr *OCR
	vrf *VRF
}

func (m *Master) Eth() *Eth {
	return m.eth
}

func (m *Master) CSA() *CSA {
	return m.csa
}

func (m *Master) OCR() *OCR {
	return m.ocr
}

func (m *Master) VRF() *VRF {
	return m.vrf
}
