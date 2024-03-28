// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Validator} from "./Validator.sol";

abstract contract GasLimitValidator is Validator {
  uint32 internal s_gasLimit;

  /// @notice emitted when gas limit on L2 is updated
  /// @param gasLimit updated gas cost
  event GasLimitUpdated(uint32 gasLimit);

  /// @param l1CrossDomainMessengerAddress address the L1CrossDomainMessenger contract address
  /// @param l2UptimeFeedAddr the address of the UptimeFeed contract address
  /// @param gasLimit the gasLimit to use for sending a message from L1 to L2
  constructor(
    address l1CrossDomainMessengerAddress,
    address l2UptimeFeedAddr,
    uint32 gasLimit
  ) Validator(l1CrossDomainMessengerAddress, l2UptimeFeedAddr) {
    s_gasLimit = gasLimit;
  }

  /// @notice sets the new gas limit when sending cross chain message
  /// @param gasLimit the updated gas cost
  function setGasLimit(uint32 gasLimit) external onlyOwner {
    s_gasLimit = gasLimit;
    emit GasLimitUpdated(gasLimit);
  }

  /// @notice fetches the gas cost of sending a cross chain message
  function getGasLimit() external view returns (uint32) {
    return s_gasLimit;
  }
}
