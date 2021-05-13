// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../interfaces/OperatorInterface.sol";
import "./ConfirmedOwner.sol";

contract OperatorForwarder is ConfirmedOwner {

  mapping(address => bool) private s_authorizedSenders;
  address[] private s_authorizedSenderList;

  address public immutable linkAddr;

  event AuthorizedSendersChanged(
    address[] senders
  );

  constructor(
    address link
  ) 
    ConfirmedOwner(msg.sender)
  {
    linkAddr = link;
  }

  /**
   * @notice Sets the fulfillment permission for a given node. Use `true` to allow, `false` to disallow.
   * @param senders The addresses of the authorized Chainlink node
   */
  function setAuthorizedSenders(
    address[] calldata senders
  )
    external
    onlyOwner()
  {
    require(senders.length > 0, "Must have at least 1 authorized sender");
    // Set previous authorized senders to false
    uint256 authorizedSendersLength = s_authorizedSenderList.length;
    for (uint256 i = 0; i < authorizedSendersLength; i++) {
      s_authorizedSenders[s_authorizedSenderList[i]] = false;
    }
    // Set new to true
    for (uint256 i = 0; i < senders.length; i++) {
      s_authorizedSenders[senders[i]] = true;
    }
    // Replace list
    s_authorizedSenderList = senders;
    emit AuthorizedSendersChanged(senders);
  }

  /**
   * @notice Retrieve a list of authorized senders
   * @return array of addresses
   */
  function getAuthorizedSenders()
    external
    view
    returns (
      address[] memory
    )
  {
    return s_authorizedSenderList;
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
    onlyAuthorizedSender()
  {
    require(to != linkAddr, "Cannot #forward to Link token");
    (bool status,) = to.call(data);
    require(status, "Forwarded call failed.");
  }

  modifier onlyAuthorizedSender() {
    require(s_authorizedSenders[msg.sender], "Not an authorized node to fulfill requests");
    _;
  }
}
