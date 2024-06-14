package monitor

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPromSolBalance(t *testing.T) {
	key := solana.PublicKey{}
	balance := uint64(1_000_000_000)

	monitor := balanceMonitor{chainID: "test-chain"}
	monitor.updateProm(key, balance)

	// happy path test
	promBalance := testutil.ToFloat64(promSolanaBalance.WithLabelValues(key.String(), monitor.chainID, "solana", "SOL"))
	assert.Equal(t, float64(balance)/float64(solana.LAMPORTS_PER_SOL), promBalance)
}
