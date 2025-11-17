package changeset

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	operations2 "github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain/changeset/operations"
)

var _ cldf.ChangeSetV2[CsRegisterNodesWithJDInput] = CsRegisterNodesWithJD{}

// CsRegisterNodesWithJD is a changeset that upserts nodes with Job Distributor for given DONs.
// If they already exist, it ensures they have the correct DON label.
type CsRegisterNodesWithJD struct{}

type CsRegisterNodesWithJDInput struct {
	Domain string               `json:"domain" yaml:"domain"`
	DONs   []offchain.DONConfig `json:"dons" yaml:"dons"`
}

func (cs CsRegisterNodesWithJD) VerifyPreconditions(_ cldf.Environment, cfg CsRegisterNodesWithJDInput) error {
	if len(cfg.DONs) == 0 {
		return errors.New("no DONs provided in input")
	}

	for _, don := range cfg.DONs {
		if don.Name == "" {
			return errors.New("DON name cannot be empty")
		}
	}

	return nil
}

func (cs CsRegisterNodesWithJD) Apply(e cldf.Environment, input CsRegisterNodesWithJDInput) (cldf.ChangesetOutput, error) {
	var terr error
	var out cldf.ChangesetOutput

	for _, don := range input.DONs {
		for _, node := range don.Nodes {
			// run the register node operation for each node in the don
			r, err := operations.ExecuteOperation(
				e.OperationsBundle,
				operations2.JDUpsertNodeOp,
				operations2.JDRegisterNodeOpDeps{Env: e},
				operations2.JDRegisterNodeOpInput{
					Domain:  input.Domain,
					Name:    node.Name,
					CSAKey:  node.CSAKey,
					P2PID:   node.P2PID,
					DONName: don.Name,
					Zone:    node.Zone,
					Labels:  map[string]string{},
				},
			)
			if err != nil {
				// Log error but continue with other nodes
				e.Logger.Errorw("failed to execute register node operation", "don", don.Name, "node", node.Name, "error", err)
				terr = errors.Join(terr, fmt.Errorf("failed to execute register node operation for node %s in don %s: %w", node.Name, don.Name, err))
			}
			out.Reports = append(out.Reports, r.ToGenericReport())
		}
	}
	return out, terr
}

var _ cldf.ChangeSetV2[CsRegisterNodesWithJDInput] = CsRegisterNodesWithJD{}

// CsRegisterNodesWithJDV2 is a changeset that registers nodes with Job Distributor for given DONs.
type CsRegisterNodesWithJDV2 struct{}

type CsRegisterNodesWithJDInputV2 struct {
	Domain string               `json:"domain" yaml:"domain"`
	DONs   []offchain.DONConfig `json:"dons" yaml:"dons"`
}

func (cs CsRegisterNodesWithJDV2) VerifyPreconditions(_ cldf.Environment, cfg CsRegisterNodesWithJDInputV2) error {
	if len(cfg.DONs) == 0 {
		return errors.New("no DONs provided in input")
	}

	for _, don := range cfg.DONs {
		if don.Name == "" {
			return errors.New("DON name cannot be empty")
		}
	}

	return nil
}

func (cs CsRegisterNodesWithJDV2) Apply(e cldf.Environment, input CsRegisterNodesWithJDInputV2) (cldf.ChangesetOutput, error) {
	var terr error
	var out cldf.ChangesetOutput

	for _, don := range input.DONs {
		for _, node := range don.Nodes {
			// run the register node operation for each node in the don
			r, err := operations.ExecuteOperation(
				e.OperationsBundle,
				operations2.JDRegisterNodeOp,
				operations2.JDRegisterNodeOpDeps{Env: e},
				operations2.JDRegisterNodeOpInput{
					Domain:  input.Domain,
					Name:    node.Name,
					CSAKey:  node.CSAKey,
					P2PID:   node.P2PID,
					DONName: don.Name,
					Zone:    node.Zone,
					Labels:  map[string]string{},
				},
			)
			if err != nil {
				// Log error but continue with other nodes
				e.Logger.Errorw("failed to execute register node operation", "don", don.Name, "node", node.Name, "error", err)
				terr = errors.Join(terr, fmt.Errorf("failed to execute register node operation for node %s in don %s: %w", node.Name, don.Name, err))
			}
			out.Reports = append(out.Reports, r.ToGenericReport())
		}
	}
	return out, terr
}
