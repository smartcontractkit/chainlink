// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {AggregatorInterface} from "../../../interfaces/AggregatorInterface.sol";
import {AggregatorV3Interface} from "../../../interfaces/AggregatorV3Interface.sol";
import {AggregatorV2V3Interface} from "../../../interfaces/AggregatorV2V3Interface.sol";
import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {OptimismSequencerUptimeFeedInterface} from "./../interfaces/OptimismSequencerUptimeFeedInterface.sol";
import {SimpleReadAccessController} from "../../../shared/access/SimpleReadAccessController.sol";
import {ConfirmedOwner} from "../../../shared/access/ConfirmedOwner.sol";
import {IL2CrossDomainMessenger} from "@eth-optimism/contracts/L2/messaging/IL2CrossDomainMessenger.sol";

/**
 * @title OptimismSequencerUptimeFeed - L2 sequencer uptime status aggregator
 * @notice L2 contract that receives status updates from a specific L1 address,
 *  records a new answer if the status changed
 */
contract OptimismSequencerUptimeFeed is
  AggregatorV2V3Interface,
  OptimismSequencerUptimeFeedInterface,
  TypeAndVersionInterface,
  SimpleReadAccessController
{
  /// @dev Round info (for uptime history)
  struct Round {
    bool status;
    uint64 startedAt;
    uint64 updatedAt;
  }

  /// @dev Packed state struct to save sloads
  struct FeedState {
    uint80 latestRoundId;
    bool latestStatus;
    uint64 startedAt;
    uint64 updatedAt;
  }

  /// @notice Sender is not the L2 messenger
  error InvalidSender();
  /// @notice Replacement for AggregatorV3Interface "No data present"
  error NoDataPresent();

  event L1SenderTransferred(address indexed from, address indexed to);
  /// @dev Emitted when an `updateStatus` call is ignored due to unchanged status or stale timestamp
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);
  /// @dev Emitted when a updateStatus is called without the status changing
  event RoundUpdated(int256 status, uint64 updatedAt);

  uint8 public constant override decimals = 0;
  string public constant override description = "L2 Sequencer Uptime Status Feed";
  uint256 public constant override version = 1;

  /// @dev L1 address
  address private s_l1Sender;
  /// @dev s_latestRoundId == 0 means this contract is uninitialized.
  FeedState private s_feedState = FeedState({latestRoundId: 0, latestStatus: false, startedAt: 0, updatedAt: 0});
  mapping(uint80 => Round) private s_rounds;

  IL2CrossDomainMessenger private immutable s_l2CrossDomainMessenger;

  /**
   * @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
   * @param l2CrossDomainMessengerAddr Address of the L2CrossDomainMessenger contract
   * @param initialStatus The initial status of the feed
   */
  constructor(address l1SenderAddress, address l2CrossDomainMessengerAddr, bool initialStatus) {
    setL1Sender(l1SenderAddress);
    s_l2CrossDomainMessenger = IL2CrossDomainMessenger(l2CrossDomainMessengerAddr);
    uint64 timestamp = uint64(block.timestamp);

    // Initialise roundId == 1 as the first round
    recordRound(1, initialStatus, timestamp);
  }

  /**
   * @notice Check if a roundId is valid in this current contract state
   * @dev Mainly used for AggregatorV2V3Interface functions
   * @param roundId Round ID to check
   */
  function isValidRound(uint256 roundId) private view returns (bool) {
    return roundId > 0 && roundId <= type(uint80).max && s_feedState.latestRoundId >= roundId;
  }

  /**
   * @notice versions:
   *
   * - OptimismSequencerUptimeFeed 1.0.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismSequencerUptimeFeed 1.0.0";
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
   * @dev Returns an AggregatorV2V3Interface compatible answer from status flag
   *
   * @param status The status flag to convert to an aggregator-compatible answer
   */
  function getStatusAnswer(bool status) private pure returns (int256) {
    return status ? int256(1) : int256(0);
  }

  /**
   * @notice Helper function to record a round and set the latest feed state.
   *
   * @param roundId The round ID to record
   * @param status Sequencer status
   * @param timestamp The L1 block timestamp of status update
   */
  function recordRound(uint80 roundId, bool status, uint64 timestamp) private {
    uint64 updatedAt = uint64(block.timestamp);
    Round memory nextRound = Round(status, timestamp, updatedAt);
    FeedState memory feedState = FeedState(roundId, status, timestamp, updatedAt);

    s_rounds[roundId] = nextRound;
    s_feedState = feedState;

    emit NewRound(roundId, msg.sender, timestamp);
    emit AnswerUpdated(getStatusAnswer(status), roundId, timestamp);
  }

  /**
   * @notice Helper function to update when a round was last updated
   *
   * @param roundId The round ID to update
   * @param status Sequencer status
   */
  function updateRound(uint80 roundId, bool status) private {
    uint64 updatedAt = uint64(block.timestamp);
    s_rounds[roundId].updatedAt = updatedAt;
    s_feedState.updatedAt = updatedAt;
    emit RoundUpdated(getStatusAnswer(status), updatedAt);
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
    if (
      msg.sender != address(s_l2CrossDomainMessenger) || s_l2CrossDomainMessenger.xDomainMessageSender() != s_l1Sender
    ) {
      revert InvalidSender();
    }

    // Ignore if latest recorded timestamp is newer
    if (feedState.startedAt > timestamp) {
      emit UpdateIgnored(feedState.latestStatus, feedState.startedAt, status, timestamp);
      return;
    }

    if (feedState.latestStatus == status) {
      updateRound(feedState.latestRoundId, status);
    } else {
      feedState.latestRoundId += 1;
      recordRound(feedState.latestRoundId, status, timestamp);
    }
  }

  /// @inheritdoc AggregatorInterface
  function latestAnswer() external view override checkAccess returns (int256) {
    FeedState memory feedState = s_feedState;
    return getStatusAnswer(feedState.latestStatus);
  }

  /// @inheritdoc AggregatorInterface
  function latestTimestamp() external view override checkAccess returns (uint256) {
    FeedState memory feedState = s_feedState;
    return feedState.startedAt;
  }

  /// @inheritdoc AggregatorInterface
  function latestRound() external view override checkAccess returns (uint256) {
    FeedState memory feedState = s_feedState;
    return feedState.latestRoundId;
  }

  /// @inheritdoc AggregatorInterface
  function getAnswer(uint256 roundId) external view override checkAccess returns (int256) {
    if (isValidRound(roundId)) {
      return getStatusAnswer(s_rounds[uint80(roundId)].status);
    }

    revert NoDataPresent();
  }

  /// @inheritdoc AggregatorInterface
  function getTimestamp(uint256 roundId) external view override checkAccess returns (uint256) {
    if (isValidRound(roundId)) {
      return s_rounds[uint80(roundId)].startedAt;
    }

    revert NoDataPresent();
  }

  /// @inheritdoc AggregatorV3Interface
  function getRoundData(
    uint80 _roundId
  )
    public
    view
    override
    checkAccess
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    if (isValidRound(_roundId)) {
      Round memory round = s_rounds[_roundId];
      answer = getStatusAnswer(round.status);
      startedAt = uint256(round.startedAt);
      roundId = _roundId;
      updatedAt = uint256(round.updatedAt);
      answeredInRound = roundId;
    } else {
      revert NoDataPresent();
    }
  }

  /// @inheritdoc AggregatorV3Interface
  function latestRoundData()
    external
    view
    override
    checkAccess
    returns (uint80 roundId, int256 answer, uint256 startedAt, uint256 updatedAt, uint80 answeredInRound)
  {
    FeedState memory feedState = s_feedState;

    roundId = feedState.latestRoundId;
    answer = getStatusAnswer(feedState.latestStatus);
    startedAt = feedState.startedAt;
    updatedAt = feedState.updatedAt;
    answeredInRound = roundId;
  }
}
