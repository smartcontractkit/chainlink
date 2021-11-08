// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../ConfirmedOwner.sol";
import "./interfaces/CrossDomainOwnableInterface.sol";

/**
 * @title The CrossDomainOwnable contract
 * @notice A contract with helpers for cross-domain contract ownership.
 */
abstract contract CrossDomainOwnable is CrossDomainOwnableInterface, ConfirmedOwner {
  address internal s_l1Owner;
  address internal s_l1PendingOwner;

  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  constructor(address newl1Owner) ConfirmedOwner(msg.sender) {
    _setL1Owner(newl1Owner);
  }

  /**
   * @notice Allows an owner to begin transferring ownership to a new address as a pending owner.
   */
  function transferL1Ownership(address to) external virtual override;

  /**
   * @notice Allows an ownership transfer to be completed by the recipient.
   * @dev The following has to be implemented per-chain because msg.sender is translated to a cross-domain messenger address.
   */
  function acceptL1Ownership() external virtual override;

  /**
   * @notice Get the current owner
   */
  function l1Owner() public view override returns (address) {
    return s_l1Owner;
  }

  /**
   * @notice validate, transfer ownership, and emit relevant events
   */
  function _transferL1Ownership(address to) internal {
    require(to != address(0), "Cannot transfer to zero address");
    require(to != msg.sender, "Cannot transfer to self");
    s_l1PendingOwner = to;

    emit L1OwnershipTransferRequested(s_l1Owner, to);
  }

  /**
   * @notice set ownership, emit relevant events. Used in acceptOwnership()
   */
  function _setL1Owner(address to) internal {
    address oldOwner = s_l1Owner;
    s_l1Owner = to;
    s_l1PendingOwner = address(0);

    emit L1OwnershipTransferred(oldOwner, to);
  }

  /**
   * @notice Reverts if called by anyone other than the L1 owner.
   */
  modifier onlyL1Owner() {
    require(msg.sender == s_l1Owner, "Only callable by L1 owner");
    _;
  }
}
