package ocrkey

import (
	"github.com/smartcontractkit/offchain-reporting-design/prototype/offchainreporting/types"
)

type PrivateKeys struct{}

var _ types.PrivateKeys = (*PrivateKeys)(nil)
