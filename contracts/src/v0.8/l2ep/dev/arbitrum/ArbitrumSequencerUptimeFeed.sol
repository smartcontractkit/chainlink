// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AddressAliasHelper} from "../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import {AggregatorInterface} from "../../../shared/interfaces/AggregatorInterface.sol";
import {AggregatorV3Interface} from "../../../shared/interfaces/AggregatorV3Interface.sol";
import {IFlags} from "../interfaces/IFlags.sol";
import {BaseSequencerUptimeFeed} from "../shared/BaseSequencerUptimeFeed.sol";

/**
 * @title ArbitrumSequencerUptimeFeed - L2 sequencer uptime status aggregator
 * @notice L2 contract that receives status updates from a specific L1 address,
 *  records a new answer if the status changed, and raises or lowers the flag on the
 *   stored Flags contract.
 */
contract ArbitrumSequencerUptimeFeed is BaseSequencerUptimeFeed {
  /**
   * @notice versions:
   *
   * - ArbitrumSequencerUptimeFeed 1.0.0: initial release
   *
   */
  string public constant override typeAndVersion = "ArbitrumSequencerUptimeFeed 1.0.0";

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

  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address public constant FLAG_L2_SEQ_OFFLINE =
    address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-seq-offline")) - 1)));

  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  uint8 public constant override decimals = 0;
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override description = "L2 Sequencer Uptime Status Feed";
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  uint256 public constant override version = 1;

  /// @dev Flags contract to raise/lower flags on, during status transitions
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  IFlags public immutable FLAGS;

  /// @dev s_latestRoundId == 0 means this contract is uninitialized.
  FeedState private s_feedState = FeedState({latestRoundId: 0, latestStatus: false, latestTimestamp: 0});
  mapping(uint80 => Round) private s_rounds;

  /**
   * @param flagsAddress Address of the Flags contract on L2
   * @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
   */
  constructor(
    address flagsAddress,
    address l1SenderAddress,
    bool initialStatus
  ) BaseSequencerUptimeFeed(l1SenderAddress, initialStatus) {
    _setL1Sender(l1SenderAddress);

    FLAGS = IFlags(flagsAddress);
  }

  /// @notice Check that this contract is initialised, otherwise throw
  function _requireInitialized(uint80 latestRoundId) private pure {
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
    _recordRound(1, currentStatus, timestamp);

    emit Initialized();
  }

  /**
   * @dev Returns an AggregatorV2V3Interface compatible answer from status flag
   *
   * @param status The status flag to convert to an aggregator-compatible answer
   */
  function _getStatusAnswer(bool status) private pure returns (int256) {
    return status ? int256(1) : int256(0);
  }

  /**
   * @notice Raise or lower the flag on the stored Flags contract.
   */
  function _forwardStatusToFlags(bool status) private {
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
  function _recordRound(uint80 roundId, bool status, uint64 timestamp) private {
    Round memory nextRound = Round(status, timestamp);
    FeedState memory feedState = FeedState(roundId, status, timestamp);

    s_rounds[roundId] = nextRound;
    s_feedState = feedState;

    emit NewRound(roundId, msg.sender, timestamp);
    emit AnswerUpdated(_getStatusAnswer(status), roundId, timestamp);
  }

  function _validateSender(address l1Sender) internal view override {
    _requireInitialized(s_feedState.latestRoundId);

    if (msg.sender != AddressAliasHelper.applyL1ToL2Alias(s_l1Sender)) {
      revert InvalidSender();
    }
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
    // _requireInitialized(feedState.latestRoundId);
    // if (msg.sender != AddressAliasHelper.applyL1ToL2Alias(s_l1Sender)) {
    //   revert InvalidSender();
    // }

    // Ignore if status did not change or latest recorded timestamp is newer
    if (feedState.latestStatus == status || feedState.latestTimestamp > timestamp) {
      emit UpdateIgnored(feedState.latestStatus, feedState.latestTimestamp, status, timestamp);
      return;
    }

    // Prepare a new round with updated status
    feedState.latestRoundId += 1;
    _recordRound(feedState.latestRoundId, status, timestamp);

    _forwardStatusToFlags(status);
  }

  /// @inheritdoc AggregatorInterface
  function latestAnswer() external view override checkAccess returns (int256) {
    FeedState memory feedState = s_feedState;
    _requireInitialized(feedState.latestRoundId);
    return _getStatusAnswer(feedState.latestStatus);
  }

  /// @inheritdoc AggregatorInterface
  function latestTimestamp() external view override checkAccess returns (uint256) {
    FeedState memory feedState = s_feedState;
    _requireInitialized(feedState.latestRoundId);
    return feedState.latestTimestamp;
  }

  /// @inheritdoc AggregatorInterface
  function latestRound() external view override checkAccess returns (uint256) {
    FeedState memory feedState = s_feedState;
    _requireInitialized(feedState.latestRoundId);
    return feedState.latestRoundId;
  }

  /// @inheritdoc AggregatorInterface
  function getAnswer(uint256 roundId) external view override checkAccess returns (int256) {
    _requireInitialized(s_feedState.latestRoundId);
    if (_isValidRound(roundId)) {
      return _getStatusAnswer(s_rounds[uint80(roundId)].status);
    }

    return 0;
  }

  /// @inheritdoc AggregatorInterface
  function getTimestamp(uint256 roundId) external view override checkAccess returns (uint256) {
    _requireInitialized(s_feedState.latestRoundId);
    if (_isValidRound(roundId)) {
      return s_rounds[uint80(roundId)].timestamp;
    }

    return 0;
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
    _requireInitialized(s_feedState.latestRoundId);

    if (_isValidRound(_roundId)) {
      Round memory round = s_rounds[_roundId];
      answer = _getStatusAnswer(round.status);
      startedAt = uint256(round.timestamp);
    } else {
      answer = 0;
      startedAt = 0;
    }
    roundId = _roundId;
    updatedAt = startedAt;
    answeredInRound = roundId;

    return (roundId, answer, startedAt, updatedAt, answeredInRound);
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
    _requireInitialized(feedState.latestRoundId);

    roundId = feedState.latestRoundId;
    answer = _getStatusAnswer(feedState.latestStatus);
    startedAt = feedState.latestTimestamp;
    updatedAt = startedAt;
    answeredInRound = roundId;

    return (roundId, answer, startedAt, updatedAt, answeredInRound);
  }
}
