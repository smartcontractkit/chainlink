pragma solidity 0.7.0;

import "./Owned.sol";
import "../interfaces/OperatorInterface.sol";

contract OperatorProxy is Owned {

  address internal immutable link;

  constructor(address linkAddress) Owned(msg.sender) {
    link = linkAddress;
  }

  function forward(address to, bytes calldata data) public
  {
    require(OperatorInterface(owner).isAuthorizedSender(msg.sender), "Not an authorized node");
    require(to != link, "Cannot send to Link token");
    (bool status,) = to.call(data);
    require(status, "Forwarded call failed.");
  }
}