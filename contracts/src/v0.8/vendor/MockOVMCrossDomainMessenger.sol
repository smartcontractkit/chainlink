// SPDX-License-Identifier: MIT

pragma solidity >=0.7.6 <0.9.0;

import "./openzeppelin-solidity/v4.7.0/contracts/utils/Address.sol";

/**
 * @title iOVM_CrossDomainMessenger
 */
interface iOVM_CrossDomainMessenger {
  /**********
   * Events *
   **********/

  event SentMessage(bytes message);
  event RelayedMessage(bytes32 msgHash);
  event FailedRelayedMessage(bytes32 msgHash);

  /*************
   * Variables *
   *************/

  function xDomainMessageSender() external view returns (address);

  /********************
   * Public Functions *
   ********************/

  /**
   * Sends a cross domain message to the target messenger.
   * @param _target Target contract address.
   * @param _message Message to send to the target.
   * @param _gasLimit Gas limit for the provided message.
   */
  function sendMessage(
    address _target,
    bytes calldata _message,
    uint32 _gasLimit
  ) external;
}

contract MockOVMCrossDomainMessenger is iOVM_CrossDomainMessenger{
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

  /********************
   * Public Functions *
   ********************/

  /**
   * Sends a cross domain message to the target messenger.
   * @param _target Target contract address.
   * @param _message Message to send to the target.
   * @param _gasLimit Gas limit for the provided message.
   */
  function sendMessage(
    address _target,
    bytes calldata _message,
    uint32 _gasLimit
  ) external override {
    Address.functionCall(_target, _message, "sendMessage reverted");
  }
}
