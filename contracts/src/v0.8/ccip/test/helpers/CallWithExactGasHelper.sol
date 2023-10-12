// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {CallWithExactGas} from "../../libraries/CallWithExactGas.sol";

contract CallWithExactGasHelper {
  function callWithExactGas(
    bytes memory payload,
    address target,
    uint256 gasLimit,
    uint16 maxReturnBytes,
    uint16 gasForCallExactCheck
  ) public returns (bool success, bytes memory retData) {
    (success, retData) = CallWithExactGas._callWithExactGas(
      payload,
      target,
      gasLimit,
      maxReturnBytes,
      gasForCallExactCheck
    );
  }
}
