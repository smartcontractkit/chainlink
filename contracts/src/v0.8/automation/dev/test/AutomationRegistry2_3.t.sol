// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../v2_3/AutomationRegistryBase2_3.sol";
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
    vm.expectRevert(abi.encodeWithSelector(Registry.InvalidBillingToken.selector));
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
    assertEq(registry.getBalance(linkUpkeepID), 0);
    vm.prank(UPKEEP_ADMIN);
    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), 1);
    vm.prank(STRANGER);
    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), 2);
  }

  function test_movesFundFromCorrectToken() public {
    vm.startPrank(UPKEEP_ADMIN);

    uint256 startBalanceLINK = linkToken.balanceOf(address(registry));
    uint256 startBalanceUSDToken = usdToken.balanceOf(address(registry));

    registry.addFunds(linkUpkeepID, 1);
    assertEq(registry.getBalance(linkUpkeepID), 1);
    assertEq(registry.getBalance(usdUpkeepID), 0);
    assertEq(linkToken.balanceOf(address(registry)), startBalanceLINK + 1);
    assertEq(usdToken.balanceOf(address(registry)), startBalanceUSDToken);

    registry.addFunds(usdUpkeepID, 2);
    assertEq(registry.getBalance(linkUpkeepID), 1);
    assertEq(registry.getBalance(usdUpkeepID), 2);
    assertEq(linkToken.balanceOf(address(registry)), startBalanceLINK + 1);
    assertEq(usdToken.balanceOf(address(registry)), startBalanceUSDToken + 2);
  }

  function test_emitsAnEvent() public {
    vm.startPrank(UPKEEP_ADMIN);
    vm.expectEmit();
    emit FundsAdded(linkUpkeepID, address(UPKEEP_ADMIN), 100);
    registry.addFunds(linkUpkeepID, 100);
  }
}

contract Withdraw is SetUp {
  address internal aMockAddress = address(0x1111111111111111111111111111111111111113);

  function testLinkAvailableForPaymentReturnsLinkBalance() public {
    //simulate a deposit of link to the liquidity pool
    _mintLink(address(registry), 1e10);

    //check there's a balance
    assertGt(linkToken.balanceOf(address(registry)), 0);

    //check the link available for payment is the link balance
    assertEq(uint256(registry.linkAvailableForPayment()), linkToken.balanceOf(address(registry)));
  }

  function testWithdrawLinkFeesRevertsBecauseOnlyFinanceAdminAllowed() public {
    vm.expectRevert(abi.encodeWithSelector(Registry.OnlyFinanceAdmin.selector));
    registry.withdrawLinkFees(aMockAddress, 1);
  }

  function testWithdrawLinkFeesRevertsBecauseOfInsufficientBalance() public {
    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(Registry.InsufficientBalance.selector, 0, 1));
    registry.withdrawLinkFees(aMockAddress, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkFeesRevertsBecauseOfInvalidRecipient() public {
    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(Registry.InvalidRecipient.selector));
    registry.withdrawLinkFees(ZERO_ADDRESS, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkFeeSuccess() public {
    //simulate a deposit of link to the liquidity pool
    _mintLink(address(registry), 1e10);

    //check there's a balance
    assertGt(linkToken.balanceOf(address(registry)), 0);

    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is a ton of link available
    registry.withdrawLinkFees(aMockAddress, 1);

    vm.stopPrank();

    assertEq(linkToken.balanceOf(address(aMockAddress)), 1);
    assertEq(linkToken.balanceOf(address(registry)), 1e10 - 1);
  }

  function testWithdrawERC20FeeSuccess() public {
    // simulate a deposit of ERC20 to the liquidity pool
    _mintERC20(address(registry), 1e10);

    // check there's a balance
    assertGt(usdToken.balanceOf(address(registry)), 0);

    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is a ton of link available
    registry.withdrawERC20Fees(address(usdToken), aMockAddress, 1);

    vm.stopPrank();

    assertEq(usdToken.balanceOf(address(aMockAddress)), 1);
    assertEq(usdToken.balanceOf(address(registry)), 1e10 - 1);
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

  function testSetConfigRevertDueToInvalidBillingToken() public {
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

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);
    // deploy registry with OFF_CHAIN payout mode
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);

    vm.expectRevert(abi.encodeWithSelector(Registry.InvalidBillingToken.selector));
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
    AutoBase.Report memory report;
    {
      uint256[] memory upkeepIds = new uint256[](1);
      uint256[] memory gasLimits = new uint256[](1);
      bytes[] memory performDatas = new bytes[](1);
      bytes[] memory triggers = new bytes[](1);
      upkeepIds[0] = id;
      gasLimits[0] = 1000000;
      triggers[0] = _encodeConditionalTrigger(
        AutoBase.ConditionalTrigger(uint32(block.number - 1), blockhash(block.number - 1))
      );
      report = AutoBase.Report(uint256(1000000000), uint256(2000000000), upkeepIds, gasLimits, triggers, performDatas);
    }
    bytes memory reportBytes = _encodeReport(report);
    (, , bytes32 configDigest) = registry.latestConfigDetails();
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];
    uint256[] memory signerPKs = new uint256[](2);
    signerPKs[0] = SIGNING_KEY0;
    signerPKs[1] = SIGNING_KEY1;
    (bytes32[] memory rs, bytes32[] memory ss, bytes32 vs) = _signReport(reportBytes, reportContext, signerPKs);

    vm.startPrank(TRANSMITTERS[0]);
    registry.transmit(reportContext, reportBytes, rs, ss, vs);

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
      (bool active, uint8 index, uint96 balance, uint96 lastCollected, ) = registry.getTransmitterInfo(TRANSMITTERS[i]);
      assertTrue(active);
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
    vm.expectRevert(Registry.InvalidBillingToken.selector);
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
    vm.expectRevert(Registry.InvalidBillingToken.selector);
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
}
