// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRouter} from "../v1_0_0/FunctionsRouter.sol";
import {CallWithExactGasZKSync} from "../../shared/call/CallWithExactGasZKSync.sol";

///
/// @title FunctionsRouterZkSync
/// @notice Specialized version of FunctionsRouter for zkSync that uses
/// CallWithExactGasZKSync to control callback gas usage.
///
contract ZKSyncFunctionsRouter is FunctionsRouter {
  constructor(address linkToken, FunctionsRouter.Config memory config) FunctionsRouter(linkToken, config) {}

  /// @dev Override the internal callback function to use CallWithExactGasZKSync
  /// for controlling and measuring gas usage on zkSync.
  function _callback(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint32 callbackGasLimit,
    address client
  ) internal override returns (CallbackResult memory) {
    if (client.code.length == 0) {
      // If there's no code at `client`, skip the callback
      return CallbackResult({success: false, gasUsed: 0, returnData: new bytes(0)});
    }
    uint256 g1 = gasleft();

    (bool success, bytes memory returnData, uint256 pubdataGasSpent) = CallWithExactGasZKSync
      ._callWithExactGasSafeReturnData(
        client,
        callbackGasLimit,
        abi.encodeWithSelector(this.getConfig().handleOracleFulfillmentSelector, requestId, response, err),
        MAX_CALLBACK_RETURN_BYTES
      );
    return CallbackResult({success: success, gasUsed: g1 - gasleft() + pubdataGasSpent, returnData: returnData});
  }
}
