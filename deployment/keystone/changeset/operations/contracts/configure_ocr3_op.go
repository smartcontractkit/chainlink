package contracts

import (
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"
)

type ConfigureOCR3OpInput = contracts.ConfigureOCR3Input
type ConfigureDKGOpInput = contracts.ConfigureDKGInput

type ConfigureOCR3OpOutput = contracts.ConfigureOCR3OpOutput
type ConfigureOCR3OpDeps = contracts.ConfigureOCR3Deps

type ConfigureDKGOpOutput = contracts.ConfigureDKGOpOutput
type ConfigureDKGOpDeps = contracts.ConfigureDKGDeps

var ConfigureOCR3Op = contracts.ConfigureOCR3

var ConfigureDKGOp = contracts.ConfigureDKG
