// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {AddressAliasHelper} from "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import {ForwarderInterface} from "./interfaces/ForwarderInterface.sol";
import {AggregatorInterface} from "../interfaces/AggregatorInterface.sol";
import {AggregatorV3Interface} from "../interfaces/AggregatorV3Interface.sol";
import {AggregatorV2V3Interface} from "../interfaces/AggregatorV2V3Interface.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {FlagsInterface} from "./interfaces/FlagsInterface.sol";
import {ArbitrumSequencerUptimeFeedInterface} from "./interfaces/ArbitrumSequencerUptimeFeedInterface.sol";
import {SimpleReadAccessController} from "../SimpleReadAccessController.sol";
import {ConfirmedOwner} from "../ConfirmedOwner.sol";

/**
 * @title ArbitrumSequencerUptimeFeed - L2 sequencer uptime status aggregator
 * @notice L2 contract that receives status updates from a specific L1 address,
 *  records a new answer if the status changed, and raises or lowers the flag on the
 *   stored Flags contract.
 */
contract ArbitrumSequencerUptimeFeed is
  AggregatorV2V3Interface,
  ArbitrumSequencerUptimeFeedInterface,
  TypeAndVersionInterface,
  SimpleReadAccessController
{
  /// @dev Round info (for uptime history)
  struct Round {
    bool status;
    uint64 timestamp;
  }

  /// @dev Packed state struct to save sloads
  struct FeedState {
    uint80 latestRoundId;
    bool latestStatus;
    uint64 latestTimestamp;
  }

  /// @notice Contract is not yet initialized
  error Uninitialized();
  /// @notice Contract is already initialized
  error AlreadyInitialized();
  /// @notice Sender is not the L2 messenger
  error InvalidSender();
  /// @notice Replacement for AggregatorV3Interface "No data present"
  error NoDataPresent();

  event Initialized();
  event L1SenderTransferred(address indexed from, address indexed to);
  /// @dev Emitted when an `updateStatus` call is ignored due to unchanged status or stale timestamp
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);

  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address public constant FLAG_L2_SEQ_OFFLINE =
    address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-seq-offline")) - 1)));

  uint8 public constant override decimals = 0;
  string public constant override description = "L2 Sequencer Uptime Status Feed";
  uint256 public constant override version = 1;

  /// @dev Flags contract to raise/lower flags on, during status transitions
  FlagsInterface public immutable FLAGS;
  /// @dev L1 address
  address private s_l1Sender;
  /// @dev s_latestRoundId == 0 means this contract is uninitialized.
  FeedState private s_feedState = FeedState({latestRoundId: 0, latestStatus: false, latestTimestamp: 0});
  mapping(uint80 => Round) private s_rounds;

  /**
   * @param flagsAddress Address of the Flags contract on L2
   * @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
   */
  constructor(address flagsAddress, address l1SenderAddress) {
    setL1Sender(l1SenderAddress);

    FLAGS = FlagsInterface(flagsAddress);
  }

  /**
   * @notice Check if a roundId is valid in this current contract state
   * @dev Mainly used for AggregatorV2V3Interface functions
   * @param roundId Round ID to check
   */
  function isValidRound(uint256 roundId) private view returns (bool) {
    return roundId > 0 && roundId <= type(uint80).max && s_feedState.latestRoundId >= roundId;
  }

  /// @notice Check that this contract is initialised, otherwise throw
  function requireInitialized(uint80 latestRoundId) private pure {
    if (latestRoundId == 0) {
      revert Uninitialized();
    }
  }

  /**
   * @notice Initialise the first round. Can't be done in the constructor,
   *    because this contract's address must be permissioned by the the Flags contract
   *    (The Flags contract itself is a SimpleReadAccessController).
   */
  function initialize() external onlyOwner {
    FeedState memory feedState = s_feedState;
    if (feedState.latestRoundId != 0) {
      revert AlreadyInitialized();
    }

    uint64 timestamp = uint64(block.timestamp);
    bool currentStatus = FLAGS.getFlag(FLAG_L2_SEQ_OFFLINE);

    // Initialise roundId == 1 as the first round
    recordRound(1, currentStatus, timestamp);

    emit Initialized();
  }

  /**
   * @notice versions:
   *
   * - ArbitrumSequencerUptimeFeed 1.0.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ArbitrumSequencerUptimeFeed 1.0.0";
  }

  /// @return L1 sender address
  function l1Sender() public view virtual returns (address) {
    return s_l1Sender;
  }

  /**
   * @notice Set the allowed L1 sender for this contract to a new L1 sender
   * @dev Can be disabled by setting the L1 sender as `address(0)`. Accessible only by owner.
   * @param to new L1 sender that will be allowed to call `updateStatus` on this contract
   */
  function transferL1Sender(address to) external virtual onlyOwner {
    setL1Sender(to);
  }

  /// @notice internal method that stores the L1 sender
  function setL1Sender(address to) private {
    address from = s_l1Sender;
    if (from != to) {
      s_l1Sender = to;
      emit L1SenderTransferred(from, to);
    }
  }

  /**
   * @notice Messages sent by the stored L1 sender will arrive on L2 with this
   *  address as the `msg.sender`
   * @return L2-aliased form of the L1 sender address
   */
  function aliasedL1MessageSender() public view returns (address) {
    return AddressAliasHelper.applyL1ToL2Alias(l1Sender());
  }

  /**
   * @dev Returns an AggregatorV2V3Interface compatible answer from status flag
   *
   * @param status The status flag to convert to an aggregator-compatible answer
   */
  function getStatusAnswer(bool status) private pure returns (int256) {
    return status ? int256(1) : int256(0);
  }

  /**
   * @notice Raise or lower the flag on the stored Flags contract.
   */
  function forwardStatusToFlags(bool status) private {
    if (status) {
      FLAGS.raiseFlag(FLAG_L2_SEQ_OFFLINE);
    } else {
      FLAGS.lowerFlag(FLAG_L2_SEQ_OFFLINE);
    }
  }

  /**
   * @notice Helper function to record a round and set the latest feed state.
   *
   * @param roundId The round ID to record
   * @param status Sequencer status
   * @param timestamp Block timestamp of status update
   */
  function recordRound(
    uint80 roundId,
    bool status,
    uint64 timestamp
  ) private {
    Round memory nextRound = Round(status, timestamp);
    FeedState memory feedState = FeedState(roundId, status, timestamp);

    s_rounds[roundId] = nextRound;
    s_feedState = feedState;

    emit NewRound(roundId, msg.sender, timestamp);
    emit AnswerUpdated(getStatusAnswer(status), roundId, timestamp);
  }

  /**
   * @notice Record a new status and timestamp if it has changed since the last round.
   * @dev This function will revert if not called from `l1Sender` via the L1->L2 messenger.
   *
   * @param status Sequencer status
   * @param timestamp Block timestamp of status update
   */
  function updateStatus(bool status, uint64 timestamp) external override {
    FeedState memory feedState = s_feedState;
    requireInitialized(feedState.latestRoundId);
    if (msg.sender != aliasedL1MessageSender()) {
      revert InvalidSender();
    }

    // Ignore if status did not change or latest recorded timestamp is newer
    if (feedState.latestStatus == status || feedState.latestTimestamp > timestamp) {
      emit UpdateIgnored(feedState.latestStatus, feedState.latestTimestamp, status, timestamp);
      return;
    }

    // Prepare a new round with updated status
    feedState.latestRoundId += 1;
    recordRound(feedState.latestRoundId, status, timestamp);

    forwardStatusToFlags(status);
  }

  /// @inheritdoc AggregatorInterface
  function latestAnswer() external view override checkAccess returns (int256) {
    FeedState memory feedState = s_feedState;
    requireInitialized(feedState.latestRoundId);
    return getStatusAnswer(feedState.latestStatus);
  }

  /// @inheritdoc AggregatorInterface
  function latestTimestamp() external view override checkAccess returns (uint256) {
    FeedState memory feedState = s_feedState;
    requireInitialized(feedState.latestRoundId);
    return feedState.latestTimestamp;
  }

  /// @inheritdoc AggregatorInterface
  function latestRound() external view override checkAccess returns (uint256) {
    FeedState memory feedState = s_feedState;
    requireInitialized(feedState.latestRoundId);
    return feedState.latestRoundId;
  }

  /// @inheritdoc AggregatorInterface
  function getAnswer(uint256 roundId) external view override checkAccess returns (int256) {
    requireInitialized(s_feedState.latestRoundId);
    if (isValidRound(roundId)) {
      return getStatusAnswer(s_rounds[uint80(roundId)].status);
    }

    return 0;
  }

  /// @inheritdoc AggregatorInterface
  function getTimestamp(uint256 roundId) external view override checkAccess returns (uint256) {
    requireInitialized(s_feedState.latestRoundId);
    if (isValidRound(roundId)) {
      return s_rounds[uint80(roundId)].timestamp;
    }

    return 0;
  }

  /// @inheritdoc AggregatorV3Interface
  function getRoundData(uint80 _roundId)
    public
    view
    override
    checkAccess
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    requireInitialized(s_feedState.latestRoundId);

    if (isValidRound(_roundId)) {
      Round memory round = s_rounds[_roundId];
      answer = getStatusAnswer(round.status);
      startedAt = uint256(round.timestamp);
    } else {
      answer = 0;
      startedAt = 0;
    }
    roundId = _roundId;
    updatedAt = startedAt;
    answeredInRound = roundId;
  }

  /// @inheritdoc AggregatorV3Interface
  function latestRoundData()
    external
    view
    override
    checkAccess
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    FeedState memory feedState = s_feedState;
    requireInitialized(feedState.latestRoundId);

    roundId = feedState.latestRoundId;
    answer = getStatusAnswer(feedState.latestStatus);
    startedAt = feedState.latestTimestamp;
    updatedAt = startedAt;
    answeredInRound = roundId;
  }
}
