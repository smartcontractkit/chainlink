// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

/* Interface Imports */
import {IL2CrossDomainMessenger} from "@eth-optimism/contracts/L2/messaging/IL2CrossDomainMessenger.sol";

contract MockOptimismL2CrossDomainMessenger is IL2CrossDomainMessenger {
  uint256 private s_nonce;
  address private s_sender;

  // slither-disable-next-line external-function
  function xDomainMessageSender() public view returns (address) {
    return s_sender;
  }

  function setSender(address newSender) external {
    s_sender = newSender;
  }

  function sendMessage(address _target, bytes memory _message, uint32 _gasLimit) public {
    emit SentMessage(_target, msg.sender, _message, s_nonce, _gasLimit);
    s_nonce++;
  }

  function relayMessage(address _target, address _sender, bytes memory _message, uint256 _messageNonce) external {}

  receive() external payable {}
}
