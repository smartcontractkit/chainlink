package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

// UpkeepTriggerID returns the identifier using the given upkeepID and trigger.
// It follows the same logic as the contract, but performs it locally.
func UpkeepTriggerID(id *big.Int, trigger []byte) string {
	idBytes := append(id.Bytes(), trigger...)
	triggerIDBytes := crypto.Keccak256(idBytes)
	return hex.EncodeToString(triggerIDBytes)
}

// UpkeepWorkID returns the identifier using the given upkeepID and trigger extension(tx hash and log index).
func UpkeepWorkID(id *big.Int, trigger ocr2keepers.Trigger) (string, error) {
	extensionBytes, err := json.Marshal(trigger.Extension)
	if err != nil {
		return "", err
	}

	// TODO (auto-4314): Ensure it works with conditionals and add unit tests
	combined := fmt.Sprintf("%s%s", id, extensionBytes)
	hash := crypto.Keccak256([]byte(combined))
	return hex.EncodeToString(hash[:]), nil
}

func NewUpkeepPayload(uid *big.Int, tp int, trigger ocr2keepers.Trigger, checkData []byte) (ocr2keepers.UpkeepPayload, error) {
	// construct payload
	p := ocr2keepers.UpkeepPayload{
		Upkeep: ocr2keepers.ConfiguredUpkeep{
			ID:     ocr2keepers.UpkeepIdentifier(uid.Bytes()),
			Type:   tp,
			Config: struct{}{}, // empty struct by default
		},
		Trigger:   trigger,
		CheckData: checkData,
	}
	// set work id based on upkeep id and trigger
	wid, err := UpkeepWorkID(uid, trigger)
	if err != nil {
		return ocr2keepers.UpkeepPayload{}, fmt.Errorf("error while generating workID: %w", err)
	}
	p.WorkID = wid
	// manually convert trigger to triggerWrapper
	triggerW := triggerWrapper{
		BlockNum:  uint32(trigger.BlockNumber),
		BlockHash: common.HexToHash(trigger.BlockHash),
	}

	switch ocr2keepers.UpkeepType(tp) {
	case ocr2keepers.LogTrigger:
		trExt, ok := trigger.Extension.(LogTriggerExtension)
		if !ok {
			return ocr2keepers.UpkeepPayload{}, fmt.Errorf("unrecognized trigger extension data")
		}
		hex, parseErr := common.ParseHexOrString(trExt.TxHash)
		if parseErr != nil {
			return ocr2keepers.UpkeepPayload{}, fmt.Errorf("tx hash parse error: %w", err)
		}
		triggerW.TxHash = common.BytesToHash(hex[:])
		triggerW.LogIndex = uint32(trExt.LogIndex)
	default:
	}

	// get trigger in bytes
	triggerBytes, err := PackTrigger(uid, triggerW)
	if err != nil {
		return ocr2keepers.UpkeepPayload{}, fmt.Errorf("%w: failed to pack trigger", err)
	}

	// set trigger id based on upkeep id and trigger
	p.ID = UpkeepTriggerID(uid, triggerBytes)

	// end
	return p, nil
}
