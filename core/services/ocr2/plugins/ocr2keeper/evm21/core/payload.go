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

func UpkeepWorkID(uid ocr2keepers.UpkeepIdentifier, trigger ocr2keepers.Trigger) string {
	var triggerExtBytes []byte
	if trigger.LogTriggerExtension != nil {
		triggerExtBytes = trigger.LogTriggerExtension.LogIdentifier()
	}
	hash := crypto.Keccak256(append(uid[:], triggerExtBytes...))
	return hex.EncodeToString(hash[:])
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
	p.WorkID = UpkeepWorkID(*uid, trigger)
	return p, nil
}
