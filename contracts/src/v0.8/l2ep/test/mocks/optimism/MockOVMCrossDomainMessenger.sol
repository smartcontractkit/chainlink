// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {iOVM_CrossDomainMessenger} from "../../../../vendor/@eth-optimism/contracts/v0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";

import {Address} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/Address.sol";

contract MockOVMCrossDomainMessenger is iOVM_CrossDomainMessenger {
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

  /**
   * Sends a cross domain message to the target messenger.
   * @param _target Target contract address.
   * @param _message Message to send to the target.
   */
  function sendMessage(address _target, bytes calldata _message, uint32) external override {
    Address.functionCall(_target, _message, "sendMessage reverted");
  }
}
