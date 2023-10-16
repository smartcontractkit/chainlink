// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
// solhint-disable-next-line no-unused-import
import {ForwarderInterface} from "../interfaces/ForwarderInterface.sol";

/* ./dev dependencies - to be moved from ./dev after audit */
import {CrossDomainForwarder} from "../CrossDomainForwarder.sol";
import {CrossDomainOwnable} from "../CrossDomainOwnable.sol";

import {iOVM_CrossDomainMessenger} from "../../../vendor/@eth-optimism/contracts/v0.4.7/contracts/optimistic-ethereum/iOVM/bridge/messaging/iOVM_CrossDomainMessenger.sol";
import {Address} from "../../../vendor/openzeppelin-solidity/v4.7.3/contracts/utils/Address.sol";

/**
 * @title OptimismCrossDomainForwarder - L1 xDomain account representation
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract OptimismCrossDomainForwarder is TypeAndVersionInterface, CrossDomainForwarder {
  // OVM_L2CrossDomainMessenger is a precompile usually deployed to 0x4200000000000000000000000000000000000007
  // solhint-disable-next-line chainlink-solidity/prefix-immutable-variables-with-i
  iOVM_CrossDomainMessenger private immutable OVM_CROSS_DOMAIN_MESSENGER;

  /**
   * @notice creates a new Optimism xDomain Forwarder contract
   * @param crossDomainMessengerAddr the xDomain bridge messenger (Optimism bridge L2) contract address
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(iOVM_CrossDomainMessenger crossDomainMessengerAddr, address l1OwnerAddr) CrossDomainOwnable(l1OwnerAddr) {
    // solhint-disable-next-line custom-errors
    require(address(crossDomainMessengerAddr) != address(0), "Invalid xDomain Messenger address");
    OVM_CROSS_DOMAIN_MESSENGER = crossDomainMessengerAddr;
  }

  /**
   * @notice versions:
   *
   * - OptimismCrossDomainForwarder 0.1.0: initial release
   * - OptimismCrossDomainForwarder 1.0.0: Use OZ Address, CrossDomainOwnable
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion() external pure virtual override returns (string memory) {
    return "OptimismCrossDomainForwarder 1.0.0";
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` being the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(address target, bytes memory data) external virtual override onlyL1Owner {
    Address.functionCall(target, data, "Forwarder call reverted");
  }

  /**
   * @notice This is always the address of the OVM_L2CrossDomainMessenger contract
   */
  function crossDomainMessenger() public view returns (address) {
    return address(OVM_CROSS_DOMAIN_MESSENGER);
  }

  /**
   * @notice The call MUST come from the L1 owner (via cross-chain message.) Reverts otherwise.
   */
  modifier onlyL1Owner() override {
    // solhint-disable-next-line custom-errors
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    // solhint-disable-next-line custom-errors
    require(
      iOVM_CrossDomainMessenger(crossDomainMessenger()).xDomainMessageSender() == l1Owner(),
      "xDomain sender is not the L1 owner"
    );
    _;
  }

  /**
   * @notice The call MUST come from the proposed L1 owner (via cross-chain message.) Reverts otherwise.
   */
  modifier onlyProposedL1Owner() override {
    address messenger = crossDomainMessenger();
    // solhint-disable-next-line custom-errors
    require(msg.sender == messenger, "Sender is not the L2 messenger");
    // solhint-disable-next-line custom-errors
    require(
      iOVM_CrossDomainMessenger(messenger).xDomainMessageSender() == s_l1PendingOwner,
      "Must be proposed L1 owner"
    );
    _;
  }
}
