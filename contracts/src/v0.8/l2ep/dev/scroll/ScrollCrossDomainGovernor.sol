// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {DelegateForwarderInterface} from "../interfaces/DelegateForwarderInterface.sol";
// solhint-disable-next-line no-unused-import
import {ForwarderInterface} from "../interfaces/ForwarderInterface.sol";

import {ScrollCrossDomainForwarder} from "./ScrollCrossDomainForwarder.sol";

import {IScrollMessenger} from "@scroll-tech/contracts/libraries/IScrollMessenger.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

///
/// @title ScrollCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Scroll
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
///   can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
///
contract ScrollCrossDomainGovernor is DelegateForwarderInterface, ScrollCrossDomainForwarder {
  ///
  /// @notice creates a new Scroll xDomain Forwarder contract
  /// @param crossDomainMessengerAddr the xDomain bridge messenger (Scroll bridge L2) contract address
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  /// @dev Empty constructor required due to inheriting from abstract contract CrossDomainForwarder
  ///
  constructor(
    IScrollMessenger crossDomainMessengerAddr,
    address l1OwnerAddr
  ) ScrollCrossDomainForwarder(crossDomainMessengerAddr, l1OwnerAddr) {}

  ///
  /// @notice versions:
  ///
  /// - ScrollCrossDomainGovernor 1.0.0: initial release
  ///
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ScrollCrossDomainGovernor 1.0.0";
  }

  ///
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  /// @inheritdoc ForwarderInterface
  ///
  function forward(address target, bytes memory data) external override onlyLocalOrCrossDomainOwner {
    Address.functionCall(target, data, "Governor call reverted");
  }

  ///
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  /// @inheritdoc DelegateForwarderInterface
  ///
  function forwardDelegate(address target, bytes memory data) external override onlyLocalOrCrossDomainOwner {
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }

  ///
  /// @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
  ///
  modifier onlyLocalOrCrossDomainOwner() {
    address messenger = crossDomainMessenger();
    // 1. The delegatecall MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    // solhint-disable-next-line custom-errors
    require(msg.sender == messenger || msg.sender == owner(), "Sender is not the L2 messenger or owner");
    // 2. The L2 Messenger's caller MUST be the L1 Owner
    if (msg.sender == messenger) {
      // solhint-disable-next-line custom-errors
      require(IScrollMessenger(messenger).xDomainMessageSender() == l1Owner(), "xDomain sender is not the L1 owner");
    }
    _;
  }
}
