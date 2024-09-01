// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {AddressAliasHelper} from "@zksync/contracts/l2-contracts/contracts/vendor/AddressAliasHelper.sol";
import {BaseSequencerUptimeFeed} from "../shared/BaseSequencerUptimeFeed.sol";

/// @title ZKSyncSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates from a specific L1 address,
///  records a new answer if the status changed
contract ZKSyncSequencerUptimeFeed is ITypeAndVersion, BaseSequencerUptimeFeed {
  string public constant override typeAndVersion = "ZKSyncSequencerUptimeFeed 1.0.0";

  /// @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
  /// @param initialStatus The initial status of the feed
  constructor(address l1SenderAddress, bool initialStatus) BaseSequencerUptimeFeed(l1SenderAddress, initialStatus) {}

  /// @notice Record a new status and timestamp if it has changed since the last round.
  /// @dev This function will revert if not called from `l1Sender` via the L1->L2 messenger.
  /// @param status Sequencer status
  /// @param timestamp Block timestamp of status update
  function updateStatus(bool status, uint64 timestamp) external override {
    address aliasedL1Sender = AddressAliasHelper.applyL1ToL2Alias(l1Sender());

    if (msg.sender != aliasedL1Sender) {
      revert InvalidSender();
    }

    FeedState memory feedState = _getFeedState();
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
