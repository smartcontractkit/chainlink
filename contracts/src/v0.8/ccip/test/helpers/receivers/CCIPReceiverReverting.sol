// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {CCIPReceiver} from "../../../applications/external/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";

contract CCIPReceiverReverting is CCIPReceiver {
  error ErrorCase();

  bool private s_simRevert;

  constructor(
    address router
  ) CCIPReceiver(router) {}

  /// @notice This function the entrypoint for this contract to process messages.
  /// @param message The message to process.
  /// @dev This example just sends the tokens to the owner of this contracts. More
  /// interesting functions could be implemented.
  /// @dev It has to be external because of the try/catch.
  function processMessage(
    Client.Any2EVMMessage calldata message
  ) external view override onlySelf isValidSender(message.sourceChainSelector, message.sender) {
    // Meant to help simulate a failed-message
    if (s_simRevert) revert ErrorCase();
  }

  // An example function to demonstrate recovery
  function setSimRevert(
    bool simRevert
  ) external {
    s_simRevert = simRevert;
  }
}
