// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {FlagsInterface} from "../interfaces/FlagsInterface.sol";

import {SequencerUptimeFeed} from "../SequencerUptimeFeed.sol";

import {AddressAliasHelper} from "../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";

/// @title ArbitrumSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates from a specific L1 address,
///  records a new answer if the status changed, and raises or lowers the flag on the
///   stored Flags contract.
contract ArbitrumSequencerUptimeFeed is SequencerUptimeFeed {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ArbitrumSequencerUptimeFeed 1.0.0";
  /// @dev Follows: https://eips.ethereum.org/EIPS/eip-1967
  address public constant FLAG_L2_SEQ_OFFLINE =
    address(bytes20(bytes32(uint256(keccak256("chainlink.flags.arbitrum-seq-offline")) - 1)));

  /// @notice Contract is already initialized
  error AlreadyInitialized();

  /// @dev Emitted when the first round is Initialized
  event Initialized();

  /// @dev Flags contract to raise/lower flags on, during status transitions
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  FlagsInterface public immutable FLAGS;

  /// @param flagsAddress Address of the Flags contract on L2
  /// @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
  constructor(address flagsAddress, address l1SenderAddress) SequencerUptimeFeed(l1SenderAddress, false) {
    FLAGS = FlagsInterface(flagsAddress);
  }

  /// @notice Reverts if the sender is not allowed to call `updateStatus`
  modifier requireValidSender() override {
    if (msg.sender != AddressAliasHelper.applyL1ToL2Alias(l1Sender())) {
      revert InvalidSender();
    }
    _;
  }

  /// @notice Initialise the first round. Can't be done in the constructor,
  ///    because this contract's address must be permissioned by the the Flags contract
  ///    (The Flags contract itself is a SimpleReadAccessController).
  function initialize() external onlyOwner {
    FeedState memory feedState = s_feedState;
    if (feedState.latestRoundId != 0) {
      revert AlreadyInitialized();
    }

    // Initialise roundId == 1 as the first round
    _recordRound(1, FLAGS.getFlag(FLAG_L2_SEQ_OFFLINE), uint64(block.timestamp));

    emit Initialized();
  }

  /// @notice Record a new status and timestamp if it has changed since the last round.
  /// @dev This function will revert if not called from `l1Sender` via the L1->L2 messenger.
  ///
  /// @param status Sequencer status
  /// @param timestamp Block timestamp of status update
  function updateStatus(bool status, uint64 timestamp) external override requireInitialized requireValidSender {
    FeedState memory feedState = s_feedState;

    // Ignore if status did not change or latest recorded timestamp is newer
    if (feedState.latestStatus == status || feedState.startedAt > timestamp) {
      emit UpdateIgnored(feedState.latestStatus, feedState.startedAt, status, timestamp);
      return;
    }

    // Prepare a new round with updated status
    feedState.latestRoundId += 1;
    _recordRound(feedState.latestRoundId, status, timestamp);

    // Raise or lower the flag on the stored Flags contract.
    if (status) {
      FLAGS.raiseFlag(FLAG_L2_SEQ_OFFLINE);
    } else {
      FLAGS.lowerFlag(FLAG_L2_SEQ_OFFLINE);
    }
  }
}
