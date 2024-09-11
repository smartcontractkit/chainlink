// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CrossDomainOwnable} from "./CrossDomainOwnable.sol";
import {IForwarder} from "./interfaces/IForwarder.sol";

/**
 * @title CrossDomainForwarder - L1 xDomain account representation
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can consider that position to be held by the `l1Owner`
 */
abstract contract CrossDomainForwarder is IForwarder, CrossDomainOwnable {}
