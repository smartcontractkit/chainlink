// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Vm} from "forge-std/Test.sol";
import {BaseTest} from "./BaseTest.t.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../v2_3/AutomationRegistryBase2_3.sol";
import {AutomationRegistrar2_3 as Registrar} from "../v2_3/AutomationRegistrar2_3.sol";
import {IAutomationRegistryMaster2_3 as Registry, AutomationRegistryBase2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {ChainModuleBase} from "../../chains/ChainModuleBase.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {IWrappedNative} from "../interfaces/v2_3/IWrappedNative.sol";

// forge test --match-path src/v0.8/automation/dev/test/AutomationRegistry2_3.t.sol

enum Trigger {
  CONDITION,
  LOG
}

contract SetUp is BaseTest {
  Registry internal registry;
  AutomationRegistryBase2_3.OnchainConfig internal config;
  bytes internal constant offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

  uint256 linkUpkeepID;
  uint256 usdUpkeepID;
  uint256 nativeUpkeepID;

  function setUp() public virtual override {
    super.setUp();

    (registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);
    config = registry.getConfig();

    vm.startPrank(OWNER);
    linkToken.approve(address(registry), type(uint256).max);
    usdToken.approve(address(registry), type(uint256).max);
    weth.approve(address(registry), type(uint256).max);
    vm.startPrank(UPKEEP_ADMIN);
    linkToken.approve(address(registry), type(uint256).max);
    usdToken.approve(address(registry), type(uint256).max);
    weth.approve(address(registry), type(uint256).max);
    vm.startPrank(STRANGER);
    linkToken.approve(address(registry), type(uint256).max);
    usdToken.approve(address(registry), type(uint256).max);
    weth.approve(address(registry), type(uint256).max);
    vm.stopPrank();

    linkUpkeepID = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );

    usdUpkeepID = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(usdToken),
      "",
      "",
      ""
    );

    nativeUpkeepID = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(weth),
      "",
      "",
      ""
    );

    vm.startPrank(OWNER);
    registry.addFunds(linkUpkeepID, registry.getMinBalanceForUpkeep(linkUpkeepID));
    registry.addFunds(usdUpkeepID, registry.getMinBalanceForUpkeep(usdUpkeepID));
    registry.addFunds(nativeUpkeepID, registry.getMinBalanceForUpkeep(nativeUpkeepID));
    vm.stopPrank();
  }
}

contract LatestConfigDetails is SetUp {
  function testGet() public {
    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = registry.latestConfigDetails();
    assertEq(configCount, 1);
    assertTrue(blockNumber > 0);
    assertNotEq(configDigest, "");
  }
}

contract CheckUpkeep is SetUp {
  function testPreventExecutionOnCheckUpkeep() public {
    uint256 id = 1;
    bytes memory triggerData = abi.encodePacked("trigger_data");

    // The tx.origin is the DEFAULT_SENDER (0x1804c8AB1F12E6bbf3894d4083f33e07309d1f38) of foundry
    // Expecting a revert since the tx.origin is not address(0)
    vm.expectRevert(abi.encodeWithSelector(Registry.OnlySimulatedBackend.selector));
    registry.checkUpkeep(id, triggerData);
  }
}

contract AddFunds is SetUp {
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);

  // when msg.value is 0, it uses the ERC20 payment path
  function testNative_msgValue0() external {
    vm.startPrank(OWNER);
    uint256 startRegistryBalance = registry.getBalance(nativeUpkeepID);
    uint256 startTokenBalance = registry.getBalance(nativeUpkeepID);
    registry.addFunds(nativeUpkeepID, 1);
    assertEq(registry.getBalance(nativeUpkeepID), startRegistryBalance + 1);
    assertEq(weth.balanceOf(address(registry)), startTokenBalance + 1);
  }

  // when msg.value is not 0, it uses the native payment path
  function testNative_msgValueNot0() external {
    uint256 startRegistryBalance = registry.getBalance(nativeUpkeepID);
    uint256 startTokenBalance = registry.getBalance(nativeUpkeepID);
    registry.addFunds{value: 1}(nativeUpkeepID, 1000); // parameter amount should be ignored
    assertEq(registry.getBalance(nativeUpkeepID), startRegistryBalance + 1);
    assertEq(weth.balanceOf(address(registry)), startTokenBalance + 1);
  }

  // it fails when the billing token is not native, but trying to pay with native
  function test_RevertsWhen_NativePaymentDoesntMatchBillingToken() external {
    vm.expectRevert(abi.encodeWithSelector(Registry.InvalidToken.selector));
    registry.addFunds{value: 1}(linkUpkeepID, 0);
  }

  function test_RevertsWhen_UpkeepDoesNotExist() public {
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.addFunds(randomNumber(), 1);
  }

  function test_RevertsWhen_UpkeepIsCanceled() public {
    registry.cancelUpkeep(linkUpkeepID);
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.addFunds(linkUpkeepID, 1);
  }

  function test_anyoneCanAddFunds() public {
    uint256 startAmount = registry.getBalance(linkUpkeepID);
    vm.prank(UPKEEP_ADMIN);
    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), startAmount + 1);
    vm.prank(STRANGER);
    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), startAmount + 2);
  }

  function test_movesFundFromCorrectToken() public {
    vm.startPrank(UPKEEP_ADMIN);

    uint256 startBalanceLINK = linkToken.balanceOf(address(registry));
    uint256 startBalanceUSDToken = usdToken.balanceOf(address(registry));
    uint256 startLinkUpkeepBalance = registry.getBalance(linkUpkeepID);
    uint256 startUSDUpkeepBalance = registry.getBalance(usdUpkeepID);

    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), startBalanceLINK + 1);
    assertEq(registry.getBalance(usdUpkeepID), startBalanceUSDToken);
    assertEq(linkToken.balanceOf(address(registry)), startLinkUpkeepBalance + 1);
    assertEq(usdToken.balanceOf(address(registry)), startUSDUpkeepBalance);

    registry.addFunds(usdUpkeepID, 2);
    assertEq(registry.getBalance(linkUpkeepID), startBalanceLINK + 1);
    assertEq(registry.getBalance(usdUpkeepID), startBalanceUSDToken + 2);
    assertEq(linkToken.balanceOf(address(registry)), startLinkUpkeepBalance + 1);
    assertEq(usdToken.balanceOf(address(registry)), startUSDUpkeepBalance + 2);
  }

  function test_emitsAnEvent() public {
    vm.startPrank(UPKEEP_ADMIN);
    vm.expectEmit();
    emit FundsAdded(linkUpkeepID, address(UPKEEP_ADMIN), 100);
    registry.addFunds(linkUpkeepID, 100);
  }
}

contract Withdraw is SetUp {
  address internal aMockAddress = randomAddress();

  function testLinkAvailableForPaymentReturnsLinkBalance() public {
    uint256 startBalance = linkToken.balanceOf(address(registry));
    int256 startLinkAvailable = registry.linkAvailableForPayment();

    //simulate a deposit of link to the liquidity pool
    _mintLink(address(registry), 1e10);

    //check there's a balance
    assertEq(linkToken.balanceOf(address(registry)), startBalance + 1e10);

    //check the link available has increased by the same amount
    assertEq(uint256(registry.linkAvailableForPayment()), uint256(startLinkAvailable) + 1e10);
  }

  function testWithdrawLinkRevertsBecauseOnlyFinanceAdminAllowed() public {
    vm.expectRevert(abi.encodeWithSelector(Registry.OnlyFinanceAdmin.selector));
    registry.withdrawLink(aMockAddress, 1);
  }

  function testWithdrawLinkRevertsBecauseOfInsufficientBalance() public {
    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(Registry.InsufficientBalance.selector, 0, 1));
    registry.withdrawLink(aMockAddress, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkRevertsBecauseOfInvalidRecipient() public {
    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(Registry.InvalidRecipient.selector));
    registry.withdrawLink(ZERO_ADDRESS, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkSuccess() public {
    //simulate a deposit of link to the liquidity pool
    _mintLink(address(registry), 1e10);
    uint256 startBalance = linkToken.balanceOf(address(registry));

    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is a ton of link available
    registry.withdrawLink(aMockAddress, 1);

    vm.stopPrank();

    assertEq(linkToken.balanceOf(address(aMockAddress)), 1);
    assertEq(linkToken.balanceOf(address(registry)), startBalance - 1);
  }

  function test_WithdrawERC20Fees_RespectsReserveAmount() public {
    assertEq(registry.getBalance(usdUpkeepID), registry.getReserveAmount(address(usdToken)));
    vm.startPrank(FINANCE_ADMIN);
    vm.expectRevert(abi.encodeWithSelector(Registry.InsufficientBalance.selector, 0, 1));
    registry.withdrawERC20Fees(address(usdToken), FINANCE_ADMIN, 1);
  }

  function test_WithdrawERC20Fees_RevertsWhen_AttemptingToWithdrawLINK() public {
    _mintLink(address(registry), 1e10);
    vm.startPrank(FINANCE_ADMIN);
    vm.expectRevert(Registry.InvalidToken.selector);
    registry.withdrawERC20Fees(address(linkToken), FINANCE_ADMIN, 1); // should revert
    registry.withdrawLink(FINANCE_ADMIN, 1); // but using link withdraw functions succeeds
  }

  function test_WithdrawERC20Fees_RevertsWhen_LinkAvailableForPaymentIsNegative() public {
    _transmit(usdUpkeepID, registry); // adds USD token to finance withdrawable, and gives NOPs a LINK balance
    require(registry.linkAvailableForPayment() < 0, "linkAvailableForPayment should be negative");
    vm.expectRevert(Registry.InsufficientLinkLiquidity.selector);
    vm.prank(FINANCE_ADMIN);
    registry.withdrawERC20Fees(address(usdToken), FINANCE_ADMIN, 1); // should revert
    _mintLink(address(registry), uint256(registry.linkAvailableForPayment() * -10)); // top up LINK liquidity pool
    vm.prank(FINANCE_ADMIN);
    registry.withdrawERC20Fees(address(usdToken), FINANCE_ADMIN, 1); // now finance can withdraw
  }

  function testWithdrawERC20FeeSuccess() public {
    // deposit excess USDToken to the registry (this goes to the "finance withdrawable" pool be default)
    uint256 startReserveAmount = registry.getReserveAmount(address(usdToken));
    uint256 startAmount = usdToken.balanceOf(address(registry));
    _mintERC20(address(registry), 1e10);

    // depositing shouldn't change reserve amount
    assertEq(registry.getReserveAmount(address(usdToken)), startReserveAmount);

    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 USDToken
    registry.withdrawERC20Fees(address(usdToken), aMockAddress, 1);

    vm.stopPrank();

    assertEq(usdToken.balanceOf(address(aMockAddress)), 1);
    assertEq(usdToken.balanceOf(address(registry)), startAmount + 1e10 - 1);
    assertEq(registry.getReserveAmount(address(usdToken)), startReserveAmount);
  }
}

contract SetConfig is SetUp {
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] signers,
    address[] transmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );

  address module = address(new ChainModuleBase());
  AutomationRegistryBase2_3.OnchainConfig cfg =
    AutomationRegistryBase2_3.OnchainConfig({
      checkGasLimit: 5_000_000,
      stalenessSeconds: 90_000,
      gasCeilingMultiplier: 0,
      maxPerformGas: 10_000_000,
      maxCheckDataSize: 5_000,
      maxPerformDataSize: 5_000,
      maxRevertDataSize: 5_000,
      fallbackGasPrice: 20_000_000_000,
      fallbackLinkPrice: 2_000_000_000, // $20
      fallbackNativePrice: 400_000_000_000, // $4,000
      transcoder: 0xB1e66855FD67f6e85F0f0fA38cd6fBABdf00923c,
      registrars: new address[](0),
      upkeepPrivilegeManager: PRIVILEGE_MANAGER,
      chainModule: module,
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });

  function testSetConfigSuccess() public {
    (uint32 configCount, uint32 blockNumber, ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    address billingTokenAddress = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = billingTokenAddress;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000
    });

    bytes memory onchainConfigBytes = abi.encode(cfg);
    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    bytes32 configDigest = _configDigestFromConfigData(
      block.chainid,
      address(registry),
      ++configCount,
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    vm.expectEmit();
    emit ConfigSet(
      blockNumber,
      configDigest,
      configCount,
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registry.getState();

    assertEq(signers, SIGNERS);
    assertEq(transmitters, TRANSMITTERS);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config = registry.getBillingTokenConfig(billingTokenAddress);
    assertEq(config.gasFeePPB, 5_000);
    assertEq(config.flatFeeMilliCents, 20_000);
    assertEq(config.priceFeed, 0x2222222222222222222222222222222222222222);
    assertEq(config.minSpend, 100_000);

    address[] memory tokens = registry.getBillingTokens();
    assertEq(tokens.length, 1);
  }

  function testSetConfigMultipleBillingConfigsSuccess() public {
    (uint32 configCount, , ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    address billingTokenAddress1 = address(0x1111111111111111111111111111111111111111);
    address billingTokenAddress2 = address(0x1111111111111111111111111111111111111112);
    address[] memory billingTokens = new address[](2);
    billingTokens[0] = billingTokenAddress1;
    billingTokens[1] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMilliCents: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221,
      fallbackPrice: 100,
      minSpend: 100
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMilliCents: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 200,
      minSpend: 200
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registry.getState();

    assertEq(signers, SIGNERS);
    assertEq(transmitters, TRANSMITTERS);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config1 = registry.getBillingTokenConfig(billingTokenAddress1);
    assertEq(config1.gasFeePPB, 5_001);
    assertEq(config1.flatFeeMilliCents, 20_001);
    assertEq(config1.priceFeed, 0x2222222222222222222222222222222222222221);
    assertEq(config1.fallbackPrice, 100);
    assertEq(config1.minSpend, 100);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registry.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMilliCents, 20_002);
    assertEq(config2.priceFeed, 0x2222222222222222222222222222222222222222);
    assertEq(config2.fallbackPrice, 200);
    assertEq(config2.minSpend, 200);

    address[] memory tokens = registry.getBillingTokens();
    assertEq(tokens.length, 2);
  }

  function testSetConfigTwiceAndLastSetOverwrites() public {
    (uint32 configCount, , ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    // BillingConfig1
    address billingTokenAddress1 = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens1 = new address[](1);
    billingTokens1[0] = billingTokenAddress1;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs1 = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs1[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMilliCents: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221,
      fallbackPrice: 100,
      minSpend: 100
    });

    bytes memory onchainConfigBytesWithBilling1 = abi.encode(cfg, billingTokens1, billingConfigs1);

    // BillingConfig2
    address billingTokenAddress2 = address(0x1111111111111111111111111111111111111112);
    address[] memory billingTokens2 = new address[](1);
    billingTokens2[0] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs2 = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs2[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMilliCents: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 200,
      minSpend: 200
    });

    bytes memory onchainConfigBytesWithBilling2 = abi.encode(cfg, billingTokens2, billingConfigs2);

    // set config once
    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling1,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    // set config twice
    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling2,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registry.getState();

    assertEq(signers, SIGNERS);
    assertEq(transmitters, TRANSMITTERS);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registry.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMilliCents, 20_002);
    assertEq(config2.priceFeed, 0x2222222222222222222222222222222222222222);
    assertEq(config2.fallbackPrice, 200);
    assertEq(config2.minSpend, 200);

    address[] memory tokens = registry.getBillingTokens();
    assertEq(tokens.length, 1);
  }

  function testSetConfigDuplicateBillingConfigFailure() public {
    (uint32 configCount, , ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    address billingTokenAddress1 = address(0x1111111111111111111111111111111111111111);
    address billingTokenAddress2 = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens = new address[](2);
    billingTokens[0] = billingTokenAddress1;
    billingTokens[1] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMilliCents: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221,
      fallbackPrice: 100,
      minSpend: 100
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMilliCents: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 200,
      minSpend: 200
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    // expect revert because of duplicate tokens
    vm.expectRevert(abi.encodeWithSelector(Registry.DuplicateEntry.selector));
    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );
  }

  function testSetConfigRevertDueToInvalidToken() public {
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = address(linkToken);

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000
    });

    // deploy registry with OFF_CHAIN payout mode
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);

    vm.expectRevert(abi.encodeWithSelector(Registry.InvalidToken.selector));
    registry.setConfigTypeSafe(
      SIGNERS,
      TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes,
      billingTokens,
      billingConfigs
    );
  }

  function testSetConfigWithNewTransmittersSuccess() public {
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);

    (uint32 configCount, uint32 blockNumber, ) = registry.latestConfigDetails();
    assertEq(configCount, 0);

    address billingTokenAddress = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = billingTokenAddress;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000
    });

    bytes memory onchainConfigBytes = abi.encode(cfg);

    bytes32 configDigest = _configDigestFromConfigData(
      block.chainid,
      address(registry),
      ++configCount,
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    vm.expectEmit();
    emit ConfigSet(
      blockNumber,
      configDigest,
      configCount,
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    registry.setConfigTypeSafe(
      SIGNERS,
      TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes,
      billingTokens,
      billingConfigs
    );

    (, , address[] memory signers, address[] memory transmitters, ) = registry.getState();
    assertEq(signers, SIGNERS);
    assertEq(transmitters, TRANSMITTERS);

    (configCount, blockNumber, ) = registry.latestConfigDetails();
    configDigest = _configDigestFromConfigData(
      block.chainid,
      address(registry),
      ++configCount,
      SIGNERS,
      NEW_TRANSMITTERS,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    vm.expectEmit();
    emit ConfigSet(
      blockNumber,
      configDigest,
      configCount,
      SIGNERS,
      NEW_TRANSMITTERS,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    registry.setConfigTypeSafe(
      SIGNERS,
      NEW_TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes,
      billingTokens,
      billingConfigs
    );

    (, , signers, transmitters, ) = registry.getState();
    assertEq(signers, SIGNERS);
    assertEq(transmitters, NEW_TRANSMITTERS);
  }

  function _configDigestFromConfigData(
    uint256 chainId,
    address contractAddress,
    uint64 configCount,
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal pure returns (bytes32) {
    uint256 h = uint256(
      keccak256(
        abi.encode(
          chainId,
          contractAddress,
          configCount,
          signers,
          transmitters,
          f,
          onchainConfig,
          offchainConfigVersion,
          offchainConfig
        )
      )
    );
    uint256 prefixMask = type(uint256).max << (256 - 16); // 0xFFFF00..00
    uint256 prefix = 0x0001 << (256 - 16); // 0x000100..00
    return bytes32((prefix & prefixMask) | (h & ~prefixMask));
  }
}

contract NOPsSettlement is SetUp {
  event NOPsSettledOffchain(address[] payees, uint256[] payments);
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);
  event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee);

  function testSettleNOPsOffchainRevertDueToUnauthorizedCaller() public {
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);

    vm.expectRevert(abi.encodeWithSelector(Registry.OnlyFinanceAdmin.selector));
    registry.settleNOPsOffchain();
  }

  function testSettleNOPsOffchainRevertDueToOffchainSettlementDisabled() public {
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    vm.prank(registry.owner());
    registry.disableOffchainPayments();

    vm.prank(FINANCE_ADMIN);
    vm.expectRevert(abi.encodeWithSelector(Registry.MustSettleOnchain.selector));
    registry.settleNOPsOffchain();
  }

  function testSettleNOPsOffchainSuccess() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    uint256[] memory payments = new uint256[](TRANSMITTERS.length);
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      payments[i] = 0;
    }

    vm.startPrank(FINANCE_ADMIN);
    vm.expectEmit();
    emit NOPsSettledOffchain(PAYEES, payments);
    registry.settleNOPsOffchain();
  }

  function testSettleNOPsOffchainSuccessTransmitterBalanceZeroed() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken), "", "", "");
    _mintERC20(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken.approve(address(registry), 1e20);
    registry.addFunds(id, 1e20);

    // manually create a transmit so transmitters earn some rewards
    _transmit(id, registry);

    // verify transmitters have positive balances
    uint256[] memory payments = new uint256[](TRANSMITTERS.length);
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      (bool active, uint8 index, uint96 balance, uint96 lastCollected, ) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      assertTrue(active);
      assertEq(i, index);
      assertTrue(balance > 0);
      assertEq(0, lastCollected);

      payments[i] = balance;
    }

    // verify offchain settlement will emit NOPs' balances
    vm.startPrank(FINANCE_ADMIN);
    vm.expectEmit();
    emit NOPsSettledOffchain(PAYEES, payments);
    registry.settleNOPsOffchain();

    // verify that transmitters balance has been zeroed out
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      (bool active, uint8 index, uint96 balance, , ) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      assertTrue(active);
      assertEq(i, index);
      assertEq(0, balance);
    }
  }

  function testSettleNOPsOffchainForDeactivatedTransmittersSuccess() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, Registrar registrar) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken), "", "", "");
    _mintERC20(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken.approve(address(registry), 1e20);
    registry.addFunds(id, 1e20);

    // manually create a transmit so TRANSMITTERS earn some rewards
    _transmit(id, registry);

    // TRANSMITTERS have positive balance now
    // configure the registry to use NEW_TRANSMITTERS
    _configureWithNewTransmitters(registry, registrar);

    _transmit(id, registry);

    // verify all transmitters have positive balances
    address[] memory expectedPayees = new address[](6);
    uint256[] memory expectedPayments = new uint256[](6);
    for (uint256 i = 0; i < NEW_TRANSMITTERS.length; i++) {
      (bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee) = registry.getTransmitterInfo(
        NEW_TRANSMITTERS[i]
      );
      assertTrue(active);
      assertEq(i, index);
      assertTrue(lastCollected > 0);
      expectedPayments[i] = balance;
      expectedPayees[i] = payee;
    }
    for (uint256 i = 2; i < TRANSMITTERS.length; i++) {
      (bool active, uint8 index, uint96 balance, uint96 lastCollected, address payee) = registry.getTransmitterInfo(
        TRANSMITTERS[i]
      );
      assertFalse(active);
      assertEq(i, index);
      assertTrue(balance > 0);
      assertTrue(lastCollected > 0);
      expectedPayments[2 + i] = balance;
      expectedPayees[2 + i] = payee;
    }

    // verify offchain settlement will emit NOPs' balances
    vm.startPrank(FINANCE_ADMIN);

    // simply expectEmit won't work here because s_deactivatedTransmitters is an enumerable set so the order of these
    // deactivated transmitters is not guaranteed. To handle this, we record logs and decode data field manually.
    vm.recordLogs();
    registry.settleNOPsOffchain();
    Vm.Log[] memory entries = vm.getRecordedLogs();

    assertEq(entries.length, 1);
    Vm.Log memory l = entries[0];
    assertEq(l.topics[0], keccak256("NOPsSettledOffchain(address[],uint256[])"));
    (address[] memory actualPayees, uint256[] memory actualPayments) = abi.decode(l.data, (address[], uint256[]));
    assertEq(actualPayees.length, 6);
    assertEq(actualPayments.length, 6);

    // first 4 payees and payments are for NEW_TRANSMITTERS and they are ordered.
    for (uint256 i = 0; i < NEW_TRANSMITTERS.length; i++) {
      assertEq(actualPayees[i], expectedPayees[i]);
      assertEq(actualPayments[i], expectedPayments[i]);
    }

    // the last 2 payees and payments for TRANSMITTERS[2] and TRANSMITTERS[3] and they are not ordered
    assertTrue(
      (actualPayments[5] == expectedPayments[5] &&
        actualPayees[5] == expectedPayees[5] &&
        actualPayments[4] == expectedPayments[4] &&
        actualPayees[4] == expectedPayees[4]) ||
        (actualPayments[5] == expectedPayments[4] &&
          actualPayees[5] == expectedPayees[4] &&
          actualPayments[4] == expectedPayments[5] &&
          actualPayees[4] == expectedPayees[5])
    );

    // verify that new transmitters balance has been zeroed out
    for (uint256 i = 0; i < NEW_TRANSMITTERS.length; i++) {
      (bool active, uint8 index, uint96 balance, , ) = registry.getTransmitterInfo(NEW_TRANSMITTERS[i]);
      assertTrue(active);
      assertEq(i, index);
      assertEq(0, balance);
    }
    // verify that deactivated transmitters (TRANSMITTERS[2] and TRANSMITTERS[3]) balance has been zeroed out
    for (uint256 i = 2; i < TRANSMITTERS.length; i++) {
      (bool active, uint8 index, uint96 balance, , ) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      assertFalse(active);
      assertEq(i, index);
      assertEq(0, balance);
    }
  }

  function testDisableOffchainPaymentsRevertDueToUnauthorizedCaller() public {
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    vm.startPrank(FINANCE_ADMIN);
    vm.expectRevert(bytes("Only callable by owner"));
    registry.disableOffchainPayments();
  }

  function testDisableOffchainPaymentsSuccess() public {
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    vm.startPrank(registry.owner());
    registry.disableOffchainPayments();

    assertEq(uint8(AutoBase.PayoutMode.ON_CHAIN), registry.getPayoutMode());
  }

  function testSinglePerformAndNodesCanWithdrawOnchain() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken), "", "", "");
    _mintERC20(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken.approve(address(registry), 1e20);
    registry.addFunds(id, 1e20);

    // manually create a transmit so transmitters earn some rewards
    _transmit(id, registry);

    // disable offchain payments
    _mintLink(address(registry), 1e19);
    vm.prank(registry.owner());
    registry.disableOffchainPayments();

    // payees should be able to withdraw onchain
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      (, , uint96 balance, , address payee) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      vm.prank(payee);
      vm.expectEmit();
      emit PaymentWithdrawn(TRANSMITTERS[i], balance, payee, payee);
      registry.withdrawPayment(TRANSMITTERS[i], payee);
    }

    // allow upkeep admin to withdraw funds
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(id);
    vm.roll(100 + block.number);
    vm.expectEmit();
    // the upkeep spent less than minimum spending limit so upkeep admin can only withdraw upkeep balance - min spend value
    emit FundsWithdrawn(id, 9.9e19, UPKEEP_ADMIN);
    registry.withdrawFunds(id, UPKEEP_ADMIN);
  }

  function testMultiplePerformsAndNodesCanWithdrawOnchain() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken), "", "", "");
    _mintERC20(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken.approve(address(registry), 1e20);
    registry.addFunds(id, 1e20);

    // manually call transmit so transmitters earn some rewards
    for (uint256 i = 0; i < 50; i++) {
      vm.roll(100 + block.number);
      _transmit(id, registry);
    }

    // disable offchain payments
    _mintLink(address(registry), 1e19);
    vm.prank(registry.owner());
    registry.disableOffchainPayments();

    // manually call transmit after offchain payment is disabled
    for (uint256 i = 0; i < 50; i++) {
      vm.roll(100 + block.number);
      _transmit(id, registry);
    }

    // payees should be able to withdraw onchain
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      (, , uint96 balance, , address payee) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      vm.prank(payee);
      vm.expectEmit();
      emit PaymentWithdrawn(TRANSMITTERS[i], balance, payee, payee);
      registry.withdrawPayment(TRANSMITTERS[i], payee);
    }

    // allow upkeep admin to withdraw funds
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(id);
    vm.roll(100 + block.number);
    uint256 balance = registry.getBalance(id);
    vm.expectEmit();
    emit FundsWithdrawn(id, balance, UPKEEP_ADMIN);
    registry.withdrawFunds(id, UPKEEP_ADMIN);
  }

  function _configureWithNewTransmitters(Registry registry, Registrar registrar) internal {
    IERC20[] memory billingTokens = new IERC20[](1);
    billingTokens[0] = IERC20(address(usdToken));
    uint256[] memory minRegistrationFees = new uint256[](billingTokens.length);
    minRegistrationFees[0] = 100000000000000000000; // 100 USD
    address[] memory billingTokenAddresses = new address[](billingTokens.length);
    for (uint256 i = 0; i < billingTokens.length; i++) {
      billingTokenAddresses[i] = address(billingTokens[i]);
    }
    AutomationRegistryBase2_3.BillingConfig[]
      memory billingTokenConfigs = new AutomationRegistryBase2_3.BillingConfig[](billingTokens.length);
    billingTokenConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 10_000_000, // 15%
      flatFeeMilliCents: 2_000, // 2 cents
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 100_000_000, // $1
      minSpend: 1000000000000000000 // 1 USD
    });

    address[] memory registrars;
    registrars = new address[](1);
    registrars[0] = address(registrar);
    AutomationRegistryBase2_3.OnchainConfig memory cfg = AutomationRegistryBase2_3.OnchainConfig({
      checkGasLimit: 5_000_000,
      stalenessSeconds: 90_000,
      gasCeilingMultiplier: 2,
      maxPerformGas: 10_000_000,
      maxCheckDataSize: 5_000,
      maxPerformDataSize: 5_000,
      maxRevertDataSize: 5_000,
      fallbackGasPrice: 20_000_000_000,
      fallbackLinkPrice: 2_000_000_000, // $20
      fallbackNativePrice: 400_000_000_000, // $4,000
      transcoder: 0xB1e66855FD67f6e85F0f0fA38cd6fBABdf00923c,
      registrars: registrars,
      upkeepPrivilegeManager: PRIVILEGE_MANAGER,
      chainModule: address(new ChainModuleBase()),
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });
    registry.setConfigTypeSafe(
      SIGNERS,
      NEW_TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      "",
      billingTokenAddresses,
      billingTokenConfigs
    );
    registry.setPayees(NEW_PAYEES);
  }
}

contract WithdrawPayment is SetUp {
  function testWithdrawPaymentRevertDueToOffchainPayoutMode() public {
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);
    vm.expectRevert(abi.encodeWithSelector(Registry.MustSettleOffchain.selector));
    vm.prank(TRANSMITTERS[0]);
    registry.withdrawPayment(TRANSMITTERS[0], TRANSMITTERS[0]);
  }
}

contract RegisterUpkeep is SetUp {
  function test_RevertsWhen_Paused() public {
    registry.pause();
    vm.expectRevert(Registry.RegistryPaused.selector);
    registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );
  }

  function test_RevertsWhen_TargetIsNotAContract() public {
    vm.expectRevert(Registry.NotAContract.selector);
    registry.registerUpkeep(
      randomAddress(),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );
  }

  function test_RevertsWhen_CalledByNonOwner() public {
    vm.prank(STRANGER);
    vm.expectRevert(Registry.OnlyCallableByOwnerOrRegistrar.selector);
    registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );
  }

  function test_RevertsWhen_ExecuteGasIsTooLow() public {
    vm.expectRevert(Registry.GasLimitOutsideRange.selector);
    registry.registerUpkeep(
      address(TARGET1),
      2299, // 1 less than min
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );
  }

  function test_RevertsWhen_ExecuteGasIsTooHigh() public {
    vm.expectRevert(Registry.GasLimitOutsideRange.selector);
    registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas + 1,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );
  }

  function test_RevertsWhen_TheBillingTokenIsNotConfigured() public {
    vm.expectRevert(Registry.InvalidToken.selector);
    registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      randomAddress(),
      "",
      "",
      ""
    );
  }

  function test_RevertsWhen_CheckDataIsTooLarge() public {
    vm.expectRevert(Registry.CheckDataExceedsLimit.selector);
    registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      randomBytes(config.maxCheckDataSize + 1),
      "",
      ""
    );
  }

  function test_Happy() public {
    bytes memory checkData = randomBytes(config.maxCheckDataSize);
    bytes memory trigggerConfig = randomBytes(100);
    bytes memory offchainConfig = randomBytes(100);

    uint256 upkeepCount = registry.getNumUpkeeps();

    uint256 upkeepID = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.LOG),
      address(linkToken),
      checkData,
      trigggerConfig,
      offchainConfig
    );

    assertEq(registry.getNumUpkeeps(), upkeepCount + 1);
    assertEq(registry.getUpkeep(upkeepID).target, address(TARGET1));
    assertEq(registry.getUpkeep(upkeepID).performGas, config.maxPerformGas);
    assertEq(registry.getUpkeep(upkeepID).checkData, checkData);
    assertEq(registry.getUpkeep(upkeepID).balance, 0);
    assertEq(registry.getUpkeep(upkeepID).admin, UPKEEP_ADMIN);
    assertEq(registry.getUpkeep(upkeepID).offchainConfig, offchainConfig);
    assertEq(registry.getUpkeepTriggerConfig(upkeepID), trigggerConfig);
    assertEq(uint8(registry.getTriggerType(upkeepID)), uint8(Trigger.LOG));
  }
}

contract OnTokenTransfer is SetUp {
  function test_RevertsWhen_NotCalledByTheLinkToken() public {
    vm.expectRevert(Registry.OnlyCallableByLINKToken.selector);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, abi.encode(linkUpkeepID));
  }

  function test_RevertsWhen_NotCalledWithExactly32Bytes() public {
    vm.startPrank(address(linkToken));
    vm.expectRevert(Registry.InvalidDataLength.selector);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, randomBytes(31));
    vm.expectRevert(Registry.InvalidDataLength.selector);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, randomBytes(33));
  }

  function test_RevertsWhen_TheUpkeepIsCancelledOrDNE() public {
    vm.startPrank(address(linkToken));
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, abi.encode(randomNumber()));
  }

  function test_RevertsWhen_TheUpkeepDoesNotUseLINKAsItsBillingToken() public {
    vm.startPrank(address(linkToken));
    vm.expectRevert(Registry.InvalidToken.selector);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, abi.encode(usdUpkeepID));
  }

  function test_Happy() public {
    vm.startPrank(address(linkToken));
    uint256 beforeBalance = registry.getBalance(linkUpkeepID);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, abi.encode(linkUpkeepID));
    assertEq(registry.getBalance(linkUpkeepID), beforeBalance + 100);
  }
}

contract GetMinBalanceForUpkeep is SetUp {
  function test_accountsForFlatFee() public {
    // set fee to 0
    AutomationRegistryBase2_3.BillingConfig memory usdTokenConfig = registry.getBillingTokenConfig(address(usdToken));
    usdTokenConfig.flatFeeMilliCents = 0;
    _updateBillingTokenConfig(registry, address(usdToken), usdTokenConfig);

    uint256 minBalanceBefore = registry.getMinBalanceForUpkeep(usdUpkeepID);

    // set fee to non-zero
    usdTokenConfig.flatFeeMilliCents = 100;
    _updateBillingTokenConfig(registry, address(usdToken), usdTokenConfig);

    uint256 minBalanceAfter = registry.getMinBalanceForUpkeep(usdUpkeepID);
    assertEq(minBalanceAfter, minBalanceBefore + (uint256(usdTokenConfig.flatFeeMilliCents) * 1e13));
  }
}

contract BillingOverrides is SetUp {
  event BillingConfigOverridden(uint256 indexed id, AutomationRegistryBase2_3.BillingOverrides overrides);
  event BillingConfigOverrideRemoved(uint256 indexed id);

  function test_RevertsWhen_NotPrivilegeManager() public {
    AutomationRegistryBase2_3.BillingOverrides memory billingOverrides = AutomationRegistryBase2_3.BillingOverrides({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000
    });

    vm.expectRevert(Registry.OnlyCallableByUpkeepPrivilegeManager.selector);
    registry.setBillingOverrides(linkUpkeepID, billingOverrides);
  }

  function test_RevertsWhen_UpkeepCancelled() public {
    AutomationRegistryBase2_3.BillingOverrides memory billingOverrides = AutomationRegistryBase2_3.BillingOverrides({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000
    });

    registry.cancelUpkeep(linkUpkeepID);

    vm.startPrank(PRIVILEGE_MANAGER);
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.setBillingOverrides(linkUpkeepID, billingOverrides);
  }

  function test_Happy_SetBillingOverrides() public {
    AutomationRegistryBase2_3.BillingOverrides memory billingOverrides = AutomationRegistryBase2_3.BillingOverrides({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000
    });

    vm.startPrank(PRIVILEGE_MANAGER);

    vm.expectEmit();
    emit BillingConfigOverridden(linkUpkeepID, billingOverrides);
    registry.setBillingOverrides(linkUpkeepID, billingOverrides);
  }

  function test_Happy_RemoveBillingOverrides() public {
    vm.startPrank(PRIVILEGE_MANAGER);

    vm.expectEmit();
    emit BillingConfigOverrideRemoved(linkUpkeepID);
    registry.removeBillingOverrides(linkUpkeepID);
  }

  function test_Happy_MaxGasPayment_WithBillingOverrides() public {
    uint96 maxPayment1 = registry.getMaxPaymentForGas(linkUpkeepID, 0, 5_000_000, address(linkToken));

    // Double the two billing values
    AutomationRegistryBase2_3.BillingOverrides memory billingOverrides = AutomationRegistryBase2_3.BillingOverrides({
      gasFeePPB: DEFAULT_GAS_FEE_PPB * 2,
      flatFeeMilliCents: DEFAULT_FLAT_FEE_MILLI_CENTS * 2
    });

    vm.startPrank(PRIVILEGE_MANAGER);
    registry.setBillingOverrides(linkUpkeepID, billingOverrides);

    // maxPayment2 should be greater than maxPayment1 after the overrides
    // The 2 numbers should follow this: maxPayment2 - maxPayment1 == 2 * recepit.premium
    // We do not apply the exact equation since we couldn't get the receipt.premium value
    uint96 maxPayment2 = registry.getMaxPaymentForGas(linkUpkeepID, 0, 5_000_000, address(linkToken));
    assertGt(maxPayment2, maxPayment1);
  }
}

contract Transmit is SetUp {
  function test_handlesMixedBatchOfBillingTokens() external {
    uint256[] memory prevUpkeepBalances = new uint256[](3);
    prevUpkeepBalances[0] = registry.getBalance(linkUpkeepID);
    prevUpkeepBalances[1] = registry.getBalance(usdUpkeepID);
    prevUpkeepBalances[2] = registry.getBalance(nativeUpkeepID);
    uint256[] memory prevTokenBalances = new uint256[](3);
    prevTokenBalances[0] = linkToken.balanceOf(address(registry));
    prevTokenBalances[1] = usdToken.balanceOf(address(registry));
    prevTokenBalances[2] = weth.balanceOf(address(registry));
    uint256[] memory prevReserveBalances = new uint256[](3);
    prevReserveBalances[0] = registry.getReserveAmount(address(linkToken));
    prevReserveBalances[1] = registry.getReserveAmount(address(usdToken));
    prevReserveBalances[2] = registry.getReserveAmount(address(weth));
    uint256[] memory upkeepIDs = new uint256[](3);
    upkeepIDs[0] = linkUpkeepID;
    upkeepIDs[1] = usdUpkeepID;
    upkeepIDs[2] = nativeUpkeepID;
    // do the thing
    _transmit(upkeepIDs, registry);
    // assert upkeep balances have decreased
    require(prevUpkeepBalances[0] > registry.getBalance(linkUpkeepID), "link upkeep balance should have decreased");
    require(prevUpkeepBalances[1] > registry.getBalance(usdUpkeepID), "usd upkeep balance should have decreased");
    require(prevUpkeepBalances[2] > registry.getBalance(nativeUpkeepID), "native upkeep balance should have decreased");
    // assert token balances have not changed
    assertEq(prevTokenBalances[0], linkToken.balanceOf(address(registry)));
    assertEq(prevTokenBalances[1], usdToken.balanceOf(address(registry)));
    assertEq(prevTokenBalances[2], weth.balanceOf(address(registry)));
    // assert reserve amounts have adjusted accordingly
    require(
      prevReserveBalances[0] < registry.getReserveAmount(address(linkToken)),
      "usd reserve amount should have increased"
    ); // link reserve amount increases in value equal to the decrease of the other reserve amounts
    require(
      prevReserveBalances[1] > registry.getReserveAmount(address(usdToken)),
      "usd reserve amount should have decreased"
    );
    require(
      prevReserveBalances[2] > registry.getReserveAmount(address(weth)),
      "native reserve amount should have decreased"
    );
  }
}

contract MigrateReceive is SetUp {
  event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination);
  event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom);

  Registry newRegistry;
  uint256[] idsToMigrate;

  function setUp() public override {
    super.setUp();
    (newRegistry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);
    idsToMigrate.push(linkUpkeepID);
    idsToMigrate.push(usdUpkeepID);
    idsToMigrate.push(nativeUpkeepID);
    registry.setPeerRegistryMigrationPermission(address(newRegistry), 1);
    newRegistry.setPeerRegistryMigrationPermission(address(registry), 2);
  }

  function test_RevertsWhen_PermissionsNotSet() external {
    // no permissions
    registry.setPeerRegistryMigrationPermission(address(newRegistry), 0);
    newRegistry.setPeerRegistryMigrationPermission(address(registry), 0);
    vm.expectRevert(Registry.MigrationNotPermitted.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));

    // only outgoing permissions
    registry.setPeerRegistryMigrationPermission(address(newRegistry), 1);
    newRegistry.setPeerRegistryMigrationPermission(address(registry), 0);
    vm.expectRevert(Registry.MigrationNotPermitted.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));

    // only incoming permissions
    registry.setPeerRegistryMigrationPermission(address(newRegistry), 0);
    newRegistry.setPeerRegistryMigrationPermission(address(registry), 2);
    vm.expectRevert(Registry.MigrationNotPermitted.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));

    // permissions opposite direction
    registry.setPeerRegistryMigrationPermission(address(newRegistry), 2);
    newRegistry.setPeerRegistryMigrationPermission(address(registry), 1);
    vm.expectRevert(Registry.MigrationNotPermitted.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));
  }

  function test_RevertsWhen_ReceivingRegistryDoesNotSupportToken() external {
    _removeBillingTokenConfig(newRegistry, address(weth));
    vm.expectRevert(Registry.InvalidToken.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));
    idsToMigrate.pop(); // remove native upkeep id
    vm.prank(UPKEEP_ADMIN);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry)); // should succeed now
  }

  function test_RevertsWhen_CalledByNonAdmin() external {
    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    vm.prank(STRANGER);
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));
  }

  function test_Success() external {
    vm.startPrank(UPKEEP_ADMIN);

    // add some changes in upkeep data to the mix
    registry.pauseUpkeep(usdUpkeepID);
    registry.setUpkeepTriggerConfig(linkUpkeepID, randomBytes(100));
    registry.setUpkeepCheckData(nativeUpkeepID, randomBytes(25));

    // record previous state
    uint256[] memory prevUpkeepBalances = new uint256[](3);
    prevUpkeepBalances[0] = registry.getBalance(linkUpkeepID);
    prevUpkeepBalances[1] = registry.getBalance(usdUpkeepID);
    prevUpkeepBalances[2] = registry.getBalance(nativeUpkeepID);
    uint256[] memory prevReserveBalances = new uint256[](3);
    prevReserveBalances[0] = registry.getReserveAmount(address(linkToken));
    prevReserveBalances[1] = registry.getReserveAmount(address(usdToken));
    prevReserveBalances[2] = registry.getReserveAmount(address(weth));
    uint256[] memory prevTokenBalances = new uint256[](3);
    prevTokenBalances[0] = linkToken.balanceOf(address(registry));
    prevTokenBalances[1] = usdToken.balanceOf(address(registry));
    prevTokenBalances[2] = weth.balanceOf(address(registry));
    bytes[] memory prevUpkeepData = new bytes[](3);
    prevUpkeepData[0] = abi.encode(registry.getUpkeep(linkUpkeepID));
    prevUpkeepData[1] = abi.encode(registry.getUpkeep(usdUpkeepID));
    prevUpkeepData[2] = abi.encode(registry.getUpkeep(nativeUpkeepID));
    bytes[] memory prevUpkeepTriggerData = new bytes[](3);
    prevUpkeepTriggerData[0] = registry.getUpkeepTriggerConfig(linkUpkeepID);
    prevUpkeepTriggerData[1] = registry.getUpkeepTriggerConfig(usdUpkeepID);
    prevUpkeepTriggerData[2] = registry.getUpkeepTriggerConfig(nativeUpkeepID);

    // event expectations
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(linkUpkeepID, prevUpkeepBalances[0], address(newRegistry));
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(usdUpkeepID, prevUpkeepBalances[1], address(newRegistry));
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(nativeUpkeepID, prevUpkeepBalances[2], address(newRegistry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(linkUpkeepID, prevUpkeepBalances[0], address(registry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(usdUpkeepID, prevUpkeepBalances[1], address(registry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(nativeUpkeepID, prevUpkeepBalances[2], address(registry));

    // do the thing
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));

    // assert upkeep balances have been migrated
    assertEq(registry.getBalance(linkUpkeepID), 0);
    assertEq(registry.getBalance(usdUpkeepID), 0);
    assertEq(registry.getBalance(nativeUpkeepID), 0);
    assertEq(newRegistry.getBalance(linkUpkeepID), prevUpkeepBalances[0]);
    assertEq(newRegistry.getBalance(usdUpkeepID), prevUpkeepBalances[1]);
    assertEq(newRegistry.getBalance(nativeUpkeepID), prevUpkeepBalances[2]);

    // assert reserve balances have been adjusted
    assertEq(newRegistry.getReserveAmount(address(linkToken)), newRegistry.getBalance(linkUpkeepID));
    assertEq(newRegistry.getReserveAmount(address(usdToken)), newRegistry.getBalance(usdUpkeepID));
    assertEq(newRegistry.getReserveAmount(address(weth)), newRegistry.getBalance(nativeUpkeepID));
    assertEq(
      newRegistry.getReserveAmount(address(linkToken)),
      prevReserveBalances[0] - registry.getReserveAmount(address(linkToken))
    );
    assertEq(
      newRegistry.getReserveAmount(address(usdToken)),
      prevReserveBalances[1] - registry.getReserveAmount(address(usdToken))
    );
    assertEq(
      newRegistry.getReserveAmount(address(weth)),
      prevReserveBalances[2] - registry.getReserveAmount(address(weth))
    );

    // assert token have been transfered
    assertEq(linkToken.balanceOf(address(newRegistry)), newRegistry.getBalance(linkUpkeepID));
    assertEq(usdToken.balanceOf(address(newRegistry)), newRegistry.getBalance(usdUpkeepID));
    assertEq(weth.balanceOf(address(newRegistry)), newRegistry.getBalance(nativeUpkeepID));
    assertEq(linkToken.balanceOf(address(registry)), prevTokenBalances[0] - linkToken.balanceOf(address(newRegistry)));
    assertEq(usdToken.balanceOf(address(registry)), prevTokenBalances[1] - usdToken.balanceOf(address(newRegistry)));
    assertEq(weth.balanceOf(address(registry)), prevTokenBalances[2] - weth.balanceOf(address(newRegistry)));

    // assert upkeep data matches
    assertEq(prevUpkeepData[0], abi.encode(newRegistry.getUpkeep(linkUpkeepID)));
    assertEq(prevUpkeepData[1], abi.encode(newRegistry.getUpkeep(usdUpkeepID)));
    assertEq(prevUpkeepData[2], abi.encode(newRegistry.getUpkeep(nativeUpkeepID)));
    assertEq(prevUpkeepTriggerData[0], newRegistry.getUpkeepTriggerConfig(linkUpkeepID));
    assertEq(prevUpkeepTriggerData[1], newRegistry.getUpkeepTriggerConfig(usdUpkeepID));
    assertEq(prevUpkeepTriggerData[2], newRegistry.getUpkeepTriggerConfig(nativeUpkeepID));

    vm.stopPrank();
  }
}
