// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IDelegateForwarderInterface} from "./interfaces/IDelegateForwarderInterface.sol";
// solhint-disable-next-line no-unused-import
import {IForwarderInterface} from "./interfaces/IForwarderInterface.sol";

import {CrossDomainForwarder} from "./CrossDomainForwarder.sol";

import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/// @title CrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Scroll
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
/// can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
abstract contract CrossDomainGovernor is IDelegateForwarderInterface, CrossDomainForwarder {
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(address l1OwnerAddr) CrossDomainForwarder(l1OwnerAddr) {}

  /// @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
  function _requireLocalOrCrossDomainOwner() internal view virtual;

  /// @inheritdoc IForwarderInterface
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  function forward(address target, bytes memory data) external override {
    _requireLocalOrCrossDomainOwner();
    Address.functionCall(target, data, "Governor call reverted");
  }

  /// @inheritdoc IDelegateForwarderInterface
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  function forwardDelegate(address target, bytes memory data) external override {
    _requireLocalOrCrossDomainOwner();
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }
}
