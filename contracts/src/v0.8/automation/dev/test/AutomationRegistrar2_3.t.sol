// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {IAutomationRegistryMaster2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {AutomationRegistrar2_3} from "../v2_3/AutomationRegistrar2_3.sol";

// forge test --match-path src/v0.8/automation/dev/test/AutomationRegistrar2_3.t.sol

contract SetUp is BaseTest {
  IAutomationRegistryMaster2_3 internal registry;
  AutomationRegistrar2_3 internal registrar;

  function setUp() public override {
    super.setUp();
    (registry, registrar) = deployAndConfigureAll();
  }
}

contract OnTokenTransfer is SetUp {}
