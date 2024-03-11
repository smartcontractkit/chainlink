package s4

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

const (
	// defaultExpiration decides how long info will be valid for.
	defaultExpiration = 10 * time.Minute
	// cleanupInterval decides when the expired items in cache will be deleted.
	cleanupInterval = 5 * time.Minute

	getSnapshotCachePrefix = "GetSnapshot"
)

// CachedORM is a cached orm wrapper that implements the ORM interface.
// It adds a cache layer in order to remove unnecessary pressure to the underlaying implementation
type CachedORM struct {
	underlayingORM ORM
	cache          *cache.Cache
	lggr           logger.Logger
}

var _ ORM = (*CachedORM)(nil)

func NewCachedORMWrapper(orm ORM, lggr logger.Logger) *CachedORM {
	return &CachedORM{
		underlayingORM: orm,
		cache:          cache.New(defaultExpiration, cleanupInterval),
		lggr:           lggr,
	}
}

func (c CachedORM) Get(address *ubig.Big, slotId uint, qopts ...pg.QOpt) (*Row, error) {
	return c.underlayingORM.Get(address, slotId, qopts...)
}

func (c CachedORM) Update(row *Row, qopts ...pg.QOpt) error {
	c.deleteRowFromSnapshotCache(row)

	return c.underlayingORM.Update(row, qopts...)
}

func (c CachedORM) DeleteExpired(limit uint, utcNow time.Time, qopts ...pg.QOpt) (int64, error) {
	deletedRows, err := c.underlayingORM.DeleteExpired(limit, utcNow, qopts...)
	if err != nil {
		return 0, err
	}

	if deletedRows > 0 {
		c.cache.Flush()
	}

	return deletedRows, nil
}

func (c CachedORM) GetSnapshot(addressRange *AddressRange, qopts ...pg.QOpt) ([]*SnapshotRow, error) {
	key := fmt.Sprintf("%s_%s_%s", getSnapshotCachePrefix, addressRange.MinAddress.String(), addressRange.MaxAddress.String())

	cached, found := c.cache.Get(key)
	if found {
		return cached.([]*SnapshotRow), nil
	}

	c.lggr.Debug("Snapshot not found in cache, fetching it from underlaying implementation")
	data, err := c.underlayingORM.GetSnapshot(addressRange, qopts...)
	if err != nil {
		return nil, err
	}
	c.cache.Set(key, data, defaultExpiration)

	return data, nil
}

func (c CachedORM) GetUnconfirmedRows(limit uint, qopts ...pg.QOpt) ([]*Row, error) {
	return c.underlayingORM.GetUnconfirmedRows(limit, qopts...)
}

// deleteRowFromSnapshotCache will clean the cache for every snapshot that would involve a given row
// in case of an error parsing a key it will also delete the key from the cache
func (c CachedORM) deleteRowFromSnapshotCache(row *Row) {
	for key := range c.cache.Items() {
		keyParts := strings.Split(key, "_")
		if len(keyParts) != 3 {
			continue
		}

		if keyParts[0] != getSnapshotCachePrefix {
			continue
		}

		minAddress, ok := new(big.Int).SetString(keyParts[1], 10)
		if !ok {
			c.lggr.Errorf("error while converting minAddress string: %s to big.Int, deleting key %q", keyParts[1], key)
			c.cache.Delete(key)
			continue
		}

		maxAddress, ok := new(big.Int).SetString(keyParts[2], 10)
		if !ok {
			c.lggr.Errorf("error while converting minAddress string: %s to big.Int, deleting key %q ", keyParts[2], key)
			c.cache.Delete(key)
			continue
		}

		if row.Address.ToInt().Cmp(minAddress) >= 0 && row.Address.ToInt().Cmp(maxAddress) <= 0 {
			c.cache.Delete(key)
		}
	}
}
