package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"
)

var (
	ErrInvalidUpkeepID = fmt.Errorf("invalid upkeepID")
)

// UpkeepWorkID returns the identifier using the given upkeepID and trigger extension(tx hash and log index).
func UpkeepWorkID(id *big.Int, trigger ocr2keepers.Trigger) (string, error) {
	var triggerExtBytes []byte
	if trigger.LogTriggerExtension != nil {
		triggerExtBytes = trigger.LogTriggerExtension.LogIdentifier()
	}
	uid := &ocr2keepers.UpkeepIdentifier{}
	ok := uid.FromBigInt(id)
	if !ok {
		return "", ErrInvalidUpkeepID
	}
	// TODO (auto-4314): Ensure it works with conditionals and add unit tests
	hash := crypto.Keccak256(append(uid[:], triggerExtBytes...))
	return hex.EncodeToString(hash[:]), nil
}

func WorkIDGenerator(u ocr2keepers.UpkeepIdentifier, trigger ocr2keepers.Trigger) string {
	// Error should not happen here since we pass in a valid upkeepID
	// TODO: Clean this up
	id, _ := UpkeepWorkID(u.BigInt(), trigger)
	return id
}

func NewUpkeepPayload(id *big.Int, trigger ocr2keepers.Trigger, checkData []byte) (ocr2keepers.UpkeepPayload, error) {
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
