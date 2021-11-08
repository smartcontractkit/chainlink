// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/TypeAndVersionInterface.sol";
import "./interfaces/DelegateForwarderInterface.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";
import "./CrossDomainDelegateForwarder.sol";

/**
 * @title ArbitrumCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Arbitrum
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract ArbitrumCrossDomainGovernor is CrossDomainDelegateForwarder, TypeAndVersionInterface {
  /**
   * @notice creates a new Arbitrum xDomain Forwarder contract
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
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
   * @notice transfer ownership of this account to a new L1 owner
   * @dev Forwarding can be disabled by setting the L1 owner as `address(0)`. Accessible only by the L1 owner (not the L2 owner.)
   * @param to new L1 owner that will be allowed to call the forward fn
   */
  function transferL1Ownership(address to) external override {
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    _setL1Owner(to);
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external override {
    // 1. The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    require(msg.sender == crossDomainMessenger() || msg.sender == owner(), "Sender is not the L2 messenger or owner");
    // 2. Make the external call
    Address.functionCall(target, data, "Governor call reverted");
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc DelegateForwarderInterface
   */
  function forwardDelegate(address target, bytes memory data) external override {
    // 1. The delegatecall MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    require(msg.sender == crossDomainMessenger() || msg.sender == owner(), "Sender is not the L2 messenger or owner");
    // 2. Make the external delegatecall
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }
}
