// SPDX-License-Identifier: MIT
pragma solidity ^0.8.16;

import {IScrollMessenger} from "@scroll-tech/contracts/libraries/IScrollMessenger.sol";

import {Address} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/Address.sol";

contract MockScrollCrossDomainMessenger is IScrollMessenger {
  address internal s_mockMessageSender;

  constructor(address sender) {
    s_mockMessageSender = sender;
  }

  function xDomainMessageSender() external view override returns (address) {
    return s_mockMessageSender;
  }

  function _setMockMessageSender(address sender) external {
    s_mockMessageSender = sender;
  }

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
