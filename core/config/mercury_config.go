package config

import ocr2models "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"

type Mercury interface {
	Credentials(credName string) *ocr2models.MercuryCredentials
}
