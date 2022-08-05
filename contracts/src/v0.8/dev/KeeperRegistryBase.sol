pragma solidity 0.8.6;

import "./vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import "./vendor/@eth-optimism/contracts/0.8.6/contracts/L2/predeploys/OVM_GasPriceOracle.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "../ConfirmedOwner.sol";
import "./ExecutionPrevention.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/KeeperCompatibleInterface.sol";
import {RegistryConfig, State} from "./interfaces/KeeperRegistryInterfaceDev.sol";
import "./interfaces/OCR2Abstract.sol";
import "../interfaces/UpkeepTranscoderInterface.sol";

/**
 * @notice Base Keeper Registry contract, contains shared logic between
 * KeeperRegistry and KeeperRegistryLogic
 */
abstract contract KeeperRegistryBase is ConfirmedOwner, ExecutionPrevention, ReentrancyGuard, Pausable, OCR2Abstract {
  address internal constant ZERO_ADDRESS = address(0);
  address internal constant IGNORE_ADDRESS = 0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF;
  bytes4 internal constant CHECK_SELECTOR = KeeperCompatibleInterface.checkUpkeep.selector;
  bytes4 internal constant PERFORM_SELECTOR = KeeperCompatibleInterface.performUpkeep.selector;
  uint256 internal constant PERFORM_GAS_MIN = 2_300;
  uint256 internal constant CANCELLATION_DELAY = 50;
  uint256 internal constant PERFORM_GAS_CUSHION = 5_000;
  uint256 internal constant PPB_BASE = 1_000_000_000;
  uint32 internal constant UINT32_MAX = 2**32 - 1;
  uint96 internal constant LINK_TOTAL_SUPPLY = 1e27;
  UpkeepFormat internal constant UPKEEP_TRANSCODER_VESION_BASE = UpkeepFormat.V1;

  // L1_FEE_DATA_PADDING includes 35 bytes for L1 data padding for Optimism
  bytes public L1_FEE_DATA_PADDING = "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  // MAX_INPUT_DATA represents the estimated max size of the sum of L1 data padding and msg.data in performUpkeep
  // function, which includes 4 bytes for function selector, 32 bytes for upkeep id, 35 bytes for data padding, and
  // 64 bytes for estimated perform data
  bytes public MAX_INPUT_DATA =
    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";

  // Upkeep storage
  EnumerableSet.UintSet internal s_upkeepIDs;
  mapping(uint256 => Upkeep) internal s_upkeep;
  mapping(uint256 => bytes) internal s_checkData;

  // Proposed payees for keepers
  mapping(address => address) internal s_proposedPayee;
  mapping(address => MigrationPermission) internal s_peerRegistryMigrationPermission;

  // TODO: optimise EVM word storage
  RegistryConfig internal s_config;

  // OCR config storage
  mapping(address => Transmitter) internal s_transmitters;
  mapping(address => Signer) internal s_signers;
  address[] internal s_signersList; // s_signersList contains the signing address of each oracle
  address[] internal s_transmittersList; // s_transmittersList contains the transmission address of each oracle
  uint8 internal s_f; // Number of faulty oracles allowed
  uint64 s_offchainConfigVersion;
  bytes s_offchainConfig;

  // OCR state
  uint32 internal s_configCount; // incremented each time a new config is posted, The count
  // is incorporated into the config digest to prevent replay attacks.
  bytes32 internal s_latestConfigDigest;
  // makes it easier for offchain systems to extract config from logs
  uint32 internal s_latestConfigBlockNumber;

  // State of registry
  uint32 internal s_nonce; // Nonce for each upkeep created
  uint96 internal s_ownerLinkBalance;
  uint256 internal s_expectedLinkBalance;

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_NATIVE_FEED;
  AggregatorV3Interface public immutable FAST_GAS_FEED;
  OVM_GasPriceOracle public immutable OPTIMISM_ORACLE = OVM_GasPriceOracle(0x420000000000000000000000000000000000000F);
  ArbGasInfo public immutable ARB_NITRO_ORACLE = ArbGasInfo(0x000000000000000000000000000000000000006C);
  PaymentModel public immutable PAYMENT_MODEL;
  uint256 public immutable REGISTRY_GAS_OVERHEAD;

  error CannotCancel();
  error UpkeepCancelled();
  error MigrationNotPermitted();
  error UpkeepNotCanceled();
  error UpkeepNotNeeded();
  error NotAContract();
  error PaymentGreaterThanAllLINK();
  error OnlyActiveKeepers();
  error InsufficientFunds();
  error KeepersMustTakeTurns();
  error ParameterLengthError();
  error OnlyCallableByOwnerOrAdmin();
  error OnlyCallableByLINKToken();
  error InvalidPayee();
  error DuplicateEntry();
  error ValueNotChanged();
  error IndexOutOfRange();
  error TranscoderNotSet();
  error ArrayHasNoEntries();
  error GasLimitOutsideRange();
  error OnlyCallableByPayee();
  error OnlyCallableByProposedPayee();
  error GasLimitCanOnlyIncrease();
  error OnlyCallableByAdmin();
  error OnlyCallableByOwnerOrRegistrar();
  error InvalidRecipient();
  error InvalidDataLength();
  error OnlyUnpausedUpkeep();
  error OnlyPausedUpkeep();
  error ConfigDisgestMismatch();
  error IncorrectNumberOfSignatures();
  error OnlyActiveSigners();
  error DuplicateSigners();
  error StaleReport();
  error TooManyOracles();
  error IncorrectNumberOfSigners();
  error IncorrectNumberOfFaultyOracles();
  error OnchainConfigNonEmpty();
  error RepeatedSigner();
  error RepeatedTransmitter();

  enum MigrationPermission {
    NONE,
    OUTGOING,
    INCOMING,
    BIDIRECTIONAL
  }

  enum PaymentModel {
    DEFAULT,
    ARBITRUM,
    OPTIMISM
  }

  enum UpkeepFailureReason {
    NONE,
    TARGET_CHECK_REVERTED,
    UPKEEP_NOT_NEEDED,
    UPKEEP_PAUSED,
    INSUFFICIENT_BALANCE
  }

  // TODO: optimise EVM word storage
  struct Upkeep {
    uint96 balance;
    uint32 executeGas;
    uint32 maxValidBlocknumber;
    uint32 lastPerformBlockNumber;
    address target;
    uint96 amountSpent;
    address admin;
    bool paused;
  }

  struct PerformParams {
    uint256 id;
    bytes performData;
    uint256 fastGasWei;
    uint256 linkNativePrice;
    uint256 maxLinkPayment;
  }

  struct Transmitter {
    bool active;
    // Index of oracle in s_signersList/s_transmittersList
    uint8 index;
    uint96 balance;
    address payee;
  }

  struct Signer {
    bool active;
    // Index of oracle in s_signersList/s_transmittersList
    uint8 index;
  }

  event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin);
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    uint256 gasUsed,
    uint32 checkBlockNumber,
    uint96 payment
  );
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);
  event UpkeepPaused(uint256 indexed id);
  event UpkeepUnpaused(uint256 indexed id);
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);
  event OwnerFundsWithdrawn(uint96 amount);
  event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination);
  event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom);
  event RegistryConfigSet(RegistryConfig params);
  event KeepersUpdated(address[] keepers, address[] payees);
  event PaymentWithdrawn(address indexed keeper, uint256 indexed amount, address indexed to, address payee);
  event PayeeshipTransferRequested(address indexed keeper, address indexed from, address indexed to);
  event PayeeshipTransferred(address indexed keeper, address indexed from, address indexed to);
  event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit);

  /**
   * @param paymentModel the payment model of default, Arbitrum, or Optimism
   * @param registryGasOverhead the gas overhead used by registry in performUpkeep
   * @param link address of the LINK Token
   * @param linkNativeFeed address of the LINK/NATIVE price feed
   * @param fastGasFeed address of the Fast Gas price feed
   */
  constructor(
    PaymentModel paymentModel,
    uint256 registryGasOverhead,
    address link,
    address linkNativeFeed,
    address fastGasFeed
  ) ConfirmedOwner(msg.sender) {
    PAYMENT_MODEL = paymentModel;
    REGISTRY_GAS_OVERHEAD = registryGasOverhead;
    LINK = LinkTokenInterface(link);
    LINK_NATIVE_FEED = AggregatorV3Interface(linkNativeFeed);
    FAST_GAS_FEED = AggregatorV3Interface(fastGasFeed);
  }

  /**
   * @dev retrieves feed data for fast gas/eth. if the feed
   * data is stale it uses the configured fallback price.
   */
  function _getFasGasFeedData() internal view returns (uint256 gasWei) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 feedValue;
    (, feedValue, , timestamp, ) = FAST_GAS_FEED.latestRoundData();
    if ((staleFallback && stalenessSeconds < block.timestamp - timestamp) || feedValue <= 0) {
      gasWei = s_config.fallbackGasPrice;
    } else {
      gasWei = uint256(feedValue);
    }
    return gasWei;
  }

  /**
   * @dev retrieves feed data for link/native price. If the feed
   * data is stale it uses the configured fallback price.
   */
  function _getLinkNativeFeedData() internal view returns (uint256 linkEth) {
    uint32 stalenessSeconds = s_config.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 feedValue;
    (, feedValue, , timestamp, ) = LINK_NATIVE_FEED.latestRoundData();
    if ((staleFallback && stalenessSeconds < block.timestamp - timestamp) || feedValue <= 0) {
      linkEth = s_config.fallbackLinkPrice;
    } else {
      linkEth = uint256(feedValue);
    }
    return linkEth;
  }

  /**
   * @dev calculates LINK paid for gas spent plus a configure premium percentage
   * @param gasLimit the amount of gas used
   * @param fastGasWei the fast gas price
   * @param linkNativePrice the exchange ratio between LINK and Native token
   * @param isExecution if this is triggered by a perform upkeep function
   */
  function _calculatePaymentAmount(
    uint256 gasLimit,
    uint256 fastGasWei,
    uint256 linkNativePrice,
    bool isExecution
  )
    internal
    view
    returns (
      uint96 total,
      uint96 gasPayment,
      uint96 premium
    )
  {
    RegistryConfig memory config = s_config;
    uint256 gasWei = fastGasWei * config.gasCeilingMultiplier;
    // in case it's actual execution use actual gas price, capped by fastGasWei * gasCeilingMultiplier
    if (isExecution && tx.gasprice < gasWei) {
      gasWei = tx.gasprice;
    }

    uint256 weiForGas = gasWei * (gasLimit + REGISTRY_GAS_OVERHEAD);
    uint256 l1CostWei = 0;
    if (PAYMENT_MODEL == PaymentModel.OPTIMISM) {
      bytes memory txCallData = new bytes(0);
      if (isExecution) {
        txCallData = bytes.concat(msg.data, L1_FEE_DATA_PADDING);
      } else {
        txCallData = MAX_INPUT_DATA;
      }
      l1CostWei = OPTIMISM_ORACLE.getL1Fee(txCallData);
    } else if (PAYMENT_MODEL == PaymentModel.ARBITRUM) {
      l1CostWei = ARB_NITRO_ORACLE.getCurrentTxL1GasFees();
    }
    // if it's not performing upkeeps, use gas ceiling multiplier to estimate the upper bound
    if (!isExecution) {
      l1CostWei = config.gasCeilingMultiplier * l1CostWei;
    }

    uint256 gasPayment = ((weiForGas + l1CostWei) * 1e18) / linkNativePrice;
    uint256 premium = (gasPayment * config.paymentPremiumPPB) / 1e9 + uint256(config.flatFeeMicroLink) * 1e12;
    uint256 total = gasPayment + premium;
    // LINK_TOTAL_SUPPLY < UINT96_MAX
    if (total > LINK_TOTAL_SUPPLY) revert PaymentGreaterThanAllLINK();
    return (uint96(total), uint96(gasPayment), uint96(premium));
  }

  /**
   * @dev ensures the upkeep is not cancelled and the caller is the upkeep admin
   */
  function requireAdminAndNotCancelled(Upkeep memory upkeep) internal view {
    if (msg.sender != upkeep.admin) revert OnlyCallableByAdmin();
    if (upkeep.maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
  }

  /**
   * @dev generates a PerformParams struct for use in _performUpkeepWithParams()
   */
  function _generatePerformParams(
    uint256 id,
    bytes memory performData,
    bool isExecution, // Whether this is an actual perform execution or just a simulation
    uint256 executionLinkNativePrice // This price is used in case of execution, otherwise price is fetched from feed
  ) internal view returns (PerformParams memory) {
    Upkeep memory upkeep = s_upkeep[id];
    uint256 fastGasWei = _getFasGasFeedData();
    uint256 linkNativePrice;
    if (isExecution) {
      linkNativePrice = executionLinkNativePrice;
    } else {
      linkNativePrice = _getLinkNativeFeedData();
    }
    (uint96 maxLinkPayment, uint96 _gasPayment, uint96 _premium) = _calculatePaymentAmount(
      s_upkeep[id].executeGas,
      fastGasWei,
      linkNativePrice,
      isExecution
    );

    return
      PerformParams({
        id: id,
        performData: performData,
        fastGasWei: fastGasWei,
        linkNativePrice: linkNativePrice,
        maxLinkPayment: maxLinkPayment
      });
  }

  function _computeAndStoreConfigDigest(
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) internal {
    uint32 previousConfigBlockNumber = s_latestConfigBlockNumber;
    s_latestConfigBlockNumber = uint32(block.number);
    s_configCount += 1;

    s_latestConfigDigest = _configDigestFromConfigData(
      block.chainid,
      address(this),
      s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );

    emit ConfigSet(
      previousConfigBlockNumber,
      s_latestConfigDigest,
      s_configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  // MODIFIERS

  /**
   * @dev ensures a upkeep is valid
   */
  modifier validUpkeep(uint256 id) {
    if (s_upkeep[id].maxValidBlocknumber <= block.number) revert UpkeepCancelled();
    _;
  }

  /**
   * @dev Reverts if called by anyone other than the admin of upkeep #id
   */
  modifier onlyUpkeepAdmin(uint256 id) {
    if (msg.sender != s_upkeep[id].admin) revert OnlyCallableByAdmin();
    _;
  }

  /**
   * @dev Reverts if called on a cancelled upkeep
   */
  modifier onlyNonCanceledUpkeep(uint256 id) {
    if (s_upkeep[id].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
    _;
  }

  /**
   * @dev ensures that burns don't accidentally happen by sending to the zero
   * address
   */
  modifier validRecipient(address to) {
    if (to == ZERO_ADDRESS) revert InvalidRecipient();
    _;
  }

  /**
   * @dev Reverts if called by anyone other than the contract owner or registrar.
   */
  modifier onlyOwnerOrRegistrar() {
    if (msg.sender != owner() && msg.sender != s_config.registrar) revert OnlyCallableByOwnerOrRegistrar();
    _;
  }
}
