// SPDX-License-Identifier: MIT

pragma solidity >=0.7.6 <0.9.0;

contract MockOVMCrossDomainMessenger {
  address internal mockMessageSender;

  constructor(address sender) {
    mockMessageSender = sender;
  }

  function xDomainMessageSender() external view returns (address) {
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
  ) external {}
}
