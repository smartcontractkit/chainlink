// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {CallWithExactGas} from "../../call/CallWithExactGas.sol";

contract CallWithExactGasHelper {
  function callWithExactGas(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 gasForCallExactCheck
  ) public returns (bool success) {
    return CallWithExactGas._callWithExactGas(payload, target, gasLimit, gasForCallExactCheck);
  }

  function callWithExactGasSafeReturnData(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 gasForCallExactCheck,
    uint16 maxReturnBytes
  ) public returns (bool success, bytes memory retData, uint256 gasUsed) {
    return
      CallWithExactGas._callWithExactGasSafeReturnData(payload, target, gasLimit, gasForCallExactCheck, maxReturnBytes);
  }

  function callWithExactGasEvenIfTargetIsNoContract(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 gasForCallExactCheck
  ) public returns (bool success, bool sufficientGas) {
    return CallWithExactGas._callWithExactGasEvenIfTargetIsNoContract(payload, target, gasLimit, gasForCallExactCheck);
  }
}
