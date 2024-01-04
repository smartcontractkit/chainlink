package rebalancer

import (
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func TestPendingTransfersCache(t *testing.T) {
	c := NewPendingTransfersCache()

	from1to2 := models.NewTransfer(models.NetworkID(1), models.NetworkID(2), big.NewInt(10))
	from1to3 := models.NewTransfer(models.NetworkID(1), models.NetworkID(3), big.NewInt(20))
	from2to3 := models.NewTransfer(models.NetworkID(2), models.NetworkID(3), big.NewInt(30))

	c.Add([]models.PendingTransfer{
		models.NewPendingTransfer(from1to2),
		models.NewPendingTransfer(from1to3),
	})
	assert.True(t, c.ContainsTransfer(from1to2))
	assert.True(t, c.ContainsTransfer(from1to3))
	assert.False(t, c.ContainsTransfer(from2to3))

	c.Add([]models.PendingTransfer{models.NewPendingTransfer(from2to3)})
	assert.True(t, c.ContainsTransfer(from1to2), "adding a new item should not affect existing items")
	assert.True(t, c.ContainsTransfer(from1to3))
	assert.True(t, c.ContainsTransfer(from2to3))

	c.Set([]models.PendingTransfer{models.NewPendingTransfer(from2to3)})
	assert.False(t, c.ContainsTransfer(from1to2), "set should delete existing items")
	assert.False(t, c.ContainsTransfer(from1to3))
	assert.True(t, c.ContainsTransfer(from2to3))
}

func TestPendingTransfersThreadSafety(t *testing.T) {
	const numWorkers = 30
	const numOps = 20

	c := NewPendingTransfersCache()

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			runRandomRegistryOperations(c, numOps)
		}()
	}
	wg.Wait()
}

func runRandomRegistryOperations(c *PendingTransfersCache, numOps int) {
	ops := []string{"add", "set", "contains"}
	for i := 0; i < numOps; i++ {
		tr1 := models.NewTransfer(models.NetworkID(rand.Intn(numOps)), models.NetworkID(rand.Intn(numOps)), big.NewInt(int64(numOps*rand.Intn(10))))
		tr2 := models.NewTransfer(models.NetworkID(rand.Intn(numOps)), models.NetworkID(rand.Intn(numOps)), big.NewInt(int64(numOps*rand.Intn(10))))
		transfers := []models.PendingTransfer{models.NewPendingTransfer(tr1), models.NewPendingTransfer(tr2)}

		switch ops[rand.Intn(len(ops))] {
		case "add":
			c.Add(transfers)
		case "set":
			c.Set(transfers)
		case "contains":
			c.ContainsTransfer(tr1)
			c.ContainsTransfer(tr2)
		}
	}
}
