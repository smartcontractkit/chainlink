// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {AggregatorValidatorInterface} from "../../shared/interfaces/AggregatorValidatorInterface.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";

import {SimpleWriteAccessController} from "../../shared/access/SimpleWriteAccessController.sol";

abstract contract Validator is TypeAndVersionInterface, AggregatorValidatorInterface, SimpleWriteAccessController {
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L2_UPTIME_FEED_ADDR;

  int256 internal constant ANSWER_SEQ_OFFLINE = 1;

  /// @param l1CrossDomainMessengerAddress address the L1CrossDomainMessenger contract address
  /// @param l2UptimeFeedAddr the address of the ScrollSequencerUptimeFeed contract address
  constructor(address l1CrossDomainMessengerAddress, address l2UptimeFeedAddr) {
    // solhint-disable-next-line custom-errors
    require(l1CrossDomainMessengerAddress != address(0), "Invalid xDomain Messenger address");
    // solhint-disable-next-line custom-errors
    require(l2UptimeFeedAddr != address(0), "Invalid ScrollSequencerUptimeFeed contract address");
    L1_CROSS_DOMAIN_MESSENGER_ADDRESS = l1CrossDomainMessengerAddress;
    L2_UPTIME_FEED_ADDR = l2UptimeFeedAddr;
  }
}
