// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../dev/OperatorForwarder.sol";
import "../dev/ConfirmedOwner.sol";

contract OperatorForwarderDeployer is ConfirmedOwner {

  address private immutable linkAddress;
  address[] private s_authorisedSenders;

  event ForwarderDeployed(
    address indexed forwarder
  );

  constructor(
    address link,
    address[] memory authorizedSenders
  ) 
    ConfirmedOwner(msg.sender)
  {
    linkAddress = link;
    s_authorisedSenders = authorizedSenders;
  }

  function createForwarder()
    external
    returns (
      address
    )
  {
    OperatorForwarder newForwarder = new OperatorForwarder(linkAddress);
    address forwarderAddress = address(newForwarder);
    emit ForwarderDeployed(forwarderAddress);
    return forwarderAddress;
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
