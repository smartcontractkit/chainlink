// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/TypeAndVersionInterface.sol";
import "./interfaces/DelegateForwarderInterface.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import "./ArbitrumCrossDomainForwarder.sol";
import "./CrossDomainDelegateForwarder.sol";

/**
 * @title ArbitrumCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Arbitrum
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract ArbitrumCrossDomainGovernor is TypeAndVersionInterface, CrossDomainDelegateForwarder {
  /**
   * @notice creates a new Arbitrum xDomain Forwarder contract
   * @param l1OwnerAddr tetherhe L1 owner address that will be allowed to call the forward fn
   */
  constructor(address l1OwnerAddr) CrossDomainDelegateForwarder(l1OwnerAddr) {
    // noop
  }

  /**
   * @notice versions:
   *
   * - ArbitrumCrossDomainGovernor 0.1.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ArbitrumCrossDomainGovernor 0.1.0";
  }

  /**
   * @notice The L2 xDomain `msg.sender`, generated from L1 sender address
   * @inheritdoc CrossDomainForwarder
   */
  function crossDomainMessenger() public view virtual override returns (address) {
    return AddressAliasHelper.applyL1ToL2Alias(l1Owner());
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external override {
    // 1. The call MUST come from the L2 Messenger (deterministically generated from the L1 xDomain sender address)
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    // 2. Make the external call
    (bool success, bytes memory res) = target.call(data);
    require(success, string(abi.encode("xDomain call failed:", res)));
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc DelegateForwarderInterface
   */
  function forwardDelegate(address target, bytes memory data) external override {
    // 1. The delegatecall MUST come from the L2 Messenger (deterministically generated from the L1 xDomain sender address)
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    // 2. Make the external delegatecall
    (bool success, bytes memory res) = target.delegatecall(data);
    require(success, string(abi.encode("xDomain delegatecall failed:", res)));
  }
}
