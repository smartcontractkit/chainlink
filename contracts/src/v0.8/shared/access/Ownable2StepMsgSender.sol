// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {Ownable2Step} from "./Ownable2Step.sol";

/// @notice Sets the msg.sender to be the owner of the contract and does not set a pending owner.
contract Ownable2StepMsgSender is Ownable2Step {
  constructor() Ownable2Step(msg.sender, address(0)) {}
}
