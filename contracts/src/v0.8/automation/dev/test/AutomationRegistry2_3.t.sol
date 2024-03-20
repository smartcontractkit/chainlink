// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../v2_3/AutomationRegistryBase2_3.sol";
import {IAutomationRegistryMaster2_3, AutomationRegistryBase2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {ChainModuleBase} from "../../chains/ChainModuleBase.sol";
import {IAutomationV21PlusCommon} from "../../interfaces/IAutomationV21PlusCommon.sol";

// forge test --match-path src/v0.8/automation/dev/test/AutomationRegistry2_3.t.sol

contract SetUp is BaseTest {
  address[] internal s_registrars;

  IAutomationRegistryMaster2_3 internal registry;
  uint256[] internal upkeepIds;
  uint256[] internal gasLimits;
  bytes[] internal performDatas;
  uint256[] internal balances;

  function setUp() public virtual override {
    super.setUp();

    s_registrars = new address[](1);
    s_registrars[0] = 0x3a0eDE26aa188BFE00b9A0C9A431A1a0CA5f7966;

    (registry, ) = deployAndConfigureAll(AutoBase.PayoutMode.ON_CHAIN);
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
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.OnlySimulatedBackend.selector));
    registry.checkUpkeep(id, triggerData);
  }
}

contract Withdraw is SetUp {
  address internal aMockAddress = address(0x1111111111111111111111111111111111111113);

  function setConfigForWithdraw() public {
    address module = address(new ChainModuleBase());
    AutomationRegistryBase2_3.OnchainConfig memory cfg = AutomationRegistryBase2_3.OnchainConfig({
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
      registrars: s_registrars,
      upkeepPrivilegeManager: 0xD9c855F08A7e460691F41bBDDe6eC310bc0593D8,
      chainModule: module,
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });
    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

    registry.setConfigTypeSafe(
      SIGNERS,
      TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes,
      new address[](0),
      new AutomationRegistryBase2_3.BillingConfig[](0)
    );
  }

  function testLinkAvailableForPaymentReturnsLinkBalance() public {
    //simulate a deposit of link to the liquidity pool
    _mintLink(address(registry), 1e10);

    //check there's a balance
    assertGt(linkToken.balanceOf(address(registry)), 0);

    //check the link available for payment is the link balance
    assertEq(registry.linkAvailableForPayment(), linkToken.balanceOf(address(registry)));
  }

  function testWithdrawLinkFeesRevertsBecauseOnlyFinanceAdminAllowed() public {
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.OnlyFinanceAdmin.selector));
    registry.withdrawLinkFees(aMockAddress, 1);
  }

  function testWithdrawLinkFeesRevertsBecauseOfInsufficientBalance() public {
    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.InsufficientBalance.selector, 0, 1));
    registry.withdrawLinkFees(aMockAddress, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkFeesRevertsBecauseOfInvalidRecipient() public {
    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.InvalidRecipient.selector));
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
    assertGt(mockERC20.balanceOf(address(registry)), 0);

    vm.startPrank(FINANCE_ADMIN);

    // try to withdraw 1 link while there is a ton of link available
    registry.withdrawERC20Fees(address(mockERC20), aMockAddress, 1);

    vm.stopPrank();

    assertEq(mockERC20.balanceOf(address(aMockAddress)), 1);
    assertEq(mockERC20.balanceOf(address(registry)), 1e10 - 1);
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
      registrars: s_registrars,
      upkeepPrivilegeManager: 0xD9c855F08A7e460691F41bBDDe6eC310bc0593D8,
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
      flatFeeMicroLink: 20_000,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000
    });

    bytes memory onchainConfigBytes = abi.encode(cfg);
    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);
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
    assertEq(config.flatFeeMicroLink, 20_000);
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
      flatFeeMicroLink: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221,
      fallbackPrice: 100,
      minSpend: 100
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMicroLink: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 200,
      minSpend: 200
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

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
    assertEq(config1.flatFeeMicroLink, 20_001);
    assertEq(config1.priceFeed, 0x2222222222222222222222222222222222222221);
    assertEq(config1.fallbackPrice, 100);
    assertEq(config1.minSpend, 100);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registry.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMicroLink, 20_002);
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
      flatFeeMicroLink: 20_001,
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
      flatFeeMicroLink: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 200,
      minSpend: 200
    });

    bytes memory onchainConfigBytesWithBilling2 = abi.encode(cfg, billingTokens2, billingConfigs2);

    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

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
    assertEq(config2.flatFeeMicroLink, 20_002);
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
      flatFeeMicroLink: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221,
      fallbackPrice: 100,
      minSpend: 100
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMicroLink: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 200,
      minSpend: 200
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

    // expect revert because of duplicate tokens
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.DuplicateEntry.selector));
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
      flatFeeMicroLink: 20_000,
      priceFeed: 0x2222222222222222222222222222222222222222,
      fallbackPrice: 2_000_000_000, // $20
      minSpend: 100_000
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);
    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);
    // deploy registry with OFF_CHAIN payout mode
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);

    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.InvalidBillingToken.selector));
    registry.setConfig(
      SIGNERS,
      TRANSMITTERS,
      F,
      onchainConfigBytesWithBilling,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
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
  event NOPsSettledOffchain(address[] transmitterList, uint256[] balances);

  function deployAndSetConfigForSettleOffchain(AutoBase.PayoutMode payoutMode) public {
    registry = deployRegistry(payoutMode);
    address module = address(new ChainModuleBase());
    AutomationRegistryBase2_3.OnchainConfig memory cfg = AutomationRegistryBase2_3.OnchainConfig({
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
      registrars: s_registrars,
      upkeepPrivilegeManager: 0xD9c855F08A7e460691F41bBDDe6eC310bc0593D8,
      chainModule: module,
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });
    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

    registry.setConfigTypeSafe(
      SIGNERS,
      TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes,
      new address[](0),
      new AutomationRegistryBase2_3.BillingConfig[](0)
    );
  }

  function testSettleNOPsOffchainRevertDueToUnauthorizedCaller() public {
    deployAndSetConfigForSettleOffchain(AutoBase.PayoutMode.ON_CHAIN);

    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.OnlyFinanceAdmin.selector));
    registry.settleNOPsOffchain();
  }

  function testSettleNOPsOffchainRevertDueToOffchainSettlementDisabled() public {
    deployAndSetConfigForSettleOffchain(AutoBase.PayoutMode.ON_CHAIN);

    vm.prank(registry.owner());
    registry.disableOffchainPayments();

    vm.prank(FINANCE_ADMIN);
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.MustSettleOnchain.selector));
    registry.settleNOPsOffchain();
  }

  function testSettleNOPsOffchainSuccess() public {
    deployAndSetConfigForSettleOffchain(AutoBase.PayoutMode.OFF_CHAIN);

    uint256[] memory balances = new uint256[](TRANSMITTERS.length);
    for (uint256 i = 0; i < TRANSMITTERS.length; i++) {
      balances[i] = 0;
    }

    vm.startPrank(FINANCE_ADMIN);
    vm.expectEmit();
    emit NOPsSettledOffchain(TRANSMITTERS, balances);
    registry.settleNOPsOffchain();
  }

  function testSettleNOPsOffchainSuccessTransmitterBalanceZeroed() public {
    // deploy and configure a registry with OFF_CHAIN payout
    (IAutomationRegistryMaster2_3 registry, ) = deployAndConfigureAll(AutoBase.PayoutMode.OFF_CHAIN);

    // register an upkeep and add funds
    uint256 id = registry.registerUpkeep(address(TARGET1), 1000000, UPKEEP_ADMIN, 0, address(mockERC20), "", "", "");
    _mintERC20(UPKEEP_ADMIN, 1e20);
    vm.startPrank(UPKEEP_ADMIN);
    mockERC20.approve(address(registry), 1e20);
    registry.addFunds(id, 1e20);

    // manually create a transmit so transmitters earn some rewards
    upkeepIds = new uint256[](1);
    gasLimits = new uint256[](1);
    performDatas = new bytes[](1);
    bytes[] memory triggers = new bytes[](1);
    upkeepIds[0] = id;
    gasLimits[0] = 1000000;
    triggers[0] = _encodeConditionalTrigger(
      AutoBase.ConditionalTrigger(uint32(block.number - 1), blockhash(block.number - 1))
    );
    AutoBase.Report memory report = AutoBase.Report(
      uint256(1000000000),
      uint256(2000000000),
      upkeepIds,
      gasLimits,
      triggers,
      performDatas
    );
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
    (bool active, uint8 index, uint96 balance, uint96 lastCollected, ) = registry.getTransmitterInfo(TRANSMITTERS[1]);
    assertTrue(active);
    assertEq(1, index);
    assertTrue(balance > 0);
    assertEq(0, lastCollected);

    balances = new uint256[](TRANSMITTERS.length);
    for (uint256 i = 0; i < balances.length; i++) {
      balances[i] = balance;
    }

    // verify offchain settlement will emit NOPs' balances
    vm.startPrank(FINANCE_ADMIN);
    vm.expectEmit();
    emit NOPsSettledOffchain(TRANSMITTERS, balances);
    registry.settleNOPsOffchain();

    // verify that transmitters balance has been zeroed out
    (active, index, balance, , ) = registry.getTransmitterInfo(TRANSMITTERS[2]);
    assertTrue(active);
    assertEq(2, index);
    assertEq(0, balance);
  }

  function testDisableOffchainPaymentsRevertDueToUnauthorizedCaller() public {
    deployAndSetConfigForSettleOffchain(AutoBase.PayoutMode.OFF_CHAIN);

    vm.startPrank(FINANCE_ADMIN);
    vm.expectRevert(bytes("Only callable by owner"));
    registry.disableOffchainPayments();
  }

  function testDisableOffchainPaymentsSuccess() public {
    deployAndSetConfigForSettleOffchain(AutoBase.PayoutMode.OFF_CHAIN);

    vm.startPrank(registry.owner());
    registry.disableOffchainPayments();

    assertEq(uint8(AutoBase.PayoutMode.ON_CHAIN), registry.getPayoutMode());
  }
}

contract WithdrawPayment is SetUp {
  function testWithdrawPaymentRevertDueToOffchainPayoutMode() public {
    registry = deployRegistry(AutoBase.PayoutMode.OFF_CHAIN);
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.MustSettleOffchain.selector));
    vm.prank(TRANSMITTERS[0]);
    registry.withdrawPayment(TRANSMITTERS[0], TRANSMITTERS[0]);
  }
}
