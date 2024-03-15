// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IL1ScrollMessenger} from "@scroll-tech/contracts/L1/IL1ScrollMessenger.sol";

contract MockScrollL1CrossDomainMessenger is IL1ScrollMessenger {
  uint256 private s_nonce;

  function xDomainMessageSender() public pure returns (address) {
    return address(0);
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

  function relayMessageWithProof(
    address from,
    address to,
    uint256 value,
    uint256 nonce,
    bytes memory message,
    L2MessageProof memory proof
  ) external override {}

  function replayMessage(
    address from,
    address to,
    uint256 value,
    uint256 messageNonce,
    bytes memory message,
    uint32 newGasLimit,
    address refundAddress
  ) external payable override {}

  function dropMessage(
    address from,
    address to,
    uint256 value,
    uint256 messageNonce,
    bytes memory message
  ) external override {}
}
