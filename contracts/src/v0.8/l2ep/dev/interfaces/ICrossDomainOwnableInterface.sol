// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

/// @title ICrossDomainOwnableInterface - A contract with helpers for cross-domain contract ownership
interface ICrossDomainOwnableInterface {
  event L1OwnershipTransferRequested(address indexed from, address indexed to);

  event L1OwnershipTransferred(address indexed from, address indexed to);

  function l1Owner() external returns (address);

  function transferL1Ownership(address recipient) external;

  function acceptL1Ownership() external;
}
