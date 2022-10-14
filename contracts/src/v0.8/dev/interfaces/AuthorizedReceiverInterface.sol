// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface AuthorizedReceiverInterface {
  function isAuthorizedSender(address sender) external view returns (bool);

  function getAuthorizedSenders() external returns (address[] memory);

  function setAuthorizedSenders(address[] calldata senders) external;
}
