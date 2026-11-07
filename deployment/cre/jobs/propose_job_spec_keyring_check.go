package jobs

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"
	"github.com/smartcontractkit/chainlink/deployment"
	job_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

// templatesWithOCRKeys are the templates whose specs embed a per-node OCR key bundle ID.
var templatesWithOCRKeys = []job_types.JobSpecTemplate{
	job_types.OCR3, job_types.Consensus, job_types.Ring,
	job_types.EVM, job_types.Aptos, job_types.Solana, job_types.Stellar,
}

// verifyOCRKeyrings checks that, for every target node, the EVM OCR2 key bundle recorded in JD has
// the onchain signing address the caller expects. It is a no-op unless ExpectedSigners is set.
func verifyOCRKeyrings(e cldf.Environment, config ProposeJobSpecInput) error {
	if len(config.ExpectedSigners) == 0 || config.SkipKeyringCheck || !slices.Contains(templatesWithOCRKeys, config.Template) {
		return nil
	}
	if e.Offchain == nil {
		return errors.New("keyring check needs a job distributor client; unset expectedSigners or set skipKeyringCheck to bypass")
	}

	nodeInfos, err := fetchTargetNodeInfos(e, config)
	if err != nil {
		return err
	}

	sel, err := specChainSelectors(config)
	if err != nil {
		return err
	}

	return checkKeyringsAgainstExpected(nodeInfos, config.ExpectedSigners, sel)
}

// fetchTargetNodeInfos mirrors the JD node selection used when the job is proposed.
func fetchTargetNodeInfos(e cldf.Environment, config ProposeJobSpecInput) (deployment.Nodes, error) {
	donName := jdDONName(config)
	filter := targetNodeFilter(config)

	jdNodes, err := offchain.FetchNodesFromJD(e.GetContext(), e.Offchain, filter)
	if err != nil {
		return nil, fmt.Errorf("keyring check: failed to fetch nodes from JD: %w", err)
	}
	if len(jdNodes) == 0 {
		return nil, fmt.Errorf("keyring check: no nodes found on JD for DON %q", donName)
	}
	ids := make([]string, len(jdNodes))
	for i, n := range jdNodes {
		ids[i] = n.Id
	}
	infos, err := deployment.NodeInfo(ids, e.Offchain)
	if err != nil {
		return nil, fmt.Errorf("keyring check: failed to fetch node infos: %w", err)
	}

	return infos, nil
}

// checkKeyringsAgainstExpected requires every non-bootstrap node to have an EVM OCR2 key bundle in
// JD with a non-empty ID and an onchain signing address equal to the expected signer for its CSA
// key. A stale JD chain config (e.g. after a key rotation) would otherwise bake the wrong bundle
// ID into the spec. For the non-EVM chain selectors the spec uses, it also requires a non-empty
// bundle ID on every node and no bundle ID shared between nodes; nodes.json has no committed
// non-EVM signers to compare against.
func checkKeyringsAgainstExpected(nodeInfos deployment.Nodes, expected map[string]string, sel specSelectors) error {
	var problems []string
	for _, n := range nodeInfos {
		if n.IsBootstrap {
			continue
		}
		signer, ok := expected[n.CSAKey]
		if !ok {
			problems = append(problems, fmt.Sprintf("%s: no expected signer for csa %s", n.Name, n.CSAKey))

			continue
		}
		want, err := normalizeSigner(signer)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", n.Name, err))

			continue
		}

		problems = append(problems, evmBundleProblems(n, want, sel.EVM)...)
	}
	problems = append(problems, nonEVMBundleProblems(nodeInfos, sel.NonEVM)...)
	if len(problems) == 0 {
		return nil
	}
	slices.Sort(problems)

	return fmt.Errorf("OCR key bundles in JD do not match the expected keyrings; re-sync JD chain configs or the expected signers (or set skipKeyringCheck):\n  %s", strings.Join(problems, "\n  "))
}

// normalizeSigner reduces a signer (20-byte address, optionally 0x-prefixed and zero-padded to 32
// bytes) to 40 lowercase hex characters.
func normalizeSigner(signer string) (string, error) {
	s := strings.ToLower(strings.TrimPrefix(signer, "0x"))
	if len(s) != 40 && len(s) != 64 {
		return "", fmt.Errorf("expected signer %q is not a 20 or 32 byte hex value", signer)
	}
	if _, err := hex.DecodeString(s); err != nil {
		return "", fmt.Errorf("expected signer %q is not hex: %w", signer, err)
	}

	return s[:40], nil
}

func jdDONName(config ProposeJobSpecInput) string {
	if config.JDDONName != "" {
		return config.JDDONName
	}

	return config.DONName
}

// targetNodeFilter reproduces the JD node selection of the operation that proposes the template, so
// the check validates the same nodes the proposal targets.
func targetNodeFilter(config ProposeJobSpecInput) *nodev1.ListNodesRequest_Filter {
	switch config.Template {
	case job_types.OCR3, job_types.Ring:
		filter := &nodev1.ListNodesRequest_Filter{}
		for _, f := range config.DONFilters {
			filter = f.AddToFilter(filter)
		}

		return offchain.TargetDONFilter{Key: "type", Value: "plugin"}.AddToFilter(filter)
	default:
		nodeType := job_ops.PluginNodeType
		filter := &nodev1.ListNodesRequest_Filter{
			Selectors: []*ptypes.Selector{
				{Key: "don-" + jdDONName(config), Op: ptypes.SelectorOp_EXIST},
				{Key: "environment", Op: ptypes.SelectorOp_EQ, Value: &config.Environment},
				{Key: "product", Op: ptypes.SelectorOp_EQ, Value: &config.Domain},
				{Key: "type", Op: ptypes.SelectorOp_EQ, Value: &nodeType},
			},
		}
		for _, f := range config.DONFilters {
			filter = f.AddToFilter(filter)
		}
		for _, f := range config.DONFilters {
			filter = f.AddToFilterIfNotPresent(filter)
		}

		return filter
	}
}

// specSelectors are the chain selectors the spec reads OCR key bundles from. EVM is zero when the
// template does not pin one, in which case any EVM bundle may match.
type specSelectors struct {
	EVM    uint64
	NonEVM []uint64
}

func specChainSelectors(config ProposeJobSpecInput) (specSelectors, error) {
	var evm pkg.ChainSelector
	var raw []pkg.ChainSelector
	switch config.Template { //nolint:exhaustive // only templates in templatesWithOCRKeys reach here
	case job_types.OCR3:
		in, err := config.Inputs.ToOCR3JobConfigInput()
		if err != nil {
			return specSelectors{}, fmt.Errorf("keyring check: failed to parse OCR3 inputs: %w", err)
		}
		evm = in.ChainSelectorEVM
		raw = []pkg.ChainSelector{in.ChainSelectorAptos, in.ChainSelectorSolana}
	case job_types.Consensus, job_types.EVM, job_types.Aptos, job_types.Solana, job_types.Stellar:
		job, err := config.Inputs.ToStandardCapabilityJob(config.JobName)
		if err != nil {
			return specSelectors{}, fmt.Errorf("keyring check: failed to parse standard capability inputs: %w", err)
		}
		// Matches generateOracleFactory: the OCR contract chain wins over the EVM selector.
		evm = job.ChainSelectorEVM
		if job.OCRChainSelector != 0 {
			evm = job.OCRChainSelector
		}
		raw = []pkg.ChainSelector{job.ChainSelectorAptos, job.ChainSelectorSolana, job.ChainSelectorStellar}
	case job_types.Ring:
		var in pkg.RingJobConfigInput
		if err := config.Inputs.UnmarshalTo(&in); err != nil {
			return specSelectors{}, fmt.Errorf("keyring check: failed to parse Ring inputs: %w", err)
		}
		evm = in.ChainSelectorEVM
	}

	out := specSelectors{EVM: uint64(evm)}
	for _, s := range raw {
		if s != 0 {
			out.NonEVM = append(out.NonEVM, uint64(s))
		}
	}

	return out, nil
}

// evmBundleProblems requires the node's EVM OCR bundle for the spec's chain selector (or, when the
// template pins none, any EVM chain) to have a non-empty ID and the wanted signing address.
func evmBundleProblems(n deployment.Node, want string, evmSel uint64) []string {
	var problems []string
	var have []string
	matched := false
	for details, cfg := range n.SelToOCRConfig {
		if evmSel != 0 && details.ChainSelector != evmSel {
			continue
		}
		family, err := chainsel.GetSelectorFamily(details.ChainSelector)
		if err != nil || family != chainsel.FamilyEVM {
			continue
		}
		if cfg.KeyBundleID == "" {
			problems = append(problems, fmt.Sprintf("%s: JD has an empty OCR key bundle ID for chain selector %d", n.Name, details.ChainSelector))

			continue
		}
		addr := hex.EncodeToString(cfg.OnchainPublicKey)
		have = append(have, addr)
		if addr == want {
			matched = true
		}
	}
	switch {
	case matched:
	case len(have) == 0:
		problems = append(problems, n.Name+": no EVM OCR key bundle in JD")
	default:
		slices.Sort(have)
		problems = append(problems, fmt.Sprintf("%s: expected signer %s matches no JD EVM bundle (JD has %s)", n.Name, want, strings.Join(have, ", ")))
	}

	return problems
}

func nonEVMBundleProblems(nodeInfos deployment.Nodes, selectors []uint64) []string {
	var problems []string
	for _, sel := range selectors {
		owners := make(map[string][]string)
		for _, n := range nodeInfos {
			if n.IsBootstrap {
				continue
			}
			cfg, ok := n.OCRConfigForChainSelector(sel)
			switch {
			case !ok:
				problems = append(problems, fmt.Sprintf("%s: no OCR key bundle in JD for chain selector %d", n.Name, sel))
			case cfg.KeyBundleID == "":
				problems = append(problems, fmt.Sprintf("%s: JD has an empty OCR key bundle ID for chain selector %d", n.Name, sel))
			default:
				owners[cfg.KeyBundleID] = append(owners[cfg.KeyBundleID], n.Name)
			}
		}
		for id, names := range owners {
			if len(names) > 1 {
				slices.Sort(names)
				problems = append(problems, fmt.Sprintf("chain selector %d: bundle ID %s is shared by nodes %s", sel, id, strings.Join(names, ", ")))
			}
		}
	}

	return problems
}
