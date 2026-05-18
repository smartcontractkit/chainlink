// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../libraries/Client.sol";

/// @notice Interface for plug-in message hook contracts that intercept OffRamp & OnRamp messages
///         and perform validations / state changes on top of the messages. The interceptor functions are expected to
///         revert on validation failures.
interface IMessageInterceptor {
  /// @notice Common error that can be thrown on validation failures and used by consumers
  /// @param errorReason abi encoded revert reason
  error MessageValidationError(bytes errorReason);

  /// @notice Intercepts & validates the given OffRamp message. Reverts on validation failure
  /// @param message to validate
  function onInboundMessage(
    Client.Any2EVMMessage memory message
  ) external;

  /// @notice Intercepts & validates the given OnRamp message. Reverts on validation failure
  /// @param destChainSelector remote destination chain selector where the message is being sent to
  /// @param message to validate
  function onOutboundMessage(uint64 destChainSelector, Client.EVM2AnyMessage memory message) external;
}
