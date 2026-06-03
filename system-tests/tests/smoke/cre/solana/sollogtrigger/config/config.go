package config

import solanago "github.com/gagliardetto/solana-go"

type Config struct {
	LogReadTestProgramID solanago.PublicKey
	ExpectedU64Value     uint64
	ExpectedStrVal       string `yaml:"expected_str_val"`
	ContractIdlJSON      string `yaml:"contract_idl_json"`
	// CPILogTrigger registers the CPI log trigger handler when true, otherwise the direct log trigger handler.
	CPILogTrigger bool `yaml:"cpi_log_trigger"`
}
