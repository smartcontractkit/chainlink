// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import "../vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol";
import "../vendor/@eth-optimism/contracts/0.8.6/contracts/L2/predeploys/OVM_GasPriceOracle.sol";
import "../ExecutionPrevention.sol";
import {OnChainConfig, State, UpkeepFailureReason} from "./interfaces/KeeperRegistryInterface2_0.sol";
import "../../ConfirmedOwner.sol";
import "../../interfaces/AggregatorV3Interface.sol";
import "../../interfaces/LinkTokenInterface.sol";
import "../../interfaces/KeeperCompatibleInterface.sol";
import "../../interfaces/UpkeepTranscoderInterface.sol";

/**
 * @notice Base Keeper Registry contract, contains shared logic between
 * KeeperRegistry and KeeperRegistryLogic
 */
abstract contract KeeperRegistryBase2_0 is ConfirmedOwner, ExecutionPrevention, ReentrancyGuard, Pausable {
  address internal constant ZERO_ADDRESS = address(0);
  address internal constant IGNORE_ADDRESS = 0xFFfFfFffFFfffFFfFFfFFFFFffFFFffffFfFFFfF;
  bytes4 internal constant CHECK_SELECTOR = KeeperCompatibleInterface.checkUpkeep.selector;
  bytes4 internal constant PERFORM_SELECTOR = KeeperCompatibleInterface.performUpkeep.selector;
  uint256 internal constant PERFORM_GAS_MIN = 2_300;
  uint256 internal constant CANCELLATION_DELAY = 50;
  uint256 internal constant PERFORM_GAS_CUSHION = 5_000;
  uint256 internal constant PPB_BASE = 1_000_000_000;
  uint32 internal constant UINT32_MAX = type(uint32).max;
  uint96 internal constant LINK_TOTAL_SUPPLY = 1e27;
  UpkeepFormat internal constant UPKEEP_TRANSCODER_VERSION_BASE = UpkeepFormat.V3;

  // L1_FEE_DATA_PADDING includes 35 bytes for L1 data padding for Optimism
  bytes public L1_FEE_DATA_PADDING = "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";
  // MAX_INPUT_DATA represents the estimated max size of the sum of L1 data padding and msg.data in performUpkeep
  // function, which includes 4 bytes for function selector, 32 bytes for upkeep id, 35 bytes for data padding, and
  // 64 bytes for estimated perform data
  bytes public MAX_INPUT_DATA =
    "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff";

  // @dev - The storage is gas optimised for one and only function - transmit. All the storage accessed in transmit
  // is stored compactly. Rest of the storage layout is not of much concern as transmit is the only hot path

  // Upkeep storage
  EnumerableSet.UintSet internal s_upkeepIDs;
  mapping(uint256 => Upkeep) internal s_upkeep; // accessed during transmit
  mapping(uint256 => address) internal s_upkeepAdmin;
  mapping(uint256 => address) internal s_proposedAdmin;
  mapping(uint256 => bytes) internal s_checkData;

  // Registry config and state
  mapping(address => Transmitter) internal s_transmitters;
  mapping(address => Signer) internal s_signers;
  address[] internal s_signersList; // s_signersList contains the signing address of each oracle
  address[] internal s_transmittersList; // s_transmittersList contains the transmission address of each oracle
  mapping(address => address) internal s_transmitterPayees; // s_payees contains the mapping from transmitter to payee.
  mapping(address => address) internal s_proposedPayee; // proposed payee for a transmitter
  HotVars internal s_hotVars; // Mixture of config and state, used in transmit
  Storage internal s_storage; // Mixture of config and state, not used in transmit
  uint64 s_offchainConfigVersion;
  bytes s_offchainConfig;
  mapping(address => MigrationPermission) internal s_peerRegistryMigrationPermission;

  LinkTokenInterface public immutable LINK;
  AggregatorV3Interface public immutable LINK_NATIVE_FEED;
  AggregatorV3Interface public immutable FAST_GAS_FEED;
  OVM_GasPriceOracle public immutable OPTIMISM_ORACLE = OVM_GasPriceOracle(0x420000000000000000000000000000000000000F);
  ArbGasInfo public immutable ARB_NITRO_ORACLE = ArbGasInfo(0x000000000000000000000000000000000000006C);
  PaymentModel public immutable PAYMENT_MODEL;
  uint256 public immutable REGISTRY_GAS_OVERHEAD;

  error ArrayHasNoEntries();
  error CannotCancel();
  error DuplicateEntry();
  error GasLimitCanOnlyIncrease();
  error GasLimitOutsideRange();
  error IndexOutOfRange();
  error InsufficientFunds();
  error InvalidDataLength();
  error InvalidPayee();
  error InvalidRecipient();
  error MigrationNotPermitted();
  error NotAContract();
  error OnlyActiveTransmitters();
  error OnlyCallableByAdmin();
  error OnlyCallableByLINKToken();
  error OnlyCallableByOwnerOrAdmin();
  error OnlyCallableByOwnerOrRegistrar();
  error OnlyCallableByPayee();
  error OnlyCallableByProposedAdmin();
  error OnlyCallableByProposedPayee();
  error OnlyPausedUpkeep();
  error OnlyUnpausedUpkeep();
  error ParameterLengthError();
  error PaymentGreaterThanAllLINK();
  error TargetCheckReverted(bytes reason);
  error TranscoderNotSet();
  error UpkeepCancelled();
  error UpkeepNotCanceled();
  error UpkeepNotNeeded();
  error ValueNotChanged();
  error ConfigDisgestMismatch();
  error IncorrectNumberOfSignatures();
  error OnlyActiveSigners();
  error DuplicateSigners();
  error StaleReport();
  error ReorgedReport();
  error TooManyOracles();
  error IncorrectNumberOfSigners();
  error IncorrectNumberOfFaultyOracles();
  error RepeatedSigner();
  error RepeatedTransmitter();
  error OnchainConfigNonEmpty();
  error CheckDataExceedsLimit();
  error MaxCheckDataSizeCanOnlyIncrease();
  error MaxPerformDataSizeCanOnlyIncrease();

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

  struct PerformPaymentParams {
    uint256 fastGasWei;
    uint256 linkNative;
    uint256 maxLinkPayment;
  }

  // Config + State storage struct which is on hot transmit path
  struct HotVars {
    uint8 f; // maximum number of faulty oracles
    bytes32 latestConfigDigest; // latest config digest which is checked against every report
    uint32 paymentPremiumPPB; // premium percentage charged to user over tx cost
    uint32 flatFeeMicroLink; // flat fee charged to user for every perform
    uint24 stalenessSeconds; // Staleness tolerance for feeds
    uint16 gasCeilingMultiplier; // multiplier on top of fast gas feed for upper bound
    // 14 bytes to 1 EVM word
  }

  // Config + State storage struct which is not on hot transmit path
  struct Storage {
    uint256 fallbackGasPrice; // Used in case feed is stale
    // 1 EVM word full
    uint256 fallbackLinkPrice; // Used in case feed is stale
    // 2 EVM word full
    uint96 minUpkeepSpend; // Minimum amount an upkeep must spend
    address transcoder; // Address of transcoder contract used in migrations
    // 3 EVM word full
    uint96 ownerLinkBalance; // Balance of owner, accumulates minUpkeepSpend in case it is not spent
    address registrar; // Address of registrar used to register upkeeps
    // 4 EVM word full
    uint256 expectedLinkBalance; // Used in case of erroneous LINK transfers to contract
    // 5 EVM word full
    uint32 checkGasLimit; // Gas limit allowed in checkUpkeep
    uint32 maxPerformGas; // Max gas an upkeep can use on this registry
    uint32 nonce; // Nonce for each upkeep created
    uint32 configCount; // incremented each time a new config is posted, The count
    // is incorporated into the config digest to prevent replay attacks.
    uint32 latestConfigBlockNumber; // makes it easier for offchain systems to extract config from logs
    uint32 maxCheckDataSize; // max length of checkData bytes
    uint32 maxPerformDataSize; // max length of performData bytes
    // 4 bytes to 6th EVM word
  }

  struct Transmitter {
    bool active;
    // Index of oracle in s_signersList/s_transmittersList
    uint8 index;
    uint96 balance;
  }

  struct Signer {
    bool active;
    // Index of oracle in s_signersList/s_transmittersList
    uint8 index;
  }

  // This struct is used to pack information about the user's check function
  struct PerformDataWrapper {
    uint32 checkBlockNumber; // Block number on which check was called
    bytes32 checkBlockhash; // blockhash of checkBlockNumber-1. Used for reorg protection
    bytes performData; // actual performData that user's check returned
  }

  /**
   * @notice relevant state of an upkeep which is used in transmit function
   * @member balance the balance of this upkeep
   * @member target the contract which needs to be serviced
   * @member amountSpent the amount this upkeep has spent
   * @member executeGas the gas limit of upkeep execution
   * @member maxValidBlocknumber until which block this upkeep is valid
   * @member lastPerformBlockNumber the last block number when this upkeep was performed
   * @member paused if this upkeep has been paused
   */
  struct Upkeep {
    uint96 balance;
    address target;
    // 1 full EVM word
    uint96 amountSpent;
    uint32 executeGas;
    uint32 maxValidBlocknumber;
    uint32 lastPerformBlockNumber;
    bool paused;
    // 7 bytes left in 2nd EVM word
  }

  event OnChainConfigSet(OnChainConfig config);
  event FundsAdded(uint256 indexed id, address indexed from, uint96 amount);
  event FundsWithdrawn(uint256 indexed id, uint256 amount, address to);
  event OwnerFundsWithdrawn(uint96 amount);
  event PayeesUpdated(address[] transmitters, address[] payees);
  event PayeeshipTransferRequested(address indexed transmitter, address indexed from, address indexed to);
  event PayeeshipTransferred(address indexed transmitter, address indexed from, address indexed to);
  event PaymentWithdrawn(address indexed transmitter, uint256 indexed amount, address indexed to, address payee);
  event UpkeepAdminTransferRequested(uint256 indexed id, address indexed from, address indexed to);
  event UpkeepAdminTransferred(uint256 indexed id, address indexed from, address indexed to);
  event UpkeepCanceled(uint256 indexed id, uint64 indexed atBlockHeight);
  event UpkeepCheckDataUpdated(uint256 indexed id, bytes newCheckData);
  event UpkeepGasLimitSet(uint256 indexed id, uint96 gasLimit);
  event UpkeepMigrated(uint256 indexed id, uint256 remainingBalance, address destination);
  event UpkeepPaused(uint256 indexed id);
  event UpkeepPerformed(
    uint256 indexed id,
    bool indexed success,
    uint32 checkBlockNumber,
    uint256 gasUsed,
    uint256 linkNative,
    uint96 gasPayment,
    uint96 totalPayment
  );
  event UpkeepReceived(uint256 indexed id, uint256 startingBalance, address importedFrom);
  event UpkeepUnpaused(uint256 indexed id);
  event UpkeepRegistered(uint256 indexed id, uint32 executeGas, address admin);

  /**
   * @param paymentModel the payment model of default, Arbitrum, or Optimism
   * @param registryGasOverhead the gas overhead used by registry in performUpkeep
   * @param link address of the LINK Token
   * @param linkNativeFeed address of the LINK/Native price feed
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
   * @dev retrieves feed data for fast gas/native and link/native prices. if the feed
   * data is stale it uses the configured fallback price. Once a price is picked
   * for gas it takes the min of gas price in the transaction or the fast gas
   * price in order to reduce costs for the upkeep clients.
   */
  function _getFeedData(HotVars memory hotVars) internal view returns (uint256 gasWei, uint256 linkNative) {
    uint32 stalenessSeconds = hotVars.stalenessSeconds;
    bool staleFallback = stalenessSeconds > 0;
    uint256 timestamp;
    int256 feedValue;
    (, feedValue, , timestamp, ) = FAST_GAS_FEED.latestRoundData();
    if (
      feedValue <= 0 || block.timestamp < timestamp || (staleFallback && stalenessSeconds < block.timestamp - timestamp)
    ) {
      gasWei = s_storage.fallbackGasPrice;
    } else {
      gasWei = uint256(feedValue);
    }
    (, feedValue, , timestamp, ) = LINK_NATIVE_FEED.latestRoundData();
    if (
      feedValue <= 0 || block.timestamp < timestamp || (staleFallback && stalenessSeconds < block.timestamp - timestamp)
    ) {
      linkNative = s_storage.fallbackLinkPrice;
    } else {
      linkNative = uint256(feedValue);
    }
    return (gasWei, linkNative);
  }

  /**
   * @dev calculates LINK paid for gas spent plus a configure premium percentage
   * @param gasLimit the amount of gas used
   * @param fastGasWei the fast gas price
   * @param linkNative the exchange ratio between LINK and Native token
   * @param isExecution if this is triggered by a perform upkeep function
   */
  function _calculatePaymentAmount(
    HotVars memory hotVars,
    uint256 gasLimit,
    uint256 fastGasWei,
    uint256 linkNative,
    bool isExecution
  ) internal view returns (uint96 gasPayment, uint96 premium) {
    uint256 gasWei = fastGasWei * hotVars.gasCeilingMultiplier;
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
      l1CostWei = hotVars.gasCeilingMultiplier * l1CostWei;
    }

    uint256 gasPayment256 = ((weiForGas + l1CostWei) * 1e18) / linkNative;
    uint256 premium256 = (gasPayment * hotVars.paymentPremiumPPB) / 1e9 + uint256(hotVars.flatFeeMicroLink) * 1e12;
    // LINK_TOTAL_SUPPLY < UINT96_MAX
    if (gasPayment + premium > LINK_TOTAL_SUPPLY) revert PaymentGreaterThanAllLINK();
    return (uint96(gasPayment256), uint96(premium256));
  }

  /**
   * @dev ensures the upkeep is not cancelled and the caller is the upkeep admin
   */
  function requireAdminAndNotCancelled(uint256 upkeepId) internal view {
    if (msg.sender != s_upkeepAdmin[upkeepId]) revert OnlyCallableByAdmin();
    if (s_upkeep[upkeepId].maxValidBlocknumber != UINT32_MAX) revert UpkeepCancelled();
  }

  /**
   * @dev generates a PerformPaymentParams struct containing payment information for an upkeep
   */
  function _generatePerformPaymentParams(
    Upkeep memory upkeep,
    HotVars memory hotVars,
    bool isExecution // Whether this is an actual perform execution or just a simulation
  ) internal view returns (PerformPaymentParams memory) {
    (uint256 fastGasWei, uint256 linkNative) = _getFeedData(hotVars);
    (uint96 gasPayment, uint96 premium) = _calculatePaymentAmount(
      hotVars,
      upkeep.executeGas,
      fastGasWei,
      linkNative,
      isExecution
    );

    return PerformPaymentParams({fastGasWei: fastGasWei, linkNative: linkNative, maxLinkPayment: gasPayment + premium});
  }
}
