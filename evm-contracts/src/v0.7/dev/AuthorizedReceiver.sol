// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../interfaces/AuthorizedReceiverInterface.sol";

abstract contract AuthorizedReceiver is
  AuthorizedReceiverInterface
{

  mapping(address => bool) private s_authorizedSenders;
  address[] private s_authorizedSenderList;

  event AuthorizedSendersChanged(
    address[] senders,
    address changedBy
  );

  /**
   * @notice Sets the fulfillment permission for a given node. Use `true` to allow, `false` to disallow.
   * @param senders The addresses of the authorized Chainlink node
   */
  function setAuthorizedSenders(
    address[] calldata senders
  )
    external
    override
    validateAuthorizedSenderSetter()
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
    emit AuthorizedSendersChanged(senders, msg.sender);
  }

  /**
   * @notice Retrieve a list of authorized senders
   * @return array of addresses
   */
  function getAuthorizedSenders()
    external
    view
    override
    returns (
      address[] memory
    )
  {
    return s_authorizedSenderList;
  }

  /**
   * @notice Use this to check if a node is authorized for fulfilling requests
   * @param sender The address of the Chainlink node
   * @return The authorization status of the node
   */
  function isAuthorizedSender(
    address sender
  )
    public
    view
    override
    returns (bool)
  {
    return s_authorizedSenders[sender];
  }

  /**
   * @notice customizable guard of who can update the authorized sender list
   * @return bool whether sender can update authorized sender list
   */
  function _canSetAuthorizedSenders()
    internal
    virtual
    returns (bool);

  /**
   * @notice validates the sender is an authorized sender
   */
  function _validateIsAuthorizedSender()
    internal
    view
  {
    require(isAuthorizedSender(msg.sender), "Not authorized sender");
  }

  /**
   * @notice prevents non-authorized addresses from calling this method
   */
  modifier validateAuthorizedSender()
  {
    _validateIsAuthorizedSender();
    _;
  }

  /**
   * @notice prevents non-authorized addresses from calling this method
   */
  modifier validateAuthorizedSenderSetter()
  {
    require(_canSetAuthorizedSenders(), "Cannot set authorized senders");
    _;
  }

}
