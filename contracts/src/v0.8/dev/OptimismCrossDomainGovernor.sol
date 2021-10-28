// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/TypeAndVersionInterface.sol";
import "./interfaces/DelegateForwarderInterface.sol";
import "./vendor/@eth-optimism/contracts/0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";
import "./CrossDomainDelegateForwarder.sol";

/**
 * @title OptimismCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Optimism
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract OptimismCrossDomainGovernor is TypeAndVersionInterface, CrossDomainDelegateForwarder {
  // OVM_L2CrossDomainMessenger is a precompile usually deployed to 0x4200000000000000000000000000000000000007
  address private immutable OVM_CROSS_DOMAIN_MESSENGER;

  /**
   * @notice creates a new Optimism xDomain Forwarder contract
   * @param crossDomainMessengerAddr the xDomain bridge messenger (Optimism bridge L2) contract address
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(address crossDomainMessengerAddr, address l1OwnerAddr) CrossDomainDelegateForwarder(l1OwnerAddr) {
    require(crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
    OVM_CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
  }

  /**
   * @notice versions:
   *
   * - OptimismCrossDomainForwarder 0.1.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismCrossDomainGovernor 0.1.0";
  }

  /**
   * @notice This is always the address of the OVM_L2CrossDomainMessenger contract
   * @inheritdoc CrossDomainForwarder
   */
  function crossDomainMessenger() public view virtual override returns (address) {
    return OVM_CROSS_DOMAIN_MESSENGER;
  }

  /**
   * @notice transfer ownership of this account to a new L1 owner
   * @dev Forwarding can be disabled by setting the L1 owner as `address(0)`. Accessible only by the L1 owner (not the L2 owner.)
   * @param to new L1 owner that will be allowed to call the forward fn
   */
  function transferL1Ownership(address to) external override {
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    require(
      iOVM_CrossDomainMessenger(OVM_CROSS_DOMAIN_MESSENGER).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _setL1Owner(to);
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external override {
    // 1. The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    require(
      msg.sender == OVM_CROSS_DOMAIN_MESSENGER || msg.sender == owner(),
      "Sender is not the L2 messenger or owner"
    );
    // 2. The L2 Messenger's caller MUST be the L1 Owner
    if (msg.sender == OVM_CROSS_DOMAIN_MESSENGER) {
      require(
        iOVM_CrossDomainMessenger(OVM_CROSS_DOMAIN_MESSENGER).xDomainMessageSender() == l1Owner(),
        "xDomain sender is not the L1 owner"
      );
    }
    // 3. Make the external call
    Address.functionCall(target, data, "Governor call reverted");
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
   * @inheritdoc DelegateForwarderInterface
   */
  function forwardDelegate(address target, bytes memory data) external override {
    // 1. The delegatecall MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    require(
      msg.sender == OVM_CROSS_DOMAIN_MESSENGER || msg.sender == owner(),
      "Sender is not the L2 messenger or owner"
    );
    // 2. The L2 Messenger's caller MUST be the L1 Owner
    if (msg.sender == OVM_CROSS_DOMAIN_MESSENGER) {
      require(
        iOVM_CrossDomainMessenger(OVM_CROSS_DOMAIN_MESSENGER).xDomainMessageSender() == l1Owner(),
        "xDomain sender is not the L1 owner"
      );
    }
    // 2. Make the external delegatecall
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }
}
