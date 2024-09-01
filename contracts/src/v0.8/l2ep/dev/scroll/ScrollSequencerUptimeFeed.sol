// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";

import {BaseSequencerUptimeFeed} from "../shared/BaseSequencerUptimeFeed.sol";

import {IL2ScrollMessenger} from "@scroll-tech/contracts/L2/IL2ScrollMessenger.sol";

/// @title ScrollSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates, and records a new answer if the status changed
contract ScrollSequencerUptimeFeed is ITypeAndVersion, BaseSequencerUptimeFeed {
  error ZeroAddress();

  string public constant override typeAndVersion = "ScrollSequencerUptimeFeed 1.0.0";

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  IL2ScrollMessenger private immutable s_l2CrossDomainMessenger;

  /// @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
  /// @param l2CrossDomainMessengerAddr Address of the L2CrossDomainMessenger contract
  /// @param initialStatus The initial status of the feed
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) BaseSequencerUptimeFeed(l1SenderAddress, initialStatus) {
    if (l2CrossDomainMessengerAddr == address(0)) {
      revert ZeroAddress();
    }

    s_l2CrossDomainMessenger = IL2ScrollMessenger(l2CrossDomainMessengerAddr);
  }

  /// @notice Record a new status and timestamp if it has changed since the last round.
  /// @dev This function will revert if not called from `l1Sender` via the L1->L2 messenger.
  ///
  /// @param status Sequencer status
  /// @param timestamp Block timestamp of status update
  function updateStatus(bool status, uint64 timestamp) external override {
    FeedState memory feedState = _getFeedState();

    if (
      msg.sender != address(s_l2CrossDomainMessenger) ||
      s_l2CrossDomainMessenger.xDomainMessageSender() != getL1Sender()
    ) {
      revert InvalidSender();
    }

    // Ignore if latest recorded timestamp is newer
    if (feedState.startedAt > timestamp) {
      emit UpdateIgnored(feedState.latestStatus, feedState.startedAt, status, timestamp);
      return;
    }

    if (feedState.latestStatus == status) {
      _updateRound(feedState.latestRoundId, status);
    } else {
      feedState.latestRoundId += 1;
      _recordRound(feedState.latestRoundId, status, timestamp);
    }
  }
}
