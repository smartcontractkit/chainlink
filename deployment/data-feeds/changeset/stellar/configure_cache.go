package stellar

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cache "github.com/smartcontractkit/chainlink-stellar/bindings/contracts/data_feeds_cache"

	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/stellar/operation"
)

// FeedPermission is the operator-facing shape of a workflow write-permission.
// AllowedSender is the on-chain caller permitted to invoke on_report — the CRE
// forwarder contract.
type FeedPermission struct {
	AllowedSender        string // Stellar address (C... or G...)
	AllowedWorkflowOwner string // 20-byte hex
	AllowedWorkflowName  string // <= 10 ASCII chars
}

// SetFeedConfigsRequest sets descriptions + workflow permissions for a batch of
// feeds. Permissions apply to every feed in the batch.
type SetFeedConfigsRequest struct {
	ChainSel     uint64
	Qualifier    string
	Version      string
	Admin        string
	DataIDs      []string
	Descriptions []string
	Permissions  []FeedPermission
}

var _ cldf.ChangeSetV2[*SetFeedConfigsRequest] = SetFeedConfigs{}

// SetFeedConfigs sets per-feed descriptions and workflow write-permissions on
// the cache.
type SetFeedConfigs struct{}

func (SetFeedConfigs) VerifyPreconditions(env cldf.Environment, req *SetFeedConfigsRequest) error {
	if err := verifyContractRef(env, req.ChainSel, CacheContract, req.Qualifier, req.Version); err != nil {
		return err
	}
	if err := validateAddress(req.Admin); err != nil {
		return err
	}
	if len(req.DataIDs) == 0 {
		return errors.New("DataIDs cannot be empty")
	}
	if len(req.DataIDs) != len(req.Descriptions) {
		return errors.New("DataIDs and Descriptions must have the same length")
	}
	if len(req.Permissions) == 0 {
		return errors.New("Permissions cannot be empty")
	}
	if _, err := dataIDsToBytes(req.DataIDs); err != nil {
		return err
	}
	if _, err := req.permissions(); err != nil {
		return err
	}
	return nil
}

// permissions converts the operator shape to the generated binding shape.
func (req *SetFeedConfigsRequest) permissions() ([]cache.WorkflowPermission, error) {
	out := make([]cache.WorkflowPermission, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		if err := validateAddress(p.AllowedSender); err != nil {
			return nil, fmt.Errorf("allowed sender: %w", err)
		}
		owner, err := workflowOwnerToBytes(p.AllowedWorkflowOwner)
		if err != nil {
			return nil, err
		}
		name, err := workflowNameToBytes(p.AllowedWorkflowName)
		if err != nil {
			return nil, err
		}
		out = append(out, cache.WorkflowPermission{
			AllowedSender:        p.AllowedSender,
			AllowedWorkflowOwner: owner,
			AllowedWorkflowName:  name,
		})
	}
	return out, nil
}

func (SetFeedConfigs) Apply(env cldf.Environment, req *SetFeedConfigsRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput
	version := semver.MustParse(req.Version)
	d, _, err := resolveContractDeps(env, req.ChainSel, CacheContract, req.Qualifier, version)
	if err != nil {
		return out, err
	}

	ids, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}
	perms, err := req.permissions()
	if err != nil {
		return out, err
	}

	entries := make([]cache.FeedConfigEntry, len(ids))
	for i, id := range ids {
		entries[i] = cache.FeedConfigEntry{
			DataId: id,
			Config: cache.FeedConfig{
				Description:         req.Descriptions[i],
				WorkflowPermissions: perms,
			},
		}
	}

	_, err = operations.ExecuteOperation(env.OperationsBundle, operation.SetFeedConfigs, d.deps, operation.SetFeedConfigsInput{
		ContractID: d.contractID,
		Admin:      req.Admin,
		Entries:    entries,
	})
	return out, err
}
