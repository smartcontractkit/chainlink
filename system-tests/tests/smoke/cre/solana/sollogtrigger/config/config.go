package config

import solanago "github.com/gagliardetto/solana-go"

type Config struct {
	// LogReadTestProgramID is the program ID of the log_read_test contract
	LogReadTestProgramID solanago.PublicKey
}
