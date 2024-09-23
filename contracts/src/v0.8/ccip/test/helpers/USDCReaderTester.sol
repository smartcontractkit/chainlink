// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

contract USDCReaderTester {
  event MessageSent(bytes);

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
