// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "forge-std/Test.sol";

import {LinkToken} from "../../../shared/token/ERC677/LinkToken.sol";
import {ERC20Mock6Decimals} from "../../mocks/ERC20Mock6Decimals.sol";
import {MockV3Aggregator} from "../../../shared/mocks/MockV3Aggregator.sol";
import {AutomationForwarderLogic} from "../../AutomationForwarderLogic.sol";
import {UpkeepTranscoder5_0 as Transcoder} from "../../v2_3/UpkeepTranscoder5_0.sol";
import {AutomationRegistry2_3} from "../../v2_3/AutomationRegistry2_3.sol";
import {AutomationRegistryBase2_3 as AutoBase} from "../../v2_3/AutomationRegistryBase2_3.sol";
import {AutomationRegistryLogicA2_3} from "../../v2_3/AutomationRegistryLogicA2_3.sol";
import {AutomationRegistryLogicB2_3} from "../../v2_3/AutomationRegistryLogicB2_3.sol";
import {AutomationRegistryLogicC2_3} from "../../v2_3/AutomationRegistryLogicC2_3.sol";
import {IAutomationRegistryMaster2_3 as Registry, AutomationRegistryBase2_3} from "../../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {AutomationRegistrar2_3} from "../../v2_3/AutomationRegistrar2_3.sol";
import {ChainModuleBase} from "../../chains/ChainModuleBase.sol";
import {MockUpkeep} from "../../mocks/MockUpkeep.sol";
import {IWrappedNative} from "../../interfaces/v2_3/IWrappedNative.sol";

import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {IERC20Metadata as IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import {WETH9} from "../../../vendor/canonical-weth/WETH9.sol";

/**
 * @title BaseTest provides basic test setup procedures and dependencies for use by other
 * unit tests
 */
contract BaseTest is Test {
  // test state (not exposed to derived tests)
  uint256 private nonce;

  // constants
  address internal constant ZERO_ADDRESS = address(0);
  uint32 internal constant DEFAULT_GAS_FEE_PPB = 10_000_000;
  uint24 internal constant DEFAULT_FLAT_FEE_MILLI_CENTS = 2_000;

  // config
  uint8 internal constant F = 1; // number of faulty nodes
  uint64 internal constant OFFCHAIN_CONFIG_VERSION = 30; // 2 for OCR2

  // contracts
  LinkToken internal linkToken;
  ERC20Mock6Decimals internal usdToken6;
  ERC20Mock internal usdToken18;
  ERC20Mock internal usdToken18_2;
  WETH9 internal weth;
  MockV3Aggregator internal LINK_USD_FEED;
  MockV3Aggregator internal NATIVE_USD_FEED;
  MockV3Aggregator internal USDTOKEN_USD_FEED;
  MockV3Aggregator internal FAST_GAS_FEED;
  MockUpkeep internal TARGET1;
  MockUpkeep internal TARGET2;
  Transcoder internal TRANSCODER;

  // roles
  address internal constant OWNER = address(uint160(uint256(keccak256("OWNER"))));
  address internal constant UPKEEP_ADMIN = address(uint160(uint256(keccak256("UPKEEP_ADMIN"))));
  address internal constant FINANCE_ADMIN = address(uint160(uint256(keccak256("FINANCE_ADMIN"))));
  address internal constant STRANGER = address(uint160(uint256(keccak256("STRANGER"))));
  address internal constant BROKE_USER = address(uint160(uint256(keccak256("BROKE_USER")))); // do not mint to this address
  address internal constant PRIVILEGE_MANAGER = address(uint160(uint256(keccak256("PRIVILEGE_MANAGER"))));

  // nodes
  uint256 internal constant SIGNING_KEY0 = 0x7b2e97fe057e6de99d6872a2ef2abf52c9b4469bc848c2465ac3fcd8d336e81d;
  uint256 internal constant SIGNING_KEY1 = 0xab56160806b05ef1796789248e1d7f34a6465c5280899159d645218cd216cee6;
  uint256 internal constant SIGNING_KEY2 = 0x6ec7caa8406a49b76736602810e0a2871959fbbb675e23a8590839e4717f1f7f;
  uint256 internal constant SIGNING_KEY3 = 0x80f14b11da94ae7f29d9a7713ea13dc838e31960a5c0f2baf45ed458947b730a;
  address[] internal SIGNERS = new address[](4);
  address[] internal TRANSMITTERS = new address[](4);
  address[] internal NEW_TRANSMITTERS = new address[](4);
  address[] internal PAYEES = new address[](4);
  address[] internal NEW_PAYEES = new address[](4);

  function setUp() public virtual {
    vm.startPrank(OWNER);
    linkToken = new LinkToken();
    linkToken.grantMintRole(OWNER);
    usdToken18 = new ERC20Mock("MOCK_ERC20_18Decimals", "MOCK_ERC20_18Decimals", OWNER, 0);
    usdToken18_2 = new ERC20Mock("Second_MOCK_ERC20_18Decimals", "Second_MOCK_ERC20_18Decimals", OWNER, 0);
    usdToken6 = new ERC20Mock6Decimals("MOCK_ERC20_6Decimals", "MOCK_ERC20_6Decimals", OWNER, 0);
    weth = new WETH9();

    LINK_USD_FEED = new MockV3Aggregator(8, 2_000_000_000); // $20
    NATIVE_USD_FEED = new MockV3Aggregator(8, 400_000_000_000); // $4,000
    USDTOKEN_USD_FEED = new MockV3Aggregator(8, 100_000_000); // $1
    FAST_GAS_FEED = new MockV3Aggregator(0, 1_000_000_000); // 1 gwei

    TARGET1 = new MockUpkeep();
    TARGET2 = new MockUpkeep();

    TRANSCODER = new Transcoder();

    SIGNERS[0] = vm.addr(SIGNING_KEY0); //0xc110458BE52CaA6bB68E66969C3218A4D9Db0211
    SIGNERS[1] = vm.addr(SIGNING_KEY1); //0xc110a19c08f1da7F5FfB281dc93630923F8E3719
    SIGNERS[2] = vm.addr(SIGNING_KEY2); //0xc110fdF6e8fD679C7Cc11602d1cd829211A18e9b
    SIGNERS[3] = vm.addr(SIGNING_KEY3); //0xc11028017c9b445B6bF8aE7da951B5cC28B326C0

    TRANSMITTERS[0] = address(uint160(uint256(keccak256("TRANSMITTER1"))));
    TRANSMITTERS[1] = address(uint160(uint256(keccak256("TRANSMITTER2"))));
    TRANSMITTERS[2] = address(uint160(uint256(keccak256("TRANSMITTER3"))));
    TRANSMITTERS[3] = address(uint160(uint256(keccak256("TRANSMITTER4"))));
    NEW_TRANSMITTERS[0] = address(uint160(uint256(keccak256("TRANSMITTER1"))));
    NEW_TRANSMITTERS[1] = address(uint160(uint256(keccak256("TRANSMITTER2"))));
    NEW_TRANSMITTERS[2] = address(uint160(uint256(keccak256("TRANSMITTER5"))));
    NEW_TRANSMITTERS[3] = address(uint160(uint256(keccak256("TRANSMITTER6"))));

    PAYEES[0] = address(100);
    PAYEES[1] = address(101);
    PAYEES[2] = address(102);
    PAYEES[3] = address(103);
    NEW_PAYEES[0] = address(100);
    NEW_PAYEES[1] = address(101);
    NEW_PAYEES[2] = address(106);
    NEW_PAYEES[3] = address(107);

    // mint funds
    vm.deal(OWNER, 100 ether);
    vm.deal(UPKEEP_ADMIN, 100 ether);
    vm.deal(FINANCE_ADMIN, 100 ether);
    vm.deal(STRANGER, 100 ether);

    linkToken.mint(OWNER, 1000e18);
    linkToken.mint(UPKEEP_ADMIN, 1000e18);
    linkToken.mint(FINANCE_ADMIN, 1000e18);
    linkToken.mint(STRANGER, 1000e18);

    usdToken18.mint(OWNER, 1000e18);
    usdToken18.mint(UPKEEP_ADMIN, 1000e18);
    usdToken18.mint(FINANCE_ADMIN, 1000e18);
    usdToken18.mint(STRANGER, 1000e18);

    usdToken18_2.mint(UPKEEP_ADMIN, 1000e18);

    usdToken6.mint(OWNER, 1000e6);
    usdToken6.mint(UPKEEP_ADMIN, 1000e6);
    usdToken6.mint(FINANCE_ADMIN, 1000e6);
    usdToken6.mint(STRANGER, 1000e6);

    deal(address(weth), OWNER, 100 ether);
    deal(address(weth), UPKEEP_ADMIN, 100 ether);
    deal(address(weth), FINANCE_ADMIN, 100 ether);
    deal(address(weth), STRANGER, 100 ether);

    vm.stopPrank();
  }

  /// @notice deploys the component parts of a registry, but nothing more
  function deployRegistry(AutoBase.PayoutMode payoutMode) internal returns (Registry) {
    AutomationForwarderLogic forwarderLogic = new AutomationForwarderLogic();
    AutomationRegistryLogicC2_3 logicC2_3 = new AutomationRegistryLogicC2_3(
      address(linkToken),
      address(LINK_USD_FEED),
      address(NATIVE_USD_FEED),
      address(FAST_GAS_FEED),
      address(forwarderLogic),
      ZERO_ADDRESS,
      payoutMode,
      address(weth)
    );
    AutomationRegistryLogicB2_3 logicB2_3 = new AutomationRegistryLogicB2_3(logicC2_3);
    AutomationRegistryLogicA2_3 logicA2_3 = new AutomationRegistryLogicA2_3(logicB2_3);
    return Registry(payable(address(new AutomationRegistry2_3(logicA2_3))));
  }

  /// @notice deploys and configures a registry, registrar, and everything needed for most tests
  function deployAndConfigureRegistryAndRegistrar(
    AutoBase.PayoutMode payoutMode
  ) internal returns (Registry, AutomationRegistrar2_3) {
    Registry registry = deployRegistry(payoutMode);

    IERC20[] memory billingTokens = new IERC20[](4);
    billingTokens[0] = IERC20(address(usdToken18));
    billingTokens[1] = IERC20(address(weth));
    billingTokens[2] = IERC20(address(linkToken));
    billingTokens[3] = IERC20(address(usdToken6));
    uint256[] memory minRegistrationFees = new uint256[](billingTokens.length);
    minRegistrationFees[0] = 100e18; // 100 USD
    minRegistrationFees[1] = 5e18; // 5 Native
    minRegistrationFees[2] = 5e18; // 5 LINK
    minRegistrationFees[3] = 100e6; // 100 USD
    address[] memory billingTokenAddresses = new address[](billingTokens.length);
    for (uint256 i = 0; i < billingTokens.length; i++) {
      billingTokenAddresses[i] = address(billingTokens[i]);
    }
    AutomationRegistryBase2_3.BillingConfig[]
      memory billingTokenConfigs = new AutomationRegistryBase2_3.BillingConfig[](billingTokens.length);
    billingTokenConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: DEFAULT_GAS_FEE_PPB, // 15%
      flatFeeMilliCents: DEFAULT_FLAT_FEE_MILLI_CENTS, // 2 cents
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 100_000_000, // $1
      minSpend: 1e18, // 1 USD
      decimals: 18
    });
    billingTokenConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: DEFAULT_GAS_FEE_PPB, // 15%
      flatFeeMilliCents: DEFAULT_FLAT_FEE_MILLI_CENTS, // 2 cents
      priceFeed: address(NATIVE_USD_FEED),
      fallbackPrice: 100_000_000, // $1
      minSpend: 5e18, // 5 Native
      decimals: 18
    });
    billingTokenConfigs[2] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: DEFAULT_GAS_FEE_PPB, // 10%
      flatFeeMilliCents: DEFAULT_FLAT_FEE_MILLI_CENTS, // 2 cents
      priceFeed: address(LINK_USD_FEED),
      fallbackPrice: 1_000_000_000, // $10
      minSpend: 1e18, // 1 LINK
      decimals: 18
    });
    billingTokenConfigs[3] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: DEFAULT_GAS_FEE_PPB, // 15%
      flatFeeMilliCents: DEFAULT_FLAT_FEE_MILLI_CENTS, // 2 cents
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 1e8, // $1
      minSpend: 1e6, // 1 USD
      decimals: 6
    });

    if (payoutMode == AutoBase.PayoutMode.OFF_CHAIN) {
      // remove LINK as a payment method if we are settling offchain
      assembly {
        mstore(billingTokens, 2)
        mstore(minRegistrationFees, 2)
        mstore(billingTokenAddresses, 2)
        mstore(billingTokenConfigs, 2)
      }
    }

    // deploy registrar
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
    AutomationRegistrar2_3 registrar = new AutomationRegistrar2_3(
      address(linkToken),
      registry,
      triggerConfigs,
      billingTokens,
      minRegistrationFees,
      IWrappedNative(address(weth))
    );

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
      transcoder: address(TRANSCODER),
      registrars: registrars,
      upkeepPrivilegeManager: PRIVILEGE_MANAGER,
      chainModule: address(new ChainModuleBase()),
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });

    registry.setConfigTypeSafe(
      SIGNERS,
      TRANSMITTERS,
      F,
      cfg,
      OFFCHAIN_CONFIG_VERSION,
      "",
      billingTokenAddresses,
      billingTokenConfigs
    );
    return (registry, registrar);
  }

  /// @notice this function updates the billing config for the provided token on the provided registry,
  /// and throws an error if the token is not found
  function _updateBillingTokenConfig(
    Registry registry,
    address billingToken,
    AutomationRegistryBase2_3.BillingConfig memory newConfig
  ) internal {
    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registry.getState();
    AutomationRegistryBase2_3.OnchainConfig memory config = registry.getConfig();
    address[] memory billingTokens = registry.getBillingTokens();

    AutomationRegistryBase2_3.BillingConfig[]
      memory billingTokenConfigs = new AutomationRegistryBase2_3.BillingConfig[](billingTokens.length);

    bool found = false;
    for (uint256 i = 0; i < billingTokens.length; i++) {
      if (billingTokens[i] == billingToken) {
        found = true;
        billingTokenConfigs[i] = newConfig;
      } else {
        billingTokenConfigs[i] = registry.getBillingTokenConfig(billingTokens[i]);
      }
    }
    require(found, "could not find billing token provided on registry");

    registry.setConfigTypeSafe(
      signers,
      transmitters,
      f,
      config,
      OFFCHAIN_CONFIG_VERSION,
      "",
      billingTokens,
      billingTokenConfigs
    );
  }

  /// @notice this function removes a billing token from the registry
  function _removeBillingTokenConfig(Registry registry, address billingToken) internal {
    (, , address[] memory signers, address[] memory transmitters, uint8 f) = registry.getState();
    AutomationRegistryBase2_3.OnchainConfig memory config = registry.getConfig();
    address[] memory billingTokens = registry.getBillingTokens();

    address[] memory newBillingTokens = new address[](billingTokens.length - 1);
    AutomationRegistryBase2_3.BillingConfig[]
      memory billingTokenConfigs = new AutomationRegistryBase2_3.BillingConfig[](billingTokens.length - 1);

    uint256 j = 0;
    for (uint256 i = 0; i < billingTokens.length; i++) {
      if (billingTokens[i] != billingToken) {
        if (j == newBillingTokens.length) revert("could not find billing token provided on registry");
        newBillingTokens[j] = billingTokens[i];
        billingTokenConfigs[j] = registry.getBillingTokenConfig(billingTokens[i]);
        j++;
      }
    }

    registry.setConfigTypeSafe(
      signers,
      transmitters,
      f,
      config,
      OFFCHAIN_CONFIG_VERSION,
      "",
      newBillingTokens,
      billingTokenConfigs
    );
  }

  // tests single upkeep, expects success
  function _transmit(uint256 id, Registry registry) internal {
    uint256[] memory ids = new uint256[](1);
    ids[0] = id;
    _handleTransmit(ids, registry, bytes4(0));
  }

  // tests multiple upkeeps, expects success
  function _transmit(uint256[] memory ids, Registry registry) internal {
    _handleTransmit(ids, registry, bytes4(0));
  }

  function _batchTransmitWithBadTriggers(uint256[] memory ids, bool[] memory goodTriggers, Registry registry) internal {
    bytes memory reportBytes;
    {
      uint256[] memory upkeepIds = new uint256[](ids.length);
      uint256[] memory gasLimits = new uint256[](ids.length);
      bytes[] memory performDatas = new bytes[](ids.length);
      bytes[] memory triggers = new bytes[](ids.length);
      for (uint256 i = 0; i < ids.length; i++) {
        upkeepIds[i] = ids[i];
        gasLimits[i] = registry.getUpkeep(ids[i]).performGas;
        performDatas[i] = new bytes(0);
        uint8 triggerType = registry.getTriggerType(ids[i]);
        if (triggerType == 0) {
          uint256 j = 0;
          if (goodTriggers[i]) {
            j = 1;
          }
          triggers[i] = _encodeConditionalTrigger(
            AutoBase.ConditionalTrigger(uint32(block.number - j), blockhash(block.number - j))
          );
        } else {
          revert("not implemented");
        }
      }

      AutoBase.Report memory report = AutoBase.Report(
        uint256(1000000000),
        uint256(2000000000),
        upkeepIds,
        gasLimits,
        triggers,
        performDatas
      );

      reportBytes = _encodeReport(report);
    }
    (, , bytes32 configDigest) = registry.latestConfigDetails();
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];
    uint256[] memory signerPKs = new uint256[](2);
    signerPKs[0] = SIGNING_KEY0;
    signerPKs[1] = SIGNING_KEY1;
    (bytes32[] memory rs, bytes32[] memory ss, bytes32 vs) = _signReport(reportBytes, reportContext, signerPKs);

    vm.startPrank(TRANSMITTERS[0]);
    registry.transmit(reportContext, reportBytes, rs, ss, vs);
    vm.stopPrank();
  }

  // tests single upkeep, expects revert
  function _transmitAndExpectRevert(uint256 id, Registry registry, bytes4 selector) internal {
    uint256[] memory ids = new uint256[](1);
    ids[0] = id;
    _handleTransmit(ids, registry, selector);
  }

  // private function not exposed to actual testing contract
  function _handleTransmit(uint256[] memory ids, Registry registry, bytes4 selector) private {
    bytes memory reportBytes;
    {
      uint256[] memory upkeepIds = new uint256[](ids.length);
      uint256[] memory gasLimits = new uint256[](ids.length);
      bytes[] memory performDatas = new bytes[](ids.length);
      bytes[] memory triggers = new bytes[](ids.length);
      for (uint256 i = 0; i < ids.length; i++) {
        upkeepIds[i] = ids[i];
        gasLimits[i] = registry.getUpkeep(ids[i]).performGas;
        performDatas[i] = new bytes(0);
        uint8 triggerType = registry.getTriggerType(ids[i]);
        if (triggerType == 0) {
          triggers[i] = _encodeConditionalTrigger(
            AutoBase.ConditionalTrigger(uint32(block.number - 1), blockhash(block.number - 1))
          );
        } else {
          revert("not implemented");
        }
      }

      AutoBase.Report memory report = AutoBase.Report(
        uint256(1000000000),
        uint256(2000000000),
        upkeepIds,
        gasLimits,
        triggers,
        performDatas
      );

      reportBytes = _encodeReport(report);
    }
    (, , bytes32 configDigest) = registry.latestConfigDetails();
    bytes32[3] memory reportContext = [configDigest, configDigest, configDigest];
    uint256[] memory signerPKs = new uint256[](2);
    signerPKs[0] = SIGNING_KEY0;
    signerPKs[1] = SIGNING_KEY1;
    (bytes32[] memory rs, bytes32[] memory ss, bytes32 vs) = _signReport(reportBytes, reportContext, signerPKs);

    vm.startPrank(TRANSMITTERS[0]);
    if (selector != bytes4(0)) {
      vm.expectRevert(selector);
    }
    registry.transmit(reportContext, reportBytes, rs, ss, vs);
    vm.stopPrank();
  }

  /// @notice Gather signatures on report data
  /// @param report - Report bytes generated from `_buildReport`
  /// @param reportContext - Report context bytes32 generated from `_buildReport`
  /// @param signerPrivateKeys - One or more addresses that will sign the report data
  /// @return rawRs - Signature rs
  /// @return rawSs - Signature ss
  /// @return rawVs - Signature vs
  function _signReport(
    bytes memory report,
    bytes32[3] memory reportContext,
    uint256[] memory signerPrivateKeys
  ) internal pure returns (bytes32[] memory, bytes32[] memory, bytes32) {
    bytes32[] memory rs = new bytes32[](signerPrivateKeys.length);
    bytes32[] memory ss = new bytes32[](signerPrivateKeys.length);
    bytes memory vs = new bytes(signerPrivateKeys.length);

    bytes32 reportDigest = keccak256(abi.encodePacked(keccak256(report), reportContext));

    for (uint256 i = 0; i < signerPrivateKeys.length; i++) {
      (uint8 v, bytes32 r, bytes32 s) = vm.sign(signerPrivateKeys[i], reportDigest);
      rs[i] = r;
      ss[i] = s;
      vs[i] = bytes1(v - 27);
    }

    return (rs, ss, bytes32(vs));
  }

  function _encodeReport(AutoBase.Report memory report) internal pure returns (bytes memory reportBytes) {
    return abi.encode(report);
  }

  function _encodeConditionalTrigger(
    AutoBase.ConditionalTrigger memory trigger
  ) internal pure returns (bytes memory triggerBytes) {
    return abi.encode(trigger.blockNum, trigger.blockHash);
  }

  /// @dev mints LINK to the recipient
  function _mintLink(address recipient, uint256 amount) internal {
    vm.prank(OWNER);
    linkToken.mint(recipient, amount);
  }

  /// @dev mints USDToken with 18 decimals to the recipient
  function _mintERC20_18Decimals(address recipient, uint256 amount) internal {
    vm.prank(OWNER);
    usdToken18.mint(recipient, amount);
  }

  /// @dev returns a pseudo-random 32 bytes
  function _random() private returns (bytes32) {
    nonce++;
    return keccak256(abi.encode(block.timestamp, nonce));
  }

  /// @dev returns a pseudo-random number
  function randomNumber() internal returns (uint256) {
    return uint256(_random());
  }

  /// @dev returns a pseudo-random address
  function randomAddress() internal returns (address) {
    return address(uint160(randomNumber()));
  }

  /// @dev returns a pseudo-random byte array
  function randomBytes(uint256 length) internal returns (bytes memory) {
    bytes memory result = new bytes(length);
    bytes32 entropy;
    for (uint256 i = 0; i < length; i++) {
      if (i % 32 == 0) {
        entropy = _random();
      }
      result[i] = entropy[i % 32];
    }
    return result;
  }
}
