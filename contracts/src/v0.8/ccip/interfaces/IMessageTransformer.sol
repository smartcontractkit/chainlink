// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Internal} from "../libraries/Internal.sol";

/// @notice Interface for message hook contracts that transform OffRamp & OnRamp messages.
/// The transformer functions are expected to revert on transform failures.
interface IMessageTransformer {
  /// @notice Common error that can be thrown on transform failures and used by consumers
  /// @param errorReason abi encoded revert reason
  error MessageTransformError(bytes errorReason);

  /// @notice Transforms the given OnRamp message. Reverts on transform failure
  /// Can be used for modifying message prior to execution
  /// e.g. adding extra metadata to the data field in CCIP message
  /// by receiver contract
  /// @param message to transform
  /// @return transformed message
  function transformOutboundMessage(
    Internal.EVM2AnyRampMessage memory message
  ) external returns (Internal.EVM2AnyRampMessage memory);

  /// @notice Transforms the given OffRamp message. Reverts on transform failure
  /// One possible use case is to add extra metadata to the data field in CCIP message
  /// before consumption by the receiver contract
  /// @param message to transform
  /// @return transformed message
  function transformInboundMessage(
    Internal.Any2EVMRampMessage memory message
  ) external returns (Internal.Any2EVMRampMessage memory);
}
