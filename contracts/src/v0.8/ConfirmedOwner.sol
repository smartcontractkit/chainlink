// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./ConfirmedOwnerWithProposal.sol";

/**
 * @title The ConfirmedOwner contract
 * @notice A contract with helpers for basic contract ownership.
 */
contract ConfirmedOwner is ConfirmedOwnerWithProposal {
  // solhint-disable-next-line no-empty-blocks
  constructor(address newOwner) ConfirmedOwnerWithProposal(newOwner, address(0)) {}
}
