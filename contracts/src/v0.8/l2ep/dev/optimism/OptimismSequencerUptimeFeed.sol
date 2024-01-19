// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {SequencerUptimeFeed} from "../SequencerUptimeFeed.sol";

import {IL2CrossDomainMessenger} from "@eth-optimism/contracts/L2/messaging/IL2CrossDomainMessenger.sol";

/// @title OptimismSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates from a specific L1 address, and records a new answer if the status changed
contract OptimismSequencerUptimeFeed is SequencerUptimeFeed {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "OptimismSequencerUptimeFeed 1.0.0";

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  IL2CrossDomainMessenger private immutable s_l2CrossDomainMessenger;

  /// @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
  /// @param l2CrossDomainMessengerAddr Address of the L2CrossDomainMessenger contract
  /// @param initialStatus The initial status of the feed
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) SequencerUptimeFeed(l1SenderAddress, true) {
    s_l2CrossDomainMessenger = IL2CrossDomainMessenger(l2CrossDomainMessengerAddr);

    _recordRound(1, initialStatus, uint64(block.timestamp));
  }

  /// @notice Reverts if the sender is not allowed to call `updateStatus`
  modifier requireValidSender() override {
    if (
      msg.sender != address(s_l2CrossDomainMessenger) || s_l2CrossDomainMessenger.xDomainMessageSender() != l1Sender()
    ) {
      revert InvalidSender();
    }
    _;
  }
}
