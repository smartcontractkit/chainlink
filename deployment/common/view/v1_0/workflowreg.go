package v1_0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	workflow_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

type WorkflowStatus uint8

const (
	WorkflowStatusActive WorkflowStatus = iota
	WorkflowStatusPaused
)

// WorkflowRegistryView is a high-fidelity view of the workflow registry contract.
type WorkflowRegistryView struct {
	types.ContractMetaData
	Workflows           []WorkflowView   `json:"workflows,omitempty"`
	AuthorizedAddresses []common.Address `json:"authorized_addresses,omitempty"`
	AllowedDONs         []uint32         `json:"allowed_dons,omitempty"`
	IsRegistryLocked    bool             `json:"is_registry_locked"`
}

type WorkflowView struct {
	WorkflowID   string         `json:"workflow_id"` // bytes32 stored as hex string (64 hex chars)
	Owner        common.Address `json:"owner"`
	DonID        uint32         `json:"don_id"`
	Status       WorkflowStatus `json:"status"`
	WorkflowName string         `json:"workflow_name"`
	BinaryURL    string         `json:"binary_url"`
	ConfigURL    string         `json:"config_url,omitempty"`
	SecretsURL   string         `json:"secrets_url,omitempty"`
}

func (ws WorkflowStatus) String() string {
	switch ws {
	case WorkflowStatusActive:
		return "ACTIVE"
	case WorkflowStatusPaused:
		return "PAUSED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", ws)
	}
}

func (ws WorkflowStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(ws.String())
}

func (ws *WorkflowStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case "ACTIVE":
		*ws = WorkflowStatusActive
	case "PAUSED":
		*ws = WorkflowStatusPaused
	default:
		return fmt.Errorf("invalid WorkflowStatus value: %q", s)
	}
	return nil
}

func NewWorkflowView(wmd *workflow_registry.WorkflowRegistryWorkflowMetadata) (WorkflowView, error) {
	if wmd == nil {
		return WorkflowView{}, nil
	}

	return WorkflowView{
		WorkflowID:   hex.EncodeToString(wmd.WorkflowID[:]),
		Owner:        wmd.Owner,
		DonID:        wmd.DonID,
		Status:       WorkflowStatus(wmd.Status),
		WorkflowName: wmd.WorkflowName,
		BinaryURL:    wmd.BinaryURL,
		ConfigURL:    wmd.ConfigURL,
		SecretsURL:   wmd.SecretsURL,
	}, nil
}

func GenerateWorkflowRegistryView(wr workflow_registry.WorkflowRegistryInterface) (WorkflowRegistryView, error) {
	// 1) Build up basic contract metadata
	md, err := types.NewContractMetaData(wr, wr.Address())
	if err != nil {
		return WorkflowRegistryView{}, fmt.Errorf("failed to build WorkflowRegistry ContractMetaData: %w", err)
	}

	// 2) Query "getAllAllowedDONs”
	donIDs, err := wr.GetAllAllowedDONs(nil)
	if err != nil {
		return WorkflowRegistryView{}, fmt.Errorf("GetAllAllowedDONs call failed: %w", err)
	}

	// 3) Query "getAllAuthorizedAddresses"
	authAddrs, err := wr.GetAllAuthorizedAddresses(nil)
	if err != nil {
		return WorkflowRegistryView{}, fmt.Errorf("GetAllAuthorizedAddresses call failed: %w", err)
	}

	// 4) Query "isRegistryLocked"
	locked, err := wr.IsRegistryLocked(nil)
	if err != nil {
		return WorkflowRegistryView{}, fmt.Errorf("IsRegistryLocked call failed: %w", err)
	}

	// 5) For each donID, gather their workflows
	var allWorkflowViews []WorkflowView
	for _, donID := range donIDs {
		// Start from index 0, fetch in pages up to the max, in a loop.
		var (
			pageSize = big.NewInt(100) // The registry's default max is 100 per page.
			start    = big.NewInt(0)
		)

		for {
			wmds, err := wr.GetWorkflowMetadataListByDON(nil, donID, start, pageSize)
			if err != nil {
				return WorkflowRegistryView{}, fmt.Errorf("GetWorkflowMetadataListByDON failed for donID %d: %w", donID, err)
			}
			if len(wmds) == 0 {
				break
			}
			// Convert each WorkflowMetadata to a local WorkflowView
			for _, wmd := range wmds {
				wv, err := NewWorkflowView(&wmd)
				if err != nil {
					return WorkflowRegistryView{}, err
				}
				allWorkflowViews = append(allWorkflowViews, wv)
			}
			// If the returned slice is smaller than pageSize, we've exhausted all results
			if len(wmds) < int(pageSize.Int64()) {
				break
			}
			start = new(big.Int).Add(start, pageSize)
		}
	}

	// 6) Build up the final struct
	return WorkflowRegistryView{
		ContractMetaData:    md,
		AllowedDONs:         donIDs,
		AuthorizedAddresses: authAddrs,
		IsRegistryLocked:    locked,
		Workflows:           allWorkflowViews,
	}, nil
}
