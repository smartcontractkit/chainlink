// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IChainModule} from "./interfaces/v2_2/IChainModule.sol";
import {Log} from "../interfaces/ILogAutomation.sol";
import {IAutomationV2Common} from "./interfaces/IAutomationV2Common.sol";

/**
 * @notice OnchainConfig of the registry v2.1
 * @member paymentPremiumPPB payment premium rate oracles receive on top of
 * being reimbursed for gas, measured in parts per billion
 * @member flatFeeMicroLink flat fee paid to oracles for performing upkeeps,
 * priced in MicroLink; can be used in conjunction with or independently of
 * paymentPremiumPPB
 * @member checkGasLimit gas limit when checking for upkeep
 * @member stalenessSeconds number of seconds that is allowed for feed data to
 * be stale before switching to the fallback pricing
 * @member gasCeilingMultiplier multiplier to apply to the fast gas feed price
 * when calculating the payment ceiling for keepers
 * @member minUpkeepSpend minimum LINK that an upkeep must spend before cancelling
 * @member maxPerformGas max performGas allowed for an upkeep on this registry
 * @member maxCheckDataSize max length of checkData bytes
 * @member maxPerformDataSize max length of performData bytes
 * @member maxRevertDataSize max length of revertData bytes
 * @member fallbackGasPrice gas price used if the gas price feed is stale
 * @member fallbackLinkPrice LINK price used if the LINK price feed is stale
 * @member transcoder address of the transcoder contract
 * @member registrars addresses of the registrar contracts
 * @member upkeepPrivilegeManager address which can set privilege for upkeeps
 */
struct OnchainConfigV21 {
  uint32 paymentPremiumPPB;
  uint32 flatFeeMicroLink; // min 0.000001 LINK, max 4294 LINK
  uint32 checkGasLimit;
  uint24 stalenessSeconds;
  uint16 gasCeilingMultiplier;
  uint96 minUpkeepSpend;
  uint32 maxPerformGas;
  uint32 maxCheckDataSize;
  uint32 maxPerformDataSize;
  uint32 maxRevertDataSize;
  uint256 fallbackGasPrice;
  uint256 fallbackLinkPrice;
  address transcoder;
  address[] registrars;
  address upkeepPrivilegeManager;
}

/**
 * @notice OnchainConfig of the registry v2.2
 * @member paymentPremiumPPB payment premium rate oracles receive on top of
 * being reimbursed for gas, measured in parts per billion
 * @member flatFeeMicroLink flat fee paid to oracles for performing upkeeps,
 * priced in MicroLink; can be used in conjunction with or independently of
 * paymentPremiumPPB
 * @member checkGasLimit gas limit when checking for upkeep
 * @member stalenessSeconds number of seconds that is allowed for feed data to
 * be stale before switching to the fallback pricing
 * @member gasCeilingMultiplier multiplier to apply to the fast gas feed price
 * when calculating the payment ceiling for keepers
 * @member minUpkeepSpend minimum LINK that an upkeep must spend before cancelling
 * @member maxPerformGas max performGas allowed for an upkeep on this registry
 * @member maxCheckDataSize max length of checkData bytes
 * @member maxPerformDataSize max length of performData bytes
 * @member maxRevertDataSize max length of revertData bytes
 * @member fallbackGasPrice gas price used if the gas price feed is stale
 * @member fallbackLinkPrice LINK price used if the LINK price feed is stale
 * @member transcoder address of the transcoder contract
 * @member registrars addresses of the registrar contracts
 * @member upkeepPrivilegeManager address which can set privilege for upkeeps
 * @member reorgProtectionEnabled if this registry enables re-org protection checks
 * @member chainModule the chain specific module
 */
struct OnchainConfigV22 {
  uint32 paymentPremiumPPB;
  uint32 flatFeeMicroLink; // min 0.000001 LINK, max 4294 LINK
  uint32 checkGasLimit;
  uint24 stalenessSeconds;
  uint16 gasCeilingMultiplier;
  uint96 minUpkeepSpend;
  uint32 maxPerformGas;
  uint32 maxCheckDataSize;
  uint32 maxPerformDataSize;
  uint32 maxRevertDataSize;
  uint256 fallbackGasPrice;
  uint256 fallbackLinkPrice;
  address transcoder;
  address[] registrars;
  address upkeepPrivilegeManager;
  IChainModule chainModule;
  bool reorgProtectionEnabled;
}

contract AutomationConvenience {
  function _onChainConfig22(OnchainConfigV22 memory) external {}
  function _onChainConfig21(OnchainConfigV21 memory) external {}
  function _report(IAutomationV2Common.Report memory) external {}

  function _logTriggerConfig(IAutomationV2Common.LogTriggerConfig memory) external {}

  function _logTrigger(IAutomationV2Common.LogTrigger memory) external {}

  function _conditionalTrigger(IAutomationV2Common.ConditionalTrigger memory) external {}

  function _log(Log memory) external {}
}
