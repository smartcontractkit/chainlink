package encoding

import (
	"math/big"

	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg/v3/types"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/mercury"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
)

const (
	// upkeep failure onchain reasons
	UpkeepFailureReasonNone                    uint8 = 0
	UpkeepFailureReasonUpkeepCancelled         uint8 = 1
	UpkeepFailureReasonUpkeepPaused            uint8 = 2
	UpkeepFailureReasonTargetCheckReverted     uint8 = 3
	UpkeepFailureReasonUpkeepNotNeeded         uint8 = 4
	UpkeepFailureReasonPerformDataExceedsLimit uint8 = 5
	UpkeepFailureReasonInsufficientBalance     uint8 = 6
	UpkeepFailureReasonMercuryCallbackReverted uint8 = 7
	UpkeepFailureReasonRevertDataExceedsLimit  uint8 = 8
	UpkeepFailureReasonRegistryPaused          uint8 = 9
	// leaving a gap here for more onchain failure reasons in the future
	// upkeep failure offchain reasons
	UpkeepFailureReasonMercuryAccessNotAllowed uint8 = 32
	UpkeepFailureReasonTxHashNoLongerExists    uint8 = 33
	UpkeepFailureReasonInvalidRevertDataInput  uint8 = 34
	UpkeepFailureReasonSimulationFailed        uint8 = 35
	UpkeepFailureReasonTxHashReorged           uint8 = 36

	// pipeline execution error
	NoPipelineError        uint8 = 0
	CheckBlockTooOld       uint8 = 1
	CheckBlockInvalid      uint8 = 2
	RpcFlakyFailure        uint8 = 3
	MercuryFlakyFailure    uint8 = 4
	PackUnpackDecodeFailed uint8 = 5
	MercuryUnmarshalError  uint8 = 6
	InvalidMercuryRequest  uint8 = 7
	InvalidMercuryResponse uint8 = 8 // this will only happen if Mercury server sends bad responses
	UpkeepNotAuthorized    uint8 = 9
)

type UpkeepInfo = iregistry21.KeeperRegistryBase21UpkeepInfo

type Packer interface {
	UnpackCheckResult(payload ocr2keepers.UpkeepPayload, raw string) (ocr2keepers.CheckResult, error)
	UnpackCheckCallbackResult(callbackResp []byte) (uint8, bool, []byte, uint8, *big.Int, error)
	UnpackPerformResult(raw string) (uint8, bool, error)
	UnpackLogTriggerConfig(raw []byte) (automation_utils_2_1.LogTriggerConfig, error)
	PackReport(report automation_utils_2_1.KeeperRegistryBase21Report) ([]byte, error)
	UnpackReport(raw []byte) (automation_utils_2_1.KeeperRegistryBase21Report, error)
	PackGetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error)
	UnpackGetUpkeepPrivilegeConfig(resp []byte) ([]byte, error)
	DecodeStreamsLookupRequest(data []byte) (*mercury.StreamsLookupError, error)
}
