// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../interfaces/IRouterClient.sol";

import {Client} from "../libraries/Client.sol";
import {CCIPClientExample} from "./CCIPClientExample.sol";

import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/token/ERC20/utils/SafeERC20.sol";
import {EnumerableMap} from "../../vendor/openzeppelin-solidity/v4.8.0/utils/structs/EnumerableMap.sol";

contract DefensiveExample is CCIPClientExample {
  using EnumerableMap for EnumerableMap.Bytes32ToUintMap;
  using SafeERC20 for IERC20;

  error OnlySelf();
  error ErrorCase();
  error MessageNotFailed(bytes32 messageId);

  event MessageFailed(bytes32 indexed messageId, bytes reason);
  event MessageSucceeded(bytes32 indexed messageId);
  event MessageRecovered(bytes32 indexed messageId);

  // Example error code, could have many different error codes.
  enum ErrorCode {
    // RESOLVED is first so that the default value is resolved.
    RESOLVED,
    // Could have any number of error codes here.
    BASIC
  }

  // The message contents of failed messages are stored here.
  mapping(bytes32 messageId => Client.Any2EVMMessage contents) public s_messageContents;

  // Contains failed messages and their state.
  EnumerableMap.Bytes32ToUintMap internal s_failedMessages;

  // This is used to simulate a revert in the processMessage function.
  bool internal s_simRevert = false;

  constructor(IRouterClient router, IERC20 feeToken) CCIPClientExample(router, feeToken) {}

  /// @notice The entrypoint for the CCIP router to call. This function should
  /// never revert, all errors should be handled internally in this contract.
  /// @param message The message to process.
  /// @dev Extremely important to ensure only router calls this.
  function ccipReceive(
    Client.Any2EVMMessage calldata message
  ) external override onlyRouter validChain(message.sourceChainSelector) {
    try this.processMessage(message) {} catch (bytes memory err) {
      // Could set different error codes based on the caught error. Each could be
      // handled differently.
      s_failedMessages.set(message.messageId, uint256(ErrorCode.BASIC));
      s_messageContents[message.messageId] = message;
      // Don't revert so CCIP doesn't revert. Emit event instead.
      // The message can be retried later without having to do manual execution of CCIP.
      emit MessageFailed(message.messageId, err);
      return;
    }
    emit MessageSucceeded(message.messageId);
  }

  /// @notice This function the entrypoint for this contract to process messages.
  /// @param message The message to process.
  /// @dev This example just sends the tokens to the owner of this contracts. More
  /// interesting functions could be implemented.
  /// @dev It has to be external because of the try/catch.
  function processMessage(
    Client.Any2EVMMessage calldata message
  ) external onlySelf validChain(message.sourceChainSelector) {
    // Simulate a revert
    if (s_simRevert) revert ErrorCase();

    // Send tokens to the owner
    for (uint256 i = 0; i < message.destTokenAmounts.length; ++i) {
      IERC20(message.destTokenAmounts[i].token).safeTransfer(owner(), message.destTokenAmounts[i].amount);
    }
    // Do other things that might revert
  }

  /// @notice This function is callable by the owner when a message has failed
  /// to unblock the tokens that are associated with that message.
  /// @dev This function is only callable by the owner.
  function retryFailedMessage(bytes32 messageId, address tokenReceiver) external onlyOwner {
    if (s_failedMessages.get(messageId) != uint256(ErrorCode.BASIC)) revert MessageNotFailed(messageId);
    // Set the error code to 0 to disallow reentry and retry the same failed message
    // multiple times.
    s_failedMessages.set(messageId, uint256(ErrorCode.RESOLVED));

    // Do stuff to retry message, potentially just releasing the associated tokens
    Client.Any2EVMMessage memory message = s_messageContents[messageId];

    // send the tokens to the receiver as escape hatch
    for (uint256 i = 0; i < message.destTokenAmounts.length; ++i) {
      IERC20(message.destTokenAmounts[i].token).safeTransfer(tokenReceiver, message.destTokenAmounts[i].amount);
    }

    emit MessageRecovered(messageId);
  }

  // An example function to demonstrate recovery
  function setSimRevert(bool simRevert) external onlyOwner {
    s_simRevert = simRevert;
  }

  modifier onlySelf() {
    if (msg.sender != address(this)) revert OnlySelf();
    _;
  }
}
