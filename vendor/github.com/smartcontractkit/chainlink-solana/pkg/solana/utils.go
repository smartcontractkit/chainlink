package solana

import "github.com/gagliardetto/solana-go"

func LamportsToSol(lamports uint64) float64 {
	return float64(lamports) / float64(solana.LAMPORTS_PER_SOL) // 1_000_000_000
}
