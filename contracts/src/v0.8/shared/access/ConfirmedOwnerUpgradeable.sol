// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {ConfirmedOwnerWithProposalUpgradeable} from "./ConfirmedOwnerWithProposalUpgradeable.sol";

/// @title The ConfirmedOwner contract
/// @notice A contract with helpers for basic contract ownership.
contract ConfirmedOwnerUpgradeable is ConfirmedOwnerWithProposalUpgradeable {
  // solhint-disable-next-line func-name-mixedcase
  function __ConfirmedOwner_init(address newOwner) internal onlyInitializing {
    __ConfirmedOwnerWithProposal_init(newOwner, address(0));
  }
}
