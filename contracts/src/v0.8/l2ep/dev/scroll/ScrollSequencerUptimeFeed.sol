// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {SequencerUptimeFeed} from "../SequencerUptimeFeed.sol";

import {IL2ScrollMessenger} from "@scroll-tech/contracts/L2/IL2ScrollMessenger.sol";

/// @title ScrollSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates, and records a new answer if the status changed
contract ScrollSequencerUptimeFeed is SequencerUptimeFeed {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ScrollSequencerUptimeFeed 1.0.0";

  /// @notice Address must not be the zero address
  error ZeroAddress();

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  IL2ScrollMessenger private immutable s_l2CrossDomainMessenger;

  /// @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
  /// @param l2CrossDomainMessengerAddr Address of the L2CrossDomainMessenger contract
  /// @param initialStatus The initial status of the feed
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) SequencerUptimeFeed(l1SenderAddress, true) {
    if (l2CrossDomainMessengerAddr == address(0)) {
      revert ZeroAddress();
    }

    s_l2CrossDomainMessenger = IL2ScrollMessenger(l2CrossDomainMessengerAddr);

    _recordRound(1, initialStatus, uint64(block.timestamp), uint64(block.timestamp));
  }

  /// @notice Reverts if the sender is not allowed to call `updateStatus`
  function _requireValidSender() internal view override {
    if (
      msg.sender != address(s_l2CrossDomainMessenger) || s_l2CrossDomainMessenger.xDomainMessageSender() != l1Sender()
    ) {
      revert InvalidSender();
    }
  }
}
