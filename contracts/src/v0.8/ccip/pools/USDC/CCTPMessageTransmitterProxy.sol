// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransmitter} from "./IMessageTransmitter.sol";
import {ITokenMessenger} from "./ITokenMessenger.sol";

import {Ownable2StepMsgSender} from "../../../shared/access/Ownable2StepMsgSender.sol";

/// @title CCTP Message Transmitter Proxy
/// @notice A proxy contract for handling messages transmitted via the Cross Chain Transfer Protocol (CCTP).
/// @dev This contract is responsible for receiving messages from the `IMessageTransmitter` and ensuring only the Token Pool can invoke it.
contract CCTPMessageTransmitterProxy is Ownable2StepMsgSender {
  /// @notice Error thrown when a function is called by an unauthorized entity.
  error Unauthorized();

  struct AllowedCallerConfigParam {
    address caller;
    bool allowed;
  }

  /// @notice Immutable reference to the `IMessageTransmitter` contract.
  IMessageTransmitter public immutable i_cctpTransmitter;

  /// @notice Addresses allowed to call `receiveMessage`.
  mapping(address => bool) private s_allowedCallers;

  /// @notice One-time cyclic dependency between TokenPool and MessageTransmitter.
  constructor(
    ITokenMessenger tokenMessenger
  ) {
    i_cctpTransmitter = IMessageTransmitter(tokenMessenger.localMessageTransmitter());
  }

  /// @notice Receives a message from the `IMessageTransmitter` contract and validates it.
  /// @dev Can only be called by the Token Pool to process incoming messages.
  /// @param message The payload of the message being received.
  /// @param attestation The cryptographic proof validating the message.
  /// @return success A boolean indicating if the message was successfully processed.
  function receiveMessage(
    bytes calldata message,
    bytes calldata attestation
  ) external onlyAllowedCaller returns (bool success) {
    return i_cctpTransmitter.receiveMessage(message, attestation);
  }

  /// @notice Configures the allowed callers for the `receiveMessage` function.
  /// @param params An array of `AllowedCallerConfigParam` structs.
  function configureAllowedCallers(
    AllowedCallerConfigParam[] calldata params
  ) external onlyOwner {
    for (uint256 i = 0; i < params.length; i++) {
      s_allowedCallers[params[i].caller] = params[i].allowed;
    }
  }

  /// @notice Checks if the caller is allowed to call the `receiveMessage` function.
  /// @return allowed A boolean indicating if the caller is allowed to call the function.
  function isAllowedCaller(
    address caller
  ) external view returns (bool allowed) {
    return s_allowedCallers[caller];
  }

  /// @notice Ensures that only the authorized Token Pool can call certain functions.
  modifier onlyAllowedCaller() {
    if (!s_allowedCallers[msg.sender]) {
      revert Unauthorized();
    }
    _;
  }
}
