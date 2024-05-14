// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {AggregatorValidatorInterface} from "../../../shared/interfaces/AggregatorValidatorInterface.sol";
import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {ScrollSequencerUptimeFeedInterface} from "../interfaces/ScrollSequencerUptimeFeedInterface.sol";

import {SimpleWriteAccessController} from "../../../shared/access/SimpleWriteAccessController.sol";

import {IL1MessageQueue} from "@scroll-tech/contracts/L1/rollup/IL1MessageQueue.sol";
import {IL1ScrollMessenger} from "@scroll-tech/contracts/L1/IL1ScrollMessenger.sol";

/// @title ScrollValidator - makes cross chain call to update the Sequencer Uptime Feed on L2
contract ScrollValidator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L2_UPTIME_FEED_ADDR;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_MSG_QUEUE_ADDR;

  string public constant override typeAndVersion = "ScrollValidator 1.0.0";
  int256 private constant ANSWER_SEQ_OFFLINE = 1;
  uint32 private s_gasLimit;

  /// @notice emitted when gas cost to spend on L2 is updated
  /// @param gasLimit updated gas cost
  event GasLimitUpdated(uint32 gasLimit);

  /// @param l1CrossDomainMessengerAddress address the L1CrossDomainMessenger contract address
  /// @param l2UptimeFeedAddr the address of the ScrollSequencerUptimeFeed contract address
  /// @param gasLimit the gasLimit to use for sending a message from L1 to L2
  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    address l1MessageQueueAddr,
    uint32 gasLimit
  ) {
    // solhint-disable-next-line gas-custom-errors
    require(l1CrossDomainMessengerAddress != address(0), "Invalid xDomain Messenger address");
    // solhint-disable-next-line gas-custom-errors
    require(l1MessageQueueAddr != address(0), "Invalid L1 message queue address");
    // solhint-disable-next-line gas-custom-errors
    require(l2UptimeFeedAddr != address(0), "Invalid ScrollSequencerUptimeFeed contract address");
    L1_CROSS_DOMAIN_MESSENGER_ADDRESS = l1CrossDomainMessengerAddress;
    L2_UPTIME_FEED_ADDR = l2UptimeFeedAddr;
    L1_MSG_QUEUE_ADDR = l1MessageQueueAddr;
    s_gasLimit = gasLimit;
  }

  /// @notice sets the new gas cost to spend when sending cross chain message
  /// @param gasLimit the updated gas cost
  function setGasLimit(uint32 gasLimit) external onlyOwner {
    s_gasLimit = gasLimit;
    emit GasLimitUpdated(gasLimit);
  }

  /// @notice fetches the gas cost of sending a cross chain message
  function getGasLimit() external view returns (uint32) {
    return s_gasLimit;
  }

  /// @notice makes this contract payable
  /// @dev receives funds:
  ///  - to use them (if configured) to pay for L2 execution on L1
  ///  - when withdrawing funds from L2 xDomain alias address (pay for L2 execution on L2)
  receive() external payable {}

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
        ScrollSequencerUptimeFeedInterface.updateStatus.selector,
        currentAnswer == ANSWER_SEQ_OFFLINE,
        uint64(block.timestamp)
      ),
      s_gasLimit
    );

    // return success
    return true;
  }
}
