package core

import (
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_convenience"
	autov2common "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_log_automation"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/streams_lookup_compatible_interface"
)

var ConvenienceABI = types.MustGetABI(automation_convenience.AutomationConvenienceABI)
var AutoV2CommonABI = types.MustGetABI(autov2common.IAutomationV21PlusCommonABI)
var StreamsCompatibleABI = types.MustGetABI(streams_lookup_compatible_interface.StreamsLookupCompatibleInterfaceABI)
var ILogAutomationABI = types.MustGetABI(i_log_automation.ILogAutomationABI)
