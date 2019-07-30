package orm

import (
	"math/big"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// RuntimeConfig represents configuration values specified while chainlink is running
type RuntimeConfig struct {
	Depot
	ORM *ORM
}

// NewRuntimeConfig returns a runtime instance
func NewRuntimeConfig(depot Depot, orm *ORM) *RuntimeConfig {
	return &RuntimeConfig{
		Depot: depot,
		ORM:   orm,
	}
}

// EthGasPriceDefault represents the default gas price for transactions.
func (r RuntimeConfig) EthGasPriceDefault() *big.Int {
	if str, err := r.ORM.GetConfigValue("EthGasPriceDefault"); err == nil {
		i, _ := new(big.Int).SetString(str, 10)
		return i
	} else if errors.Cause(err) == gorm.ErrRecordNotFound {
		return r.Depot.EthGasPriceDefault()
	}

	// TODO: Return error? Panic? Return default?
	return r.Depot.EthGasPriceDefault()
}
