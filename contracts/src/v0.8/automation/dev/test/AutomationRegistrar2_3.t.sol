// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {IAutomationRegistryMaster2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {AutomationRegistrar2_3} from "../v2_3/AutomationRegistrar2_3.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../v2_3/AutomationRegistryBase2_3.sol";

// forge test --match-path src/v0.8/automation/dev/test/AutomationRegistrar2_3.t.sol

contract SetUp is BaseTest {
  IAutomationRegistryMaster2_3 internal registry;
  AutomationRegistrar2_3 internal registrar;

  function setUp() public override {
    super.setUp();
    (registry, registrar) = deployAndConfigureAll(AutoBase.PayoutMode.ON_CHAIN);
    vm.stopPrank(); // reset identity at the start of each test
  }
}

contract RegisterUpkeep is SetUp {
  function testLink_autoApproveOff_happy() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(linkToken))));
    linkToken.approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: IERC20(address(linkToken)),
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(linkToken.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);
  }

  function testUSDToken_autoApproveOff_happy() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(mockERC20));
    mockERC20.approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: mockERC20,
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(mockERC20.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);
  }

  function testLink_autoApproveOn_happy() external {
    registrar.setTriggerConfig(0, AutomationRegistrar2_3.AutoApproveType.ENABLED_ALL, 1000);

    vm.startPrank(UPKEEP_ADMIN);
    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(linkToken))));
    linkToken.approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: IERC20(address(linkToken)),
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(linkToken.balanceOf(address(registrar)), 0);
    assertEq(linkToken.balanceOf(address(registry)), amount);
    assertEq(registry.getNumUpkeeps(), 1);
  }

  function testUSDToken_autoApproveOn_happy() external {
    registrar.setTriggerConfig(0, AutomationRegistrar2_3.AutoApproveType.ENABLED_ALL, 1000);

    vm.startPrank(UPKEEP_ADMIN);
    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(mockERC20));
    mockERC20.approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: mockERC20,
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(mockERC20.balanceOf(address(registrar)), 0);
    assertEq(mockERC20.balanceOf(address(registry)), amount);
    assertEq(registry.getNumUpkeeps(), 1);
  }
}
