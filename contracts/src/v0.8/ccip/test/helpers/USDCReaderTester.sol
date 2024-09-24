// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

contract USDCReaderTester {
  event MessageSent(bytes);

  // emitMessageSent reflects the logic from Circle's MessageTransmitter emitting MeseageSent(bytes) events
  // https://github.com/circlefin/evm-cctp-contracts/blob/377c9bd813fb86a42d900ae4003599d82aef635a/src/MessageTransmitter.sol#L41
  // https://github.com/circlefin/evm-cctp-contracts/blob/377c9bd813fb86a42d900ae4003599d82aef635a/src/MessageTransmitter.sol#L365
  function emitMessageSent(
    uint32 version,
    uint32 sourceDomain,
    uint32 destinationDomain,
    bytes32 recipient,
    bytes32 destinationCaller,
    bytes32 sender,
    uint64 nonce,
    bytes calldata messageBody
  ) external {
    bytes memory _message =
      _formatMessage(version, sourceDomain, destinationDomain, nonce, sender, recipient, destinationCaller, messageBody);
    emit MessageSent(_message);
  }

  /**
   * @notice Returns formatted (packed) message with provided fields
   * It's a copy paste of the Message._formatMessage() call in MessageTransmitter.sol
   * https://github.com/circlefin/evm-cctp-contracts/blob/377c9bd813fb86a42d900ae4003599d82aef635a/src/messages/Message.sol#L54C1-L65C9
   * Check the chainlink-ccip repo for the offchain implementation of matching this format
   * @param _msgVersion the version of the message format
   * @param _msgSourceDomain Domain of home chain
   * @param _msgDestinationDomain Domain of destination chain
   * @param _msgNonce Destination-specific nonce
   * @param _msgSender Address of sender on source chain as bytes32
   * @param _msgRecipient Address of recipient on destination chain as bytes32
   * @param _msgDestinationCaller Address of caller on destination chain as bytes32
   * @param _msgRawBody Raw bytes of message body
   * @return Formatted message
   *
   */
  function _formatMessage(
    uint32 _msgVersion,
    uint32 _msgSourceDomain,
    uint32 _msgDestinationDomain,
    uint64 _msgNonce,
    bytes32 _msgSender,
    bytes32 _msgRecipient,
    bytes32 _msgDestinationCaller,
    bytes memory _msgRawBody
  ) internal pure returns (bytes memory) {
    return abi.encodePacked(
      _msgVersion,
      _msgSourceDomain,
      _msgDestinationDomain,
      _msgNonce,
      _msgSender,
      _msgRecipient,
      _msgDestinationCaller,
      _msgRawBody
    );
  }
}
