package config

import solanago "github.com/gagliardetto/solana-go"

type Config struct {
	LogReadTestProgramID solanago.PublicKey
	ExpectedU64Value     uint64
}
