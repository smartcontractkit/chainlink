// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationForwarderLogic} from "../../AutomationForwarderLogic.sol";
import {BaseTest} from "./BaseTest.t.sol";
import {AutomationRegistry2_3} from "../v2_3/AutomationRegistry2_3.sol";
import {AutomationRegistryLogicA2_3} from "../v2_3/AutomationRegistryLogicA2_3.sol";
import {AutomationRegistryLogicB2_3} from "../v2_3/AutomationRegistryLogicB2_3.sol";
import {IAutomationRegistryMaster2_3, AutomationRegistryBase2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {ChainModuleBase} from "../../chains/ChainModuleBase.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";

contract AutomationRegistry2_3_SetUp is BaseTest {
  address internal LINK_USD_FEED;
  address internal NATIVE_USD_FEED;
  address internal FAST_GAS_FEED;
  address internal constant FINANCE_ADMIN_ADDRESS = 0x1111111111111111111111111111111111111114;
  address internal constant ZERO_ADDRESS = address(0);
  address internal constant UPKEEP_ADMIN = address(uint160(uint256(keccak256("ADMIN"))));

  // Signer private keys used for these test
  uint256 internal constant PRIVATE0 = 0x7b2e97fe057e6de99d6872a2ef2abf52c9b4469bc848c2465ac3fcd8d336e81d;
  uint256 internal constant PRIVATE1 = 0xab56160806b05ef1796789248e1d7f34a6465c5280899159d645218cd216cee6;
  uint256 internal constant PRIVATE2 = 0x6ec7caa8406a49b76736602810e0a2871959fbbb675e23a8590839e4717f1f7f;
  uint256 internal constant PRIVATE3 = 0x80f14b11da94ae7f29d9a7713ea13dc838e31960a5c0f2baf45ed458947b730a;

  uint64 internal constant OFFCHAIN_CONFIG_VERSION = 30; // 2 for OCR2
  uint8 internal constant F = 1;

  address[] internal s_valid_signers;
  address[] internal s_valid_transmitters;
  address[] internal s_registrars;

  IAutomationRegistryMaster2_3 internal registryMaster;
  ERC20Mock internal link; // the link token
  ERC20Mock internal mockERC20; // the supported ERC20 tokens except link

  function setUp() public override {
    LINK_USD_FEED = address(new MockV3Aggregator(8, 2_000_000_000)); // $20
    NATIVE_USD_FEED = address(new MockV3Aggregator(8, 400_000_000_000)); // $4,000
    FAST_GAS_FEED = address(new MockV3Aggregator(0, 1_000_000_000)); // 1 gwei

    link = new ERC20Mock("LINK", "LINK", UPKEEP_ADMIN, 0);
    mockERC20 = new ERC20Mock("MOCK_ERC20", "MOCK_ERC20", UPKEEP_ADMIN, 0);

    s_valid_transmitters = new address[](4);
    for (uint160 i = 0; i < 4; ++i) {
      s_valid_transmitters[i] = address(4 + i);
    }

    s_valid_signers = new address[](4);
    s_valid_signers[0] = vm.addr(PRIVATE0); //0xc110458BE52CaA6bB68E66969C3218A4D9Db0211
    s_valid_signers[1] = vm.addr(PRIVATE1); //0xc110a19c08f1da7F5FfB281dc93630923F8E3719
    s_valid_signers[2] = vm.addr(PRIVATE2); //0xc110fdF6e8fD679C7Cc11602d1cd829211A18e9b
    s_valid_signers[3] = vm.addr(PRIVATE3); //0xc11028017c9b445B6bF8aE7da951B5cC28B326C0

    s_registrars = new address[](1);
    s_registrars[0] = 0x3a0eDE26aa188BFE00b9A0C9A431A1a0CA5f7966;

    AutomationForwarderLogic forwarderLogic = new AutomationForwarderLogic();
    AutomationRegistryLogicB2_3 logicB2_3 = new AutomationRegistryLogicB2_3(
      address(link),
      LINK_USD_FEED,
      NATIVE_USD_FEED,
      FAST_GAS_FEED,
      address(forwarderLogic),
      ZERO_ADDRESS
    );
    AutomationRegistryLogicA2_3 logicA2_3 = new AutomationRegistryLogicA2_3(logicB2_3);
    registryMaster = IAutomationRegistryMaster2_3(
      address(new AutomationRegistry2_3(AutomationRegistryLogicB2_3(address(logicA2_3))))
    );
  }
}

contract AutomationRegistry2_3_LatestConfigDetails is AutomationRegistry2_3_SetUp {
  function testGet() public {
    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);
    assertEq(blockNumber, 0);
    assertEq(configDigest, "");
  }
}

contract AutomationRegistry2_3_CheckUpkeep is AutomationRegistry2_3_SetUp {
  function testPreventExecutionOnCheckUpkeep() public {
    uint256 id = 1;
    bytes memory triggerData = abi.encodePacked("trigger_data");

    // The tx.origin is the DEFAULT_SENDER (0x1804c8AB1F12E6bbf3894d4083f33e07309d1f38) of foundry
    // Expecting a revert since the tx.origin is not address(0)
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.OnlySimulatedBackend.selector));
    registryMaster.checkUpkeep(id, triggerData);
  }
}

contract AutomationRegistry2_3_Withdraw is AutomationRegistry2_3_SetUp {
  address internal aMockAddress = address(0x1111111111111111111111111111111111111113);

  function mintLink(address recipient, uint256 amount) public {
    vm.prank(UPKEEP_ADMIN);
    //mint the link to the recipient
    link.mint(recipient, amount);
  }

  function mintERC20(address recipient, uint256 amount) public {
    vm.prank(UPKEEP_ADMIN);
    //mint the ERC20 to the recipient
    mockERC20.mint(recipient, amount);
  }

  function setConfigForWithdraw() public {
    address module = address(new ChainModuleBase());
    AutomationRegistryBase2_3.OnchainConfig memory cfg = AutomationRegistryBase2_3.OnchainConfig({
      paymentPremiumPPB: 10_000,
      flatFeeMicroLink: 40_000,
      checkGasLimit: 5_000_000,
      stalenessSeconds: 90_000,
      gasCeilingMultiplier: 0,
      minUpkeepSpend: 0,
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
      financeAdmin: FINANCE_ADMIN_ADDRESS
    });
    bytes memory offchainConfigBytes = abi.encode(1234, ZERO_ADDRESS);

    registryMaster.setConfigTypeSafe(
      s_valid_signers,
      s_valid_transmitters,
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
    mintLink(address(registryMaster), 1e10);

    //check there's a balance
    assertGt(link.balanceOf(address(registryMaster)), 0);

    //check the link available for payment is the link balance
    assertEq(registryMaster.linkAvailableForPayment(), link.balanceOf(address(registryMaster)));
  }

  function testWithdrawLinkFeesRevertsBecauseOnlyFinanceAdminAllowed() public {
    // set config with the finance admin
    setConfigForWithdraw();

    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.OnlyFinanceAdmin.selector));
    registryMaster.withdrawLinkFees(aMockAddress, 1);
  }

  function testWithdrawLinkFeesRevertsBecauseOfInsufficientBalance() public {
    // set config with the finance admin
    setConfigForWithdraw();

    vm.startPrank(FINANCE_ADMIN_ADDRESS);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.InsufficientBalance.selector, 0, 1));
    registryMaster.withdrawLinkFees(aMockAddress, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkFeesRevertsBecauseOfInvalidRecipient() public {
    // set config with the finance admin
    setConfigForWithdraw();

    vm.startPrank(FINANCE_ADMIN_ADDRESS);

    // try to withdraw 1 link while there is 0 balance
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.InvalidRecipient.selector));
    registryMaster.withdrawLinkFees(ZERO_ADDRESS, 1);

    vm.stopPrank();
  }

  function testWithdrawLinkFeeSuccess() public {
    // set config with the finance admin
    setConfigForWithdraw();

    //simulate a deposit of link to the liquidity pool
    mintLink(address(registryMaster), 1e10);

    //check there's a balance
    assertGt(link.balanceOf(address(registryMaster)), 0);

    vm.startPrank(FINANCE_ADMIN_ADDRESS);

    // try to withdraw 1 link while there is a ton of link available
    registryMaster.withdrawLinkFees(aMockAddress, 1);

    vm.stopPrank();

    assertEq(link.balanceOf(address(aMockAddress)), 1);
    assertEq(link.balanceOf(address(registryMaster)), 1e10 - 1);
  }

  function testWithdrawERC20FeeSuccess() public {
    // set config with the finance admin
    setConfigForWithdraw();

    // simulate a deposit of ERC20 to the liquidity pool
    mintERC20(address(registryMaster), 1e10);

    // check there's a balance
    assertGt(mockERC20.balanceOf(address(registryMaster)), 0);

    vm.startPrank(FINANCE_ADMIN_ADDRESS);

    // try to withdraw 1 link while there is a ton of link available
    registryMaster.withdrawERC20Fees(address(mockERC20), aMockAddress, 1);

    vm.stopPrank();

    assertEq(mockERC20.balanceOf(address(aMockAddress)), 1);
    assertEq(mockERC20.balanceOf(address(registryMaster)), 1e10 - 1);
  }
}

contract AutomationRegistry2_3_SetConfig is AutomationRegistry2_3_SetUp {
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
      paymentPremiumPPB: 10_000,
      flatFeeMicroLink: 40_000,
      checkGasLimit: 5_000_000,
      stalenessSeconds: 90_000,
      gasCeilingMultiplier: 0,
      minUpkeepSpend: 0,
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
      financeAdmin: FINANCE_ADMIN_ADDRESS
    });

  function testSetConfigSuccess() public {
    (uint32 configCount, , ) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);

    address billingTokenAddress = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens = new address[](1);
    billingTokens[0] = billingTokenAddress;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_000,
      flatFeeMicroLink: 20_000,
      priceFeed: 0x2222222222222222222222222222222222222222
    });

    bytes memory onchainConfigBytes = abi.encode(cfg);
    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    uint256 a = 1234;
    address b = ZERO_ADDRESS;
    bytes memory offchainConfigBytes = abi.encode(a, b);
    bytes32 configDigest = _configDigestFromConfigData(
      block.chainid,
      address(registryMaster),
      ++configCount,
      s_valid_signers,
      s_valid_transmitters,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    vm.expectEmit();
    emit ConfigSet(
      0,
      configDigest,
      configCount,
      s_valid_signers,
      s_valid_transmitters,
      F,
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    registryMaster.setConfig(
      s_valid_signers,
      s_valid_transmitters,
      F,
      onchainConfigBytesWithBilling,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registryMaster.getState();

    assertEq(signers, s_valid_signers);
    assertEq(transmitters, s_valid_transmitters);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config = registryMaster.getBillingTokenConfig(billingTokenAddress);
    assertEq(config.gasFeePPB, 5_000);
    assertEq(config.flatFeeMicroLink, 20_000);
    assertEq(config.priceFeed, 0x2222222222222222222222222222222222222222);

    address[] memory tokens = registryMaster.getBillingTokens();
    assertEq(tokens.length, 1);
  }

  function testSetConfigMultipleBillingConfigsSuccess() public {
    (uint32 configCount, , ) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);

    address billingTokenAddress1 = address(0x1111111111111111111111111111111111111111);
    address billingTokenAddress2 = address(0x1111111111111111111111111111111111111112);
    address[] memory billingTokens = new address[](2);
    billingTokens[0] = billingTokenAddress1;
    billingTokens[1] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMicroLink: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMicroLink: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    uint256 a = 1234;
    address b = ZERO_ADDRESS;
    bytes memory offchainConfigBytes = abi.encode(a, b);

    registryMaster.setConfig(
      s_valid_signers,
      s_valid_transmitters,
      F,
      onchainConfigBytesWithBilling,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registryMaster.getState();

    assertEq(signers, s_valid_signers);
    assertEq(transmitters, s_valid_transmitters);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config1 = registryMaster.getBillingTokenConfig(billingTokenAddress1);
    assertEq(config1.gasFeePPB, 5_001);
    assertEq(config1.flatFeeMicroLink, 20_001);
    assertEq(config1.priceFeed, 0x2222222222222222222222222222222222222221);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registryMaster.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMicroLink, 20_002);
    assertEq(config2.priceFeed, 0x2222222222222222222222222222222222222222);

    address[] memory tokens = registryMaster.getBillingTokens();
    assertEq(tokens.length, 2);
  }

  function testSetConfigTwiceAndLastSetOverwrites() public {
    (uint32 configCount, , ) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);

    // BillingConfig1
    address billingTokenAddress1 = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens1 = new address[](1);
    billingTokens1[0] = billingTokenAddress1;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs1 = new AutomationRegistryBase2_3.BillingConfig[](1);
    billingConfigs1[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMicroLink: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221
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
      priceFeed: 0x2222222222222222222222222222222222222222
    });

    bytes memory onchainConfigBytesWithBilling2 = abi.encode(cfg, billingTokens2, billingConfigs2);

    uint256 a = 1234;
    address b = ZERO_ADDRESS;
    bytes memory offchainConfigBytes = abi.encode(a, b);

    // set config once
    registryMaster.setConfig(
      s_valid_signers,
      s_valid_transmitters,
      F,
      onchainConfigBytesWithBilling1,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    // set config twice
    registryMaster.setConfig(
      s_valid_signers,
      s_valid_transmitters,
      F,
      onchainConfigBytesWithBilling2,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registryMaster.getState();

    assertEq(signers, s_valid_signers);
    assertEq(transmitters, s_valid_transmitters);
    assertEq(f, F);

    AutomationRegistryBase2_3.BillingConfig memory config2 = registryMaster.getBillingTokenConfig(billingTokenAddress2);
    assertEq(config2.gasFeePPB, 5_002);
    assertEq(config2.flatFeeMicroLink, 20_002);
    assertEq(config2.priceFeed, 0x2222222222222222222222222222222222222222);

    address[] memory tokens = registryMaster.getBillingTokens();
    assertEq(tokens.length, 1);
  }

  function testSetConfigDuplicateBillingConfigFailure() public {
    (uint32 configCount, , ) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);

    address billingTokenAddress1 = address(0x1111111111111111111111111111111111111111);
    address billingTokenAddress2 = address(0x1111111111111111111111111111111111111111);
    address[] memory billingTokens = new address[](2);
    billingTokens[0] = billingTokenAddress1;
    billingTokens[1] = billingTokenAddress2;

    AutomationRegistryBase2_3.BillingConfig[] memory billingConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_001,
      flatFeeMicroLink: 20_001,
      priceFeed: 0x2222222222222222222222222222222222222221
    });
    billingConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 5_002,
      flatFeeMicroLink: 20_002,
      priceFeed: 0x2222222222222222222222222222222222222222
    });

    bytes memory onchainConfigBytesWithBilling = abi.encode(cfg, billingTokens, billingConfigs);

    uint256 a = 1234;
    address b = ZERO_ADDRESS;
    bytes memory offchainConfigBytes = abi.encode(a, b);

    // expect revert because of duplicate tokens
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster2_3.DuplicateEntry.selector));
    registryMaster.setConfig(
      s_valid_signers,
      s_valid_transmitters,
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
