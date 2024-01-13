// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CrossDomainGovernor} from "../CrossDomainGovernor.sol";

import {AddressAliasHelper} from "../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";

/// @title ArbitrumCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Arbitrum
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
/// can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
contract ArbitrumCrossDomainGovernor is CrossDomainGovernor {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "ArbitrumCrossDomainGovernor 1.0.0";

  /// @notice creates a new xDomain governor contract
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  /// @dev Empty constructor required due to inheriting from abstract contract CrossDomainGovernor
  constructor(address l1OwnerAddr) CrossDomainGovernor(l1OwnerAddr) {}

  /// @notice The L2 xDomain `msg.sender`, generated from L1 sender address
  function crossDomainMessenger() external view override returns (address) {
    return AddressAliasHelper.applyL1ToL2Alias(l1Owner());
  }

  /// @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
  modifier onlyLocalOrCrossDomainOwner() override {
    // solhint-disable-next-line custom-errors
    require(
      msg.sender == this.crossDomainMessenger() || msg.sender == owner(),
      "Sender is not the L2 messenger or owner"
    );
    _;
  }

  /// @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyL1Owner() override {
    // solhint-disable-next-line custom-errors
    require(msg.sender == this.crossDomainMessenger(), "Sender is not the L2 messenger");
    _;
  }

  /// @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
  modifier onlyProposedL1Owner() override {
    // solhint-disable-next-line custom-errors
    require(msg.sender == AddressAliasHelper.applyL1ToL2Alias(s_l1PendingOwner), "Must be proposed L1 owner");
    _;
  }
}
