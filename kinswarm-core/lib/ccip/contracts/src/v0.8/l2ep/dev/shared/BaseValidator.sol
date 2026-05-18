// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {AggregatorValidatorInterface} from "../../../shared/interfaces/AggregatorValidatorInterface.sol";
import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";

import {SimpleWriteAccessController} from "../../../shared/access/SimpleWriteAccessController.sol";

abstract contract BaseValidator is SimpleWriteAccessController, AggregatorValidatorInterface, ITypeAndVersion {
  /// @notice emitted when gas cost to spend on L2 is updated
  /// @param gasLimit updated gas cost
  event GasLimitUpdated(uint32 gasLimit);

  error L1CrossDomainMessengerAddressZero();
  error L2UptimeFeedAddrZero();

  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L1_CROSS_DOMAIN_MESSENGER_ADDRESS;
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  address public immutable L2_UPTIME_FEED_ADDR;

  int256 internal constant ANSWER_SEQ_OFFLINE = 1;

  uint32 internal s_gasLimit;

  /// @param l1CrossDomainMessengerAddress address the L1CrossDomainMessenger contract address
  /// @param l2UptimeFeedAddr the address of the SequencerUptimeFeed contract address
  /// @param gasLimit the gasLimit to use for sending a message from L1 to L2
  constructor(address l1CrossDomainMessengerAddress, address l2UptimeFeedAddr, uint32 gasLimit) {
    if (l1CrossDomainMessengerAddress == address(0)) {
      revert L1CrossDomainMessengerAddressZero();
    }
    if (l2UptimeFeedAddr == address(0)) {
      revert L2UptimeFeedAddrZero();
    }

    L1_CROSS_DOMAIN_MESSENGER_ADDRESS = l1CrossDomainMessengerAddress;
    L2_UPTIME_FEED_ADDR = l2UptimeFeedAddr;
    s_gasLimit = gasLimit;
  }

  /// @notice fetches the gas cost of sending a cross chain message
  function getGasLimit() external view returns (uint32) {
    return s_gasLimit;
  }

  /// @notice sets the new gas cost to spend when sending cross chain message
  /// @param gasLimit the updated gas cost
  function setGasLimit(uint32 gasLimit) external onlyOwner {
    s_gasLimit = gasLimit;
    emit GasLimitUpdated(gasLimit);
  }

  /// @notice makes this contract payable
  /// @dev receives funds:
  ///  - to use them (if configured) to pay for L2 execution on L1
  ///  - when withdrawing funds from L2 xDomain alias address (pay for L2 execution on L2)
  receive() external payable {}
}
