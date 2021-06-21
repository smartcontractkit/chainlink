// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../interfaces/OperatorInterface.sol";
import "./ConfirmedOwnerWithProposal.sol";
import "./AuthorizedReceiver.sol";

contract AuthorizedForwarder is
  ConfirmedOwnerWithProposal,
  AuthorizedReceiver
{

  address public immutable getChainlinkToken;

  event OwnershipTransferRequestedWithMessage(
    address indexed from,
    address indexed to,
    bytes message
  );

  constructor(
    address link,
    address owner,
    address recipient,
    bytes memory message
  )
    ConfirmedOwnerWithProposal(owner, recipient)
  {
    getChainlinkToken = link;
    if (recipient != address(0)) {
      emit OwnershipTransferRequestedWithMessage(owner, recipient, message);
    }
  }

  /**
   * @notice Forward a call to another contract
   * @dev Only callable by an authorized sender
   * @param to address
   * @param data to forward
   */
  function forward(
    address to,
    bytes calldata data
  )
    external
    validateAuthorizedSender()
  {
    require(to != getChainlinkToken, "Cannot #forward to Link token");
    (bool status,) = to.call(data);
    require(status, "Forwarded call failed");
  }

  /**
   * @notice Forward a call to another contract
   * @dev Only callable by the owner
   * @param to address
   * @param data to forward
   */
  function ownerForward(
    address to,
    bytes calldata data
  )
    external
    onlyOwner()
  {
    (bool status,) = to.call(data);
    require(status, "Forwarded call failed");
  }

  /**
   * @notice Transfer ownership with instructions for recipient
   * @param to address proposed recipeint of ownership
   * @param message instructions for recipient upon accepting ownership
   */
  function transferOwnershipWithMessage(
    address to,
    bytes memory message
  )
    public
  {
    transferOwnership(to);
    emit OwnershipTransferRequestedWithMessage(msg.sender, to, message);
  }

  /**
   * @notice concrete implementation of AuthorizedReceiver
   * @return bool of whether sender is authorized
   */
  function _canSetAuthorizedSenders()
    internal
    view
    override
    returns (bool)
  {
    return owner() == msg.sender;
  }

}
