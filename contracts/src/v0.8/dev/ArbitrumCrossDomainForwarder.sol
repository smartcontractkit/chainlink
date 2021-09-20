// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../interfaces/TypeAndVersionInterface.sol";
import "./vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import "./CrossDomainForwarder.sol";

/**
 * @title ArbitrumCrossDomainForwarder - L1 xDomain account representation
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
contract ArbitrumCrossDomainForwarder is TypeAndVersionInterface, CrossDomainForwarder {
  /**
   * @notice creates a new Arbitrum xDomain Forwarder contract
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(
    address l1OwnerAddr
  )
    CrossDomainForwarder(l1OwnerAddr)
  {
    // noop
  }

  /**
   * @notice versions:
   *
   * - ArbitrumCrossDomainForwarder 0.1.0: initial release
   *
   * @inheritdoc TypeAndVersionInterface
   */
  function typeAndVersion()
    external
    pure
    override
    virtual
    returns (
      string memory
    )
  {
    return "ArbitrumCrossDomainForwarder 0.1.0";
  }

  /**
   * @notice The L2 xDomain `msg.sender`, generated from L1 sender address
   * @inheritdoc CrossDomainForwarder
   */
  function crossDomainMessenger()
    public
    view
    override
    virtual
    returns (address)
  {
    return AddressAliasHelper.applyL1ToL2Alias(l1Owner());
  }

  /**
   * @dev forwarded only if L2 Messenger calls with `xDomainMessageSender` beeing the L1 owner address
   * @inheritdoc ForwarderInterface
   */
  function forward(
    address target,
    bytes memory data
  )
    override
    external
  {
    // 1. The call MUST come from the L2 Messenger (deterministically generated from the L1 xDomain sender address)
    require(msg.sender == crossDomainMessenger(), "Sender is not the L2 messenger");
    // 2. Make the external call
    (bool success, bytes memory res) = target.call(data);
    require(success, string(abi.encode("xDomain call failed:", res)));
  }
}
