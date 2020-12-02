pragma solidity 0.7.0;

import "./ConfirmedOwner.sol";
import "../interfaces/OperatorInterface.sol";

contract OperatorProxy is ConfirmedOwner {

  address internal immutable link;

  /**
   * @notice Deploy using the address of the LINK token
   * @dev The msg.sender is set as the owner of this contract
   * @param linkAddress Address of deployed LINK token
   */
  constructor(address linkAddress) ConfirmedOwner(msg.sender) {
    link = linkAddress;
  }

  /**
   * @notice Forward a call on to another address, checking that the
   * msg.sender is authorized to do so from the owner
   * @param to Target address
   * @param data Data to send to the target address
   */
  function forward(address to, bytes calldata data) public
  {
    require(OperatorInterface(owner()).isAuthorizedSender(msg.sender), "Not an authorized node");
    require(to != link, "Cannot send to Link token");
    (bool status,) = to.call(data);
    require(status, "Forwarded call failed.");
  }
}