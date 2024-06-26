package config

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
)

type ocr2Automation struct {
	c toml.Automation
}

func (o *ocr2Automation) GasLimit() uint32 {
	return *o.c.GasLimit
}

func (o *ocr2Automation) BlockRate() uint32 {
	return *o.c.BlockRate
}

func (o *ocr2Automation) LogLimit() uint32 {
	return *o.c.LogLimit
}

type ocr2Config struct {
	c toml.OCR2
}

func (o *ocr2Config) Automation() OCR2Automation {
	return &ocr2Automation{c: o.c.Automation}
}

func (o *ocr2Config) ContractConfirmations() uint16 {
	return uint16(*o.c.Automation.GasLimit)
}
