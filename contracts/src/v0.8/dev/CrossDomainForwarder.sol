// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../ConfirmedOwner.sol";
import "./interfaces/ForwarderInterface.sol";

/**
 * @title CrossDomainForwarder - L1 xDomain account representation
 * @notice L2 Contract which receives messages from a specific L1 address and transparently forwards them to the destination.
 * @dev Any other L2 contract which uses this contract's address as a privileged position,
 *   can be considered to be owned by the `l1Owner`
 */
abstract contract CrossDomainForwarder is ForwarderInterface, ConfirmedOwner {
  address private s_l1Owner;

  event L1OwnershipTransferred(
    address indexed from,
    address indexed to
  );

  /**
   * @notice creates a new xDomain Forwarder contract
   * @dev Forwarding can be disabled by setting the L1 owner as `address(0)`.
   * @param l1OwnerAddr the L1 owner address that will be allowed to call the forward fn
   */
  constructor(
    address l1OwnerAddr
  )
    ConfirmedOwner(msg.sender)
  {
    _setL1Owner(l1OwnerAddr);
  }

  /// @return xDomain messenger address (L2 `msg.sender`)
  function crossDomainMessenger()
    public
    view
    virtual
    returns (address);

  /// @return L1 owner address
  function l1Owner()
    view
    public
    virtual
    returns (address)
  {
    return s_l1Owner;
  }

  /**
   * @notice transfer ownership of this account to a new L1 owner
   * @dev Forwarding can be disabled by setting the L1 owner as `address(0)`. Accessible only by owner.
   * @param to new L1 owner that will be allowed to call the forward fn
   */
  function transferL1Ownership(
    address to
  )
    external
    virtual
    onlyOwner
  {
    _setL1Owner(to);
  }

  /// @notice internal method that stores the L1 owner
  function _setL1Owner(
    address to
  )
    internal
  {
    address from = s_l1Owner;
    if (from != to) {
      s_l1Owner = to;
      emit L1OwnershipTransferred(from, to);
    }
  }
}
