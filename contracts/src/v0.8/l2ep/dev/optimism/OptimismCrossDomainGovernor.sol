// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CrossDomainGovernor} from "../CrossDomainGovernor.sol";

import {iOVM_CrossDomainMessenger} from "../../../vendor/@eth-optimism/contracts/v0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";

/// @title OptimismCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Optimism
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
/// can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
contract OptimismCrossDomainGovernor is CrossDomainGovernor {
  string public constant override typeAndVersion = "OptimismCrossDomainGovernor 1.0.0";

  /// The cross domain messenger address
  address internal immutable i_crossDomainMessengerAddr;

  /// @notice creates a new xDomain governor contract
  /// @param crossDomainMessengerAddr the xDomain bridge messenger (Optimism bridge L2) contract address
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  /// @dev Empty constructor required due to inheriting from abstract contract CrossDomainGovernor
  constructor(
    iOVM_CrossDomainMessenger crossDomainMessengerAddr,
    address l1OwnerAddr
  ) CrossDomainGovernor(l1OwnerAddr) {
    i_crossDomainMessengerAddr = address(crossDomainMessengerAddr);

    // solhint-disable-next-line gas-custom-errors
    require(i_crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
  }

  /// @notice The address of the Cross Domain Messenger contract
  function crossDomainMessenger() external view override returns (address) {
    return i_crossDomainMessengerAddr;
  }

  /// @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
  function _requireLocalOrCrossDomainOwner() internal view override {
    // 1. The delegatecall MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    // solhint-disable-next-line gas-custom-errors
    require(
      msg.sender == i_crossDomainMessengerAddr || msg.sender == owner(),
      "Sender is not the L2 messenger or owner"
    );
    // 2. The L2 Messenger's caller MUST be the L1 Owner
    if (msg.sender == i_crossDomainMessengerAddr) {
      // solhint-disable-next-line gas-custom-errors
      require(
        iOVM_CrossDomainMessenger(i_crossDomainMessengerAddr).xDomainMessageSender() == l1Owner(),
        "xDomain sender is not the L1 owner"
      );
    }
  }

  /// @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyL1Owner() override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == i_crossDomainMessengerAddr, "Sender is not the L2 messenger");
    // solhint-disable-next-line gas-custom-errors
    require(
      iOVM_CrossDomainMessenger(i_crossDomainMessengerAddr).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }

  /// @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyProposedL1Owner() override {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == i_crossDomainMessengerAddr, "Sender is not the L2 messenger");
    // solhint-disable-next-line gas-custom-errors
    require(
      iOVM_CrossDomainMessenger(i_crossDomainMessengerAddr).xDomainMessageSender() == s_l1PendingOwner,
      "Must be proposed L1 owner"
    );
    _;
  }
}
