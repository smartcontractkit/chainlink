// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./interfaces/DelegateForwarderInterface.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";
import "./ArbitrumCrossDomainForwarder.sol";

/**
 * @title ArbitrumCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Arbitrum
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract ArbitrumCrossDomainGovernor is DelegateForwarderInterface, ArbitrumCrossDomainForwarder {
  /**
   * @notice creates a new Arbitrum xDomain Forwarder contract
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(address l1OwnerAddr) ArbitrumCrossDomainForwarder(l1OwnerAddr) {}

  /**
   * @notice versions:
   *
   * - ArbitrumCrossDomainGovernor 1.0.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "ArbitrumCrossDomainGovernor 1.0.0";
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external override onlyCrossDomainOwners {
    Address.functionCall(target, data, "Governor call reverted");
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc DelegateForwarderInterface
   */
  function forwardDelegate(address target, bytes memory data) external override onlyCrossDomainOwners {
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }

  /**
   * @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
   */
  modifier onlyCrossDomainOwners() {
    require(msg.sender == crossDomainMessenger() || msg.sender == owner(), "Sender is not the L2 messenger or owner");
    _;
  }
}
