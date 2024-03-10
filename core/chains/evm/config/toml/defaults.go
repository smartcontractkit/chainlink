package toml

import (
	"bytes"
	"embed"
	"log"
	"path/filepath"
	"slices"
	"strings"

	cconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/common/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

var (
	//go:embed defaults/*.toml
	defaultsFS   embed.FS
	fallback     Chain
	defaults     = map[string]Chain{}
	defaultNames = map[string]string{}

	// DefaultIDs is the set of chain ids which have defaults.
	DefaultIDs []*big.Big
)

func init() {
	fes, err := defaultsFS.ReadDir("defaults")
	if err != nil {
		log.Fatalf("failed to read defaults/: %v", err)
	}
	for _, fe := range fes {
		path := filepath.Join("defaults", fe.Name())
		b, err := defaultsFS.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read %q: %v", path, err)
		}
		var config = struct {
			ChainID *big.Big
			Chain
		}{}

		if err := cconfig.DecodeTOML(bytes.NewReader(b), &config); err != nil {
			log.Fatalf("failed to decode %q: %v", path, err)
		}
		if fe.Name() == "fallback.toml" {
			if config.ChainID != nil {
				log.Fatalf("fallback ChainID must be nil, not: %s", config.ChainID)
			}
			fallback = config.Chain
			continue
		}
		if config.ChainID == nil {
			log.Fatalf("missing ChainID: %s", path)
		}
		DefaultIDs = append(DefaultIDs, config.ChainID)
		id := config.ChainID.String()
		if _, ok := defaults[id]; ok {
			log.Fatalf("%q contains duplicate ChainID: %s", path, id)
		}
		defaults[id] = config.Chain
		defaultNames[id] = strings.ReplaceAll(strings.TrimSuffix(fe.Name(), ".toml"), "_", " ")
	}
	slices.SortFunc(DefaultIDs, func(a, b *big.Big) int {
		return a.Cmp(b)
	})
}

// DefaultsNamed returns the default Chain values, optionally for the given chainID, as well as a name if the chainID is known.
func DefaultsNamed(chainID *big.Big) (c Chain, name string) {
	c.SetFrom(&fallback)
	if chainID == nil {
		return
	}
	s := chainID.String()
	if d, ok := defaults[s]; ok {
		c.SetFrom(&d)
		name = defaultNames[s]
	}
	return
}

// Defaults returns a Chain based on the defaults for chainID and fields from with, applied in order so later Chains
// override earlier ones.
func Defaults(chainID *big.Big, with ...*Chain) Chain {
	c, _ := DefaultsNamed(chainID)
	for _, w := range with {
		c.SetFrom(w)
	}
	return c
}

func ChainTypeForID(chainID *big.Big) (config.ChainType, bool) {
	s := chainID.String()
	if d, ok := defaults[s]; ok {
		if d.ChainType == nil {
			return "", true
		}
		return config.ChainType(*d.ChainType), true
	}
	return "", false
}

// SetFrom updates c with any non-nil values from f.
func (c *Chain) SetFrom(f *Chain) {
	if v := f.AutoCreateKey; v != nil {
		c.AutoCreateKey = v
	}
	if v := f.BlockBackfillDepth; v != nil {
		c.BlockBackfillDepth = v
	}
	if v := f.BlockBackfillSkip; v != nil {
		c.BlockBackfillSkip = v
	}
	if v := f.ChainType; v != nil {
		c.ChainType = v
	}
	if v := f.FinalityDepth; v != nil {
		c.FinalityDepth = v
	}
	if v := f.FinalityTagEnabled; v != nil {
		c.FinalityTagEnabled = v
	}
	if v := f.FlagsContractAddress; v != nil {
		c.FlagsContractAddress = v
	}
	if v := f.LinkContractAddress; v != nil {
		c.LinkContractAddress = v
	}
	if v := f.LogBackfillBatchSize; v != nil {
		c.LogBackfillBatchSize = v
	}
	if v := f.LogPollInterval; v != nil {
		c.LogPollInterval = v
	}
	if v := f.LogKeepBlocksDepth; v != nil {
		c.LogKeepBlocksDepth = v
	}
	if v := f.LogPrunePageSize; v != nil {
		c.LogPrunePageSize = v
	}
	if v := f.BackupLogPollerBlockDelay; v != nil {
		c.BackupLogPollerBlockDelay = v
	}
	if v := f.MinIncomingConfirmations; v != nil {
		c.MinIncomingConfirmations = v
	}
	if v := f.MinContractPayment; v != nil {
		c.MinContractPayment = v
	}
	if v := f.NonceAutoSync; v != nil {
		c.NonceAutoSync = v
	}
	if v := f.NoNewHeadsThreshold; v != nil {
		c.NoNewHeadsThreshold = v
	}
	if v := f.OperatorFactoryAddress; v != nil {
		c.OperatorFactoryAddress = v
	}
	if v := f.RPCDefaultBatchSize; v != nil {
		c.RPCDefaultBatchSize = v
	}
	if v := f.RPCBlockQueryDelay; v != nil {
		c.RPCBlockQueryDelay = v
	}

	c.Transactions.setFrom(&f.Transactions)
	c.BalanceMonitor.setFrom(&f.BalanceMonitor)
	c.GasEstimator.setFrom(&f.GasEstimator)

	if ks := f.KeySpecific; ks != nil {
		for i := range ks {
			v := ks[i]
			if i := slices.IndexFunc(c.KeySpecific, func(k KeySpecific) bool { return k.Key == v.Key }); i == -1 {
				c.KeySpecific = append(c.KeySpecific, v)
			} else {
				c.KeySpecific[i].GasEstimator.setFrom(&v.GasEstimator)
			}
		}
	}

	c.HeadTracker.setFrom(&f.HeadTracker)
	c.NodePool.setFrom(&f.NodePool)
	c.OCR.setFrom(&f.OCR)
	c.OCR2.setFrom(&f.OCR2)
	c.ChainWriter.setFrom(&f.ChainWriter)
}
