// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CallWithExactGas} from "../../call/CallWithExactGas.sol";

contract CallWithExactGasHelper {
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
}
