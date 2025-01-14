// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../../interfaces/IRouterClient.sol";
import {Client} from "../../libraries/Client.sol";
import {CCIPReceiver} from "./CCIPReceiver.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {EnumerableMap} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableMap.sol";

/// @title CCIPReceiverWithACK
/// @notice Acts as a CCIP receiver, but upon receiving an incoming message, attempts to send a response back to the
/// sender with an ACK packet indicating they received and processed the initial request correctly.
/// @dev Messages received by this contract must be of special formatting in which any arbitrary data is first wrapped
/// inside a MessagePayload struct, and must be processed first to ensure conformity.
contract CCIPReceiverWithACK is CCIPReceiver {
  using SafeERC20 for IERC20;
  using EnumerableMap for EnumerableMap.Bytes32ToUintMap;

  error InvalidAckMessageHeader();
  error MessageAlreadyAcknowledged(bytes32 messageId);

  event MessageAckSent(bytes32 incomingMessageId);
  event MessageSent(bytes32 indexed incomingMessageId, bytes32 indexed ACKMessageId);

  event MessageAckReceived(bytes32 messageId);
  event FeeTokenUpdated(address indexed oldToken, address indexed newToken);

  enum MessageType {
    OUTGOING, // Indicates that a message is being sent for the first time to its recipient.
    ACK // Indicates that another message of type "OUTGOING" has already been received, and an acknowledgement is being
      // returned to the original sender, by the original recipient.

  }

  enum MessageStatus {
    QUIET, // A message which has not been sent yet, the default status for any messageId
    SENT, // Indicates a message has been sent through CCIP but not yet received an ACK response from the recipient
    ACKNOWLEDGED // The original SENT message was received and processed by the recipient, and confirmation of
      // reception was received by this via the returned ACK message sent in response.

  }

  struct MessagePayload {
    bytes version; // An optional byte string which can be used to denote the ACK version formatting or how to decode remaining data.
    bytes data; // The Arbitrary data initially meant to be received by this contract and sent from the source chain.
    MessageType messageType; // Denotes whether the incoming message is being received for the first time, or is an
      // acknowledgement that the initial outgoing correspondence was successfully received.
  }

  // keccak256("MESSAGE_ACKNOWLEDGED_)"
  bytes32 public constant ACK_MESSAGE_HEADER = 0x1c778f21871bcc06cfebd177c4d0360c2f3550962fb071f69ed007e4f55f23b2;

  // Current feeToken
  IERC20 internal s_feeToken;

  mapping(bytes32 messageId => MessageStatus status) public s_messageStatus;

  constructor(address router, IERC20 feeToken) CCIPReceiver(router) {
    s_feeToken = feeToken;

    // If fee token is in LINK, then approve router to transfer
    if (address(feeToken) != address(0)) {
      feeToken.safeIncreaseAllowance(router, type(uint256).max);
    }
  }

  function updateFeeToken(address token) external onlyOwner {
    // If the current fee token is not-native, zero out the allowance to the router for safety
    if (address(s_feeToken) != address(0)) {
      s_feeToken.safeApprove(getRouter(), 0);
    }

    address oldFeeToken = address(s_feeToken);
    s_feeToken = IERC20(token);

    // Approve the router to spend the new fee token
    if (token != address(0)) {
      s_feeToken.safeIncreaseAllowance(getRouter(), type(uint256).max);
    }

    emit FeeTokenUpdated(oldFeeToken, token);
  }

  /// @notice Application-specific logic for incoming ccip messages.
  /// @dev Function does NOT require the status of an incoming ACK be "sent" because this implementation does not send, only receives
  /// Any MessageType encoding is implemented by the sender contract, and is not natively part of CCIP messages.
  function processMessage(
    Client.Any2EVMMessage calldata message
  ) external virtual override onlySelf isValidSender(message.sourceChainSelector, message.sender) {
    (MessagePayload memory payload) = abi.decode(message.data, (MessagePayload));

    // message type is a concept with ClientWithACK
    if (payload.messageType == MessageType.OUTGOING) {
      _processIncomingMessage(message);

      // If the message was outgoing on the source chain, then send an ack response.
      _sendAck(message);
      return;
    } else if (payload.messageType == MessageType.ACK) {
      // Decode message into the message header and the messageId to ensure the message is encoded correctly
      (bytes32 messageHeader, bytes32 messageId) = abi.decode(payload.data, (bytes32, bytes32));

      // Ensure Ack Message contains proper message header. Must abi.encode() before hashing since its of the string type
      if (messageHeader != ACK_MESSAGE_HEADER) {
        revert InvalidAckMessageHeader();
      }

      // Make sure the ACK message has not already been acknowledged
      if (s_messageStatus[messageId] == MessageStatus.ACKNOWLEDGED) revert MessageAlreadyAcknowledged(messageId);

      // Mark the message has finalized from a proper ack-message.
      s_messageStatus[messageId] = MessageStatus.ACKNOWLEDGED;

      emit MessageAckReceived(messageId);
    }
  }

  /// @notice Contains the arbitrary logic for processing incoming messages from an authorized sender & source-chain
  function _processIncomingMessage(Client.Any2EVMMessage calldata incomingMessage) internal virtual {}

  /// @notice Sends the acknowledgement message back through CCIP to original sender contract
  function _sendAck(Client.Any2EVMMessage calldata incomingMessage) internal {
    // Build the outgoing ACK message, with no tokens, with data being the concatenation of the acknowledgement header
    // and incoming-messageId
    Client.EVM2AnyMessage memory outgoingMessage = Client.EVM2AnyMessage({
      receiver: incomingMessage.sender,
      data: abi.encode(ACK_MESSAGE_HEADER, incomingMessage.messageId),
      tokenAmounts: new Client.EVMTokenAmount[](0),
      extraArgs: s_chainConfigs[incomingMessage.sourceChainSelector].extraArgsBytes,
      feeToken: address(s_feeToken)
    });

    uint256 feeAmount = IRouterClient(s_ccipRouter).getFee(incomingMessage.sourceChainSelector, outgoingMessage);

    bytes32 ACKMessageId = IRouterClient(s_ccipRouter).ccipSend{
      value: address(s_feeToken) == address(0) ? feeAmount : 0
    }(incomingMessage.sourceChainSelector, outgoingMessage);

    emit MessageSent(incomingMessage.messageId, ACKMessageId);
  }

  /// @notice returns the address of the fee token.
  /// @dev the zero address indicates to pay fees in native-tokens instead.
  function getFeeToken() public view virtual returns (address feeToken) {
    return address(s_feeToken);
  }
}
