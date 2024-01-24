// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from ".././../interfaces/TypeAndVersionInterface.sol";
import {DelegateForwarderInterface} from "./interfaces/DelegateForwarderInterface.sol";
import {ForwarderInterface} from "./interfaces/ForwarderInterface.sol";

import {CrossDomainOwnable} from "./CrossDomainOwnable.sol";

import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/// @title CrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Scroll
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
/// can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
abstract contract CrossDomainGovernor is
  DelegateForwarderInterface,
  TypeAndVersionInterface,
  ForwarderInterface,
  CrossDomainOwnable
{
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {}

  /// @notice The address of the Cross Domain Messenger contract
  function crossDomainMessenger() external view virtual returns (address);

  /// @inheritdoc ForwarderInterface
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  function forward(address target, bytes memory data) external override {
    _onlyLocalOrCrossDomainOwner();
    Address.functionCall(target, data, "Governor call reverted");
  }

  /// @inheritdoc DelegateForwarderInterface
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  function forwardDelegate(address target, bytes memory data) external override {
    _onlyLocalOrCrossDomainOwner();
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }

  function _onlyLocalOrCrossDomainOwner() internal view virtual;
}
