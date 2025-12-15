package config

import solanago "github.com/gagliardetto/solana-go"

type Config struct {
	Receiver           solanago.PublicKey
	ForwarderState     solanago.PublicKey
	ForwarderProgramID solanago.PublicKey
}
