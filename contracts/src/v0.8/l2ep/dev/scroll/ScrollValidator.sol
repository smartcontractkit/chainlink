// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ISequencerUptimeFeed} from "../interfaces/ISequencerUptimeFeed.sol";

import {BaseValidator} from "../shared/BaseValidator.sol";

import {IL1MessageQueue} from "@scroll-tech/contracts/L1/rollup/IL1MessageQueue.sol";
import {IL1ScrollMessenger} from "@scroll-tech/contracts/L1/IL1ScrollMessenger.sol";

/// @title ScrollValidator - makes cross chain call to update the Sequencer Uptime Feed on L2
contract ScrollValidator is BaseValidator {
  string public constant override typeAndVersion = "ScrollValidator 1.1.0-dev";

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_MSG_QUEUE_ADDR;

  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    address l1MessageQueueAddr,
    uint32 gasLimit
  ) BaseValidator(l1CrossDomainMessengerAddress, l2UptimeFeedAddr, gasLimit) {
    // solhint-disable-next-line gas-custom-errors
    require(l1MessageQueueAddr != address(0), "Invalid L1 message queue address");
    L1_MSG_QUEUE_ADDR = l1MessageQueueAddr;
  }

  /// @notice validate method sends an xDomain L2 tx to update Uptime Feed contract on L2.
  /// @dev A message is sent using the L1CrossDomainMessenger. This method is accessed controlled.
  /// @param currentAnswer new aggregator answer - value of 1 considers the sequencer offline.
  function validate(
    uint256 /* previousRoundId */,
    int256 /* previousAnswer */,
    uint256 /* currentRoundId */,
    int256 currentAnswer
  ) external override checkAccess returns (bool) {
    // Make the xDomain call
    IL1ScrollMessenger(L1_CROSS_DOMAIN_MESSENGER_ADDRESS).sendMessage{
      value: IL1MessageQueue(L1_MSG_QUEUE_ADDR).estimateCrossDomainMessageFee(s_gasLimit)
    }(
      L2_UPTIME_FEED_ADDR,
      0,
      abi.encodeWithSelector(
        ISequencerUptimeFeed.updateStatus.selector,
        currentAnswer == ANSWER_SEQ_OFFLINE,
        uint64(block.timestamp)
      ),
      s_gasLimit
    );

    // return success
    return true;
  }
}
