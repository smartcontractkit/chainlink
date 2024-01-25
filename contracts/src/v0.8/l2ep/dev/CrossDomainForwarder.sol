// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CrossDomainOwnableInterface} from "./interfaces/CrossDomainOwnableInterface.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {ForwarderInterface} from "./interfaces/ForwarderInterface.sol";

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";

import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/// @title CrossDomainForwarder - L1 xDomain account representation
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
///   can consider that position to be held by the `l1Owner`
abstract contract CrossDomainForwarder is
  ITypeAndVersion,
  ForwarderInterface,
  CrossDomainOwnableInterface,
  ConfirmedOwner
{
  address internal s_l1Owner;
  address internal s_l1PendingOwner;

  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(address l1OwnerAddr) ConfirmedOwner(msg.sender) {
    _setL1Owner(l1OwnerAddr);
  }

  /// @notice Reverts if called by anyone other than the L1 owner.
  modifier onlyL1Owner() virtual {
    // solhint-disable-next-line custom-errors
    require(msg.sender == s_l1Owner, "Only callable by L1 owner");
    _;
  }

  /// @notice Reverts if called by anyone other than the L1 owner.
  modifier onlyProposedL1Owner() virtual {
    // solhint-disable-next-line custom-errors
    require(msg.sender == s_l1PendingOwner, "Only callable by proposed L1 owner");
    _;
  }

  /// @notice The address of the Cross Domain Messenger contract
  function crossDomainMessenger() external view virtual returns (address);

  /// @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
  /// @inheritdoc ForwarderInterface
  function forward(address target, bytes memory data) external virtual override onlyL1Owner {
    Address.functionCall(target, data, "Forwarder call reverted");
  }

  /// @notice transfer ownership of this account to a new L1 owner
  /// @param to new L1 owner that will be allowed to call the forward fn
  function transferL1Ownership(address to) public virtual override onlyL1Owner {
    _transferL1Ownership(to);
  }

  /// @notice accept ownership of this account to a new L1 owner
  function acceptL1Ownership() public virtual override onlyProposedL1Owner {
    _setL1Owner(s_l1PendingOwner);
  }

  /// @notice Get the current owner
  function l1Owner() public view override returns (address) {
    return s_l1Owner;
  }

  /// @notice validate, transfer ownership, and emit relevant events
  function _transferL1Ownership(address to) internal {
    // solhint-disable-next-line custom-errors
    require(to != msg.sender, "Cannot transfer to self");

    s_l1PendingOwner = to;

    emit L1OwnershipTransferRequested(s_l1Owner, to);
  }

  /// @notice set ownership, emit relevant events. Used in acceptOwnership()
  function _setL1Owner(address to) internal {
    address oldOwner = s_l1Owner;
    s_l1Owner = to;
    s_l1PendingOwner = address(0);

    emit L1OwnershipTransferred(oldOwner, to);
  }
}
