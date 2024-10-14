// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.19;

import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";

contract ReentrancyAbuserMultiRamp is CCIPReceiver {
  event ReentrancySucceeded();

  bool internal s_ReentrancyDone = false;
  Internal.ExecutionReportSingleChain internal s_payload;
  OffRamp internal s_offRamp;

  constructor(address router, OffRamp offRamp) CCIPReceiver(router) {
    s_offRamp = offRamp;
  }

  function setPayload(Internal.ExecutionReportSingleChain calldata payload) public {
    s_payload = payload;
  }

  function _ccipReceive(Client.Any2EVMMessage memory) internal override {
    // Use original message gas limits in manual execution
    uint256 numMsgs = s_payload.messages.length;
    OffRamp.GasLimitOverride[][] memory gasOverrides = new OffRamp.GasLimitOverride[][](1);
    gasOverrides[0] = new OffRamp.GasLimitOverride[](numMsgs);
    for (uint256 i = 0; i < numMsgs; ++i) {
      gasOverrides[0][i].receiverExecutionGasLimit = 0;
      gasOverrides[0][i].tokenGasOverrides = new uint32[](s_payload.messages[i].tokenAmounts.length);
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
