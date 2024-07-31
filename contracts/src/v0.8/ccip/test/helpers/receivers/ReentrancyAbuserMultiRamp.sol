// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.19;

import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {EVM2EVMMultiOffRamp} from "../../../offRamp/EVM2EVMMultiOffRamp.sol";

contract ReentrancyAbuserMultiRamp is CCIPReceiver {
  event ReentrancySucceeded();

  bool internal s_ReentrancyDone = false;
  Internal.ExecutionReportSingleChain internal s_payload;
  EVM2EVMMultiOffRamp internal s_offRamp;

  constructor(address router, EVM2EVMMultiOffRamp offRamp) CCIPReceiver(router) {
    s_offRamp = offRamp;
  }

  function setPayload(Internal.ExecutionReportSingleChain calldata payload) public {
    s_payload = payload;
  }

  function _ccipReceive(Client.Any2EVMMessage memory) internal override {
    // Use original message gas limits in manual execution
    uint256 numMsgs = s_payload.messages.length;
    uint256[][] memory gasOverrides = new uint256[][](1);
    gasOverrides[0] = new uint256[](numMsgs);
    for (uint256 i = 0; i < numMsgs; ++i) {
      gasOverrides[0][i] = 0;
    }

    Internal.ExecutionReportSingleChain[] memory batchPayload = new Internal.ExecutionReportSingleChain[](1);
    batchPayload[0] = s_payload;

    if (!s_ReentrancyDone) {
      // Could do more rounds but a PoC one is enough
      s_ReentrancyDone = true;
      s_offRamp.manuallyExecute(batchPayload, gasOverrides);
    } else {
      emit ReentrancySucceeded();
    }
  }
}
