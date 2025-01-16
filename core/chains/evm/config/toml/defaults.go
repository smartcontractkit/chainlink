package toml

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	cconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config/env"
)

var (
	//go:embed defaults/*.toml
	defaultsFS   embed.FS
	fallback     Chain
	defaults     = map[string]Chain{}
	defaultNames = map[string]string{}

	customDefaults = map[string]Chain{}

	// DefaultIDs is the set of chain ids which have defaults.
	DefaultIDs []*big.Big
)

func init() {
	var (
		fb  *Chain
		err error
	)

	// read all default configs
	DefaultIDs, defaultNames, defaults, fb, err = initDefaults(defaultsFS.ReadDir, defaultsFS.ReadFile, "defaults")
	if err != nil {
		log.Fatalf("failed to read defaults: %s", err)
	}

	if fb == nil {
		log.Fatal("failed to set fallback chain config")
	}

	fallback = *fb

	// check for and apply any overrides
	// read the custom defaults overrides
	dir := env.CustomDefaults.Get()
	if dir == "" {
		// short-circuit; no default overrides provided
		return
	}

	// use evm overrides specifically
	_, _, customDefaults, fb, err = initDefaults(os.ReadDir, os.ReadFile, dir+"/evm")
	if err != nil {
		log.Fatalf("failed to read custom overrides: %s", err)
	}

	if fb != nil {
		fallback = *fb
	}
}

func initDefaults(
	dirReader func(name string) ([]fs.DirEntry, error),
	fileReader func(name string) ([]byte, error),
	root string,
) ([]*big.Big, map[string]string, map[string]Chain, *Chain, error) {
	entries, err := dirReader(root)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var fb *Chain

	ids := make([]*big.Big, 0)
	configs := make(map[string]Chain)
	names := make(map[string]string)

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip directories
			continue
		}

		// read the file to bytes
		path := filepath.Join(root, entry.Name())

		chainID, chain, err := readConfig(path, fileReader)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		if entry.Name() == "fallback.toml" {
			if chainID != nil {
				return nil, nil, nil, nil, fmt.Errorf("fallback ChainID must be nil: found: %s", chainID)
			}

			fb = &chain

			continue
		}

		// ensure ChainID is set
		if chainID == nil {
			return nil, nil, nil, nil, fmt.Errorf("missing ChainID: %s", path)
		}

		ids = append(ids, chainID)

		// ChainID as a default should not be duplicated
		id := chainID.String()
		if _, ok := configs[id]; ok {
			log.Fatalf("%q contains duplicate ChainID: %s", path, id)
		}

		// set lookups
		configs[id] = chain
		names[id] = strings.ReplaceAll(strings.TrimSuffix(entry.Name(), ".toml"), "_", " ")
	}

	// sort IDs in numeric order
	slices.SortFunc(ids, func(a, b *big.Big) int {
		return a.Cmp(b)
	})

	return ids, names, configs, fb, nil
}

func readConfig(path string, reader func(name string) ([]byte, error)) (*big.Big, Chain, error) {
	bts, err := reader(path)
	if err != nil {
		return nil, Chain{}, fmt.Errorf("error reading file: %w", err)
	}

	var config = struct {
		ChainID *big.Big
		Chain
	}{}

	// decode from toml to a chain config
	if err := cconfig.DecodeTOML(bytes.NewReader(bts), &config); err != nil {
		return nil, Chain{}, fmt.Errorf("error in TOML decoding %s: %w", path, err)
	}

	return config.ChainID, config.Chain, nil
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
	if overrides, ok := customDefaults[s]; ok {
		c.SetFrom(&overrides)
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

func ChainTypeForID(chainID *big.Big) (chaintype.ChainType, bool) {
	s := chainID.String()
	if d, ok := defaults[s]; ok {
		return d.ChainType.ChainType(), true
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
	if v := f.LogBroadcasterEnabled; v != nil {
		c.LogBroadcasterEnabled = v
	}
	if v := f.RPCDefaultBatchSize; v != nil {
		c.RPCDefaultBatchSize = v
	}
	if v := f.RPCBlockQueryDelay; v != nil {
		c.RPCBlockQueryDelay = v
	}
	if v := f.FinalizedBlockOffset; v != nil {
		c.FinalizedBlockOffset = v
	}
	if v := f.NoNewFinalizedHeadsThreshold; v != nil {
		c.NoNewFinalizedHeadsThreshold = v
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
	c.Workflow.setFrom(&f.Workflow)
}
