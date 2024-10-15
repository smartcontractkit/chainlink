// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseSequencerUptimeFeed} from "../base/BaseSequencerUptimeFeed.sol";

import {IL2ScrollMessenger} from "@scroll-tech/contracts/L2/IL2ScrollMessenger.sol";

/// @title ScrollSequencerUptimeFeed - L2 sequencer uptime status aggregator
/// @notice L2 contract that receives status updates, and records a new answer if the status changed
contract ScrollSequencerUptimeFeed is BaseSequencerUptimeFeed {
  error ZeroAddress();

  string public constant override typeAndVersion = "ScrollSequencerUptimeFeed 1.1.0-dev";

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

  function _validateSender(address l1Sender) internal view override {
    if (
      msg.sender != address(s_l2CrossDomainMessenger) || s_l2CrossDomainMessenger.xDomainMessageSender() != l1Sender
    ) {
      revert InvalidSender();
    }
  }
}
