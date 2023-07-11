package forwarders

import (
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Forwarder is the struct for Forwarder Addresses
type Forwarder struct {
	ID         int64
	Address    common.Address
	EVMChainID utils.Big
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var SupportedPlugins = []job.OCR2PluginType{job.Median, job.DKG, job.OCR2VRF, job.OCR2Keeper, job.OCR2Functions}
