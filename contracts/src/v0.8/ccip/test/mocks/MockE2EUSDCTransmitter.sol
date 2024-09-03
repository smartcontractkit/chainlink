/*
 * Copyright (c) 2022, Circle Internet Financial Limited.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
pragma solidity ^0.8.0;

import {IMessageTransmitterWithRelay} from "./interfaces/IMessageTransmitterWithRelay.sol";

import {BurnMintERC677} from "../../../shared/token/ERC677/BurnMintERC677.sol";

contract MockE2EUSDCTransmitter is IMessageTransmitterWithRelay {
  // Indicated whether the receiveMessage() call should succeed.
  bool public s_shouldSucceed;
  uint32 private immutable i_version;
  uint32 private immutable i_localDomain;
  // Next available nonce from this source domain
  uint64 public nextAvailableNonce;

  BurnMintERC677 internal immutable i_token;

  /**
   * @notice Emitted when a new message is dispatched
   * @param message Raw bytes of message
   */
  event MessageSent(bytes message);

  constructor(uint32 _version, uint32 _localDomain, address token) {
    i_version = _version;
    i_localDomain = _localDomain;
    s_shouldSucceed = true;

    i_token = BurnMintERC677(token);
  }

  /// @param message The original message on the source chain
  ///     * Message format:
  ///     * Field                 Bytes      Type       Index
  ///     * version               4          uint32     0
  ///     * sourceDomain          4          uint32     4
  ///     * destinationDomain     4          uint32     8
  ///     * nonce                 8          uint64     12
  ///     * sender                32         bytes32    20
  ///     * recipient             32         bytes32    52
  ///     * destinationCaller     32         bytes32    84
  ///     * messageBody           dynamic    bytes      116
  function receiveMessage(bytes calldata message, bytes calldata) external returns (bool success) {
    // The receiver of the funds is the _mintRecipient in the following encoded format
    //   function _formatMessage(
    //    uint32 _version,             4
    //    bytes32 _burnToken,         32
    //    bytes32 _mintRecipient,     32, first 12 empty for EVM addresses
    //    uint256 _amount,
    //    bytes32 _messageSender
    //  ) internal pure returns (bytes memory) {
    //    return abi.encodePacked(_version, _burnToken, _mintRecipient, _amount, _messageSender);
    //  }
    address recipient = address(bytes20(message[116 + 36 + 12:116 + 36 + 12 + 20]));
    // We always mint 1 token to not complicate the test.
    i_token.mint(recipient, 1);

    return s_shouldSucceed;
  }

  function setShouldSucceed(bool shouldSucceed) external {
    s_shouldSucceed = shouldSucceed;
  }

  function version() external view returns (uint32) {
    return i_version;
  }

  function localDomain() external view returns (uint32) {
    return i_localDomain;
  }

  /**
   * This is based on similar function in https://github.com/circlefin/evm-cctp-contracts/blob/master/src/MessageTransmitter.sol
   * @notice Send the message to the destination domain and recipient
   * @dev Increment nonce, format the message, and emit `MessageSent` event with message information.
   * @param destinationDomain Domain of destination chain
   * @param recipient Address of message recipient on destination chain as bytes32
   * @param messageBody Raw bytes content of message
   * @return nonce reserved by message
   */
  function sendMessage(
    uint32 destinationDomain,
    bytes32 recipient,
    bytes calldata messageBody
  ) external returns (uint64) {
    bytes32 _emptyDestinationCaller = bytes32(0);
    uint64 _nonce = _reserveAndIncrementNonce();
    bytes32 _messageSender = bytes32(uint256(uint160((msg.sender))));

    _sendMessage(destinationDomain, recipient, _emptyDestinationCaller, _messageSender, _nonce, messageBody);

    return _nonce;
  }

  /**
   * @notice Send the message to the destination domain and recipient, for a specified `destinationCaller` on the
   * destination domain.
   * @dev Increment nonce, format the message, and emit `MessageSent` event with message information.
   * WARNING: if the `destinationCaller` does not represent a valid address, then it will not be possible
   * to broadcast the message on the destination domain. This is an advanced feature, and the standard
   * sendMessage() should be preferred for use cases where a specific destination caller is not required.
   * @param destinationDomain Domain of destination chain
   * @param recipient Address of message recipient on destination domain as bytes32
   * @param destinationCaller caller on the destination domain, as bytes32
   * @param messageBody Raw bytes content of message
   * @return nonce reserved by message
   */
  function sendMessageWithCaller(
    uint32 destinationDomain,
    bytes32 recipient,
    bytes32 destinationCaller,
    bytes calldata messageBody
  ) external returns (uint64) {
    require(destinationCaller != bytes32(0), "Destination caller must be nonzero");

    uint64 _nonce = _reserveAndIncrementNonce();
    bytes32 _messageSender = bytes32(uint256(uint160((msg.sender))));

    _sendMessage(destinationDomain, recipient, destinationCaller, _messageSender, _nonce, messageBody);

    return _nonce;
  }

  /**
   * Reserve and increment next available nonce
   * @return nonce reserved
   */
  function _reserveAndIncrementNonce() internal returns (uint64) {
    uint64 _nonceReserved = nextAvailableNonce;
    nextAvailableNonce = nextAvailableNonce + 1;
    return _nonceReserved;
  }

  /**
   * @notice Send the message to the destination domain and recipient. If `_destinationCaller` is not equal to bytes32(0),
   * the message can only be received on the destination chain when called by `_destinationCaller`.
   * @dev Format the message and emit `MessageSent` event with message information.
   * @param _destinationDomain Domain of destination chain
   * @param _recipient Address of message recipient on destination domain as bytes32
   * @param _destinationCaller caller on the destination domain, as bytes32
   * @param _sender message sender, as bytes32
   * @param _nonce nonce reserved for message
   * @param _messageBody Raw bytes content of message
   */
  function _sendMessage(
    uint32 _destinationDomain,
    bytes32 _recipient,
    bytes32 _destinationCaller,
    bytes32 _sender,
    uint64 _nonce,
    bytes calldata _messageBody
  ) internal {
    require(_recipient != bytes32(0), "Recipient must be nonzero");
    // serialize message
    bytes memory _message = abi.encodePacked(
      i_version, i_localDomain, _destinationDomain, _nonce, _sender, _recipient, _destinationCaller, _messageBody
    );

    // Emit MessageSent event
    emit MessageSent(_message);
  }
}
