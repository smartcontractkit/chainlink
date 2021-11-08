// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./CrossDomainForwarder.sol";
import "./interfaces/ForwarderInterface.sol";
import "./interfaces/DelegateForwarderInterface.sol";

/**
 * @title CrossDomainDelegateForwarder - L1 xDomain account representation (with delegatecall support)
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can consider that position to be held by the `l1Owner`
 */
abstract contract CrossDomainDelegateForwarder is DelegateForwarderInterface, CrossDomainOwnable {
  /**
   * @notice creates a new xDomain DelegateForwarder contract
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {}
}
