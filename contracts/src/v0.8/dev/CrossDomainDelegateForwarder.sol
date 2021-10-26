// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./CrossDomainForwarder.sol";
import "./interfaces/ForwarderInterface.sol";
import "./interfaces/DelegateForwarderInterface.sol";

/**
 * @title CrossDomainDelegateForwarder - L1 xDomain account representation (with delegatecall support)
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
abstract contract CrossDomainDelegateForwarder is DelegateForwarderInterface, CrossDomainForwarder {
  /**
   * @notice creates a new xDomain Forwarder contract
   * @dev Forwarding can be disabled by setting the L1 owner as `address(0)`.
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(address l1OwnerAddr) CrossDomainForwarder(l1OwnerAddr) {
    // noop
  }
}
