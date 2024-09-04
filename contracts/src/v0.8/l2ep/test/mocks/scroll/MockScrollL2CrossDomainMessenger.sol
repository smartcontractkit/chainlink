// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IL2ScrollMessenger} from "@scroll-tech/contracts/L2/IL2ScrollMessenger.sol";

contract MockScrollL2CrossDomainMessenger is IL2ScrollMessenger {
  uint256 private s_nonce;
  address private s_sender;

  function xDomainMessageSender() public view returns (address) {
    return s_sender;
  }

  function sendMessage(
    address _target,
    uint256 _value,
    bytes calldata _message,
    uint256 _gasLimit
  ) external payable override {
    emit SentMessage(msg.sender, _target, _value, s_nonce, _gasLimit, _message);
    s_nonce++;
  }

  function sendMessage(
    address _target,
    uint256 _value,
    bytes calldata _message,
    uint256 _gasLimit,
    address
  ) external payable override {
    emit SentMessage(msg.sender, _target, _value, s_nonce, _gasLimit, _message);
    s_nonce++;
  }

  function relayMessage(
    address from,
    address to,
    uint256 value,
    uint256 nonce,
    bytes calldata message
  ) external override {}

  /// Needed for backwards compatibility in Hardhat tests
  function setSender(address newSender) external {
    s_sender = newSender;
  }

  /// Needed for backwards compatibility in Hardhat tests
  receive() external payable {}
}
