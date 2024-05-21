// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Log} from "./interfaces/ILogAutomation.sol";
import {IAutomationV21PlusCommon} from "./interfaces/IAutomationV21PlusCommon.sol";

contract AutomationCompatibleUtils {
  function _report(IAutomationV21PlusCommon.Report memory) external {}

  function _logTriggerConfig(IAutomationV21PlusCommon.LogTriggerConfig memory) external {}

  function _logTrigger(IAutomationV21PlusCommon.LogTrigger memory) external {}

  function _conditionalTrigger(IAutomationV21PlusCommon.ConditionalTrigger memory) external {}

  function _log(Log memory) external {}
}
