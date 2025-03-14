// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransmitter} from "./IMessageTransmitter.sol";
import {ITokenMessenger} from "./ITokenMessenger.sol";

import {Ownable2StepMsgSender} from "../../../shared/access/Ownable2StepMsgSender.sol";

import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @title CCTP Message Transmitter Proxy
/// @notice A proxy contract for handling messages transmitted via the Cross Chain Transfer Protocol (CCTP).
/// @dev This contract is responsible for sending messages to the `IMessageTransmitter` and ensuring only allowed callers can invoke it.
contract CCTPMessageTransmitterProxy is Ownable2StepMsgSender {
  using EnumerableSet for EnumerableSet.AddressSet;

  /// @notice Error thrown when a function is called by an unauthorized address.
  error Unauthorized(address caller);

  /// @notice Emitted when an allowed caller is added.
  event AllowedCallerAdded(address indexed caller);
  /// @notice Emitted when an allowed caller is removed.
  event AllowedCallerRemoved(address indexed caller);

  struct AllowedCallerConfigArgs {
    address caller;
    bool allowed;
  }

  /// @notice Immutable reference to the `IMessageTransmitter` contract.
  IMessageTransmitter public immutable i_cctpTransmitter;

  /// @notice Enumerable set of addresses allowed to call `receiveMessage`.
  EnumerableSet.AddressSet private s_allowedCallers;

  /// @notice One-time cyclic dependency between TokenPool and MessageTransmitter.
  constructor(
    ITokenMessenger tokenMessenger
  ) {
    i_cctpTransmitter = IMessageTransmitter(tokenMessenger.localMessageTransmitter());
  }

  /// @notice Receives a message from the `IMessageTransmitter` contract and validates it.
  /// @dev Can only be called by an allowed caller to process incoming messages.
  /// @param message The payload of the message being received.
  /// @param attestation The cryptographic proof validating the message.
  /// @return success A boolean indicating if the message was successfully processed.
  function receiveMessage(bytes calldata message, bytes calldata attestation) external returns (bool success) {
    if (!s_allowedCallers.contains(msg.sender)) {
      revert Unauthorized(msg.sender);
    }
    return i_cctpTransmitter.receiveMessage(message, attestation);
  }

  /// @notice Configures the allowed callers for the `receiveMessage` function.
  /// @param configArgs An array of `AllowedCallerConfigArgs` structs.
  function configureAllowedCallers(
    AllowedCallerConfigArgs[] calldata configArgs
  ) external onlyOwner {
    for (uint256 i = 0; i < configArgs.length; ++i) {
      if (configArgs[i].allowed) {
        if (s_allowedCallers.add(configArgs[i].caller)) {
          emit AllowedCallerAdded(configArgs[i].caller);
        }
      } else {
        if (s_allowedCallers.remove(configArgs[i].caller)) {
          emit AllowedCallerRemoved(configArgs[i].caller);
        }
      }
    }
  }

  /// @notice Checks if the caller is allowed to call the `receiveMessage` function.
  /// @param caller The address to check.
  /// @return allowed A boolean indicating if the caller is allowed.
  function isAllowedCaller(
    address caller
  ) external view returns (bool allowed) {
    return s_allowedCallers.contains(caller);
  }

  /// @notice Returns an array of all allowed callers.
  /// @return allowedCallers An array of allowed caller addresses.
  function getAllowedCallers() external view returns (address[] memory allowedCallers) {
    return s_allowedCallers.values();
  }
}
