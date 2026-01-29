package contracts

import (
	contracts3_1 "github.com/smartcontractkit/chainlink/deployment/cre/ocr3/ocr3_1/changeset/operations/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type ConfigureOCR3OpInput = contracts.ConfigureOCR3Input
type ConfigureDKGOpInput = contracts3_1.ConfigureDKGInput

type ConfigureOCR3OpOutput = contracts.ConfigureOCR3OpOutput
type ConfigureOCR3OpDeps = contracts.ConfigureOCR3Deps

type ConfigureDKGOpOutput = contracts3_1.ConfigureDKGOpOutput
type ConfigureDKGOpDeps = contracts3_1.ConfigureDKGDeps

var ConfigureOCR3Op = contracts.ConfigureOCR3

var ConfigureDKGOp = contracts3_1.ConfigureDKG
