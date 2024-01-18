package rebalancer

import (
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/models"
)

func TestPendingTransfersCache(t *testing.T) {
	c := NewPendingTransfersCache()

	date := time.Now()
	from1to2 := models.NewTransfer(models.NetworkSelector(1), models.NetworkSelector(2), big.NewInt(10), date, []byte{})
	from1to3 := models.NewTransfer(models.NetworkSelector(1), models.NetworkSelector(3), big.NewInt(20), date, []byte{})
	from2to3 := models.NewTransfer(models.NetworkSelector(2), models.NetworkSelector(3), big.NewInt(30), date, []byte{})

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
	date := time.Now()
	for i := 0; i < numOps; i++ {
		tr1 := models.NewTransfer(models.NetworkSelector(rand.Intn(numOps)), models.NetworkSelector(rand.Intn(numOps)), big.NewInt(int64(numOps*rand.Intn(10))), date, []byte{})
		tr2 := models.NewTransfer(models.NetworkSelector(rand.Intn(numOps)), models.NetworkSelector(rand.Intn(numOps)), big.NewInt(int64(numOps*rand.Intn(10))), date, []byte{})
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
