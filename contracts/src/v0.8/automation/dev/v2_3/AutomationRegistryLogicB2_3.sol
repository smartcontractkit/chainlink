// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationRegistryBase2_3} from "./AutomationRegistryBase2_3.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/structs/EnumerableSet.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import {AutomationRegistryLogicC2_3} from "./AutomationRegistryLogicC2_3.sol";
import {Chainable} from "../../Chainable.sol";

contract AutomationRegistryLogicB2_3 is AutomationRegistryBase2_3, Chainable {
  using Address for address;
  using EnumerableSet for EnumerableSet.UintSet;
  using EnumerableSet for EnumerableSet.AddressSet;

  /**
   * @param logicC the address of the third logic contract
   */
  constructor(
    AutomationRegistryLogicC2_3 logicC
  )
    AutomationRegistryBase2_3(
      logicC.getLinkAddress(),
      logicC.getLinkUSDFeedAddress(),
      logicC.getNativeUSDFeedAddress(),
      logicC.getFastGasFeedAddress(),
      logicC.getAutomationForwarderLogic(),
      logicC.getAllowedReadOnlyAddress()
    )
    Chainable(address(logicC))
  {}
}
