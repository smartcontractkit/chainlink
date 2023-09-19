// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IOwnable} from "../../../../shared/interfaces/IOwnable.sol";
import {Initializable} from "../../../../vendor/openzeppelin-contracts-upgradeable/v4.8.1/proxy/utils/Initializable.sol";

/**
 * @title The ConfirmedOwnerUpgradeable contract
 * @notice An upgrade compatible contract with helpers for basic contract ownership.
 */
contract ConfirmedOwnerUpgradeable is Initializable, IOwnable {
  address private s_owner;
  address private s_pendingOwner;

  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);

  error OwnerMustBeSet();
  error NotProposedOwner();
  error CannotSelfTransfer();
  error OnlyCallableByOwner();

  /**
   * @dev Initializes the contract in unpaused state.
   */
  function __ConfirmedOwner_initialize(address newOwner, address pendingOwner) internal onlyInitializing {
    if (newOwner == address(0)) {
      revert OwnerMustBeSet();
    }

    s_owner = newOwner;
    if (pendingOwner != address(0)) {
      _transferOwnership(pendingOwner);
    }
  }

  /**
   * @notice Allows an owner to begin transferring ownership to a new address,
   * pending.
   */
  function transferOwnership(address to) public override onlyOwner {
    _transferOwnership(to);
  }

  /**
   * @notice Allows an ownership transfer to be completed by the recipient.
   */
  function acceptOwnership() external override {
    if (msg.sender != s_pendingOwner) {
      revert NotProposedOwner();
    }

    address oldOwner = s_owner;
    s_owner = msg.sender;
    s_pendingOwner = address(0);

    emit OwnershipTransferred(oldOwner, msg.sender);
  }

  /**
   * @notice Get the current owner
   */
  function owner() public view override returns (address) {
    return s_owner;
  }

  /**
   * @notice validate, transfer ownership, and emit relevant events
   */
  function _transferOwnership(address to) private {
    if (to == msg.sender) {
      revert CannotSelfTransfer();
    }

    s_pendingOwner = to;

    emit OwnershipTransferRequested(s_owner, to);
  }

  /**
   * @notice validate access
   */
  function _validateOwnership() internal view {
    if (msg.sender != s_owner) {
      revert OnlyCallableByOwner();
    }
  }

  /**
   * @notice Reverts if called by anyone other than the contract owner.
   */
  modifier onlyOwner() {
    _validateOwnership();
    _;
  }
}
