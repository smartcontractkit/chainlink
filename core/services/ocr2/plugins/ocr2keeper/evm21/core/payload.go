package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

var (
	ErrInvalidUpkeepID = fmt.Errorf("invalid upkeepID")
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
	var triggerExtBytes []byte
	if trigger.LogTriggerExtension != nil {
		triggerExtBytes = trigger.LogTriggerExtension.LogIdentifier()
	}
	return UpkeepWorkIDFromTriggerBytes(id, triggerExtBytes)
}

func UpkeepWorkIDFromTriggerBytes(id *big.Int, triggerBytes []byte) (string, error) {
	uid := &ocr2keepers.UpkeepIdentifier{}
	ok := uid.FromBigInt(id)
	if !ok {
		return "", ErrInvalidUpkeepID
	}
	// TODO (auto-4314): Ensure it works with conditionals and add unit tests
	hash := crypto.Keccak256(append(uid[:], triggerBytes...))
	return hex.EncodeToString(hash[:]), nil
}

func NewUpkeepPayload(id *big.Int, tp int, trigger ocr2keepers.Trigger, checkData []byte) (ocr2keepers.UpkeepPayload, error) {
	uid := &ocr2keepers.UpkeepIdentifier{}
	ok := uid.FromBigInt(id)
	if !ok {
		return ocr2keepers.UpkeepPayload{}, ErrInvalidUpkeepID
	}
	p := ocr2keepers.UpkeepPayload{
		UpkeepID:  *uid,
		Trigger:   trigger,
		CheckData: checkData,
	}
	// set work id based on upkeep id and trigger
	wid, err := UpkeepWorkID(id, trigger)
	if err != nil {
		return ocr2keepers.UpkeepPayload{}, fmt.Errorf("error while generating workID: %w", err)
	}
	p.WorkID = wid

	return p, nil
}

func UpkeepTriggerIDFromPayload(p ocr2keepers.UpkeepPayload) (string, error) {
	trigger := p.Trigger
	// manually convert trigger to triggerWrapper
	triggerW := triggerWrapper{
		BlockNum:  uint32(trigger.BlockNumber),
		BlockHash: common.Hash(trigger.BlockHash),
	}

	if trigger.LogTriggerExtension != nil {
		triggerW.TxHash = common.Hash(trigger.LogTriggerExtension.TxHash)
		triggerW.LogIndex = trigger.LogTriggerExtension.Index
	}

	// get trigger in bytes
	uid := p.UpkeepID.BigInt()
	triggerBytes, err := PackTrigger(uid, triggerW)
	if err != nil {
		return "", fmt.Errorf("%w: failed to pack trigger", err)
	}

	return UpkeepTriggerID(uid, triggerBytes), nil
}
