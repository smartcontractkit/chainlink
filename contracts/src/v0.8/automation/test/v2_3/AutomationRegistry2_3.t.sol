// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Vm} from "forge-std/Test.sol";
import {BaseTest} from "./BaseTest.t.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../../v2_3/AutomationRegistryBase2_3.sol";
import {AutomationRegistrar2_3 as Registrar} from "../../v2_3/AutomationRegistrar2_3.sol";
import {IAutomationRegistryMaster2_3 as Registry, AutomationRegistryBase2_3, IAutomationV21PlusCommon} from "../../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {ChainModuleBase} from "../../chains/ChainModuleBase.sol";
import {IERC20Metadata as IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {IWrappedNative} from "../../interfaces/v2_3/IWrappedNative.sol";

// forge test --match-path src/v0.8/automation/test/v2_3/AutomationRegistry2_3.t.sol

enum Trigger {
  CONDITION,
  LOG
}

contract SetUp is BaseTest {
  Registry internal registry;
  AutomationRegistryBase2_3.OnchainConfig internal config;
  bytes internal constant offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

  uint256 internal linkUpkeepID;
  uint256 internal linkUpkeepID2; // 2 upkeeps use the same billing token (LINK) to test migration scenario
  uint256 internal usdUpkeepID18; // 1 upkeep uses ERC20 token with 18 decimals
  uint256 internal usdUpkeepID6; // 1 upkeep uses ERC20 token with 6 decimals
  uint256 internal nativeUpkeepID;

  function setUp() public virtual override {
    super.setUp();
    (registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);
    config = registry.getConfig();

    vm.startPrank(OWNER);
    linkToken.approve(address(registry), type(uint256).max);
    usdToken6.approve(address(registry), type(uint256).max);
    usdToken18.approve(address(registry), type(uint256).max);
    weth.approve(address(registry), type(uint256).max);
    vm.startPrank(UPKEEP_ADMIN);
    linkToken.approve(address(registry), type(uint256).max);
    usdToken6.approve(address(registry), type(uint256).max);
    usdToken18.approve(address(registry), type(uint256).max);
    weth.approve(address(registry), type(uint256).max);
    vm.startPrank(STRANGER);
    linkToken.approve(address(registry), type(uint256).max);
    usdToken6.approve(address(registry), type(uint256).max);
    usdToken18.approve(address(registry), type(uint256).max);
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

    linkUpkeepID2 = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(linkToken),
      "",
      "",
      ""
    );

    usdUpkeepID18 = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(usdToken18),
      "",
      "",
      ""
    );

    usdUpkeepID6 = registry.registerUpkeep(
      address(TARGET1),
      config.maxPerformGas,
      UPKEEP_ADMIN,
      uint8(Trigger.CONDITION),
      address(usdToken6),
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
    registry.addFunds(linkUpkeepID2, registry.getMinBalanceForUpkeep(linkUpkeepID2));
    registry.addFunds(usdUpkeepID18, registry.getMinBalanceForUpkeep(usdUpkeepID18));
    registry.addFunds(usdUpkeepID6, registry.getMinBalanceForUpkeep(usdUpkeepID6));
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

contract WithdrawFunds is SetUp {
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);

  function test_RevertsWhen_CalledByNonAdmin() external {
    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    vm.prank(STRANGER);
    registry.withdrawFunds(linkUpkeepID, STRANGER);
  }

  function test_RevertsWhen_InvalidRecipient() external {
    vm.expectRevert(Registry.InvalidRecipient.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.withdrawFunds(linkUpkeepID, ZERO_ADDRESS);
  }

  function test_RevertsWhen_UpkeepNotCanceled() external {
    vm.expectRevert(Registry.UpkeepNotCanceled.selector);
    vm.prank(UPKEEP_ADMIN);
    registry.withdrawFunds(linkUpkeepID, UPKEEP_ADMIN);
  }

  function test_Happy_Link() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);
    vm.roll(100 + block.number);

    uint256 startUpkeepAdminBalance = linkToken.balanceOf(UPKEEP_ADMIN);
    uint256 startLinkReserveAmountBalance = registry.getReserveAmount(address(linkToken));

    uint256 upkeepBalance = registry.getBalance(linkUpkeepID);
    vm.expectEmit();
    emit FundsWithdrawn(linkUpkeepID, upkeepBalance, address(UPKEEP_ADMIN));
    registry.withdrawFunds(linkUpkeepID, UPKEEP_ADMIN);

    assertEq(registry.getBalance(linkUpkeepID), 0);
    assertEq(linkToken.balanceOf(UPKEEP_ADMIN), startUpkeepAdminBalance + upkeepBalance);
    assertEq(registry.getReserveAmount(address(linkToken)), startLinkReserveAmountBalance - upkeepBalance);
  }

  function test_Happy_USDToken() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(usdUpkeepID6);
    vm.roll(100 + block.number);

    uint256 startUpkeepAdminBalance = usdToken6.balanceOf(UPKEEP_ADMIN);
    uint256 startUSDToken6ReserveAmountBalance = registry.getReserveAmount(address(usdToken6));

    uint256 upkeepBalance = registry.getBalance(usdUpkeepID6);
    vm.expectEmit();
    emit FundsWithdrawn(usdUpkeepID6, upkeepBalance, address(UPKEEP_ADMIN));
    registry.withdrawFunds(usdUpkeepID6, UPKEEP_ADMIN);

    assertEq(registry.getBalance(usdUpkeepID6), 0);
    assertEq(usdToken6.balanceOf(UPKEEP_ADMIN), startUpkeepAdminBalance + upkeepBalance);
    assertEq(registry.getReserveAmount(address(usdToken6)), startUSDToken6ReserveAmountBalance - upkeepBalance);
  }
}

contract AddFunds is SetUp {
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);

  // when msg.value is 0, it uses the ERC20 payment path
  function test_HappyWhen_NativeUpkeep_WithMsgValue0() external {
    vm.startPrank(OWNER);
    uint256 startRegistryBalance = registry.getBalance(nativeUpkeepID);
    uint256 startTokenBalance = registry.getBalance(nativeUpkeepID);
    registry.addFunds(nativeUpkeepID, 1);
    assertEq(registry.getBalance(nativeUpkeepID), startRegistryBalance + 1);
    assertEq(weth.balanceOf(address(registry)), startTokenBalance + 1);
    assertEq(registry.getAvailableERC20ForPayment(address(weth)), 0);
  }

  // when msg.value is not 0, it uses the native payment path
  function test_HappyWhen_NativeUpkeep_WithMsgValueNot0() external {
    uint256 startRegistryBalance = registry.getBalance(nativeUpkeepID);
    uint256 startTokenBalance = registry.getBalance(nativeUpkeepID);
    registry.addFunds{value: 1}(nativeUpkeepID, 1000); // parameter amount should be ignored
    assertEq(registry.getBalance(nativeUpkeepID), startRegistryBalance + 1);
    assertEq(weth.balanceOf(address(registry)), startTokenBalance + 1);
    assertEq(registry.getAvailableERC20ForPayment(address(weth)), 0);
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

    uint256 startLINKRegistryBalance = linkToken.balanceOf(address(registry));
    uint256 startUSDRegistryBalance = usdToken18.balanceOf(address(registry));
    uint256 startLinkUpkeepBalance = registry.getBalance(linkUpkeepID);
    uint256 startUSDUpkeepBalance = registry.getBalance(usdUpkeepID18);

    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), startLinkUpkeepBalance + 1);
    assertEq(registry.getBalance(usdUpkeepID18), startUSDRegistryBalance);
    assertEq(linkToken.balanceOf(address(registry)), startLINKRegistryBalance + 1);
    assertEq(usdToken18.balanceOf(address(registry)), startUSDUpkeepBalance);

    registry.addFunds(usdUpkeepID18, 2);
    assertEq(registry.getBalance(linkUpkeepID), startLinkUpkeepBalance + 1);
    assertEq(registry.getBalance(usdUpkeepID18), startUSDRegistryBalance + 2);
    assertEq(linkToken.balanceOf(address(registry)), startLINKRegistryBalance + 1);
    assertEq(usdToken18.balanceOf(address(registry)), startUSDUpkeepBalance + 2);
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
    assertEq(registry.getBalance(usdUpkeepID18), registry.getReserveAmount(address(usdToken18)));
    vm.startPrank(FINANCE_ADMIN);
    vm.expectRevert(abi.encodeWithSelector(Registry.InsufficientBalance.selector, 0, 1));
    registry.withdrawERC20Fees(address(usdToken18), FINANCE_ADMIN, 1);
  }

  function test_WithdrawERC20Fees_RevertsWhen_AttemptingToWithdrawLINK() public {
    _mintLink(address(registry), 1e10);
    vm.startPrank(FINANCE_ADMIN);
    vm.expectRevert(Registry.InvalidToken.selector);
    registry.withdrawERC20Fees(address(linkToken), FINANCE_ADMIN, 1); // should revert
    registry.withdrawLink(FINANCE_ADMIN, 1); // but using link withdraw functions succeeds
  }

  // default is ON_CHAIN mode
  function test_WithdrawERC20Fees_RevertsWhen_LinkAvailableForPaymentIsNegative() public {
    _transmit(usdUpkeepID18, registry); // adds USD token to finance withdrawable, and gives NOPs a LINK balance
    require(registry.linkAvailableForPayment() < 0, "linkAvailableForPayment should be negative");
    require(
      registry.getAvailableERC20ForPayment(address(usdToken18)) > 0,
      "ERC20AvailableForPayment should be positive"
    );
    vm.expectRevert(Registry.InsufficientLinkLiquidity.selector);
    vm.prank(FINANCE_ADMIN);
    registry.withdrawERC20Fees(address(usdToken18), FINANCE_ADMIN, 1); // should revert
    _mintLink(address(registry), uint256(registry.linkAvailableForPayment() * -10)); // top up LINK liquidity pool
    vm.prank(FINANCE_ADMIN);
    registry.withdrawERC20Fees(address(usdToken18), FINANCE_ADMIN, 1); // now finance can withdraw
  }

  function test_WithdrawERC20Fees_InOffChainMode_Happy() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken18), "", "", "");
    _mintERC20_18Decimals(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken18.approve(address(registry), 1e20);
    registry.addFunds(id, 1e20);

    // manually create a transmit so transmitters earn some rewards
    _transmit(id, registry);
    require(registry.linkAvailableForPayment() < 0, "linkAvailableForPayment should be negative");
    vm.prank(FINANCE_ADMIN);
    registry.withdrawERC20Fees(address(usdToken18), aMockAddress, 1); // finance can withdraw

    // recipient should get the funds
    assertEq(usdToken18.balanceOf(address(aMockAddress)), 1);
  }

  function testWithdrawERC20FeeSuccess() public {
    // deposit excess USDToken to the registry (this goes to the "finance withdrawable" pool be default)
    uint256 startReserveAmount = registry.getReserveAmount(address(usdToken18));
    uint256 startAmount = usdToken18.balanceOf(address(registry));
    _mintERC20_18Decimals(address(registry), 1e10);

    // depositing shouldn't change reserve amount
    assertEq(registry.getReserveAmount(address(usdToken18)), startReserveAmount);

    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 USDToken
    registry.withdrawERC20Fees(address(usdToken18), aMockAddress, 1);

    vm.stopPrank();

    assertEq(usdToken18.balanceOf(address(aMockAddress)), 1);
    assertEq(usdToken18.balanceOf(address(registry)), startAmount + 1e10 - 1);
    assertEq(registry.getReserveAmount(address(usdToken18)), startReserveAmount);
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
      registrars: _getRegistrars(),
      upkeepPrivilegeManager: PRIVILEGE_MANAGER,
      chainModule: module,
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });

  function testSetConfigSuccess() public {
    (uint32 configCount, uint32 blockNumber, ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    address billingTokenAddress = address(usdToken18);
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = billingTokenAddress;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000,
      decimals: 18
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
    assertEq(config.priceFeed, address(USDTOKEN_USD_FEED));
    assertEq(config.minSpend, 100_000);

    address[] memory tokens = registry.getBillingTokens();
    assertEq(tokens.length, 1);
  }

  function testSetConfigMultipleBillingConfigsSuccess() public {
    (uint32 configCount, , ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    address billingTokenAddress1 = address(linkToken);
    address billingTokenAddress2 = address(usdToken18);
    address[] memory billingTokens = new address[](2);
    billingTokens[0] = billingTokenAddress1;
    billingTokens[1] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMilliCents: 20_001,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 100,
      minSpend: 100,
      decimals: 18
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMilliCents: 20_002,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 200,
      minSpend: 200,
      decimals: 18
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
    assertEq(config1.priceFeed, address(USDTOKEN_USD_FEED));
    assertEq(config1.fallbackPrice, 100);
    assertEq(config1.minSpend, 100);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registry.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMilliCents, 20_002);
    assertEq(config2.priceFeed, address(USDTOKEN_USD_FEED));
    assertEq(config2.fallbackPrice, 200);
    assertEq(config2.minSpend, 200);

    address[] memory tokens = registry.getBillingTokens();
    assertEq(tokens.length, 2);
  }

  function testSetConfigTwiceAndLastSetOverwrites() public {
    (uint32 configCount, , ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    // BillingConfig1
    address billingTokenAddress1 = address(usdToken18);
    address[] memory billingTokens1 = new address[](1);
    billingTokens1[0] = billingTokenAddress1;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs1 = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs1[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMilliCents: 20_001,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 100,
      minSpend: 100,
      decimals: 18
    });

    // the first time uses the default onchain config with 2 registrars
    bytes memory onchainConfigBytesWithBilling1 = abi.encode(cfg, billingTokens1, billingConfigs1);

    // set config once
    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling1,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, IAutomationV21PlusCommon.OnchainConfigLegacy memory onchainConfig1, , , ) = registry.getState();
    assertEq(onchainConfig1.registrars.length, 2);

    // BillingConfig2
    address billingTokenAddress2 = address(usdToken18);
    address[] memory billingTokens2 = new address[](1);
    billingTokens2[0] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs2 = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs2[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMilliCents: 20_002,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 200,
      minSpend: 200,
      decimals: 18
    });

    address[] memory newRegistrars = new address[](3);
    newRegistrars[0] = address(uint160(uint256(keccak256("newRegistrar1"))));
    newRegistrars[1] = address(uint160(uint256(keccak256("newRegistrar2"))));
    newRegistrars[2] = address(uint160(uint256(keccak256("newRegistrar3"))));

    // new onchain config with 3 new registrars, all other fields stay the same as the default
    AutomationRegistryBase2_3.OnchainConfig memory cfg2 = AutomationRegistryBase2_3.OnchainConfig({
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
      registrars: newRegistrars,
      upkeepPrivilegeManager: PRIVILEGE_MANAGER,
      chainModule: module,
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });

    // the second time uses the new onchain config with 3 new registrars and also new billing tokens/configs
    bytes memory onchainConfigBytesWithBilling2 = abi.encode(cfg2, billingTokens2, billingConfigs2);

    // set config twice
    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling2,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (
      ,
      IAutomationV21PlusCommon.OnchainConfigLegacy memory onchainConfig2,
      address[] memory signers,
      address[] memory transmitters,
      uint8 f
    ) = registry.getState();

    assertEq(onchainConfig2.registrars.length, 3);
    for (uint256 i = 0; i < newRegistrars.length; i++) {
      assertEq(newRegistrars[i], onchainConfig2.registrars[i]);
    }
    assertEq(signers, SIGNERS);
    assertEq(transmitters, TRANSMITTERS);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registry.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMilliCents, 20_002);
    assertEq(config2.priceFeed, address(USDTOKEN_USD_FEED));
    assertEq(config2.fallbackPrice, 200);
    assertEq(config2.minSpend, 200);

    address[] memory tokens = registry.getBillingTokens();
    assertEq(tokens.length, 1);
  }

  function testSetConfigDuplicateBillingConfigFailure() public {
    (uint32 configCount, , ) = registry.latestConfigDetails();
    assertEq(configCount, 1);

    address billingTokenAddress1 = address(linkToken);
    address billingTokenAddress2 = address(linkToken);
    address[] memory billingTokens = new address[](2);
    billingTokens[0] = billingTokenAddress1;
    billingTokens[1] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMilliCents: 20_001,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 100,
      minSpend: 100,
      decimals: 18
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMilliCents: 20_002,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 200,
      minSpend: 200,
      decimals: 18
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
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000,
      decimals: 18
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

  function testSetConfigRevertDueToInvalidDecimals() public {
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = address(linkToken);

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000,
      decimals: 6 // link token should have 18 decimals
    });

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

  function testSetConfigOnTransmittersAndPayees() public {
    registry.setPayees(PAYEES);
    AutomationRegistryBase2_3.TransmitterPayeeInfo[] memory transmitterPayeeInfos = registry
      .getTransmittersWithPayees();
    assertEq(transmitterPayeeInfos.length, TRANSMITTERS.length);

    for (uint256 i = 0; i < transmitterPayeeInfos.length; i++) {
      address transmitterAddress = transmitterPayeeInfos[i].transmitterAddress;
      address payeeAddress = transmitterPayeeInfos[i].payeeAddress;

      address expectedTransmitter = TRANSMITTERS[i];
      address expectedPayee = PAYEES[i];

      assertEq(transmitterAddress, expectedTransmitter);
      assertEq(payeeAddress, expectedPayee);
    }
  }

  function testSetConfigWithNewTransmittersSuccess() public {
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);

    (uint32 configCount, uint32 blockNumber, ) = registry.latestConfigDetails();
    assertEq(configCount, 0);

    address billingTokenAddress = address(usdToken18);
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = billingTokenAddress;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMilliCents: 20_000,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000,
      decimals: 18
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

  function _getRegistrars() private pure returns (address[] memory) {
    address[] memory registrars = new address[](2);
    registrars[0] = address(uint160(uint256(keccak256("registrar1"))));
    registrars[1] = address(uint160(uint256(keccak256("registrar2"))));
    return registrars;
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
    registry.setPayees(PAYEES);

    uint256[] memory payments = new uint256[](TRANSMITTERS.length);
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      payments[i] = 0;
    }

    vm.startPrank(FINANCE_ADMIN);
    vm.expectEmit();
    emit NOPsSettledOffchain(PAYEES, payments);
    registry.settleNOPsOffchain();
  }

  // 1. transmitter balance zeroed after settlement, 2. admin can withdraw ERC20, 3. switch to onchain mode, 4. link amount owed to NOPs stays the same
  function testSettleNOPsOffchainSuccessWithERC20MultiSteps() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);
    registry.setPayees(PAYEES);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken18), "", "", "");
    _mintERC20_18Decimals(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken18.approve(address(registry), 1e20);
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

    // after the offchain settlement, the total reserve amount of LINK should be 0
    assertEq(registry.getReserveAmount(address(linkToken)), 0);
    // should have some ERC20s in registry after transmit
    uint256 erc20ForPayment1 = registry.getAvailableERC20ForPayment(address(usdToken18));
    require(erc20ForPayment1 > 0, "ERC20AvailableForPayment should be positive");

    vm.startPrank(UPKEEP_ADMIN);
    vm.roll(100 + block.number);
    // manually create a transmit so transmitters earn some rewards
    _transmit(id, registry);

    uint256 erc20ForPayment2 = registry.getAvailableERC20ForPayment(address(usdToken18));
    require(erc20ForPayment2 > erc20ForPayment1, "ERC20AvailableForPayment should be greater after another transmit");

    // finance admin comes to withdraw all available ERC20s
    vm.startPrank(FINANCE_ADMIN);
    registry.withdrawERC20Fees(address(usdToken18), FINANCE_ADMIN, erc20ForPayment2);

    uint256 erc20ForPayment3 = registry.getAvailableERC20ForPayment(address(usdToken18));
    require(erc20ForPayment3 == 0, "ERC20AvailableForPayment should be 0 now after withdrawal");

    uint256 reservedLink = registry.getReserveAmount(address(linkToken));
    require(reservedLink > 0, "Reserve amount of LINK should be positive since there was another transmit");

    // owner comes to disable offchain mode
    vm.startPrank(registry.owner());
    registry.disableOffchainPayments();

    // finance admin comes to withdraw all available ERC20s, should revert bc of insufficient link liquidity
    vm.startPrank(FINANCE_ADMIN);
    uint256 erc20ForPayment4 = registry.getAvailableERC20ForPayment(address(usdToken18));
    vm.expectRevert(abi.encodeWithSelector(Registry.InsufficientLinkLiquidity.selector));
    registry.withdrawERC20Fees(address(usdToken18), FINANCE_ADMIN, erc20ForPayment4);

    // reserved link amount to NOPs should stay the same after switching to onchain mode
    assertEq(registry.getReserveAmount(address(linkToken)), reservedLink);
    // available ERC20 for payment should be 0 since finance admin withdrew all already
    assertEq(erc20ForPayment4, 0);
  }

  function testSettleNOPsOffchainForDeactivatedTransmittersSuccess() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (Registry registry, Registrar registrar) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken18), "", "", "");
    _mintERC20_18Decimals(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken18.approve(address(registry), 1e20);
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

    // after the offchain settlement, the total reserve amount of LINK should be 0
    assertEq(registry.getReserveAmount(address(linkToken)), 0);
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
    registry.setPayees(PAYEES);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken18), "", "", "");
    _mintERC20_18Decimals(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken18.approve(address(registry), 1e20);
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
    registry.setPayees(PAYEES);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(usdToken18), "", "", "");
    _mintERC20_18Decimals(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken18.approve(address(registry), 1e20);
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
    billingTokens[0] = IERC20(address(usdToken18));

    uint256[] memory minRegistrationFees = new uint256[](billingTokens.length);
    minRegistrationFees[0] = 100e18; // 100 USD

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
      fallbackPrice: 1e8, // $1
      minSpend: 1e18, // 1 USD
      decimals: 18
    });

    address[] memory registrars = new address[](1);
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
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, abi.encode(usdUpkeepID18));
  }

  function test_Happy() public {
    vm.startPrank(address(linkToken));
    uint256 beforeBalance = registry.getBalance(linkUpkeepID);
    registry.onTokenTransfer(UPKEEP_ADMIN, 100, abi.encode(linkUpkeepID));
    assertEq(registry.getBalance(linkUpkeepID), beforeBalance + 100);
  }
}

contract GetMinBalanceForUpkeep is SetUp {
  function test_accountsForFlatFee_with18Decimals() public {
    // set fee to 0
    AutomationRegistryBase2_3.BillingConfig memory usdTokenConfig = registry.getBillingTokenConfig(address(usdToken18));
    usdTokenConfig.flatFeeMilliCents = 0;
    _updateBillingTokenConfig(registry, address(usdToken18), usdTokenConfig);

    uint256 minBalanceBefore = registry.getMinBalanceForUpkeep(usdUpkeepID18);

    // set fee to non-zero
    usdTokenConfig.flatFeeMilliCents = 100;
    _updateBillingTokenConfig(registry, address(usdToken18), usdTokenConfig);

    uint256 minBalanceAfter = registry.getMinBalanceForUpkeep(usdUpkeepID18);
    assertEq(
      minBalanceAfter,
      minBalanceBefore + ((uint256(usdTokenConfig.flatFeeMilliCents) * 1e13) / 10 ** (18 - usdTokenConfig.decimals))
    );
  }

  function test_accountsForFlatFee_with6Decimals() public {
    // set fee to 0
    AutomationRegistryBase2_3.BillingConfig memory usdTokenConfig = registry.getBillingTokenConfig(address(usdToken6));
    usdTokenConfig.flatFeeMilliCents = 0;
    _updateBillingTokenConfig(registry, address(usdToken6), usdTokenConfig);

    uint256 minBalanceBefore = registry.getMinBalanceForUpkeep(usdUpkeepID6);

    // set fee to non-zero
    usdTokenConfig.flatFeeMilliCents = 100;
    _updateBillingTokenConfig(registry, address(usdToken6), usdTokenConfig);

    uint256 minBalanceAfter = registry.getMinBalanceForUpkeep(usdUpkeepID6);
    assertEq(
      minBalanceAfter,
      minBalanceBefore + ((uint256(usdTokenConfig.flatFeeMilliCents) * 1e13) / 10 ** (18 - usdTokenConfig.decimals))
    );
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
  function test_transmitRevertWithExtraBytes() external {
    bytes32[3] memory exampleReportContext = [
      bytes32(0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef),
      bytes32(0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890),
      bytes32(0x7890abcdef1234567890abcdef1234567890abcdef1234567890abcdef123456)
    ];

    bytes memory exampleRawReport = hex"deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef";

    bytes32[] memory exampleRs = new bytes32[](3);
    exampleRs[0] = bytes32(0x1234561234561234561234561234561234561234561234561234561234561234);
    exampleRs[1] = bytes32(0x1234561234561234561234561234561234561234561234561234561234561234);
    exampleRs[2] = bytes32(0x7890789078907890789078907890789078907890789078907890789078907890);

    bytes32[] memory exampleSs = new bytes32[](3);
    exampleSs[0] = bytes32(0x1234561234561234561234561234561234561234561234561234561234561234);
    exampleSs[1] = bytes32(0x1234561234561234561234561234561234561234561234561234561234561234);
    exampleSs[2] = bytes32(0x1234561234561234561234561234561234561234561234561234561234561234);

    bytes32 exampleRawVs = bytes32(0x1234561234561234561234561234561234561234561234561234561234561234);

    bytes memory transmitData = abi.encodeWithSelector(
      registry.transmit.selector,
      exampleReportContext,
      exampleRawReport,
      exampleRs,
      exampleSs,
      exampleRawVs
    );
    bytes memory badTransmitData = bytes.concat(transmitData, bytes1(0x00)); // add extra data
    vm.startPrank(TRANSMITTERS[0]);
    (bool success, bytes memory returnData) = address(registry).call(badTransmitData); // send the bogus transmit
    assertFalse(success, "Call did not revert as expected");
    assertEq(returnData, abi.encodePacked(Registry.InvalidDataLength.selector));
    vm.stopPrank();
  }

  function test_handlesMixedBatchOfBillingTokens() external {
    uint256[] memory prevUpkeepBalances = new uint256[](3);
    prevUpkeepBalances[0] = registry.getBalance(linkUpkeepID);
    prevUpkeepBalances[1] = registry.getBalance(usdUpkeepID18);
    prevUpkeepBalances[2] = registry.getBalance(nativeUpkeepID);
    uint256[] memory prevTokenBalances = new uint256[](3);
    prevTokenBalances[0] = linkToken.balanceOf(address(registry));
    prevTokenBalances[1] = usdToken18.balanceOf(address(registry));
    prevTokenBalances[2] = weth.balanceOf(address(registry));
    uint256[] memory prevReserveBalances = new uint256[](3);
    prevReserveBalances[0] = registry.getReserveAmount(address(linkToken));
    prevReserveBalances[1] = registry.getReserveAmount(address(usdToken18));
    prevReserveBalances[2] = registry.getReserveAmount(address(weth));
    uint256[] memory upkeepIDs = new uint256[](3);
    upkeepIDs[0] = linkUpkeepID;
    upkeepIDs[1] = usdUpkeepID18;
    upkeepIDs[2] = nativeUpkeepID;

    // withdraw-able by finance team should be 0
    require(registry.getAvailableERC20ForPayment(address(usdToken18)) == 0, "ERC20AvailableForPayment should be 0");
    require(registry.getAvailableERC20ForPayment(address(weth)) == 0, "ERC20AvailableForPayment should be 0");

    // do the thing
    _transmit(upkeepIDs, registry);

    // withdraw-able by the finance team should be positive
    require(
      registry.getAvailableERC20ForPayment(address(usdToken18)) > 0,
      "ERC20AvailableForPayment should be positive"
    );
    require(registry.getAvailableERC20ForPayment(address(weth)) > 0, "ERC20AvailableForPayment should be positive");

    // assert upkeep balances have decreased
    require(prevUpkeepBalances[0] > registry.getBalance(linkUpkeepID), "link upkeep balance should have decreased");
    require(prevUpkeepBalances[1] > registry.getBalance(usdUpkeepID18), "usd upkeep balance should have decreased");
    require(prevUpkeepBalances[2] > registry.getBalance(nativeUpkeepID), "native upkeep balance should have decreased");
    // assert token balances have not changed
    assertEq(prevTokenBalances[0], linkToken.balanceOf(address(registry)));
    assertEq(prevTokenBalances[1], usdToken18.balanceOf(address(registry)));
    assertEq(prevTokenBalances[2], weth.balanceOf(address(registry)));
    // assert reserve amounts have adjusted accordingly
    require(
      prevReserveBalances[0] < registry.getReserveAmount(address(linkToken)),
      "usd reserve amount should have increased"
    ); // link reserve amount increases in value equal to the decrease of the other reserve amounts
    require(
      prevReserveBalances[1] > registry.getReserveAmount(address(usdToken18)),
      "usd reserve amount should have decreased"
    );
    require(
      prevReserveBalances[2] > registry.getReserveAmount(address(weth)),
      "native reserve amount should have decreased"
    );
  }

  function test_handlesInsufficientBalanceWithUSDToken18() external {
    // deploy and configure a registry with ON_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);

    // register an upkeep and add funds
    uint256 upkeepID = registry.registerUpkeep(
      address(TARGET1),
      1000000,
      UPKEEP_ADMIN,
      0,
      address(usdToken18),
      "",
      "",
      ""
    );
    _mintERC20_18Decimals(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    usdToken18.approve(address(registry), 1e20);
    registry.addFunds(upkeepID, 1); // smaller than gasCharge
    uint256 balance = registry.getBalance(upkeepID);

    // manually create a transmit
    vm.recordLogs();
    _transmit(upkeepID, registry);
    Vm.Log[] memory entries = vm.getRecordedLogs();

    assertEq(entries.length, 3);
    Vm.Log memory l1 = entries[1];
    assertEq(
      l1.topics[0],
      keccak256("UpkeepCharged(uint256,(uint96,uint96,uint96,uint96,address,uint96,uint96,uint96))")
    );
    (
      uint96 gasChargeInBillingToken,
      uint96 premiumInBillingToken,
      uint96 gasReimbursementInJuels,
      uint96 premiumInJuels,
      address billingToken,
      uint96 linkUSD,
      uint96 nativeUSD,
      uint96 billingUSD
    ) = abi.decode(l1.data, (uint96, uint96, uint96, uint96, address, uint96, uint96, uint96));

    assertEq(gasChargeInBillingToken, balance);
    assertEq(gasReimbursementInJuels, (balance * billingUSD) / linkUSD);
    assertEq(premiumInJuels, 0);
    assertEq(premiumInBillingToken, 0);
  }

  function test_handlesInsufficientBalanceWithUSDToken6() external {
    // deploy and configure a registry with ON_CHAIN payout
    (Registry registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);

    // register an upkeep and add funds
    uint256 upkeepID = registry.registerUpkeep(
      address(TARGET1),
      1000000,
      UPKEEP_ADMIN,
      0,
      address(usdToken6),
      "",
      "",
      ""
    );
    vm.prank(OWNER);
    usdToken6.mint(UPKEEP_ADMIN, 1e20);

    vm.startPrank(UPKEEP_ADMIN);
    usdToken6.approve(address(registry), 1e20);
    registry.addFunds(upkeepID, 100); // this is greater than gasCharge but less than (gasCharge + premium)
    uint256 balance = registry.getBalance(upkeepID);

    // manually create a transmit
    vm.recordLogs();
    _transmit(upkeepID, registry);
    Vm.Log[] memory entries = vm.getRecordedLogs();

    assertEq(entries.length, 3);
    Vm.Log memory l1 = entries[1];
    assertEq(
      l1.topics[0],
      keccak256("UpkeepCharged(uint256,(uint96,uint96,uint96,uint96,address,uint96,uint96,uint96))")
    );
    (
      uint96 gasChargeInBillingToken,
      uint96 premiumInBillingToken,
      uint96 gasReimbursementInJuels,
      uint96 premiumInJuels,
      address billingToken,
      uint96 linkUSD,
      uint96 nativeUSD,
      uint96 billingUSD
    ) = abi.decode(l1.data, (uint96, uint96, uint96, uint96, address, uint96, uint96, uint96));

    assertEq(premiumInJuels, (balance * billingUSD * 1e12) / linkUSD - gasReimbursementInJuels); // scale to 18 decimals
    assertEq(premiumInBillingToken, (premiumInJuels * linkUSD + (billingUSD * 1e12 - 1)) / (billingUSD * 1e12));
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
    idsToMigrate.push(linkUpkeepID2);
    idsToMigrate.push(usdUpkeepID18);
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
    registry.pauseUpkeep(usdUpkeepID18);
    registry.setUpkeepTriggerConfig(linkUpkeepID, randomBytes(100));
    registry.setUpkeepCheckData(nativeUpkeepID, randomBytes(25));

    // record previous state
    uint256[] memory prevUpkeepBalances = new uint256[](4);
    prevUpkeepBalances[0] = registry.getBalance(linkUpkeepID);
    prevUpkeepBalances[1] = registry.getBalance(linkUpkeepID2);
    prevUpkeepBalances[2] = registry.getBalance(usdUpkeepID18);
    prevUpkeepBalances[3] = registry.getBalance(nativeUpkeepID);
    uint256[] memory prevReserveBalances = new uint256[](3);
    prevReserveBalances[0] = registry.getReserveAmount(address(linkToken));
    prevReserveBalances[1] = registry.getReserveAmount(address(usdToken18));
    prevReserveBalances[2] = registry.getReserveAmount(address(weth));
    uint256[] memory prevTokenBalances = new uint256[](3);
    prevTokenBalances[0] = linkToken.balanceOf(address(registry));
    prevTokenBalances[1] = usdToken18.balanceOf(address(registry));
    prevTokenBalances[2] = weth.balanceOf(address(registry));
    bytes[] memory prevUpkeepData = new bytes[](4);
    prevUpkeepData[0] = abi.encode(registry.getUpkeep(linkUpkeepID));
    prevUpkeepData[1] = abi.encode(registry.getUpkeep(linkUpkeepID2));
    prevUpkeepData[2] = abi.encode(registry.getUpkeep(usdUpkeepID18));
    prevUpkeepData[3] = abi.encode(registry.getUpkeep(nativeUpkeepID));
    bytes[] memory prevUpkeepTriggerData = new bytes[](4);
    prevUpkeepTriggerData[0] = registry.getUpkeepTriggerConfig(linkUpkeepID);
    prevUpkeepTriggerData[1] = registry.getUpkeepTriggerConfig(linkUpkeepID2);
    prevUpkeepTriggerData[2] = registry.getUpkeepTriggerConfig(usdUpkeepID18);
    prevUpkeepTriggerData[3] = registry.getUpkeepTriggerConfig(nativeUpkeepID);

    // event expectations
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(linkUpkeepID, prevUpkeepBalances[0], address(newRegistry));
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(linkUpkeepID2, prevUpkeepBalances[1], address(newRegistry));
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(usdUpkeepID18, prevUpkeepBalances[2], address(newRegistry));
    vm.expectEmit(address(registry));
    emit UpkeepMigrated(nativeUpkeepID, prevUpkeepBalances[3], address(newRegistry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(linkUpkeepID, prevUpkeepBalances[0], address(registry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(linkUpkeepID2, prevUpkeepBalances[1], address(registry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(usdUpkeepID18, prevUpkeepBalances[2], address(registry));
    vm.expectEmit(address(newRegistry));
    emit UpkeepReceived(nativeUpkeepID, prevUpkeepBalances[3], address(registry));

    // do the thing
    registry.migrateUpkeeps(idsToMigrate, address(newRegistry));

    // assert upkeep balances have been migrated
    assertEq(registry.getBalance(linkUpkeepID), 0);
    assertEq(registry.getBalance(linkUpkeepID2), 0);
    assertEq(registry.getBalance(usdUpkeepID18), 0);
    assertEq(registry.getBalance(nativeUpkeepID), 0);
    assertEq(newRegistry.getBalance(linkUpkeepID), prevUpkeepBalances[0]);
    assertEq(newRegistry.getBalance(linkUpkeepID2), prevUpkeepBalances[1]);
    assertEq(newRegistry.getBalance(usdUpkeepID18), prevUpkeepBalances[2]);
    assertEq(newRegistry.getBalance(nativeUpkeepID), prevUpkeepBalances[3]);

    // assert reserve balances have been adjusted
    assertEq(
      newRegistry.getReserveAmount(address(linkToken)),
      newRegistry.getBalance(linkUpkeepID) + newRegistry.getBalance(linkUpkeepID2)
    );
    assertEq(newRegistry.getReserveAmount(address(usdToken18)), newRegistry.getBalance(usdUpkeepID18));
    assertEq(newRegistry.getReserveAmount(address(weth)), newRegistry.getBalance(nativeUpkeepID));
    assertEq(
      newRegistry.getReserveAmount(address(linkToken)),
      prevReserveBalances[0] - registry.getReserveAmount(address(linkToken))
    );
    assertEq(
      newRegistry.getReserveAmount(address(usdToken18)),
      prevReserveBalances[1] - registry.getReserveAmount(address(usdToken18))
    );
    assertEq(
      newRegistry.getReserveAmount(address(weth)),
      prevReserveBalances[2] - registry.getReserveAmount(address(weth))
    );

    // assert token have been transferred
    assertEq(
      linkToken.balanceOf(address(newRegistry)),
      newRegistry.getBalance(linkUpkeepID) + newRegistry.getBalance(linkUpkeepID2)
    );
    assertEq(usdToken18.balanceOf(address(newRegistry)), newRegistry.getBalance(usdUpkeepID18));
    assertEq(weth.balanceOf(address(newRegistry)), newRegistry.getBalance(nativeUpkeepID));
    assertEq(linkToken.balanceOf(address(registry)), prevTokenBalances[0] - linkToken.balanceOf(address(newRegistry)));
    assertEq(
      usdToken18.balanceOf(address(registry)),
      prevTokenBalances[1] - usdToken18.balanceOf(address(newRegistry))
    );
    assertEq(weth.balanceOf(address(registry)), prevTokenBalances[2] - weth.balanceOf(address(newRegistry)));

    // assert upkeep data matches
    assertEq(prevUpkeepData[0], abi.encode(newRegistry.getUpkeep(linkUpkeepID)));
    assertEq(prevUpkeepData[1], abi.encode(newRegistry.getUpkeep(linkUpkeepID2)));
    assertEq(prevUpkeepData[2], abi.encode(newRegistry.getUpkeep(usdUpkeepID18)));
    assertEq(prevUpkeepData[3], abi.encode(newRegistry.getUpkeep(nativeUpkeepID)));
    assertEq(prevUpkeepTriggerData[0], newRegistry.getUpkeepTriggerConfig(linkUpkeepID));
    assertEq(prevUpkeepTriggerData[1], newRegistry.getUpkeepTriggerConfig(linkUpkeepID2));
    assertEq(prevUpkeepTriggerData[2], newRegistry.getUpkeepTriggerConfig(usdUpkeepID18));
    assertEq(prevUpkeepTriggerData[3], newRegistry.getUpkeepTriggerConfig(nativeUpkeepID));

    vm.stopPrank();
  }
}

contract Pause is SetUp {
  function test_RevertsWhen_CalledByNonOwner() external {
    vm.expectRevert(bytes("Only callable by owner"));
    vm.prank(STRANGER);
    registry.pause();
  }

  function test_CalledByOwner_success() external {
    vm.startPrank(registry.owner());
    registry.pause();

    (IAutomationV21PlusCommon.StateLegacy memory state, , , , ) = registry.getState();
    assertTrue(state.paused);
  }

  function test_revertsWhen_transmitInPausedRegistry() external {
    vm.startPrank(registry.owner());
    registry.pause();

    _transmitAndExpectRevert(usdUpkeepID18, registry, Registry.RegistryPaused.selector);
  }
}

contract Unpause is SetUp {
  function test_RevertsWhen_CalledByNonOwner() external {
    vm.startPrank(registry.owner());
    registry.pause();

    vm.expectRevert(bytes("Only callable by owner"));
    vm.startPrank(STRANGER);
    registry.unpause();
  }

  function test_CalledByOwner_success() external {
    vm.startPrank(registry.owner());
    registry.pause();
    (IAutomationV21PlusCommon.StateLegacy memory state1, , , , ) = registry.getState();
    assertTrue(state1.paused);

    registry.unpause();
    (IAutomationV21PlusCommon.StateLegacy memory state2, , , , ) = registry.getState();
    assertFalse(state2.paused);
  }
}

contract CancelUpkeep is SetUp {
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);

  function test_RevertsWhen_IdIsInvalid_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);
    vm.expectRevert(Registry.CannotCancel.selector);
    registry.cancelUpkeep(1111111);
  }

  function test_RevertsWhen_IdIsInvalid_CalledByOwner() external {
    vm.startPrank(registry.owner());
    vm.expectRevert(Registry.CannotCancel.selector);
    registry.cancelUpkeep(1111111);
  }

  function test_RevertsWhen_NotCalledByOwnerOrAdmin() external {
    vm.startPrank(STRANGER);
    vm.expectRevert(Registry.OnlyCallableByOwnerOrAdmin.selector);
    registry.cancelUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceledByAdmin_CalledByOwner() external {
    uint256 bn = block.number;
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.startPrank(registry.owner());
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.cancelUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceledByOwner_CalledByAdmin() external {
    uint256 bn = block.number;
    vm.startPrank(registry.owner());
    registry.cancelUpkeep(linkUpkeepID);

    vm.startPrank(UPKEEP_ADMIN);
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.cancelUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceledByAdmin_CalledByAdmin() external {
    uint256 bn = block.number;
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.cancelUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceledByOwner_CalledByOwner() external {
    uint256 bn = block.number;
    vm.startPrank(registry.owner());
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.cancelUpkeep(linkUpkeepID);
  }

  function test_CancelUpkeep_SetMaxValidBlockNumber_CalledByAdmin() external {
    uint256 bn = block.number;
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    uint256 maxValidBlocknumber = uint256(registry.getUpkeep(linkUpkeepID).maxValidBlocknumber);

    // 50 is the cancellation delay
    assertEq(bn + 50, maxValidBlocknumber);
  }

  function test_CancelUpkeep_SetMaxValidBlockNumber_CalledByOwner() external {
    uint256 bn = block.number;
    vm.startPrank(registry.owner());
    registry.cancelUpkeep(linkUpkeepID);

    uint256 maxValidBlocknumber = uint256(registry.getUpkeep(linkUpkeepID).maxValidBlocknumber);

    // cancellation by registry owner is immediate and no cancellation delay is applied
    assertEq(bn, maxValidBlocknumber);
  }

  function test_CancelUpkeep_EmitEvent_CalledByAdmin() external {
    uint256 bn = block.number;
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectEmit();
    emit UpkeepCanceled(linkUpkeepID, uint64(bn + 50));
    registry.cancelUpkeep(linkUpkeepID);
  }

  function test_CancelUpkeep_EmitEvent_CalledByOwner() external {
    uint256 bn = block.number;
    vm.startPrank(registry.owner());

    vm.expectEmit();
    emit UpkeepCanceled(linkUpkeepID, uint64(bn));
    registry.cancelUpkeep(linkUpkeepID);
  }
}

contract SetPeerRegistryMigrationPermission is SetUp {
  function test_SetPeerRegistryMigrationPermission_CalledByOwner() external {
    address peer = randomAddress();
    vm.startPrank(registry.owner());

    uint8 permission = registry.getPeerRegistryMigrationPermission(peer);
    assertEq(0, permission);

    registry.setPeerRegistryMigrationPermission(peer, 1);
    permission = registry.getPeerRegistryMigrationPermission(peer);
    assertEq(1, permission);

    registry.setPeerRegistryMigrationPermission(peer, 2);
    permission = registry.getPeerRegistryMigrationPermission(peer);
    assertEq(2, permission);

    registry.setPeerRegistryMigrationPermission(peer, 0);
    permission = registry.getPeerRegistryMigrationPermission(peer);
    assertEq(0, permission);
  }

  function test_RevertsWhen_InvalidPermission_CalledByOwner() external {
    address peer = randomAddress();
    vm.startPrank(registry.owner());

    vm.expectRevert();
    registry.setPeerRegistryMigrationPermission(peer, 100);
  }

  function test_RevertsWhen_CalledByNonOwner() external {
    address peer = randomAddress();
    vm.startPrank(STRANGER);

    vm.expectRevert(bytes("Only callable by owner"));
    registry.setPeerRegistryMigrationPermission(peer, 1);
  }
}

contract SetUpkeepPrivilegeConfig is SetUp {
  function test_RevertsWhen_CalledByNonManager() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByUpkeepPrivilegeManager.selector);
    registry.setUpkeepPrivilegeConfig(linkUpkeepID, hex"1233");
  }

  function test_UpkeepHasEmptyConfig() external {
    bytes memory cfg = registry.getUpkeepPrivilegeConfig(linkUpkeepID);
    assertEq(cfg, bytes(""));
  }

  function test_SetUpkeepPrivilegeConfig_CalledByManager() external {
    vm.startPrank(PRIVILEGE_MANAGER);

    registry.setUpkeepPrivilegeConfig(linkUpkeepID, hex"1233");

    bytes memory cfg = registry.getUpkeepPrivilegeConfig(linkUpkeepID);
    assertEq(cfg, hex"1233");
  }
}

contract SetAdminPrivilegeConfig is SetUp {
  function test_RevertsWhen_CalledByNonManager() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByUpkeepPrivilegeManager.selector);
    registry.setAdminPrivilegeConfig(randomAddress(), hex"1233");
  }

  function test_UpkeepHasEmptyConfig() external {
    bytes memory cfg = registry.getAdminPrivilegeConfig(randomAddress());
    assertEq(cfg, bytes(""));
  }

  function test_SetAdminPrivilegeConfig_CalledByManager() external {
    vm.startPrank(PRIVILEGE_MANAGER);
    address admin = randomAddress();

    registry.setAdminPrivilegeConfig(admin, hex"1233");

    bytes memory cfg = registry.getAdminPrivilegeConfig(admin);
    assertEq(cfg, hex"1233");
  }
}

contract TransferUpkeepAdmin is SetUp {
  event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to);

  function test_RevertsWhen_NotCalledByAdmin() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.transferUpkeepAdmin(linkUpkeepID, randomAddress());
  }

  function test_RevertsWhen_TransferToSelf() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.ValueNotChanged.selector);
    registry.transferUpkeepAdmin(linkUpkeepID, UPKEEP_ADMIN);
  }

  function test_RevertsWhen_UpkeepCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);

    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.transferUpkeepAdmin(linkUpkeepID, randomAddress());
  }

  function test_DoesNotChangeUpkeepAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.transferUpkeepAdmin(linkUpkeepID, randomAddress());

    assertEq(registry.getUpkeep(linkUpkeepID).admin, UPKEEP_ADMIN);
  }

  function test_EmitEvent_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);
    address newAdmin = randomAddress();

    vm.expectEmit();
    emit UpkeepAdminTransferRequested(linkUpkeepID, UPKEEP_ADMIN, newAdmin);
    registry.transferUpkeepAdmin(linkUpkeepID, newAdmin);

    // transferring to the same propose admin won't yield another event
    vm.recordLogs();
    registry.transferUpkeepAdmin(linkUpkeepID, newAdmin);
    Vm.Log[] memory entries = vm.getRecordedLogs();
    assertEq(0, entries.length);
  }

  function test_CancelTransfer_ByTransferToEmptyAddress() external {
    vm.startPrank(UPKEEP_ADMIN);
    address newAdmin = randomAddress();

    vm.expectEmit();
    emit UpkeepAdminTransferRequested(linkUpkeepID, UPKEEP_ADMIN, newAdmin);
    registry.transferUpkeepAdmin(linkUpkeepID, newAdmin);

    vm.expectEmit();
    emit UpkeepAdminTransferRequested(linkUpkeepID, UPKEEP_ADMIN, address(0));
    registry.transferUpkeepAdmin(linkUpkeepID, address(0));
  }
}

contract AcceptUpkeepAdmin is SetUp {
  event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to);

  function test_RevertsWhen_NotCalledByProposedAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);
    address newAdmin = randomAddress();
    registry.transferUpkeepAdmin(linkUpkeepID, newAdmin);

    vm.startPrank(STRANGER);
    vm.expectRevert(Registry.OnlyCallableByProposedAdmin.selector);
    registry.acceptUpkeepAdmin(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    address newAdmin = randomAddress();
    registry.transferUpkeepAdmin(linkUpkeepID, newAdmin);

    registry.cancelUpkeep(linkUpkeepID);

    vm.startPrank(newAdmin);
    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.acceptUpkeepAdmin(linkUpkeepID);
  }

  function test_UpkeepAdminChanged() external {
    vm.startPrank(UPKEEP_ADMIN);
    address newAdmin = randomAddress();
    registry.transferUpkeepAdmin(linkUpkeepID, newAdmin);

    vm.startPrank(newAdmin);
    vm.expectEmit();
    emit UpkeepAdminTransferred(linkUpkeepID, UPKEEP_ADMIN, newAdmin);
    registry.acceptUpkeepAdmin(linkUpkeepID);

    assertEq(newAdmin, registry.getUpkeep(linkUpkeepID).admin);
  }
}

contract PauseUpkeep is SetUp {
  event UpkeepPaused(uint256 indexed id);

  function test_RevertsWhen_NotCalledByUpkeepAdmin() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.pauseUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_InvalidUpkeepId() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.pauseUpkeep(linkUpkeepID + 1);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.pauseUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepAlreadyPaused() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.pauseUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.OnlyUnpausedUpkeep.selector);
    registry.pauseUpkeep(linkUpkeepID);
  }

  function test_EmitEvent_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectEmit();
    emit UpkeepPaused(linkUpkeepID);
    registry.pauseUpkeep(linkUpkeepID);
  }
}

contract UnpauseUpkeep is SetUp {
  event UpkeepUnpaused(uint256 indexed id);

  function test_RevertsWhen_InvalidUpkeepId() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.unpauseUpkeep(linkUpkeepID + 1);
  }

  function test_RevertsWhen_UpkeepIsNotPaused() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyPausedUpkeep.selector);
    registry.unpauseUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.pauseUpkeep(linkUpkeepID);

    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.unpauseUpkeep(linkUpkeepID);
  }

  function test_RevertsWhen_NotCalledByUpkeepAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.pauseUpkeep(linkUpkeepID);

    vm.startPrank(STRANGER);
    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.unpauseUpkeep(linkUpkeepID);
  }

  function test_UnpauseUpkeep_CalledByUpkeepAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.pauseUpkeep(linkUpkeepID);

    uint256[] memory ids1 = registry.getActiveUpkeepIDs(0, 0);

    vm.expectEmit();
    emit UpkeepUnpaused(linkUpkeepID);
    registry.unpauseUpkeep(linkUpkeepID);

    uint256[] memory ids2 = registry.getActiveUpkeepIDs(0, 0);
    assertEq(ids1.length + 1, ids2.length);
  }
}

contract SetUpkeepCheckData is SetUp {
  event UpkeepCheckDataSet(uint256 indexed id, bytes newCheckData);

  function test_RevertsWhen_InvalidUpkeepId() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepCheckData(linkUpkeepID + 1, hex"1234");
  }

  function test_RevertsWhen_UpkeepAlreadyCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.setUpkeepCheckData(linkUpkeepID, hex"1234");
  }

  function test_RevertsWhen_NewCheckDataTooLarge() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.CheckDataExceedsLimit.selector);
    registry.setUpkeepCheckData(linkUpkeepID, new bytes(10_000));
  }

  function test_RevertsWhen_NotCalledByAdmin() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepCheckData(linkUpkeepID, hex"1234");
  }

  function test_UpdateOffchainConfig_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectEmit();
    emit UpkeepCheckDataSet(linkUpkeepID, hex"1234");
    registry.setUpkeepCheckData(linkUpkeepID, hex"1234");

    assertEq(registry.getUpkeep(linkUpkeepID).checkData, hex"1234");
  }

  function test_UpdateOffchainConfigOnPausedUpkeep_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);

    registry.pauseUpkeep(linkUpkeepID);

    vm.expectEmit();
    emit UpkeepCheckDataSet(linkUpkeepID, hex"1234");
    registry.setUpkeepCheckData(linkUpkeepID, hex"1234");

    assertTrue(registry.getUpkeep(linkUpkeepID).paused);
    assertEq(registry.getUpkeep(linkUpkeepID).checkData, hex"1234");
  }
}

contract SetUpkeepGasLimit is SetUp {
  event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit);

  function test_RevertsWhen_InvalidUpkeepId() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepGasLimit(linkUpkeepID + 1, 1230000);
  }

  function test_RevertsWhen_UpkeepAlreadyCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.setUpkeepGasLimit(linkUpkeepID, 1230000);
  }

  function test_RevertsWhen_NewGasLimitOutOfRange() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.GasLimitOutsideRange.selector);
    registry.setUpkeepGasLimit(linkUpkeepID, 300);

    vm.expectRevert(Registry.GasLimitOutsideRange.selector);
    registry.setUpkeepGasLimit(linkUpkeepID, 3000000000);
  }

  function test_RevertsWhen_NotCalledByAdmin() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepGasLimit(linkUpkeepID, 1230000);
  }

  function test_UpdateGasLimit_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectEmit();
    emit UpkeepGasLimitSet(linkUpkeepID, 1230000);
    registry.setUpkeepGasLimit(linkUpkeepID, 1230000);

    assertEq(registry.getUpkeep(linkUpkeepID).performGas, 1230000);
  }
}

contract SetUpkeepOffchainConfig is SetUp {
  event UpkeepOffchainConfigSet(uint256 indexed id, bytes offchainConfig);

  function test_RevertsWhen_InvalidUpkeepId() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepOffchainConfig(linkUpkeepID + 1, hex"1233");
  }

  function test_RevertsWhen_UpkeepAlreadyCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.setUpkeepOffchainConfig(linkUpkeepID, hex"1234");
  }

  function test_RevertsWhen_NotCalledByAdmin() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepOffchainConfig(linkUpkeepID, hex"1234");
  }

  function test_UpdateOffchainConfig_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectEmit();
    emit UpkeepOffchainConfigSet(linkUpkeepID, hex"1234");
    registry.setUpkeepOffchainConfig(linkUpkeepID, hex"1234");

    assertEq(registry.getUpkeep(linkUpkeepID).offchainConfig, hex"1234");
  }
}

contract SetUpkeepTriggerConfig is SetUp {
  event UpkeepTriggerConfigSet(uint256 indexed id, bytes triggerConfig);

  function test_RevertsWhen_InvalidUpkeepId() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepTriggerConfig(linkUpkeepID + 1, hex"1233");
  }

  function test_RevertsWhen_UpkeepAlreadyCanceled() external {
    vm.startPrank(UPKEEP_ADMIN);
    registry.cancelUpkeep(linkUpkeepID);

    vm.expectRevert(Registry.UpkeepCancelled.selector);
    registry.setUpkeepTriggerConfig(linkUpkeepID, hex"1234");
  }

  function test_RevertsWhen_NotCalledByAdmin() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByAdmin.selector);
    registry.setUpkeepTriggerConfig(linkUpkeepID, hex"1234");
  }

  function test_UpdateTriggerConfig_CalledByAdmin() external {
    vm.startPrank(UPKEEP_ADMIN);

    vm.expectEmit();
    emit UpkeepTriggerConfigSet(linkUpkeepID, hex"1234");
    registry.setUpkeepTriggerConfig(linkUpkeepID, hex"1234");

    assertEq(registry.getUpkeepTriggerConfig(linkUpkeepID), hex"1234");
  }
}

contract TransferPayeeship is SetUp {
  event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to);

  function test_RevertsWhen_NotCalledByPayee() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(Registry.OnlyCallableByPayee.selector);
    registry.transferPayeeship(TRANSMITTERS[0], randomAddress());
  }

  function test_RevertsWhen_TransferToSelf() external {
    registry.setPayees(PAYEES);
    vm.startPrank(PAYEES[0]);

    vm.expectRevert(Registry.ValueNotChanged.selector);
    registry.transferPayeeship(TRANSMITTERS[0], PAYEES[0]);
  }

  function test_Transfer_DoesNotChangePayee() external {
    registry.setPayees(PAYEES);

    vm.startPrank(PAYEES[0]);

    registry.transferPayeeship(TRANSMITTERS[0], randomAddress());

    (, , , , address payee) = registry.getTransmitterInfo(TRANSMITTERS[0]);
    assertEq(PAYEES[0], payee);
  }

  function test_EmitEvent_CalledByPayee() external {
    registry.setPayees(PAYEES);

    vm.startPrank(PAYEES[0]);
    address newPayee = randomAddress();

    vm.expectEmit();
    emit PayeeshipTransferRequested(TRANSMITTERS[0], PAYEES[0], newPayee);
    registry.transferPayeeship(TRANSMITTERS[0], newPayee);

    // transferring to the same propose payee won't yield another event
    vm.recordLogs();
    registry.transferPayeeship(TRANSMITTERS[0], newPayee);
    Vm.Log[] memory entries = vm.getRecordedLogs();
    assertEq(0, entries.length);
  }
}

contract AcceptPayeeship is SetUp {
  event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to);

  function test_RevertsWhen_NotCalledByProposedPayee() external {
    registry.setPayees(PAYEES);

    vm.startPrank(PAYEES[0]);
    address newPayee = randomAddress();
    registry.transferPayeeship(TRANSMITTERS[0], newPayee);

    vm.startPrank(STRANGER);
    vm.expectRevert(Registry.OnlyCallableByProposedPayee.selector);
    registry.acceptPayeeship(TRANSMITTERS[0]);
  }

  function test_PayeeChanged() external {
    registry.setPayees(PAYEES);

    vm.startPrank(PAYEES[0]);
    address newPayee = randomAddress();
    registry.transferPayeeship(TRANSMITTERS[0], newPayee);

    vm.startPrank(newPayee);
    vm.expectEmit();
    emit PayeeshipTransferred(TRANSMITTERS[0], PAYEES[0], newPayee);
    registry.acceptPayeeship(TRANSMITTERS[0]);

    (, , , , address payee) = registry.getTransmitterInfo(TRANSMITTERS[0]);
    assertEq(newPayee, payee);
  }
}

contract SetPayees is SetUp {
  event PayeesUpdated(address[] transmitters, address[] payees);

  function test_RevertsWhen_NotCalledByOwner() external {
    vm.startPrank(STRANGER);

    vm.expectRevert(bytes("Only callable by owner"));
    registry.setPayees(NEW_PAYEES);
  }

  function test_RevertsWhen_PayeesLengthError() external {
    vm.startPrank(registry.owner());

    address[] memory payees = new address[](5);
    vm.expectRevert(Registry.ParameterLengthError.selector);
    registry.setPayees(payees);
  }

  function test_RevertsWhen_InvalidPayee() external {
    vm.startPrank(registry.owner());

    NEW_PAYEES[0] = address(0);
    vm.expectRevert(Registry.InvalidPayee.selector);
    registry.setPayees(NEW_PAYEES);
  }

  function test_SetPayees_WhenExistingPayeesAreEmpty() external {
    (registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);

    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      (, , , , address payee) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      assertEq(address(0), payee);
    }

    vm.startPrank(registry.owner());

    vm.expectEmit();
    emit PayeesUpdated(TRANSMITTERS, PAYEES);
    registry.setPayees(PAYEES);
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      (bool active, , , , address payee) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      assertTrue(active);
      assertEq(PAYEES[i], payee);
    }
  }

  function test_DotNotSetPayeesToIgnoredAddress() external {
    address IGNORE_ADDRESS = 0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF;
    (registry, ) = deployAndConfigureRegistryAndRegistrar(AutoBase.PayoutMode.ON_CHAIN);
    PAYEES[0] = IGNORE_ADDRESS;

    registry.setPayees(PAYEES);
    (bool active, , , , address payee) = registry.getTransmitterInfo(TRANSMITTERS[0]);
    assertTrue(active);
    assertEq(address(0), payee);

    (active, , , , payee) = registry.getTransmitterInfo(TRANSMITTERS[1]);
    assertTrue(active);
    assertEq(PAYEES[1], payee);
  }
}

contract GetActiveUpkeepIDs is SetUp {
  function test_RevertsWhen_IndexOutOfRange() external {
    vm.expectRevert(Registry.IndexOutOfRange.selector);
    registry.getActiveUpkeepIDs(5, 0);

    vm.expectRevert(Registry.IndexOutOfRange.selector);
    registry.getActiveUpkeepIDs(6, 0);
  }

  function test_ReturnsAllUpkeeps_WhenMaxCountIsZero() external {
    uint256[] memory uids = registry.getActiveUpkeepIDs(0, 0);
    assertEq(5, uids.length);

    uids = registry.getActiveUpkeepIDs(2, 0);
    assertEq(3, uids.length);
  }

  function test_ReturnsAllRemainingUpkeeps_WhenMaxCountIsTooLarge() external {
    uint256[] memory uids = registry.getActiveUpkeepIDs(2, 20);
    assertEq(3, uids.length);
  }

  function test_ReturnsUpkeeps_BoundByMaxCount() external {
    uint256[] memory uids = registry.getActiveUpkeepIDs(1, 2);
    assertEq(2, uids.length);
    assertEq(uids[0], linkUpkeepID2);
    assertEq(uids[1], usdUpkeepID18);
  }
}
