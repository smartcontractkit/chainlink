// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface CrossDomainOwnableInterface {
  function l1Owner() external returns (address);

  function transferL1Ownership(address recipient) external;

  function acceptL1Ownership() external;
}
