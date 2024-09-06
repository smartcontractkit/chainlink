// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseSequencerUptimeFeed} from "../shared/BaseSequencerUptimeFeed.sol";

import {IL2CrossDomainMessenger} from "@eth-optimism/contracts/L2/messaging/IL2CrossDomainMessenger.sol";

/**
 * @title OptimismSequencerUptimeFeed - L2 sequencer uptime status aggregator
 * @notice L2 contract that receives status updates from a specific L1 address,
 *  records a new answer if the status changed
 */
contract OptimismSequencerUptimeFeed is BaseSequencerUptimeFeed {
  string public constant override typeAndVersion = "OptimismSequencerUptimeFeed 1.1.0-dev";

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  IL2CrossDomainMessenger private immutable s_l2CrossDomainMessenger;

  /**
   * @param l1SenderAddress Address of the L1 contract that is permissioned to call this contract
   * @param l2CrossDomainMessengerAddr Address of the L2CrossDomainMessenger contract
   * @param initialStatus The initial status of the feed
   */
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) BaseSequencerUptimeFeed(l1SenderAddress, initialStatus) {
    s_l2CrossDomainMessenger = IL2CrossDomainMessenger(l2CrossDomainMessengerAddr);
  }

  function _validateSender(address l1Sender) internal view override {
    if (
      msg.sender != address(s_l2CrossDomainMessenger) || s_l2CrossDomainMessenger.xDomainMessageSender() != l1Sender
    ) {
      revert InvalidSender();
    }
  }
}
