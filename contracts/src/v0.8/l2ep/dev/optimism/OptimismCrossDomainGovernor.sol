// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/DelegateForwarderInterface.sol";
import "../../../vendor/@eth-optimism/contracts/v0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";
import "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";
import "./OptimismCrossDomainForwarder.sol";

/**
 * @title OptimismCrossDomainGovernor - L1 xDomain account representation (with delegatecall support) for Optimism
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be simultaneously owned by the `l1Owner` and L2 `owner`
 */
contract OptimismCrossDomainGovernor is DelegateForwarderInterface, OptimismCrossDomainForwarder {
  /**
   * @notice creates a new Optimism xDomain Forwarder contract
   * @param crossDomainMessengerAddr the xDomain bridge messenger (Optimism bridge L2) contract address
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   * @dev Empty constructor required due to inheriting from abstract contract CrossDomainForwarder
   */
  constructor(
    iOVM_CrossDomainMessenger crossDomainMessengerAddr,
    address l1OwnerAddr
  ) OptimismCrossDomainForwarder(crossDomainMessengerAddr, l1OwnerAddr) {}

  /**
   * @notice versions:
   *
   * - OptimismCrossDomainForwarder 1.0.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismCrossDomainGovernor 1.0.0";
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external override onlyLocalOrCrossDomainOwner {
    Address.functionCall(target, data, "Governor call reverted");
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `msg.sender` being the L1 owner address, or called by the L2 owner
   * @inheritdoc DelegateForwarderInterface
   */
  function forwardDelegate(address target, bytes memory data) external override onlyLocalOrCrossDomainOwner {
    Address.functionDelegateCall(target, data, "Governor delegatecall reverted");
  }

  /**
   * @notice The call MUST come from either the L1 owner (via cross-chain message) or the L2 owner. Reverts otherwise.
   */
  modifier onlyLocalOrCrossDomainOwner() {
    address messenger = crossDomainMessenger();
    // 1. The delegatecall MUST come from either the L1 owner (via cross-chain message) or the L2 owner
    require(msg.sender == messenger || msg.sender == owner(), "Sender is not the L2 messenger or owner");
    // 2. The L2 Messenger's caller MUST be the L1 Owner
    if (msg.sender == messenger) {
      require(
        iOVM_CrossDomainMessenger(messenger).xDomainMessageSender() == l1Owner(),
        "xDomain sender is not the L1 owner"
      );
    }
    _;
  }
}
