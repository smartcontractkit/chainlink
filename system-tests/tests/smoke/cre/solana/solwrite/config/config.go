package config

import solanago "github.com/gagliardetto/solana-go"

type Config struct {
	Receiver           solanago.PublicKey
	ReceiverState      solanago.PublicKey
	WFOwner            [20]byte
	WFName             [10]byte
	FeedID             [16]byte
	ForwarderState     solanago.PublicKey
	ForwarderProgramID solanago.PublicKey
}
