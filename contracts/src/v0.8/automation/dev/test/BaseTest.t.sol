// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import "forge-std/Test.sol";

// import {LinkTokenInterface} from "../../../shared/interfaces/LinkTokenInterface.sol";
import {LinkToken} from "../../../shared/token/ERC677/LinkToken.sol";
import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {MockV3Aggregator} from "../../../tests/MockV3Aggregator.sol";
import {AutomationForwarderLogic} from "../../AutomationForwarderLogic.sol";
import {AutomationRegistry2_3} from "../v2_3/AutomationRegistry2_3.sol";
import {AutomationRegistryLogicA2_3} from "../v2_3/AutomationRegistryLogicA2_3.sol";
import {AutomationRegistryLogicB2_3} from "../v2_3/AutomationRegistryLogicB2_3.sol";
import {IAutomationRegistryMaster2_3, AutomationRegistryBase2_3} from "../interfaces/v2_3/IAutomationRegistryMaster2_3.sol";
import {AutomationRegistrar2_3} from "../v2_3/AutomationRegistrar2_3.sol";
import {ChainModuleBase} from "../../chains/ChainModuleBase.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

/**
 * @title BaseTest provides basic test setup procedures and dependancies for use by other
 * unit tests
 */
contract BaseTest is Test {
  // constants
  address internal constant ZERO_ADDRESS = address(0);

  // config
  uint8 internal constant F = 1; // number of faulty nodes
  uint64 internal constant OFFCHAIN_CONFIG_VERSION = 30; // 2 for OCR2

  // contracts
  LinkToken internal linkToken;
  ERC20Mock internal mockERC20;
  MockV3Aggregator internal LINK_USD_FEED;
  MockV3Aggregator internal NATIVE_USD_FEED;
  MockV3Aggregator internal USDTOKEN_USD_FEED;
  MockV3Aggregator internal FAST_GAS_FEED;

  // roles
  address internal constant OWNER = address(uint160(uint256(keccak256("OWNER"))));
  address internal constant UPKEEP_ADMIN = address(uint160(uint256(keccak256("UPKEEP_ADMIN"))));
  address internal constant FINANCE_ADMIN = address(uint160(uint256(keccak256("FINANCE_ADMIN"))));

  // nodes
  uint256 internal constant SIGNING_KEY0 = 0x7b2e97fe057e6de99d6872a2ef2abf52c9b4469bc848c2465ac3fcd8d336e81d;
  uint256 internal constant SIGNING_KEY1 = 0xab56160806b05ef1796789248e1d7f34a6465c5280899159d645218cd216cee6;
  uint256 internal constant SIGNING_KEY2 = 0x6ec7caa8406a49b76736602810e0a2871959fbbb675e23a8590839e4717f1f7f;
  uint256 internal constant SIGNING_KEY3 = 0x80f14b11da94ae7f29d9a7713ea13dc838e31960a5c0f2baf45ed458947b730a;
  address[] internal SIGNERS = new address[](4);
  address[] internal TRANSMITTERS = new address[](4);

  function setUp() public virtual {
    vm.startPrank(OWNER);
    linkToken = new LinkToken();
    linkToken.grantMintRole(OWNER);
    mockERC20 = new ERC20Mock("MOCK_ERC20", "MOCK_ERC20", OWNER, 0);

    LINK_USD_FEED = new MockV3Aggregator(8, 2_000_000_000); // $20
    NATIVE_USD_FEED = new MockV3Aggregator(8, 400_000_000_000); // $4,000
    USDTOKEN_USD_FEED = new MockV3Aggregator(8, 100_000_000); // $1
    FAST_GAS_FEED = new MockV3Aggregator(0, 1_000_000_000); // 1 gwei

    SIGNERS[0] = vm.addr(SIGNING_KEY0); //0xc110458BE52CaA6bB68E66969C3218A4D9Db0211
    SIGNERS[1] = vm.addr(SIGNING_KEY1); //0xc110a19c08f1da7F5FfB281dc93630923F8E3719
    SIGNERS[2] = vm.addr(SIGNING_KEY2); //0xc110fdF6e8fD679C7Cc11602d1cd829211A18e9b
    SIGNERS[3] = vm.addr(SIGNING_KEY3); //0xc11028017c9b445B6bF8aE7da951B5cC28B326C0

    TRANSMITTERS[0] = address(uint160(uint256(keccak256("TRANSMITTER1"))));
    TRANSMITTERS[1] = address(uint160(uint256(keccak256("TRANSMITTER2"))));
    TRANSMITTERS[2] = address(uint160(uint256(keccak256("TRANSMITTER3"))));
    TRANSMITTERS[3] = address(uint160(uint256(keccak256("TRANSMITTER4"))));

    vm.stopPrank();
  }

  /**
   * @notice deploys the component parts of a registry, but nothing more
   */
  function deployRegistry() internal returns (IAutomationRegistryMaster2_3) {
    AutomationForwarderLogic forwarderLogic = new AutomationForwarderLogic();
    AutomationRegistryLogicB2_3 logicB2_3 = new AutomationRegistryLogicB2_3(
      address(linkToken),
      address(LINK_USD_FEED),
      address(NATIVE_USD_FEED),
      address(FAST_GAS_FEED),
      address(forwarderLogic),
      ZERO_ADDRESS
    );
    AutomationRegistryLogicA2_3 logicA2_3 = new AutomationRegistryLogicA2_3(logicB2_3);
    return
      IAutomationRegistryMaster2_3(address(new AutomationRegistry2_3(AutomationRegistryLogicB2_3(address(logicA2_3))))); // wow this line is hilarious
  }

  /**
   * @notice deploys and configures a regisry, registrar, and everything needed for most tests
   */
  function deployAndConfigureAll() internal returns (IAutomationRegistryMaster2_3, AutomationRegistrar2_3) {
    IAutomationRegistryMaster2_3 registry = deployRegistry();
    // deploy & configure registrar
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
    IERC20[] memory billingTokens = new IERC20[](2);
    billingTokens[0] = IERC20(address(linkToken));
    billingTokens[1] = IERC20(address(mockERC20));
    uint256[] memory minRegistrationFees = new uint256[](billingTokens.length);
    minRegistrationFees[0] = 5000000000000000000; // 5 LINK
    minRegistrationFees[1] = 100000000000000000000; // 100 USD
    AutomationRegistrar2_3 registrar = new AutomationRegistrar2_3(
      address(linkToken),
      registry,
      triggerConfigs,
      billingTokens,
      minRegistrationFees
    );
    // configure registry
    address[] memory registrars = new address[](1);
    registrars[0] = address(registrar);
    address[] memory billingTokenAddresses = new address[](billingTokens.length);
    for (uint256 i = 0; i < billingTokens.length; i++) {
      billingTokenAddresses[i] = address(billingTokens[i]);
    }
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
      registrars: registrars,
      upkeepPrivilegeManager: 0xD9c855F08A7e460691F41bBDDe6eC310bc0593D8,
      chainModule: address(new ChainModuleBase()),
      reorgProtectionEnabled: true,
      financeAdmin: FINANCE_ADMIN
    });
    AutomationRegistryBase2_3.BillingConfig[]
      memory billingTokenConfigs = new AutomationRegistryBase2_3.BillingConfig[](2);
    billingTokenConfigs[0] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 10_000_000, // 10%
      flatFeeMicroLink: 100_000,
      priceFeed: address(LINK_USD_FEED),
      fallbackPrice: 1_000_000_000, // $10
      minSpend: 5000000000000000000 // 5 LINK
    });
    billingTokenConfigs[1] = AutomationRegistryBase2_3.BillingConfig({
      gasFeePPB: 10_000_000, // 15%
      flatFeeMicroLink: 100_000,
      priceFeed: address(USDTOKEN_USD_FEED),
      fallbackPrice: 100_000_000, // $1
      minSpend: 100000000000000000000 // 100 USD
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
}
