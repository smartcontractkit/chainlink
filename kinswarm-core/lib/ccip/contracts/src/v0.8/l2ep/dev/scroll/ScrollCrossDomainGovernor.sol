// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ITypeAndVersion} from "../../../shared/interfaces/ITypeAndVersion.sol";
import {IDelegateForwarder} from "../interfaces/IDelegateForwarder.sol";
// solhint-disable-next-line no-unused-import
import {IForwarder} from "../interfaces/IForwarder.sol";

import {CrossDomainForwarder} from "../CrossDomainForwarder.sol";
import {CrossDomainOwnable} from "../CrossDomainOwnable.sol";

import {IScrollMessenger} from "@scroll-tech/contracts/libraries/IScrollMessenger.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/// @title ScrollCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Scroll
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
/// can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
contract ScrollCrossDomainGovernor is IDelegateForwarder, ITypeAndVersion, CrossDomainForwarder {
  string public constant override typeAndVersion = "ScrollCrossDomainGovernor 1.0.0";

  address internal immutable i_scrollCrossDomainMessenger;

  /// @param crossDomainMessengerAddr the xDomain bridge messenger (Scroll bridge L2) contract address
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(IScrollMessenger crossDomainMessengerAddr, address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {
    // solhint-disable-next-line gas-custom-errors
    require(address(crossDomainMessengerAddr) != address(0), "Invalid xDomain Messenger address");
    i_scrollCrossDomainMessenger = address(crossDomainMessengerAddr);
  }

  /// @inheritdoc IForwarder
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  function forward(address target, bytes memory data) external override onlyLocalOrCrossDomainOwner {
    Address.functionCall(target, data, "Governor call reverted");
  }

  /// @inheritdoc IDelegateForwarder
  /// @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
  function forwardDelegate(address target, bytes memory data) external override onlyLocalOrCrossDomainOwner {
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }

  /// @notice The address of the Scroll Cross Domain Messenger contract
  function crossDomainMessenger() external view returns (address) {
    return address(i_scrollCrossDomainMessenger);
  }

  /// @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyL1Owner() override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == i_scrollCrossDomainMessenger, "Sender is not the L2 messenger");
    // solhint-disable-next-line gas-custom-errors
    require(
      IScrollMessenger(i_scrollCrossDomainMessenger).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }

  /// @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
  modifier onlyLocalOrCrossDomainOwner() {
    // 1. The delegatecall MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    // solhint-disable-next-line gas-custom-errors
    require(
      msg.sender == i_scrollCrossDomainMessenger || msg.sender == owner(),
      "Sender is not the L2 messenger or owner"
    );
    // 2. The L2 Messenger's caller MUST be the L1 Owner
    if (msg.sender == i_scrollCrossDomainMessenger) {
      // solhint-disable-next-line gas-custom-errors
      require(
        IScrollMessenger(i_scrollCrossDomainMessenger).xDomainMessageSender() == l1Owner(),
        "xDomain sender is not the L1 owner"
      );
    }
    _;
  }

  /// @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyProposedL1Owner() override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == i_scrollCrossDomainMessenger, "Sender is not the L2 messenger");
    // solhint-disable-next-line gas-custom-errors
    require(
      IScrollMessenger(i_scrollCrossDomainMessenger).xDomainMessageSender() == s_l1PendingOwner,
      "Must be proposed L1 owner"
    );
    _;
  }
}
