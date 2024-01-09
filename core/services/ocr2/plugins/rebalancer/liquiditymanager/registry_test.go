package liquiditymanager

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()
	assert.Empty(t, r.GetAll())

	netID1 := models.NetworkID(1)
	addr1 := models.Address(utils.RandomAddress())

	netID2 := models.NetworkID(2)
	addr2 := models.Address(utils.RandomAddress())

	r.Add(netID1, addr1)
	assert.Len(t, r.GetAll(), 1)
	addr, exists := r.Get(netID1)
	assert.True(t, exists)
	assert.Equal(t, addr1, addr)

	_, exists = r.Get(netID2)
	assert.False(t, exists)

	r.Add(netID1, addr2)
	assert.Len(t, r.GetAll(), 1)
	addr, exists = r.Get(netID1)
	assert.True(t, exists)
	assert.Equal(t, addr2, addr, "address should be overwritten")
}

func TestRegistryThreadSafety(t *testing.T) {
	const numWorkers = 30
	const numOps = 20

	r := NewRegistry()

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			runRandomRegistryOperations(r, numOps)
		}()
	}
	wg.Wait()
}

func runRandomRegistryOperations(r *Registry, numOps int) {
	ops := []string{"add", "get", "getAll"}
	for i := 0; i < numOps; i++ {
		switch ops[rand.Intn(len(ops))] {
		case "add":
			r.Add(models.NetworkID(rand.Intn(numOps)), models.Address(utils.RandomAddress()))
		case "get":
			_, _ = r.Get(models.NetworkID(rand.Intn(numOps)))
		case "getAll":
			r.GetAll()
		}
	}
}
