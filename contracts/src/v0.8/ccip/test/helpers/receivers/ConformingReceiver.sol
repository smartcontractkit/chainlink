// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCIPReceiver} from "../../../applications/external/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";

contract ConformingReceiver is CCIPReceiver {
  event MessageReceived();

  constructor(address router, address feeToken) CCIPReceiver(router) {}

  function processMessage(
    Client.Any2EVMMessage calldata
  ) external virtual override {
    emit MessageReceived();
  }

  modifier isValidChain(uint64 chainSelector) virtual override {
    _;
  }

  modifier isValidSender(uint64 chainSelector, bytes memory sender) virtual override {
    _;
  }
}
