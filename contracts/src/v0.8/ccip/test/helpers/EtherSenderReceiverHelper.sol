// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {EtherSenderReceiver} from "../../applications/EtherSenderReceiver.sol";
import {Client} from "../../libraries/Client.sol";

contract EtherSenderReceiverHelper is EtherSenderReceiver {
  constructor(
    address router
  ) EtherSenderReceiver(router) {}

  function validatedMessage(
    Client.EVM2AnyMessage calldata message
  ) public view returns (Client.EVM2AnyMessage memory) {
    return _validatedMessage(message);
  }

  function validateFeeToken(
    Client.EVM2AnyMessage calldata message
  ) public payable {
    _validateFeeToken(message);
  }

  function publicCcipReceive(
    Client.Any2EVMMessage memory message
  ) public {
    _ccipReceive(message);
  }
}
