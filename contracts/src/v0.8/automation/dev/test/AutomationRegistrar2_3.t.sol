// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {IAutomationRegistryMaster2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {AutomationRegistrar2_3} from "../v2_3/AutomationRegistrar2_3.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

// forge test --match-path src/v0.8/automation/dev/test/AutomationRegistrar2_3.t.sol

contract SetUp is BaseTest {
  IAutomationRegistryMaster2_3 internal registry;
  AutomationRegistrar2_3 internal registrar;

  function setUp() public override {
    super.setUp();
    registry = deployRegistry();
    AutomationRegistrar2_3.InitialTriggerConfig[]
      memory triggerConfigs = new AutomationRegistrar2_3.InitialTriggerConfig[](2);
    triggerConfigs[0] = AutomationRegistrar2_3.InitialTriggerConfig({
      triggerType: 0, // condition
      autoApproveType: AutomationRegistrar2_3.AutoApproveType.DISABLED,
      autoApproveMaxAllowed: 0
    });
    triggerConfigs[1] = AutomationRegistrar2_3.InitialTriggerConfig({
      triggerType: 1, // log
      autoApproveType: AutomationRegistrar2_3.AutoApproveType.DISABLED,
      autoApproveMaxAllowed: 0
    });
    IERC20[] memory billingTokens;
    uint256[] memory minRegistrationFees;
    registrar = new AutomationRegistrar2_3(
      address(linkToken),
      registry,
      triggerConfigs,
      billingTokens,
      minRegistrationFees
    );
  }
}

contract OnTokenTransfer is SetUp {}
