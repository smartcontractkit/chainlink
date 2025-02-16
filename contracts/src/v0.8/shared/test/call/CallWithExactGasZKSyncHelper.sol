// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {CallWithExactGasZKSync} from "../../call/CallWithExactGasZKSync.sol";

/**
 * @notice This helper contract exposes the `_callWithExactGasSafeReturnData` function from the
 * CallWithExactGasZKSync library so it can be called easily in unit tests.
 */
contract CallWithExactGasZKSyncHelper {
  function callWithExactGasSafeReturnData(
    address _to,
    uint256 _maxTotalGas,
    bytes memory _data,
    uint16 _maxReturnBytes
  ) external returns (bool success, bytes memory retData, uint256 pubdataGasSpent) {
    return CallWithExactGasZKSync._callWithExactGasSafeReturnData(_to, _maxTotalGas, _data, _maxReturnBytes);
  }
}
