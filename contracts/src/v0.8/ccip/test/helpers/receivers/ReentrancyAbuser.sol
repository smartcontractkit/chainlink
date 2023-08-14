// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";
import {EVM2EVMOffRamp} from "../../../offRamp/EVM2EVMOffRamp.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";

contract ReentrancyAbuser is CCIPReceiver {
  event ReentrancySucceeded();

  bool internal s_ReentrancyDone = false;
  Internal.ExecutionReport internal s_payload;
  EVM2EVMOffRamp internal s_offRamp;

  constructor(address router, EVM2EVMOffRamp offRamp) CCIPReceiver(router) {
    s_offRamp = offRamp;
  }

  function setPayload(Internal.ExecutionReport calldata payload) public {
    s_payload = payload;
  }

  function _ccipReceive(Client.Any2EVMMessage memory) internal override {
    // Use original message gas limits in manual execution
    uint256 numMsgs = s_payload.messages.length;
    uint256[] memory gasOverrides = new uint256[](numMsgs);
    for (uint256 i = 0; i < numMsgs; ++i) {
      gasOverrides[i] = 0;
    }

    if (!s_ReentrancyDone) {
      // Could do more rounds but a PoC one is enough
      s_ReentrancyDone = true;
      s_offRamp.manuallyExecute(s_payload, gasOverrides);
    } else {
      emit ReentrancySucceeded();
    }
  }
}
