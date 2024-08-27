// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {CCIPReceiver} from "../../../applications/CCIPReceiver.sol";
import {Client} from "../../../libraries/Client.sol";
import {Internal} from "../../../libraries/Internal.sol";
import {EVM2EVMOffRamp} from "../../../offRamp/EVM2EVMOffRamp.sol";

contract ReentrancyAbuser is CCIPReceiver {
  event ReentrancySucceeded();

  uint32 internal constant DEFAULT_TOKEN_DEST_GAS_OVERHEAD = 144_000;

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
    EVM2EVMOffRamp.GasLimitOverride[] memory gasOverrides = _getGasLimitsFromMessages(s_payload.messages);

    if (!s_ReentrancyDone) {
      // Could do more rounds but a PoC one is enough
      s_ReentrancyDone = true;
      s_offRamp.manuallyExecute(s_payload, gasOverrides);
    } else {
      emit ReentrancySucceeded();
    }
  }

  function _getGasLimitsFromMessages(
    Internal.EVM2EVMMessage[] memory messages
  ) internal pure returns (EVM2EVMOffRamp.GasLimitOverride[] memory) {
    EVM2EVMOffRamp.GasLimitOverride[] memory gasLimitOverrides = new EVM2EVMOffRamp.GasLimitOverride[](messages.length);
    for (uint256 i = 0; i < messages.length; ++i) {
      gasLimitOverrides[i].receiverExecutionGasLimit = messages[i].gasLimit;
      gasLimitOverrides[i].tokenGasOverrides = new uint32[](messages[i].tokenAmounts.length);

      for (uint256 j = 0; j < messages[i].tokenAmounts.length; ++j) {
        gasLimitOverrides[i].tokenGasOverrides[j] = DEFAULT_TOKEN_DEST_GAS_OVERHEAD + 1;
      }
    }

    return gasLimitOverrides;
  }
}
