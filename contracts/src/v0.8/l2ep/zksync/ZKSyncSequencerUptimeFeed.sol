// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseSequencerUptimeFeed} from "../base/BaseSequencerUptimeFeed.sol";

import {AddressAliasHelper} from "../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";

/// @title ZKSyncSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates from a specific L1 address,
///  records a new answer if the status changed
contract ZKSyncSequencerUptimeFeed is BaseSequencerUptimeFeed {
  string public constant override typeAndVersion = "ZKSyncSequencerUptimeFeed 1.1.0-dev";

  /// @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
  /// @param initialStatus The initial status of the feed
  constructor(address l1SenderAddress, bool initialStatus) BaseSequencerUptimeFeed(l1SenderAddress, initialStatus) {}

  function _validateSender(address l1Sender) internal view override {
    address aliasedL1Sender = AddressAliasHelper.applyL1ToL2Alias(l1Sender);

    if (msg.sender != aliasedL1Sender) {
      revert InvalidSender();
    }
  }
}
