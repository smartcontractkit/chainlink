// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../dev/OperatorForwarder.sol";
import "../dev/ConfirmedOwner.sol";

contract OperatorForwarderDeployer is ConfirmedOwner {

  address private immutable linkAddress;
  bytes32 private immutable salt;
  address[] private s_authorisedSenders;
  OperatorForwarder public forwarder;

  event ForwarderDeployed(address indexed forwarder);

  constructor(
    address link,
    address[] memory authorizedSenders
  ) 
    ConfirmedOwner(msg.sender)
  {
    linkAddress = link;
    setAuthorizedSenders(authorizedSenders);
    salt = bytes32("1");
  }

  function createForwarder()
    external
    returns (
      address
    )
  {
    forwarder = new OperatorForwarder{salt: salt}(linkAddress);
    address forwarderAddress = address(forwarder);
    emit ForwarderDeployed(forwarderAddress);
    return forwarderAddress;
  }

  function destroyForwarder()
    external
  {
    forwarder.destroy();
  }

  function setAuthorizedSenders(
    address[] memory authorizedSenders
  )
    public
  {
    s_authorisedSenders = authorizedSenders;
  }

  function getAuthorizedSenders()
    external
    view
    returns (
      address[] memory
    )
  {
    return s_authorisedSenders;
  }
}
