// Package presenters allow for the specification and result
// of a Job, its associated TaskSpecs, and every JobRun and TaskRun
// to be returned in a user friendly human readable format.
package presenters

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/logger"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/tidwall/gjson"
)

func LogListeningAddress(address common.Address) string {
	if address == utils.ZeroAddress {
		return "[all]"
	}
	return address.String()
}

func ShowEthBalance(store *store.Store) (string, error) {
	if !store.KeyStore.HasAccounts() {
		logger.Panic("KeyStore must have an account in order to show balance")
	}
	address := store.KeyStore.GetAccount().Address
	balance, err := store.TxManager.GetEthBalance(address)
	if err != nil {
		return "", err
	}
	result := fmt.Sprintf("ETH Balance for %v: %v", address.Hex(), balance)
	if balance == 0 {
		return result, errors.New("0 Balance. Chainlink node not fully functional, please deposit ETH into your address: " + address.Hex())
	}
	return result, nil
}

func ShowLinkBalance(store *store.Store) (string, error) {
	if !store.KeyStore.HasAccounts() {
		logger.Panic("KeyStore must have an account in order to show balance")
	}
	address := store.KeyStore.GetAccount().Address
	linkContractAddress := common.HexToAddress(store.Config.LinkContractAddress)
	balance, err := store.TxManager.GetERC20Balance(address, linkContractAddress)
	if err != nil {
		return "", err
	}
	// Because Eth and Link both use 1e18 precision, we can correct using the same facility
	linkBalance := utils.WeiToEth(balance)
	result := fmt.Sprintf("Link Balance for %v: %v", address.Hex(), linkBalance)
	if balance == big.NewInt(0) {
		return result, errors.New("0 Balance. Chainlink node not fully functional, please deposit LINK into your address: " + address.Hex())
	}
	return result, nil
}

// JobSpec holds the JobSpec definition and each run associated with that Job.
type JobSpec struct {
	models.JobSpec
	Runs []models.JobRun `json:"runs,omitempty"`
}

// MarshalJSON returns the JSON data of the Job and its Initiators.
func (j JobSpec) MarshalJSON() ([]byte, error) {
	type Alias JobSpec
	pis := make([]Initiator, len(j.Initiators))
	for i, modelInitr := range j.Initiators {
		pis[i] = Initiator{modelInitr}
	}
	return json.Marshal(&struct {
		Initiators []Initiator `json:"initiators"`
		Alias
	}{
		pis,
		Alias(j),
	})
}

// FriendlyCreatedAt returns a human-readable string of the Job's
// CreatedAt field.
func (job JobSpec) FriendlyCreatedAt() string {
	return job.CreatedAt.HumanString()
}

// FriendlyStartAt returns a human-readable string of the Job's
// StartAt field.
func (job JobSpec) FriendlyStartAt() string {
	if job.StartAt.Valid {
		return utils.ISO8601UTC(job.StartAt.Time)
	}
	return ""
}

// FriendlyEndAt returns a human-readable string of the Job's
// EndAt field.
func (job JobSpec) FriendlyEndAt() string {
	if job.EndAt.Valid {
		return utils.ISO8601UTC(job.EndAt.Time)
	}
	return ""
}

// FriendlyInitiators returns the list of Initiator types as
// a comma separated string.
func (job JobSpec) FriendlyInitiators() string {
	var initrs []string
	for _, i := range job.Initiators {
		initrs = append(initrs, i.Type)
	}
	return strings.Join(initrs, "\n")
}

// FriendlyTasks returns the list of Task types as a comma
// separated string.
func (job JobSpec) FriendlyTasks() string {
	var tasks []string
	for _, t := range job.Tasks {
		tasks = append(tasks, t.Type)
	}

	return strings.Join(tasks, "\n")
}

// Initiator holds the Job definition's Initiator.
type Initiator struct {
	models.Initiator
}

// MarshalJSON returns the JSON data of the Initiator based
// on its Initiator Type.
func (i Initiator) MarshalJSON() ([]byte, error) {
	switch i.Type {
	case models.InitiatorWeb:
		return json.Marshal(&struct {
			Type string `json:"type"`
		}{
			models.InitiatorWeb,
		})
	case models.InitiatorCron:
		return json.Marshal(&struct {
			Type     string      `json:"type"`
			Schedule models.Cron `json:"schedule"`
		}{
			models.InitiatorCron,
			i.Schedule,
		})
	case models.InitiatorRunAt:
		return json.Marshal(&struct {
			Type string      `json:"type"`
			Time models.Time `json:"time"`
			Ran  bool        `json:"ran"`
		}{
			models.InitiatorRunAt,
			i.Time,
			i.Ran,
		})
	case models.InitiatorEthLog:
		return json.Marshal(&struct {
			Type    string         `json:"type"`
			Address common.Address `json:"address"`
		}{
			models.InitiatorEthLog,
			i.Address,
		})
	case models.InitiatorRunLog:
		return json.Marshal(&struct {
			Type    string         `json:"type"`
			Address common.Address `json:"address"`
		}{
			models.InitiatorRunLog,
			i.Address,
		})
	default:
		return nil, fmt.Errorf("Cannot marshal unsupported initiator type %v", i.Type)
	}
}

// FriendlyRunAt returns a human-readable string for Cron Initiator types.
func (i Initiator) FriendlyRunAt() string {
	if i.Type == models.InitiatorRunAt {
		return i.Time.HumanString()
	}
	return ""
}

var empty_address = common.Address{}.String()

// FriendlyAddress returns the Ethereum address if present, and a blank
// string if not.
func (i Initiator) FriendlyAddress() string {
	if i.IsLogInitiated() {
		return LogListeningAddress(i.Address)
	}
	return ""
}

// TaskSpec holds a task specified in the Job definition.
type TaskSpec struct {
	models.TaskSpec
}

// FriendlyParams returns a map of the TaskSpec's parameters.
func (t TaskSpec) FriendlyParams() (string, string) {
	keys := []string{}
	values := []string{}
	t.Params.ForEach(func(key, value gjson.Result) bool {
		if key.String() != "type" {
			keys = append(keys, key.String())
			values = append(values, value.String())
		}
		return true
	})
	return strings.Join(keys, "\n"), strings.Join(values, "\n")
}
