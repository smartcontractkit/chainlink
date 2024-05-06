// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {IAutomationRegistryMaster2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {AutomationRegistrar2_3} from "../v2_3/AutomationRegistrar2_3.sol";
import {IERC20Metadata as IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../v2_3/AutomationRegistryBase2_3.sol";
import {IWrappedNative} from "../interfaces/v2_3/IWrappedNative.sol";

// forge test --match-path src/v0.8/automation/dev/test/AutomationRegistrar2_3.t.sol

contract SetUp is BaseTest {
  IAutomationRegistryMaster2_3 internal registry;
  AutomationRegistrar2_3 internal registrar;

  function setUp() public override {
    super.setUp();
    vm.startPrank(OWNER);
    (registry, registrar) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);
    vm.stopPrank(); // reset identity at the start of each test
  }
}

contract CancelUpkeep is SetUp {
  function testUSDToken_happy() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(usdToken18));
    usdToken18.approve(address(registrar), amount);

    AutomationRegistrar2_3.RegistrationParams memory registrationParams = AutomationRegistrar2_3.RegistrationParams({
      upkeepContract: address(TARGET1),
      amount: amount,
      adminAddress: UPKEEP_ADMIN,
      gasLimit: 10_000,
      triggerType: 0,
      billingToken: usdToken18,
      name: "foobar",
      encryptedEmail: "",
      checkData: bytes("check data"),
      triggerConfig: "",
      offchainConfig: ""
    });

    // default is auto approve off
    registrar.registerUpkeep(registrationParams);

    assertEq(usdToken18.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);

    uint256 startRegistrarBalance = usdToken18.balanceOf(address(registrar));
    uint256 startUpkeepAdminBalance = usdToken18.balanceOf(UPKEEP_ADMIN);

    // cancel the upkeep
    vm.startPrank(OWNER);
    bytes32 hash = keccak256(abi.encode(registrationParams));
    registrar.cancel(hash);

    uint256 endRegistrarBalance = usdToken18.balanceOf(address(registrar));
    uint256 endUpkeepAdminBalance = usdToken18.balanceOf(UPKEEP_ADMIN);

    assertEq(startRegistrarBalance - amount, endRegistrarBalance);
    assertEq(startUpkeepAdminBalance + amount, endUpkeepAdminBalance);
  }
}

contract ApproveUpkeep is SetUp {
  function testUSDToken_happy() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(usdToken18));
    usdToken18.approve(address(registrar), amount);

    AutomationRegistrar2_3.RegistrationParams memory registrationParams = AutomationRegistrar2_3.RegistrationParams({
      upkeepContract: address(TARGET1),
      amount: amount,
      adminAddress: UPKEEP_ADMIN,
      gasLimit: 10_000,
      triggerType: 0,
      billingToken: usdToken18,
      name: "foobar",
      encryptedEmail: "",
      checkData: bytes("check data"),
      triggerConfig: "",
      offchainConfig: ""
    });

    // default is auto approve off
    registrar.registerUpkeep(registrationParams);

    assertEq(usdToken18.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);

    uint256 startRegistrarBalance = usdToken18.balanceOf(address(registrar));
    uint256 startRegistryBalance = usdToken18.balanceOf(address(registry));

    // approve the upkeep
    vm.startPrank(OWNER);
    registrar.approve(registrationParams);

    uint256 endRegistrarBalance = usdToken18.balanceOf(address(registrar));
    uint256 endRegistryBalance = usdToken18.balanceOf(address(registry));

    assertEq(startRegistrarBalance - amount, endRegistrarBalance);
    assertEq(startRegistryBalance + amount, endRegistryBalance);
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

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(usdToken18));
    usdToken18.approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: usdToken18,
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(usdToken18.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);
  }

  function testLink_autoApproveOn_happy() external {
    vm.startPrank(OWNER);
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
    vm.startPrank(OWNER);
    registrar.setTriggerConfig(0, AutomationRegistrar2_3.AutoApproveType.ENABLED_ALL, 1000);

    vm.startPrank(UPKEEP_ADMIN);
    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(usdToken18));
    usdToken18.approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: usdToken18,
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(usdToken18.balanceOf(address(registrar)), 0);
    assertEq(usdToken18.balanceOf(address(registry)), amount);
    assertEq(registry.getNumUpkeeps(), 1);
  }

  function testNative_autoApproveOn_happy() external {
    vm.startPrank(OWNER);
    registrar.setTriggerConfig(0, AutomationRegistrar2_3.AutoApproveType.ENABLED_ALL, 1000);

    vm.startPrank(UPKEEP_ADMIN);
    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(weth))));
    IWrappedNative(address(weth)).approve(address(registrar), amount);

    registrar.registerUpkeep{value: amount}(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: 0,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: IERC20(address(weth)),
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(weth.balanceOf(address(registrar)), 0);
    assertEq(weth.balanceOf(address(registry)), amount);
    assertEq(registry.getNumUpkeeps(), 1);
  }

  // when msg.value is 0, it uses the ERC20 payment path
  function testNative_autoApproveOff_msgValue0() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(weth))));
    IWrappedNative(address(weth)).approve(address(registrar), amount);

    registrar.registerUpkeep(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: amount,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: IERC20(address(weth)),
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(weth.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);
  }

  // when msg.value is not 0, it uses the native payment path
  function testNative_autoApproveOff_msgValueNot0() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(weth))));
    IWrappedNative(address(weth)).approve(address(registrar), amount);

    registrar.registerUpkeep{value: amount}(
      AutomationRegistrar2_3.RegistrationParams({
        upkeepContract: address(TARGET1),
        amount: 0,
        adminAddress: UPKEEP_ADMIN,
        gasLimit: 10_000,
        triggerType: 0,
        billingToken: IERC20(address(weth)),
        name: "foobar",
        encryptedEmail: "",
        checkData: bytes("check data"),
        triggerConfig: "",
        offchainConfig: ""
      })
    );

    assertEq(weth.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);
  }

  function testLink_autoApproveOff_revertOnDuplicateEntry() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(linkToken))));
    linkToken.approve(address(registrar), amount * 2);

    AutomationRegistrar2_3.RegistrationParams memory params = AutomationRegistrar2_3.RegistrationParams({
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
    });

    registrar.registerUpkeep(params);

    assertEq(linkToken.balanceOf(address(registrar)), amount);
    assertEq(registry.getNumUpkeeps(), 0);

    // attempt to register the same upkeep again
    vm.expectRevert(AutomationRegistrar2_3.DuplicateEntry.selector);
    registrar.registerUpkeep(params);
  }

  function test_revertOnInsufficientPayment() external {
    vm.startPrank(UPKEEP_ADMIN);

    // slightly less than the minimum amount
    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(linkToken))) - 1);
    linkToken.approve(address(registrar), amount);

    AutomationRegistrar2_3.RegistrationParams memory params = AutomationRegistrar2_3.RegistrationParams({
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
    });

    // attempt to register but revert bc of insufficient payment
    vm.expectRevert(AutomationRegistrar2_3.InsufficientPayment.selector);
    registrar.registerUpkeep(params);
  }

  function test_revertOnInvalidAdminAddress() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = uint96(registrar.getMinimumRegistrationAmount(IERC20(address(linkToken))));
    linkToken.approve(address(registrar), amount);

    AutomationRegistrar2_3.RegistrationParams memory params = AutomationRegistrar2_3.RegistrationParams({
      upkeepContract: address(TARGET1),
      amount: amount,
      adminAddress: ZERO_ADDRESS, // zero address is invalid
      gasLimit: 10_000,
      triggerType: 0,
      billingToken: IERC20(address(linkToken)),
      name: "foobar",
      encryptedEmail: "",
      checkData: bytes("check data"),
      triggerConfig: "",
      offchainConfig: ""
    });

    // attempt to register but revert bc of invalid admin address
    vm.expectRevert(AutomationRegistrar2_3.InvalidAdminAddress.selector);
    registrar.registerUpkeep(params);
  }

  function test_revertOnInvalidBillingToken() external {
    vm.startPrank(UPKEEP_ADMIN);

    uint96 amount = 1;
    usdToken18_2.approve(address(registrar), amount);

    AutomationRegistrar2_3.RegistrationParams memory params = AutomationRegistrar2_3.RegistrationParams({
      upkeepContract: address(TARGET1),
      amount: amount,
      adminAddress: UPKEEP_ADMIN,
      gasLimit: 10_000,
      triggerType: 0,
      billingToken: IERC20(address(usdToken18_2)), // unsupported billing token
      name: "foobar",
      encryptedEmail: "",
      checkData: bytes("check data"),
      triggerConfig: "",
      offchainConfig: ""
    });

    // attempt to register but revert bc of invalid admin address
    vm.expectRevert(AutomationRegistrar2_3.InvalidBillingToken.selector);
    registrar.registerUpkeep(params);
  }
}
