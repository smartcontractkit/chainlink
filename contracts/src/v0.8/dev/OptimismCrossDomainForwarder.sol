// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/TypeAndVersionInterface.sol";

/* ./dev dependencies - to be moved from ./dev after audit */
import "./CrossDomainForwarder.sol";
import "./vendor/@eth-optimism/contracts/0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";
import "./vendor/openzeppelin-solidity/v4.3.1/contracts/utils/Address.sol";

/**
 * @title OptimismCrossDomainForwarder - L1 xDomain account representation
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract OptimismCrossDomainForwarder is TypeAndVersionInterface, CrossDomainForwarder {
  // OVM_L2CrossDomainMessenger is a precompile usually deployed to 0x4200000000000000000000000000000000000007
  address private immutable OVM_CROSS_DOMAIN_MESSENGER;

  /**
   * @notice creates a new Optimism xDomain Forwarder contract
   * @param crossDomainMessengerAddr the xDomain bridge messenger (Optimism bridge L2) contract address
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(address crossDomainMessengerAddr, address l1OwnerAddr) CrossDomainForwarder(l1OwnerAddr) {
    require(crossDomainMessengerAddr != address(0), "Invalid xDomain Messenger address");
    OVM_CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
  }

  /**
   * @notice versions:
   *
   * - OptimismCrossDomainForwarder 0.1.0: initial release
   * - OptimismCrossDomainForwarder 0.2.0: Use OZ Address, CrossDomainOwnable
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismCrossDomainForwarder 0.2.0";
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external virtual override onlyCrossDomainMessenger {
    Address.functionCall(target, data, "Forwarder call reverted");
  }

  /**
   * @notice This is always the address of the OVM_L2CrossDomainMessenger contract
   */
  function crossDomainMessenger() public view returns (address) {
    return OVM_CROSS_DOMAIN_MESSENGER;
  }

  /**
   * @notice transfer ownership of this account to a new L1 owner
   * @param to new L1 owner that will be allowed to call the forward fn
   */
  function transferL1Ownership(address to) external override onlyCrossDomainMessenger {
    super._transferL1Ownership(to);
  }

  /**
   * @notice accept ownership of this account to a new L1 owner
   */
  function acceptL1Ownership() external virtual override {
    require(
      iOVM_CrossDomainMessenger(OVM_CROSS_DOMAIN_MESSENGER).xDomainMessageSender() == s_l1PendingOwner,
      "Must be proposed owner"
    );
    super._setL1Owner(s_l1PendingOwner);
  }

  /**
   * @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
   */
  modifier onlyCrossDomainMessenger() {
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    require(
      iOVM_CrossDomainMessenger(crossDomainMessenger()).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }
}
