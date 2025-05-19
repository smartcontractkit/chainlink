package v1_0

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
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

// WorkflowRegistryError is a custom error type for errors that occur while building the workflow registry view.
type WorkflowRegistryError struct {
	Operation string
	Err       error
}

func (e *WorkflowRegistryError) Error() string {
	return fmt.Sprintf("%s: %v", e.Operation, e.Err)
}

func (e *WorkflowRegistryError) Unwrap() error {
	return e.Err
}

func NewWorkflowView(wmd *workflow_registry.WorkflowRegistryWorkflowMetadata) (WorkflowView, error) {
	if wmd == nil {
		return WorkflowView{}, &WorkflowRegistryError{
			Operation: "converting workflow metadata",
			Err:       errors.New("workflow metadata is nil"),
		}
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

// GenerateWorkflowRegistryView builds a WorkflowRegistryView from the provided workflow registry interface.
// Instead of aborting on the first error, it collects errors in a slice.
func GenerateWorkflowRegistryView(wr workflow_registry.WorkflowRegistryInterface) (WorkflowRegistryView, []error) {
	var errs []error

	// 1) Build up basic contract metadata.
	md, err := types.NewContractMetaData(wr, wr.Address())
	if err != nil {
		errs = append(errs, &WorkflowRegistryError{
			Operation: "failed to build WorkflowRegistry ContractMetaData",
			Err:       err,
		})
	}

	// 2) Query "getAllAllowedDONs”.
	donIDs, err := wr.GetAllAllowedDONs(nil)
	if err != nil {
		errs = append(errs, &WorkflowRegistryError{
			Operation: "GetAllAllowedDONs call failed",
			Err:       err,
		})
	}

	// 3) Query "getAllAuthorizedAddresses".
	authAddrs, err := wr.GetAllAuthorizedAddresses(nil)
	if err != nil {
		errs = append(errs, &WorkflowRegistryError{
			Operation: "GetAllAuthorizedAddresses call failed",
			Err:       err,
		})
	}

	// 4) Query "isRegistryLocked".
	locked, err := wr.IsRegistryLocked(nil)
	if err != nil {
		errs = append(errs, &WorkflowRegistryError{
			Operation: "IsRegistryLocked call failed",
			Err:       err,
		})
	}

	// 5) For each DON ID, gather their workflows.
	var allWorkflowViews []WorkflowView
	for _, donID := range donIDs {
		// Start from index 0; fetch in pages up to the max.
		// The registry's default max is 100 per page.
		pageSize := big.NewInt(100)
		start := big.NewInt(0)

		for {
			wmds, err := wr.GetWorkflowMetadataListByDON(nil, donID, start, pageSize)
			if err != nil {
				errs = append(errs, &WorkflowRegistryError{
					Operation: fmt.Sprintf("GetWorkflowMetadataListByDON failed for donID %d", donID),
					Err:       err,
				})
				break
			}
			if len(wmds) == 0 {
				break
			}
			// Convert each WorkflowMetadata to a local WorkflowView.
			for _, wmd := range wmds {
				wv, err := NewWorkflowView(&wmd)
				if err != nil {
					errs = append(errs, &WorkflowRegistryError{
						Operation: "failed to convert workflow metadata",
						Err:       err,
					})
					continue
				}
				allWorkflowViews = append(allWorkflowViews, wv)
			}
			// If the returned slice is smaller than pageSize, we've exhausted all results.
			if len(wmds) < int(pageSize.Int64()) {
				break
			}
			start = new(big.Int).Add(start, pageSize)
		}
	}

	// 6) Build up the final struct.
	view := WorkflowRegistryView{
		ContractMetaData:    md,
		AllowedDONs:         donIDs,
		AuthorizedAddresses: authAddrs,
		IsRegistryLocked:    locked,
		Workflows:           allWorkflowViews,
	}
	return view, errs
}
