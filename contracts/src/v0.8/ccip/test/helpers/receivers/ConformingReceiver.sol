// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";

contract ConformingReceiver is CCIPReceiver {
  event MessageReceived();

  constructor(address router, address feeToken) CCIPReceiver(router) {}

  function _ccipReceive(
    Client.Any2EVMMessage memory
  ) internal virtual override {
    emit MessageReceived();
  }
}
