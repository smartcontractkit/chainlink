// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";

contract CCIPDummyReceiver is CCIPReceiver {
  event MessageReceived(bytes32 messageId, uint64 sourceChainSelector, bytes data);

  constructor(
    address router
  ) CCIPReceiver(router) {}

  function _ccipReceive(
    Client.Any2EVMMessage memory message
  ) internal virtual override {
    emit MessageReceived(message.messageId, message.sourceChainSelector, message.data);
  }
}
