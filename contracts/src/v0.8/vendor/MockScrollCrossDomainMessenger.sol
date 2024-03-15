// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

import "./openzeppelin-solidity/v4.8.3/contracts/utils/Address.sol";

/// sourced from: https://github.com/scroll-tech/scroll/blob/develop/contracts/src/libraries/IScrollMessenger.sol
interface IScrollMessenger {
  /// **********
  /// * Events *
  /// **********

  /// @notice Emitted when a cross domain message is sent.
  /// @param sender The address of the sender who initiates the message.
  /// @param target The address of target contract to call.
  /// @param value The amount of value passed to the target contract.
  /// @param messageNonce The nonce of the message.
  /// @param gasLimit The optional gas limit passed to L1 or L2.
  /// @param message The calldata passed to the target contract.
  event SentMessage(
    address indexed sender,
    address indexed target,
    uint256 value,
    uint256 messageNonce,
    uint256 gasLimit,
    bytes message
  );

  /// @notice Emitted when a cross domain message is relayed successfully.
  /// @param messageHash The hash of the message.
  event RelayedMessage(bytes32 indexed messageHash);

  /// @notice Emitted when a cross domain message is failed to relay.
  /// @param messageHash The hash of the message.
  event FailedRelayedMessage(bytes32 indexed messageHash);

  ///  *************************
  ///  * Public View Functions *
  ///  *************************

  /// @notice Return the sender of a cross domain message.
  function xDomainMessageSender() external view returns (address);

  /// *****************************
  /// * Public Mutating Functions *
  /// *****************************

  /// @notice Send cross chain message from L1 to L2 or L2 to L1.
  /// @param target The address of account who receive the message.
  /// @param value The amount of ether passed when call target contract.
  /// @param message The content of the message.
  /// @param gasLimit Gas limit required to complete the message relay on corresponding chain.
  function sendMessage(address target, uint256 value, bytes calldata message, uint256 gasLimit) external payable;

  /// @notice Send cross chain message from L1 to L2 or L2 to L1.
  /// @param target The address of account who receive the message.
  /// @param value The amount of ether passed when call target contract.
  /// @param message The content of the message.
  /// @param gasLimit Gas limit required to complete the message relay on corresponding chain.
  /// @param refundAddress The address of account who will receive the refunded fee.
  function sendMessage(
    address target,
    uint256 value,
    bytes calldata message,
    uint256 gasLimit,
    address refundAddress
  ) external payable;
}

contract MockScrollCrossDomainMessenger is IScrollMessenger {
  address internal mockMessageSender;

  constructor(address sender) {
    mockMessageSender = sender;
  }

  function xDomainMessageSender() external view override returns (address) {
    return mockMessageSender;
  }

  function _setMockMessageSender(address sender) external {
    mockMessageSender = sender;
  }

  /// *****************************
  /// * Public Mutating Functions *
  /// *****************************

  /// @notice Send cross chain message from L1 to L2 or L2 to L1.
  /// @param _target The address of account who receive the message.
  /// @param _message The content of the message.
  function sendMessage(address _target, uint256, bytes calldata _message, uint256) external payable override {
    Address.functionCall(_target, _message, "sendMessage reverted");
  }

  /// @notice Send cross chain message from L1 to L2 or L2 to L1.
  /// @param _target The address of account who receive the message.
  /// @param _message The content of the message.
  function sendMessage(address _target, uint256, bytes calldata _message, uint256, address) external payable override {
    Address.functionCall(_target, _message, "sendMessage reverted");
  }
}
