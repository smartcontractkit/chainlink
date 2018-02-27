package store_test

import (
	"math/big"
	"syscall"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestGracefulShutdown(t *testing.T) {
	RegisterTestingT(t)
	store, cleanup := cltest.NewStore()
	defer cleanup()

	var completed bool
	store.Exiter = func(code int) {
		completed = true
	}

	store.Start()
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)

	Eventually(func() bool {
		return completed
	}).Should(BeTrue())
}

func TestConfigDefaults(t *testing.T) {
	config := strpkg.NewConfig()
	assert.Equal(t, uint64(0), config.ChainID)
	assert.Equal(t, *big.NewInt(20000000000), config.EthGasPriceDefault)
}

func TestHeadTracker_New(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	assert.Nil(t, store.Save(&models.BlockHeader{cltest.BigHexInt(1)}))
	last := models.BlockHeader{cltest.BigHexInt(10)}
	assert.Nil(t, store.Save(&last))
	assert.Nil(t, store.Save(&models.BlockHeader{cltest.BigHexInt(2)}))

	ht, err := strpkg.NewHeadTracker(store.ORM)
	assert.Nil(t, err)
	assert.Equal(t, last.Number, ht.Get().Number)
}

func TestHeadTracker_Get(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	initial := models.BlockHeader{cltest.BigHexInt(1)}
	assert.Nil(t, store.Save(&initial))

	tests := []struct {
		name      string
		toSave    *models.BlockHeader
		want      hexutil.Big
		wantError bool
	}{
		// order matters
		{"greater", &models.BlockHeader{cltest.BigHexInt(2)}, cltest.BigHexInt(2), false},
		{"less than", &models.BlockHeader{cltest.BigHexInt(1)}, cltest.BigHexInt(2), false},
		{"zero", &models.BlockHeader{cltest.BigHexInt(0)}, cltest.BigHexInt(2), true},
		{"nil", nil, cltest.BigHexInt(2), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ht, err := strpkg.NewHeadTracker(store.ORM)
			assert.Nil(t, err)
			err = ht.Save(test.toSave)
			if test.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, test.want, ht.Get().Number)
		})
	}
}
