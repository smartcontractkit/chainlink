package v2

import (
	"bytes"
	"embed"
	"log"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	//go:embed defaults/*.toml
	defaultsFS   embed.FS
	fallback     Chain
	defaults     = map[string]Chain{}
	defaultNames = map[string]string{}

	// DefaultIDs is the set of chain ids which have defaults.
	DefaultIDs []*utils.Big
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
			ChainID *utils.Big
			Chain
		}{}
		d := toml.NewDecoder(bytes.NewReader(b)).DisallowUnknownFields()
		if err := d.Decode(&config); err != nil {
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
	slices.SortFunc(DefaultIDs, func(a, b *utils.Big) bool {
		return a.Cmp(b) < 0
	})
}

// Defaults returns the default Chain values, optionally for the given chainID, as well as a name if the chainID is known.
func Defaults(chainID *utils.Big) (c Chain, name string) {
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

// SetFrom updates c with any non-nil values from f.
func (c *Chain) SetFrom(f *Chain) {
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
	if v := f.MaxInFlightTransactions; v != nil {
		c.MaxInFlightTransactions = v
	}
	if v := f.MaxQueuedTransactions; v != nil {
		c.MaxQueuedTransactions = v
	}
	if v := f.MinIncomingConfirmations; v != nil {
		c.MinIncomingConfirmations = v
	}
	if v := f.MinimumContractPayment; v != nil {
		c.MinimumContractPayment = v
	}
	if v := f.NonceAutoSync; v != nil {
		c.NonceAutoSync = v
	}
	if v := f.OperatorFactoryAddress; v != nil {
		c.OperatorFactoryAddress = v
	}
	if v := f.RPCDefaultBatchSize; v != nil {
		c.RPCDefaultBatchSize = v
	}
	if v := f.TxReaperInterval; v != nil {
		c.TxReaperInterval = v
	}
	if v := f.TxReaperThreshold; v != nil {
		c.TxReaperThreshold = v
	}
	if v := f.TxResendAfterThreshold; v != nil {
		c.TxResendAfterThreshold = v
	}
	if v := f.UseForwarders; v != nil {
		c.UseForwarders = v
	}
	if b := f.BalanceMonitor; b != nil {
		if c.BalanceMonitor == nil {
			c.BalanceMonitor = &BalanceMonitor{}
		}
		if v := b.Enabled; v != nil {
			c.BalanceMonitor.Enabled = v
		}
		if v := b.BlockDelay; v != nil {
			c.BalanceMonitor.BlockDelay = v
		}
	}
	if g := f.GasEstimator; g != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		if v := g.Mode; v != nil {
			c.GasEstimator.Mode = v
		}
		if v := g.EIP1559DynamicFees; v != nil {
			c.GasEstimator.EIP1559DynamicFees = v
		}
		if v := g.BumpPercent; v != nil {
			c.GasEstimator.BumpPercent = v
		}
		if v := g.BumpThreshold; v != nil {
			c.GasEstimator.BumpThreshold = v
		}
		if v := g.BumpTxDepth; v != nil {
			c.GasEstimator.BumpTxDepth = v
		}
		if v := g.BumpMin; v != nil {
			c.GasEstimator.BumpMin = v
		}
		if v := g.FeeCapDefault; v != nil {
			c.GasEstimator.FeeCapDefault = v
		}
		if v := g.LimitDefault; v != nil {
			c.GasEstimator.LimitDefault = v
		}
		if v := g.LimitMultiplier; v != nil {
			c.GasEstimator.LimitMultiplier = v
		}
		if v := g.LimitTransfer; v != nil {
			c.GasEstimator.LimitTransfer = v
		}
		if v := g.LimitOCRJobType; v != nil {
			c.GasEstimator.LimitOCRJobType = v
		}
		if v := g.LimitDRJobType; v != nil {
			c.GasEstimator.LimitDRJobType = v
		}
		if v := g.LimitVRFJobType; v != nil {
			c.GasEstimator.LimitVRFJobType = v
		}
		if v := g.LimitFMJobType; v != nil {
			c.GasEstimator.LimitFMJobType = v
		}
		if v := g.LimitKeeperJobType; v != nil {
			c.GasEstimator.LimitKeeperJobType = v
		}
		if v := g.PriceDefault; v != nil {
			c.GasEstimator.PriceDefault = v
		}
		if v := g.TipCapDefault; v != nil {
			c.GasEstimator.TipCapDefault = v
		}
		if v := g.TipCapMinimum; v != nil {
			c.GasEstimator.TipCapMinimum = v
		}
		if v := g.PriceMax; v != nil {
			c.GasEstimator.PriceMax = v
		}
		if v := g.PriceMin; v != nil {
			c.GasEstimator.PriceMin = v
		}
		if b := g.BlockHistory; b != nil {
			if c.GasEstimator.BlockHistory == nil {
				c.GasEstimator.BlockHistory = &BlockHistoryEstimator{}
			}
			if v := b.BatchSize; v != nil {
				c.GasEstimator.BlockHistory.BatchSize = v
			}
			if v := b.BlockDelay; v != nil {
				c.GasEstimator.BlockHistory.BlockDelay = v
			}
			if v := b.BlockHistorySize; v != nil {
				c.GasEstimator.BlockHistory.BlockHistorySize = v
			}
			if v := b.EIP1559FeeCapBufferBlocks; v != nil {
				c.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks = v
			}
			if v := b.TransactionPercentile; v != nil {
				c.GasEstimator.BlockHistory.TransactionPercentile = v
			}
		}
	}
	if ks := f.KeySpecific; ks != nil {
		for _, v := range ks {
			if i := slices.IndexFunc(c.KeySpecific, func(k KeySpecific) bool { return k.Key == v.Key }); i == -1 {
				c.KeySpecific = append(c.KeySpecific, v)
			} else {
				if v := v.GasEstimator; v != nil {
					c.KeySpecific[i].GasEstimator = v
				}
			}
		}
	}
	if h := f.HeadTracker; h != nil {
		if c.HeadTracker == nil {
			c.HeadTracker = &HeadTracker{}
		}
		if v := h.BlockEmissionIdleWarningThreshold; v != nil {
			c.HeadTracker.BlockEmissionIdleWarningThreshold = v
		}
		if v := h.HistoryDepth; v != nil {
			c.HeadTracker.HistoryDepth = v
		}
		if v := h.MaxBufferSize; v != nil {
			c.HeadTracker.MaxBufferSize = v
		}
		if v := h.SamplingInterval; v != nil {
			c.HeadTracker.SamplingInterval = v
		}
	}
	if n := f.NodePool; n != nil {
		if c.NodePool == nil {
			c.NodePool = &NodePool{}
		}
		if v := n.NoNewHeadsThreshold; v != nil {
			c.NodePool.NoNewHeadsThreshold = v
		}
		if v := n.PollFailureThreshold; v != nil {
			c.NodePool.PollFailureThreshold = v
		}
		if v := n.PollInterval; v != nil {
			c.NodePool.PollInterval = v
		}
	}
	if o := f.OCR; o != nil {
		if c.OCR == nil {
			c.OCR = &OCR{}
		}
		if v := o.ContractConfirmations; v != nil {
			c.OCR.ContractConfirmations = v
		}
		if v := o.ContractTransmitterTransmitTimeout; v != nil {
			c.OCR.ContractTransmitterTransmitTimeout = v
		}
		if v := o.DatabaseTimeout; v != nil {
			c.OCR.DatabaseTimeout = v
		}
		if v := o.ObservationTimeout; v != nil {
			c.OCR.ObservationTimeout = v
		}
		if v := o.ObservationGracePeriod; v != nil {
			c.OCR.ObservationGracePeriod = v
		}
	}
}
