// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../interfaces/OperatorInterface.sol";

contract OperatorForwarder {
  address public immutable authorizedSender1;
  address public immutable authorizedSender2;
  address public immutable authorizedSender3;

  constructor() {
    address[] memory authorizedSenders = OperatorInterface(msg.sender).getAuthorizedSenders();
    authorizedSender1 = authorizedSenders[0];
    authorizedSender2 = (authorizedSenders.length > 1) ? authorizedSenders[1] : address(0);
    authorizedSender3 = (authorizedSenders.length > 2) ? authorizedSenders[2] : address(0);
  }
}
