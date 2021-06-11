package keystore

import (
	"github.com/smartcontractkit/chainlink/core/utils"
	"gorm.io/gorm"
)

func New(db *gorm.DB, scryptParams utils.ScryptParams) *Master {
	return &Master{
		Eth: newEthKeyStore(db, scryptParams),
		CSA: newCSAKeyStore(db, scryptParams),
		OCR: newOCRKeyStore(db, scryptParams),
		VRF: newVRFKeyStore(db, scryptParams),
	}
}

type Master struct {
	Eth *Eth
	CSA *CSA
	OCR *OCR
	VRF *VRF
}
