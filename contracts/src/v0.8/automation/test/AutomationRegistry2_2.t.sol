// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {AutomationForwarderLogic} from "../AutomationForwarderLogic.sol";
import {BaseTest} from "./BaseTest.t.sol";
import {AutomationRegistry2_2} from "../v2_2/AutomationRegistry2_2.sol";
import {AutomationRegistryBase2_2} from "../v2_2/AutomationRegistryBase2_2.sol";
import {AutomationRegistryLogicA2_2} from "../v2_2/AutomationRegistryLogicA2_2.sol";
import {AutomationRegistryLogicB2_2} from "../v2_2/AutomationRegistryLogicB2_2.sol";
import {IAutomationRegistryMaster} from "../interfaces/v2_2/IAutomationRegistryMaster.sol";
import {ChainModuleBase} from "../chains/ChainModuleBase.sol";

contract AutomationRegistry2_2_SetUp is BaseTest {
  address internal constant LINK_ETH_FEED = 0x1111111111111111111111111111111111111110;
  address internal constant FAST_GAS_FEED = 0x1111111111111111111111111111111111111112;
  address internal constant LINK_TOKEN = 0x1111111111111111111111111111111111111113;
  address internal constant ZERO_ADDRESS = address(0);

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

  IAutomationRegistryMaster internal registryMaster;

  function setUp() public override {
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
    AutomationRegistryLogicB2_2 logicB2_2 = new AutomationRegistryLogicB2_2(
      LINK_TOKEN,
      LINK_ETH_FEED,
      FAST_GAS_FEED,
      address(forwarderLogic),
      ZERO_ADDRESS
    );
    AutomationRegistryLogicA2_2 logicA2_2 = new AutomationRegistryLogicA2_2(logicB2_2);
    registryMaster = IAutomationRegistryMaster(
      address(new AutomationRegistry2_2(AutomationRegistryLogicB2_2(address(logicA2_2))))
    );
  }
}

contract AutomationRegistry2_2_LatestConfigDetails is AutomationRegistry2_2_SetUp {
  function testGet() public {
    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);
    assertEq(blockNumber, 0);
    assertEq(configDigest, "");
  }
}

contract AutomationRegistry2_2_CheckUpkeep is AutomationRegistry2_2_SetUp {
  function testPreventExecutionOnCheckUpkeep() public {
    uint256 id = 1;
    bytes memory triggerData = abi.encodePacked("trigger_data");

    // The tx.origin is the DEFAULT_SENDER (0x1804c8AB1F12E6bbf3894d4083f33e07309d1f38) of foundry
    // Expecting a revert since the tx.origin is not address(0)
    vm.expectRevert(abi.encodeWithSelector(IAutomationRegistryMaster.OnlySimulatedBackend.selector));
    registryMaster.checkUpkeep(id, triggerData);
  }
}

contract AutomationRegistry2_2_SetConfig is AutomationRegistry2_2_SetUp {
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

  function testSetConfigSuccess() public {
    (uint32 configCount, , ) = registryMaster.latestConfigDetails();
    assertEq(configCount, 0);
    ChainModuleBase module = new ChainModuleBase();

    AutomationRegistryBase2_2.OnchainConfig memory cfg = AutomationRegistryBase2_2.OnchainConfig({
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
      fallbackLinkPrice: 200_000_000_000,
      transcoder: 0xB1e66855FD67f6e85F0f0fA38cd6fBABdf00923c,
      registrars: s_registrars,
      upkeepPrivilegeManager: 0xD9c855F08A7e460691F41bBDDe6eC310bc0593D8,
      chainModule: module,
      reorgProtectionEnabled: true
    });
    bytes memory onchainConfigBytes = abi.encode(cfg);

    uint256 a = 1234;
    address b = address(0);
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
      onchainConfigBytes,
      OFFCHAIN_CONFIG_VERSION,
      offchainConfigBytes
    );

    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registryMaster.getState();

    assertEq(signers, s_valid_signers);
    assertEq(transmitters, s_valid_transmitters);
    assertEq(f, F);
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
