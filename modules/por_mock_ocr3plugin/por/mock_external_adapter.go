package por

import (
	"context"
	"math/big"
	"math/rand"
	"time"
)

// A mock implementation of the ExternalAdapter interface for testing purposes.
type MockExternalAdapterImpl struct {
	latestBlocks Blocks
	chains       []ChainSelector
	counter      int
}

func NewMockExternalAdapterImpl() *MockExternalAdapterImpl {
	chains := []ChainSelector{
		8953668971247136127, // "bitcoin-testnet-rootstock"
		729797994450396300,  // "telos-evm-testnet"
	}

	latestBlocks := make(Blocks)
	for _, chain := range chains {
		// Initialize the latest block for each chain to a random number between 1 and 100.
		latestBlocks[chain] = 0
	}
	return &MockExternalAdapterImpl{
		latestBlocks,
		chains,
		0,
	}
}

func (m *MockExternalAdapterImpl) GetChains(ctx context.Context) ([]ChainSelector, error) {
	return m.chains, nil
}

func (m *MockExternalAdapterImpl) GetPayload(ctx context.Context, blocks Blocks) (ExternalAdapterPayload, error) {
	if m.counter == 10 {
		newSelector := ChainSelector(555555555555555555)
		m.chains = append(m.chains, newSelector)                     // Example of adding a new chain dynamically.
		m.latestBlocks[newSelector] = BlockNumber(rand.Intn(10) + 1) // Assign a random block number to the new chain.
	}
	m.counter++

	mintables := make(Mintables)

	sameChains := (len(blocks) == len(m.chains))
	for _, chain := range m.chains {
		if _, exists := blocks[chain]; !exists {
			sameChains = false
			break
		}
	}

	if !sameChains {
		mintables = nil // If the blocks do not match the chains, return nil mintables.
	} else {
		// Simulate mintable amounts by generating deterministic values based on the block number.
		for chain, block := range blocks {
			mintables[chain] = BlockMintablePair{
				Block:    block,
				Mintable: big.NewInt(int64(block)), // Example: mintable amount is the block number itself.
			}
		}
	}

	reserveInfo := ReserveInfo{
		ReserveAmount: big.NewInt(1000), // Example reserve amount.
		Timestamp:     time.Now(),       // Current time as the reserve timestamp.
	}

	for chain, block := range m.latestBlocks {
		m.latestBlocks[chain] = block + 1 + BlockNumber(rand.Int()%2)
	}

	payload := ExternalAdapterPayload{
		Mintables:    mintables,
		ReserveInfo:  reserveInfo,
		LatestBlocks: m.latestBlocks,
	}

	return payload, nil
}
