// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../libraries/Client.sol";
import {CCIPBase} from "./CCIPBase.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableSet} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";

/// @title CCIPReceiver
/// @notice This contract is capable of receiving incoming messages from the CCIP Router.
/// @dev This contract implements various "defensive" measures to enhance security and efficiency. These include the
/// ability to implement custom-retry logic and ownership-based token-recovery functions.
contract CCIPReceiver is CCIPBase {
  using SafeERC20 for IERC20;
  using EnumerableSet for EnumerableSet.Bytes32Set;

  error OnlySelf();
  error MessageNotFailed(bytes32 messageId);

  event MessageFailed(bytes32 indexed messageId, bytes reason);
  event MessageSucceeded(bytes32 indexed messageId);
  event MessageRecovered(bytes32 indexed messageId);
  event MessageAbandoned(bytes32 indexed messageId, address tokenReceiver);

  mapping(bytes32 messageId => Client.Any2EVMMessage contents) internal s_messageContents;

  // Contains the set of all messages in s_messageContents which failed to process properly.
  // When a message is retried or abandoned it is removed from this set.
  EnumerableSet.Bytes32Set internal s_failedMessages;

  constructor(address router) CCIPBase(router) {}

  // ================================================================
  // │                  Incoming Message Processing                 |
  // ================================================================

  /// @notice The entrypoint for the CCIP router to call. This function should
  /// not revert, all errors should be handled internally in this contract.
  /// @param message The message to process.
  /// @dev Extremely important to ensure only router calls this.
  function ccipReceive(
    Client.Any2EVMMessage calldata message
  ) external virtual onlyRouter isValidChain(message.sourceChainSelector) {
    try this.processMessage(message) {}
    catch (bytes memory err) {
      // If custom retry logic is desired, plus granting the owner the ability to extract tokens as a last resort for
      // recovery, use this try-catch pattern in ccipReceiver. It will make the message appear as a success to CCIP, and
      // actual message state and any residual errors can be tracked within the dapp.
      s_failedMessages.add(message.messageId);

      // Store the message contents in case it needs to be retried or abandoned
      s_messageContents[message.messageId] = message;

      // Don't revert because CCIPRouter doesn't revert. Emit event instead.
      // The message can be retried or abandoned later without having to do manual execution of CCIP, which should
      // be reserved for retrying with a higher gas limit.
      emit MessageFailed(message.messageId, err);
      return;
    }

    emit MessageSucceeded(message.messageId);
  }

  /// @notice Contains arbitrary application-logic for incoming CCIP messages.
  /// @dev It has to be external because of the try/catch of ccipReceive() which invokes it
  function processMessage(
    Client.Any2EVMMessage calldata message
  ) external virtual onlySelf isValidSender(message.sourceChainSelector, message.sender) {}

  // ================================================================
  // │                  Failed Message Processing                   |
  // ================================================================

  /// @notice Executes a message that failed initial delivery, but with different logic specifically for re-execution.
  /// @dev Since the function invoked _retryFailedMessage(), which is marked onlyOwner, this may only be called by the
  /// Owner as well. The function will revert if the messageId was not already stored as failed during initial execution
  /// @param messageId the unique ID of the CCIP-message which failed initial-execution.
  function retryFailedMessage(bytes32 messageId) external virtual {
    if (!s_failedMessages.contains(messageId)) revert MessageNotFailed(messageId);

    // Allow developer to implement arbitrary functionality on retried messages, such as just releasing the associated
    // tokens
    Client.Any2EVMMessage memory message = s_messageContents[messageId];

    // Set remove the message from storage to disallow reentry and retry the same failed message multiple times.
    delete s_messageContents[messageId];
    s_failedMessages.remove(messageId);

    // Allow the user override the implementation, since different workflow may be desired for retrying a message
    _retryFailedMessage(message);

    emit MessageRecovered(messageId);
  }

  /// @notice A function that should contain any special logic needed to "retry" processing of a previously failed message.
  /// @dev If the owner wants to retrieve tokens without special logic, then abandonFailedMessage(), withdrawNativeTokens(), or withdrawTokens() should be used instead
  /// This function is marked onlyOwner, but is virtual. Allowing permissionless execution is not recommended but may be allowed if function is overridden
  function _retryFailedMessage(Client.Any2EVMMessage memory message) internal virtual onlyOwner {
    this.processMessage(message);
  }

  /// @notice Should be used to recover tokens from a failed message, while ensuring the message cannot be retried
  /// @dev function will send tokens to destination, but will NOT invoke any arbitrary logic afterwards.
  /// function is only callable by the contract owner
  function abandonFailedMessage(bytes32 messageId, address receiver) external onlyOwner {
    if (!s_failedMessages.contains(messageId)) revert MessageNotFailed(messageId);

    Client.EVMTokenAmount[] memory tokenAmounts = s_messageContents[messageId].destTokenAmounts;

    // Follow CEI and remove failed message from state before transferring in case of ERC-667 external calls
    delete s_messageContents[messageId];
    s_failedMessages.remove(messageId);

    for (uint256 i = 0; i < tokenAmounts.length; ++i) {
      IERC20(tokenAmounts[i].token).safeTransfer(receiver, tokenAmounts[i].amount);
    }

    emit MessageAbandoned(messageId, receiver);
  }

  // ================================================================
  // │                  Message Tracking                            │
  // ================================================================

  /// @param messageId the ID of the message delivered by the CCIP Router
  /// @return Any2EVMMessage a standard CCIP message for EVM-compatible networks
  function getMessageContents(bytes32 messageId) public view returns (Client.Any2EVMMessage memory) {
    return s_messageContents[messageId];
  }

  /// @notice Retrieve whether a message delivered by the CCIP router failed to process properly.
  /// @dev Querying this function with message which was successfully retried or abandoned will return false
  /// @param messageId the ID of the message delivered by the CCIP Router
  /// @return bool Whether the previously-delivered message failed to process.
  function isFailedMessage(bytes32 messageId) public view returns (bool) {
    return s_failedMessages.contains(messageId);
  }

  modifier onlySelf() {
    if (msg.sender != address(this)) revert OnlySelf();
    _;
  }
}
