// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IOwnable} from "../interfaces/IOwnable.sol";

/// @notice A minimal contract that implements 2-step ownership transfer and nothing more. It's made to be minimal
/// to reduce the impact of the bytecode size on any contract that inherits from it.
contract Ownable2Step is IOwnable {
  /// @notice The pending owner is the address to which ownership may be transferred.
  address private s_pendingOwner;
  /// @notice The owner is the current owner of the contract.
  /// @dev The owner is the second storage variable so any implementing contract could pack other state with it
  /// instead of the much less used s_pendingOwner.
  address private s_owner;

  error OwnerCannotBeZero();
  error MustBeProposedOwner();
  error CannotTransferToSelf();
  error OnlyCallableByOwner();

  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);

  constructor(address newOwner, address pendingOwner) {
    if (newOwner == address(0)) {
      revert OwnerCannotBeZero();
    }

    s_owner = newOwner;
    if (pendingOwner != address(0)) {
      _transferOwnership(pendingOwner);
    }
  }

  /// @notice Get the current owner
  function owner() public view override returns (address) {
    return s_owner;
  }

  /// @notice Allows an owner to begin transferring ownership to a new address. The new owner needs to call
  /// `acceptOwnership` to accept the transfer before any permissions are changed.
  /// @param to The address to which ownership will be transferred.
  function transferOwnership(address to) public override onlyOwner {
    _transferOwnership(to);
  }

  /// @notice validate, transfer ownership, and emit relevant events
  /// @param to The address to which ownership will be transferred.
  function _transferOwnership(address to) private {
    if (to == msg.sender) {
      revert CannotTransferToSelf();
    }

    s_pendingOwner = to;

    emit OwnershipTransferRequested(s_owner, to);
  }

  /// @notice Allows an ownership transfer to be completed by the recipient.
  function acceptOwnership() external override {
    if (msg.sender != s_pendingOwner) {
      revert MustBeProposedOwner();
    }

    address oldOwner = s_owner;
    s_owner = msg.sender;
    s_pendingOwner = address(0);

    emit OwnershipTransferred(oldOwner, msg.sender);
  }

  /// @notice validate access
  function _validateOwnership() internal view {
    if (msg.sender != s_owner) {
      revert OnlyCallableByOwner();
    }
  }

  /// @notice Reverts if called by anyone other than the contract owner.
  modifier onlyOwner() {
    _validateOwnership();
    _;
  }
}
