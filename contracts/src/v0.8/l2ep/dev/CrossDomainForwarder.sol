// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {ForwarderInterface} from "./interfaces/ForwarderInterface.sol";

import {CrossDomainOwnable} from "./CrossDomainOwnable.sol";

import {Address} from "../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/// @title CrossDomainForwarder - L1 xDomain account representation
/// @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
/// @dev Any other L2 contract which uses this contract's address as a privileged position,
///   can consider that position to be held by the `l1Owner`
abstract contract CrossDomainForwarder is TypeAndVersionInterface, ForwarderInterface, CrossDomainOwnable {
  /// @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
  constructor(address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {}

  /// @notice The address of the Cross Domain Messenger contract
  function crossDomainMessenger() external view virtual returns (address);

  /// @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
  /// @inheritdoc ForwarderInterface
  function forward(address target, bytes memory data) external override onlyL1Owner {
    Address.functionCall(target, data, "Forwarder call reverted");
  }
}
